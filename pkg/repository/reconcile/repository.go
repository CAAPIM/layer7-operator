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
	//authType := "none"
	authType := securityv1.RepositoryAuthTypeNone
	if err != nil {
		params.Log.Info("repository unavailable", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		_ = s.RemoveByTag(params.Instance.Name + "-sync-repository")
		return nil
	}

	params.Instance = &repository

	repoStatus := repository.Status
	if !repository.Spec.Enabled {
		params.Log.Info("repository not enabled", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	if repository.Spec.Auth != (securityv1.RepositoryAuth{}) {

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
		commit, err = util.DownloadArtifact(repository.Spec.Endpoint, username, token, repository.Spec.Name, forceUpdate)
		if err != nil {
			if err == util.ErrInvalidFileFormatError || err == util.ErrInvalidTarArchive || err == util.ErrInvalidZipArchive {
				params.Log.Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)
				backoffAttempts := 5
				syncCache.Update(util.SyncRequest{RequestName: backoffRequestCacheEntry, Attempts: backoffAttempts}, time.Now().Add(360*time.Second).Unix())
				err = setRepoStatus(ctx, params, patch)
				if err != nil {
					params.Log.V(2).Error(err, "failed to patch repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
				}
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
	default:
		commit, err = util.CloneRepository(repository.Spec.Endpoint, username, token, sshKey, sshKeyPass, repository.Spec.Branch, repository.Spec.Tag, repository.Spec.RemoteName, repository.Spec.Name, repository.Spec.Auth.Vendor, string(authType), knownHosts)
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
			return nil
		}
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
			params.Log.Info("failed to update repository status", "namespace", repository.Namespace, "name", repository.Name, "error", err.Error())
		}
		params.Log.Info("reconciled", "name", repository.Name, "namespace", repository.Namespace, "commit", commit)
	}

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
