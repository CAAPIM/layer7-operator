package gateway

import (
	"reflect"
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func gatewayAPIExample() *securityv1.Gateway {
	return &securityv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"}, Spec: securityv1.GatewaySpec{App: securityv1.App{KubernetesGateway: securityv1.KubernetesGateway{
		Enabled: true, Gateway: securityv1.KubernetesGatewayResource{Create: true, GatewayClassName: "contour", TLS: securityv1.KubernetesGatewayTLS{ExistingSecretName: "frontend"}},
		HTTPRoute:        securityv1.KubernetesHTTPRoute{Rules: []securityv1.KubernetesHTTPRouteRule{{Hostname: "api.example.com", Port: 8443}, {Hostname: "api.example.com", Port: 8443}, {Port: 8443}}},
		BackendTLSPolicy: securityv1.KubernetesBackendTLSPolicy{Validation: securityv1.KubernetesBackendTLSValidation{CACertificateRefs: []gatewayv1.LocalObjectReference{{Group: "", Kind: "ConfigMap", Name: "backend-ca"}}}},
	}}}}
}

func TestKubernetesGatewayResources(t *testing.T) {
	gw := gatewayAPIExample()
	original := gw.DeepCopy()
	objects, err := NewKubernetesGatewayResources(gw)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, gw) {
		t.Fatal("builder mutated the CR")
	}
	if len(objects.Gateways) != 1 || len(objects.HTTPRoutes) != 2 || len(objects.BackendTLSPolicies) != 1 || len(objects.TLSRoutes) != 0 {
		t.Fatalf("expected Gateway, two HTTPRoutes and policy; got %+v", objects)
	}
	kg := objects.Gateways[0]
	if len(kg.Spec.Listeners) != 2 || kg.Spec.Listeners[0].Hostname != nil {
		t.Fatalf("unexpected listeners: %+v", kg.Spec.Listeners)
	}
	for _, l := range kg.Spec.Listeners {
		if l.TLS.CertificateRefs[0].Name != "frontend" || *l.TLS.Mode != gatewayv1.TLSModeTerminate {
			t.Fatal("wrong listener TLS")
		}
	}
	rules := 0
	for _, r := range objects.HTTPRoutes {
		rules += len(r.Spec.Rules)
		if r.Spec.ParentRefs[0].Name != "test-gateway" {
			t.Fatal("wrong parent")
		}
		if r.Spec.Rules[0].BackendRefs[0].Name != "test" {
			t.Fatal("wrong backend")
		}
	}

	if rules != 3 {
		t.Fatalf("lost grouped rules: %d", rules)
	}
	p := objects.BackendTLSPolicies[0]
	if p.Spec.Validation.Hostname != "test.default.svc.cluster.local" || p.Spec.Validation.CACertificateRefs[0].Name != "backend-ca" {
		t.Fatal("wrong validation")
	}
}

func TestKubernetesGatewayPassthroughAndShared(t *testing.T) {
	gw := gatewayAPIExample()
	c := &gw.Spec.App.KubernetesGateway
	c.HTTPRoute.Enabled = new(false)
	c.BackendTLSPolicy.Enabled = new(false)
	c.TLSRoute = securityv1.KubernetesTLSRoute{Enabled: true, Rules: []securityv1.KubernetesTLSRouteRule{{Hostname: "tls.example.com", Port: 8443}}}
	objs, err := NewKubernetesGatewayResources(gw)
	if err != nil {
		t.Fatal(err)
	}
	l := objs.Gateways[0].Spec.Listeners[0]
	if l.Protocol != gatewayv1.TLSProtocolType || *l.TLS.Mode != gatewayv1.TLSModePassthrough || len(l.TLS.CertificateRefs) != 0 {
		t.Fatal("incorrect passthrough listener")
	}
	c.Gateway.Create = false
	c.Gateway.ExistingRef = securityv1.KubernetesGatewayReference{Name: "shared", Namespace: "infra"}
	c.HTTPRoute.Enabled = new(true)
	objs, err = NewKubernetesGatewayResources(gw)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs.Gateways) != 0 {
		t.Fatal("created external Gateway")
	}
	for _, route := range objs.HTTPRoutes {
		p := route.Spec.ParentRefs[0]
		if p.Name != "shared" || *p.Namespace != "infra" {
			t.Fatal("wrong shared HTTPRoute parent")
		}
	}
	for _, route := range objs.TLSRoutes {
		p := route.Spec.ParentRefs[0]
		if p.Name != "shared" || *p.Namespace != "infra" {
			t.Fatal("wrong shared TLSRoute parent")
		}
	}

}

func TestKubernetesGatewayManagementAndPolicyPorts(t *testing.T) {
	gw := gatewayAPIExample()
	gw.Spec.App.Management.Service.Enabled = true
	gw.Spec.App.Service.Ports = []securityv1.Ports{{Name: "https", Port: 8443}, {Name: "http", Port: 8080}}
	gw.Spec.App.Management.Service.Ports = []securityv1.Ports{{Name: "management", Port: 9443}}
	c := &gw.Spec.App.KubernetesGateway
	c.HTTPRoute.Rules = []securityv1.KubernetesHTTPRouteRule{{Hostname: "management.example.com", Backend: "management", Port: 9443}}
	c.BackendTLSPolicy.ExcludePorts = []int32{8080}
	objs, err := NewKubernetesGatewayResources(gw)
	if err != nil {
		t.Fatal(err)
	}
	r := objs.HTTPRoutes[0]
	if r.Spec.Rules[0].BackendRefs[0].Name != "test-management-service" {
		t.Fatal("wrong management service")
	}
	p := objs.BackendTLSPolicies[0]
	if len(p.Spec.TargetRefs) != 2 || *p.Spec.TargetRefs[0].SectionName != "https" || p.Spec.TargetRefs[1].Name != "test-management-service" {
		t.Fatalf("incorrect targets: %+v", p.Spec.TargetRefs)
	}
}

func TestKubernetesGatewayValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*securityv1.KubernetesGateway)
	}{
		{"missing class", func(c *securityv1.KubernetesGateway) { c.Gateway.GatewayClassName = "" }},
		{"missing certificate", func(c *securityv1.KubernetesGateway) { c.Gateway.TLS.ExistingSecretName = "" }},
		{"missing parent", func(c *securityv1.KubernetesGateway) { c.Gateway.Create = false }},
		{"missing CA", func(c *securityv1.KubernetesGateway) { c.BackendTLSPolicy.Validation.CACertificateRefs = nil }},
		{"conflicting CA", func(c *securityv1.KubernetesGateway) {
			c.BackendTLSPolicy.Validation.WellKnownCACertificates = new(gatewayv1.WellKnownCACertificatesSystem)
		}},
		{"invalid CA kind", func(c *securityv1.KubernetesGateway) {
			c.BackendTLSPolicy.Validation.CACertificateRefs[0].Kind = "Secret"
		}},
		{"both modes managed", func(c *securityv1.KubernetesGateway) { c.TLSRoute.Enabled = true }},
		{"missing HTTP rules", func(c *securityv1.KubernetesGateway) { c.HTTPRoute.Rules = nil }},
		{"invalid port", func(c *securityv1.KubernetesGateway) { c.HTTPRoute.Rules[0].Port = 0 }},
		{"invalid backend", func(c *securityv1.KubernetesGateway) { c.HTTPRoute.Rules[0].Backend = "other" }},
		{"management disabled", func(c *securityv1.KubernetesGateway) { c.HTTPRoute.Rules[0].Backend = "management" }},
		{"missing TLS hostname", func(c *securityv1.KubernetesGateway) {
			c.HTTPRoute.Enabled = new(false)
			c.TLSRoute.Enabled = true
			c.TLSRoute.Rules = []securityv1.KubernetesTLSRouteRule{{Port: 8443}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw := gatewayAPIExample()
			tt.change(&gw.Spec.App.KubernetesGateway)
			if _, err := NewKubernetesGatewayResources(gw); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
	gw := gatewayAPIExample()
	gw.Spec.App.KubernetesGateway.Enabled = false
	gw.Spec.App.KubernetesGateway.Gateway = securityv1.KubernetesGatewayResource{}
	if objs, err := NewKubernetesGatewayResources(gw); err != nil || len(objs.Gateways) != 0 || len(objs.HTTPRoutes) != 0 || len(objs.TLSRoutes) != 0 || len(objs.BackendTLSPolicies) != 0 {
		t.Fatal("disabled configuration should be ignored")
	}
}

func TestKubernetesGatewayHTTPRouteLimits(t *testing.T) {
	for _, host := range []gatewayv1.Hostname{"api.example.com", ""} {
		t.Run(string(host), func(t *testing.T) {
			gw := gatewayAPIExample()
			c := &gw.Spec.App.KubernetesGateway
			c.HTTPRoute.Rules = nil
			for range 16 {
				c.HTTPRoute.Rules = append(c.HTTPRoute.Rules,
					securityv1.KubernetesHTTPRouteRule{Hostname: host, Port: 8443},
					securityv1.KubernetesHTTPRouteRule{Hostname: "other.example.com", Port: 8443})
			}
			objects, err := NewKubernetesGatewayResources(gw)
			if err != nil {
				t.Fatalf("32 rules across two hosts must be accepted: %v", err)
			}
			if len(objects.HTTPRoutes) != 2 {
				t.Fatalf("expected two routes, got %d", len(objects.HTTPRoutes))
			}
			for _, route := range objects.HTTPRoutes {
				if len(route.Spec.Rules) != 16 {
					t.Fatalf("expected 16 rules per route, got %d", len(route.Spec.Rules))
				}
			}
			c.HTTPRoute.Rules = append(c.HTTPRoute.Rules, securityv1.KubernetesHTTPRouteRule{Hostname: host, Port: 8443})
			if _, err := NewKubernetesGatewayResources(gw); err == nil || !strings.Contains(err.Error(), "16 rules per HTTPRoute") {
				t.Fatalf("expected per-route rule limit error, got %v", err)
			}
		})
	}

	t.Run("matches per route", func(t *testing.T) {
		gw := gatewayAPIExample()
		c := &gw.Spec.App.KubernetesGateway
		c.HTTPRoute.Rules = nil
		for range 2 {
			c.HTTPRoute.Rules = append(c.HTTPRoute.Rules, securityv1.KubernetesHTTPRouteRule{
				Hostname: "api.example.com", Port: 8443, Matches: make([]gatewayv1.HTTPRouteMatch, 64),
			})
		}
		c.HTTPRoute.Rules = append(c.HTTPRoute.Rules, securityv1.KubernetesHTTPRouteRule{Hostname: "other.example.com", Port: 8443})
		if _, err := NewKubernetesGatewayResources(gw); err != nil {
			t.Fatalf("match limit must apply per route: %v", err)
		}
		c.HTTPRoute.Rules[2].Hostname = "api.example.com"
		if _, err := NewKubernetesGatewayResources(gw); err == nil || !strings.Contains(err.Error(), "128 matches per HTTPRoute") {
			t.Fatalf("expected default match to count towards route limit, got %v", err)
		}
	})
}

func TestGatewayAPINames(t *testing.T) {
	if routeName("test", "httproute", "*.example.com") == routeName("test", "httproute", "wildcard.example.com") {
		t.Fatal("hostname collision")
	}
	a := routeName("test", "httproute", "a.b.example.com")
	b := routeName("test", "httproute", "a-b.example.com")
	if a == b {
		t.Fatal("punctuation collision")
	}
	if len(routeName(strings.Repeat("a", 253), "httproute", "example.com")) > 253 {
		t.Fatal("name exceeds limit")
	}
}
