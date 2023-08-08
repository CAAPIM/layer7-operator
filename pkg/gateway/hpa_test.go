package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"testing"
)

func TestNewHPA(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Spec.App.Autoscaling = securityv1.Autoscaling{}
	gateway.Spec.App.Autoscaling.HPA = securityv1.HPA{}
	minReplicas := int32(3)
	gateway.Spec.App.Autoscaling.HPA.MinReplicas = &minReplicas
	gateway.Spec.App.Autoscaling.HPA.MaxReplicas = 5
	gateway.Spec.App.Autoscaling.HPA.Metrics = []autoscalingv2.MetricSpec{}
	gateway.Spec.App.Autoscaling.HPA.Behavior = autoscalingv2.HorizontalPodAutoscalerBehavior{}
	gateway.Spec.App.Autoscaling.HPA.Behavior.ScaleDown = &autoscalingv2.HPAScalingRules{}
	gateway.Spec.App.Autoscaling.HPA.Behavior.ScaleUp = &autoscalingv2.HPAScalingRules{}

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
