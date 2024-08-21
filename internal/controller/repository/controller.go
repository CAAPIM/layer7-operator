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

package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/repository/reconcile"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	creconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// RepositoryReconciler reconciles a Repository object
type RepositoryReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

type ReconcileOperations struct {
	Run  func(context.Context, reconcile.Params) error
	Name string
}

func (r *RepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	log := r.Log.WithValues("repository", req.NamespacedName)

	ops := []ReconcileOperations{
		{reconcile.Secret, "secrets"},
		{reconcile.LocalReference, "local repositories"},
		{reconcile.ScheduledJobs, "scheduled jobs"},
	}

	l7repository := &securityv1.Repository{}

	err := r.Get(ctx, req.NamespacedName, l7repository)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: l7repository,
	}

	for _, op := range ops {
		err = op.Run(ctx, params)
		if err != nil {
			log.Error(err, fmt.Sprintf("failed to reconcile %s", op.Name))
			_ = captureMetrics(ctx, params, start, req.NamespacedName.Namespace)
			return ctrl.Result{}, err
		}
	}

	_ = captureMetrics(ctx, params, start, req.NamespacedName.Namespace)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	builder := ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.Repository{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{})

	builder.WatchesMetadata(&corev1.Secret{},
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []creconcile.Request {
			repoList := &securityv1.RepositoryList{}
			listOpts := []client.ListOption{
				client.InNamespace(a.GetNamespace()),
			}
			err := r.List(ctx, repoList, listOpts...)

			if err != nil {
				if k8serrors.IsNotFound(err) {
					return []creconcile.Request{}
				}
			}
			req := []creconcile.Request{}
			for _, repo := range repoList.Items {
				if strings.ToLower(repo.Spec.Type) == "local" {
					req = append(req, creconcile.Request{NamespacedName: types.NamespacedName{Namespace: repo.Namespace, Name: repo.Name}})
				}
			}
			return req
		}),
	)

	return builder.Complete(r)
}

func captureMetrics(ctx context.Context, params reconcile.Params, start time.Time, namespace string) error {

	otelEnabled, err := util.GetOtelEnabled()
	if err != nil {
		params.Log.Info("could not determine if OTel is enabled")
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

	meter := otel.Meter("layer7-operator-repository-controller-metrics")
	reconcileLatency, err := meter.Float64Histogram(otelMetricPrefix+"operator_repo_reconciler_latency",
		metric.WithDescription("repository controller reconcile latency"), metric.WithUnit("ms"))
	if err != nil {
		return err
	}

	duration := time.Since(start)
	reconcileLatency.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("k8s.namespace.name", namespace),
			attribute.String("k8s.pod.name", hostname),
		))

	return nil
}
