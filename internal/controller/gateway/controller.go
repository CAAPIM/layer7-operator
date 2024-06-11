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
	"sync"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway/reconcile"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
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

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: gw,
		Platform: r.Platform,
	}

	err = reconcile.GatewayLicense(ctx, params)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	err = reconcile.Secrets(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Services(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.ServiceAccount(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Ingress(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Route(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.HorizontalPodAutoscaler(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.PodDisruptionBudget(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.GatewayStatus(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.ConfigMaps(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Deployment(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.ScheduledJobs(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {

	builder := ctrl.NewControllerManagedBy(mgr).For(&securityv1.Gateway{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&appsv1.Deployment{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{})
		//.Watches(&corev1.Secret{}, )

	if r.Platform == "openshift" {
		builder.Owns(&routev1.Route{})
	}

	return builder.Complete(r)
}
