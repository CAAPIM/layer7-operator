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
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if string(got.Data["USERNAME"]) != params.Instance.Spec.Auth.Username {
			t.Errorf("expected %s, actual %s", params.Instance.Spec.Auth.Username, string(got.Data["USERNAME"]))
		}
		if string(got.Data["PASSWORD"]) != params.Instance.Spec.Auth.Password {
			t.Errorf("expected %s, actual %s", params.Instance.Spec.Auth.Password, string(got.Data["PASSWORD"]))
		}

		if string(got.Data["TOKEN"]) != params.Instance.Spec.Auth.Token {
			t.Errorf("expected %s, actual %s", params.Instance.Spec.Auth.Token, string(got.Data["TOKEN"]))
		}
	})
}
