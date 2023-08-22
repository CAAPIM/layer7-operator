package reconcile

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewSecret(t *testing.T) {
	t.Run("should create secret", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Secret(ctx, params)
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
