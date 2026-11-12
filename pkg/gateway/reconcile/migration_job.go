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

// GatewayMigrationJob manages the lifecycle of the pre-upgrade database migration
// job. It uses Gateway.Status.MigrationStatus as the authoritative source of truth
// for whether the migration has completed, so that:
//
//   - Completed migrations are not re-run on 12-hour reconciles or operator restarts.
//   - Completed migrations survive Job deletion/GC without triggering a re-run.
//   - A spec change (image upgrade, jdbcUrl correction, clearLocks toggle, etc.)
//     resets the status and triggers a fresh migration run automatically.
//
// The function always returns nil unless it encounters a real error. Waiting states
// (Job created, Job active, Job being replaced) return nil and rely on the
// controller's Owns(&batchv1.Job{}) watch to re-trigger reconciliation as soon as
// the Job's status changes. The Deployment step is only unblocked when
// Status.MigrationStatus.Complete is true.
func GatewayMigrationJob(ctx context.Context, params Params) error {
	migrationEnabled := params.Instance.Spec.App.Management.Database.MigrationJob.Enabled
	databaseEnabled := params.Instance.Spec.App.Management.Database.Enabled

	// When migration is disabled (or database is disabled), clean up any leftover
	// job so that toggling enabled→false doesn't orphan a job in the namespace.
	if !migrationEnabled || !databaseEnabled {
		staleJob := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{
			Name:      gateway.MigrationJobName(params.Instance),
			Namespace: params.Instance.Namespace,
		}}
		propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationBackground)
		_ = client.IgnoreNotFound(params.Client.Delete(ctx, staleJob, propagationPolicy))
		return nil
	}

	desiredHash := migrationSpecHash(params.Instance)
	currentStatus := params.Instance.Status.MigrationStatus

	// Fast path: migration already completed for the current spec.
	// This allows 12-hour reconciles and operator restarts to skip migration
	// management entirely — the status persists regardless of Job existence.
	if currentStatus.SpecHash == desiredHash && currentStatus.Complete {
		params.Log.V(2).Info("migration already complete for current spec, skipping",
			"hash", desiredHash)
		return nil
	}

	// Spec changed (image upgrade, jdbcUrl correction, clearLocks/secretName/
	// disklessConfig toggle, etc.): delete the old Job first, then update status.
	//
	// Ordering matters: writing status first would update SpecHash to
	// desiredHash before deletion, making this branch unreachable on the next
	// reconcile if the deletion failed — turning a retryable error into a silent
	// permanent skip. Instead we delete first and only write status after the
	// delete succeeds or the job is already gone.
	//
	// No preceding Get (Fix 6): Delete with IgnoreNotFound is one fewer API
	// round-trip and removes a swallowed-error opportunity.
	//
	// Background propagation (Fix 2): removes the Job object from etcd immediately
	// so a subsequent Get cannot return the old (stale) object while it is
	// mid-teardown, closing the foreground-deletion timing race.
	if currentStatus.SpecHash != "" && currentStatus.SpecHash != desiredHash {
		params.Log.Info("migration spec changed, replacing job",
			"oldHash", currentStatus.SpecHash,
			"newHash", desiredHash)
		staleJob := &batchv1.Job{ObjectMeta: metav1.ObjectMeta{
			Name:      gateway.MigrationJobName(params.Instance),
			Namespace: params.Instance.Namespace,
		}}
		propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationBackground)
		delErr := params.Client.Delete(ctx, staleJob, propagationPolicy)
		if delErr != nil && !k8serrors.IsNotFound(delErr) {
			return fmt.Errorf("failed deleting stale migration job: %w", delErr)
		}
		if err := setMigrationStatus(ctx, params, securityv1.MigrationStatus{SpecHash: desiredHash}); err != nil {
			return err
		}
		// Keep the local copy in sync so the first-enable status write further
		// down (which checks for an empty SpecHash) doesn't run redundantly.
		currentStatus.SpecHash = desiredHash

		if delErr == nil {
			// A Job existed and was just deleted. Rely on the deletion event (Owns
			// watch) to retrigger reconciliation rather than reading it back here —
			// an immediate Get could still return the stale cached object before
			// the informer observes the delete.
			return nil
		}
		// No Job existed to delete (e.g. it was manually removed, or garbage
		// collected, after a previous migration completed). There is no deletion
		// event to wait for in that case, so fall through and create the
		// replacement Job in this same pass instead of stalling until the next
		// spec change or the 12-hour periodic reconcile.
	}

	// Build the desired Job now — needed both for creation and for the live
	// sanity check (Fix 4+5) further below.
	desiredJob := gateway.NewMigrationJob(params.Instance)
	if err := controllerutil.SetControllerReference(params.Instance, desiredJob, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on migration job: %w", err)
	}

	// Record the desired hash in status on first enable (SpecHash is empty).
	if currentStatus.SpecHash == "" {
		if err := setMigrationStatus(ctx, params, securityv1.MigrationStatus{SpecHash: desiredHash}); err != nil {
			return err
		}
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
		// Return nil; the Job's status changes will re-trigger reconciliation via
		// the Owns(&batchv1.Job{}) watch.
		return nil
	} else if err != nil {
		return err
	}

	// Don't trust Status on a Job that has already been requested for
	// deletion — its Succeeded/Failed counts still reflect the previous run.
	// Background deletion removes the object immediately, but a webhook delay or
	// brief informer lag can still deliver a stale cached object for a short window.
	// Return nil; the Job's full deletion event will re-trigger reconciliation.
	if currentJob.DeletionTimestamp != nil {
		return nil
	}

	// Sanity-check the live Job against the desired spec before reading
	// its Status. Re-adds the zero-container guard that existed in the previous
	// implementation, and extends it with an image check.
	//
	// If the container count is zero the Job is malformed (manual edit, webhook
	// interference) and must be replaced. If the image differs the Job carries a
	// stale spec that somehow escaped hash-based detection (e.g., a field not yet
	// covered by migrationSpecHash); replacing it is the safe fallback.
	//
	// In both cases: delete with Background propagation and return nil so the
	// Job deletion event re-triggers reconciliation to create a fresh Job.
	if len(currentJob.Spec.Template.Spec.Containers) == 0 ||
		currentJob.Spec.Template.Spec.Containers[0].Image != desiredJob.Spec.Template.Spec.Containers[0].Image {
		params.Log.Info("replacing malformed or stale migration job",
			"name", currentJob.Name, "namespace", currentJob.Namespace)
		propagationPolicy := client.PropagationPolicy(metav1.DeletePropagationBackground)
		_ = client.IgnoreNotFound(params.Client.Delete(ctx, currentJob, propagationPolicy))
		return nil
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
	// pod-2 (the retry) is running. Declaring failure then would block a healthy
	// retry and require unnecessary manual intervention.
	if currentJob.Status.Active == 0 && currentJob.Status.Failed > 0 {
		params.Log.Info("migration job failed — delete the job to retry after resolving the root cause",
			"name", currentJob.Name,
			"namespace", currentJob.Namespace)
		return fmt.Errorf("migration job %s/%s failed — Gateway deployment blocked pending manual intervention", currentJob.Namespace, currentJob.Name)
	}

	// Job is still running (or briefly between pod attempts). Return nil so the
	// reconcile loop completes normally; the Job completion event will re-trigger.
	params.Log.V(2).Info("migration job is active, waiting for completion",
		"name", currentJob.Name, "namespace", currentJob.Namespace)
	return nil
}

// migrationComplete reports whether the pre-upgrade migration job (if enabled)
// has finished for the Gateway's current spec. The Deployment step calls this
// directly so that it — not the reconcile loop's control flow — is what decides
// whether the Deployment is created/updated. When
// migration is disabled there is nothing to wait for, so it reports complete.
func migrationComplete(gw *securityv1.Gateway) bool {
	if !gw.Spec.App.Management.Database.MigrationJob.Enabled || !gw.Spec.App.Management.Database.Enabled {
		return true
	}
	status := gw.Status.MigrationStatus
	return status.Complete && status.SpecHash == migrationSpecHash(gw)
}

// migrationSpecHash returns a 16-character hex hash of all Gateway spec fields
// that affect the Job built by NewMigrationJob. When any of these change, the
// hash changes, status is reset, and a fresh migration job is triggered.
//
// Fields hashed:
//   - Image: new version = new Liquibase changesets to apply
//   - effective jdbcUrl: different target DB requires its own migration run
//   - clearLocks: changes the schema-update mode flag passed to the container
//   - activeDeadlineSeconds: changing this is typically to fix a too-short timeout
//   - DisklessConfig.Disabled: determines Secret-mounting strategy (envFrom vs Volume)
//   - effective secretName: controls which Secret the Job container references
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

	// Include the two fields NewMigrationJob also reads but that were
	// previously absent from the hash. Omitting them meant that rotating the
	// Gateway Secret or flipping DisklessConfig after a completed migration would
	// go undetected and the operator would keep trusting the old Complete: true status.
	h.Write([]byte(strconv.FormatBool(gw.Spec.App.Management.DisklessConfig.Disabled)))

	secretName := gw.Name
	if gw.Spec.App.Management.SecretName != "" {
		secretName = gw.Spec.App.Management.SecretName
	}
	h.Write([]byte(secretName))

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
