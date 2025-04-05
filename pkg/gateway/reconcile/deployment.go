package reconcile

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"maps"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway"
	"github.com/caapim/layer7-operator/pkg/util"
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
	maps.Copy(updatedDeployment.Spec.Template.ObjectMeta.Annotations, desiredDeployment.Spec.Template.ObjectMeta.Annotations)

	for k, v := range desiredDeployment.ObjectMeta.Labels {
		updatedDeployment.ObjectMeta.Labels[k] = v
	}

	for k, v := range desiredDeployment.Spec.Template.ObjectMeta.Labels {
		updatedDeployment.Spec.Template.ObjectMeta.Labels[k] = v
	}

	desiredDeployment, err = setStateStoreConfig(ctx, params, desiredDeployment)
	if err != nil {
		return err
	}

	updatedDeployment.Spec.Template.Spec.InitContainers = desiredDeployment.Spec.Template.Spec.InitContainers
	updatedDeployment.Spec.Template.Spec.Volumes = desiredDeployment.Spec.Template.Spec.Volumes

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
			configMaps = append(configMaps, params.Instance.Name+"-otk-shared-init-config", params.Instance.Name+"-otk-install-init-config")
			if params.Instance.Spec.App.Otk.Database.Type != securityv1.OtkDatabaseTypeCassandra {
				configMaps = append(configMaps, params.Instance.Name+"-otk-db-init-config")
			}
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
		//if !params.Instance.Spec.App.Management.Database.Enabled {
		if params.Instance.Spec.App.Management.SecretName == "" {
			if !params.Instance.Spec.App.Management.DisklessConfig.Disabled {
				secrets = append(secrets, params.Instance.Name+"-node-properties")
			} else {
				secrets = append(secrets, params.Instance.Name)
			}
		}
		//}
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
	commits := ""
	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		for _, repoStatus := range params.Instance.Status.RepositoryStatus {
			if repoRef.Name == repoStatus.Name && repoRef.Type == "static" && repoRef.Enabled {
				commits = commits + repoStatus.Commit
			}
		}
	}

	h := sha1.New()
	h.Write([]byte(commits))
	commits = fmt.Sprintf("%x", h.Sum(nil))

	dep.ObjectMeta.Labels["security.brcmlabs.com/static-repositories-checksum"] = commits

	return dep, nil
}

func setStateStoreConfig(ctx context.Context, params Params, dep *appsv1.Deployment) (*appsv1.Deployment, error) {
	defaultMode := int32(0755)
	optional := false
	if len(params.Instance.Spec.App.RepositoryReferences) > 0 {
		stateStores := []string{}
		for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {

			repo := securityv1.Repository{}
			err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, &repo)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve repository: %s", repoRef.Name)
			}
			if repo.Spec.StateStoreReference != "" {
				if !util.Contains(stateStores, repo.Spec.StateStoreReference) {
					stateStores = append(stateStores, repo.Spec.StateStoreReference)
				}
			}
		}

		for i, ic := range dep.Spec.Template.Spec.InitContainers {
			if strings.Contains(ic.Name, "graphman-static-init") {
				for _, stateStore := range stateStores {
					dep.Spec.Template.Spec.InitContainers[i].VolumeMounts = append(dep.Spec.Template.Spec.InitContainers[i].VolumeMounts, corev1.VolumeMount{
						Name:      stateStore + "-secret",
						MountPath: "/graphman/statestore-secret/" + stateStore,
					})
					vs := corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
						SecretName:  stateStore + "-secret",
						DefaultMode: &defaultMode,
						Optional:    &optional,
					}}

					dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
						Name:         stateStore + "-secret",
						VolumeSource: vs,
					})

					dep.Spec.Template.Spec.InitContainers[i].VolumeMounts = append(dep.Spec.Template.Spec.InitContainers[i].VolumeMounts, corev1.VolumeMount{
						Name:      stateStore + "-config-secret",
						MountPath: "/graphman/statestore-config/" + stateStore + "/config.json",
						SubPath:   "config.json",
					})
					vs = corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
						SecretName:  stateStore + "-config-secret",
						DefaultMode: &defaultMode,
						Optional:    &optional,
						Items: []corev1.KeyToPath{{
							Path: "config.json",
							Key:  "config.json"},
						}},
					}

					dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
						Name:         stateStore + "-config-secret",
						VolumeSource: vs,
					})

				}
			}
		}
	}
	return dep, nil
}
