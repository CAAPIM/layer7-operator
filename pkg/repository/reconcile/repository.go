/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
package reconcile

import (
	"context"
	"net/url"
	"reflect"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-git/go-git/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	corev1 "k8s.io/api/core/v1"
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
	var statestore securityv1alpha1.L7StateStore

	authType := repository.Spec.Auth.Type
	if err != nil {
		params.Log.Info("repository unavailable", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		_ = s.RemoveByTag(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-repository")
		return nil
	}

	params.Instance = &repository

	// repoStatus := repository.Status
	if !repository.Spec.Enabled {
		return nil
	}
	start := time.Now()

	if repository.Spec.Auth != (securityv1.RepositoryAuth{}) && repository.Spec.Auth.Type != securityv1.RepositoryAuthTypeNone {

		repositorySecret, err = getRepositorySecret(ctx, repository, params)

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

	if repository.Spec.StateStoreReference != "" {
		statestore, err = getStateStore(ctx, params)
		if err != nil {
			params.Log.V(2).Error(err, "failed to retrieve statestore", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			return nil
		}
	}

	patch := []byte(`[{"op": "replace", "path": "/status/ready", "value": false}]`)
	switch strings.ToLower(string(repository.Spec.Type)) {
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
				err = setRepoReady(ctx, params, patch)
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
			err = setRepoReady(ctx, params, patch)
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
		if params.Instance.Status.Commit == commit {
			params.Log.V(5).Info("already up-to-date", "name", repository.Name, "namespace", repository.Namespace)
			return nil
		}
	case "git":
		commit, err = util.CloneRepository(repository.Spec.Endpoint, username, token, sshKey, sshKeyPass, repository.Spec.Branch, repository.Spec.Tag, repository.Spec.RemoteName, repository.Name, repository.Spec.Auth.Vendor, string(authType), knownHosts, repository.Namespace)

		stateStoreSynced := true
		if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
			params.Log.V(5).Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)

			if repository.Spec.StateStoreReference != "" && !repository.Status.StateStoreSynced {
				err = StateStorage(ctx, params, statestore)
				if err != nil {
					params.Log.V(2).Info("failed to reconcile state storage", "name", repository.Name+"-repository", "namespace", repository.Namespace, "error", err.Error())
					stateStoreSynced = false
				}
			}
			params.Instance.Status.StateStoreSynced = stateStoreSynced

			updateErr := updateStatus(ctx, params, commit, storageSecretName)
			if updateErr != nil {
				_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
				params.Log.Info("failed to update repository status", "namespace", repository.Namespace, "name", repository.Name, "error", err.Error())
			}
			return nil
		}

		if err != nil {
			params.Log.Info("repository error", "name", repository.Name, "namespace", repository.Namespace, "error", err.Error())
			attempts := syncRequest.Attempts + 1
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: attempts}, time.Now().Add(30*time.Second).Unix())
			err = setRepoReady(ctx, params, patch)
			if err != nil {
				params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			}
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			return nil
		}
	case "statestore":
		storageSecretName = "_"
		commit, err = GetStateStoreChecksum(ctx, params, statestore)
		if err != nil {
			params.Log.Info("repository error", "name", repository.Name, "namespace", repository.Namespace, "error", err.Error())
			attempts := syncRequest.Attempts + 1
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: attempts}, time.Now().Add(30*time.Second).Unix())
			err = setRepoReady(ctx, params, patch)
			if err != nil {
				params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			}
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			return nil
		}

		if params.Instance.Status.Commit == commit {
			params.Log.V(5).Info("already up-to-date", "name", repository.Name, "namespace", repository.Namespace)
			return nil
		}

	case "local":
		return nil
	default:
		params.Log.Info("repository type not set or not supported", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	if strings.ToLower(string(repository.Spec.Type)) != "statestore" {
		if repository.Spec.StateStoreReference != "" {
			err = StateStorage(ctx, params, statestore)
			if err != nil {
				params.Log.V(2).Info("failed to reconcile state storage", "name", repository.Name+"-repository", "namespace", repository.Namespace, "error", err.Error())
				params.Instance.Status.StateStoreSynced = false
			}
		} else {
			err = StorageSecret(ctx, params)
			if err != nil {
				// add a check here that prevents the Operator trying to sync the secret repeatedly
				params.Log.V(2).Info("failed to reconcile storage secret", "name", repository.Name+"-repository", "namespace", repository.Namespace, "error", err.Error())
				storageSecretName = ""
				if err.Error() == "exceededMaxSize" {
					storageSecretName = "_"
				}
			}
		}
	}

	err = updateStatus(ctx, params, commit, storageSecretName)
	if err != nil {
		_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
		params.Log.Info("failed to update repository status", "namespace", repository.Namespace, "name", repository.Name, "error", err.Error())
		return nil
	}

	_ = captureRepositorySyncMetrics(ctx, params, start, commit, false)
	params.Log.Info("reconciled", "name", repository.Name, "namespace", repository.Namespace, "commit", commit)
	return nil
}

func getRepositorySecret(ctx context.Context, repository securityv1.Repository, params Params) (*corev1.Secret, error) {
	repositorySecret := &corev1.Secret{}
	name := repository.Name

	if repository.Spec.Auth.ExistingSecretName != "" {
		name = repository.Spec.Auth.ExistingSecretName
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: repository.Namespace}, repositorySecret)
	if err != nil {
		return repositorySecret, err
	}
	return repositorySecret, nil
}

func getRepository(ctx context.Context, params Params) (securityv1.Repository, error) {
	repository := securityv1.Repository{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, &repository)
	if err != nil {
		return repository, err
	}
	return repository, nil
}

func updateStatus(ctx context.Context, params Params, commit string, storageSecretName string) (err error) {

	if params.Instance.Status.StorageSecretName == "_" {
		storageSecretName = "_"
	}

	rs := params.Instance.Status
	r := params.Instance

	rs.Commit = commit
	rs.Name = r.Name
	rs.Vendor = r.Spec.Auth.Vendor
	rs.Ready = true

	if r.Spec.Type == "http" {
		rs.Summary = r.Spec.Endpoint
	}

	rs.StorageSecretName = storageSecretName

	if !reflect.DeepEqual(rs, r.Status) {
		params.Log.Info("syncing repository", "name", r.Name, "namespace", r.Namespace)
		rs.Updated = time.Now().String()

		if rs.StateStoreSynced {
			if r.Spec.StateStoreReference != "" {
				if rs.StateStoreVersion == 0 {
					rs.StateStoreVersion = 1
				} else {
					rs.StateStoreVersion = rs.StateStoreVersion + 1
				}
			}
		}

		r.Status = rs
		err = params.Client.Status().Update(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func setRepoReady(ctx context.Context, params Params, patch []byte) error {

	if !params.Instance.Status.Ready {
		return nil
	}

	if err := params.Client.Status().Patch(ctx, params.Instance,
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
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("repository_namespace", gateway.Namespace)))

	repoSyncTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("repository_namespace", gateway.Namespace)))

	if hasError {
		repoSyncFailure.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("repository_namespace", gateway.Namespace),
				attribute.String("repository_type", string(params.Instance.Spec.Type)),
				attribute.String("repository_name", params.Instance.Name),
				attribute.String("commit_id", commitId)))
	} else {
		repoSyncSuccess.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("repository_namespace", gateway.Namespace),
				attribute.String("repository_type", string(params.Instance.Spec.Type)),
				attribute.String("repository_name", params.Instance.Name),
				attribute.String("commit_id", commitId)))
	}

	return nil
}
