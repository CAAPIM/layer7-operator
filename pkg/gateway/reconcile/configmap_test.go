package reconcile

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestConfigMap(t *testing.T) {
	t.Run("should create configmap", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := ConfigMaps(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that ConfigMap is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.ConfigMap{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "test" {
			t.Errorf("Expected %s, Actual %s", "test", got.Name)
		}
	})
}
