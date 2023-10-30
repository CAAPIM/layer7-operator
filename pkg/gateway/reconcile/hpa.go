package reconcile

import (
	"context"
	"fmt"
	"reflect"

	"github.com/caapim/layer7-operator/pkg/gateway"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// / TODO: Should auto-delete if not used...
func HorizontalPodAutoscaler(ctx context.Context, params Params) error {

	if !params.Instance.Spec.App.Autoscaling.Enabled {
		return nil
	}

	desiredHpa := gateway.NewHPA(params.Instance)
	currentHpa := &autoscalingv2.HorizontalPodAutoscaler{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredHpa, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, currentHpa)

	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredHpa); err != nil {
			return fmt.Errorf("failed creating horizontal pod autoscaler updates: %w", err)
		}
		params.Log.Info("created horizontal pod autoscaler updates", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	if err != nil {
		return err
	}

	if reflect.DeepEqual(currentHpa.Spec.Behavior, desiredHpa.Spec.Behavior) && reflect.DeepEqual(currentHpa.Spec.MaxReplicas, desiredHpa.Spec.MaxReplicas) && reflect.DeepEqual(currentHpa.Spec.MinReplicas, desiredHpa.Spec.MinReplicas) && reflect.DeepEqual(currentHpa.Spec.Metrics, desiredHpa.Spec.Metrics) {
		params.Log.V(2).Info("no horizontal pod autoscaler updates needed", "name", desiredHpa.Name, "namespace", desiredHpa.Namespace)
		return nil
	}

	updated := currentHpa.DeepCopy()
	updated.Spec = desiredHpa.Spec

	updated.ObjectMeta.OwnerReferences = desiredHpa.ObjectMeta.OwnerReferences

	for k, v := range desiredHpa.ObjectMeta.Annotations {
		updated.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desiredHpa.ObjectMeta.Labels {
		updated.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(currentHpa)

	if err := params.Client.Patch(ctx, updated, patch); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	params.Log.Info("updated horizontal pod autoscaler", "name", desiredHpa.Name, "namespace", desiredHpa.Namespace)

	return nil

}
