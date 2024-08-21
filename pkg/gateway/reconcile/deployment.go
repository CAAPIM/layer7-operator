package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Deployment(ctx context.Context, params Params) error {
	desiredDeployment := gateway.NewDeployment(params.Instance, params.Platform)
	currentDeployment := &appsv1.Deployment{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredDeployment, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, currentDeployment)

	if err != nil && k8serrors.IsNotFound(err) {

		if err = params.Client.Create(ctx, desiredDeployment); err != nil {
			return fmt.Errorf("failed creating deployment: %w", err)
		}
		params.Log.Info("created deployment", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	if err != nil {
		return err
	}

	if params.Instance.Spec.App.Autoscaling.Enabled {
		desiredDeployment.Spec.Replicas = currentDeployment.Spec.Replicas
	}

	updatedDeployment := currentDeployment.DeepCopy()
	updatedDeployment.Spec = desiredDeployment.Spec

	updatedDeployment.ObjectMeta.OwnerReferences = desiredDeployment.ObjectMeta.OwnerReferences

	if params.Instance.Spec.App.Autoscaling.Enabled {
		updatedDeployment.Spec.Replicas = currentDeployment.Spec.Replicas
	}

	for k, v := range desiredDeployment.ObjectMeta.Annotations {
		updatedDeployment.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desiredDeployment.Spec.Template.ObjectMeta.Annotations {
		updatedDeployment.Spec.Template.ObjectMeta.Annotations[k] = v
	}

	for k, v := range desiredDeployment.ObjectMeta.Labels {
		updatedDeployment.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(currentDeployment)

	if err := params.Client.Patch(ctx, updatedDeployment, patch); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	params.Log.V(2).Info("updated deployment", "name", desiredDeployment.Name, "namespace", desiredDeployment.Namespace)

	return nil

}
