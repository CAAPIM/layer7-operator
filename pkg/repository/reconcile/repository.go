package reconcile

import (
	"context"
	"net/url"
	"reflect"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-git/go-git/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func syncRepository(ctx context.Context, params Params) error {
	repository, err := getRepository(ctx, params)
	var commit string
	var username string
	var token string
	var repositorySecret *corev1.Secret
	var sshKey []byte
	var sshKeyPass string
	var knownHosts []byte

	authType := repository.Spec.Auth.Type
	if err != nil {
		params.Log.Info("repository unavailable", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		_ = s.RemoveByTag(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-repository")
		return nil
	}

	params.Instance = &repository

	repoStatus := repository.Status
	if !repository.Spec.Enabled {
		return nil
	}
	start := time.Now()

	if repository.Spec.Auth != (securityv1.RepositoryAuth{}) && repository.Spec.Auth.Type != securityv1.RepositoryAuthTypeNone {

		repositorySecret, err = getSecret(ctx, repository, params)

		if err != nil {
			params.Log.Info("secret unavailable", "name", repository.Name, "namespace", repository.Namespace, "error", err.Error())
			return nil
		}

		token = string(repositorySecret.Data["TOKEN"])
		if token == "" {
			token = string(repositorySecret.Data["PASSWORD"])
		}

		username = string(repositorySecret.Data["USERNAME"])
		sshKey = repositorySecret.Data["SSH_KEY"]
		sshKeyPass = string(repositorySecret.Data["SSH_KEY_PASS"])
		knownHosts = repositorySecret.Data["KNOWN_HOSTS"]
		authType = repository.Spec.Auth.Type

		if authType == "" && username != "" && token != "" {
			authType = securityv1.RepositoryAuthTypeBasic
		}

		if authType == "" && username == "" && sshKey != nil {
			authType = securityv1.RepositoryAuthTypeSSH
		}

	}

	ext := repository.Spec.Branch
	if ext == "" {
		ext = repository.Spec.Tag
	}

	storageSecretName := repository.Name + "-repository-" + ext

	requestCacheEntry := repository.Name
	syncRequest, _ := syncCache.Read(requestCacheEntry)

	backoffRequestCacheEntry := repository.Name + "-backoff"
	backoffSyncRequest, _ := syncCache.Read(backoffRequestCacheEntry)

	if backoffSyncRequest.Attempts > 4 {
		params.Log.Info("several attempts to sync this repository have failed, please check your configuration", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	if syncRequest.Attempts > 2 {
		params.Log.V(2).Info("request has failed more than 3 times in the last 30 seconds", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	patch := []byte(`[{"op": "replace", "path": "/status/ready", "value": false}]`)

	switch strings.ToLower(repository.Spec.Type) {
	case "http":
		forceUpdate := false
		if repository.Status.Summary != repository.Spec.Endpoint {
			forceUpdate = true
		}
		commit, err = util.DownloadArtifact(repository.Spec.Endpoint, username, token, repository.Name, forceUpdate, repository.Namespace)
		if err != nil {
			if err == util.ErrInvalidFileFormatError || err == util.ErrInvalidTarArchive || err == util.ErrInvalidZipArchive {
				params.Log.Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)
				backoffAttempts := 5
				syncCache.Update(util.SyncRequest{RequestName: backoffRequestCacheEntry, Attempts: backoffAttempts}, time.Now().Add(360*time.Second).Unix())
				err = setRepoStatus(ctx, params, patch)
				if err != nil {
					params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
				}
				_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
				return nil
			}
			params.Log.Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)
			attempts := syncRequest.Attempts + 1
			backoffAttempts := backoffSyncRequest.Attempts + 1
			syncCache.Update(util.SyncRequest{RequestName: backoffRequestCacheEntry, Attempts: backoffAttempts}, time.Now().Add(360*time.Second).Unix())
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: attempts}, time.Now().Add(30*time.Second).Unix())
			err = setRepoStatus(ctx, params, patch)
			if err != nil {
				params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			}
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			return nil
		}
		fileURL, _ := url.Parse(repository.Spec.Endpoint)
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]
		ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		if ext == "gz" && strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
			folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
		}
		storageSecretName = repository.Name + "-repository-" + folderName

	case "git":
		commit, err = util.CloneRepository(repository.Spec.Endpoint, username, token, sshKey, sshKeyPass, repository.Spec.Branch, repository.Spec.Tag, repository.Spec.RemoteName, repository.Name, repository.Spec.Auth.Vendor, string(authType), knownHosts, repository.Namespace)
		if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
			params.Log.V(2).Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)
			return nil
		}

		if err != nil {
			params.Log.Info("repository error", "name", repository.Name, "namespace", repository.Namespace, "error", err.Error())
			attempts := syncRequest.Attempts + 1
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: attempts}, time.Now().Add(30*time.Second).Unix())
			err = setRepoStatus(ctx, params, patch)
			if err != nil {
				params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			}
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			return nil
		}
	default:
		params.Log.Info("repository type not set", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	err = StorageSecret(ctx, params)
	if err != nil {
		params.Log.V(2).Info("failed to reconcile storage secret", "name", repository.Name+"-repository", "namespace", repository.Namespace, "error", err.Error())
		storageSecretName = ""
	}

	repoStatus.Commit = commit
	repoStatus.Name = repository.Name
	repoStatus.Vendor = repository.Spec.Auth.Vendor
	repoStatus.Ready = true

	if repository.Spec.Type == "http" {
		// future usage will include filesize for tracking changes in remote
		// or use a different status field.
		repoStatus.Summary = repository.Spec.Endpoint
	}

	repoStatus.StorageSecretName = storageSecretName

	if !reflect.DeepEqual(repoStatus, repository.Status) {
		params.Log.Info("syncing repository", "name", repository.Name, "namespace", repository.Namespace)
		repoStatus.Updated = time.Now().String()
		repository.Status = repoStatus
		err = params.Client.Status().Update(ctx, &repository)
		if err != nil {
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			params.Log.Info("failed to update repository status", "namespace", repository.Namespace, "name", repository.Name, "error", err.Error())
		}
		params.Log.Info("reconciled", "name", repository.Name, "namespace", repository.Namespace, "commit", commit)
	}
	_ = captureRepositorySyncMetrics(ctx, params, start, commit, false)
	return nil
}

func getSecret(ctx context.Context, repository securityv1.Repository, params Params) (*corev1.Secret, error) {
	repositorySecret := &corev1.Secret{}
	name := repository.Name

	if repository.Spec.Auth.ExistingSecretName != "" {
		name = repository.Spec.Auth.ExistingSecretName
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: repository.Namespace}, repositorySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return repositorySecret, err
			}
		}
	}
	return repositorySecret, nil
}

func getRepository(ctx context.Context, params Params) (securityv1.Repository, error) {
	repository := securityv1.Repository{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, &repository)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return repository, err
			}
		}
	}
	return repository, nil
}

func setRepoStatus(ctx context.Context, params Params, patch []byte) error {

	if !params.Instance.Status.Ready {
		return nil
	}

	if err := params.Client.Status().Patch(context.Background(), params.Instance,
		client.RawPatch(types.JSONPatchType, patch)); err != nil {
		return err
	}
	params.Log.Info("repository status has been updated", "namespace", params.Instance.Namespace, "name", params.Instance.Name)

	return nil
}

func captureRepositorySyncMetrics(ctx context.Context, params Params, start time.Time, commitId string, hasError bool) error {
	operatorNamespace, err := util.GetOperatorNamespace()
	if err != nil {
		params.Log.Info("could not determine operator namespace")
		return err
	}
	gateway := params.Instance
	otelEnabled, err := util.GetOtelEnabled()
	if err != nil {
		params.Log.Info("could not determine if OTel is enabled")
		return err
	}

	if !otelEnabled {
		return nil
	}

	otelMetricPrefix, err := util.GetOtelMetricPrefix()
	if err != nil {
		params.Log.Info("could not determine otel metric prefix")
		return err
	}

	if otelMetricPrefix == "" {
		otelMetricPrefix = "layer7_"
	}

	hostname, err := util.GetHostname()
	if err != nil {
		params.Log.Error(err, "failed to retrieve operator hostname")
		return err
	}
	if err != nil {
		params.Log.Error(err, "failed to retrieve operator namespace")
		return err
	}

	meter := otel.Meter("layer7-operator-repository-sync-metrics")
	repoSyncLatency, err := meter.Float64Histogram(otelMetricPrefix+"operator_repository_sync_latency",
		metric.WithDescription("repository sync latency"), metric.WithUnit("ms"))
	if err != nil {
		return err
	}

	repoSyncSuccess, err := meter.Int64Counter(otelMetricPrefix+"operator_repository_sync_success",
		metric.WithDescription("graphman request success"))
	if err != nil {
		return err
	}

	repoSyncFailure, err := meter.Int64Counter(otelMetricPrefix+"operator_repository_sync_failure",
		metric.WithDescription("graphman request failure"))
	if err != nil {
		return err
	}

	repoSyncTotal, err := meter.Int64Counter(otelMetricPrefix+"operator_repository_sync_total",
		metric.WithDescription("graphman request total"))
	if err != nil {
		return err
	}

	duration := time.Since(start)
	repoSyncLatency.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("repository_namespace", gateway.Namespace)))

	repoSyncTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("repository_namespace", gateway.Namespace)))

	if hasError {
		repoSyncFailure.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("pod", hostname),
				attribute.String("namespace", operatorNamespace),
				attribute.String("repository_namespace", gateway.Namespace),
				attribute.String("repository_type", params.Instance.Spec.Type),
				attribute.String("repository_name", params.Instance.Name),
				attribute.String("commit_id", commitId)))
	} else {
		repoSyncSuccess.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("pod", hostname),
				attribute.String("namespace", operatorNamespace),
				attribute.String("repository_namespace", gateway.Namespace),
				attribute.String("repository_type", params.Instance.Spec.Type),
				attribute.String("repository_name", params.Instance.Name),
				attribute.String("commit_id", commitId)))
	}

	return nil
}
