package reconcile

import (
	"context"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	//"context"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewHpa(t *testing.T) {
	t.Run("should create hpa", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := HorizontalPodAutoscaler(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Hpa is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &autoscalingv2.HorizontalPodAutoscaler{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.MaxReplicas != 3 {
			t.Errorf("Expected %d, Actual %d", 3, got.Spec.MaxReplicas)
		}
	})

	t.Run("should update hpa", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := HorizontalPodAutoscaler(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.Autoscaling.HPA.MaxReplicas = 5
		err = HorizontalPodAutoscaler(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that HorizontalPodAutoscaler is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &autoscalingv2.HorizontalPodAutoscaler{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.MaxReplicas != 5 {
			t.Errorf("Expected %d, Actual %d", 5, got.Spec.MaxReplicas)
		}
	})
}
