package gateway

import (
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewHPA(t *testing.T) {
	minReplicas := int32(3)

	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			License: securityv1.License{
				Accept: true,
			},
			App: securityv1.App{
				Autoscaling: securityv1.Autoscaling{
					Enabled: true,
					HPA: securityv1.HPA{
						MinReplicas: &minReplicas,
						MaxReplicas: 5,
						Metrics:     []autoscalingv2.MetricSpec{},
						Behavior: autoscalingv2.HorizontalPodAutoscalerBehavior{
							ScaleUp:   &autoscalingv2.HPAScalingRules{},
							ScaleDown: &autoscalingv2.HPAScalingRules{},
						},
					},
				},
			},
		},
	}

	hpa := NewHPA(&gateway)
	if *hpa.Spec.Behavior.ScaleUp.SelectPolicy != autoscalingv2.ScalingPolicySelect("Max") {
		t.Errorf("expected %s, actual %s", autoscalingv2.ScalingPolicySelect("Max"), *hpa.Spec.Behavior.ScaleUp.SelectPolicy)
	}

	if *hpa.Spec.Behavior.ScaleDown.SelectPolicy != autoscalingv2.ScalingPolicySelect("Max") {
		t.Errorf("expected %s, actual %s", autoscalingv2.ScalingPolicySelect("Max"), *hpa.Spec.Behavior.ScaleDown.SelectPolicy)
	}

	if *hpa.Spec.MinReplicas != 3 {
		t.Errorf("expected %d, actual %d", 3, hpa.Spec.MinReplicas)
	}

	if hpa.Spec.MaxReplicas != 5 {
		t.Errorf("expected %d, actual %d", 5, hpa.Spec.MaxReplicas)
	}
}
