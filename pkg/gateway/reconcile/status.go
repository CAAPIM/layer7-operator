package reconcile

import (
	"context"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
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
		return err
	}

	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("repository not found", "name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
			return err
		}

		for i, repoStatus := range gatewayStatus.RepositoryStatus {
			if repoStatus.Name == repository.Name {
				gatewayStatus.RepositoryStatus[i].Enabled = repoRef.Enabled
			}
		}
	}

	for _, p := range podList.Items {
		if p.ObjectMeta.Labels["management-access"] == "leader" {
			gatewayStatus.ManagementPod = p.Name
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
