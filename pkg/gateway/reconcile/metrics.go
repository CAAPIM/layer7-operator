package reconcile

import (
	"context"
	"strings"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func captureGraphmanMetrics(ctx context.Context, params Params, start time.Time, podName string, bundleType string, bundleName string, sha1sum string, hasError bool) error {
	operatorNamespace, err := util.GetOperatorNamespace()
	if err != nil {
		params.Log.Info("could not determine operator namespace")
		return err
	}
	gateway := params.Instance
	otelEnabled, err := util.GetOtelEnabled()
	if err != nil {
		params.Log.Info("could not determine if OTel is enabled")
		return err
	}

	if !otelEnabled {
		return nil
	}

	otelMetricPrefix, err := util.GetOtelMetricPrefix()
	if err != nil {
		params.Log.Info("could not determine otel metric prefix")
		return err
	}

	if otelMetricPrefix == "" {
		otelMetricPrefix = "layer7_"
	}

	hostname, err := util.GetHostname()
	if err != nil {
		params.Log.Error(err, "failed to retrieve operator hostname")
		return err
	}

	meter := otel.Meter("layer7-operator-graphman-metrics")

	graphmanApplyLatency, err := meter.Float64Histogram(otelMetricPrefix+"operator_graphman_latency",
		metric.WithDescription("gateway controller reconcile latency"), metric.WithUnit("ms"))
	if err != nil {
		return err
	}

	graphmanRequestSuccess, err := meter.Int64Counter(otelMetricPrefix+"operator_graphman_request_success",
		metric.WithDescription("graphman request success"))
	if err != nil {
		return err
	}

	graphmanRequestFailure, err := meter.Int64Counter(otelMetricPrefix+"operator_graphman_request_failure",
		metric.WithDescription("graphman request failure"))
	if err != nil {
		return err
	}

	graphmanRequestTotal, err := meter.Int64Counter(otelMetricPrefix+"operator_graphman_request_total",
		metric.WithDescription("graphman request total"))
	if err != nil {
		return err
	}

	duration := time.Since(start)
	graphmanApplyLatency.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_pod", podName),
			attribute.String("bundle_type", bundleType),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	graphmanRequestTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("k8s.pod.name", hostname),
			attribute.String("k8s.namespace.name", operatorNamespace),
			attribute.String("gateway_namespace", gateway.Namespace),
			attribute.String("gateway_name", gateway.Name),
			attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1])))

	if hasError {
		graphmanRequestFailure.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("sha1sum", sha1sum),
				attribute.String("bundle_type", bundleType),
				attribute.String("bundle_name", bundleName),
				attribute.String("gateway_pod", podName),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1]),
				attribute.String("applied_time", time.Now().UTC().String())))
	} else {
		graphmanRequestSuccess.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("k8s.pod.name", hostname),
				attribute.String("k8s.namespace.name", operatorNamespace),
				attribute.String("gateway_namespace", gateway.Namespace),
				attribute.String("sha1sum", sha1sum),
				attribute.String("bundle_type", bundleType),
				attribute.String("bundle_name", bundleName),
				attribute.String("gateway_pod", podName),
				attribute.String("gateway_name", gateway.Name),
				attribute.String("gateway_version", strings.Split(gateway.Spec.App.Image, ":")[1]),
				attribute.String("applied_time", time.Now().UTC().String())))
	}

	return nil
}
