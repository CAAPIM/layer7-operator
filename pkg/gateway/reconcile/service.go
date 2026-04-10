package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Services(ctx context.Context, params Params) error {
	managementServiceName := params.Instance.Name + "-management-service"
	if !params.Instance.Spec.App.Management.Service.Enabled {
		currentMgmt := &corev1.Service{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: managementServiceName, Namespace: params.Instance.Namespace}, currentMgmt)
		if err == nil && controllerutil.HasControllerReference(currentMgmt) {
			if err := params.Client.Delete(ctx, currentMgmt); err != nil {
				return fmt.Errorf("failed to remove management service: %w", err)
			}
			params.Log.Info("removed management service", "name", managementServiceName, "namespace", params.Instance.Namespace)
		} else if err != nil && !k8serrors.IsNotFound(err) {
			return err
		}
	}

	desiredServices := []*corev1.Service{
		gateway.NewService(params.Instance),
	}

	if params.Instance.Spec.App.Management.Service.Enabled {
		desiredServices = append(desiredServices, gateway.NewManagementService(params.Instance))
	}

	if err := reconcileServices(ctx, params, desiredServices); err != nil {
		return fmt.Errorf("failed to reconcile services: %w", err)
	}

	return nil
}

// metadataKeysMatch returns true if current has every label and annotation key/value from desired
// (extra keys on current are allowed).
func metadataKeysMatch(current, desired metav1.ObjectMeta) bool {
	for k, v := range desired.Labels {
		if current.Labels == nil || current.Labels[k] != v {
			return false
		}
	}
	for k, v := range desired.Annotations {
		if current.Annotations == nil || current.Annotations[k] != v {
			return false
		}
	}
	return true
}

func reconcileServices(ctx context.Context, params Params, desiredServices []*corev1.Service) error {
	for _, ds := range desiredServices {
		desiredService := ds
		if err := controllerutil.SetControllerReference(params.Instance, desiredService, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		currentService := corev1.Service{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: desiredService.Name, Namespace: params.Instance.Namespace}, &currentService)
		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredService); err != nil {
				return err
			}
			params.Log.Info("created service", "name", desiredService.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}

		updated := currentService.DeepCopy()

		updated.Spec = desiredService.Spec
		updated.ObjectMeta.OwnerReferences = desiredService.ObjectMeta.OwnerReferences

		if updated.ObjectMeta.Labels == nil {
			updated.ObjectMeta.Labels = make(map[string]string)
		}
		if updated.ObjectMeta.Annotations == nil {
			updated.ObjectMeta.Annotations = make(map[string]string)
		}
		for k, v := range desiredService.ObjectMeta.Annotations {
			updated.ObjectMeta.Annotations[k] = v
		}
		for k, v := range desiredService.ObjectMeta.Labels {
			updated.ObjectMeta.Labels[k] = v
		}

		if apiequality.Semantic.DeepEqual(currentService.Spec, desiredService.Spec) &&
			apiequality.Semantic.DeepEqual(currentService.ObjectMeta.OwnerReferences, desiredService.ObjectMeta.OwnerReferences) &&
			metadataKeysMatch(currentService.ObjectMeta, desiredService.ObjectMeta) {
			params.Log.V(2).Info("no service updates needed", "name", desiredService.Name, "namespace", desiredService.Namespace)
			continue
		}

		patch := client.MergeFrom(&currentService)

		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply updates: %w", err)
		}

		params.Log.V(2).Info("updated service", "name", desiredService.Name, "namespace", desiredService.Namespace)
	}

	return nil
}
