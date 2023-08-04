package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/repository"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Secret(ctx context.Context, params Params) error {

	if params.Instance.Spec.Auth.ExistingSecretName != "" {
		return nil
	}

	data := map[string][]byte{
		"USERNAME": []byte(params.Instance.Spec.Auth.Username),
		"PASSWORD": []byte(params.Instance.Spec.Auth.Password),
		"TOKEN":    []byte(params.Instance.Spec.Auth.Token),
	}
	desiredSecret := repository.NewSecret(params.Instance, params.Instance.Name, data)

	if err := reconcileSecret(ctx, params, desiredSecret); err != nil {
		return fmt.Errorf("failed to reconcile secrets: %w", err)
	}

	return nil
}

func StorageSecret(ctx context.Context, params Params) error {

	storageSecretName := params.Instance.Name + "-repository"
	bundleGzip, err := util.CompressGraphmanBundle("/tmp/" + params.Instance.Spec.Name)
	if err != nil {
		return err
	}

	data := map[string][]byte{
		params.Instance.Name + ".gz": bundleGzip,
	}

	desiredSecret := repository.NewSecret(params.Instance, storageSecretName, data)

	if err := reconcileSecret(ctx, params, desiredSecret); err != nil {
		return fmt.Errorf("failed to reconcile secrets: %w", err)
	}

	return nil
}

func reconcileSecret(ctx context.Context, params Params, desiredSecret *corev1.Secret) error {

	if err := controllerutil.SetControllerReference(params.Instance, desiredSecret, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	currentSecret := corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredSecret.Name, Namespace: params.Instance.Namespace}, &currentSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredSecret); err != nil {
			return err
		}
		params.Log.Info("created secret", "name", desiredSecret.Name, "namespace", params.Instance.Namespace)
	}
	if err != nil {
		return err
	}

	if desiredSecret.ObjectMeta.Annotations["checksum/data"] != currentSecret.ObjectMeta.Annotations["checksum/data"] {
		patch := client.MergeFrom(&currentSecret)
		if err := params.Client.Patch(ctx, desiredSecret, patch); err != nil {
			return err
		}
		params.Log.V(2).Info("secret updated", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
	}

	return nil
}
