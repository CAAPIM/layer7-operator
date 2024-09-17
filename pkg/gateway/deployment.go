package gateway

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"

	"github.com/caapim/layer7-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewDeployment(gw *securityv1.Gateway, platform string) *appsv1.Deployment {
	var image string = gw.Spec.App.Image
	defaultMode := int32(0755)
	optional := false
	ports := []corev1.ContainerPort{}

	defaultUser := int64(1001)
	defaultGroup := int64(1001)
	runAsNonRoot := true

	ocPodSecurityContext := corev1.PodSecurityContext{
		RunAsUser:    &defaultUser,
		RunAsGroup:   &defaultGroup,
		RunAsNonRoot: &runAsNonRoot,
	}
	ocContainerSecurityContext := corev1.SecurityContext{
		RunAsUser:    &defaultUser,
		RunAsNonRoot: &runAsNonRoot,
		Capabilities: &corev1.Capabilities{Drop: []corev1.Capability{"ALL"}},
	}

	for p := range gw.Spec.App.Service.Ports {
		ports = append(ports, corev1.ContainerPort{
			Name:          gw.Spec.App.Service.Ports[p].Name,
			ContainerPort: gw.Spec.App.Service.Ports[p].TargetPort,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	if gw.Spec.App.Management.Service.Enabled {
		for p := range gw.Spec.App.Management.Service.Ports {
			ports = append(ports, corev1.ContainerPort{
				Name:          gw.Spec.App.Management.Service.Ports[p].Name,
				ContainerPort: gw.Spec.App.Management.Service.Ports[p].TargetPort,
				Protocol:      corev1.ProtocolTCP,
			})
		}
	}

	livenessProbe := corev1.Probe{

		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"/bin/bash", "/opt/docker/rc.d/diagnostic/health_check.sh"},
			},
		},
		InitialDelaySeconds: 30,
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
		InitialDelaySeconds: 30,
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
				SecretName: gw.Spec.License.SecretName,
				Items: []corev1.KeyToPath{{
					Path: "license.xml",
					Key:  "license.xml"},
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
	}}

	if gw.Spec.App.System.Properties != "" {
		volumes = append(volumes, corev1.Volume{
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
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "system-properties",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/conf/system.properties",
			SubPath:   "system.properties",
		})
	}

	if gw.Spec.App.AutoMountServiceAccountToken {
		volumes = append(volumes, corev1.Volume{
			Name: "service-account-token-script",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
					Items: []corev1.KeyToPath{{
						Path: "load-service-account-token.sh",
						Key:  "load-service-account-token"},
					},
					DefaultMode: &defaultMode,
					Optional:    &optional,
				},
			},
		}, corev1.Volume{
			Name: "service-account-token-template",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
					Items: []corev1.KeyToPath{{
						Path: "update-service-account-token.xml",
						Key:  "service-account-token-template"},
					},
					DefaultMode: &defaultMode,
					Optional:    &optional,
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "service-account-token-script",
			MountPath: "/opt/docker/rc.d/load-service-account-token.sh",
			SubPath:   "load-service-account-token.sh",
		}, corev1.VolumeMount{
			Name:      "service-account-token-template",
			MountPath: "/opt/docker/rc.d/base/update-service-account-token.xml",
			SubPath:   "update-service-account-token.xml",
		})
	}

	if gw.Spec.App.Redis.Enabled {
		secretName := gw.Name + "-shared-state-client-configuration"
		if gw.Spec.App.Redis.ExistingSecret != "" {
			secretName = gw.Spec.App.Redis.ExistingSecret
		}

		items := []corev1.KeyToPath{{Key: "sharedstate_client.yaml", Path: "sharedstate_client.yaml"}}

		if gw.Spec.App.Redis.ExistingSecret == "" {
			if gw.Spec.App.Redis.Default.Ssl.Enabled {
				certSecretName := secretName
				if gw.Spec.App.Redis.Default.Ssl.ExistingSecretName != "" {
					certSecretName = gw.Spec.App.Redis.Default.Ssl.ExistingSecretName
				}
				key := "redis.crt"
				if gw.Spec.App.Redis.Default.Ssl.ExistingSecretKey != "" {
					key = gw.Spec.App.Redis.Default.Ssl.ExistingSecretKey
				}

				volumes = append(volumes, corev1.Volume{
					Name: "default-redis-ssl",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: certSecretName,
							Optional:   &optional,
							Items: []corev1.KeyToPath{{
								Key:  key,
								Path: "redis.crt",
							}},
						},
					},
				})
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      "default-redis-ssl",
					MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/providers/redis.crt",
					SubPath:   "redis.crt",
				})

			}

			for _, ac := range gw.Spec.App.Redis.AdditionalConfigs {
				if ac.Enabled && ac.Ssl.Enabled {
					certSecretName := secretName
					if ac.Ssl.ExistingSecretKey != "" {
						certSecretName = ac.Ssl.ExistingSecretName
					}
					key := ac.Name + "-redis.crt"
					if ac.Ssl.ExistingSecretKey != "" {
						key = ac.Ssl.ExistingSecretKey
					}
					volumes = append(volumes, corev1.Volume{
						Name: ac.Name + "redis-ssl",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: certSecretName,
								Optional:   &optional,
								Items: []corev1.KeyToPath{{
									Key:  key,
									Path: ac.Name + "-redis.crt",
								}},
							},
						},
					})
					volumeMounts = append(volumeMounts, corev1.VolumeMount{
						Name:      ac.Name + "redis-ssl",
						MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/providers/" + ac.Name + "-redis.crt",
						SubPath:   ac.Name + "-redis.crt",
					})
				}
			}

		} else {
			for i, certSecret := range gw.Spec.App.Redis.CertSecrets {
				if certSecret.Enabled {
					volumes = append(volumes, corev1.Volume{
						Name: "redis-ssl-" + strconv.Itoa(i),
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: certSecret.SecretName,
								Optional:   &optional,
								Items: []corev1.KeyToPath{{
									Key:  certSecret.Key,
									Path: certSecret.Key,
								}},
							},
						},
					})
					volumeMounts = append(volumeMounts, corev1.VolumeMount{
						Name:      "redis-ssl-" + strconv.Itoa(i),
						MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/providers/" + certSecret.Key,
						SubPath:   certSecret.Key,
					})
				}
			}
		}

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "sharedstate-client-config",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/providers/sharedstate_client.yaml",
			SubPath:   "sharedstate_client.yaml",
		})

		volumes = append(volumes, corev1.Volume{
			Name: "sharedstate-client-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
					Optional:   &optional,
					Items:      items,
				},
			},
		})

	}

	if gw.Spec.App.Log.Override {
		volumes = append(volumes, corev1.Volume{
			Name: "log-override-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
					Items: []corev1.KeyToPath{{
						Path: "log-override.properties",
						Key:  "log-override-properties"},
					},
					DefaultMode: &defaultMode,
					Optional:    &optional,
				},
			},
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "log-override-config",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/conf/log-override.properties",
			SubPath:   "log-override.properties",
		})

	}

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
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
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
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
					Items: []corev1.KeyToPath{{
						Path: "003-parse-custom-files.sh",
						Key:  "003-parse-custom-files"},
					},
					DefaultMode: &defaultMode,
					Optional:    &optional,
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
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name + "-gateway-files"},
					Items: []corev1.KeyToPath{{
						Path: "hazelcast-client.xml",
						Key:  "hazelcast-client.xml"},
					},
				},
			},
		})
	}
	i := 2
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
			if reflect.DeepEqual(gw.Spec.App.Bundle[v].CSI, securityv1.CSI{}) {
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
					VolumeAttributes: gw.Spec.App.Bundle[v].CSI.VolumeAttributes,
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
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name},
					DefaultMode:          &defaultMode,
					Optional:             &optional,
				}}

				volumes = append(volumes, corev1.Volume{
					Name:         gw.Spec.App.CustomConfig.Mounts[v].ConfigRef.Name,
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
		if platform == "openshift" && ic.SecurityContext == nil {
			ic.SecurityContext = &ocContainerSecurityContext
			if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
				ic.SecurityContext = &gw.Spec.App.ContainerSecurityContext
			}
		}
		if ic.ImagePullPolicy == "" {
			ic.ImagePullPolicy = corev1.PullIfNotPresent
		}
		initContainers = append(initContainers, ic)
	}

	graphmanInitContainer := false
	commits := ""
	gmanInitContainerVolumeMounts := []corev1.VolumeMount{}
	for _, staticRepository := range gw.Status.RepositoryStatus {
		if staticRepository.Enabled && staticRepository.Type == "static" {
			commits = commits + staticRepository.Commit
			graphmanInitContainer = true

			if staticRepository.SecretName != "" {
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
			}

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

	if graphmanInitContainer {
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

		h := sha1.New()
		h.Write([]byte(commits))
		commits = fmt.Sprintf("%x", h.Sum(nil))

		graphmanInitContainerImage := "docker.io/caapim/graphman-static-init:1.0.2"
		graphmanInitContainerImagePullPolicy := corev1.PullIfNotPresent
		graphmanInitContainerSecurityContext := corev1.SecurityContext{}

		if gw.Spec.App.Management.Graphman.InitContainerImage != "" {
			graphmanInitContainerImage = gw.Spec.App.Management.Graphman.InitContainerImage
		}

		if gw.Spec.App.Management.Graphman.InitContainerImagePullPolicy != "" {
			graphmanInitContainerImagePullPolicy = gw.Spec.App.Management.Graphman.InitContainerImagePullPolicy
		}

		if platform == "openshift" {
			graphmanInitContainerSecurityContext = ocContainerSecurityContext
		}

		if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
			graphmanInitContainerSecurityContext = gw.Spec.App.ContainerSecurityContext
		}

		if gw.Spec.App.Management.Graphman.InitContainerSecurityContext != (corev1.SecurityContext{}) {
			graphmanInitContainerSecurityContext = gw.Spec.App.Management.Graphman.InitContainerSecurityContext
		}

		initContainers = append(initContainers, corev1.Container{
			Name:            "graphman-static-init-" + commits[30:],
			Image:           graphmanInitContainerImage,
			ImagePullPolicy: graphmanInitContainerImagePullPolicy,
			SecurityContext: &graphmanInitContainerSecurityContext,
			VolumeMounts:    gmanInitContainerVolumeMounts,
			Env: []corev1.EnvVar{{
				Name:  "BOOTSTRAP_BASE",
				Value: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/0",
			}},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		})
	}

	if gw.Spec.App.PortalReference.Enabled {
		portalInitContainerVolumeMounts := []corev1.VolumeMount{}
		portalInitContainerVolumeMounts = append(portalInitContainerVolumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-portal-init-config",
			MountPath: "/portal/config.json",
			SubPath:   "config.json",
		})

		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-portal-init-config",
			VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: gw.Spec.App.PortalReference.PortalName + "-api-summary",
				},
				DefaultMode: &defaultMode,
				Optional:    &optional,
				Items: []corev1.KeyToPath{{
					Key:  "apis",
					Path: "config.json",
				}},
			}},
		})

		portalInitContainerVolumeMounts = append(portalInitContainerVolumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-portal-init-dest",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/1",
		})
		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-portal-init-dest",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})

		volumeMounts = append(volumeMounts, portalInitContainerVolumeMounts...)

		portalInitContainerImage := "docker.io/layer7api/portal-bulk-init:0.0.1"
		portalInitContainerImagePullPolicy := corev1.PullIfNotPresent
		portalInitContainerSecurityContext := corev1.SecurityContext{}

		if gw.Spec.App.PortalReference.InitContainerImage != "" {
			portalInitContainerImage = gw.Spec.App.PortalReference.InitContainerImage
		}

		if gw.Spec.App.PortalReference.InitContainerImagePullPolicy != "" {
			portalInitContainerImagePullPolicy = gw.Spec.App.PortalReference.InitContainerImagePullPolicy
		}

		if platform == "openshift" {
			portalInitContainerSecurityContext = ocContainerSecurityContext
		}

		if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
			portalInitContainerSecurityContext = gw.Spec.App.ContainerSecurityContext
		}

		if gw.Spec.App.PortalReference.InitContainerSecurityContext != (corev1.SecurityContext{}) {
			portalInitContainerSecurityContext = gw.Spec.App.PortalReference.InitContainerSecurityContext
		}
		initContainers = append(initContainers, corev1.Container{
			Name:            "portal-init",
			Image:           portalInitContainerImage,
			ImagePullPolicy: portalInitContainerImagePullPolicy,
			SecurityContext: &portalInitContainerSecurityContext,
			VolumeMounts:    portalInitContainerVolumeMounts,
			Env: []corev1.EnvVar{{
				Name:  "BOOTSTRAP_BASE",
				Value: "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/1",
			}},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		})
	}

	otkInstallInitContainer := false
	otkDbInitContainer := false
	otkBootstrapDirectory := "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/000OTK"
	otkInitContainerVolumeMounts := []corev1.VolumeMount{}
	otkInitContainerImage := "docker.io/caapim/otk-install:4.6.2_202402"
	otkInitContainerImagePullPolicy := corev1.PullIfNotPresent
	otkInitContainerSecurityContext := corev1.SecurityContext{}

	otkInitContainerSecret := gw.Name + "-otk-db-credentials"

	if gw.Spec.App.Otk.Database.Auth.ExistingSecret != "" {
		otkInitContainerSecret = gw.Spec.App.Otk.Database.Auth.ExistingSecret
	}

	if gw.Spec.App.Otk.InitContainerImage != "" {
		otkInitContainerImage = gw.Spec.App.Otk.InitContainerImage
	}

	if gw.Spec.App.Otk.InitContainerImagePullPolicy != "" {
		otkInitContainerImagePullPolicy = gw.Spec.App.Otk.InitContainerImagePullPolicy
	}

	if platform == "openshift" {
		otkInitContainerSecurityContext = ocContainerSecurityContext
	}

	if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
		otkInitContainerSecurityContext = gw.Spec.App.ContainerSecurityContext
	}

	if gw.Spec.App.Otk.InitContainerSecurityContext != (corev1.SecurityContext{}) {
		otkInitContainerSecurityContext = gw.Spec.App.Otk.InitContainerSecurityContext
	}

	if gw.Spec.App.Otk.Overrides.Enabled {
		if gw.Spec.App.Otk.Overrides.BootstrapDirectory != "" {
			otkBootstrapDirectory = gw.Spec.App.Otk.Overrides.BootstrapDirectory
		}
	}

	if gw.Spec.App.Otk.Enabled {
		otkInstallInitContainer = true
		otkInitContainerVolumeMounts = append(otkInitContainerVolumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-otk-bundle-dest",
			MountPath: otkBootstrapDirectory,
		})

		volumeMounts = append(volumeMounts, otkInitContainerVolumeMounts...)
		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-otk-bundle-dest",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		if gw.Spec.App.Otk.Database.Type == securityv1.OtkDatabaseTypeMySQL || gw.Spec.App.Otk.Database.Type == securityv1.OtkDatabaseTypeOracle {
			if gw.Spec.App.Otk.Database.Create {
				otkDbInitContainer = true
			}
		}
	}

	if otkInstallInitContainer {
		initContainers = append(initContainers, corev1.Container{
			Name:            "otk-install-init",
			Image:           otkInitContainerImage,
			ImagePullPolicy: otkInitContainerImagePullPolicy,
			SecurityContext: &otkInitContainerSecurityContext,
			VolumeMounts:    otkInitContainerVolumeMounts,
			EnvFrom: []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: gw.Name + "-otk-shared-init-config",
						},
					},
				},
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: gw.Name + "-otk-install-init-config",
						},
						Optional: &optional,
					},
				},
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: otkInitContainerSecret,
						},
						Optional: &optional,
					},
				},
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		})
	}

	if otkDbInitContainer && (gw.Spec.App.Otk.Type == securityv1.OtkTypeInternal || gw.Spec.App.Otk.Type == securityv1.OtkTypeSingle) {
		initContainers = append(initContainers, corev1.Container{
			Name:            "otk-db-init",
			Image:           otkInitContainerImage,
			ImagePullPolicy: otkInitContainerImagePullPolicy,
			SecurityContext: &otkInitContainerSecurityContext,
			EnvFrom: []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: gw.Name + "-otk-db-init-config",
						},
						Optional: &optional,
					},
				},
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: otkInitContainerSecret,
						},
						Optional: &optional,
					},
				},
			},
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

	gatewayContainerSecurityContext := corev1.SecurityContext{}
	podSecurityContext := corev1.PodSecurityContext{}

	if platform == "openshift" {
		gatewayContainerSecurityContext = ocContainerSecurityContext
		podSecurityContext = ocPodSecurityContext
	}

	if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
		gatewayContainerSecurityContext = gw.Spec.App.ContainerSecurityContext
	}

	if !reflect.DeepEqual(gw.Spec.App.PodSecurityContext, corev1.PodSecurityContext{}) {
		podSecurityContext = gw.Spec.App.PodSecurityContext
	}
	gatewaySecretName := gw.Name
	if gw.Spec.App.Management.DisklessConfig.Disabled {
		gatewaySecretName = gw.Name + "-node-properties"
	}
	if gw.Spec.App.Management.SecretName != "" {
		gatewaySecretName = gw.Spec.App.Management.SecretName
	}

	gateway := corev1.Container{
		Image:                    image,
		ImagePullPolicy:          imagePullPolicy,
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		SecurityContext:          &gatewayContainerSecurityContext,
		Name:                     "gateway",
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
				}},
		},
		Ports:          ports,
		LivenessProbe:  &livenessProbe,
		ReadinessProbe: &readinessProbe,
		Resources:      resources,
		Lifecycle:      &lifecycleHooks,
	}

	secretRef := corev1.SecretEnvSource{}
	if !gw.Spec.App.Management.DisklessConfig.Disabled {
		secretRef = corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: gatewaySecretName},
		}
		gateway.EnvFrom = append(gateway.EnvFrom, corev1.EnvFromSource{SecretRef: &secretRef})
	} else {
		vs := corev1.VolumeSource{}

		if !reflect.DeepEqual(gw.Spec.App.Management.DisklessConfig.Csi, securityv1.CSI{}) {

			vs = corev1.VolumeSource{CSI: &corev1.CSIVolumeSource{
				Driver:           gw.Spec.App.Management.DisklessConfig.Csi.Driver,
				ReadOnly:         &gw.Spec.App.Management.DisklessConfig.Csi.ReadOnly,
				VolumeAttributes: gw.Spec.App.Management.DisklessConfig.Csi.VolumeAttributes,
			}}
		} else {
			vs = corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{
				SecretName:  gatewaySecretName,
				DefaultMode: &defaultMode,
				Optional:    &optional,
			}}
		}

		volumes = append(volumes, corev1.Volume{
			Name:         gw.Name + "-node-properties",
			VolumeSource: vs,
		})

		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-node-properties",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/conf/node.properties",
			SubPath:   "node.properties",
		})

	}

	gateway.VolumeMounts = volumeMounts

	sidecars := []corev1.Container{}

	for _, sc := range gw.Spec.App.Sidecars {
		if platform == "openshift" && sc.SecurityContext == nil {
			sc.SecurityContext = &ocContainerSecurityContext
			if gw.Spec.App.ContainerSecurityContext != (corev1.SecurityContext{}) {
				sc.SecurityContext = &gw.Spec.App.ContainerSecurityContext
			}
		}
		sidecars = append(sidecars, sc)
	}

	containers = append(containers, gateway)
	containers = append(containers, sidecars...)

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	revisionHistoryLimit := int32(10)
	progressDeadlineSeconds := int32(600)

	serviceAccountName := gw.Spec.App.ServiceAccount.Name
	if gw.Spec.App.ServiceAccount.Name == "" {
		serviceAccountName = gw.Name
	}

	if gw.Spec.App.ServiceAccount == (securityv1.ServiceAccount{}) {
		serviceAccountName = "default"
	}

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
					Labels:      gw.Spec.App.PodLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            serviceAccountName,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					SecurityContext:               &podSecurityContext,
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

	if gw.Spec.App.CustomHosts.Enabled {
		if len(gw.Spec.App.CustomHosts.HostAliases) > 0 {
			dep.Spec.Template.Spec.HostAliases = gw.Spec.App.CustomHosts.HostAliases
		}
	}

	if len(gw.Spec.App.Annotations) != 0 {
		dep.ObjectMeta.Annotations = gw.Spec.App.Annotations
	}

	if len(initContainers) > 0 {
		dep.Spec.Template.Spec.InitContainers = initContainers
	}

	dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, gw.Spec.App.ImagePullSecrets...)
	dep.Spec.Template.Labels = ls

	if !gw.Spec.App.Autoscaling.Enabled {
		dep.Spec.Replicas = &gw.Spec.App.Replicas
	}

	return dep
}
