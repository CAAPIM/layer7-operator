package monitoring

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	otelv1alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewOtelCollector(gw *securityv1.Gateway) *otelv1alpha1.OpenTelemetryCollector {

	ports := []corev1.ServicePort{}
	ports = append(ports, v1.ServicePort{Port: 8889, Name: "prometheus"})
	config := `
receivers:
  otlp:
    protocols:
      grpc:
      http:
processors:
  batch:
exporters:
  logging:
    loglevel: warn 
  prometheus:
    endpoint: "0.0.0.0:8889"
    const_labels:
      name: ssg
service:
  telemetry:
    logs:
      level: "debug"
    metrics:
      address: "0.0.0.0:8888"
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [logging]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus, logging]
    logs: 
      receivers: [otlp]
      exporters: [logging]
extensions:
  health_check:
  pprof:
    endpoint: 0.0.0.0:1777
  zpages:
    endpoint: 0.0.0.0:55679
`

	otelSpec := otelv1alpha1.OpenTelemetryCollectorSpec{
		Mode:   otelv1alpha1.ModeSidecar,
		Ports:  ports,
		Config: config,
	}

	otelCollector := otelv1alpha1.OpenTelemetryCollector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1alpha1",
			Kind:       "OpenTelemetryCollector",
		},
		Spec: otelSpec,
	}

	return &otelCollector
}
