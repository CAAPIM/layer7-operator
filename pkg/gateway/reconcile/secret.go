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

func Secrets(ctx context.Context, params Params) error {

	desiredSecrets := []*corev1.Secret{}

	if params.Instance.Spec.App.Management.SecretName == "" {
		desiredSecrets = append(desiredSecrets, gateway.NewSecret(params.Instance, params.Instance.Name))
	}

	if params.Instance.Spec.App.Otk.Enabled && params.Instance.Spec.App.Otk.Database.Auth != (securityv1.OtkDatabaseAuth{}) && params.Instance.Spec.App.Otk.Database.Auth.ExistingSecret == "" {
		desiredSecrets = append(desiredSecrets, gateway.NewSecret(params.Instance, params.Instance.Name+"-otk-db-credentials"))
	}

	if err := reconcileSecrets(ctx, params, desiredSecrets); err != nil {
		return fmt.Errorf("failed to reconcile secrets: %w", err)
	}

	return nil
}

func reconcileSecrets(ctx context.Context, params Params, desiredSecrets []*corev1.Secret) error {
	for _, ds := range desiredSecrets {
		desiredSecret := ds
		if err := controllerutil.SetControllerReference(params.Instance, desiredSecret, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		currentSecret := corev1.Secret{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: desiredSecret.Name, Namespace: params.Instance.Namespace}, &currentSecret)
		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredSecret); err != nil {
				return err
			}
			params.Log.Info("created secret", "name", desiredSecret.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}

		if desiredSecret.ObjectMeta.Annotations["checksum/data"] != currentSecret.ObjectMeta.Annotations["checksum/data"] {
			patch := client.MergeFrom(&currentSecret)
			if err := params.Client.Patch(ctx, desiredSecret, patch); err != nil {
				return err
			}
			params.Log.V(2).Info("secret updated", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
		}
	}

	return nil
}
