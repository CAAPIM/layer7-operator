package reconcile

import (
	"context"
	"reflect"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncRepository(ctx context.Context, params Params) error {

	repoStatus := params.Instance.Status
	if !params.Instance.Spec.Enabled {
		params.Log.Info("repository not enabled", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	rSecret, err := getSecret(ctx, params)

	if err != nil {
		params.Log.Info("secret unavailable", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		return nil
	}

	token := string(rSecret.Data["TOKEN"])
	if token == "" {
		token = string(rSecret.Data["PASSWORD"])
	}

	commit, err := util.CloneRepository(params.Instance.Spec.Endpoint, string(rSecret.Data["USERNAME"]), token, params.Instance.Spec.Branch, params.Instance.Spec.Name, params.Instance.Spec.Auth.Vendor)

	if err != nil {
		return err
	}

	storageSecretName := params.Instance.Name + "-repository"

	err = StorageSecret(ctx, params)
	if err != nil {
		params.Log.V(2).Info("failed to reconcile storage secret", "name", params.Instance.Name+"-repository", "namespace", params.Instance.Namespace, "error", err.Error())
		storageSecretName = ""
	}

	repoStatus.Commit = commit
	repoStatus.Name = params.Instance.Name
	repoStatus.Vendor = params.Instance.Spec.Auth.Vendor

	repoStatus.StorageSecretName = storageSecretName

	if !reflect.DeepEqual(repoStatus, params.Instance.Status) {
		params.Log.Info("syncing repository", "name", params.Instance.Name, "namespace", params.Instance.Namespace)

		repoStatus.Updated = time.Now().String()
		params.Instance.Status = repoStatus
		err = params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.Info("failed to update repository status", "namespace", params.Instance.Namespace, "name", params.Instance.Name, "error", err.Error())
		}
		params.Log.Info("reconciled", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "commit", commit)
	}

	return nil
}

func getSecret(ctx context.Context, params Params) (*corev1.Secret, error) {
	repositorySecret := &corev1.Secret{}
	name := params.Instance.Name
	if params.Instance.Spec.Auth.ExistingSecretName != "" {
		name = params.Instance.Spec.Auth.ExistingSecretName
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, repositorySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return repositorySecret, err
			}
		}
	}
	return repositorySecret, nil
}
