package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"
	"testing"
)

func TestDefaultJVMHeapSize(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Java = securityv1.Java{}
	gateway.Spec.App.Java.JVMHeap = securityv1.JVMHeap{}
	gateway.Spec.App.Java.JVMHeap.Calculate = false

	resultHeap := setJVMHeapSize(&gateway)

	if resultHeap != "2g" {
		t.Errorf("Default heap should be 2g")
	}
}

func TestConfiguredDefaultJVMHeapSize(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Java = securityv1.Java{}
	gateway.Spec.App.Java.JVMHeap = securityv1.JVMHeap{}
	gateway.Spec.App.Java.JVMHeap.Calculate = true
	gateway.Spec.App.Java.JVMHeap.Default = "3g"

	gateway.Spec.App.Resources = securityv1.PodResources{}
	gateway.Spec.App.Resources.Limits = corev1.ResourceList{}

	resultHeap := setJVMHeapSize(&gateway)

	if resultHeap != "3g" {
		t.Errorf("Expected 3g, but Actual %s", resultHeap)
	}
}

func TestGivenJVMHeapSize(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Java = securityv1.Java{}
	gateway.Spec.App.Java.JVMHeap = securityv1.JVMHeap{}
	gateway.Spec.App.Java.JVMHeap.Calculate = true
	gateway.Spec.App.Java.JVMHeap.Default = "3g"
	gateway.Spec.App.Java.JVMHeap.Percentage = 50

	gateway.Spec.App.Resources = securityv1.PodResources{}
	gateway.Spec.App.Resources.Limits = corev1.ResourceList{"memory": resource.MustParse("4Gi")}

	resultHeap := setJVMHeapSize(&gateway)

	if resultHeap != "2048m" {
		t.Errorf("Expected 2048m, but Actual %s", resultHeap)
	}
}

func TestSystemPropertiesConfigMap(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Java = securityv1.Java{}
	gateway.Spec.App.Java.JVMHeap = securityv1.JVMHeap{}
	gateway.Spec.App.Java.JVMHeap.Calculate = false
	gateway.Spec.App.Java.ExtraArgs = []string{}

	expected := "testProperty"
	gateway.Spec.App.System = securityv1.System{expected}
	gateway.Name = "test"
	configMap := NewConfigMap(&gateway, "test-system")
	propValue := configMap.Data["system.properties"]
	if propValue != expected {
		t.Errorf("Expected %s, but Actual %s", expected, propValue)
	}
}

func TestGatewayPropertiesConfigMap(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Java = securityv1.Java{}
	gateway.Spec.App.Java.JVMHeap = securityv1.JVMHeap{}
	gateway.Spec.App.Java.JVMHeap.Calculate = false
	gateway.Spec.App.Java.ExtraArgs = []string{}

	gateway.Spec.License = securityv1.License{}
	gateway.Spec.License.Accept = true

	gateway.Spec.App.Management = securityv1.Management{}
	gateway.Spec.App.Management.Cluster = securityv1.Cluster{}
	gateway.Spec.App.Management.Cluster.Hostname = "testHost"
	gateway.Spec.App.Management.Database = securityv1.Database{}
	gateway.Spec.App.Management.Database.Enabled = true
	gateway.Spec.App.Management.Database.JDBCUrl = "jdbc:mysql:localhost:3606"

	gateway.Spec.App.Hazelcast = securityv1.Hazelcast{}
	gateway.Spec.App.Hazelcast.External = true
	gateway.Spec.App.Hazelcast.Endpoint = "hazelcasthost"

	gateway.Name = "test"
	configMap := NewConfigMap(&gateway, "test")
	acceptLicense := configMap.Data["ACCEPT_LICENSE"]
	if acceptLicense != "true" {
		t.Errorf("Expected %s, but Actual %s", "true", acceptLicense)
	}
	ssgClusterHost := configMap.Data["SSG_CLUSTER_HOST"]
	if ssgClusterHost != "testHost" {
		t.Errorf("Expected %s, but Actual %s", "testHost", ssgClusterHost)
	}
	jvmHeap := configMap.Data["SSG_JVM_HEAP"]
	if jvmHeap != "2g" {
		t.Errorf("Expected %s, but Actual %s", "2g", jvmHeap)
	}
	jvmExtraArgs := configMap.Data["EXTRA_JAVA_ARGS"]
	if !strings.Contains(jvmExtraArgs, "-Dcom.l7tech.server.extension.sharedCounterProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedKeyValueStoreProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedClusterInfoProvider=externalhazelcast") {
		t.Errorf("jvmExtraArgs %s should containe hazelcast property", jvmExtraArgs)
	}
	expectHazelcaseClientConfig := "<hazelcast-client xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.hazelcast.com/schema/client-config http://www.hazelcast.com/schema/client-config/hazelcast-client-config-3.10.xsd\" xmlns=\"http://www.hazelcast.com/schema/client-config\"><instance-name>test-hazelcast-client</instance-name><network><cluster-members><address>hazelcasthost</address></cluster-members><connection-attempt-limit>10</connection-attempt-limit><redo-operation>true</redo-operation></network><connection-strategy async-start=\"false\" reconnect-mode=\"ON\" /></hazelcast-client>"
	hazelcaseClientConfig := configMap.Data["hazelcast-client.xml"]
	if hazelcaseClientConfig != expectHazelcaseClientConfig {
		t.Errorf("Expected %s, but Actual %s", expectHazelcaseClientConfig, hazelcaseClientConfig)
	}
}

func TestCWPPropertiesConfigMap(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Spec.App.ClusterProperties = securityv1.ClusterProperties{}
	gateway.Spec.App.ClusterProperties.Properties = []securityv1.Property{securityv1.Property{"cwp1", "test1"}}
	configMap := NewConfigMap(&gateway, "test-cwp-bundle")
	cwpBundle := configMap.Data["cwp.bundle"]
	if !strings.Contains(cwpBundle, "cwp1") {
		t.Errorf("cwpBundle %s should contain property %s", cwpBundle, "cwp1")
	}
}

func TestDefaultListenPortBundle(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Spec.App.ListenPorts = securityv1.ListenPorts{}
	gateway.Spec.App.ListenPorts.Custom = securityv1.CustomListenPort{false}
	configMap := NewConfigMap(&gateway, "test-listen-port-bundle")
	listenPortBundle := configMap.Data["listen-ports.bundle"]
	if !strings.Contains(listenPortBundle, "8080") {
		t.Errorf("listenPortBundle %s should contain port %s", listenPortBundle, "8080")
	}
	if !strings.Contains(listenPortBundle, "8443") {
		t.Errorf("listenPortBundle %s should contain port %s", listenPortBundle, "8443")
	}
	if !strings.Contains(listenPortBundle, "9443") {
		t.Errorf("listenPortBundle %s should contain port %s", listenPortBundle, "9443")
	}
}

func TestCustomListenPortBundle(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Spec.App.ListenPorts = securityv1.ListenPorts{}
	gateway.Spec.App.ListenPorts.Custom = securityv1.CustomListenPort{true}

	port := securityv1.ListenPort{}
	port.Name = "custom1"
	port.Enabled = true
	port.Protocol = "http"
	port.Port = "9090"
	gateway.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{port}

	configMap := NewConfigMap(&gateway, "test-listen-port-bundle")
	listenPortBundle := configMap.Data["listen-ports.bundle"]
	if !strings.Contains(listenPortBundle, "9090") {
		t.Errorf("listenPortBundle %s should contain port %s", listenPortBundle, "9090")
	}

}

func TestRepositoryConfigWithAuth(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Status = securityv1.GatewayStatus{}
	gatewayRepositoryStatus := securityv1.GatewayRepositoryStatus{}
	gatewayRepositoryStatus.Enabled = true
	gatewayRepositoryStatus.Name = "testrepo"
	gatewayRepositoryStatus.Commit = "1234"
	gatewayRepositoryStatus.Type = "static"
	gatewayRepositoryStatus.SecretName = "testSecret"
	gatewayRepositoryStatus.Branch = "testBranch"
	gatewayRepositoryStatus.Endpoint = "github.com"

	gateway.Status.RepositoryStatus = []securityv1.GatewayRepositoryStatus{gatewayRepositoryStatus}

	configMap := NewConfigMap(&gateway, "-repository-init-config")
	repositoryConfig := configMap.Data["config.json"]
	if !strings.Contains(repositoryConfig, "/graphman/secrets/testrepo") {
		t.Errorf("repositoryConfig %s should contain auth %s", repositoryConfig, "/graphman/secrets/testrepo")
	}
}

func TestRepositoryConfigWithLocalRef(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Status = securityv1.GatewayStatus{}
	gatewayRepositoryStatus := securityv1.GatewayRepositoryStatus{}
	gatewayRepositoryStatus.Enabled = true
	gatewayRepositoryStatus.Name = "testrepo"
	gatewayRepositoryStatus.Type = "static"
	gatewayRepositoryStatus.StorageSecretName = "testStorageSecret"

	gateway.Status.RepositoryStatus = []securityv1.GatewayRepositoryStatus{gatewayRepositoryStatus}

	configMap := NewConfigMap(&gateway, "-repository-init-config")
	repositoryConfig := configMap.Data["config.json"]
	if !strings.Contains(repositoryConfig, "/graphman/localref/testStorageSecret/testrepo.gz") {
		t.Errorf("repositoryConfig %s should contain auth %s", repositoryConfig, "/graphman/secrets/testrepo")
	}
}
