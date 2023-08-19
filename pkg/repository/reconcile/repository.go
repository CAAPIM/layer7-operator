package reconcile

import (
	"context"
	"reflect"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-git/go-git/v5"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncRepository(ctx context.Context, params Params) error {
	// TODO: Deregister job if repository is removed..
	repository, err := getRepository(ctx, params)
	if err != nil {
		params.Log.Info("repository unavailable", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		return nil
	}

	params.Instance = &repository

	ext := repository.Spec.Branch
	if ext == "" {
		ext = repository.Spec.Tag
	}

	repoStatus := repository.Status
	if !repository.Spec.Enabled {
		params.Log.Info("repository not enabled", "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	rSecret, err := getSecret(ctx, repository, params)

	if err != nil {
		params.Log.Info("secret unavailable", "name", repository.Name, "namespace", repository.Namespace, "error", err.Error())
		return nil
	}

	token := string(rSecret.Data["TOKEN"])
	if token == "" {
		token = string(rSecret.Data["PASSWORD"])
	}

	username := string(rSecret.Data["USERNAME"])

	commit, err := util.CloneRepository(repository.Spec.Endpoint, username, token, repository.Spec.Branch, repository.Spec.Tag, repository.Spec.RemoteName, repository.Spec.Name, repository.Spec.Auth.Vendor)

	if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
		params.Log.V(2).Info(err.Error(), "name", repository.Name, "namespace", repository.Namespace)
		return nil
	}

	if err != nil {
		params.Log.Error(err, "repository err", "name", repository.Name, "namespace", repository.Namespace)
		return err
	}

	storageSecretName := repository.Name + "-repository-" + ext

	err = StorageSecret(ctx, params)
	if err != nil {
		params.Log.V(2).Info("failed to reconcile storage secret", "name", repository.Name+"-repository", "namespace", repository.Namespace, "error", err.Error())
		storageSecretName = ""
	}

	repoStatus.Commit = commit
	repoStatus.Name = repository.Name
	repoStatus.Vendor = repository.Spec.Auth.Vendor

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
