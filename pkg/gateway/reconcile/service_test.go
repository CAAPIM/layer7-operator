package reconcile

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("should create secret", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Services(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that service is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Service{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "test" {
			t.Errorf("Expected %s, Actual %s", "test", got.Name)
		}
	})
}
