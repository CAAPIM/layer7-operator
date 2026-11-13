package reconcile

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	gatewaybuilder "github.com/caapim/layer7-operator/pkg/gateway"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestKubernetesGatewayMissingAPIs(t *testing.T) {
	for _, response := range []string{"no match", "endpoint not found"} {
		t.Run(response, func(t *testing.T) {
			scheme := runtime.NewScheme()
			if err := securityv1.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			if err := gatewayv1.Install(scheme); err != nil {
				t.Fatal(err)
			}
			unavailable := func(obj runtime.Object) error {
				t.Helper()
				kinds, _, err := scheme.ObjectKinds(obj)
				if err != nil {
					t.Fatal(err)
				}
				kind := strings.TrimSuffix(kinds[0].Kind, "List")
				if response == "no match" {
					return &meta.NoKindMatchError{GroupKind: schema.GroupKind{Group: gatewayv1.GroupName, Kind: kind}, SearchedVersions: []string{"v1"}}
				}
				resources := map[string]string{"Gateway": "gateways", "HTTPRoute": "httproutes", "TLSRoute": "tlsroutes", "BackendTLSPolicy": "backendtlspolicies"}
				return apierrors.NewNotFound(schema.GroupResource{Group: gatewayv1.GroupName, Resource: resources[kind]}, "")
			}
			// Simulate API-server availability changing without replacing the reconciler
			// or Params. Availability belongs to the server, not to controller state.
			httpAvailable := false
			cl := fake.NewClientBuilder().WithScheme(scheme).WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, cl client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					if _, ok := obj.(*gatewayv1.HTTPRoute); ok && httpAvailable {
						return cl.Get(ctx, key, obj, opts...)
					}
					return unavailable(obj)
				},
				Create: func(ctx context.Context, cl client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
					if _, ok := obj.(*gatewayv1.HTTPRoute); ok && httpAvailable {
						return cl.Create(ctx, obj, opts...)
					}
					return unavailable(obj)
				},
				List: func(ctx context.Context, cl client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
					if _, ok := list.(*gatewayv1.HTTPRouteList); ok && httpAvailable {
						return cl.List(ctx, list, opts...)
					}
					return unavailable(list)
				},
			}).Build()
			p := Params{Client: cl, Scheme: scheme, Instance: &securityv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: "http-only", Namespace: "default", UID: "http-only-owner"}}}
			if err := KubernetesGateway(context.Background(), p); err != nil {
				t.Fatalf("disabled with missing APIs: %v", err)
			}
			p.Instance.Spec.App.KubernetesGateway = securityv1.KubernetesGateway{Enabled: true, Gateway: securityv1.KubernetesGatewayResource{ExistingRef: securityv1.KubernetesGatewayReference{Name: "shared"}}, HTTPRoute: securityv1.KubernetesHTTPRoute{Rules: []securityv1.KubernetesHTTPRouteRule{{Port: 8443}}}, BackendTLSPolicy: securityv1.KubernetesBackendTLSPolicy{Enabled: new(false)}}
			if err := KubernetesGateway(context.Background(), p); err == nil || !strings.Contains(err.Error(), "HTTPRoutes") || (!meta.IsNoMatchError(err) && !apierrors.IsNotFound(err)) {
				t.Fatalf("expected unavailable HTTPRoute API error, got %v", err)
			}
			httpAvailable = true
			if err := KubernetesGateway(context.Background(), p); err != nil {
				t.Fatalf("newly available API must work without cached state: %v", err)
			}
			routes := &gatewayv1.HTTPRouteList{}
			if err := cl.List(context.Background(), routes); err != nil || len(routes.Items) != 1 {
				t.Fatalf("expected one HTTPRoute, got %d routes and error %v", len(routes.Items), err)
			}
			p.Instance.Spec.App.KubernetesGateway.Enabled = false
			if err := KubernetesGateway(context.Background(), p); err != nil {
				t.Fatalf("cleanup with only HTTPRoute available: %v", err)
			}
			if err := cl.List(context.Background(), routes); err != nil || len(routes.Items) != 0 {
				t.Fatalf("expected HTTPRoute cleanup, got %d routes and error %v", len(routes.Items), err)
			}
			httpAvailable = false
			if err := KubernetesGateway(context.Background(), p); err != nil {
				t.Fatalf("disabled after API removal: %v", err)
			}
			listErr := errors.New("API list failed")
			p.Client = interceptor.NewClient(cl, interceptor.Funcs{List: func(context.Context, client.WithWatch, client.ObjectList, ...client.ListOption) error { return listErr }})
			if err := KubernetesGateway(context.Background(), p); !errors.Is(err, listErr) {
				t.Fatalf("cleanup errors must be retried, got %v", err)
			}
		})
	}
}

func TestKubernetesGatewayLifecycle(t *testing.T) {
	// Use the actual CRDs from the pinned Go dependency to exercise validation and defaulting.
	module, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "sigs.k8s.io/gateway-api").Output()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(strings.TrimSpace(string(module)), "config", "crd", "standard")
	paths := []string{}
	for _, name := range []string{"gateways", "httproutes", "tlsroutes", "backendtlspolicies"} {
		paths = append(paths, filepath.Join(dir, "gateway.networking.k8s.io_"+name+".yaml"))
	}
	if _, err := envtest.InstallCRDs(testEnv.Config, envtest.CRDInstallOptions{Paths: paths, ErrorIfPathMissing: true}); err != nil {
		t.Fatal(err)
	}
	if err := gatewayv1.Install(testScheme); err != nil {
		t.Fatal(err)
	}
	// A fresh mapper picks up the CRDs installed after TestMain.
	cl, err := client.New(testEnv.Config, client.Options{Scheme: testScheme})
	if err != nil {
		t.Fatal(err)
	}
	p := newParams()
	p.Client = cl
	p.Instance.Name = "gateway-api-test"
	p.Instance.UID = ""
	p.Instance.Spec.App.KubernetesGateway = securityv1.KubernetesGateway{Enabled: true, Gateway: securityv1.KubernetesGatewayResource{Create: true, GatewayClassName: "contour", TLS: securityv1.KubernetesGatewayTLS{ExistingSecretName: "frontend"}}, HTTPRoute: securityv1.KubernetesHTTPRoute{Rules: []securityv1.KubernetesHTTPRouteRule{{Hostname: "api.example.com", Port: 8443}, {Hostname: "old.example.com", Port: 8443}}}, BackendTLSPolicy: securityv1.KubernetesBackendTLSPolicy{Validation: securityv1.KubernetesBackendTLSValidation{CACertificateRefs: []gatewayv1.LocalObjectReference{{Group: "", Kind: "ConfigMap", Name: "ca"}}}}}
	if err := cl.Create(ctx, p.Instance); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cl.Delete(ctx, p.Instance) })
	run := func() {
		t.Helper()
		if err := KubernetesGateway(ctx, p); err != nil {
			t.Fatal(err)
		}
	}
	run()
	routes := &gatewayv1.HTTPRouteList{}
	if err := cl.List(ctx, routes); err != nil {
		t.Fatal(err)
	}
	if len(routes.Items) != 2 {
		t.Fatalf("expected two routes, got %d", len(routes.Items))
	}
	kg := &gatewayv1.Gateway{}
	key := types.NamespacedName{Name: p.Instance.Name + "-gateway", Namespace: p.Instance.Namespace}
	if err := cl.Get(ctx, key, kg); err != nil {
		t.Fatal(err)
	}
	version := kg.ResourceVersion
	routeVersion := routes.Items[0].ResourceVersion
	policy := &gatewayv1.BackendTLSPolicy{}
	policyKey := types.NamespacedName{Name: p.Instance.Name + "-backendtlspolicy", Namespace: p.Instance.Namespace}
	if err := cl.Get(ctx, policyKey, policy); err != nil {
		t.Fatal(err)
	}
	policyVersion := policy.ResourceVersion
	run()
	if err := cl.Get(ctx, key, kg); err != nil {
		t.Fatal(err)
	}
	if version != kg.ResourceVersion {
		raw, _ := json.Marshal(kg.Spec)
		wanted, _ := gatewaybuilder.NewKubernetesGatewayResources(p.Instance)
		want, _ := json.Marshal(wanted.Gateways[0])
		t.Fatalf("Gateway changed on no-op reconcile: %s -> %s; current=%s desired=%s", version, kg.ResourceVersion, raw, want)
	}
	route := &gatewayv1.HTTPRoute{}
	if err := cl.Get(ctx, client.ObjectKeyFromObject(&routes.Items[0]), route); err != nil {
		t.Fatal(err)
	}
	if routeVersion != route.ResourceVersion {
		t.Fatal("HTTPRoute changed on no-op reconcile")
	}
	if err := cl.Get(ctx, policyKey, policy); err != nil {
		t.Fatal(err)
	}
	if policyVersion != policy.ResourceVersion {
		t.Fatal("BackendTLSPolicy changed on no-op reconcile")
	}
	// Preserve externally injected metadata and controller status while updating a spec.
	kg.Annotations = map[string]string{"external": "keep"}
	if err := cl.Update(ctx, kg); err != nil {
		t.Fatal(err)
	}
	kg.Status.Conditions = []metav1.Condition{{Type: "Accepted", Status: metav1.ConditionTrue, Reason: "Accepted", Message: "accepted", LastTransitionTime: metav1.Now()}}
	if err := cl.Status().Update(ctx, kg); err != nil {
		t.Fatal(err)
	}
	p.Instance.Spec.App.KubernetesGateway.HTTPRoute.Rules = p.Instance.Spec.App.KubernetesGateway.HTTPRoute.Rules[:1]
	run()
	if err := cl.Get(ctx, key, kg); err != nil {
		t.Fatal(err)
	}
	if kg.Annotations["external"] != "keep" || len(kg.Status.Conditions) != 1 {
		t.Fatal("lost controller metadata or status")
	}
	if err := cl.List(ctx, routes); err != nil {
		t.Fatal(err)
	}
	if len(routes.Items) != 1 {
		t.Fatal("obsolete hostname route not removed")
	}
	// Exercise typed spec and metadata updates for HTTPRoute and BackendTLSPolicy.
	p.Instance.Spec.App.KubernetesGateway.HTTPRoute.Rules[0].Port = 8444
	p.Instance.Spec.App.KubernetesGateway.BackendTLSPolicy.Validation.Hostname = "backend.example.com"
	p.Instance.Spec.App.KubernetesGateway.Labels = map[string]string{"revision": "2"}
	run()
	if err := cl.Get(ctx, client.ObjectKeyFromObject(&routes.Items[0]), route); err != nil {
		t.Fatal(err)
	}
	if *route.Spec.Rules[0].BackendRefs[0].Port != 8444 || route.Labels["revision"] != "2" {
		t.Fatal("HTTPRoute spec or metadata was not updated")
	}
	if err := cl.Get(ctx, policyKey, policy); err != nil {
		t.Fatal(err)
	}
	if policy.Spec.Validation.Hostname != "backend.example.com" || policy.Labels["revision"] != "2" {
		t.Fatal("BackendTLSPolicy spec or metadata was not updated")
	}
	// Both endpoints coexist; disabling either leaves the other intact.
	if err := Ingress(ctx, p); err != nil {
		t.Fatal(err)
	}
	p.Instance.Spec.App.Ingress.Enabled = false
	if err := Ingress(ctx, p); err != nil {
		t.Fatal(err)
	}
	if err := cl.Get(ctx, key, kg); err != nil {
		t.Fatalf("disabling Ingress affected Gateway API: %v", err)
	}
	p.Instance.Spec.App.Ingress.Enabled = true
	if err := Ingress(ctx, p); err != nil {
		t.Fatal(err)
	}
	p.Instance.Spec.App.KubernetesGateway.Enabled = false
	run()
	ingress := &networkingv1.Ingress{}
	if err := cl.Get(ctx, client.ObjectKeyFromObject(p.Instance), ingress); err != nil {
		t.Fatalf("disabling Gateway API affected Ingress: %v", err)
	}
	if err := cl.Get(ctx, key, kg); !apierrors.IsNotFound(err) {
		t.Fatalf("managed Gateway not removed: %v", err)
	}
	if err := cl.List(ctx, routes); err != nil || len(routes.Items) != 0 {
		t.Fatal("routes not removed")
	}
	policies := &gatewayv1.BackendTLSPolicyList{}
	if err := cl.List(ctx, policies); err != nil || len(policies.Items) != 0 {
		t.Fatal("policy not removed")
	}
	// Switch to passthrough, then to an external Gateway; external resources remain untouched.
	c := &p.Instance.Spec.App.KubernetesGateway
	c.Enabled = true
	c.HTTPRoute.Enabled = new(false)
	c.BackendTLSPolicy.Enabled = new(false)
	c.TLSRoute = securityv1.KubernetesTLSRoute{Enabled: true, Rules: []securityv1.KubernetesTLSRouteRule{{Hostname: "tls.example.com", Port: 8443}}}
	run()
	tlsRoutes := &gatewayv1.TLSRouteList{}
	if err := cl.List(ctx, tlsRoutes); err != nil || len(tlsRoutes.Items) != 1 {
		t.Fatal("expected one passthrough route")
	}
	tlsKey := client.ObjectKeyFromObject(&tlsRoutes.Items[0])
	tlsVersion := tlsRoutes.Items[0].ResourceVersion
	run()
	tlsRoute := &gatewayv1.TLSRoute{}
	if err := cl.Get(ctx, tlsKey, tlsRoute); err != nil {
		t.Fatal(err)
	}
	if tlsRoute.ResourceVersion != tlsVersion {
		t.Fatal("TLSRoute changed on no-op reconcile")
	}
	c.TLSRoute.Rules[0].Port = 9443
	c.Gateway.Create = false
	c.Gateway.ExistingRef = securityv1.KubernetesGatewayReference{Name: "shared", Namespace: "infra"}
	run()
	if err := cl.Get(ctx, key, kg); !apierrors.IsNotFound(err) {
		t.Fatalf("old managed Gateway not cleaned up: %v", err)
	}
	if err := cl.List(ctx, tlsRoutes); err != nil || len(tlsRoutes.Items) != 1 {
		t.Fatal("expected one passthrough route")
	}
	if *tlsRoutes.Items[0].Spec.ParentRefs[0].Namespace != "infra" {
		t.Fatal("shared Gateway reference not updated")
	}
	if *tlsRoutes.Items[0].Spec.Rules[0].BackendRefs[0].Port != 9443 {
		t.Fatal("TLSRoute backend port was not updated")
	}
	// A resource with a colliding name must not be adopted or deleted.
	c.Gateway.Create = true
	foreign := &gatewayv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: key.Name, Namespace: key.Namespace}, Spec: gatewayv1.GatewaySpec{GatewayClassName: "external", Listeners: []gatewayv1.Listener{{Name: "http", Port: 80, Protocol: gatewayv1.HTTPProtocolType}}}}
	if err := cl.Create(ctx, foreign); err != nil {
		t.Fatal(err)
	}
	if err := KubernetesGateway(ctx, p); err == nil || !strings.Contains(err.Error(), "not owned") {
		t.Fatalf("expected ownership conflict, got %v", err)
	}
	c.Enabled = false
	run()
	if err := cl.Get(ctx, key, foreign); err != nil {
		t.Fatalf("deleted unowned Gateway: %v", err)
	}

	t.Run("HTTP rules span multiple routes", func(t *testing.T) {
		p := newParams()
		p.Client = cl
		p.Instance.Name = "many-http-rules"
		p.Instance.UID = ""
		c := &p.Instance.Spec.App.KubernetesGateway
		*c = securityv1.KubernetesGateway{
			Enabled:          true,
			Gateway:          securityv1.KubernetesGatewayResource{ExistingRef: securityv1.KubernetesGatewayReference{Name: "shared"}},
			BackendTLSPolicy: securityv1.KubernetesBackendTLSPolicy{Enabled: new(false)},
		}
		for _, host := range []gatewayv1.Hostname{"api.example.com", "other.example.com"} {
			for range 16 {
				c.HTTPRoute.Rules = append(c.HTTPRoute.Rules, securityv1.KubernetesHTTPRouteRule{
					Hostname: host, Port: 8443,
					Matches: []gatewayv1.HTTPRouteMatch{{Path: &gatewayv1.HTTPPathMatch{Type: new(gatewayv1.PathMatchPathPrefix), Value: new("/api")}}},
					Filters: []gatewayv1.HTTPRouteFilter{{Type: gatewayv1.HTTPRouteFilterRequestHeaderModifier,
						RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{Set: []gatewayv1.HTTPHeader{{Name: "X-Test", Value: "retained"}}}}},
				})
			}
		}
		if err := cl.Create(ctx, p.Instance); err != nil {
			t.Fatalf("Layer7 CRD must accept more than 16 rules: %v", err)
		}
		t.Cleanup(func() { _ = cl.Delete(ctx, p.Instance) })
		// Read back to exercise schema pruning/defaulting before reconciliation.
		if err := cl.Get(ctx, client.ObjectKeyFromObject(p.Instance), p.Instance); err != nil {
			t.Fatal(err)
		}
		if err := KubernetesGateway(ctx, p); err != nil {
			t.Fatalf("generated routes must pass upstream admission: %v", err)
		}
		desired, err := gatewaybuilder.NewKubernetesGatewayResources(p.Instance)
		if err != nil {
			t.Fatal(err)
		}
		if len(desired.HTTPRoutes) != 2 {
			t.Fatalf("expected two routes, got %d", len(desired.HTTPRoutes))
		}
		for _, want := range desired.HTTPRoutes {
			actual := &gatewayv1.HTTPRoute{}
			if err := cl.Get(ctx, client.ObjectKeyFromObject(want), actual); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = cl.Delete(ctx, actual) })
			if len(actual.Spec.Rules) != 16 || *actual.Spec.Rules[0].Matches[0].Path.Value != "/api" || actual.Spec.Rules[0].Filters[0].RequestHeaderModifier.Set[0].Value != "retained" {
				t.Fatalf("rules, matches or filters were not preserved: %+v", actual.Spec.Rules)
			}
		}
		// Inherited CEL is still enforced by the actual HTTPRoute API.
		p.Instance.Spec.App.KubernetesGateway.HTTPRoute.Rules[0].Matches[0].Path.Value = new("relative-path")
		if err := KubernetesGateway(ctx, p); !apierrors.IsInvalid(err) {
			t.Fatalf("expected upstream rejection of invalid path, got %v", err)
		}
	})
}
