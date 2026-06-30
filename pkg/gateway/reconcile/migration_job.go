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
	"errors"
	"fmt"
	"reflect"

	"github.com/caapim/layer7-operator/pkg/gateway"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ErrMigrationPending is returned by GatewayMigrationJob when the migration job
// exists and has not yet completed (created, replaced, or still active). The
// controller treats this as a normal "not ready yet" state: it requeues with a
// fixed short interval rather than applying exponential backoff or incrementing
// failure metrics.
var ErrMigrationPending = errors.New("migration pending")

func GatewayMigrationJob(ctx context.Context, params Params) error {
	migrationEnabled := params.Instance.Spec.App.Management.Database.MigrationJob.Enabled
	databaseEnabled := params.Instance.Spec.App.Management.Database.Enabled

	// Fix #2/#8: when the migration job is disabled (or database is disabled),
	// clean up any leftover job so that toggling enabled→false doesn't orphan
	// a completed or failed Job in the namespace.
	if !migrationEnabled || !databaseEnabled {
		existingJob := &batchv1.Job{}
		jobName := params.Instance.Name + "-db-migration"
		if err := params.Client.Get(ctx, types.NamespacedName{
			Name:      jobName,
			Namespace: params.Instance.Namespace,
		}, existingJob); err == nil {
			propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationForeground)
			_ = params.Client.Delete(ctx, existingJob, propagationPolicy)
		}
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
		return fmt.Errorf("%w: migration job started, waiting for completion", ErrMigrationPending)
	} else if err != nil {
		return err
	}

	// Replace the job if any migration-relevant spec field changed.
	// This covers image, env vars (jdbcUrl, clearLocks mode), EnvFrom (secretName
	// rotation), Volumes (non-diskless secret path), and activeDeadlineSeconds.
	if migrationJobSpecChanged(currentJob, desiredJob) {
		params.Log.Info("migration job spec changed, replacing", "name", currentJob.Name)
		propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationForeground)
		if err := params.Client.Delete(ctx, currentJob, propagationPolicy); err != nil {
			return fmt.Errorf("failed deleting stale migration job: %w", err)
		}
		return fmt.Errorf("%w: deleted stale migration job, requeueing to create new one", ErrMigrationPending)
	}

	if currentJob.Status.Succeeded > 0 {
		params.Log.V(2).Info("migration job completed successfully")
		return nil
	}

	// Only declare terminal failure when no pod attempt is still active.
	// With backoffLimit=1, Kubernetes briefly shows Active=1, Failed=1 while
	// pod-2 (the retry) is running. Declaring failure at that point would kill
	// a healthy retry and require unnecessary manual intervention.
	if currentJob.Status.Active == 0 && currentJob.Status.Failed > 0 {
		params.Log.Info("migration job failed — Gateway deployment is blocked. Investigate logs then delete the job to retry",
			"name", currentJob.Name,
			"namespace", currentJob.Namespace,
			"fix", "kubectl delete job "+currentJob.Name+" -n "+currentJob.Namespace)
		return fmt.Errorf("migration job %s/%s failed — Gateway deployment blocked pending manual intervention", currentJob.Namespace, currentJob.Name)
	}

	// Job is still active (or briefly between pod attempts). Return ErrMigrationPending
	// so the controller requeues at a fixed interval without incrementing failure metrics.
	params.Log.Info("migration job is active, waiting...")
	return fmt.Errorf("%w: migration job is active", ErrMigrationPending)
}

// migrationJobSpecChanged returns true if any migration-relevant field differs
// between the current (in-cluster) job and what the operator would create now.
// Compares image, env vars, EnvFrom (covers secretName changes), Volumes (covers
// non-diskless secret path changes), and activeDeadlineSeconds.
func migrationJobSpecChanged(current, desired *batchv1.Job) bool {
	// A zero-container current job is malformed — force replacement so
	// it doesn't permanently block the Deployment step.
	if len(current.Spec.Template.Spec.Containers) == 0 {
		return true
	}
	// desired is always built by NewMigrationJob and always has one container.
	if len(desired.Spec.Template.Spec.Containers) == 0 {
		return false
	}
	cc := current.Spec.Template.Spec.Containers[0]
	dc := desired.Spec.Template.Spec.Containers[0]
	if cc.Image != dc.Image {
		return true
	}
	if envVarsChanged(cc.Env, dc.Env) {
		return true
	}
	// Compare EnvFrom so a secretName rotation triggers job replacement.
	if !reflect.DeepEqual(cc.EnvFrom, dc.EnvFrom) {
		return true
	}
	// Compare Volumes so non-diskless secret path changes are detected.
	if !reflect.DeepEqual(current.Spec.Template.Spec.Volumes, desired.Spec.Template.Spec.Volumes) {
		return true
	}
	// Detect activeDeadlineSeconds changes so users can tune slow migrations
	// without having to manually delete the failed job.
	if !int64PtrEqual(
		current.Spec.Template.Spec.ActiveDeadlineSeconds,
		desired.Spec.Template.Spec.ActiveDeadlineSeconds,
	) {
		return true
	}
	return false
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

// int64PtrEqual returns true if both pointers are nil or both point to equal values.
func int64PtrEqual(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
