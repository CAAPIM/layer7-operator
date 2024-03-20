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

package portal

import (
	"context"
	"sync"

	corev1 "k8s.io/api/core/v1"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/portal/reconcile"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/// TODO: refactor

// L7PortalReconciler reconciles a Gateway object
type L7PortalReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

func (r *L7PortalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("portal", req.NamespacedName)

	portal := &securityv1alpha1.L7Portal{}
	err := r.Get(ctx, req.NamespacedName, portal)
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
		Instance: portal,
	}

	err = reconcile.Jobs(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Status(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *L7PortalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.L7Portal{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
