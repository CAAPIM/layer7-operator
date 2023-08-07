package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"testing"
)

func TestNewService(t *testing.T) {
	gateway := getGatewayWitApp()
	ports := securityv1.Ports{}
	ports.Name = "httpPort"
	ports.Port = 443
	ports.TargetPort = 8443
	ports.Protocol = "http"

	gateway.Spec.App.Service = securityv1.Service{}
	gateway.Spec.App.Service.Ports = []securityv1.Ports{ports}
	service := NewService(&gateway)

	if service.Spec.Ports[0].Port != 443 {
		t.Errorf("expected %d, actual %d", 443, service.Spec.Ports[0].Port)
	}
}

func TestNewManagementService(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	ports := securityv1.Ports{}
	ports.Name = "httpPort"
	ports.Port = 9443
	ports.TargetPort = 9443
	ports.Protocol = "http"
	gateway.Spec.App.Management = securityv1.Management{}
	gateway.Spec.App.Management.Service = securityv1.Service{}
	gateway.Spec.App.Management.Service.Ports = []securityv1.Ports{ports}
	service := NewManagementService(&gateway)

	if service.Spec.Ports[0].Port != 9443 {
		t.Errorf("expected %d, actual %d", 9443, service.Spec.Ports[0].Port)
	}

	if service.Name != "test-management-service" {
		t.Errorf("expected %s, actual %s", "test-management-service", service.Name)
	}
}
