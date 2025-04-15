package reconcile

import (
	"context"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
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
