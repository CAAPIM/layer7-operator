package util

import (
	"fmt"
	"io/ioutil"
	"os"
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
	nsBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(nsBytes)), nil
}
