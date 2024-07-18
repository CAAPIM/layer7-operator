package util

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DefaultLabels(name string, additionalLabels map[string]string) map[string]string {

	labels := map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/managed-by": "layer7-operator",
		"app.kubernetes.io/created-by": "layer7-operator",
		"app.kubernetes.io/part-of":    name,
	}

	for k, v := range additionalLabels {
		labels[k] = v
	}

	return labels
}

// Contains returns true if string array contains string
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

var (
	WatchNamespaceEnvVar    = "WATCH_NAMESPACE"
	OperatorNamespaceEnvVar = "OPERATOR_NAMESPACE"
	EnableWebHookEnvVar     = "ENABLE_WEBHOOK"
	EnableOtelEnvVar        = "ENABLE_OTEL"
	OtelCollectorUrlEnvVar  = "OTEL_EXPORTER_OTLP_ENDPOINT"
	OtelMetricPrefixEnvVar  = "OTEL_METRIC_PREFIX"
	HostNameEnvVar          = "HOSTNAME"
)

// GetWatchNamespace returns the namespace the operator should be watching for changes
func GetWatchNamespace() (string, error) {
	ns, found := os.LookupEnv(WatchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", WatchNamespaceEnvVar)
	}
	return ns, nil
}

// GetOperatorNamespace returns the namespace of the operator pod
func GetOperatorNamespace() (string, error) {
	nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if os.IsNotExist(err) {
		ns, found := os.LookupEnv(OperatorNamespaceEnvVar)
		if !found {
			return "", fmt.Errorf("%s must be set to run locally", OperatorNamespaceEnvVar)
		}
		return ns, nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(nsBytes)), nil
}

func GetWebhookEnabled() (bool, error) {
	wh, found := os.LookupEnv(EnableWebHookEnvVar)
	if !found {
		return false, nil
	}
	enabled, err := strconv.ParseBool(wh)

	if err != nil {
		return false, err
	}

	return enabled, nil
}

func GetOtelEnabled() (bool, error) {
	o, found := os.LookupEnv(EnableOtelEnvVar)
	if !found {
		return false, nil
	}
	enabled, err := strconv.ParseBool(o)

	if err != nil {
		return false, err
	}

	return enabled, nil
}

func GetOtelCollectorUrl() (string, error) {
	collectorUrl, found := os.LookupEnv(OtelCollectorUrlEnvVar)
	if !found {
		return "", nil
	}

	return collectorUrl, nil
}

func GetHostname() (string, error) {
	hostname, found := os.LookupEnv(HostNameEnvVar)
	if !found {
		return "", nil
	}

	return hostname, nil
}

func GetOtelMetricPrefix() (string, error) {
	otelMetricPrefix, found := os.LookupEnv(OtelMetricPrefixEnvVar)
	if !found {
		return "", nil
	}

	return otelMetricPrefix, nil
}

func InitOTelProvider(collectorURL string, ctx context.Context) (func(context.Context) error, error) {

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hostname: %w", err)
	}

	ns, err := GetOperatorNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to operator namespace: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attribute.String("host.name", hostname), attribute.String("namespace", ns)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// TODO: Expand to support TLS in the future.
	conn, err := grpc.NewClient(collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a meter exporter
	meterExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn), otlpmetricgrpc.WithCompressor("gzip"))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(metric.NewPeriodicReader(meterExporter)))

	otel.SetMeterProvider(meterProvider)

	// Shutdown will flush any remaining spans and shut down the exporter.
	return meterProvider.Shutdown, nil
}
