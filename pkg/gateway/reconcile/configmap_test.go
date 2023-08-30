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
		if got.Data["SSG_CLUSTER_HOST"] != "testHost" {
			t.Errorf("Expected %s, Actual %s", "testHost", got.Data["SSG_CLUSTER_HOST"])
		}
	})

	t.Run("should update configmap", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := ConfigMaps(ctx, params)
		if err != nil {
			t.Fatal(err)
		}

		params.Instance.Spec.App.Management.Cluster.Hostname = "testing.com"
		ctx = context.Background()
		err = ConfigMaps(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that ConfigMap is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.ConfigMap{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Data["SSG_CLUSTER_HOST"] != "testing.com" {
			t.Errorf("Expected %s, Actual %s", "testing.com", got.Data["SSG_CLUSTER_HOST"])
		}
	})
}
