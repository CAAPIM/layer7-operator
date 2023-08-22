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

		params := newParams()
		k8sClient, err = client.New(cfg, client.Options{Scheme: params.Scheme})
		if err != nil {
			fmt.Printf("failed to setup a Kubernetes client: %v", err)
			os.Exit(1)
		}
		params.Client = k8sClient
		err = Secret(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that secret is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Secret{}
		err = params.Client.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
	})
}
