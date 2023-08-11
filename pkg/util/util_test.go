package util

import (
	"os"
	"testing"
)

func TestGetWatchNamespace(t *testing.T) {
	os.Setenv(WatchNamespaceEnvVar, "test")
	actual, err := GetWatchNamespace()
	expected := "test"
	if err != nil {
		t.Errorf("%s must be set", WatchNamespaceEnvVar)
	}
	if actual != expected {
		t.Errorf("actual %s, expected %s", actual, expected)
	}
}

func TestContains(t *testing.T) {
	array := []string{"test1", "test2", "test3"}
	found := Contains(array, "test1")
	if !found {
		t.Errorf("%s must be present", "test1")
	}
}

func TestDefaultLabels(t *testing.T) {
	additionalLabels := make(map[string]string)
	additionalLabelValue := "label1"
	additionalLabels["additionalLabel1"] = additionalLabelValue
	testName := "test"
	labels := DefaultLabels(testName, additionalLabels)
	value, ok := labels["app.kubernetes.io/name"]
	if !ok {
		t.Errorf("key is missing")
	}
	if value != testName {
		t.Errorf("value is wrong for given key %s", "app.kubernetes.io/name")
	}

	value, ok = labels["additionalLabel1"]
	if !ok {
		t.Errorf("key is missing")
	}
	if value != additionalLabelValue {
		t.Errorf("value is wrong for given key %s", "additionalLabel1")
	}

}
