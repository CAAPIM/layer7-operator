package config

import (
	"encoding/json"
	"os"
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
	Name           string `json:"name"`
	Endpoint       string `json:"endpoint"`
	Branch         string `json:"branch"`
	Auth           string `json:"auth"`
	LocalReference string `json:"localReference,omitempty"`
}

// NewConfigMap
func NewConfigMap(gw *securityv1.Gateway, name string) *corev1.ConfigMap {
	javaArgs := strings.Join(gw.Spec.App.Java.ExtraArgs, " ")
	data := make(map[string]string)
	jvmHeap := setJVMHeapSize(gw)
	switch name {
	case gw.Name + "-system":
		data["system.properties"] = gw.Spec.App.System.Properties
	case gw.Name:
		data["ACCEPT_LICENSE"] = strconv.FormatBool(gw.Spec.License.Accept)
		data["SSG_CLUSTER_HOST"] = gw.Spec.App.Management.Cluster.Hostname
		data["SSG_JVM_HEAP"] = jvmHeap
		data["EXTRA_JAVA_ARGS"] = javaArgs

		if gw.Spec.App.Bootstrap.Script.Enabled {
			f, _ := os.ReadFile("./003-parse-custom-files.sh")
			data["003-parse-custom-files.sh"] = string(f)
		}
		if gw.Spec.App.Management.Database.Enabled {
			data["SSG_DATABASE_JDBC_URL"] = gw.Spec.App.Management.Database.JDBCUrl
		}

		if gw.Spec.App.Hazelcast.External {
			/// external hazelcast Dcom - update
			data["hazelcast-client.xml"] = `<hazelcast-client xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.hazelcast.com/schema/client-config http://www.hazelcast.com/schema/client-config/hazelcast-client-config-3.10.xsd" xmlns="http://www.hazelcast.com/schema/client-config"><instance-name>` + gw.Name + `-hazelcast-client</instance-name><network><cluster-members><address>` + gw.Spec.App.Hazelcast.Endpoint + `</address></cluster-members><connection-attempt-limit>10</connection-attempt-limit><redo-operation>true</redo-operation></network><connection-strategy async-start="false" reconnect-mode="ON" /></hazelcast-client>`
			data["EXTRA_JAVA_ARGS"] = javaArgs + " -Dcom.l7tech.server.extension.sharedCounterProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedKeyValueStoreProvider=externalhazelcast -Dcom.l7tech.server.extension.sharedClusterInfoProvider=externalhazelcast"
		}
	case gw.Name + "-cwp-bundle":
		props := map[string]string{}

		for _, p := range gw.Spec.App.ClusterProperties.Properties {
			props[p.Name] = p.Value
		}
		bundle, _ := util.BuildCWPBundle(props)
		data["cwp.bundle"] = string(bundle)
	case gw.Name + "-listen-port-bundle":
		bundle := []byte{}

		if !gw.Spec.App.ListenPorts.Custom.Enabled {
			bundle, _ = util.BuildDefaultListenPortBundle()

		} else {
			bundle, _ = util.BuildCustomListenPortBundle(gw)
		}

		data["listen-ports.bundle"] = string(bundle)
	case gw.Name + "-repository-init-config":
		initContainerStaticConfig := InitContainerStaticConfig{}
		initContainerStaticConfig.Version = "1.0"
		for i := range gw.Status.RepositoryStatus {
			if gw.Status.RepositoryStatus[i].Enabled && gw.Status.RepositoryStatus[i].Type == "static" {
				var localRef string
				if gw.Status.RepositoryStatus[i].StorageSecretName != "" {
					localRef = "/graphman/localref/" + gw.Status.RepositoryStatus[i].StorageSecretName + "/" + gw.Status.RepositoryStatus[i].Name + ".gz"
					initContainerStaticConfig.Repositories = append(initContainerStaticConfig.Repositories, RepositoryConfig{
						Name:           gw.Status.RepositoryStatus[i].Name,
						LocalReference: localRef,
					})
				} else {
					initContainerStaticConfig.Repositories = append(initContainerStaticConfig.Repositories, RepositoryConfig{
						Name:     gw.Status.RepositoryStatus[i].Name,
						Endpoint: gw.Status.RepositoryStatus[i].Endpoint,
						Branch:   gw.Status.RepositoryStatus[i].Branch,
						Auth:     "/graphman/secrets/" + gw.Status.RepositoryStatus[i].Name,
					})
				}

			}
		}

		initContainerStaticConfigBytes, _ := json.Marshal(initContainerStaticConfig)
		data["config.json"] = string(initContainerStaticConfigBytes)
	}

	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: gw.Namespace,
			Labels:    util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: data,
	}
	return cmap
}

func setJVMHeapSize(gw *securityv1.Gateway) string {
	var jvmHeap string
	memLimit := gw.Spec.App.Resources.Limits.Memory()

	if gw.Spec.App.Java.JVMHeap.Calculate && memLimit.IsZero() && gw.Spec.App.Java.JVMHeap.Default != "" {
		jvmHeap = gw.Spec.App.Java.JVMHeap.Default
	}

	if gw.Spec.App.Java.JVMHeap.Calculate && !memLimit.IsZero() {
		memMB := float64(memLimit.Value()) * 0.00000095367432 //binary conversion
		heapPercntg := float64(gw.Spec.App.Java.JVMHeap.Percentage) / 100.0
		heapMb := strconv.FormatInt(int64(memMB*heapPercntg), 10)
		jvmHeap = heapMb + "m"
	}

	if jvmHeap == "" {
		jvmHeap = "2g"
	}

	return jvmHeap
}
