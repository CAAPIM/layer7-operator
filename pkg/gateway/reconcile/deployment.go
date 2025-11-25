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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package reconcile

import (
	"context"
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
	"strings"

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

// getOpenShiftUIDRange extracts the UID/GID range from OpenShift namespace annotations
// Returns the minimum UID and GID from the assigned range, or nil if not found/not OpenShift
func getOpenShiftUIDRange(ctx context.Context, k8sClient client.Client, namespace string) (*int64, *int64, error) {
	// Get the namespace
	ns := &corev1.Namespace{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: namespace}, ns); err != nil {
		return nil, nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	// Parse OpenShift UID annotation
	// Example format: "openshift.io/sa.scc.uid-range: 1001620000/10000"
	uidRange := ns.Annotations["openshift.io/sa.scc.uid-range"]
	if uidRange == "" {
		// Not OpenShift or no range set
		return nil, nil, nil
	}

	// Parse "1001620000/10000" format
	parts := strings.Split(uidRange, "/")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid uid-range format: %s", uidRange)
	}

	minUID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse UID: %w", err)
	}

	// Use the same value for GID (OpenShift typically assigns matching ranges)
	return &minUID, &minUID, nil
}

// applyDefaultCapabilities ensures all init containers and sidecars drop ALL capabilities by default
// This is a security best practice unless explicitly overridden by the user
func applyDefaultCapabilities(dep *appsv1.Deployment) {
	// Apply to all init containers
	for i := range dep.Spec.Template.Spec.InitContainers {
		if dep.Spec.Template.Spec.InitContainers[i].SecurityContext == nil {
			dep.Spec.Template.Spec.InitContainers[i].SecurityContext = &corev1.SecurityContext{}
		}
		if dep.Spec.Template.Spec.InitContainers[i].SecurityContext.Capabilities == nil {
			dep.Spec.Template.Spec.InitContainers[i].SecurityContext.Capabilities = &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			}
		}
	}

	// Apply to main container (gateway)
	for i := range dep.Spec.Template.Spec.Containers {
		if dep.Spec.Template.Spec.Containers[i].SecurityContext == nil {
			dep.Spec.Template.Spec.Containers[i].SecurityContext = &corev1.SecurityContext{}
		}
		if dep.Spec.Template.Spec.Containers[i].SecurityContext.Capabilities == nil {
			dep.Spec.Template.Spec.Containers[i].SecurityContext.Capabilities = &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			}
		}
	}
}

// applyOpenShiftSecurityDefaults applies OpenShift-specific security context defaults
// including UID/GID auto-detection and capability dropping
func applyOpenShiftSecurityDefaults(ctx context.Context, params Params) {
	// If either PodSecurityContext or ContainerSecurityContext has RunAsUser set, use it for both
	// Otherwise, auto-detect from OpenShift namespace annotations
	if params.Instance.Spec.App.PodSecurityContext.RunAsUser == nil || params.Instance.Spec.App.ContainerSecurityContext.RunAsUser == nil {
		var uid, gid *int64
		var runAsNonRootPtr *bool

		// Check if one is already set and use it (only UID/GID/RunAsNonRoot)
		if params.Instance.Spec.App.PodSecurityContext.RunAsUser != nil {
			uid = params.Instance.Spec.App.PodSecurityContext.RunAsUser
			gid = params.Instance.Spec.App.PodSecurityContext.RunAsGroup
			runAsNonRootPtr = params.Instance.Spec.App.PodSecurityContext.RunAsNonRoot
		} else if params.Instance.Spec.App.ContainerSecurityContext.RunAsUser != nil {
			// Copy only UID/GID/RunAsNonRoot from ContainerSecurityContext to PodSecurityContext
			// Don't copy container-specific fields like Capabilities or AllowPrivilegeEscalation
			uid = params.Instance.Spec.App.ContainerSecurityContext.RunAsUser
			gid = params.Instance.Spec.App.ContainerSecurityContext.RunAsGroup
			runAsNonRootPtr = params.Instance.Spec.App.ContainerSecurityContext.RunAsNonRoot
		} else {
			// Neither is set, auto-detect from OpenShift
			var err error
			uid, gid, err = getOpenShiftUIDRange(ctx, params.Client, params.Instance.Namespace)
			if err != nil {
				params.Log.V(2).Info("failed to detect OpenShift UID range, using defaults",
					"namespace", params.Instance.Namespace,
					"error", err.Error())
				uid = nil
			} else if uid != nil && gid != nil {
				params.Log.V(2).Info("auto-detected OpenShift UID/GID from namespace annotations",
					"namespace", params.Instance.Namespace,
					"uid", *uid,
					"gid", *gid)
				runAsNonRoot := true
				runAsNonRootPtr = &runAsNonRoot
			}
		}

		// Apply the UID/GID to both contexts if we have values
		if uid != nil {
			// Set PodSecurityContext if RunAsUser is not set
			if params.Instance.Spec.App.PodSecurityContext.RunAsUser == nil {
				params.Instance.Spec.App.PodSecurityContext.RunAsUser = uid
				params.Instance.Spec.App.PodSecurityContext.RunAsGroup = gid
				params.Instance.Spec.App.PodSecurityContext.RunAsNonRoot = runAsNonRootPtr
			}

			// Set ContainerSecurityContext if RunAsUser is not set
			if params.Instance.Spec.App.ContainerSecurityContext.RunAsUser == nil {
				params.Instance.Spec.App.ContainerSecurityContext.RunAsUser = uid
				params.Instance.Spec.App.ContainerSecurityContext.RunAsGroup = gid
				params.Instance.Spec.App.ContainerSecurityContext.RunAsNonRoot = runAsNonRootPtr
			}

			// Set security context for init containers if not already set
			for i := range params.Instance.Spec.App.InitContainers {
				if params.Instance.Spec.App.InitContainers[i].SecurityContext == nil {
					params.Instance.Spec.App.InitContainers[i].SecurityContext = &corev1.SecurityContext{
						RunAsUser:    uid,
						RunAsGroup:   gid,
						RunAsNonRoot: runAsNonRootPtr,
					}
				}
			}

			// Set security context for sidecar containers if not already set
			for i := range params.Instance.Spec.App.Sidecars {
				if params.Instance.Spec.App.Sidecars[i].SecurityContext == nil {
					params.Instance.Spec.App.Sidecars[i].SecurityContext = &corev1.SecurityContext{
						RunAsUser:    uid,
						RunAsGroup:   gid,
						RunAsNonRoot: runAsNonRootPtr,
					}
				}
			}
		}
	}

	// Drop all capabilities by default if not explicitly set
	if params.Instance.Spec.App.ContainerSecurityContext.Capabilities == nil {
		params.Instance.Spec.App.ContainerSecurityContext.Capabilities = &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		}
	}
}

func Deployment(ctx context.Context, params Params) error {
	// Auto-detect and set OpenShift UID/GID if on OpenShift platform
	// and user hasn't explicitly set RunAsUser
	if params.Platform == "openshift" {
		applyOpenShiftSecurityDefaults(ctx, params)
	}

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
		desiredDeployment, err = setStateStoreConfig(ctx, params, desiredDeployment)
		if err != nil {
			return err
		}

		desiredDeployment, err = setGmanInitContainerVolumeMounts(ctx, params, desiredDeployment)
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

	desiredDeployment, err = setLabels(ctx, params, desiredDeployment)
	if err != nil {
		return err
	}
	desiredDeployment, err = setStateStoreConfig(ctx, params, desiredDeployment)
	if err != nil {
		return err
	}
	desiredDeployment, err = setGmanInitContainerVolumeMounts(ctx, params, desiredDeployment)
	if err != nil {
		return err
	}

	// Start with current deployment (preserves API server defaults) and merge in desired changes
	updatedDeployment := currentDeployment.DeepCopy()
	updatedDeployment.Spec = desiredDeployment.Spec
	updatedDeployment.ObjectMeta.OwnerReferences = desiredDeployment.ObjectMeta.OwnerReferences

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

	// Check if there are actual changes
	patchData, err := patch.Data(updatedDeployment)
	if err != nil {
		return fmt.Errorf("failed to generate patch: %w", err)
	}

	// Skip empty patches
	if string(patchData) == "{}" || string(patchData) == "{\"metadata\":{\"creationTimestamp\":null}}" {
		params.Log.V(2).Info("no deployment changes detected, skipping patch",
			"name", params.Instance.Name,
			"namespace", params.Instance.Namespace)
		return nil
	}

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
		if params.Instance.Spec.App.Management.SecretName == "" {
			if params.Instance.Spec.App.Management.DisklessConfig.Disabled {
				secrets = append(secrets, params.Instance.Name+"-node-properties")
			} else {
				secrets = append(secrets, params.Instance.Name)
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
	commits := ""
	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		for _, repoStatus := range params.Instance.Status.RepositoryStatus {
			if repoRef.Name == repoStatus.Name && repoRef.Type == securityv1.RepositoryReferenceTypeStatic && repoRef.Enabled {
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

func setGmanInitContainerVolumeMounts(ctx context.Context, params Params, dep *appsv1.Deployment) (*appsv1.Deployment, error) {

	var (
		gw                            = params.Instance
		defaultMode                   = int32(444)
		optional                      = false
		gmanInitContainerVolumeMounts = []corev1.VolumeMount{}
		repoRefStatuses               = []string{}
	)

	for _, repoRef := range gw.Status.RepositoryStatus {
		if repoRef.Type == "static" || (gw.Spec.App.RepositoryReferenceBootstrap.Enabled && !gw.Spec.App.Management.Database.Enabled) {
			repoRefStatuses = append(repoRefStatuses, repoRef.Name)
			// if the repository compressed is less than 1mb in size it will be
			// available as an existing Kubernetes secret which reduces reliance on an external Git repository for Gateway boot.
			// these secrets are managed by the Repository controller.
			// if the storageSecret is available we don't need to mount the secrets.
			if repoRef.StorageSecretName != "" && repoRef.StorageSecretName != "_" {
				gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
					Name:      repoRef.StorageSecretName,
					MountPath: "/graphman/localref/" + repoRef.StorageSecretName,
				})

				existingVolume := false
				for _, v := range dep.Spec.Template.Spec.Volumes {
					if v.Name == repoRef.StorageSecretName {
						existingVolume = true
					}
				}

				if !existingVolume {
					dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
						Name: repoRef.StorageSecretName,
						VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
							SecretName:  repoRef.StorageSecretName,
							DefaultMode: &defaultMode,
							Optional:    &optional,
						}},
					})
				}
			} else {
				if repoRef.SecretName != "" {
					gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
						Name:      repoRef.SecretName,
						MountPath: "/graphman/secrets/" + repoRef.Name,
					})

					existingVolume := false
					for _, v := range dep.Spec.Template.Spec.Volumes {
						if v.Name == repoRef.SecretName {
							existingVolume = true
						}
					}
					if !existingVolume {
						dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
							Name: repoRef.SecretName,
							VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
								SecretName:  repoRef.SecretName,
								DefaultMode: &defaultMode,
								Optional:    &optional,
							}},
						})
					}
				}
			}
		}
	}

	// look at repoRefs in spec. If in spec but not in status then add or replace volumes
	for _, repoRefSpec := range gw.Spec.App.RepositoryReferences {
		found := false
		for _, repoRefStatus := range repoRefStatuses {
			if repoRefStatus == repoRefSpec.Name {
				found = true

			}
		}
		if !found {
			repo := securityv1.Repository{}
			err := params.Client.Get(ctx, types.NamespacedName{Name: repoRefSpec.Name, Namespace: params.Instance.Namespace}, &repo)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve repository: %s", repoRefSpec.Name)
			}

			if repoRefSpec.Type == "static" || (gw.Spec.App.RepositoryReferenceBootstrap.Enabled && !gw.Spec.App.Management.Database.Enabled) {
				repoRefStatuses = append(repoRefStatuses, repoRefSpec.Name)
				// if the repository compressed is less than 1mb in size it will be
				// available as an existing Kubernetes secret which reduces reliance on an external Git repository for Gateway boot.
				// these secrets are managed by the Repository controller.
				// if the storageSecret is available we don't need to mount the secrets.
				if repo.Status.StorageSecretName != "_" {
					gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
						Name:      repo.Status.StorageSecretName,
						MountPath: "/graphman/localref/" + repo.Status.StorageSecretName,
					})

					existingVolume := false
					for _, v := range dep.Spec.Template.Spec.Volumes {
						if v.Name == repo.Status.StorageSecretName {
							existingVolume = true
						}
					}

					if !existingVolume {
						dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
							Name: repo.Status.StorageSecretName,
							VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
								SecretName:  repo.Status.StorageSecretName,
								DefaultMode: &defaultMode,
								Optional:    &optional,
							}},
						})
					}
				} else {

					secretName := repo.Name
					if repo.Spec.Auth.ExistingSecretName != "" {
						secretName = repo.Spec.Auth.ExistingSecretName
					}

					if repo.Spec.Auth == (securityv1.RepositoryAuth{}) {
						secretName = ""
					}

					if secretName != "" {
						gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
							Name:      secretName,
							MountPath: "/graphman/secrets/" + repoRefSpec.Name,
						})

						existingVolume := false
						for _, v := range dep.Spec.Template.Spec.Volumes {
							if v.Name == secretName {
								existingVolume = true
							}
						}
						if !existingVolume {
							dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
								Name: secretName,
								VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
									SecretName:  secretName,
									DefaultMode: &defaultMode,
									Optional:    &optional,
								}},
							})
						}
					}
				}
			}
		}
	}

	for index, ic := range dep.Spec.Template.Spec.InitContainers {
		if ic.Name == "graphman-static-init" {

			for _, nvm := range gmanInitContainerVolumeMounts {
				found := false
				for _, vm := range dep.Spec.Template.Spec.InitContainers[index].VolumeMounts {
					if nvm.Name == vm.Name {
						found = true
					}
				}
				if !found {
					dep.Spec.Template.Spec.InitContainers[index].VolumeMounts = append(dep.Spec.Template.Spec.InitContainers[index].VolumeMounts, nvm)
				}
			}

			// Sort volume mounts to ensure consistent ordering
			sort.Slice(dep.Spec.Template.Spec.InitContainers[index].VolumeMounts, func(i, j int) bool {
				return dep.Spec.Template.Spec.InitContainers[index].VolumeMounts[i].Name <
					dep.Spec.Template.Spec.InitContainers[index].VolumeMounts[j].Name
			})
		}
	}

	// Sort volumes to ensure consistent ordering across reconciliations
	sort.Slice(dep.Spec.Template.Spec.Volumes, func(i, j int) bool {
		return dep.Spec.Template.Spec.Volumes[i].Name < dep.Spec.Template.Spec.Volumes[j].Name
	})

	// Also sort volume mounts on the main gateway container
	for i, container := range dep.Spec.Template.Spec.Containers {
		if container.Name == "gateway" {
			sort.Slice(dep.Spec.Template.Spec.Containers[i].VolumeMounts, func(j, k int) bool {
				return dep.Spec.Template.Spec.Containers[i].VolumeMounts[j].Name <
					dep.Spec.Template.Spec.Containers[i].VolumeMounts[k].Name
			})
			break
		}
	}

	return dep, nil
}
