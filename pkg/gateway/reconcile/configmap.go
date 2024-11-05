package reconcile

import (
	"context"
	"fmt"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ConfigMaps(ctx context.Context, params Params) error {
	desiredConfigMaps := []*corev1.ConfigMap{
		gateway.NewConfigMap(params.Instance, params.Instance.Name),
		gateway.NewConfigMap(params.Instance, params.Instance.Name+"-system"),
		gateway.NewConfigMap(params.Instance, params.Instance.Name+"-repository-init-config"),
		gateway.NewConfigMap(params.Instance, params.Instance.Name+"-gateway-files"),
	}

	if params.Instance.Spec.App.ClusterProperties.Enabled {
		desiredConfigMaps = append(desiredConfigMaps, gateway.NewConfigMap(params.Instance, params.Instance.Name+"-cwp-bundle"))
	}

	if params.Instance.Spec.App.ListenPorts.Harden || params.Instance.Spec.App.ListenPorts.Custom.Enabled {
		desiredConfigMaps = append(desiredConfigMaps, gateway.NewConfigMap(params.Instance, params.Instance.Name+"-listen-port-bundle"))
	}

	if params.Instance.Spec.App.Otk.Enabled {
		desiredConfigMaps = append(desiredConfigMaps, gateway.NewConfigMap(params.Instance, params.Instance.Name+"-otk-install-init-config"))
		desiredConfigMaps = append(desiredConfigMaps, gateway.NewConfigMap(params.Instance, params.Instance.Name+"-otk-shared-init-config"))
	}

	if params.Instance.Spec.App.Otk.Database.Type == securityv1.OtkDatabaseTypeMySQL || params.Instance.Spec.App.Otk.Database.Type == securityv1.OtkDatabaseTypeOracle {
		if params.Instance.Spec.App.Otk.Database.Sql.ManageSchema {
			desiredConfigMaps = append(desiredConfigMaps, gateway.NewConfigMap(params.Instance, params.Instance.Name+"-otk-db-init-config"))
		}
	}

	if err := reconcileConfigMaps(ctx, params, desiredConfigMaps); err != nil {
		return fmt.Errorf("failed to reconcile configMaps: %w", err)
	}

	return nil
}

func reconcileConfigMaps(ctx context.Context, params Params, desiredConfigMaps []*corev1.ConfigMap) error {
	for _, dcm := range desiredConfigMaps {
		desiredConfigMap := dcm
		if err := controllerutil.SetControllerReference(params.Instance, desiredConfigMap, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		currentConfigMap := corev1.ConfigMap{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: desiredConfigMap.Name, Namespace: params.Instance.Namespace}, &currentConfigMap)
		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredConfigMap); err != nil {
				return err
			}
			params.Log.Info("created configMap", "name", desiredConfigMap.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}

		if desiredConfigMap.ObjectMeta.Annotations["checksum/data"] != currentConfigMap.ObjectMeta.Annotations["checksum/data"] {
			patch := client.MergeFrom(&currentConfigMap)
			if err := params.Client.Patch(ctx, desiredConfigMap, patch); err != nil {
				return err
			}
			params.Log.V(2).Info("configMap updated", "name", desiredConfigMap.Name, "namespace", desiredConfigMap.Namespace)
		}
	}

	return nil
}
