package hpa

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHPA(gw *securityv1.Gateway) *autoscalingv2.HorizontalPodAutoscaler {

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
		hpaSpec.Behavior = &gw.Spec.App.Autoscaling.HPA.Behavior
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
