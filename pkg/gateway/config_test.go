package gateway

import (
	"strconv"
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaultJVMHeapSize(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				Java: securityv1.Java{
					JVMHeap: securityv1.JVMHeap{
						Calculate: false,
					},
				},
			},
		}}

	legacyResultHeap := setJVMHeapSize(&gateway, "", gateway.Spec.App.Java.JVMHeap.Percentage)
	resultMinHeap := setJVMHeapSize(&gateway, "min", gateway.Spec.App.Java.JVMHeap.MinPercentage)
	resultMaxHeap := setJVMHeapSize(&gateway, "max", gateway.Spec.App.Java.JVMHeap.MaxPercentage)

	if legacyResultHeap != "3g" {
		t.Errorf("Default heap should be 3g")
	}
	if resultMinHeap != "1g" {
		t.Errorf("Default min heap should be 3g")
	}
	if resultMaxHeap != "3g" {
		t.Errorf("Default max heap should be 3g")
	}
}

func TestConfiguredDefaultJVMHeapSize(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				Java: securityv1.Java{
					JVMHeap: securityv1.JVMHeap{
						Calculate:  false,
						Default:    "5g",
						MinDefault: "2g",
						MaxDefault: "4g",
					},
				},
				Resources: securityv1.PodResources{
					Limits:   corev1.ResourceList{},
					Requests: corev1.ResourceList{},
				},
			},
		}}

	legacyResultHeap := setJVMHeapSize(&gateway, "", gateway.Spec.App.Java.JVMHeap.Percentage)
	resultMinHeap := setJVMHeapSize(&gateway, "min", gateway.Spec.App.Java.JVMHeap.MinPercentage)
	resultMaxHeap := setJVMHeapSize(&gateway, "max", gateway.Spec.App.Java.JVMHeap.MaxPercentage)

	if legacyResultHeap != "5g" {
		t.Errorf("Default heap should be 5g")
	}
	if resultMinHeap != "2g" {
		t.Errorf("Default min heap should be 2g")
	}
	if resultMaxHeap != "4g" {
		t.Errorf("Default max heap should be 4g")
	}
}

func TestGivenJVMHeapSize(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				Java: securityv1.Java{
					JVMHeap: securityv1.JVMHeap{
						Calculate:     true,
						MinPercentage: 25,
						MaxPercentage: 75,
						Percentage:    75,
					},
				},
				Resources: securityv1.PodResources{
					Limits: corev1.ResourceList{
						"memory": resource.MustParse("4Gi"),
					},
					Requests: corev1.ResourceList{
						"memory": resource.MustParse("4Gi"),
					},
				},
			},
		},
	}

	resultHeap := setJVMHeapSize(&gateway, "", gateway.Spec.App.Java.JVMHeap.Percentage)

	legacyResultHeap := setJVMHeapSize(&gateway, "", gateway.Spec.App.Java.JVMHeap.Percentage)
	resultMinHeap := setJVMHeapSize(&gateway, "min", gateway.Spec.App.Java.JVMHeap.MinPercentage)
	resultMaxHeap := setJVMHeapSize(&gateway, "max", gateway.Spec.App.Java.JVMHeap.MaxPercentage)

	if legacyResultHeap != "3072m" {
		t.Errorf("Expected 3072m, but Actual %s", resultHeap)
	}
	if resultMinHeap != "1024m" {
		t.Errorf("Expected 1024m, but Actual %s", resultHeap)
	}
	if resultMaxHeap != "3072m" {
		t.Errorf("Expected 3072m, but Actual %s", resultHeap)
	}
}

func TestSystemPropertiesConfigMap(t *testing.T) {
	expectedSystemProp := "testProperty"
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				Java: securityv1.Java{
					JVMHeap: securityv1.JVMHeap{
						Calculate: false,
					},
					ExtraArgs: []string{},
				},
				System: securityv1.System{
					Properties: expectedSystemProp,
				},
			},
		}}

	configMap := NewConfigMap(&gateway, gateway.Name+"-system")

	if configMap.Data["system.properties"] != expectedSystemProp {
		t.Errorf("Expected %s, but Actual %s", expectedSystemProp, configMap.Data["system.properties"])
	}
}

func TestGatewayPropertiesConfigMap(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			License: securityv1.License{
				Accept: true,
			},
			App: securityv1.App{
				Java: securityv1.Java{
					JVMHeap: securityv1.JVMHeap{
						Calculate: false,
					},
					ExtraArgs: []string{},
				},
				Management: securityv1.Management{
					Cluster: securityv1.Cluster{
						Hostname: "testHost",
						Password: "7layer",
					},
					Database: securityv1.Database{
						Enabled: true,
						JDBCUrl: "jdbc:mysql:localhost:3306",
					},
				},
				Hazelcast: securityv1.Hazelcast{
					External: true,
					Endpoint: "hazelcasthost:5701",
				},
			},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name)

	expectedConfigData := map[string]string{
		"ACCEPT_LICENSE":        strconv.FormatBool(gateway.Spec.License.Accept),
		"SSG_CLUSTER_HOST":      gateway.Spec.App.Management.Cluster.Hostname,
		"SSG_CLUSTER_PASS":      gateway.Spec.App.Management.Cluster.Password,
		"SSG_JVM_HEAP":          "3g",
		"EXTRA_JAVA_ARGS":       " -Dcom.l7tech.server.extension.sharedCounterProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedKeyValueStoreProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedClusterInfoProvider=externalhazelcast",
		"LIQUIBASE_LOG_LEVEL":   "off",
		"DISKLESS_CONFIG":       "true",
		"SSG_DATABASE_JDBC_URL": "jdbc:mysql:localhost:3306",
		"hazelcast-client.xml":  `<hazelcast-client xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.hazelcast.com/schema/client-config http://www.hazelcast.com/schema/client-config/hazelcast-client-config-5.2.xsd" xmlns="http://www.hazelcast.com/schema/client-config"><instance-name>test-hazelcast-client</instance-name><network><cluster-members><address>hazelcasthost:5701</address></cluster-members><redo-operation>true</redo-operation></network><connection-strategy async-start="false" reconnect-mode="ON"><connection-retry><cluster-connect-timeout-millis>-1</cluster-connect-timeout-millis></connection-retry></connection-strategy></hazelcast-client>`,
	}

	for i := range configMap.Data {
		if configMap.Data[i] != expectedConfigData[i] {
			t.Errorf("expected %s, actual %s", expectedConfigData[i], configMap.Data[i])

		}
	}
}

func TestCWPPropertiesConfigMap(t *testing.T) {
	expectedCwpName := "cwptestprop1"
	expectedCwpValue := "cwptestval1"
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				ClusterProperties: securityv1.ClusterProperties{
					Enabled: true,
					Properties: []securityv1.Property{{
						Name:  expectedCwpName,
						Value: expectedCwpValue,
					}},
				},
			},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name+"-cwp-bundle")
	cwpBundle := configMap.Data["cwp.json"]

	if !strings.Contains(configMap.Data["cwp.json"], expectedCwpName) {
		t.Errorf("cwpBundle %s should contain property %s", cwpBundle, expectedCwpName)
	}

	if !strings.Contains(configMap.Data["cwp.json"], expectedCwpValue) {
		t.Errorf("cwpBundle %s should contain property %s", cwpBundle, expectedCwpValue)
	}
}

func TestDefaultListenPortBundle(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				ListenPorts: securityv1.ListenPorts{
					Custom: securityv1.CustomListenPort{
						Enabled: false,
					},
				},
			},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name+"-listen-port-bundle")
	listenPortBundle := configMap.Data["listen-ports.json"]
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
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: securityv1.GatewaySpec{
			App: securityv1.App{
				ListenPorts: securityv1.ListenPorts{
					Custom: securityv1.CustomListenPort{
						Enabled: true,
					},
					Ports: []securityv1.ListenPort{{
						Name:     "custom1",
						Port:     9090,
						Protocol: "http",
						Enabled:  true,
					}},
				},
			},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name+"-listen-port-bundle")
	listenPortBundle := configMap.Data["listen-ports.json"]
	if !strings.Contains(listenPortBundle, "9090") {
		t.Errorf("listenPortBundle %s should contain port %s", listenPortBundle, "9090")
	}

}

func TestRepositoryConfigWithAuth(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Status: securityv1.GatewayStatus{
			RepositoryStatus: []securityv1.GatewayRepositoryStatus{{
				Enabled:    true,
				Name:       "testrepo",
				Commit:     "1234",
				Type:       "static",
				SecretName: "testSecret",
				Branch:     "testBranch",
				Endpoint:   "github.com",
			}},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name+"-repository-init-config")
	repositoryConfig := configMap.Data["config.json"]
	if !strings.Contains(repositoryConfig, "/graphman/secrets/testrepo") {
		t.Errorf("repositoryConfig %s should contain auth %s", repositoryConfig, "/graphman/secrets/testrepo")
	}
}

func TestRepositoryConfigWithLocalRef(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Status: securityv1.GatewayStatus{
			RepositoryStatus: []securityv1.GatewayRepositoryStatus{{
				Enabled:           true,
				Name:              "testrepo",
				Type:              "static",
				StorageSecretName: "testStorageSecret",
			}},
		},
	}

	configMap := NewConfigMap(&gateway, gateway.Name+"-repository-init-config")
	repositoryConfig := configMap.Data["config.json"]
	if !strings.Contains(repositoryConfig, "/graphman/localref/testStorageSecret/testrepo.gz") {
		t.Errorf("repositoryConfig %s should contain auth %s", repositoryConfig, "/graphman/secrets/testrepo")
	}
}
