package reconcile

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/caapim/layer7-operator/pkg/gateway"
	routev1 "github.com/openshift/api/route/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Route(ctx context.Context, params Params) error {

	//Potentially delete the route if it exists.
	if !params.Instance.Spec.App.Ingress.Enabled || strings.ToLower(params.Platform) != "openshift" || strings.ToLower(params.Instance.Spec.App.Ingress.Type) != "route" {
		return nil
	}

	desiredRoute := gateway.NewRoute(params.Instance)
	currentRoute := routev1.Route{}

	if err := controllerutil.SetControllerReference(params.Instance, &desiredRoute, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredRoute.Name, Namespace: params.Instance.Namespace}, &currentRoute)
	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, &desiredRoute); err != nil {
			return err
		}
		params.Log.Info("created route", "name", desiredRoute.Name, "namespace", params.Instance.Namespace)
		return nil
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(currentRoute.Spec, desiredRoute.Spec) {
		params.Log.V(2).Info("no route updates needed", "name", desiredRoute.Name, "namespace", desiredRoute.Namespace)
		return nil
	}

	updatedRoute := currentRoute.DeepCopy()
	updatedRoute.ObjectMeta.OwnerReferences = desiredRoute.ObjectMeta.OwnerReferences
	updatedRoute.Spec.To = desiredRoute.Spec.To
	updatedRoute.Spec.TLS = desiredRoute.Spec.TLS
	updatedRoute.Spec.AlternateBackends = desiredRoute.Spec.AlternateBackends
	updatedRoute.Spec.WildcardPolicy = desiredRoute.Spec.WildcardPolicy

	for k, v := range desiredRoute.ObjectMeta.Annotations {
		updatedRoute.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desiredRoute.ObjectMeta.Labels {
		updatedRoute.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(&currentRoute)

	if err := params.Client.Patch(ctx, updatedRoute, patch); err != nil {
		return err
	}
	params.Log.Info("route updated", "name", desiredRoute.Name, "namespace", desiredRoute.Namespace)

	return nil
}
