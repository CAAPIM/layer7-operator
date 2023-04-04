package monitoring

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServiceMonitor(gw *securityv1.Gateway) *monitoringv1.ServiceMonitor {
	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	//endpoints := []monitoringv1.Endpoint{}
	//endpoint :=
	serviceMonitorSpec := monitoringv1.ServiceMonitorSpec{
		JobLabel: gw.Name,
		NamespaceSelector: monitoringv1.NamespaceSelector{
			MatchNames: []string{gw.Namespace},
		},
		Selector: metav1.LabelSelector{
			MatchLabels: ls,
		},
		Endpoints: []monitoringv1.Endpoint{{Port: "monitoring", Path: "/metrics", Interval: monitoringv1.Duration("10s"), Scheme: "HTTP"}},
	}

	serviceMonitor := monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceMonitor",
		},
		Spec: serviceMonitorSpec,
	}

	return &serviceMonitor
}
