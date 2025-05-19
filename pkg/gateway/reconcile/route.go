package reconcile

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/caapim/layer7-operator/pkg/gateway"
	routev1 "github.com/openshift/api/route/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Routes(ctx context.Context, params Params) error {
	desiredRoutes := []*routev1.Route{}
	if strings.ToLower(params.Instance.Spec.App.Ingress.Type) == "ingress" {
		return nil
	}
	for index, route := range params.Instance.Spec.App.Ingress.Routes {
		if route.Port != nil {
			if route.Port.TargetPort == intstr.FromString("management") {
				desiredRoutes = append(desiredRoutes, gateway.NewRoute(params.Instance, route, "management", true))
				continue
			}
		}
		desiredRoutes = append(desiredRoutes, gateway.NewRoute(params.Instance, route, strconv.Itoa(index), false))
	}

	if err := reconcileRoutes(ctx, params, desiredRoutes); err != nil {
		return fmt.Errorf("failed to reconcile configMaps: %w", err)
	}
	return nil
}

func reconcileRoutes(ctx context.Context, params Params, desiredRoutes []*routev1.Route) error {
	for _, desiredRoute := range desiredRoutes {
		currentRoute := routev1.Route{}

		if err := controllerutil.SetControllerReference(params.Instance, desiredRoute, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		err := params.Client.Get(ctx, types.NamespacedName{Name: desiredRoute.Name, Namespace: params.Instance.Namespace}, &currentRoute)

		if !params.Instance.Spec.App.Ingress.Enabled && !k8serrors.IsNotFound(err) && controllerutil.HasControllerReference(&currentRoute) {
			if err := params.Client.Delete(ctx, &currentRoute); err != nil {
				return err
			}
			params.Log.Info("removed route", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
			continue
		}

		if !params.Instance.Spec.App.Ingress.Enabled || strings.ToLower(params.Platform) != "openshift" || strings.ToLower(params.Instance.Spec.App.Ingress.Type) != "route" {
			return nil
		}

		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredRoute); err != nil {
				return err
			}
			params.Log.Info("created route", "name", desiredRoute.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}

		if reflect.DeepEqual(currentRoute.Spec, desiredRoute.Spec) {
			params.Log.V(2).Info("no route updates needed", "name", desiredRoute.Name, "namespace", desiredRoute.Namespace)
			continue
		}

		updatedRoute := currentRoute.DeepCopy()
		updatedRoute.ObjectMeta.OwnerReferences = desiredRoute.ObjectMeta.OwnerReferences
		updatedRoute.Spec.To = desiredRoute.Spec.To
		updatedRoute.Spec.TLS = desiredRoute.Spec.TLS
		updatedRoute.Spec.Host = desiredRoute.Spec.Host
		updatedRoute.Spec.Port = desiredRoute.Spec.Port
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
	}
	return nil
}
