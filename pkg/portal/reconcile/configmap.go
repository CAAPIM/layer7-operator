package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/portal"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ConfigMap(ctx context.Context, params Params, apiSummary []byte) error {
	desiredConfigMap := portal.NewConfigMap(params.Instance, apiSummary)

	err := reconcileConfigMap(ctx, params, desiredConfigMap)
	if err != nil {
		return fmt.Errorf("failed to reconcile configMaps: %w", err)
	}

	return nil
}

func reconcileConfigMap(ctx context.Context, params Params, desiredConfigMap *corev1.ConfigMap) error {
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
		return nil
	}
	if err != nil {
		return err
	}

	if desiredConfigMap.ObjectMeta.Annotations["checksum/data"] != currentConfigMap.ObjectMeta.Annotations["checksum/data"] {
		patch := client.MergeFrom(&currentConfigMap)
		if err := params.Client.Patch(ctx, desiredConfigMap, patch); err != nil {
			return err
		}
		params.Log.Info("configMap updated", "name", desiredConfigMap.Name, "namespace", desiredConfigMap.Namespace)
		return nil
	}

	return nil
}
