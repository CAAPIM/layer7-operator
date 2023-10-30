package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// / TODO: Should auto-delete if not used...
func ServiceAccount(ctx context.Context, params Params) error {

	if !params.Instance.Spec.App.ServiceAccount.Create {
		return nil
	}

	serviceAccountName := params.Instance.Spec.App.ServiceAccount.Name
	if params.Instance.Spec.App.ServiceAccount.Name == "" {
		serviceAccountName = params.Instance.Name
	}

	desiredSA := gateway.NewServiceAccount(params.Instance)
	currentSA := &corev1.ServiceAccount{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredSA, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: serviceAccountName, Namespace: params.Instance.Namespace}, currentSA)

	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredSA); err != nil {
			return fmt.Errorf("failed creating service account: %w", err)
		}
		params.Log.Info("created service account", "name", serviceAccountName, "namespace", params.Instance.Namespace)
		return nil
	}

	if err != nil {
		return err
	}

	return nil

}
