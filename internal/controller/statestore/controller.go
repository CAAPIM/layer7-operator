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

	err = reconcile.RedisStateStore(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Secret(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *L7StateStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.L7StateStore{}).
		// add a watch for repositories
		Complete(r)
}
