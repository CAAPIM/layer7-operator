package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

const WatchNamespaceEnvVar = "WATCH_NAMESPACE"
const OperatorNamespaceEnvVar = "OPERATOR_NAMESPACE"
const EnableWebHookEnvVar = "ENABLE_WEBHOOK"

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
