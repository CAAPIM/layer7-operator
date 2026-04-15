package reconcile

import (
	"context"
	"strings"
	"testing"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestRecoverReadyIfStuck_SetsReadyWhenNotReady(t *testing.T) {
	ctx := context.Background()
	name := "recover-ready-" + strings.ReplaceAll(t.Name(), "/", "-")

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
	name := "recover-ready-" + strings.ReplaceAll(t.Name(), "/", "-")

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
