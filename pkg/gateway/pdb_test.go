package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestNewPDB(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.PodDisruptionBudget = securityv1.PodDisruptionBudgetSpec{}
	minIntOrString := intstr.IntOrString{}
	minIntOrString.IntVal = 3
	gateway.Spec.App.PodDisruptionBudget.MinAvailable = minIntOrString
	maxIntOrString := intstr.IntOrString{}
	maxIntOrString.IntVal = 2
	gateway.Spec.App.PodDisruptionBudget.MaxUnavailable = maxIntOrString
	pdb := NewPDB(&gateway)
	if *pdb.Spec.MinAvailable != gateway.Spec.App.PodDisruptionBudget.MinAvailable {
		t.Errorf("expected %d", gateway.Spec.App.PodDisruptionBudget.MinAvailable.IntVal)
	}
	if *pdb.Spec.MaxUnavailable != gateway.Spec.App.PodDisruptionBudget.MaxUnavailable {
		t.Errorf("expected %d", gateway.Spec.App.PodDisruptionBudget.MaxUnavailable.IntVal)
	}
}
