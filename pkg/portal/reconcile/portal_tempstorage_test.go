package reconcile

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func writePortalAPI(t *testing.T, namespace, portal, file, apiName string) {
	t.Helper()
	dir := portalTempPath(namespace, portal)
	if err := os.MkdirAll(dir, 0760); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(templategen.PortalAPI{Name: apiName})
	if err := os.WriteFile(filepath.Join(dir, file), b, 0660); err != nil {
		t.Fatal(err)
	}
}

func TestPortalTempPathPartitionedByNamespace(t *testing.T) {
	orig := tempDirectoryBase
	tempDirectoryBase = "/base/"
	defer func() { tempDirectoryBase = orig }()

	if got := portalTempPath("ns-a", "portal"); got != "/base/ns-a/portal" {
		t.Errorf("expected %q, actual %q", "/base/ns-a/portal", got)
	}
	if portalTempPath("ns-a", "portal") == portalTempPath("ns-b", "portal") {
		t.Error("same portal name in different namespaces must not share temp storage")
	}
}

func TestSyncPortalReadsOnlyOwnNamespace(t *testing.T) {
	orig := tempDirectoryBase
	tempDirectoryBase = t.TempDir() + "/"
	defer func() { tempDirectoryBase = orig }()
	ctx := context.Background()

	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	portalA := &v1alpha1.L7Portal{
		ObjectMeta: metav1.ObjectMeta{Name: "portal", Namespace: "ns-a"},
		Spec:       v1alpha1.L7PortalSpec{Enabled: true},
	}
	portalB := &v1alpha1.L7Portal{
		ObjectMeta: metav1.ObjectMeta{Name: "portal", Namespace: "ns-b"},
		Spec:       v1alpha1.L7PortalSpec{Enabled: true},
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(portalA, portalB).Build()

	writePortalAPI(t, "ns-a", "portal", "api-a.json", "api-a")
	writePortalAPI(t, "ns-b", "portal", "api-b.json", "api-b")

	syncPortal(ctx, Params{Client: c, Scheme: scheme, Log: logr.Discard(), Instance: portalA})

	cm := &corev1.ConfigMap{}
	if err := c.Get(ctx, types.NamespacedName{Name: "portal-api-summary", Namespace: "ns-a"}, cm); err != nil {
		t.Fatalf("api summary configmap not created: %v", err)
	}
	gz, err := base64.StdEncoding.DecodeString(cm.Data["apis"])
	if err != nil {
		t.Fatal(err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	apis := []templategen.PortalAPI{}
	if err := json.Unmarshal(raw, &apis); err != nil {
		t.Fatal(err)
	}
	if len(apis) != 1 || apis[0].Name != "api-a" {
		t.Errorf("expected only api-a in ns-a summary, actual %+v", apis)
	}
}
