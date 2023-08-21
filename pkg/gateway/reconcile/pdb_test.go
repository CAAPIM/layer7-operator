package reconcile

import (
	"context"
	"testing"
)

func TestNewPDB(t *testing.T) {
	t.Run("should create PDB", func(t *testing.T) {
		ctx := context.Background()

		params, err := newParams()
		params.Instance.Name = "test"
		params.Instance.Namespace = "default"
		if err != nil {
			t.Fatal(err)
		}
		err = PodDisruptionBudget(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
	})
}
