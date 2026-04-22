package gateway

import (
	"reflect"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func gatewayForDeploymentProbeTest() securityv1.Gateway {
	gw := getGatewayWitApp()
	gw.Name = "test"
	gw.Namespace = "testNamespace"
	gw.Spec.App.Image = "gateway:test"
	gw.Spec.App.Replicas = 1
	gw.Spec.App.ServiceAccount.Name = "testServiceAccount"
	gw.Spec.App.ServiceAccount.Create = true
	gw.Spec.App.PodSecurityContext = corev1.PodSecurityContext{}
	gw.Spec.License.SecretName = "license-secret"
	ports := securityv1.Ports{}
	ports.Name = "https"
	ports.Port = 8443
	ports.TargetPort = 8443
	ports.Protocol = "TCP"
	gw.Spec.App.Service = securityv1.Service{}
	gw.Spec.App.Service.Ports = []securityv1.Ports{ports}
	return gw
}

func TestDeploymentStartupProbe(t *testing.T) {
	t.Run("omitted leaves StartupProbe nil", func(t *testing.T) {
		gw := gatewayForDeploymentProbeTest()
		dep := NewDeployment(&gw, "kubernetes")
		c := dep.Spec.Template.Spec.Containers[0]
		if c.StartupProbe != nil {
			t.Fatalf("StartupProbe: want nil, got %+v", c.StartupProbe)
		}
	})
	t.Run("set when spec defines StartupProbe", func(t *testing.T) {
		gw := gatewayForDeploymentProbeTest()
		want := corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/health",
					Port: intstr.FromInt32(8443),
				},
			},
			FailureThreshold: 18,
		}
		gw.Spec.App.StartupProbe = want
		dep := NewDeployment(&gw, "kubernetes")
		c := dep.Spec.Template.Spec.Containers[0]
		if c.StartupProbe == nil {
			t.Fatal("StartupProbe: want non-nil")
		}
		if !reflect.DeepEqual(*c.StartupProbe, want) {
			t.Errorf("StartupProbe mismatch\nwant: %+v\ngot:  %+v", want, *c.StartupProbe)
		}
	})
}

func TestDeploymentWithPorts(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Namespace = "testNamespace"
	gateway.Spec.App.Image = "testImage"
	gateway.Spec.App.Replicas = 5
	gateway.Spec.App.ServiceAccount.Name = "testServiceAccount"
	gateway.Spec.App.ServiceAccount.Create = true
	gateway.Spec.App.PodSecurityContext = corev1.PodSecurityContext{}
	gateway.Spec.App.TopologySpreadConstraints = []corev1.TopologySpreadConstraint{{TopologyKey: "testTopology"}}
	gateway.Spec.App.Tolerations = []corev1.Toleration{{Key: "testKey", Value: "testValue"}}
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
	platform := "kubernetes"

	dep := NewDeployment(&gateway, platform)

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
