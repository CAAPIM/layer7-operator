package reconcile

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/types"
)

func TestNewDeployment(t *testing.T) {
	t.Run("should create Deployment", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Deployment is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &appsv1.Deployment{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if *got.Spec.Replicas != int32(5) {
			t.Errorf("Expected %d, Actual %d", int32(5), *got.Spec.Replicas)
		}
	})

	t.Run("should update Deployment", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.ServiceAccount.Name = "modified"
		err = Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Deployment is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &appsv1.Deployment{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Template.Spec.ServiceAccountName != "modified" {
			t.Errorf("Expected %s, Actual %s", "modified", got.Spec.Template.Spec.ServiceAccountName)
		}
	})
}
