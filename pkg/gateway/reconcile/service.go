package reconcile

import (
	"context"
	"fmt"
	"reflect"

	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Services(ctx context.Context, params Params) error {
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

		for k, v := range desiredService.ObjectMeta.Annotations {
			updated.ObjectMeta.Annotations[k] = v
		}
		for k, v := range desiredService.ObjectMeta.Labels {
			updated.ObjectMeta.Labels[k] = v
		}

		if reflect.DeepEqual(currentService.Spec.Ports, desiredService.Spec.Ports) && reflect.DeepEqual(currentService.Spec.Type, desiredService.Spec.Type) {
			params.Log.V(2).Info("no service updates needed", "name", desiredService.Name, "namespace", desiredService.Namespace)
			return nil
		}

		patch := client.MergeFrom(&currentService)

		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply updates: %w", err)
		}

		params.Log.V(2).Info("updated service", "name", desiredService.Name, "namespace", desiredService.Namespace)
	}

	return nil
}
