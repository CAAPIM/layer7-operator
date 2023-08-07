package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"testing"
)

func TestNewIngress(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	ports := securityv1.Ports{}
	ports.Name = "httpPort"
	ports.Port = 8443
	ports.TargetPort = 8443
	ports.Protocol = "http"

	gateway.Spec.App.Service = securityv1.Service{}
	gateway.Spec.App.Service.Ports = []securityv1.Ports{ports}

	gateway.Spec.App.Ingress = securityv1.Ingress{}
	gateway.Spec.App.Ingress.Annotations = map[string]string{"test1": "testAnnotation"}
	gateway.Spec.App.Ingress.IngressClassName = "ingressClass"
	ingressRule := networkingv1.IngressRule{}
	ingressRule.Host = "testing.com"
	gateway.Spec.App.Ingress.Rules = []networkingv1.IngressRule{ingressRule}

	ingress := NewIngress(&gateway)
	if ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name != gateway.Spec.App.Service.Ports[0].Name {
		t.Errorf("expected %s, actual %s", ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name, gateway.Spec.App.Service.Ports[0].Name)
	}

}
