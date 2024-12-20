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

package api

import (
	"context"
	"sync"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/api/reconcile"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const apiFinalizer = "security.brcmlabs.com/finalizer"

// L7ApiReconciler reconciles a L7Api object
type L7ApiReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

func (r *L7ApiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := r.Log.WithValues("L7Api", req.NamespacedName)

	l7Api := &securityv1alpha1.L7Api{}
	err := r.Get(ctx, req.NamespacedName, l7Api)
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
		Instance: l7Api,
	}

	// Add a finalizer for the L7Api
	if !controllerutil.ContainsFinalizer(l7Api, apiFinalizer) {
		controllerutil.AddFinalizer(l7Api, apiFinalizer)
		err = r.Update(ctx, l7Api)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconcile.Status(ctx, params)
	if err != nil {
		r.Log.Info("error reconciling status", "error", err.Error())
		return ctrl.Result{}, err
	}

	err = reconcile.WriteTempStorage(ctx, params)
	if err != nil {
		r.Log.Info("error reconciling temp storage", "error", err.Error())
		return ctrl.Result{}, err
	}

	err = reconcile.L7Portal(ctx, params)
	if err != nil {
		r.Log.Info("error reconciling L7Portal", "error", err.Error())
		return ctrl.Result{}, err
	}

	err = reconcile.Gateway(ctx, params)
	if err != nil {
		r.Log.Info("error reconciling gateway", "error", err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *L7ApiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.L7Api{}).
		Complete(r)
}
