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

package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway/reconcile"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	creconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
	Platform string
}

type ReconcileOperations struct {
	Run  func(context.Context, reconcile.Params) error
	Name string
}

func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("gateway", req.NamespacedName)
	gw := &securityv1.Gateway{}
	err := r.Get(ctx, req.NamespacedName, gw)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	ops := []ReconcileOperations{
		{reconcile.GatewayLicense, "gateway license"},
		{reconcile.Secrets, "secrets"},
		{reconcile.Services, "services"},
		{reconcile.ServiceAccount, "service account"},
		{reconcile.Ingress, "ingress"},
		{reconcile.Routes, "openshift routes"},
		{reconcile.HorizontalPodAutoscaler, "horizontalPodAutoscaler"},
		{reconcile.PodDisruptionBudget, "podDisruptionBudget"},
		{reconcile.GatewayStatus, "gatewayStatus"},
		{reconcile.ConfigMaps, "configMaps"},
		{reconcile.Deployment, "deployment"},
		{reconcile.ManagementPod, "management pod"},
		{reconcile.ClusterProperties, "cluster properties"},
		{reconcile.ListenPorts, "listen ports"},
		{reconcile.ExternalRepository, "repository references"},
		{reconcile.ExternalSecrets, "external secrets"},
		{reconcile.ExternalKeys, "external keys"},
		{reconcile.ExternalCerts, "external certs"},
	}

	if gw.Spec.App.Otk.Enabled {
		if gw.Spec.App.Otk.Type == securityv1.OtkTypeDMZ || gw.Spec.App.Otk.Type == securityv1.OtkTypeInternal {
			ops = append(ops, ReconcileOperations{reconcile.ScheduledJobs, "scheduled jobs"})
		}
		if gw.Spec.App.Otk.Type == securityv1.OtkTypeSingle && !gw.Spec.App.Management.Database.Enabled {
			ops = append(ops, ReconcileOperations{reconcile.OTKDatabaseMaintenanceTasks, "otk-db-maintenance-tasks"})
		}
	}

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: gw,
		Platform: r.Platform,
	}

	start := time.Now()
	for _, op := range ops {
		r.muTasks.Lock()
		err = op.Run(ctx, params)
		if err != nil {
			_ = captureMetrics(ctx, params, start, true, op.Name)
			// record failures here
			r.muTasks.Unlock()
			return ctrl.Result{}, err
		}
		r.muTasks.Unlock()
	}

	_ = captureMetrics(ctx, params, start, false, "")

	return ctrl.Result{RequeueAfter: 12 * time.Hour}, nil
}

func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {

	builder := ctrl.NewControllerManagedBy(mgr).For(&securityv1.Gateway{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&appsv1.Deployment{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{})

	repo := &metav1.PartialObjectMetadata{}
	repo.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "security.brcmlabs.com",
		Version: "v1",
		Kind:    "repository",
	})

	builder.WatchesMetadata(repo,
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []creconcile.Request {
			rb, err := json.Marshal(a.DeepCopyObject())
			var repository securityv1.Repository
			if err != nil {
				return []creconcile.Request{}
			}

			err = json.Unmarshal(rb, &repository)
			if err != nil {
				return []creconcile.Request{}
			}

			gatewayList := &securityv1.GatewayList{}
			listOpts := []client.ListOption{
				client.InNamespace(a.GetNamespace()),
			}
			err = r.List(ctx, gatewayList, listOpts...)

			if err != nil {
				if k8serrors.IsNotFound(err) {
					return []creconcile.Request{}
				}
			}
			req := []creconcile.Request{}
			for _, gateway := range gatewayList.Items {
				for _, repoRef := range gateway.Spec.App.RepositoryReferences {
					if repoRef.Name == repository.Name {
						req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: gateway.Namespace, Name: gateway.Name}})
					}
				}
			}
			return req
		}),
	)

	s := &metav1.PartialObjectMetadata{}
	s.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "secret",
	})
	builder.WatchesMetadata(s,
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []creconcile.Request {
			gatewayList := &securityv1.GatewayList{}
			listOpts := []client.ListOption{
				client.InNamespace(a.GetNamespace()),
			}
			err := r.List(ctx, gatewayList, listOpts...)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					return []creconcile.Request{}
				}
			}
			req := []creconcile.Request{}
			for _, gateway := range gatewayList.Items {
				for _, secretRef := range gateway.Spec.App.ExternalSecrets {
					if secretRef.Name == a.GetName() {
						req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: gateway.Namespace, Name: gateway.Name}})
					}
				}
				for _, keyRef := range gateway.Spec.App.ExternalKeys {
					if keyRef.Name == a.GetName() {
						req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: gateway.Namespace, Name: gateway.Name}})
					}
				}
				for _, keyRef := range gateway.Spec.App.ExternalCerts {
					if keyRef.Name == a.GetName() {
						req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: gateway.Namespace, Name: gateway.Name}})
					}
				}
			}
			return req
		}),
	).WithEventFilter(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew.GetObjectKind().GroupVersionKind().Kind == "gateway" {
				oldGwGen := e.ObjectOld.GetGeneration()
				newGwGen := e.ObjectNew.GetGeneration()
				return oldGwGen != newGwGen
			}
			objType := fmt.Sprintf("%v", reflect.TypeOf(e.ObjectNew))
			if e.ObjectNew.GetObjectKind().GroupVersionKind().Kind == "" {
				if objType == "*v1.Gateway" {
					oldGw := securityv1.Gateway{}
					newGw := securityv1.Gateway{}
					oldGwB, err := json.Marshal(e.ObjectOld.DeepCopyObject())
					if err != nil {
						return true
					}
					newGwB, err := json.Marshal(e.ObjectNew.DeepCopyObject())
					if err != nil {
						return true
					}

					err = json.Unmarshal(oldGwB, &oldGw)
					if err != nil {
						return true
					}
					err = json.Unmarshal(newGwB, &newGw)
					if err != nil {
						return true
					}

					if reflect.DeepEqual(oldGw.Spec, newGw.Spec) {
						return false
					}
				}

				if objType == "*v1.Deployment" {
					oldDep := appsv1.Deployment{}
					newDep := appsv1.Deployment{}
					oldDepB, err := json.Marshal(e.ObjectOld.DeepCopyObject())
					if err != nil {
						return true
					}
					newDepB, err := json.Marshal(e.ObjectNew.DeepCopyObject())
					if err != nil {
						return true
					}

					err = json.Unmarshal(oldDepB, &oldDep)
					if err != nil {
						return true
					}
					err = json.Unmarshal(newDepB, &newDep)
					if err != nil {
						return true
					}
					if oldDep.Status.ReadyReplicas == newDep.Status.ReadyReplicas || newDep.Status.ReadyReplicas == 0 { //|| oldDep.Status.ReadyReplicas > newDep.Status.ReadyReplicas
						return false
					}
					return true
				}

				if objType == "*v1.PodDisruptionBudget" {
					oldPdb := policyv1.PodDisruptionBudget{}
					newPdb := policyv1.PodDisruptionBudget{}
					oldPdbB, err := json.Marshal(e.ObjectOld.DeepCopyObject())
					if err != nil {
						return true
					}
					newPdbB, err := json.Marshal(e.ObjectNew.DeepCopyObject())
					if err != nil {
						return true
					}

					err = json.Unmarshal(oldPdbB, &oldPdb)
					if err != nil {
						return true
					}
					err = json.Unmarshal(newPdbB, &newPdb)
					if err != nil {
						return true
					}
					if reflect.DeepEqual(oldPdb.Spec, newPdb.Spec) {
						return false
					}

					return true
				}
				if objType == "*v2.HorizontalPodAutoscaler" {
					oldHpa := autoscalingv2.HorizontalPodAutoscaler{}
					newHpa := autoscalingv2.HorizontalPodAutoscaler{}
					oldHpaB, err := json.Marshal(e.ObjectOld.DeepCopyObject())
					if err != nil {
						return true
					}
					newHpaB, err := json.Marshal(e.ObjectNew.DeepCopyObject())
					if err != nil {
						return true
					}

					err = json.Unmarshal(oldHpaB, &oldHpa)
					if err != nil {
						return true
					}
					err = json.Unmarshal(newHpaB, &newHpa)
					if err != nil {
						return true
					}
					if reflect.DeepEqual(oldHpa.Spec, newHpa.Spec) {
						return false
					}

					return true
				}

			}
			return true
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return true
		},
	})

	if r.Platform == "openshift" {
		builder.Owns(&routev1.Route{})
	}
	return builder.Complete(r)
}

func captureMetrics(ctx context.Context, params reconcile.Params, start time.Time, hasError bool, opName string) error {

	gateway := params.Instance
	operatorNamespace, err := util.GetOperatorNamespace()
	if err != nil {
		params.Log.Info("could not determine operator namespace")
		return err
	}
	otelEnabled, err := util.GetOtelEnabled()
	if err != nil {
		params.Log.Info("could not determine if OTel is enabled")
		return err
	}

	if !otelEnabled {
		return nil
	}

	otelMetricPrefix, err := util.GetOtelMetricPrefix()
	if err != nil {
		params.Log.Info("could not determine otel metric prefix")
		return err
	}

	if otelMetricPrefix == "" {
		otelMetricPrefix = "layer7_"
	}

	hostname, err := util.GetHostname()
	if err != nil {
		params.Log.Error(err, "failed to retrieve operator hostname")
		return err
	}

	meter := otel.Meter("layer7-operator-gateway-controller-metrics")

	reconcileLatency, err := meter.Float64Histogram(otelMetricPrefix+"operator_gw_reconciler_latency",
		metric.WithDescription("gateway controller reconcile latency"), metric.WithUnit("ms"))
	if err != nil {
		return err
	}

	gatewayReconcileTotal, err := meter.Int64Counter(otelMetricPrefix+"operator_gw_reconcile_total",
		metric.WithDescription("gateway reconcile total"))
	if err != nil {
		return err
	}

	gatewayReconcileSuccess, err := meter.Int64Counter(otelMetricPrefix+"operator_gw_reconcile_success",
		metric.WithDescription("gateway reconcile success"))
	if err != nil {
		return err
	}

	gatewayReconcileFailure, err := meter.Int64Counter(otelMetricPrefix+"operator_gw_reconcile_failure",
		metric.WithDescription("gateway reconcile failure"))
	if err != nil {
		return err
	}

	gatewayExternalRefs, err := meter.Int64Gauge(otelMetricPrefix+"operator_gw_external_references",
		metric.WithDescription("operator managed gateway repository references"))
	if err != nil {
		return err
	}

	duration := time.Since(start)
	reconcileLatency.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
		))

	externalRefs := len(gateway.Spec.App.RepositoryReferences) + len(gateway.Spec.App.ExternalSecrets) + len(gateway.Spec.App.ExternalKeys)

	gatewayExternalRefs.Record(ctx, int64(externalRefs),
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.Int("repository_references", len(gateway.Spec.App.RepositoryReferences)),
			attribute.Int("external_keys", len(gateway.Spec.App.ExternalKeys)),
			attribute.Int("external_secrets", len(gateway.Spec.App.ExternalSecrets)),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	gatewayReconcileTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	if hasError {
		gatewayReconcileFailure.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("operation", opName),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))
	} else {
		gatewayReconcileSuccess.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))
	}
	return nil
}
