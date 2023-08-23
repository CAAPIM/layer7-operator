package reconcile

import (
	"context"
	networkingv1 "k8s.io/api/networking/v1"
	//"context"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestNewIngress(t *testing.T) {
	t.Run("should create Ingress", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Ingress(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Ingress is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &networkingv1.Ingress{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Rules[0].Host != "localhost" {
			t.Errorf("Expected %s, Actual %s", "localhost", got.Spec.Rules[0].Host)
		}
	})

	/*t.Run("should update Ingress", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Ingress(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.Ingress.Rules[0].Host = "testing.com"
		err = Ingress(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Ingress is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &networkingv1.Ingress{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Rules[0].Host != "testing.com" {
			t.Errorf("Expected %s, Actual %s", "testing.com", got.Spec.Rules[0].Host)
		}
	})*/
}
