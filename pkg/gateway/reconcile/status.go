package reconcile

import (
	"context"
	"fmt"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func GatewayStatus(ctx context.Context, params Params) error {
	gatewayStatus := params.Instance.Status
	gatewayStatus.Host = params.Instance.Spec.App.Management.Cluster.Hostname
	gatewayStatus.Image = params.Instance.Spec.App.Image
	gatewayStatus.Version = params.Instance.Spec.Version
	gatewayStatus.Gateway = []securityv1.GatewayState{}

	dep, err := getGatewayDeployment(ctx, params)
	if err != nil {
		params.Log.V(2).Info("deployment hasn't been created yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	} else {
		gatewayStatus.Replicas = dep.Status.Replicas
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		params.Log.V(2).Info("pods aren't available yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	if len(params.Instance.Spec.App.RepositoryReferences) < len(gatewayStatus.RepositoryStatus) {
		for i, repoStatus := range gatewayStatus.RepositoryStatus {
			found := false
			for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
				if repoStatus.Name == repoRef.Name {
					gatewayStatus.RepositoryStatus[i].Enabled = repoRef.Enabled
					found = true
				}
			}
			if !found {
				gatewayStatus.RepositoryStatus[i].Enabled = false
				gatewayStatus.RepositoryStatus[i].Commit = ""
			}
		}
	}

	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("repository not found", "name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
			return err
		}

		if repository.Status.Commit == "" {
			return fmt.Errorf("repository %s is not ready yet", repository.Name)
		}

		found := false
		for i, repoStatus := range gatewayStatus.RepositoryStatus {
			if repoStatus.Name == repository.Name {
				gatewayStatus.RepositoryStatus[i].Enabled = repoRef.Enabled
				gatewayStatus.RepositoryStatus[i].Commit = repository.Status.Commit
				found = true
			}
		}

		if !found {
			secretName := repository.Name
			if repository.Spec.Auth.ExistingSecretName != "" {
				secretName = repository.Spec.Auth.ExistingSecretName
			}

			if repository.Spec.Auth == (securityv1.RepositoryAuth{}) {
				secretName = ""
			}

			rs := securityv1.GatewayRepositoryStatus{
				Commit:            repository.Status.Commit,
				Enabled:           repoRef.Enabled,
				Name:              repoRef.Name,
				Type:              string(repoRef.Type),
				SecretName:        secretName,
				StorageSecretName: repository.Status.StorageSecretName,
				Endpoint:          repository.Spec.Endpoint,
			}

			if repository.Spec.Tag != "" && repository.Spec.Branch == "" {
				rs.Tag = repository.Spec.Tag
			}

			if repository.Spec.Branch != "" {
				rs.Branch = repository.Spec.Branch
			}

			rs.RemoteName = "origin"
			if repository.Spec.RemoteName != "" {
				rs.RemoteName = repository.Spec.RemoteName
			}

			if repository.Spec.StateStoreReference != "" {
				rs.StateStoreReference = repository.Spec.StateStoreReference
				statestore := &securityv1alpha1.L7StateStore{}
				err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.StateStoreReference, Namespace: params.Instance.Namespace}, statestore)
				if err != nil && k8serrors.IsNotFound(err) {
					params.Log.Info("state store not found", "name", repository.Spec.StateStoreReference, "repository", repository.Name, "namespace", params.Instance.Namespace)
					return err
				}
				rs.StateStoreKey = statestore.Spec.Redis.GroupName + ":" + statestore.Spec.Redis.StoreId + ":" + "repository" + ":" + repository.Status.StorageSecretName + ":latest"
				if repository.Spec.StateStoreKey != "" {
					rs.StateStoreKey = repository.Spec.StateStoreKey
				}
			}
			gatewayStatus.RepositoryStatus = append(gatewayStatus.RepositoryStatus, rs)
		}
	}

	if podList != nil {
		for _, p := range podList.Items {
			if p.ObjectMeta.Labels["management-access"] == "leader" {
				gatewayStatus.ManagementPod = p.Name
			}
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
