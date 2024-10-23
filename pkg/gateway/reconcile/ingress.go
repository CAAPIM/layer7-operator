package reconcile

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/caapim/layer7-operator/pkg/gateway"
	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Ingress(ctx context.Context, params Params) error {
	if strings.ToLower(params.Instance.Spec.App.Ingress.Type) == "route" {
		return nil
	}

	desiredIngress := gateway.NewIngress(params.Instance)
	currentIngress := &networkingv1.Ingress{}

	if desiredIngress != nil {
		if err := controllerutil.SetControllerReference(params.Instance, desiredIngress, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
	}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, currentIngress)

	if !params.Instance.Spec.App.Ingress.Enabled && !k8serrors.IsNotFound(err) && controllerutil.HasControllerReference(currentIngress) {
		if err := params.Client.Delete(ctx, currentIngress); err != nil {
			return err
		}
		params.Log.Info("removed ingress", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	if !params.Instance.Spec.App.Ingress.Enabled {
		return nil
	}

	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredIngress); err != nil {
			return err
		}
		params.Log.Info("created ingress", "name", desiredIngress.Name, "namespace", params.Instance.Namespace)
		return nil
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(currentIngress.Spec, desiredIngress.Spec) && reflect.DeepEqual(currentIngress.ObjectMeta.Annotations, desiredIngress.ObjectMeta.Annotations) {
		params.Log.V(2).Info("no ingress updates needed", "name", desiredIngress.Name, "namespace", desiredIngress.Namespace)
		return nil
	}

	updatedIngress := currentIngress.DeepCopy()

	updatedIngress.ObjectMeta.OwnerReferences = desiredIngress.ObjectMeta.OwnerReferences
	updatedIngress.Spec.Rules = desiredIngress.Spec.Rules
	updatedIngress.Spec.TLS = desiredIngress.Spec.TLS
	updatedIngress.Spec.DefaultBackend = desiredIngress.Spec.DefaultBackend
	updatedIngress.Spec.IngressClassName = desiredIngress.Spec.IngressClassName
	updatedIngress.ObjectMeta.Annotations = desiredIngress.ObjectMeta.Annotations
	updatedIngress.ObjectMeta.Labels = desiredIngress.ObjectMeta.Labels

	patch := client.MergeFrom(currentIngress)
	if err := params.Client.Patch(ctx, updatedIngress, patch); err != nil {
		return err
	}
	params.Log.Info("ingress updated", "name", desiredIngress.Name, "namespace", desiredIngress.Namespace)

	return nil
}
