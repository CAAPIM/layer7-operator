package util

import (
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
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
	bundleBytes, sha1, err := BuildDefaultListenPortBundle(true)
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
	port.Port = 9090
	gateway.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{port}

	bundleBytes, sha1, err := BuildCustomListenPortBundle(&gateway, false)
	if err != nil {
		t.Errorf("Error getting default listen port bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "9090") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "9090")
	}
}

func TestListenPortNamesExcludedFromGraphmanSync(t *testing.T) {
	gw := securityv1.Gateway{}
	gw.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{
		{Name: "a", Port: 8443},
		{Name: "graphman-lp", Port: 9443},
	}
	ex := ListenPortNamesExcludedFromGraphmanSync(&gw, 9443)
	if _, ok := ex["graphman-lp"]; !ok {
		t.Fatalf("expected graphman-lp excluded for port 9443, got %v", ex)
	}
	if _, ok := ex["a"]; ok {
		t.Fatalf("did not expect port 8443 in excluded set")
	}
}

func TestFilterListenPortBundleForGraphmanSync(t *testing.T) {
	b := graphman.Bundle{
		ListenPorts: []*graphman.ListenPortInput{
			{Name: "keep", Port: 8443},
			{Name: "drop", Port: 9443},
		},
	}
	FilterListenPortBundleForGraphmanSync(&b, 9443)
	if len(b.ListenPorts) != 1 || b.ListenPorts[0].Name != "keep" {
		t.Fatalf("expected single listen port keep, got %+v", b.ListenPorts)
	}
}
