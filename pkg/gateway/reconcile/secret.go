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

func Secret(ctx context.Context, params Params) error {
	if params.Instance.Spec.App.Management.SecretName != "" {
		return nil
	}

	desiredSecret := gateway.NewSecret(params.Instance)
	currentSecret := corev1.Secret{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredSecret, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredSecret.Name, Namespace: params.Instance.Namespace}, &currentSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredSecret); err != nil {
			return err
		}
		params.Log.Info("created secret", "name", desiredSecret.Name, "namespace", params.Instance.Namespace)
		return nil
	}
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(currentSecret.Data, desiredSecret.Data) {
		patch := client.MergeFrom(&currentSecret)
		if err := params.Client.Patch(ctx, desiredSecret, patch); err != nil {
			return err
		}
		params.Log.V(2).Info("secret updated", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
	}

	return nil
}
