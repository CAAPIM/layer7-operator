package util

import (
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
)

func TestBuildCWPBundle(t *testing.T) {
	cwps := []securityv1.Property{{Name: "test1", Value: "value1"}}
	bundleBytes, sha1, err := BuildCWPBundle(cwps)
	if err != nil {
		t.Errorf("Error getting cwp bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "test1") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "test1")
	}
}

func TestBuildDefaultListenPortBundle(t *testing.T) {
	bundleBytes, sha1, err := BuildDefaultListenPortBundle()
	if err != nil {
		t.Errorf("Error getting default listen port bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "8443") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "8443")
	}
	if !strings.Contains(bundle, "9443") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "9443")
	}

	if !strings.Contains(bundle, "8080") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "8080")
	}
}

func TestBuildCustomListenPortBundle(t *testing.T) {
	gateway := securityv1.Gateway{}
	gateway.Spec = securityv1.GatewaySpec{}
	gateway.Spec.App = securityv1.App{}
	gateway.Spec.App.ListenPorts = securityv1.ListenPorts{}
	gateway.Spec.App.ListenPorts.Custom = securityv1.CustomListenPort{Enabled: true} //ecurityv1.CustomListenPort{true}

	port := securityv1.ListenPort{}
	port.Name = "custom1"
	port.Enabled = true
	port.Protocol = "http"
	port.Port = "9090"
	gateway.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{port}

	bundleBytes, sha1, err := BuildCustomListenPortBundle(&gateway)
	if err != nil {
		t.Errorf("Error getting default listen port bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "9090") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "9090")
	}

}
