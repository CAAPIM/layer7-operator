package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHPA(gw *securityv1.Gateway) *autoscalingv2.HorizontalPodAutoscaler {

	behavior := &gw.Spec.App.Autoscaling.HPA.Behavior
	selectPolicyDefault := autoscalingv2.ScalingPolicySelect("Max")

	if behavior.ScaleDown.SelectPolicy == nil {
		behavior.ScaleDown.SelectPolicy = &selectPolicyDefault
	}

	if behavior.ScaleUp.SelectPolicy == nil {
		behavior.ScaleUp.SelectPolicy = &selectPolicyDefault
	}

	hpaSpec := autoscalingv2.HorizontalPodAutoscalerSpec{
		ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       gw.Name,
		},
		MinReplicas: gw.Spec.App.Autoscaling.HPA.MinReplicas,
		MaxReplicas: gw.Spec.App.Autoscaling.HPA.MaxReplicas,
		Metrics:     gw.Spec.App.Autoscaling.HPA.Metrics,
	}

	if gw.Spec.App.Autoscaling.HPA.Behavior != (autoscalingv2.HorizontalPodAutoscalerBehavior{}) {
		hpaSpec.Behavior = behavior
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v2",
			Kind:       "HorizontalPodAutoscaler",
		},
		Spec: hpaSpec,
	}
	return hpa
}
