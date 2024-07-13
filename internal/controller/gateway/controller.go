/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gateway

import (
	"context"
	"fmt"
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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
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
		{reconcile.Route, "openshift route"},
		{reconcile.HorizontalPodAutoscaler, "horizontalPodAutoscaler"},
		{reconcile.PodDisruptionBudget, "podDisruptionBudget"},
		{reconcile.GatewayStatus, "gatewayStatus"},
		{reconcile.ConfigMaps, "configMaps"},
		{reconcile.Deployment, "deployment"},
		{reconcile.ManagementPod, "management pod"},
		{reconcile.ExternalRepository, "repository references"},
		{reconcile.ExternalSecrets, "external secrets"},
		{reconcile.ExternalKeys, "external keys"},
	}

	if gw.Spec.App.Otk.Enabled {
		ops = append(ops, ReconcileOperations{reconcile.ScheduledJobs, "scheduled jobs"})
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
		err = op.Run(ctx, params)
		if err != nil {
			log.Error(err, fmt.Sprintf("failed to reconcile %s", op.Name))
			_ = captureMetrics(ctx, params, start, true)
			return ctrl.Result{}, err
		}
	}

	_ = captureMetrics(ctx, params, start, false)

	return ctrl.Result{RequeueAfter: 12 * time.Hour}, nil
}

func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {

	builder := ctrl.NewControllerManagedBy(mgr).For(&securityv1.Gateway{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&appsv1.Deployment{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{})

	builder.WatchesMetadata(&securityv1.Repository{},
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
				for _, repoRef := range gateway.Spec.App.RepositoryReferences {
					if repoRef.Name == a.GetName() {
						req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: gateway.Namespace, Name: gateway.Name}})
					}
				}
			}
			return req
		}),
	)

	builder.WatchesMetadata(&corev1.Secret{},
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
			}
			return req
		}),
	)

	if r.Platform == "openshift" {
		builder.Owns(&routev1.Route{})
	}
	return builder.Complete(r)
}

func captureMetrics(ctx context.Context, params reconcile.Params, start time.Time, hasError bool) error {

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
	if err != nil {
		params.Log.Error(err, "failed to retrieve operator namespace")
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

	gatewayRepoRefs, err := meter.Int64Gauge(otelMetricPrefix+"operator_gw_repository_references",
		metric.WithDescription("operator managed gateway repository references"))
	if err != nil {
		return err
	}

	gatewaySecretRefs, err := meter.Int64Gauge(otelMetricPrefix+"operator_gw_secret_references",
		metric.WithDescription("operator managed gateway repository references"))
	if err != nil {
		return err
	}

	gatewayKeyRefs, err := meter.Int64Gauge(otelMetricPrefix+"operator_gw_key_references",
		metric.WithDescription("operator managed gateway repository references"))
	if err != nil {
		return err
	}

	duration := time.Since(start)
	reconcileLatency.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
		))

	gatewayRepoRefs.Record(ctx, int64(len(gateway.Spec.App.RepositoryReferences)),
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	gatewaySecretRefs.Record(ctx, int64(len(gateway.Spec.App.ExternalSecrets)),
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	gatewayKeyRefs.Record(ctx, int64(len(gateway.Spec.App.ExternalKeys)),
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	gatewayReconcileTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("pod", hostname),
			attribute.String("namespace", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	if hasError {
		gatewayReconcileFailure.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("pod", hostname),
				attribute.String("namespace", operatorNamespace),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))
	} else {
		gatewayReconcileSuccess.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("pod", hostname),
				attribute.String("namespace", operatorNamespace),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))
	}

	return nil
}
