package reconcile

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTempStorageParams(name, namespace, portal, apiName string) Params {
	return Params{
		Instance: &v1alpha1.L7Api{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec: v1alpha1.L7ApiSpec{
				PortalPublished: true,
				L7Portal:        portal,
				PortalMeta:      v1alpha1.PortalMeta{Name: apiName},
			},
		},
	}
}

func readPortalMetaName(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	meta := v1alpha1.PortalMeta{}
	if err := json.Unmarshal(b, &meta); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", path, err)
	}
	return meta.Name
}

func TestTempStoragePartitionedByNamespace(t *testing.T) {
	orig := portalTempDirectory
	portalTempDirectory = t.TempDir() + "/"
	defer func() { portalTempDirectory = orig }()
	ctx := context.Background()

	// same L7Api name and same L7Portal name in two namespaces
	victim := newTempStorageParams("api", "ns-a", "portal", "victim-api")
	attacker := newTempStorageParams("api", "ns-b", "portal", "attacker-api")

	if err := WriteTempStorage(ctx, victim); err != nil {
		t.Fatal(err)
	}
	if err := WriteTempStorage(ctx, attacker); err != nil {
		t.Fatal(err)
	}

	victimFile := filepath.Join(portalTempDirectory, "ns-a", "portal", "api.json")
	attackerFile := filepath.Join(portalTempDirectory, "ns-b", "portal", "api.json")
	if got := readPortalMetaName(t, victimFile); got != "victim-api" {
		t.Errorf("ns-a portal api was overwritten, expected %q, actual %q", "victim-api", got)
	}
	if got := readPortalMetaName(t, attackerFile); got != "attacker-api" {
		t.Errorf("expected %q, actual %q", "attacker-api", got)
	}
	entries, err := os.ReadDir(filepath.Join(portalTempDirectory, "ns-a", "portal"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 api in ns-a portal directory, actual %d", len(entries))
	}

	// removing the same-named L7Api in ns-b must not remove the ns-a entry
	if err := RemoveTempStorage(ctx, attacker); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(attackerFile); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed", attackerFile)
	}
	if _, err := os.Stat(victimFile); err != nil {
		t.Errorf("ns-a portal api was removed by ns-b L7Api: %v", err)
	}
}

func TestTempStorageRejectsInvalidPathSegments(t *testing.T) {
	orig := portalTempDirectory
	portalTempDirectory = t.TempDir() + "/"
	defer func() { portalTempDirectory = orig }()
	ctx := context.Background()

	for _, p := range []Params{
		newTempStorageParams("api", "ns-a", "../portal", "x"),
		newTempStorageParams("api", "ns-a", "a/b", "x"),
		newTempStorageParams("api", "", "portal", "x"),
		newTempStorageParams("api", "..", "portal", "x"),
		newTempStorageParams("api", "ns-a", "..", "x"),
	} {
		if err := WriteTempStorage(ctx, p); err == nil {
			t.Errorf("expected error for namespace %q portal %q", p.Instance.Namespace, p.Instance.Spec.L7Portal)
		}
		if err := RemoveTempStorage(ctx, p); err == nil {
			t.Errorf("expected error for namespace %q portal %q", p.Instance.Namespace, p.Instance.Spec.L7Portal)
		}
	}

	// Kubernetes object names may contain dots.
	if err := WriteTempStorage(ctx, newTempStorageParams("api", "ns-a", "dev.portal", "x")); err != nil {
		t.Errorf("dotted portal name rejected: %v", err)
	}
}
