package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewRoute(gw *securityv1.Gateway) routev1.Route {

	labels := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)

	defaultPort := routev1.RoutePort{TargetPort: intstr.FromString("https")}
	defaultTls := routev1.TLSConfig{
		Termination:                   routev1.TLSTerminationPassthrough,
		InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
	}
	defaultWildcardPolicy := routev1.WildcardPolicyNone

	routeSpec := routev1.RouteSpec{
		To: routev1.RouteTargetReference{
			Kind: "Service",
			Name: gw.Name,
		},
		Port: &defaultPort,
		// default to passthrough
		// will be updated to reencrypt in the future
		// can be overriden
		TLS:            &defaultTls,
		WildcardPolicy: defaultWildcardPolicy,
	}

	routeOverrides := gw.Spec.App.Ingress.Route

	if routeOverrides.Host != "" {
		routeSpec.Host = routeOverrides.Host
	}

	if routeOverrides.Path != "" {
		routeSpec.Path = routeOverrides.Path
	}

	if routeOverrides.To != (routev1.RouteTargetReference{}) {
		routeSpec.To = routeOverrides.To
	}

	if routeOverrides.AlternateBackends != nil {
		routeSpec.AlternateBackends = routeOverrides.AlternateBackends
	}

	if routeOverrides.Port != nil {
		routeSpec.Port = routeOverrides.Port
	}

	if routeOverrides.TLS != nil {
		routeSpec.TLS = routeOverrides.TLS
	}

	if routeOverrides.WildcardPolicy != "" {
		routeSpec.WildcardPolicy = routeOverrides.WildcardPolicy
	}

	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name,
			Namespace:   gw.Namespace,
			Labels:      labels,
			Annotations: gw.Spec.App.Ingress.Annotations,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "route.openshift.io/v1",
			Kind:       "Route",
		},
		Spec: routeSpec,
	}

	return *route

}
