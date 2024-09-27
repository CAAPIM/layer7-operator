package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Deployment(ctx context.Context, params Params) error {
	desiredDeployment := gateway.NewDeployment(params.Instance, params.Platform)
	currentDeployment := &appsv1.Deployment{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredDeployment, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, currentDeployment)

	if err != nil && k8serrors.IsNotFound(err) {

		desiredDeployment, err = setLabels(ctx, params, desiredDeployment)
		if err != nil {
			return err
		}

		if err = params.Client.Create(ctx, desiredDeployment); err != nil {
			return fmt.Errorf("failed creating deployment: %w", err)
		}

		params.Log.Info("created deployment", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	if err != nil {
		return err
	}

	if params.Instance.Spec.App.Autoscaling.Enabled {
		desiredDeployment.Spec.Replicas = currentDeployment.Spec.Replicas
	}

	updatedDeployment := currentDeployment.DeepCopy()
	updatedDeployment.Spec = desiredDeployment.Spec

	updatedDeployment.ObjectMeta.OwnerReferences = desiredDeployment.ObjectMeta.OwnerReferences

	desiredDeployment, err = setLabels(ctx, params, desiredDeployment)
	if err != nil {
		return err
	}

	if params.Instance.Spec.App.Autoscaling.Enabled {
		updatedDeployment.Spec.Replicas = currentDeployment.Spec.Replicas
	}

	for k, v := range desiredDeployment.ObjectMeta.Annotations {
		updatedDeployment.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desiredDeployment.Spec.Template.ObjectMeta.Annotations {
		updatedDeployment.Spec.Template.ObjectMeta.Annotations[k] = v
	}

	for k, v := range desiredDeployment.ObjectMeta.Labels {
		updatedDeployment.ObjectMeta.Labels[k] = v
	}

	for k, v := range desiredDeployment.Spec.Template.ObjectMeta.Labels {
		updatedDeployment.Spec.Template.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(currentDeployment)

	if err := params.Client.Patch(ctx, updatedDeployment, patch); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	params.Log.V(2).Info("updated deployment", "name", desiredDeployment.Name, "namespace", desiredDeployment.Namespace)

	return nil

}

func setLabels(ctx context.Context, params Params, dep *appsv1.Deployment) (*appsv1.Deployment, error) {
	restartOnConfigChange := false
	if params.Instance.Spec.App.RestartOnConfigChange {
		restartOnConfigChange = params.Instance.Spec.App.RestartOnConfigChange
	}

	if restartOnConfigChange {
		configMaps := []string{params.Instance.Name, params.Instance.Name + "-system", params.Instance.Name + "-gateway-files"}

		if params.Instance.Spec.App.Otk.Enabled && !params.Instance.Spec.App.Management.Database.Enabled {
			configMaps = append(configMaps, params.Instance.Name+"-otk-shared-init-config", params.Instance.Name+"-otk-install-init-config", params.Instance.Name+"-otk-db-init-config")
		}

		for _, cmName := range configMaps {
			cm := corev1.ConfigMap{}
			err := params.Client.Get(ctx, types.NamespacedName{Name: cmName, Namespace: params.Instance.Namespace}, &cm)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve configmap: %s", cmName)
			}
			for k, v := range cm.ObjectMeta.Annotations {
				if k == "checksum/data" {
					dep.ObjectMeta.Labels[cmName+"-checksum"] = v
				}
			}
		}

		secrets := []string{}
		if !params.Instance.Spec.App.Management.Database.Enabled {
			if params.Instance.Spec.App.Management.SecretName == "" {
				if !params.Instance.Spec.App.Management.DisklessConfig.Disabled {
					secrets = append(secrets, params.Instance.Name+"-node-properties")
				} else {
					secrets = append(secrets, params.Instance.Name)
				}
			}
		}
		if params.Instance.Spec.App.Redis.Enabled && params.Instance.Spec.App.Redis.ExistingSecret == "" {
			secrets = append(secrets, params.Instance.Name+"-shared-state-client-configuration")
		}

		for _, secretName := range secrets {
			secret := corev1.Secret{}
			err := params.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: params.Instance.Namespace}, &secret)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve configmap: %s", secretName)
			}
			for k, v := range secret.ObjectMeta.Annotations {
				if k == "checksum/data" {
					dep.ObjectMeta.Labels[secretName+"-checksum"] = v
				}
			}
		}
	}
	return dep, nil
}
