package gateway

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"

	"github.com/caapim/layer7-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewDeployment(gw *securityv1.Gateway) *appsv1.Deployment {
	var image string = gw.Spec.App.Image
	defaultMode := int32(0755)
	optional := false
	ports := []corev1.ContainerPort{}

	for p := range gw.Spec.App.Service.Ports {
		ports = append(ports, corev1.ContainerPort{
			Name:          gw.Spec.App.Service.Ports[p].Name,
			ContainerPort: gw.Spec.App.Service.Ports[p].Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	if gw.Spec.App.Management.Service.Enabled {
		for p := range gw.Spec.App.Management.Service.Ports {
			ports = append(ports, corev1.ContainerPort{
				Name:          gw.Spec.App.Management.Service.Ports[p].Name,
				ContainerPort: gw.Spec.App.Management.Service.Ports[p].Port,
				Protocol:      corev1.ProtocolTCP,
			})
		}
	}

	secretName := gw.Name
	if gw.Spec.App.Management.SecretName != "" {
		secretName = gw.Spec.App.Management.SecretName
	}

	livenessProbe := corev1.Probe{

		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/bash", "/opt/docker/rc.d/diagnostic/health_check.sh"},
			},
		},
		InitialDelaySeconds: 45,
		TimeoutSeconds:      1,
		PeriodSeconds:       15,
		FailureThreshold:    25,
		SuccessThreshold:    1,
	}

	readinessProbe := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/bash", "/opt/docker/rc.d/diagnostic/health_check.sh"},
			},
		},
		InitialDelaySeconds: 45,
		TimeoutSeconds:      1,
		PeriodSeconds:       15,
		FailureThreshold:    25,
		SuccessThreshold:    1,
	}

	if gw.Spec.App.LivenessProbe != (corev1.Probe{}) {
		livenessProbe = gw.Spec.App.LivenessProbe
	}

	if gw.Spec.App.ReadinessProbe != (corev1.Probe{}) {
		readinessProbe = gw.Spec.App.ReadinessProbe
	}

	terminationGracePeriodSeconds := int64(30)
	if gw.Spec.App.TerminationGracePeriodSeconds != 0 {
		terminationGracePeriodSeconds = gw.Spec.App.TerminationGracePeriodSeconds
	}

	// As in the Gateway Helm Chart, if lifecycle hooks are defined they take precendence over the
	// pre stop script. Termination grace period seconds is automatically set to timeoutSeconds + 30
	lifecycleHooks := corev1.Lifecycle{}

	if gw.Spec.App.PreStopScript.Enabled {
		lifecycleHooks = corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"/bin/bash", "/opt/docker/graceful-shutdown.sh", strconv.Itoa(gw.Spec.App.PreStopScript.TimeoutSeconds), strconv.Itoa(gw.Spec.App.PreStopScript.PeriodSeconds)},
				},
			},
		}
		for _, port := range gw.Spec.App.PreStopScript.ExcludedPorts {
			// ignore 2124 and 8777 as they are manually set
			if port != 2124 && port != 8777 {
				lifecycleHooks.PreStop.Exec.Command = append(lifecycleHooks.PreStop.Exec.Command, strconv.Itoa(port))
			}
		}

		lifecycleHooks.PreStop.Exec.Command = append(lifecycleHooks.PreStop.Exec.Command, "2124")
		lifecycleHooks.PreStop.Exec.Command = append(lifecycleHooks.PreStop.Exec.Command, "8777")

		terminationGracePeriodSeconds = int64(gw.Spec.App.PreStopScript.TimeoutSeconds) + 30
	}

	if gw.Spec.App.LifecycleHooks != (corev1.Lifecycle{}) {
		lifecycleHooks = gw.Spec.App.LifecycleHooks
		terminationGracePeriodSeconds = gw.Spec.App.TerminationGracePeriodSeconds
	}

	volumes := []corev1.Volume{{
		Name: "gateway-license",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "gateway-license",
				Items: []corev1.KeyToPath{{
					Path: "license.xml",
					Key:  "license.xml"},
				},
				DefaultMode: &defaultMode,
				Optional:    &optional,
			},
		},
	}, {
		Name: "system-properties",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-system"},
				Items: []corev1.KeyToPath{{
					Path: "system.properties",
					Key:  "system.properties"},
				},
				DefaultMode: &defaultMode,
				Optional:    &optional,
			},
		},
	}}

	volumeMounts := []corev1.VolumeMount{{
		Name:      "gateway-license",
		MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/license/license.xml",
		SubPath:   "license.xml",
	}, {
		Name:      "system-properties",
		MountPath: "/opt/SecureSpan/Gateway/node/default/etc/conf/system.properties",
		SubPath:   "system.properties",
	}}

	if gw.Spec.App.ClusterProperties.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-cwp-bundle",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/" + gw.Name + "-cwp-bundle",
		})

		vs := corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-cwp-bundle"},
			DefaultMode:          &defaultMode,
			Optional:             &optional,
		}}

		volumes = append(volumes, corev1.Volume{
			Name:         gw.Name + "-cwp-bundle",
			VolumeSource: vs,
		})
	}

	if gw.Spec.App.ListenPorts.Harden || gw.Spec.App.ListenPorts.Custom.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-listen-port-bundle",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/" + gw.Name + "-listen-port-bundle",
		})

		vs := corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-listen-port-bundle"},
			DefaultMode:          &defaultMode,
			Optional:             &optional,
		}}

		volumes = append(volumes, corev1.Volume{
			Name:         gw.Name + "-listen-port-bundle",
			VolumeSource: vs,
		})
	}

	if gw.Spec.App.Management.Restman.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "restman",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/services/restman",
		})

		volumes = append(volumes, corev1.Volume{
			Name:         "restman",
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}

	if gw.Spec.App.Management.Graphman.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "graphman",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/services/graphman",
		})
		volumes = append(volumes, corev1.Volume{
			Name:         "graphman",
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}

	if gw.Spec.App.PreStopScript.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-graceful-shutdown",
			MountPath: "/opt/docker/graceful-shutdown.sh",
			SubPath:   "graceful-shutdown.sh",
		})
		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-graceful-shutdown",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
					Items: []corev1.KeyToPath{{
						Path: "graceful-shutdown.sh",
						Key:  "graceful-shutdown"},
					},
					DefaultMode: &defaultMode,
				},
			},
		})
	}

	if gw.Spec.App.Bootstrap.Script.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-parse-custom-files-script",
			MountPath: "/opt/docker/rc.d/003-parse-custom-files.sh",
			SubPath:   "003-parse-custom-files.sh",
		})
		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-parse-custom-files-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
					Items: []corev1.KeyToPath{{
						Path: "003-parse-custom-files.sh",
						Key:  "003-parse-custom-files"},
					},
					DefaultMode: &defaultMode,
				},
			},
		})
	}

	if gw.Spec.App.Hazelcast.External {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "hazelcast-client",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/assertions/ExternalHazelcastSharedStateProviderAssertion/hazelcast-client.xml",
			SubPath:   "hazelcast-client.xml",
		})
		volumes = append(volumes, corev1.Volume{
			Name: "hazelcast-client",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
					Items: []corev1.KeyToPath{{
						Path: "hazelcast-client.xml",
						Key:  "hazelcast-client.xml"},
					},
				},
			},
		})
	}
	i := 1
	for v := range gw.Spec.App.Bundle {
		defaultMode := int32(444)
		optional := false
		switch strings.ToLower(gw.Spec.App.Bundle[v].Source) {

		case "configmap":
			baseFolder := gw.Spec.App.Bundle[v].Name
			if gw.Spec.App.Bundle[v].Type == "graphman" {
				baseFolder = "graphman/" + strconv.Itoa(i)
				i = i + 1
			}
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      gw.Spec.App.Bundle[v].Name,
				MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/" + baseFolder,
			})

			vs := corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: gw.Spec.App.Bundle[v].Name},
				DefaultMode:          &defaultMode,
				Optional:             &optional,
			}}

			volumes = append(volumes, corev1.Volume{
				Name:         gw.Spec.App.Bundle[v].Name,
				VolumeSource: vs,
			})
		case "secret":
			baseFolder := gw.Spec.App.Bundle[v].Name
			if gw.Spec.App.Bundle[v].Type == "graphman" {
				baseFolder = "graphman/" + strconv.Itoa(i)
				i = i + 1
			}
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      gw.Spec.App.Bundle[v].Name,
				MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/" + baseFolder,
			})

			if gw.Spec.App.Bundle[v].CSI == (securityv1.CSI{}) {
				volumes = append(volumes, corev1.Volume{
					Name: gw.Spec.App.Bundle[v].Name,
					VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
						SecretName:  gw.Spec.App.Bundle[v].Name,
						DefaultMode: &defaultMode,
					}},
				})
			} else {
				vs := corev1.CSIVolumeSource{
					Driver:           gw.Spec.App.Bundle[v].CSI.Driver,
					ReadOnly:         &gw.Spec.App.Bundle[v].CSI.ReadOnly,
					VolumeAttributes: map[string]string{"secretProviderClass": gw.Spec.App.Bundle[v].CSI.VolumeAttributes.SecretProviderClass},
				}
				volumes = append(volumes, corev1.Volume{
					Name:         gw.Spec.App.Bundle[v].Name,
					VolumeSource: corev1.VolumeSource{CSI: &vs},
				})
			}
		}
	}

	if gw.Spec.App.CustomConfig.Enabled {
		for v := range gw.Spec.App.CustomConfig.Mounts {
			defaultMode := int32(444)
			optional := false
			switch strings.ToLower(gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Type) {
			case "configmap":
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name,
					MountPath: gw.Spec.App.CustomConfig.Mounts[v].MountPath,
					SubPath:   gw.Spec.App.CustomConfig.Mounts[v].SubPath,
				})

				vs := corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Spec.App.CustomConfig.Mounts[v].Name},
					DefaultMode:          &defaultMode,
					Optional:             &optional,
				}}

				volumes = append(volumes, corev1.Volume{
					Name:         gw.Spec.App.CustomConfig.Mounts[v].Name,
					VolumeSource: vs,
				})
			case "secret":
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name,
					MountPath: gw.Spec.App.CustomConfig.Mounts[v].MountPath,
					SubPath:   gw.Spec.App.CustomConfig.Mounts[v].SubPath,
				})

				vs := corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
					SecretName:  gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name,
					DefaultMode: &defaultMode,
				}}

				volumes = append(volumes, corev1.Volume{
					Name:         gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name,
					VolumeSource: vs,
				})
			}
		}
	}

	for vm := range gw.Spec.App.InitContainers {
		volumeMounts = append(volumeMounts, gw.Spec.App.InitContainers[vm].VolumeMounts...)
		for v := range gw.Spec.App.InitContainers[vm].VolumeMounts {
			volumes = append(volumes, corev1.Volume{
				Name: gw.Spec.App.InitContainers[vm].VolumeMounts[v].Name,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			})
		}
	}

	for vm := range gw.Spec.App.Sidecars {
		volumeMounts = append(volumeMounts, gw.Spec.App.Sidecars[vm].VolumeMounts...)
		for v := range gw.Spec.App.Sidecars[vm].VolumeMounts {
			volumes = append(volumes, corev1.Volume{
				Name: gw.Spec.App.Sidecars[vm].VolumeMounts[v].Name,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			})
		}
	}

	strategy := appsv1.DeploymentStrategy{}

	if gw.Spec.App.UpdateStrategy != (securityv1.UpdateStrategy{}) {
		switch gw.Spec.App.UpdateStrategy.Type {
		case "rollingUpdate":
			strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
			strategy.RollingUpdate = &gw.Spec.App.UpdateStrategy.RollingUpdate
		case "recreate":
			strategy.Type = appsv1.RecreateDeploymentStrategyType
		}
	} else {
		strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
		strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
			MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
		}
	}

	containers := []corev1.Container{}
	initContainers := []corev1.Container{}
	for _, ic := range gw.Spec.App.InitContainers {
		ic.TerminationMessagePath = corev1.TerminationMessagePathDefault
		ic.TerminationMessagePolicy = corev1.TerminationMessageReadFile
		initContainers = append(initContainers, ic)
	}

	graphmanInitContainer := false
	commits := ""
	gmanInitContainerVolumeMounts := []corev1.VolumeMount{}
	for _, staticRepository := range gw.Status.RepositoryStatus {
		if staticRepository.Enabled && staticRepository.Type == "static" {
			commits = commits + staticRepository.Commit
			graphmanInitContainer = true
			gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
				Name:      staticRepository.SecretName,
				MountPath: "/graphman/secrets/" + staticRepository.Name,
			})
			volumes = append(volumes, corev1.Volume{
				Name: staticRepository.SecretName,
				VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
					SecretName:  staticRepository.SecretName,
					DefaultMode: &defaultMode,
					Optional:    &optional,
				}},
			})

			// if the repository compressed is less than 1mb in size it will be
			// available as an existing Kubernetes secret which reduces reliance on an external Git repository for Gateway boot.
			// these secrets are managed by the Repository controller.
			if staticRepository.StorageSecretName != "" {
				gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
					Name:      staticRepository.StorageSecretName,
					MountPath: "/graphman/localref/" + staticRepository.StorageSecretName,
				})
				volumes = append(volumes, corev1.Volume{
					Name: staticRepository.StorageSecretName,
					VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
						SecretName:  staticRepository.StorageSecretName,
						DefaultMode: &defaultMode,
						Optional:    &optional,
					}},
				})
			}
		}
	}
	// Config Mount
	gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
		Name:      gw.Name + "-repository-init-config",
		MountPath: "/graphman/config.json",
		SubPath:   "config.json",
	})

	volumes = append(volumes, corev1.Volume{
		Name: gw.Name + "-repository-init-config",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: gw.Name + "-repository-init-config",
			},
			DefaultMode: &defaultMode,
			Optional:    &optional,
		}},
	})

	// Target Bootstrap Mount
	gmanInitContainerVolumeMounts = append(gmanInitContainerVolumeMounts, corev1.VolumeMount{
		Name:      gw.Name + "-repository-bundle-dest",
		MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/0",
	})
	volumes = append(volumes, corev1.Volume{
		Name: gw.Name + "-repository-bundle-dest",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumeMounts = append(volumeMounts, gmanInitContainerVolumeMounts...)

	if graphmanInitContainer {
		h := sha1.New()
		h.Write([]byte(commits))
		commits = fmt.Sprintf("%x", h.Sum(nil))
		initContainers = append(initContainers, corev1.Container{
			Name:            "graphman-static-init-" + commits[30:],
			Image:           gw.Spec.App.Management.Graphman.InitContainerImage,
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts:    gmanInitContainerVolumeMounts,
			Env: []corev1.EnvVar{{
				Name:  "BOOTSTRAP_BASE",
				Value: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/0",
			}},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		})
	}

	resources := corev1.ResourceRequirements{
		Requests: gw.Spec.App.Resources.Requests,
		Limits:   gw.Spec.App.Resources.Limits,
	}

	imagePullPolicy := corev1.PullIfNotPresent

	if gw.Spec.App.ImagePullPolicy != "" {
		imagePullPolicy = gw.Spec.App.ImagePullPolicy
	}

	gateway := corev1.Container{
		Image:                    image,
		ImagePullPolicy:          imagePullPolicy,
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		SecurityContext:          &gw.Spec.App.ContainerSecurityContext,
		Name:                     "gateway",
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
				}},

			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
				}},
		},
		Ports:          ports,
		VolumeMounts:   volumeMounts,
		LivenessProbe:  &livenessProbe,
		ReadinessProbe: &readinessProbe,
		Resources:      resources,
		Lifecycle:      &lifecycleHooks,
	}

	containers = append(containers, gateway)
	containers = append(containers, gw.Spec.App.Sidecars...)

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	revisionHistoryLimit := int32(10)
	progressDeadlineSeconds := int32(600)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Replicas:                &gw.Spec.App.Replicas,
			Strategy:                strategy,
			RevisionHistoryLimit:    &revisionHistoryLimit,
			ProgressDeadlineSeconds: &progressDeadlineSeconds,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: gw.Spec.App.PodAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            gw.Spec.App.ServiceAccountName,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					SecurityContext:               &gw.Spec.App.PodSecurityContext,
					TopologySpreadConstraints:     gw.Spec.App.TopologySpreadConstraints,
					Tolerations:                   gw.Spec.App.Tolerations,
					DNSPolicy:                     corev1.DNSClusterFirst,
					RestartPolicy:                 corev1.RestartPolicyAlways,
					Affinity:                      &gw.Spec.App.Affinity,
					NodeSelector:                  gw.Spec.App.NodeSelector,
					SchedulerName:                 "default-scheduler",
					Containers:                    containers,
					Volumes:                       volumes,
				},
			},
		},
	}

	if len(gw.Spec.App.Annotations) != 0 {
		dep.ObjectMeta.Annotations = gw.Spec.App.Annotations
	}

	if len(initContainers) > 0 {
		dep.Spec.Template.Spec.InitContainers = initContainers
	}

	dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, gw.Spec.App.ImagePullSecrets...)
	dep.Spec.Template.Labels = ls

	// if gw.Spec.App.Repository.Enabled {
	// 	dep.Spec.Template.Annotations = map[string]string{"commitId": gw.Status.CommitID}
	// }

	if !gw.Spec.App.Autoscaling.Enabled {
		dep.Spec.Replicas = &gw.Spec.App.Replicas
	}

	return dep
}
