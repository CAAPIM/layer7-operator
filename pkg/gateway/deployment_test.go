package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestDeploymentWithPorts(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Namespace = "testNamespace"
	gateway.Spec.App.Image = "testImage"
	gateway.Spec.App.Replicas = 5
	gateway.Spec.App.ServiceAccountName = "testServiceAccount"
	gateway.Spec.App.PodSecurityContext = corev1.PodSecurityContext{}
	gateway.Spec.App.TopologySpreadConstraints = []corev1.TopologySpreadConstraint{
		corev1.TopologySpreadConstraint{TopologyKey: "testTopology"}}
	gateway.Spec.App.Tolerations = []corev1.Toleration{corev1.Toleration{Key: "testKey", Value: "testValue"}}
	gateway.Spec.App.NodeSelector = map[string]string{"testLabel": "testValue"}

	ports := securityv1.Ports{}
	ports.Name = "httpPort"
	ports.Port = 8443
	ports.TargetPort = 8443
	ports.Protocol = "http"

	gateway.Spec.App.Service = securityv1.Service{}
	gateway.Spec.App.Service.Ports = []securityv1.Ports{ports}

	gateway.Spec.App.Management.Service.Enabled = true
	managementPorts := securityv1.Ports{}
	managementPorts.Name = "httpPort"
	managementPorts.Port = 9443
	managementPorts.TargetPort = 9443
	managementPorts.Protocol = "http"
	gateway.Spec.App.Management.Service.Ports = []securityv1.Ports{managementPorts}

	dep := NewDeployment(&gateway)

	if dep.ObjectMeta.Namespace != "testNamespace" {
		t.Errorf("expected %s, actual %s", "testNamespace", dep.ObjectMeta.Namespace)
	}

	if dep.ObjectMeta.Name != "test" {
		t.Errorf("expected %s, actual %s", "test", dep.ObjectMeta.Name)
	}

	if *dep.Spec.Replicas != 5 {
		t.Errorf("expected %d, actual %d", 5, *dep.Spec.Replicas)
	}

	if dep.Spec.Template.Spec.ServiceAccountName != "testServiceAccount" {
		t.Errorf("expected %s, actual %s", "testServiceAccount", dep.Spec.Template.Spec.ServiceAccountName)
	}

	if dep.Spec.Template.Spec.TopologySpreadConstraints[0].TopologyKey != "testTopology" {
		t.Errorf("expected %s, actual %s", "testTopology", dep.Spec.Template.Spec.TopologySpreadConstraints[0].TopologyKey)
	}
	if dep.Spec.Template.Spec.Tolerations[0].Key != "testKey" {
		t.Errorf("expected %s, actual %s", "testKey", dep.Spec.Template.Spec.Tolerations[0].Key)
	}
	if dep.Spec.Template.Spec.Tolerations[0].Value != "testValue" {
		t.Errorf("expected %s, actual %s", "testValue", dep.Spec.Template.Spec.Tolerations[0].Value)
	}
	if dep.Spec.Template.Spec.NodeSelector["testLabel"] != "testValue" {
		t.Errorf("expected %s, actual %s", "testValue", dep.Spec.Template.Spec.NodeSelector["testLabel"])
	}
	if dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort != 8443 {
		t.Errorf("expected %d, actual %d", 8443, dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
	}
	if dep.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort != 9443 {
		t.Errorf("expected %d, actual %d", 9443, dep.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort)
	}
}
