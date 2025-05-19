package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewRoute(gw *securityv1.Gateway, routeSpec securityv1.RouteSpec, suffix string, managementRoute bool) *routev1.Route {

	labels := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)

	defaultPort := &routev1.RoutePort{TargetPort: intstr.FromString("https")}
	defaultWeight := int32(100)
	serviceName := gw.Name

	if managementRoute {
		defaultPort = &routev1.RoutePort{TargetPort: intstr.FromString("management")}
		serviceName = gw.Name + "-management-service"
	}

	if routeSpec.Port != nil {
		if routeSpec.Port.TargetPort.StrVal != "" || routeSpec.Port.TargetPort.IntVal != 0 {
			defaultPort = routeSpec.Port
		}
	}

	defaultTls := routev1.TLSConfig{
		Termination:                   routev1.TLSTerminationPassthrough,
		InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
	}

	if routeSpec.TLS != nil {
		defaultTls = *routeSpec.TLS
	}

	defaultWildcardPolicy := routev1.WildcardPolicyNone

	defaultRouteSpec := routev1.RouteSpec{
		To: routev1.RouteTargetReference{
			Kind:   "Service",
			Name:   serviceName,
			Weight: &defaultWeight,
		},
		Port:           defaultPort,
		TLS:            &defaultTls,
		WildcardPolicy: defaultWildcardPolicy,
	}

	routeOverrides := routeSpec

	if routeOverrides.To != nil {
		if routeOverrides.To.Name != "" {
			defaultRouteSpec.To.Name = routeSpec.To.Name
		}
		if routeOverrides.To.Weight != nil {
			defaultRouteSpec.To.Weight = routeSpec.To.Weight
		}
	}

	if routeOverrides.Host != "" {
		defaultRouteSpec.Host = routeOverrides.Host
	}

	if routeOverrides.Path != "" {
		defaultRouteSpec.Path = routeOverrides.Path
	}

	if routeOverrides.Port != nil {
		defaultRouteSpec.Port = routeOverrides.Port
	}

	if routeOverrides.TLS != nil {
		defaultRouteSpec.TLS = routeOverrides.TLS
	}

	if routeOverrides.WildcardPolicy != "" {
		defaultRouteSpec.WildcardPolicy = routeOverrides.WildcardPolicy
	}

	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name + "-" + suffix,
			Namespace:   gw.Namespace,
			Labels:      labels,
			Annotations: gw.Spec.App.Ingress.Annotations,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "route.openshift.io/v1",
			Kind:       "Route",
		},
		Spec: defaultRouteSpec,
	}

	return route

}
