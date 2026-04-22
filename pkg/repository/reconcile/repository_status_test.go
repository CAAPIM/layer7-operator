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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestRecoverReadyIfStuck_SetsReadyWhenNotReady(t *testing.T) {
	ctx := context.Background()
	name := "recover-ready-" + strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))

	repo := &securityv1.Repository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "security.brcmlabs.com/v1",
			Kind:       "Repository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: securityv1.RepositorySpec{
			Enabled:  true,
			Type:     securityv1.RepositoryTypeGit,
			Branch:   "main",
			Endpoint: "https://example.com/repo.git",
		},
	}
	if err := k8sClient.Create(ctx, repo); err != nil {
		t.Fatal(err)
	}

	commit := "abc123def"
	storageName := name + "-repository-main"
	repo.Status = securityv1.RepositoryStatus{
		Name:              name,
		Ready:             false,
		Commit:            commit,
		StorageSecretName: storageName,
		StateStoreSynced:  true,
	}
	if err := k8sClient.Status().Update(ctx, repo); err != nil {
		t.Fatal(err)
	}

	// Seed repo cache so BuildRepositoryCache hits the commit file fast path (no clone on disk in unit tests).
	cachePath := util.RepoCacheDir(name, "default")
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cachePath, commit+".json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	params := Params{
		Client:   k8sClient,
		Log:      logger,
		Scheme:   testScheme,
		Instance: repo,
	}
	ensureReadyOnIdleSync(ctx, params, time.Now(), commit, storageName, true)

	if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: "default"}, repo); err != nil {
		t.Fatal(err)
	}
	if !repo.Status.Ready {
		t.Errorf("expected Ready true after ensureReadyOnIdleSync")
	}
	if repo.Status.Commit != commit {
		t.Errorf("commit: expected %q, got %q", commit, repo.Status.Commit)
	}
}

func TestRecoverReadyIfStuck_NoOpWhenAlreadyReady(t *testing.T) {
	ctx := context.Background()
	name := "recover-ready-" + strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))

	repo := &securityv1.Repository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "security.brcmlabs.com/v1",
			Kind:       "Repository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: securityv1.RepositorySpec{
			Enabled:  true,
			Type:     securityv1.RepositoryTypeGit,
			Branch:   "main",
			Endpoint: "https://example.com/repo.git",
		},
	}
	if err := k8sClient.Create(ctx, repo); err != nil {
		t.Fatal(err)
	}

	commit := "xyz789"
	repo.Status = securityv1.RepositoryStatus{
		Name:              name,
		Ready:             true,
		Commit:            commit,
		StorageSecretName: name + "-repository-main",
		StateStoreSynced:  true,
	}
	if err := k8sClient.Status().Update(ctx, repo); err != nil {
		t.Fatal(err)
	}
	updated := repo.Status.Updated

	params := Params{
		Client:   k8sClient,
		Log:      logger,
		Scheme:   testScheme,
		Instance: repo,
	}
	ensureReadyOnIdleSync(ctx, params, time.Now(), commit, repo.Status.StorageSecretName, true)

	if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: "default"}, repo); err != nil {
		t.Fatal(err)
	}
	if !repo.Status.Ready {
		t.Errorf("expected Ready to remain true")
	}
	if repo.Status.Updated != updated {
		t.Errorf("expected status Updated to be unchanged when already Ready")
	}
}
