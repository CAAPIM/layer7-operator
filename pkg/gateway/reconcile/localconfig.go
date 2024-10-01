package reconcile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
)

func ClusterProperties(ctx context.Context, params Params) error {
	cleanUpDbbacked := false
	gateway := params.Instance
	if !gateway.Spec.App.ClusterProperties.Enabled {
		if len(gateway.Status.LastAppliedClusterProperties) == 0 {
			return nil
		}
		if !gateway.Spec.App.Management.Database.Enabled {
			gateway.Status.LastAppliedClusterProperties = []string{}
			if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
				return fmt.Errorf("failed to remove cluster properties status: %w", err)
			}
			return nil
		}
		cleanUpDbbacked = true
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

	cm, err := getGatewayConfigMap(ctx, params, params.Instance.Name+"-cwp-bundle")
	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/" + params.Instance.Name + "-cwp-bundle"

	bundle := graphman.Bundle{}
	err = json.Unmarshal([]byte(cm.Data["cwp.json"]), &bundle)
	if err != nil {
		return err
	}

	notFound := []string{}

	if !cleanUpDbbacked {
		for _, sCwp := range params.Instance.Status.LastAppliedClusterProperties {
			found := false
			for _, cwp := range bundle.ClusterProperties {
				if cwp.Name == sCwp {
					found = true
				}
			}
			if !found {
				notFound = append(notFound, sCwp)
			}
		}
	} else {
		notFound = append(notFound, params.Instance.Status.LastAppliedClusterProperties...)
	}

	bundle.Properties = &graphman.BundleProperties{}
	for _, deletedCwp := range notFound {
		mappingSource := MappingSource{Name: deletedCwp}
		bundle.ClusterProperties = append(bundle.ClusterProperties, &graphman.ClusterPropertyInput{
			Name:  deletedCwp,
			Value: "to be deleted",
		})

		bundle.Properties.Mappings.ClusterProperties = append(bundle.Properties.Mappings.ClusterProperties, &graphman.MappingInstructionInput{
			Action: graphman.MappingActionDelete,
			Source: mappingSource,
		})
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "cluster properties", *podList, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "cluster properties", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		err = ReconcileDBGateway(ctx, params, "cluster properties", gatewayDeployment, gateway, gwSecret, "", annotation, cm.ObjectMeta.Annotations["checksum/data"], false, "cluster properties", bundleBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListenPorts(ctx context.Context, params Params) error {
	cleanUpDbbacked := false
	gateway := params.Instance
	if !gateway.Spec.App.ListenPorts.Custom.Enabled {
		if len(gateway.Status.LastAppliedListenPorts) == 0 {
			return nil
		}
		if !gateway.Spec.App.Management.Database.Enabled {
			gateway.Status.LastAppliedListenPorts = []string{}
			if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
				return fmt.Errorf("failed to remove listen ports status: %w", err)
			}
			return nil
		}
		cleanUpDbbacked = true
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

	if gateway.Spec.App.ListenPorts.Custom.Enabled {
		annotation := "security.brcmlabs.com/" + params.Instance.Name + "-listen-port-bundle"

		refreshOnKeyChanges := false
		dataCheckSum := ""
		var bundleBytes []byte
		if params.Instance.Spec.App.ListenPorts.RefreshOnKeyChanges {
			refreshOnKeyChanges = true
		}
		if !params.Instance.Spec.App.ListenPorts.Custom.Enabled {
			bundleBytes, dataCheckSum, err = util.BuildDefaultListenPortBundle(refreshOnKeyChanges)
			if err != nil {
				return err
			}
		} else {
			bundleBytes, dataCheckSum, err = util.BuildCustomListenPortBundle(params.Instance, refreshOnKeyChanges)
			if err != nil {
				return err
			}
		}

		bundle := graphman.Bundle{}
		err = json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}

		grapmanDynamicSyncPort := 9443
		if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
			grapmanDynamicSyncPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
		}

		notFound := []string{}
		if !cleanUpDbbacked {
			for _, slistenPort := range params.Instance.Status.LastAppliedListenPorts {
				found := false
				for _, listenPort := range bundle.ListenPorts {
					if listenPort.Name == slistenPort {
						found = true
					}
					// anti-lockout
					if listenPort.Name == slistenPort && listenPort.Port == grapmanDynamicSyncPort {
						found = true
					}
				}
				if !found {
					notFound = append(notFound, slistenPort)
				}
			}
		} else {
			notFound = append(notFound, params.Instance.Status.LastAppliedListenPorts...)
		}

		bundle.Properties = &graphman.BundleProperties{}
		for _, deletedListenPort := range notFound {
			mappingSource := MappingSource{Name: deletedListenPort}
			bundle.ListenPorts = append(bundle.ListenPorts, &graphman.ListenPortInput{
				Name:     deletedListenPort,
				Port:     1,
				Enabled:  false,
				Protocol: "HTTP",
				EnabledFeatures: []graphman.ListenPortFeature{
					graphman.ListenPortFeaturePublishedServiceMessageInput,
				},
			})

			bundle.Properties.Mappings.ListenPorts = append(bundle.Properties.Mappings.ListenPorts, &graphman.MappingInstructionInput{
				Action: graphman.MappingActionDelete,
				Source: mappingSource,
			})
		}

		bundleBytes, err := json.Marshal(bundle)
		if err != nil {
			return err
		}

		if !gateway.Spec.App.Management.Database.Enabled {
			podList, err := getGatewayPods(ctx, params)
			if err != nil {
				return err
			}
			err = ReconcileEphemeralGateway(ctx, params, "listen ports", *podList, gateway, gwSecret, "", annotation, dataCheckSum, false, "listen ports", bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "listen ports", gatewayDeployment, gateway, gwSecret, "", annotation, dataCheckSum, false, "listen ports", bundleBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
