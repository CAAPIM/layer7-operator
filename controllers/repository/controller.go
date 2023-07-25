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
	"reflect"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/repository/secrets"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RepositoryReconciler reconciles a Repository object
type RepositoryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = r.Log.WithValues("repository", req.NamespacedName)

	repository := &securityv1.Repository{}

	err := r.Get(ctx, req.NamespacedName, repository)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if repository.Spec.Auth.ExistingSecretName != "" {
		_, err := getSecret(r, ctx, repository)
		if err != nil {
			r.Log.Info("Secret not found", "Name", repository.Name, "Namespace", repository.Namespace, "Secret", repository.Spec.Auth.ExistingSecretName)
			return ctrl.Result{}, nil
		}
	} else {
		err = reconcileSecret(r, ctx, repository)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconcileRepository(r, ctx, repository)
	if err != nil {
		r.Log.Error(err, "Reconcile Failure", "Namespace", repository.Namespace, "Name", repository.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 5}, nil
}

func reconcileSecret(r *RepositoryReconciler, ctx context.Context, repository *securityv1.Repository) error {
	currSecret := &corev1.Secret{}
	secret := secrets.NewSecret(repository)
	err := r.Get(ctx, types.NamespacedName{Name: repository.Name, Namespace: repository.Namespace}, currSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Secret", "Name", repository.Name, "Namespace", repository.Namespace)

		ctrl.SetControllerReference(repository, secret, r.Scheme)
		err = r.Create(ctx, secret)
		if err != nil {
			r.Log.Error(err, "Failed creating Secret", "Name", repository.Name, "Namespace", repository.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(currSecret.Data, secret.Data) {
		ctrl.SetControllerReference(repository, secret, r.Scheme)
		return r.Update(ctx, secret)
	}
	return nil
}

func reconcileRepository(r *RepositoryReconciler, ctx context.Context, repository *securityv1.Repository) error {
	repoStatus := repository.Status
	if !repository.Spec.Enabled {
		r.Log.Info("Repository not enabled", "Name", repository.Name, "Namespace", repository.Namespace)
		return nil
	}

	rSecret, err := getSecret(r, ctx, repository)

	if err != nil {
		r.Log.Info("Secret unavailable", "Name", repository.Name, "Namespace", repository.Namespace, "Error", err.Error())
		return nil
	}

	token := string(rSecret.Data["TOKEN"])
	if token == "" {
		token = string(rSecret.Data["PASSWORD"])
	}

	commit, err := util.CloneRepository(repository.Spec.Endpoint, string(rSecret.Data["USERNAME"]), token, repository.Spec.Branch, repository.Spec.Name)

	if err != nil {
		return err
	}

	storageSecretName := repository.Name + "-repository"

	err = reconcileRepositoryStorageSecret(r, ctx, repository)
	if err != nil {
		r.Log.Info("Failed to reconcile storage secret", "Name", repository.Name, "Namespace", repository.Namespace, "Error", err.Error())
		storageSecretName = ""
	}

	repoStatus.Commit = commit
	repoStatus.Name = repository.Name
	repoStatus.Vendor = repository.Spec.Auth.Vendor

	repoStatus.StorageSecretName = storageSecretName

	if !reflect.DeepEqual(repoStatus, repository.Status) {
		repoStatus.Updated = time.Now().String()
		repository.Status = repoStatus
		err = r.Client.Status().Update(ctx, repository)
		if err != nil {
			r.Log.Info("Failed to update repository status", "Namespace", repository.Namespace, "Name", repository.Name, "Message", err.Error())
		}
		r.Log.Info("Reconciled", "Name", repository.Name, "Namespace", repository.Namespace, "Commit", commit)
	}

	return nil
}

func getSecret(r *RepositoryReconciler, ctx context.Context, repository *securityv1.Repository) (*corev1.Secret, error) {
	repositorySecret := &corev1.Secret{}
	name := repository.Name
	if repository.Spec.Auth.ExistingSecretName != "" {
		name = repository.Spec.Auth.ExistingSecretName
	}

	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: repository.Namespace}, repositorySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return repositorySecret, err
			}
		}
	}
	return repositorySecret, nil
}

func reconcileRepositoryStorageSecret(r *RepositoryReconciler, ctx context.Context, repository *securityv1.Repository) error {
	repositoryStorageSecret := &corev1.Secret{}
	bundleGzip, err := util.CompressGraphmanBundle("/tmp/" + repository.Spec.Name)
	if err != nil {
		return err
	}
	secret := secrets.NewStorageSecret(repository, bundleGzip)
	name := repository.Name + "-repository"

	err = r.Get(ctx, types.NamespacedName{Name: name, Namespace: repository.Namespace}, repositoryStorageSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Storage Secret", "Name", name, "Namespace", repository.Namespace)

		ctrl.SetControllerReference(repository, secret, r.Scheme)
		err = r.Create(ctx, secret)
		if err != nil {
			r.Log.Error(err, "Failed creating Secret", "Name", repository.Name, "Namespace", repository.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(repositoryStorageSecret.Data, secret.Data) {
		r.Log.Info("Storage Secret Updated", "Name", name, "Namespace", repository.Namespace)
		ctrl.SetControllerReference(repository, secret, r.Scheme)
		return r.Update(ctx, secret)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.Repository{}).
		Complete(r)
}
