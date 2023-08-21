package reconcile

import (
	"context"
	"testing"
)

func TestNewSecret(t *testing.T) {
	t.Run("should create secret", func(t *testing.T) {
		ctx := context.Background()

		params, err := newParams()
		params.Instance.Name = "test"
		params.Instance.Namespace = "default"
		if err != nil {
			t.Fatal(err)
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
