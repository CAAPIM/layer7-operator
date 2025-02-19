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

package statestore

import (
	"context"
	"sync"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/statestore/reconcile"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
)

// L7StateStoreReconciler reconciles a L7StateStore object
type L7StateStoreReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

func (r *L7StateStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("L7StateStore", req.NamespacedName)
	log.Info("connecting to statestore", "statestore", "stateStore.Name")
	stateStore := &securityv1alpha1.L7StateStore{}
	err := r.Get(ctx, req.NamespacedName, stateStore)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.Info("connecting to statestore", "statestore", stateStore.Name)

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: stateStore,
	}

	// listRepositories
	// if repository has statestore reference then sync that (checksum).
	// otherwise do nothing..

	err = reconcile.RedisStateStore(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}
	// just test the connection???

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *L7StateStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.L7StateStore{}).
		// add a watch for repositories
		Complete(r)
}
