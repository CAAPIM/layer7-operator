package reconcile

import (
	"context"
	policyv1 "k8s.io/api/policy/v1"
	//"context"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewPDB(t *testing.T) {
	t.Run("should create PDB", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := PodDisruptionBudget(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that PDB is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &policyv1.PodDisruptionBudget{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.MaxUnavailable.IntVal != 5 {
			t.Errorf("Expected %d, Actual %d", 5, got.Spec.MaxUnavailable.IntVal)
		}
	})

	t.Run("should update PDB", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := PodDisruptionBudget(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.PodDisruptionBudget.MaxUnavailable.IntVal = 3
		err = PodDisruptionBudget(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that PDB is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &policyv1.PodDisruptionBudget{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.MaxUnavailable.IntVal != 3 {
			t.Errorf("Expected %d, Actual %d", 3, got.Spec.MaxUnavailable.IntVal)
		}
	})
}
