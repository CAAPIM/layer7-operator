// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package gateway

import (
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// KubernetesGatewayResources contains the typed resources managed for a Layer7 Gateway.
type KubernetesGatewayResources struct {
	Gateways           []*gatewayv1.Gateway
	HTTPRoutes         []*gatewayv1.HTTPRoute
	TLSRoutes          []*gatewayv1.TLSRoute
	BackendTLSPolicies []*gatewayv1.BackendTLSPolicy
}

// gatewayAPIName bounds generated names and retains a hash to distinguish wildcard,
// punctuation and long hostname variants that otherwise normalize to the same name.
func gatewayAPIName(base, suffix string) string {
	name := base + "-" + suffix
	if len(name) <= 253 {
		return name
	}
	return fmt.Sprintf("%s-%x", strings.TrimRight(name[:236], "-."), sha256.Sum256([]byte(name)))[:253]
}

func routeName(base, kind, hostname string) string {
	if hostname == "" {
		return gatewayAPIName(base, kind)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(hostname)))[:12]
	host := strings.NewReplacer("*", "wildcard", ".", "-").Replace(hostname)
	if len(host) > 48 {
		host = host[:48]
	}
	return gatewayAPIName(base, kind+"-"+strings.TrimRight(host, "-")+"-"+hash)
}

// mergeGatewayAPIMaps returns a new map with resource values overriding shared values.
func mergeGatewayAPIMaps(shared, resource map[string]string) map[string]string {
	merged := make(map[string]string, len(shared)+len(resource))
	maps.Copy(merged, shared)
	maps.Copy(merged, resource)
	return merged
}

// NewKubernetesGatewayResources builds the complete desired set without mutating the CR.
func NewKubernetesGatewayResources(gw *securityv1.Gateway) (*KubernetesGatewayResources, error) {
	if err := securityv1.ValidateKubernetesGateway(gw); err != nil {
		return nil, err
	}
	c := gw.Spec.App.KubernetesGateway
	if !c.Enabled {
		return &KubernetesGatewayResources{}, nil
	}
	sharedLabels := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	maps.Copy(sharedLabels, c.Labels)
	result := &KubernetesGatewayResources{}
	parentName, parentNS := c.Gateway.ExistingRef.Name, c.Gateway.ExistingRef.Namespace
	if parentNS == "" {
		parentNS = gw.Namespace
	}
	if c.Gateway.Create {
		parentName, parentNS = gatewayAPIName(gw.Name, "gateway"), gw.Namespace
	}
	group, kind, ns := gatewayv1.Group(gatewayv1.GroupName), gatewayv1.Kind("Gateway"), gatewayv1.Namespace(parentNS)
	parents := []gatewayv1.ParentReference{
		{
			Group:     &group,
			Kind:      &kind,
			Name:      gatewayv1.ObjectName(parentName),
			Namespace: &ns,
		},
	}
	httpGroups := map[string][]gatewayv1.HTTPRouteRule{}
	tlsGroups := map[string][]gatewayv1.TLSRouteRule{}
	management := false
	backend := func(selector string, port int32) gatewayv1.BackendRef {
		name := gw.Name
		if selector == "management" {
			name += "-management-service"
			management = true
		}
		p := gatewayv1.PortNumber(port)
		g, k, weight := gatewayv1.Group(""), gatewayv1.Kind("Service"), int32(1)
		return gatewayv1.BackendRef{
			BackendObjectReference: gatewayv1.BackendObjectReference{
				Group: &g,
				Kind:  &k,
				Name:  gatewayv1.ObjectName(name),
				Port:  &p,
			},
			Weight: &weight,
		}
	}
	if c.HTTPEnabled() {
		for _, rule := range c.HTTPRoute.Rules {
			host := string(rule.Hostname)
			r := gatewayv1.HTTPRouteRule{
				Matches: rule.Matches,
				Filters: rule.Filters,
				BackendRefs: []gatewayv1.HTTPBackendRef{
					{BackendRef: backend(rule.Backend, rule.Port)},
				},
			}
			if len(r.Matches) == 0 {
				typ, value := gatewayv1.PathMatchPathPrefix, "/"
				r.Matches = []gatewayv1.HTTPRouteMatch{
					{
						Path: &gatewayv1.HTTPPathMatch{
							Type:  &typ,
							Value: &value,
						},
					},
				}
			}
			httpGroups[host] = append(httpGroups[host], r)
		}
		for _, host := range slices.Sorted(maps.Keys(httpGroups)) {
			route := &gatewayv1.HTTPRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: gatewayv1.GroupVersion.String(),
					Kind:       "HTTPRoute",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        routeName(gw.Name, "httproute", host),
					Namespace:   gw.Namespace,
					Labels:      mergeGatewayAPIMaps(sharedLabels, c.HTTPRoute.Labels),
					Annotations: mergeGatewayAPIMaps(c.Annotations, c.HTTPRoute.Annotations),
				},
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: parents,
					},
					Rules: httpGroups[host],
				},
			}
			if host != "" {
				route.Spec.Hostnames = []gatewayv1.Hostname{gatewayv1.Hostname(host)}
			}
			result.HTTPRoutes = append(result.HTTPRoutes, route)
		}
	}
	if c.TLSRoute.Enabled {
		for _, rule := range c.TLSRoute.Rules {
			host := string(rule.Hostname)
			tlsGroups[host] = append(tlsGroups[host], gatewayv1.TLSRouteRule{
				BackendRefs: []gatewayv1.BackendRef{backend(rule.Backend, rule.Port)},
			})
		}
		for _, host := range slices.Sorted(maps.Keys(tlsGroups)) {
			result.TLSRoutes = append(result.TLSRoutes, &gatewayv1.TLSRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: gatewayv1.GroupVersion.String(),
					Kind:       "TLSRoute",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:        routeName(gw.Name, "tlsroute", host),
					Namespace:   gw.Namespace,
					Labels:      mergeGatewayAPIMaps(sharedLabels, c.TLSRoute.Labels),
					Annotations: mergeGatewayAPIMaps(c.Annotations, c.TLSRoute.Annotations),
				},
				Spec: gatewayv1.TLSRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: parents,
					},
					Hostnames: []gatewayv1.Hostname{gatewayv1.Hostname(host)},
					Rules:     tlsGroups[host],
				},
			})
		}
	}
	if c.Gateway.Create {
		hosts := slices.Sorted(maps.Keys(httpGroups))
		if c.TLSRoute.Enabled {
			hosts = slices.Sorted(maps.Keys(tlsGroups))
		}
		listeners := []gatewayv1.Listener{}
		for i, host := range hosts {
			from := gatewayv1.NamespacesFromSame
			listener := gatewayv1.Listener{
				Name:     gatewayv1.SectionName(fmt.Sprintf("https-%d", i)),
				Port:     443,
				Protocol: gatewayv1.HTTPSProtocolType,
				AllowedRoutes: &gatewayv1.AllowedRoutes{
					Namespaces: &gatewayv1.RouteNamespaces{
						From: &from,
					},
				},
			}
			if host != "" {
				hostname := gatewayv1.Hostname(host)
				listener.Hostname = &hostname
			}
			mode := gatewayv1.TLSModeTerminate
			listener.TLS = &gatewayv1.ListenerTLSConfig{
				Mode: &mode,
				CertificateRefs: []gatewayv1.SecretObjectReference{
					{
						Group: new(gatewayv1.Group("")),
						Kind:  new(gatewayv1.Kind("Secret")),
						Name:  gatewayv1.ObjectName(c.Gateway.TLS.ExistingSecretName),
					},
				},
			}
			if c.TLSRoute.Enabled {
				mode = gatewayv1.TLSModePassthrough
				listener.Name = gatewayv1.SectionName(fmt.Sprintf("tls-passthrough-%d", i))
				listener.Protocol = gatewayv1.TLSProtocolType
				listener.TLS.CertificateRefs = nil
			}
			listeners = append(listeners, listener)
		}
		if len(listeners) > 64 {
			return nil, fmt.Errorf("kubernetesGateway: a managed Gateway supports at most 64 distinct hostnames")
		}
		result.Gateways = append(result.Gateways, &gatewayv1.Gateway{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "Gateway",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        parentName,
				Namespace:   gw.Namespace,
				Labels:      mergeGatewayAPIMaps(sharedLabels, c.Gateway.Labels),
				Annotations: mergeGatewayAPIMaps(c.Annotations, c.Gateway.Annotations),
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: gatewayv1.ObjectName(c.Gateway.GatewayClassName),
				Addresses:        c.Gateway.Addresses,
				Listeners:        listeners,
			},
		})
	}
	if c.BackendTLSEnabled() {
		targets := []gatewayv1.LocalPolicyTargetReferenceWithSectionName{}
		addTargets := func(name string, ports []securityv1.Ports) {
			if len(c.BackendTLSPolicy.ExcludePorts) == 0 {
				targets = append(targets, gatewayv1.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gatewayv1.LocalPolicyTargetReference{
						Group: "",
						Kind:  "Service",
						Name:  gatewayv1.ObjectName(name),
					},
				})
				return
			}
			for _, port := range ports {
				if !slices.Contains(c.BackendTLSPolicy.ExcludePorts, port.Port) {
					section := gatewayv1.SectionName(port.Name)
					targets = append(targets, gatewayv1.LocalPolicyTargetReferenceWithSectionName{
						LocalPolicyTargetReference: gatewayv1.LocalPolicyTargetReference{
							Group: "",
							Kind:  "Service",
							Name:  gatewayv1.ObjectName(name),
						},
						SectionName: &section,
					})
				}
			}
		}
		addTargets(gw.Name, gw.Spec.App.Service.Ports)
		if management {
			addTargets(gw.Name+"-management-service", gw.Spec.App.Management.Service.Ports)
		}
		if len(targets) == 0 || len(targets) > 16 {
			return nil, fmt.Errorf("kubernetesGateway: backendTLSPolicy requires between 1 and 16 targets after excludePorts")
		}
		v := c.BackendTLSPolicy.Validation
		hostname := v.Hostname
		if hostname == "" {
			hostname = gatewayv1.PreciseHostname(fmt.Sprintf("%s.%s.svc.cluster.local", gw.Name, gw.Namespace))
		}
		result.BackendTLSPolicies = append(result.BackendTLSPolicies, &gatewayv1.BackendTLSPolicy{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "BackendTLSPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        gatewayAPIName(gw.Name, "backendtlspolicy"),
				Namespace:   gw.Namespace,
				Labels:      mergeGatewayAPIMaps(sharedLabels, c.BackendTLSPolicy.Labels),
				Annotations: mergeGatewayAPIMaps(c.Annotations, c.BackendTLSPolicy.Annotations),
			},
			Spec: gatewayv1.BackendTLSPolicySpec{
				TargetRefs: targets,
				Validation: gatewayv1.BackendTLSPolicyValidation{
					Hostname:                hostname,
					CACertificateRefs:       v.CACertificateRefs,
					WellKnownCACertificates: v.WellKnownCACertificates,
				},
			},
		})
	}
	return result, nil
}
