package reconcile

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewService(t *testing.T) {
	// Exercises both main and management Services in one reconcile pass (regression: main service
	// must not return early and skip the management Service).
	t.Run("should create services", func(t *testing.T) {
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

		//verify that management service is created
		nns = types.NamespacedName{Namespace: "default", Name: "test-management-service"}
		got = &corev1.Service{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Ports[0].Port != 9443 {
			t.Errorf("Expected %d, Actual %d", 9443, got.Spec.Ports[0].Port)
		}

	})

	t.Run("should update service", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Services(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.Service.Ports[0].Port = 1234
		err = Services(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that service is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Service{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Ports[0].Port != 1234 {
			t.Errorf("Expected %d, Actual %d", 1234, got.Spec.Ports[0].Port)
		}
	})
}

func TestMetadataKeysMatch(t *testing.T) {
	t.Run("extra keys on current allowed", func(t *testing.T) {
		cur := metav1.ObjectMeta{
			Labels:      map[string]string{"app.kubernetes.io/name": "gw", "extra": "x"},
			Annotations: map[string]string{"a": "b"},
		}
		des := metav1.ObjectMeta{
			Labels:      map[string]string{"app.kubernetes.io/name": "gw"},
			Annotations: map[string]string{"a": "b"},
		}
		if !metadataKeysMatch(cur, des) {
			t.Fatal("expected match")
		}
	})
	t.Run("missing desired label value", func(t *testing.T) {
		cur := metav1.ObjectMeta{Labels: map[string]string{"k": "wrong"}}
		des := metav1.ObjectMeta{Labels: map[string]string{"k": "v"}}
		if metadataKeysMatch(cur, des) {
			t.Fatal("expected no match")
		}
	})
}
