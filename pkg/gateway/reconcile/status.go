package reconcile

import (
	"context"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func GatewayStatus(ctx context.Context, params Params) error {
	gatewayStatus := params.Instance.Status
	gatewayStatus.RepositoryStatus = []securityv1.GatewayRepositoryStatus{}
	gatewayStatus.Host = params.Instance.Spec.App.Management.Cluster.Hostname
	gatewayStatus.Image = params.Instance.Spec.App.Image
	gatewayStatus.Version = params.Instance.Spec.Version
	gatewayStatus.Gateway = []securityv1.GatewayState{}

	dep, err := getGatewayDeployment(ctx, params)
	if err != nil || k8serrors.IsNotFound(err) {
		params.Log.V(2).Info("deployment hasn't been created yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	} else {
		gatewayStatus.Replicas = dep.Status.Replicas
	}

	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("repository not found", "name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
			return err
		}

		secretName := repository.Name
		if repository.Spec.Auth.ExistingSecretName != "" {
			secretName = repository.Spec.Auth.ExistingSecretName
		}

		if repository.Spec.Auth == (securityv1.RepositoryAuth{}) {
			secretName = ""
		}

		commit := repository.Status.Commit

		newRepoStatus := securityv1.GatewayRepositoryStatus{
			Commit:            commit,
			Enabled:           repoRef.Enabled,
			Name:              repoRef.Name,
			Type:              string(repoRef.Type),
			SecretName:        secretName,
			StorageSecretName: repository.Status.StorageSecretName,
			Endpoint:          repository.Spec.Endpoint,
		}

		if repository.Spec.Tag != "" && repository.Spec.Branch == "" {
			newRepoStatus.Tag = repository.Spec.Tag
		}

		if repository.Spec.Branch != "" {
			newRepoStatus.Branch = repository.Spec.Branch
		}

		newRepoStatus.RemoteName = "origin"
		if repository.Spec.RemoteName != "" {
			newRepoStatus.RemoteName = repository.Spec.RemoteName
		}

		if repository.Spec.StateStoreReference != "" {
			newRepoStatus.StateStoreReference = repository.Spec.StateStoreReference

			statestore := &securityv1alpha1.L7StateStore{}

			err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.StateStoreReference, Namespace: params.Instance.Namespace}, statestore)
			if err != nil && k8serrors.IsNotFound(err) {
				params.Log.Info("state store not found", "name", repository.Spec.StateStoreReference, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
				return err
			}

			newRepoStatus.StateStoreKey = statestore.Spec.Redis.GroupName + ":" + statestore.Spec.Redis.StoreId + ":" + "repository" + ":" + repository.Status.StorageSecretName + ":latest"
			if repository.Spec.StateStoreKey != "" {
				newRepoStatus.StateStoreKey = repository.Spec.StateStoreKey
			}
		}

		gatewayStatus.RepositoryStatus = append(gatewayStatus.RepositoryStatus, newRepoStatus)

	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	for _, p := range podList.Items {
		if p.ObjectMeta.Labels["management-access"] == "leader" {
			gatewayStatus.ManagementPod = p.Name
			// range over repository status and check if the delete annotation has been added to pods
		}
	}

	if !reflect.DeepEqual(gatewayStatus, params.Instance.Status) {
		params.Instance.Status = gatewayStatus
		err = params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.V(2).Info("failed to update gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
		params.Log.V(2).Info("updated gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace)

	}
	return nil
}
