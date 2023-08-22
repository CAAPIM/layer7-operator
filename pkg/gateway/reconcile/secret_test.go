package reconcile

import (
	"context"
	"fmt"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"testing"
)

func TestNewSecret(t *testing.T) {
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

		params, err := newParams()
		params.Instance.Name = "test"
		params.Instance.Namespace = "default"
		params.Scheme.AddKnownTypes(securityv1.GroupVersion, params.Instance)
		k8sClient, err = client.New(cfg, client.Options{Scheme: params.Scheme})
		params.Client = k8sClient
		if err != nil {
			fmt.Printf("failed to setup a Kubernetes client: %v", err)
			os.Exit(1)
		}
		err = Secret(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		/*nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Secret{}
		err = params.Client.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}*/
	})
}
