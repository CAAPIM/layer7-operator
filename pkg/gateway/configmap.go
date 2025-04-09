package gateway

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InitContainerStaticConfig struct {
	Version      string             `json:"version"`
	Repositories []RepositoryConfig `json:"repositories,omitempty"`
}

type RepositoryConfig struct {
	Name                string `json:"name"`
	Endpoint            string `json:"endpoint"`
	Branch              string `json:"branch"`
	Auth                string `json:"auth"`
	LocalReference      string `json:"localReference,omitempty"`
	Tag                 string `json:"tag,omitempty"`
	RemoteName          string `json:"remoteName,omitempty"`
	StateStoreReference string `json:"stateStoreReference,omitempty"`
	StateStoreKey       string `json:"stateStoreKey,omitempty"`
	SingletonExtraction bool   `json:"singletonExtraction,omitempty"`
}

// NewConfigMap
func NewConfigMap(gw *securityv1.Gateway, name string) *corev1.ConfigMap {
	javaArgs := strings.Join(gw.Spec.App.Java.ExtraArgs, " ")
	data := make(map[string]string)
	jvmHeap := setJVMHeapSize(gw, "", gw.Spec.App.Java.JVMHeap.Percentage)
	dataCheckSum := ""
	disklessConfig := "true"
	if gw.Spec.App.Management.DisklessConfig.Disabled {
		disklessConfig = "false"
	}
	liquibaseLogLevel := "off"
	if gw.Spec.App.Management.Database.LiquibaseLogLevel != "" {
		liquibaseLogLevel = string(gw.Spec.App.Management.Database.LiquibaseLogLevel)
	}
	switch name {
	case gw.Name + "-system":
		data["system.properties"] = gw.Spec.App.System.Properties
	case gw.Name:
		data["ACCEPT_LICENSE"] = strconv.FormatBool(gw.Spec.License.Accept)
		data["SSG_CLUSTER_HOST"] = gw.Spec.App.Management.Cluster.Hostname
		data["SSG_JVM_HEAP"] = jvmHeap
		data["LIQUIBASE_LOG_LEVEL"] = liquibaseLogLevel
		data["EXTRA_JAVA_ARGS"] = javaArgs
		data["DISKLESS_CONFIG"] = disklessConfig
		if (gw.Spec.App.Java.JVMHeap.MinPercentage != 0 && gw.Spec.App.Java.JVMHeap.MaxPercentage != 0) || (gw.Spec.App.Java.JVMHeap.MinDefault != "" && gw.Spec.App.Java.JVMHeap.MaxDefault != "") {
			minJvmHeap := setJVMHeapSize(gw, "min", gw.Spec.App.Java.JVMHeap.MinPercentage)
			maxJvmHeap := setJVMHeapSize(gw, "max", gw.Spec.App.Java.JVMHeap.MaxPercentage)
			data["SSG_JVM_MIN_HEAP"] = minJvmHeap
			data["SSG_JVM_MAX_HEAP"] = maxJvmHeap
		}
		if gw.Spec.App.Management.Database.Enabled && !gw.Spec.App.Management.DisklessConfig.Disabled {
			data["SSG_DATABASE_JDBC_URL"] = gw.Spec.App.Management.Database.JDBCUrl
		}
		if gw.Spec.App.Hazelcast.External {
			data["EXTRA_JAVA_ARGS"] = javaArgs + " -Dcom.l7tech.server.extension.sharedCounterProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedKeyValueStoreProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedClusterInfoProvider=externalhazelcast"
		}

		if gw.Spec.App.Otk.Enabled && gw.Spec.App.Otk.HealthCheck.Enabled {
			otkHealthCheckPort := 8443
			if gw.Spec.App.Otk.HealthCheck.Port != 0 {
				otkHealthCheckPort = gw.Spec.App.Otk.HealthCheck.Port
			}
			data["OTK_HEALTHCHECK_PORT"] = strconv.Itoa(otkHealthCheckPort)
		}

	case gw.Name + "-gateway-files":
		if gw.Spec.App.Bootstrap.Script.Enabled {
			f, _ := os.ReadFile("./003-parse-custom-files.sh")
			data["003-parse-custom-files"] = string(f)
		}
		if gw.Spec.App.PreStopScript.Enabled {
			f, _ := os.ReadFile("./graceful-shutdown.sh")
			data["graceful-shutdown"] = string(f)
		}

		if gw.Spec.App.Otk.Enabled && gw.Spec.App.Otk.HealthCheck.Enabled {
			f, _ := os.ReadFile("./otk-healthcheck.sh")
			data["otk-healthcheck"] = string(f)
		}

		if gw.Spec.App.AutoMountServiceAccountToken {
			f, _ := os.ReadFile("./load-service-account-token.sh")
			data["load-service-account-token"] = string(f)
			f, _ = os.ReadFile("./update-service-account-token.xml")
			data["service-account-token-template"] = string(f)
		}
		if gw.Spec.App.Log.Override {
			data["log-override-properties"] = gw.Spec.App.Log.Properties
		}
		if gw.Spec.App.Hazelcast.External {
			data["hazelcast-client.xml"] = `<hazelcast-client xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.hazelcast.com/schema/client-config http://www.hazelcast.com/schema/client-config/hazelcast-client-config-5.2.xsd" xmlns="http://www.hazelcast.com/schema/client-config"><instance-name>` + gw.Name + `-hazelcast-client</instance-name><network><cluster-members><address>` + gw.Spec.App.Hazelcast.Endpoint + `</address></cluster-members><redo-operation>true</redo-operation></network><connection-strategy async-start="false" reconnect-mode="ON"><connection-retry><cluster-connect-timeout-millis>-1</cluster-connect-timeout-millis></connection-retry></connection-strategy></hazelcast-client>`
		}

	case gw.Name + "-cwp-bundle":
		var bundle []byte
		bundle, dataCheckSum, _ = util.BuildCWPBundle(gw.Spec.App.ClusterProperties.Properties)
		data["cwp.json"] = string(bundle)
	case gw.Name + "-listen-port-bundle":
		var bundle []byte
		refreshOnKeyChanges := false

		if gw.Spec.App.ListenPorts.RefreshOnKeyChanges {
			refreshOnKeyChanges = true
		}
		if !gw.Spec.App.ListenPorts.Custom.Enabled {
			bundle, dataCheckSum, _ = util.BuildDefaultListenPortBundle(refreshOnKeyChanges)
		} else {
			bundle, dataCheckSum, _ = util.BuildCustomListenPortBundle(gw, refreshOnKeyChanges)
		}
		data["listen-ports.json"] = string(bundle)
	case gw.Name + "-repository-init-config":
		initContainerStaticConfig := InitContainerStaticConfig{}
		initContainerStaticConfig.Version = "1.0"
		for i := range gw.Status.RepositoryStatus {
			var localRef string
			if gw.Status.RepositoryStatus[i].Enabled && gw.Status.RepositoryStatus[i].Type == "static" {
				//if gw.Status.RepositoryStatus[i].Type == "static" {
				if gw.Status.RepositoryStatus[i].StateStoreReference == "" {
					if gw.Status.RepositoryStatus[i].StorageSecretName != "" {
						localRef = "/graphman/localref/" + gw.Status.RepositoryStatus[i].StorageSecretName + "/" + gw.Status.RepositoryStatus[i].Name + ".gz"
						initContainerStaticConfig.Repositories = append(initContainerStaticConfig.Repositories, RepositoryConfig{
							Name:                gw.Status.RepositoryStatus[i].Name,
							LocalReference:      localRef,
							SingletonExtraction: gw.Spec.App.SingletonExtraction,
						})
					} else {
						if !gw.Spec.App.Management.Database.Enabled {
							initContainerStaticConfig.Repositories = append(initContainerStaticConfig.Repositories, RepositoryConfig{
								Name:                gw.Status.RepositoryStatus[i].Name,
								Endpoint:            gw.Status.RepositoryStatus[i].Endpoint,
								Branch:              gw.Status.RepositoryStatus[i].Branch,
								RemoteName:          gw.Status.RepositoryStatus[i].RemoteName,
								Tag:                 gw.Status.RepositoryStatus[i].Tag,
								Auth:                "/graphman/secrets/" + gw.Status.RepositoryStatus[i].Name,
								SingletonExtraction: gw.Spec.App.SingletonExtraction,
							})
						}
					}
				} else {
					if gw.Status.RepositoryStatus[i].StateStoreReference != "" && !gw.Spec.App.Management.Database.Enabled {
						initContainerStaticConfig.Repositories = append(initContainerStaticConfig.Repositories, RepositoryConfig{
							Name:                gw.Status.RepositoryStatus[i].Name,
							StateStoreReference: gw.Status.RepositoryStatus[i].StateStoreReference,
							StateStoreKey:       gw.Status.RepositoryStatus[i].StateStoreKey,
							SingletonExtraction: gw.Spec.App.SingletonExtraction,
						})
					}
				}
				//	}
			}
		}
		initContainerStaticConfigBytes, _ := json.Marshal(initContainerStaticConfig)
		data["config.json"] = string(initContainerStaticConfigBytes)
	case gw.Name + "-otk-shared-init-config":
		// parse properties in raw JSON Object without []
		data["OTK_TYPE"] = strings.ToUpper(string(securityv1.OtkTypeSingle))
		data["OTK_SK_UPGRADE"] = "true" // only applies when this runs as a job
		data["OTK_UPDATE_DATABASE_CONNECTION"] = "true"
		data["OTK_DATABASE_PROPERTIES"] = "na"
		data["OTK_RO_DATABASE_PROPERTIES"] = "na"

		if gw.Spec.App.Otk.Type != "" {
			data["OTK_TYPE"] = strings.ToUpper(string(gw.Spec.App.Otk.Type))
		}
		if gw.Spec.App.Otk.Database.Properties != nil {
			dbPropertyBytes, _ := json.Marshal(gw.Spec.App.Otk.Database.Properties)
			dbPropertyString := strings.ReplaceAll(string(dbPropertyBytes), "[", "")
			dbPropertyString = strings.ReplaceAll(dbPropertyString, "]", "")
			data["OTK_DATABASE_PROPERTIES"] = base64.StdEncoding.EncodeToString([]byte(dbPropertyString))
			if gw.Spec.App.Otk.Database.CreateReadOnlySqlConnection {
				data["OTK_RO_DATABASE_PROPERTIES"] = base64.StdEncoding.EncodeToString([]byte(dbPropertyString))
			}
		}

		switch gw.Spec.App.Otk.Database.Type {
		case securityv1.OtkDatabaseTypeMySQL, securityv1.OtkDatabaseTypeOracle:

			/// client read only connection
			data["OTK_DATABASE_CONN_PROPERTIES"] = "na"
			data["OTK_DATABASE_TYPE"] = string(gw.Spec.App.Otk.Database.Type)
			data["OTK_DATABASE_NAME"] = gw.Spec.App.Otk.Database.Sql.DatabaseName
			data["OTK_JDBC_URL"] = gw.Spec.App.Otk.Database.Sql.JDBCUrl
			data["OTK_JDBC_DRIVER_CLASS"] = "com.mysql.cj.jdbc.Driver"
			if gw.Spec.App.Otk.Database.Sql.ConnectionProperties != nil {
				dbPropertyBytes, _ := json.Marshal(gw.Spec.App.Otk.Database.Sql.ConnectionProperties)
				dbPropertyString := strings.Replace(string(dbPropertyBytes), "[", "", 1)
				dbPropertyString = strings.Replace(dbPropertyString, "]", "", 1)
				data["OTK_DATABASE_CONN_PROPERTIES"] = base64.StdEncoding.EncodeToString([]byte(dbPropertyString))
			}

			if gw.Spec.App.Otk.Database.Sql.JDBCDriverClass != "" {
				data["OTK_JDBC_DRIVER_CLASS"] = string(gw.Spec.App.Otk.Database.Sql.JDBCDriverClass)
			}

		case securityv1.OtkDatabaseTypeCassandra:
			data["OTK_CASSANDRA_DRIVER_CONFIG"] = "na"
			data["OTK_DATABASE_TYPE"] = string(gw.Spec.App.Otk.Database.Type)
			data["OTK_CASSANDRA_CONNECTION_POINTS"] = gw.Spec.App.Otk.Database.Cassandra.ConnectionPoints
			data["OTK_CASSANDRA_PORT"] = strconv.Itoa(gw.Spec.App.Otk.Database.Cassandra.Port)
			data["OTK_CASSANDRA_KEYSPACE"] = gw.Spec.App.Otk.Database.Cassandra.Keyspace
			if gw.Spec.App.Otk.Database.Cassandra.DriverConfig != "" {
				data["OTK_CASSANDRA_DRIVER_CONFIG"] = base64.StdEncoding.EncodeToString([]byte(gw.Spec.App.Otk.Database.Cassandra.DriverConfig))
			}
		}

		if gw.Spec.App.Otk.Database.CreateReadOnlySqlConnection {
			data["OTK_CREATE_RO_DATABASE_CONN"] = "false"
			if !reflect.DeepEqual(gw.Spec.App.Otk.Database.SqlReadOnly, securityv1.OtkSql{}) {
				data["OTK_RO_DATABASE_CONNECTION_NAME"] = "OAuth_ReadOnly"

				data["OTK_RO_DATABASE_CONN_PROPERTIES"] = "na"
				data["OTK_CREATE_RO_DATABASE_CONN"] = "true"
				data["OTK_RO_DATABASE_NAME"] = gw.Spec.App.Otk.Database.SqlReadOnly.DatabaseName
				data["OTK_RO_JDBC_URL"] = gw.Spec.App.Otk.Database.SqlReadOnly.JDBCUrl
				data["OTK_RO_JDBC_DRIVER_CLASS"] = gw.Spec.App.Otk.Database.SqlReadOnly.JDBCDriverClass

				if gw.Spec.App.Otk.Database.SqlReadOnlyConnectionName != "" {
					data["OTK_RO_DATABASE_CONNECTION_NAME"] = gw.Spec.App.Otk.Database.SqlReadOnlyConnectionName
				}
				if gw.Spec.App.Otk.Database.SqlReadOnly.ConnectionProperties != nil {
					roDbPropertyBytes, _ := json.Marshal(gw.Spec.App.Otk.Database.SqlReadOnly.ConnectionProperties)
					dbPropertyString := strings.Replace("[", string(roDbPropertyBytes), "", 1)
					dbPropertyString = strings.Replace("]", dbPropertyString, "", 1)
					data["OTK_RO_DATABASE_CONN_PROPERTIES"] = base64.StdEncoding.EncodeToString([]byte(dbPropertyString))
				}
			}
		}

	case gw.Name + "-otk-install-init-config":
		data["OTK_INSTALL_MODE"] = "initContainer"
		data["BOOTSTRAP_DIR"] = "/opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/000OTK"
		data["OTK_INTEGRATE_WITH_PORTAL"] = "false"
		data["OTK_SKIP_INTERNAL_SERVER_TOOLS"] = "false"
		data["OTK_SKIP_POST_INSTALLATION_TASKS"] = "false"
		data["OTK_DATABASE_UPGRADE"] = strconv.FormatBool(gw.Spec.App.Otk.Database.DbUpgrade)
		data["OTK_INTERNAL_CERT_ENCODED"] = ""
		data["OTK_INTERNAL_CERT_ISS"] = ""
		data["OTK_INTERNAL_CERT_SERIAL"] = "12345"
		data["OTK_INTERNAL_CERT_SUB"] = ""
		data["OTK_DMZ_CERT_ENCODED"] = ""
		data["OTK_DMZ_CERT_ISS"] = ""
		data["OTK_DMZ_CERT_SERIAL"] = "12345"
		data["OTK_DMZ_CERT_SUB"] = ""
		data["OTK_DATABASE_CONNECTION_NAME"] = "OAuth"

		data["OTK_INTERNAL_GW_HOST"] = ""
		data["OTK_INTERNAL_GW_PORT"] = ""
		data["OTK_DMZ_GW_HOST"] = ""
		data["OTK_DMZ_GW_PORT"] = ""
		if gw.Spec.App.Otk.Overrides.Enabled {
			if gw.Spec.App.Otk.Overrides.BootstrapDirectory != "" {
				data["BOOTSTRAP_DIR"] = gw.Spec.App.Otk.Overrides.BootstrapDirectory
			}
			if gw.Spec.App.Otk.Overrides.SkipInternalServerTools {
				data["OTK_SKIP_INTERNAL_SERVER_TOOLS"] = "true"
			}
			if gw.Spec.App.Otk.Overrides.EnablePortalIntegration {
				data["OTK_INTEGRATE_WITH_PORTAL"] = "true"
			}
		}

		if gw.Spec.App.Otk.Database.ConnectionName != "" {
			data["OTK_DATABASE_CONNECTION_NAME"] = gw.Spec.App.Otk.Database.ConnectionName
		}

	case gw.Name + "-otk-db-init-config":
		data["OTK_TYPE"] = strings.ToUpper(string(securityv1.OtkTypeSingle))
		data["OTK_SK_UPGRADE"] = "false"

		if gw.Spec.App.Otk.Type != "" {
			data["OTK_TYPE"] = strings.ToUpper(string(gw.Spec.App.Otk.Type))
		}

		switch gw.Spec.App.Otk.Database.Type {
		case securityv1.OtkDatabaseTypeMySQL, securityv1.OtkDatabaseTypeOracle:
			//data["OTK_DATABASE_CONN_PROPERTIES"] = "na"
			data["OTK_DATABASE_TYPE"] = string(gw.Spec.App.Otk.Database.Type)
			data["OTK_DATABASE_NAME"] = gw.Spec.App.Otk.Database.Sql.DatabaseName
			data["OTK_JDBC_URL"] = gw.Spec.App.Otk.Database.Sql.JDBCUrl
			data["OTK_JDBC_DRIVER_CLASS"] = "com.mysql.cj.jdbc.Driver"
			if gw.Spec.App.Otk.Database.Sql.JDBCDriverClass != "" {
				data["OTK_JDBC_DRIVER_CLASS"] = string(gw.Spec.App.Otk.Database.Sql.JDBCDriverClass)

			}
			data["OTK_DATABASE_UPGRADE"] = strconv.FormatBool(gw.Spec.App.Otk.Database.DbUpgrade)
			data["OTK_CREATE_TEST_CLIENTS"] = "false"
			data["OTK_TEST_CLIENTS_REDIRECT_URL_PREFIX"] = ""
			//data["OTK_LIQUIBASE_OPERATION"] = "changelogSync"
			if gw.Spec.App.Otk.Overrides.Enabled {
				if gw.Spec.App.Otk.Overrides.CreateTestClients {
					data["OTK_CREATE_TEST_CLIENTS"] = "true"
					data["OTK_TEST_CLIENTS_REDIRECT_URL_PREFIX"] = gw.Spec.App.Otk.Overrides.TestClientsRedirectUrlPrefix
				}
			}
			data["OTK_DATABASE_WAIT_TIMEOUT"] = strconv.Itoa(gw.Spec.App.Otk.Database.Sql.DatabaseWaitTimeout)
		}
	}

	if dataCheckSum == "" {
		dataBytes, _ := json.Marshal(data)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		dataCheckSum = sha1Sum
	}

	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   gw.Namespace,
			Labels:      util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
			Annotations: map[string]string{"checksum/data": dataCheckSum},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: data,
	}
	return cmap
}

func setJVMHeapSize(gw *securityv1.Gateway, heapType string, percentage int) string {
	var jvmHeap string
	memLimit := gw.Spec.App.Resources.Limits.Memory()

	switch heapType {
	case "min":
		if gw.Spec.App.Java.JVMHeap.Calculate && memLimit.IsZero() && gw.Spec.App.Java.JVMHeap.MinDefault != "" {
			jvmHeap = gw.Spec.App.Java.JVMHeap.MinDefault
		}
	case "max":
		if gw.Spec.App.Java.JVMHeap.Calculate && memLimit.IsZero() && gw.Spec.App.Java.JVMHeap.MaxDefault != "" {
			jvmHeap = gw.Spec.App.Java.JVMHeap.MaxDefault
		}
	default:
		if gw.Spec.App.Java.JVMHeap.Calculate && memLimit.IsZero() && gw.Spec.App.Java.JVMHeap.Default != "" {
			jvmHeap = gw.Spec.App.Java.JVMHeap.Default
		}

	}

	if gw.Spec.App.Java.JVMHeap.Calculate && percentage != 0 && !memLimit.IsZero() {
		memMB := float64(memLimit.Value()) * 0.00000095367432 //binary conversion
		heapPercntg := float64(percentage) / 100.0
		heapMb := strconv.FormatInt(int64(memMB*heapPercntg), 10)
		jvmHeap = heapMb + "m"
	}
	if jvmHeap == "" {
		switch heapType {
		case "min":
			jvmHeap = "1g"
			if gw.Spec.App.Java.JVMHeap.MinDefault != "" {
				jvmHeap = gw.Spec.App.Java.JVMHeap.MinDefault
			}
		case "max":
			jvmHeap = "3g"
			if gw.Spec.App.Java.JVMHeap.MaxDefault != "" {
				jvmHeap = gw.Spec.App.Java.JVMHeap.MaxDefault
			}
		default:
			jvmHeap = "3g"
			if gw.Spec.App.Java.JVMHeap.Default != "" {
				jvmHeap = gw.Spec.App.Java.JVMHeap.Default
			}
		}
	}

	return jvmHeap
}
