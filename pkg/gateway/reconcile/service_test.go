package reconcile

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"testing"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

func TestNewService(t *testing.T) {
	t.Run("should create secret", func(t *testing.T) {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
			ErrorIfCRDPathMissing: true,
		}
		cfg, err := testEnv.Start()
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()

		params := newParams()
		k8sClient, err = client.New(cfg, client.Options{Scheme: params.Scheme})
		if err != nil {
			fmt.Printf("failed to setup a Kubernetes client: %v", err)
			os.Exit(1)
		}
		params.Client = k8sClient
		err = Services(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that service is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Service{}
		err = params.Client.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "test" {
			t.Errorf("Expected %s, Actual %s", "test", got.Name)
		}
	})
}

/*var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}
	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})*/
