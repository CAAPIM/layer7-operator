package reconcile

import (
	"context"
)

func ClusterProperties(ctx context.Context, params Params) error {
	gateway := params.Instance
	if !gateway.Spec.App.ClusterProperties.Enabled {
		return nil
	}
	name := gateway.Name
	if gateway.Spec.App.Management.DisklessConfig.Disabled {
		name = gateway.Name + "-node-properties"
	}
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	if gateway.Spec.App.ClusterProperties.Enabled {
		cm, err := getGatewayConfigMap(ctx, params, params.Instance.Name+"-cwp-bundle")
		if err != nil {
			return err
		}

		annotation := "security.brcmlabs.com/" + params.Instance.Name + "-cwp-bundle"
		if !gateway.Spec.App.Management.Database.Enabled {
			podList, err := getGatewayPods(ctx, params)
			if err != nil {
				return err
			}
			err = ReconcileEphemeralGateway(ctx, params, "cluster properties", *podList, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "cluster properties", []byte(cm.Data["cwp.json"]))
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "cluster properties", gatewayDeployment, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "cluster properties", []byte(cm.Data["cwp.json"]))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ListenPorts(ctx context.Context, params Params) error {
	gateway := params.Instance
	if !gateway.Spec.App.ClusterProperties.Enabled {
		return nil
	}
	name := gateway.Name
	if gateway.Spec.App.Management.DisklessConfig.Disabled {
		name = gateway.Name + "-node-properties"
	}
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	if gateway.Spec.App.ListenPorts.Harden || gateway.Spec.App.ListenPorts.Custom.Enabled {
		cm, err := getGatewayConfigMap(ctx, params, params.Instance.Name+"-listen-port-bundle")
		if err != nil {
			return err
		}

		annotation := "security.brcmlabs.com/" + params.Instance.Name + "-listen-port-bundle"
		if !gateway.Spec.App.Management.Database.Enabled {
			podList, err := getGatewayPods(ctx, params)
			if err != nil {
				return err
			}
			err = ReconcileEphemeralGateway(ctx, params, "listen ports", *podList, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "listen ports", []byte(cm.Data["listen-ports.json"]))
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "listen ports", gatewayDeployment, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "listen ports", []byte(cm.Data["listen-ports.json"]))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
