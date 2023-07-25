package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewPDB - Pod Disruption Budget
func NewPDB(gw *securityv1.Gateway) *policyv1.PodDisruptionBudget {
	labels := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)

	pdb := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1",
			Kind:       "PodDisruptionBudget",
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	if gw.Spec.App.PodDisruptionBudget.MinAvailable != (intstr.IntOrString{}) {
		pdb.Spec.MinAvailable = &gw.Spec.App.PodDisruptionBudget.MinAvailable
	}

	if gw.Spec.App.PodDisruptionBudget.MaxUnavailable != (intstr.IntOrString{}) {
		pdb.Spec.MaxUnavailable = &gw.Spec.App.PodDisruptionBudget.MaxUnavailable
	}

	return pdb
}
