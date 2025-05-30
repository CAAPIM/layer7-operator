/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
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
