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
	"errors"
	"reflect"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// seedRepoDirectories persists a Gateway whose single RepositoryStatus entry
// records the given directories, and returns the repo name used.
func seedRepoDirectories(t *testing.T, params Params, repoName string, dirs []string) {
	t.Helper()
	createGateway(t, params)
	params.Instance.Status.RepositoryStatus = []securityv1.GatewayRepositoryStatus{
		{Enabled: true, Name: repoName, Commit: "c1", Type: "dynamic", Directories: dirs},
	}
	if err := k8sClient.Status().Update(ctx, params.Instance); err != nil {
		t.Fatalf("seedRepoDirectories: %v", err)
	}
}

// readRepoDirectories re-reads the Gateway from the API server and returns the
// recorded directories for the named repository.
func readRepoDirectories(t *testing.T, params Params, repoName string) []string {
	t.Helper()
	got := &securityv1.Gateway{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, got); err != nil {
		t.Fatalf("readRepoDirectories: %v", err)
	}
	for _, rs := range got.Status.RepositoryStatus {
		if rs.Name == repoName {
			return rs.Directories
		}
	}
	t.Fatalf("readRepoDirectories: repo %q not found in status", repoName)
	return nil
}

func newDirRepo(repoName string) (securityv1.Repository, securityv1.RepositoryReference) {
	repository := securityv1.Repository{
		ObjectMeta: metav1.ObjectMeta{Name: repoName, Namespace: "default"},
		Spec:       securityv1.RepositorySpec{Type: securityv1.RepositoryTypeGit, Endpoint: "github.com"},
	}
	repoRef := securityv1.RepositoryReference{
		Name:        repoName,
		Type:        securityv1.RepositoryReferenceTypeDynamic,
		Directories: []string{"/dir2"},
	}
	return repository, repoRef
}

// TestUpdateRepoRefStatus_directoryPreservedOnApplyError is the regression test
// for the directory-apply bug: when a directory change fails to apply, the
// recorded directories must NOT advance to the new value. Advancing them would
// make the next reconcile see previous == current, skip the directoryChange
// path, and silently drop the failed change. The old directories must be
// preserved so directoryChanged stays true and the change is retried.
func TestUpdateRepoRefStatus_directoryPreservedOnApplyError(t *testing.T) {
	repoName := "dirbug-fail"
	params := newParams()
	params.Instance.Name = repoName

	seedRepoDirectories(t, params, repoName, []string{"/dir1"})
	repository, repoRef := newDirRepo(repoName)

	// Directory change from /dir1 -> /dir2 that fails to apply.
	if err := updateRepoRefStatus(ctx, params, repository, repoRef, "c2", errors.New("graphman apply failed"), false); err != nil {
		t.Fatalf("updateRepoRefStatus returned error: %v", err)
	}

	got := readRepoDirectories(t, params, repoName)
	if !reflect.DeepEqual(got, []string{"/dir1"}) {
		t.Fatalf("directories advanced on apply failure: got %v, want [/dir1] preserved so the change is retried", got)
	}
}

// TestUpdateRepoRefStatus_directoryAdvancesOnSuccess guards against
// over-correction: on a successful apply the directories must advance to the
// new value so the change is not re-applied forever.
func TestUpdateRepoRefStatus_directoryAdvancesOnSuccess(t *testing.T) {
	repoName := "dirbug-ok"
	params := newParams()
	params.Instance.Name = repoName

	seedRepoDirectories(t, params, repoName, []string{"/dir1"})
	repository, repoRef := newDirRepo(repoName)

	// Directory change from /dir1 -> /dir2 that applies successfully.
	if err := updateRepoRefStatus(ctx, params, repository, repoRef, "c2", nil, false); err != nil {
		t.Fatalf("updateRepoRefStatus returned error: %v", err)
	}

	got := readRepoDirectories(t, params, repoName)
	if !reflect.DeepEqual(got, []string{"/dir2"}) {
		t.Fatalf("directories did not advance on success: got %v, want [/dir2]", got)
	}
}
