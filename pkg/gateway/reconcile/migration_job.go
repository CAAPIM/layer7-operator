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
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway"
	batchv1 "k8s.io/api/batch/v1"
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

// GatewayMigrationJob manages the lifecycle of the pre-upgrade database migration
// job. It uses Gateway.Status.MigrationStatus as the authoritative source of truth
// for whether the migration has completed, so that:
//
//   - Completed migrations are not re-run on daily reconciles or operator restarts.
//   - Completed migrations survive Job deletion/GC without triggering a re-run.
//   - A spec change (image upgrade, jdbcUrl correction, clearLocks toggle) resets
//     the status and triggers a fresh migration run automatically.
//
// The Deployment step is only unblocked when Status.MigrationStatus.Complete is true.
func GatewayMigrationJob(ctx context.Context, params Params) error {
	migrationEnabled := params.Instance.Spec.App.Management.Database.MigrationJob.Enabled
	databaseEnabled := params.Instance.Spec.App.Management.Database.Enabled

	// When migration is disabled (or database is disabled), clean up any leftover
	// job so that toggling enabled→false doesn't orphan a job in the namespace.
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

	desiredHash := migrationSpecHash(params.Instance)
	currentStatus := params.Instance.Status.MigrationStatus

	// Fast path: migration already completed for the current spec.
	// This allows daily reconciles and operator restarts to skip migration
	// management entirely — the status persists regardless of Job existence.
	if currentStatus.SpecHash == desiredHash && currentStatus.Complete {
		params.Log.V(2).Info("migration already complete for current spec, skipping",
			"hash", desiredHash)
		return nil
	}

	// Spec changed (image upgrade, jdbcUrl correction, clearLocks toggle, etc.):
	// reset the status and delete the old job so a fresh migration runs.
	if currentStatus.SpecHash != "" && currentStatus.SpecHash != desiredHash {
		params.Log.Info("migration spec changed, resetting status and replacing job",
			"oldHash", currentStatus.SpecHash,
			"newHash", desiredHash)
		if err := setMigrationStatus(ctx, params, securityv1.MigrationStatus{SpecHash: desiredHash}); err != nil {
			return err
		}
		oldJob := &batchv1.Job{}
		jobName := params.Instance.Name + "-db-migration"
		if getErr := params.Client.Get(ctx, types.NamespacedName{Name: jobName, Namespace: params.Instance.Namespace}, oldJob); getErr == nil {
			propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationForeground)
			if delErr := params.Client.Delete(ctx, oldJob, propagationPolicy); delErr != nil {
				return fmt.Errorf("failed deleting stale migration job: %w", delErr)
			}
		}
		return fmt.Errorf("%w: migration spec changed, requeueing to create new job", ErrMigrationPending)
	}

	// Record the desired hash in status if this is the first time we see it.
	if currentStatus.SpecHash == "" {
		if err := setMigrationStatus(ctx, params, securityv1.MigrationStatus{SpecHash: desiredHash}); err != nil {
			return err
		}
	}

	// Build the desired Job and look for the current in-cluster Job.
	desiredJob := gateway.NewMigrationJob(params.Instance)
	if err := controllerutil.SetControllerReference(params.Instance, desiredJob, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on migration job: %w", err)
	}

	currentJob := &batchv1.Job{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredJob.Name, Namespace: desiredJob.Namespace}, currentJob)

	if err != nil && k8serrors.IsNotFound(err) {
		// Job does not exist — create it. This covers both the initial run and
		// the user-deletes-failed-job-to-retry flow.
		if err = params.Client.Create(ctx, desiredJob); err != nil {
			return fmt.Errorf("failed creating migration job: %w", err)
		}
		params.Log.Info("created migration job", "name", desiredJob.Name, "namespace", desiredJob.Namespace)
		return fmt.Errorf("%w: migration job started, waiting for completion", ErrMigrationPending)
	} else if err != nil {
		return err
	}

	// Job succeeded: persist completion in status so future reconciles skip migration.
	// Because the controller owns the Job (Owns(&batchv1.Job{})), Job completion
	// triggers a reconcile event, which is how we reach this branch promptly.
	if currentJob.Status.Succeeded > 0 {
		params.Log.Info("migration job completed successfully, updating status",
			"name", currentJob.Name, "namespace", currentJob.Namespace)
		if err := setMigrationStatus(ctx, params, securityv1.MigrationStatus{
			SpecHash: desiredHash,
			Complete: true,
		}); err != nil {
			return err
		}
		return nil
	}

	// Only declare terminal failure when no pod attempt is still active.
	// With backoffLimit=1, Kubernetes briefly shows Active=1, Failed=1 while
	// pod-2 (the retry) is running. Declaring failure then would kill a healthy
	// retry and require unnecessary manual intervention.
	if currentJob.Status.Active == 0 && currentJob.Status.Failed > 0 {
		params.Log.Info("migration job failed — Gateway deployment is blocked. Investigate logs then delete the job to retry",
			"name", currentJob.Name,
			"namespace", currentJob.Namespace,
			"fix", "kubectl delete job "+currentJob.Name+" -n "+currentJob.Namespace)
		return fmt.Errorf("migration job %s/%s failed — Gateway deployment blocked pending manual intervention", currentJob.Namespace, currentJob.Name)
	}

	// Job is still running (or briefly between pod attempts). Return ErrMigrationPending
	// so the controller requeues at a fixed interval without incrementing failure metrics.
	params.Log.Info("migration job is active, waiting...")
	return fmt.Errorf("%w: migration job is active", ErrMigrationPending)
}

// migrationSpecHash returns a 16-character hex hash of the migration-relevant spec
// fields: image, effective jdbcUrl, clearLocks flag, and activeDeadlineSeconds.
// When any of these change the hash changes and a fresh migration job is triggered.
func migrationSpecHash(gw *securityv1.Gateway) string {
	h := sha256.New()
	h.Write([]byte(gw.Spec.App.Image))

	// Use the migration-specific jdbcUrl when set; otherwise fall back to the
	// main database URL. This mirrors what NewMigrationJob passes to the container.
	jdbcUrl := gw.Spec.App.Management.Database.JDBCUrl
	if gw.Spec.App.Management.Database.MigrationJob.JDBCUrl != "" {
		jdbcUrl = gw.Spec.App.Management.Database.MigrationJob.JDBCUrl
	}
	h.Write([]byte(jdbcUrl))

	h.Write([]byte(strconv.FormatBool(gw.Spec.App.Management.Database.MigrationJob.ClearLocks)))

	if gw.Spec.App.Management.Database.MigrationJob.ActiveDeadlineSeconds != nil {
		h.Write([]byte(strconv.FormatInt(*gw.Spec.App.Management.Database.MigrationJob.ActiveDeadlineSeconds, 10)))
	}

	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

// setMigrationStatus persists the given MigrationStatus into Gateway.Status.
// It mutates params.Instance.Status in-place so that the running reconcile sees
// the updated value without a re-fetch.
func setMigrationStatus(ctx context.Context, params Params, status securityv1.MigrationStatus) error {
	params.Instance.Status.MigrationStatus = status
	if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}
	return nil
}
