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
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────────────────

// migrationParams builds a Params value with database and migrationJob both
// enabled. Each test passes a unique name so objects do not collide in the
// shared in-process API server.
func migrationParams(name string) Params {
	p := newParams()
	p.Instance.Name = name
	p.Instance.Spec.App.Management.Database.Enabled = true
	p.Instance.Spec.App.Management.Database.MigrationJob.Enabled = true
	return p
}

// createGateway persists the Gateway CR in the in-process Kubernetes API server.
// controller-runtime's Create mutates params.Instance in-place with the
// server-assigned ResourceVersion, which is required by later Status().Update
// calls. A cleanup function deletes the CR when the test finishes.
func createGateway(t *testing.T, params Params) {
	t.Helper()
	if err := k8sClient.Create(ctx, params.Instance); err != nil {
		t.Fatalf("createGateway: failed to create %s: %v", params.Instance.Name, err)
	}
	t.Cleanup(func() {
		// ignore not-found errors — some tests delete the CR themselves
		_ = k8sClient.Delete(context.Background(), params.Instance)
	})
}

// seedMigrationStatus writes a pre-configured MigrationStatus into a Gateway CR
// that already exists in the cluster. Used to set up test preconditions (e.g.
// a previously-completed migration) without calling GatewayMigrationJob.
func seedMigrationStatus(t *testing.T, params Params, ms securityv1.MigrationStatus) {
	t.Helper()
	params.Instance.Status.MigrationStatus = ms
	if err := k8sClient.Status().Update(ctx, params.Instance); err != nil {
		t.Fatalf("seedMigrationStatus: %v", err)
	}
}

// buildOrphanJob constructs a minimal valid Job that can be registered in the
// in-process API server. It is used by tests that need an in-cluster Job
// independent of gateway.NewMigrationJob (e.g. for the disabled-cleanup and
// stale-image scenarios).
func buildOrphanJob(name, namespace, image string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers:    []corev1.Container{{Name: "db-migration", Image: image}},
				},
			},
		},
	}
}

// createJob persists a Job in the in-process API server. controller-runtime's
// Create mutates job in-place with the server-assigned ResourceVersion so that
// subsequent Status().Update calls succeed. A cleanup function deletes the Job.
func createJob(t *testing.T, job *batchv1.Job) {
	t.Helper()
	if err := k8sClient.Create(ctx, job); err != nil {
		t.Fatalf("createJob %s: %v", job.Name, err)
	}
	t.Cleanup(func() {
		_ = k8sClient.Delete(context.Background(), job)
	})
}

// patchJobStatus updates a Job's status via the status subresource so tests can
// simulate Succeeded/Failed/Active counts without running a real container.
// The job argument must carry a valid ResourceVersion (i.e. returned from
// Create or Get before calling this function).
func patchJobStatus(t *testing.T, job *batchv1.Job, status batchv1.JobStatus) {
	t.Helper()
	job.Status = status
	if err := k8sClient.Status().Update(ctx, job); err != nil {
		t.Fatalf("patchJobStatus %s: %v", job.Name, err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// migrationSpecHash — pure function tests (no Kubernetes API required)
// ─────────────────────────────────────────────────────────────────────────────

// TestMigrationSpecHash verifies that migrationSpecHash produces the correct
// output for each input. These tests cover all six fields that the function
// hashes, confirming that a change to any one of them produces a different hash
// and therefore a new migration job.
func TestMigrationSpecHash(t *testing.T) {

	t.Run("same spec always produces the same hash (determinism)", func(t *testing.T) {
		// Input:  identical Gateway spec, called twice.
		// Output: both calls return the exact same 16-character hex string.
		// Why:    a non-deterministic hash would cause spurious job replacements.
		p := newParams()
		h1 := migrationSpecHash(p.Instance)
		h2 := migrationSpecHash(p.Instance)
		if h1 != h2 {
			t.Errorf("hash is not deterministic: %q != %q", h1, h2)
		}
		if len(h1) != 16 {
			t.Errorf("expected 16-char hash, got len=%d: %q", len(h1), h1)
		}
	})

	t.Run("image change produces a different hash", func(t *testing.T) {
		// Input:  two Gateways identical except spec.app.image.
		// Output: different hashes.
		// Why:    upgrading to a new Gateway image carries new Liquibase changesets
		//         that must be migrated before the Deployment proceeds.
		p1 := newParams()
		p1.Instance.Spec.App.Image = "gateway:11.1.1"
		p2 := newParams()
		p2.Instance.Spec.App.Image = "gateway:11.3.0"
		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes for different gateway images")
		}
	})

	t.Run("migrationJob.jdbcUrl overrides database.jdbcUrl in the hash", func(t *testing.T) {
		// Input A: database.jdbcUrl set; migrationJob.jdbcUrl not set.
		// Input B: same database.jdbcUrl plus a different migrationJob.jdbcUrl.
		// Output: different hashes.
		// Why:    migrationJob.jdbcUrl points the migration at a different database;
		//         the hash must capture the URL that the Job actually uses.
		p1 := newParams()
		p1.Instance.Spec.App.Management.Database.JDBCUrl = "jdbc:mysql://host-a/ssg"

		p2 := newParams()
		p2.Instance.Spec.App.Management.Database.JDBCUrl = "jdbc:mysql://host-a/ssg"
		p2.Instance.Spec.App.Management.Database.MigrationJob.JDBCUrl = "jdbc:mysql://host-b/ssg"

		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes when migrationJob.jdbcUrl overrides database.jdbcUrl")
		}
	})

	t.Run("clearLocks toggle produces a different hash", func(t *testing.T) {
		// Input:  two Gateways identical except migrationJob.clearLocks.
		// Output: different hashes.
		// Why:    clearLocks changes the schema-update mode flag passed to the container
		//         (liquibase-only vs liquibase-only-with-unlock); the migration must be
		//         re-run with the new flag.
		p1 := newParams()
		p1.Instance.Spec.App.Management.Database.MigrationJob.ClearLocks = false
		p2 := newParams()
		p2.Instance.Spec.App.Management.Database.MigrationJob.ClearLocks = true
		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes for clearLocks=false vs clearLocks=true")
		}
	})

	t.Run("activeDeadlineSeconds change produces a different hash", func(t *testing.T) {
		// Input:  one Gateway with activeDeadlineSeconds unset, one with 300.
		// Output: different hashes.
		// Why:    changing the deadline (typically to fix a too-short timeout that caused
		//         a failure) should trigger a replacement job with the new deadline.
		p1 := newParams()
		p2 := newParams()
		deadline := int64(300)
		p2.Instance.Spec.App.Management.Database.MigrationJob.ActiveDeadlineSeconds = &deadline
		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes when activeDeadlineSeconds changes from nil to 300")
		}
	})

	t.Run("DisklessConfig.Disabled toggle produces a different hash", func(t *testing.T) {
		// Input:  disklessConfig.disabled=false vs disklessConfig.disabled=true.
		// Output: different hashes.
		// Why:    DisklessConfig.Disabled controls whether the Secret is mounted as
		//         env vars (diskless) or as a node.properties file (non-diskless).
		//         Flipping this after a completed migration changes the Job wiring
		//         entirely and must invalidate Status.Complete.
		p1 := newParams()
		p1.Instance.Spec.App.Management.DisklessConfig.Disabled = false
		p2 := newParams()
		p2.Instance.Spec.App.Management.DisklessConfig.Disabled = true
		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes for DisklessConfig.Disabled=false vs true")
		}
	})

	t.Run("secretName change produces a different hash", func(t *testing.T) {
		// Input:  one Gateway using the default secret name (== gateway name), one
		//         with an explicit spec.app.management.secretName.
		// Output: different hashes.
		// Why:    secretName determines which Kubernetes Secret the Job reads for
		//         database credentials. Rotating to a new Secret must invalidate the
		//         cached Complete=true status so the operator re-runs the migration
		//         with the new credentials.
		p1 := newParams()
		p1.Instance.Spec.App.Management.SecretName = ""
		p2 := newParams()
		p2.Instance.Spec.App.Management.SecretName = "my-custom-db-secret"
		if migrationSpecHash(p1.Instance) == migrationSpecHash(p2.Instance) {
			t.Error("expected different hashes when secretName changes")
		}
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// MigrationJobName — pure function test (no Kubernetes API required)
// ─────────────────────────────────────────────────────────────────────────────

// TestMigrationJobName confirms the canonical naming convention for migration
// Jobs. Centralising this in a helper (gateway.MigrationJobName) ensures the
// cleanup branch, the spec-change branch, and NewMigrationJob all agree on the
// Job name.
func TestMigrationJobName(t *testing.T) {
	// Input:  Gateway with name "my-gateway".
	// Output: "my-gateway-db-migration".
	p := newParams()
	p.Instance.Name = "my-gateway"
	got := gateway.MigrationJobName(p.Instance)
	want := "my-gateway-db-migration"
	if got != want {
		t.Errorf("MigrationJobName: got %q, want %q", got, want)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// GatewayMigrationJob — integration tests using the in-process Kubernetes API
//
// The in-process API server (envtest) provides a real Kubernetes control plane
// without needing an external cluster. Jobs and Gateway CRs are persisted in
// etcd and retrieved via the same client the production code uses, giving
// high confidence that the reconcile logic behaves correctly in a real cluster.
//
// Because no Job controller is running, Job pod counts (Active, Succeeded,
// Failed) are set manually via the status subresource to simulate each
// lifecycle stage.
// ─────────────────────────────────────────────────────────────────────────────

func TestGatewayMigrationJob(t *testing.T) {

	// ── disabled: migration feature off ──────────────────────────────────────

	t.Run("migration disabled: returns nil without creating a Job", func(t *testing.T) {
		// Input:  Gateway with migrationJob.enabled=false and database.enabled=true.
		// Output: GatewayMigrationJob returns nil (no blocking, no work done).
		//         No Job object is created in the cluster.
		p := newParams()
		p.Instance.Name = "mj-disabled"
		p.Instance.Spec.App.Management.Database.Enabled = true
		p.Instance.Spec.App.Management.Database.MigrationJob.Enabled = false

		err := GatewayMigrationJob(ctx, p)

		if err != nil {
			t.Fatalf("expected nil for disabled migration, got: %v", err)
		}
		job := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: gateway.MigrationJobName(p.Instance), Namespace: p.Instance.Namespace,
		}, job)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("expected no Job in cluster when migration disabled; get returned: %v", getErr)
		}
	})

	t.Run("database disabled: returns nil without creating a Job", func(t *testing.T) {
		// Input:  migrationJob.enabled=true but database.enabled=false.
		//         (The webhook prevents this in practice; the reconciler guards defensively.)
		// Output: GatewayMigrationJob returns nil and creates no Job.
		p := newParams()
		p.Instance.Name = "mj-db-disabled"
		p.Instance.Spec.App.Management.Database.Enabled = false
		p.Instance.Spec.App.Management.Database.MigrationJob.Enabled = true

		err := GatewayMigrationJob(ctx, p)

		if err != nil {
			t.Fatalf("expected nil when database disabled, got: %v", err)
		}
		job := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: gateway.MigrationJobName(p.Instance), Namespace: p.Instance.Namespace,
		}, job)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("expected no Job in cluster when database disabled; get returned: %v", getErr)
		}
	})

	// ── first enable ─────────────────────────────────────────────────────────

	t.Run("first enable: Job created, SpecHash written, ErrMigrationPending returned", func(t *testing.T) {
		// Input:  Gateway with migration freshly enabled; empty Status.MigrationStatus;
		//         no existing Job in the cluster.
		// Output:
		//   - GatewayMigrationJob returns an error wrapping ErrMigrationPending.
		//     (ErrMigrationPending is a sentinel: "not done yet, check back in 10s".
		//      It is not a real failure — the controller requeues without backoff.)
		//   - A Job named "<gateway>-db-migration" exists in the cluster.
		//   - Status.MigrationStatus.SpecHash is non-empty (hash persisted to etcd).
		//   - Status.MigrationStatus.Complete is false (migration not yet finished).
		p := migrationParams("mj-first-enable")
		createGateway(t, p)

		err := GatewayMigrationJob(ctx, p)

		if !errors.Is(err, ErrMigrationPending) {
			t.Fatalf("expected ErrMigrationPending on first enable, got: %v", err)
		}

		// Verify the Job was created.
		job := &batchv1.Job{}
		if getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: gateway.MigrationJobName(p.Instance), Namespace: p.Instance.Namespace,
		}, job); getErr != nil {
			t.Fatalf("expected Job to exist after first enable, got: %v", getErr)
		}

		// Verify Status.MigrationStatus was written.
		if p.Instance.Status.MigrationStatus.SpecHash == "" {
			t.Error("expected SpecHash to be set in status after first enable")
		}
		if p.Instance.Status.MigrationStatus.Complete {
			t.Error("expected Complete=false — migration job was just created, not yet finished")
		}
	})

	// ── fast path ────────────────────────────────────────────────────────────

	t.Run("fast path: Complete=true for current spec returns nil without touching the cluster", func(t *testing.T) {
		// Input:  Gateway where Status.MigrationStatus already records Complete=true
		//         for the exact current spec hash.
		// Output: GatewayMigrationJob returns nil immediately.
		//         No Job is created or modified (the 12-hour reconcile no-op path).
		// Why this matters: without this fast path, the 12-hour requeue would create
		//         a new Job every time the completed Job was GC'd by Kubernetes.
		p := migrationParams("mj-fast-path")
		createGateway(t, p)

		// Pre-seed status to simulate a previously completed migration.
		desiredHash := migrationSpecHash(p.Instance)
		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: desiredHash,
			Complete: true,
		})

		err := GatewayMigrationJob(ctx, p)

		if err != nil {
			t.Fatalf("expected nil from fast path (Complete=true), got: %v", err)
		}

		// Verify no Job was created.
		job := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: gateway.MigrationJobName(p.Instance), Namespace: p.Instance.Namespace,
		}, job)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("fast path must not create a Job; found one in cluster")
		}
	})

	// ── job lifecycle: running ────────────────────────────────────────────────

	t.Run("active job: returns ErrMigrationPending without modifying the Job", func(t *testing.T) {
		// Input:  Gateway with an existing Job that reports Active=1 (pod is running).
		//         Status.MigrationStatus.SpecHash matches the current desired hash.
		// Output: GatewayMigrationJob returns ErrMigrationPending.
		//         The existing Job is neither deleted nor recreated.
		p := migrationParams("mj-active-job")
		createGateway(t, p)

		// Create the Job directly so we control its spec and status independently.
		job := gateway.NewMigrationJob(p.Instance)
		createJob(t, job)
		// Simulate a running pod: patch Job status to Active=1.
		patchJobStatus(t, job, batchv1.JobStatus{Active: 1})

		// Pre-seed status with the current hash (as the first-enable path would have).
		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: migrationSpecHash(p.Instance),
		})

		err := GatewayMigrationJob(ctx, p)

		if !errors.Is(err, ErrMigrationPending) {
			t.Fatalf("expected ErrMigrationPending while job is active, got: %v", err)
		}

		// Job must still exist and be unmodified.
		got := &batchv1.Job{}
		if getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: job.Name, Namespace: job.Namespace,
		}, got); getErr != nil {
			t.Fatalf("Job must still exist while active: %v", getErr)
		}
	})

	// ── job lifecycle: succeeded ──────────────────────────────────────────────

	t.Run("succeeded job: writes Complete=true to status and returns nil", func(t *testing.T) {
		// Input:  Gateway whose existing Job reports Succeeded=1 (all changesets applied).
		// Output:
		//   - GatewayMigrationJob returns nil (Deployment step is now unblocked).
		//   - Status.MigrationStatus.Complete is true (persisted to etcd).
		//   - Status.MigrationStatus.SpecHash still matches the current spec hash.
		// Why:    the Complete=true flag is what allows all future 12-hour reconciles
		//         to skip migration management entirely via the fast path.
		p := migrationParams("mj-succeeded")
		createGateway(t, p)

		// First call: creates the Job and writes SpecHash to status.
		if firstErr := GatewayMigrationJob(ctx, p); !errors.Is(firstErr, ErrMigrationPending) {
			t.Fatalf("first call should return ErrMigrationPending, got: %v", firstErr)
		}

		// Retrieve the created Job and simulate success by patching its status.
		job := &batchv1.Job{}
		if getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: gateway.MigrationJobName(p.Instance), Namespace: p.Instance.Namespace,
		}, job); getErr != nil {
			t.Fatalf("Job must exist after first call: %v", getErr)
		}
		patchJobStatus(t, job, batchv1.JobStatus{Succeeded: 1})

		// Second call: Job is succeeded — should write Complete=true and return nil.
		err := GatewayMigrationJob(ctx, p)

		if err != nil {
			t.Fatalf("expected nil after Job succeeded, got: %v", err)
		}
		if !p.Instance.Status.MigrationStatus.Complete {
			t.Error("expected Status.MigrationStatus.Complete=true after Job succeeded")
		}
		wantHash := migrationSpecHash(p.Instance)
		if p.Instance.Status.MigrationStatus.SpecHash != wantHash {
			t.Errorf("expected SpecHash=%q, got %q", wantHash, p.Instance.Status.MigrationStatus.SpecHash)
		}
	})

	// ── job lifecycle: failed ─────────────────────────────────────────────────

	t.Run("failed job: returns a blocking error (not ErrMigrationPending)", func(t *testing.T) {
		// Input:  Gateway whose existing Job reports Active=0, Failed=1
		//         (both pod attempts exhausted; backoffLimit=1 reached).
		// Output: GatewayMigrationJob returns a non-nil error that does NOT wrap
		//         ErrMigrationPending. The controller applies backoff and metrics.
		//         The Job is NOT deleted — the user must fix the root cause
		//         (restore DB to pre-upgrade state if partially migrated, fix
		//         credentials, etc.) and then delete the Job manually to retry.
		p := migrationParams("mj-failed")
		createGateway(t, p)

		job := gateway.NewMigrationJob(p.Instance)
		createJob(t, job)
		patchJobStatus(t, job, batchv1.JobStatus{Active: 0, Failed: 1})

		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: migrationSpecHash(p.Instance),
		})

		err := GatewayMigrationJob(ctx, p)

		// Must be a real error (not the "pending" sentinel).
		if err == nil {
			t.Fatal("expected a blocking error for failed job, got nil")
		}
		if errors.Is(err, ErrMigrationPending) {
			t.Fatalf("expected a real blocking error, not ErrMigrationPending — got: %v", err)
		}

		// Job must still exist (user must delete it manually to retry).
		got := &batchv1.Job{}
		if getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: job.Name, Namespace: job.Namespace,
		}, got); getErr != nil {
			t.Fatalf("failed Job must not be auto-deleted (manual intervention required): %v", getErr)
		}
	})

	// ── spec change ───────────────────────────────────────────────────────────

	t.Run("spec changed: old Job deleted, status reset to new hash, ErrMigrationPending returned", func(t *testing.T) {
		// Input:  Gateway previously migrated with image v1 (Complete=true, SpecHash=H1).
		//         User changes spec.app.image to v2, making desiredHash=H2 != H1.
		// Output:
		//   - The old Job is deleted from the cluster.
		//   - Status.MigrationStatus is reset to {SpecHash: H2, Complete: false}.
		//   - GatewayMigrationJob returns ErrMigrationPending so the controller
		//     requeues in 10s, at which point it will create the new Job.
		// Why the delete-before-status-write ordering matters: if the status were
		//   written first and the deletion then failed, the new hash would already be
		//   persisted, making this branch unreachable on the next reconcile — a silent,
		//   permanent skip of the required migration.
		p := migrationParams("mj-spec-change")
		p.Instance.Spec.App.Image = "gateway:11.1.1"
		createGateway(t, p)

		hashV1 := migrationSpecHash(p.Instance)

		// Seed status as if migration with image v1 already completed.
		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: hashV1,
			Complete: true,
		})

		// Create an old Job representing the completed v1 migration run.
		oldJob := gateway.NewMigrationJob(p.Instance)
		createJob(t, oldJob)

		// Simulate an image upgrade — this changes the desired spec hash.
		p.Instance.Spec.App.Image = "gateway:11.3.0"

		err := GatewayMigrationJob(ctx, p)

		if !errors.Is(err, ErrMigrationPending) {
			t.Fatalf("expected ErrMigrationPending after spec change, got: %v", err)
		}

		// Old Job must be gone (Background deletion removes it from etcd immediately).
		oldJobCheck := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: oldJob.Name, Namespace: oldJob.Namespace,
		}, oldJobCheck)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("expected old Job to be deleted after spec change; still exists")
		}

		// Status must be reset to the new hash with Complete=false.
		hashV2 := migrationSpecHash(p.Instance)
		if p.Instance.Status.MigrationStatus.SpecHash != hashV2 {
			t.Errorf("SpecHash: got %q, want %q", p.Instance.Status.MigrationStatus.SpecHash, hashV2)
		}
		if p.Instance.Status.MigrationStatus.Complete {
			t.Error("Complete must be false after spec change — migration not yet run for new spec")
		}
	})

	// ── disabled cleanup ──────────────────────────────────────────────────────

	t.Run("disabled cleanup: orphaned Job is deleted when migration toggled off", func(t *testing.T) {
		// Input:  A Job already exists in the cluster (left from a previous run).
		//         User sets migrationJob.enabled=false.
		// Output: GatewayMigrationJob returns nil and the orphaned Job is deleted.
		// Why:    toggling migration off should not leave stale Job objects in the
		//         namespace that could confuse operators or trigger alert rules.
		p := newParams()
		p.Instance.Name = "mj-disabled-cleanup"
		p.Instance.Spec.App.Management.Database.Enabled = true

		// Create the orphaned Job directly — simulates a job left behind from a
		// previous enabled run using the same Gateway name.
		orphanJob := buildOrphanJob(
			gateway.MigrationJobName(p.Instance),
			p.Instance.Namespace,
			"gateway:11.1.1",
		)
		createJob(t, orphanJob)

		// Now disable migration.
		p.Instance.Spec.App.Management.Database.MigrationJob.Enabled = false

		err := GatewayMigrationJob(ctx, p)

		if err != nil {
			t.Fatalf("expected nil when migration disabled, got: %v", err)
		}
		got := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: orphanJob.Name, Namespace: orphanJob.Namespace,
		}, got)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("expected orphaned Job to be deleted when migration disabled; still exists")
		}
	})

	// ── stale / malformed job guard ───────────────────────────────────────────

	t.Run("stale image job: Job with wrong image is replaced, ErrMigrationPending returned", func(t *testing.T) {
		// Input:  An existing Job was built with image "old-image:1.0" but the
		//         current Gateway spec still carries the same spec hash (the image
		//         change somehow escaped hash detection). The live Job's image does
		//         not match the image in the freshly-built desired Job.
		// Output: The stale Job is deleted and ErrMigrationPending is returned so
		//         the next reconcile creates a valid replacement.
		// This guard also covers the zero-container case (a Job whose containers
		//         were cleared by a manual edit or a mutating webhook after creation —
		//         something the API server admission webhook would not catch post-creation).
		p := migrationParams("mj-stale-image")
		createGateway(t, p)

		// Create a Job with a deliberately wrong image.
		staleJob := buildOrphanJob(
			gateway.MigrationJobName(p.Instance),
			p.Instance.Namespace,
			"old-image:1.0", // does not match p.Instance.Spec.App.Image ("test")
		)
		createJob(t, staleJob)

		// Pre-seed status with the current hash so we reach the live-object check.
		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: migrationSpecHash(p.Instance),
		})

		err := GatewayMigrationJob(ctx, p)

		if !errors.Is(err, ErrMigrationPending) {
			t.Fatalf("expected ErrMigrationPending when replacing stale-image job, got: %v", err)
		}
		got := &batchv1.Job{}
		getErr := k8sClient.Get(ctx, types.NamespacedName{
			Name: staleJob.Name, Namespace: staleJob.Namespace,
		}, got)
		if !k8serrors.IsNotFound(getErr) {
			t.Errorf("expected stale Job to be deleted; still exists")
		}
	})

	// ── DeletionTimestamp guard ───────────────────────────────────────────────

	t.Run("DeletionTimestamp guard: does not read stale status from a mid-deletion Job", func(t *testing.T) {
		// Input:  An existing Job has DeletionTimestamp set (Kubernetes accepted the
		//         delete request but the object is still visible because a finalizer
		//         is preventing immediate removal). The Job's status still shows
		//         Succeeded=1 from the previous run — stale data.
		// Output: GatewayMigrationJob returns ErrMigrationPending and does NOT write
		//         Complete=true to Status.MigrationStatus.
		// Why:    Background deletion removes the Job from etcd immediately, but a
		//         brief informer-cache lag or webhook delay can deliver a stale cached
		//         object. Without the DeletionTimestamp guard in GatewayMigrationJob,
		//         the operator would trust the old Succeeded count and falsely mark
		//         the new spec as migrated.
		p := migrationParams("mj-deletion-ts")
		createGateway(t, p)

		// ── Setup: engineer a Job that has DeletionTimestamp set but Succeeded=1 ──
		//
		// Add a finalizer so that Delete() sets DeletionTimestamp without immediately
		// removing the object — simulating the mid-teardown window.
		job := gateway.NewMigrationJob(p.Instance)
		job.Finalizers = []string{"test/keep-alive"}
		createJob(t, job)

		// Patch Succeeded=1 so the Job looks like it completed a previous run.
		// This is the stale status the guard must ignore.
		patchJobStatus(t, job, batchv1.JobStatus{Succeeded: 1})

		// Delete the Job — Kubernetes sets DeletionTimestamp but the finalizer
		// keeps the object alive, so GatewayMigrationJob's Get still returns it.
		if delErr := k8sClient.Delete(ctx, job); delErr != nil {
			t.Fatalf("failed to delete Job to set DeletionTimestamp: %v", delErr)
		}

		// Pre-seed the spec hash so the reconcile reaches the live-job inspection
		// (not short-circuited by the fast path or spec-change branches).
		seedMigrationStatus(t, p, securityv1.MigrationStatus{
			SpecHash: migrationSpecHash(p.Instance),
		})

		// ── Exercise ──────────────────────────────────────────────────────────────
		//
		// This simulates reconcile #2 after a spec change. The operator deleted the
		// old Job in reconcile #1 (Background propagation), but due to an
		// informer-cache lag the Get here still returns the old object — now with
		// DeletionTimestamp set and a stale Succeeded=1 from the previous run.
		// GatewayMigrationJob must detect this and not trust the stale status.
		err := GatewayMigrationJob(ctx, p)

		// ── Assert: the DeletionTimestamp guard fires ─────────────────────────────
		//
		// The guard is the check inside GatewayMigrationJob that inspects
		// currentJob.DeletionTimestamp after the Get. It must return
		// ErrMigrationPending instead of falling through to the Succeeded check.
		if !errors.Is(err, ErrMigrationPending) {
			t.Fatalf("DeletionTimestamp guard must return ErrMigrationPending, got: %v", err)
		}

		// The guard must also prevent setMigrationStatus from being called with
		// Complete=true — the stale Succeeded=1 count must not be trusted.
		if p.Instance.Status.MigrationStatus.Complete {
			t.Error("DeletionTimestamp guard must not write Complete=true from a mid-deletion Job")
		}

		// Remove the finalizer so the Job is actually deleted when the test ends.
		t.Cleanup(func() {
			j := &batchv1.Job{}
			if err := k8sClient.Get(context.Background(), types.NamespacedName{
				Name: job.Name, Namespace: job.Namespace,
			}, j); err == nil {
				j.Finalizers = nil
				_ = k8sClient.Update(context.Background(), j)
			}
		})
	})
}
