/*
* Copyright (c) 2026 Broadcom. All rights reserved.
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
package gateway

import (
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MigrationJobName returns the canonical name of the migration Job for a Gateway.
// Exported so pkg/gateway/reconcile can derive the same name without duplicating
// the "-db-migration" suffix convention.
func MigrationJobName(gw *securityv1.Gateway) string {
	return gw.Name + "-db-migration"
}

func NewMigrationJob(gw *securityv1.Gateway) *batchv1.Job {
	gatewaySecretName := gw.Name
	if gw.Spec.App.Management.DisklessConfig.Disabled {
		gatewaySecretName = gw.Name + "-node-properties"
	}
	if gw.Spec.App.Management.SecretName != "" {
		gatewaySecretName = gw.Spec.App.Management.SecretName
	}

	schemaUpdateMode := "liquibase-only"
	if gw.Spec.App.Management.Database.MigrationJob.ClearLocks {
		schemaUpdateMode = "liquibase-only-with-unlock"
	}

	// EXTRA_JAVA_ARGS is always set explicitly to control the schema update mode.
	// SSG_DATABASE_JDBC_URL and Secret mounting depend on the diskless mode — see below.
	envVars := []corev1.EnvVar{
		{
			Name:  "EXTRA_JAVA_ARGS",
			Value: "-Dgateway.db.schema-update.mode=" + schemaUpdateMode,
		},
	}

	envFrom := []corev1.EnvFromSource{
		{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: gw.Name},
			},
		},
	}

	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount

	// Secret mounting strategy mirrors the main Gateway deployment and matches Helm behavior:
	//
	// Diskless mode (disklessConfig.disabled: false, the default):
	//   The Secret is exposed as environment variables via envFrom. The entrypoint reads
	//   SSG_DATABASE_JDBC_URL directly from env, so migrationJob.jdbcUrl overrides the
	//   main database.jdbcUrl via Kubernetes env precedence (explicit Env beats envFrom).
	//
	// Non-diskless mode (disklessConfig.disabled: true):
	//   The Secret is mounted as a node.properties file at the well-known path the
	//   entrypoint expects. With DISKLESS_CONFIG=false, the entrypoint calls
	//   loadFromNodeProperties() and reads l7.mysql.connection.url from the file.
	//   migrationJob.jdbcUrl is NOT used in this mode — node.properties always wins,
	//   consistent with Helm behavior when disklessConfig is disabled.
	if !gw.Spec.App.Management.DisklessConfig.Disabled {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: gatewaySecretName},
			},
		})
		if gw.Spec.App.Management.Database.MigrationJob.JDBCUrl != "" {
			// Append createDatabaseIfNotExist=false so that a misconfigured or mistyped
			// jdbcUrl fails fast rather than silently creating an unintended database.
			// entrypoint.sh only appends createDatabaseIfNotExist=true when the parameter
			// is absent, so an explicit =false here takes precedence.
			// Use '?' to start the query string if none exists yet, '&' otherwise.
			jdbcUrl := gw.Spec.App.Management.Database.MigrationJob.JDBCUrl
			if !strings.Contains(jdbcUrl, "createDatabaseIfNotExist") {
				if strings.Contains(jdbcUrl, "?") {
					jdbcUrl += "&createDatabaseIfNotExist=false"
				} else {
					jdbcUrl += "?createDatabaseIfNotExist=false"
				}
			}
			envVars = append(envVars, corev1.EnvVar{
				Name:  "SSG_DATABASE_JDBC_URL",
				Value: jdbcUrl,
			})
		}
	} else {
		// optional must be false to match the main Gateway deployment.
		// With optional: true, Kubernetes schedules the pod even when the Secret
		// is missing, causing a silent empty mount and a confusing DB connection
		// error instead of a clear "secret not found" Kubernetes event.
		optional := false
		defaultMode := int32(0444)
		volumes = append(volumes, corev1.Volume{
			Name: gw.Name + "-node-properties",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  gatewaySecretName,
					DefaultMode: &defaultMode,
					Optional:    &optional,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      gw.Name + "-node-properties",
			MountPath: "/opt/SecureSpan/Gateway/node/default/etc/conf/node.properties",
			SubPath:   "node.properties",
		})
	}

	container := corev1.Container{
		Name:            "db-migration",
		Image:           gw.Spec.App.Image,
		ImagePullPolicy: gw.Spec.App.ImagePullPolicy,
		EnvFrom:         envFrom,
		Env:             envVars,
		VolumeMounts:    volumeMounts,
	}

	var backoffLimit int32 = 1
	var activeDeadlineSeconds int64 = 300

	if gw.Spec.App.Management.Database.MigrationJob.ActiveDeadlineSeconds != nil {
		activeDeadlineSeconds = *gw.Spec.App.Management.Database.MigrationJob.ActiveDeadlineSeconds
	}

	// apply the same ServiceAccount fallback logic as the main deployment
	// so the migration pod runs under the same identity as Gateway pods.
	// Without this, pods default to the namespace "default" service account when
	// spec.app.serviceAccount.name is not explicitly set, which can fail in clusters
	// with restrictive default SA policies.
	serviceAccountName := gw.Spec.App.ServiceAccount.Name
	if gw.Spec.App.ServiceAccount.Name == "" {
		serviceAccountName = gw.Name
	}
	if gw.Spec.App.ServiceAccount == (securityv1.ServiceAccount{}) {
		serviceAccountName = "default"
	}

	jobName := MigrationJobName(gw)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: gw.Namespace,
			Labels:    map[string]string{"app": jobName},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": jobName},
				},
				Spec: corev1.PodSpec{
					ActiveDeadlineSeconds: &activeDeadlineSeconds,
					ServiceAccountName:    serviceAccountName,
					RestartPolicy:         corev1.RestartPolicyNever,
					Containers:            []corev1.Container{container},
					ImagePullSecrets:      gw.Spec.App.ImagePullSecrets,
					Volumes:               volumes,
				},
			},
		},
	}

	return job
}
