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
package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func GatewayMigrationJob(ctx context.Context, params Params) error {
	if !params.Instance.Spec.App.Management.Database.MigrationJob.Enabled {
		return nil
	}

	desiredJob := gateway.NewMigrationJob(params.Instance)
	if err := controllerutil.SetControllerReference(params.Instance, desiredJob, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on migration job: %w", err)
	}

	currentJob := &batchv1.Job{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredJob.Name, Namespace: desiredJob.Namespace}, currentJob)

	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredJob); err != nil {
			return fmt.Errorf("failed creating migration job: %w", err)
		}
		params.Log.Info("created migration job", "name", desiredJob.Name, "namespace", desiredJob.Namespace)
		return fmt.Errorf("migration job started, waiting for completion")
	} else if err != nil {
		return err
	}

	// Replace the job if any migration-relevant spec field changed: image, jdbcUrl,
	// clearLocks, or any other env var. This ensures that applying an updated CR
	// (e.g. fixing a wrong jdbcUrl or upgrading the image) always triggers a fresh run.
	if migrationJobSpecChanged(currentJob, desiredJob) {
		params.Log.Info("migration job spec changed, replacing", "name", currentJob.Name)
		propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationForeground)
		if err := params.Client.Delete(ctx, currentJob, propagationPolicy); err != nil {
			return fmt.Errorf("failed deleting stale migration job: %w", err)
		}
		return fmt.Errorf("deleted stale migration job, requeueing to create new one")
	}

	if currentJob.Status.Succeeded > 0 {
		params.Log.V(2).Info("migration job completed successfully")
		return nil
	}

	if currentJob.Status.Failed > 0 {
		params.Log.Info("migration job failed — Gateway deployment is blocked. Investigate logs then delete the job to retry",
			"name", currentJob.Name,
			"namespace", currentJob.Namespace,
			"fix", "kubectl delete job "+currentJob.Name+" -n "+currentJob.Namespace)
		return fmt.Errorf("migration job %s/%s failed — Gateway deployment blocked pending manual intervention", currentJob.Namespace, currentJob.Name)
	}

	params.Log.Info("migration job is active, waiting...")
	return fmt.Errorf("migration job is active, waiting...")
}

// migrationJobSpecChanged returns true if any migration-relevant container fields
// differ between the current (in-cluster) job and what the operator would create now.
// Comparing env vars by name→value catches jdbcUrl, clearLocks, and mode changes
// in addition to image upgrades.
func migrationJobSpecChanged(current, desired *batchv1.Job) bool {
	if len(current.Spec.Template.Spec.Containers) == 0 || len(desired.Spec.Template.Spec.Containers) == 0 {
		return false
	}
	cc := current.Spec.Template.Spec.Containers[0]
	dc := desired.Spec.Template.Spec.Containers[0]
	if cc.Image != dc.Image {
		return true
	}
	return envVarsChanged(cc.Env, dc.Env)
}

// envVarsChanged compares two env var slices by name→value, ignoring order.
func envVarsChanged(current, desired []corev1.EnvVar) bool {
	cm := make(map[string]string, len(current))
	for _, e := range current {
		cm[e.Name] = e.Value
	}
	if len(cm) != len(desired) {
		return true
	}
	for _, e := range desired {
		if cv, ok := cm[e.Name]; !ok || cv != e.Value {
			return true
		}
	}
	return false
}
