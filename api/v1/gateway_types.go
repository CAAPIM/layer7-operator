/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GatewaySpec defines the desired state of Gateway
type GatewaySpec struct {
	License License `json:"license,omitempty"`
	App     App     `json:"app,omitempty"`
	// Version references the Gateway release that this Operator is intended to be used with
	// while all supported container gateway versions will work, some functionality will not be available
	// like bootstrapping graphman bundles which is currently unique to 10.1.00_CR3
	Version string `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Gateway is the Schema for the Gateways API
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

// GatewayState tracks the status of Gateway Resources
type GatewayState struct {
	Name      string          `json:"name,omitempty"`
	Phase     corev1.PodPhase `json:"phase,omitempty"`
	Ready     bool            `json:"ready"`
	StartTime string          `json:"startTime,omitempty"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	//+operator-sdk:csv:customresourcedefinitions:type=status
	Host                string                       `json:"host,omitempty"`
	Conditions          []appsv1.DeploymentCondition `json:"conditions,omitempty"`
	Phase               corev1.PodPhase              `json:"phase,omitempty"`
	Gateway             []GatewayState               `json:"gateway,omitempty"`
	ObservedGeneration  int64                        `json:"observedGeneration,omitempty"`
	Ready               int32                        `json:"ready,omitempty"`
	State               corev1.PodConditionType      `json:"state,omitempty"`
	Replicas            int32                        `json:"replicas,omitempty"`
	Version             string                       `json:"version,omitempty"`
	Image               string                       `json:"image,omitempty"`
	LabelSelectorPath   string                       `json:"labelSelectorPath,omitempty"`
	ManagementPod       string                       `json:"managementPod,omitempty"`
	RepositoryStatus    []GatewayRepositoryStatus    `json:"repositoryStatus,omitempty"`
	PortalSyncStatus    []PortalSyncStatus           `json:"portalSyncStatus,omitempty"`
	ApiSyncStatus       ApiSyncStatus                `json:"ApiSyncStatus,omitempty"`
	PortalApiSyncStatus ApiSyncStatus                `json:"PortalApiSyncStatus,omitempty"`
}

// GatewayRepositoryStatus tracks the status of which Graphman repositories have been applied to the Gateway Resource.
type GatewayRepositoryStatus struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name,omitempty"`
	Commit  string `json:"commit,omitempty"`
	Type    string `json:"type,omitempty"`
	//SecretName is used to mount the correct repository secret to the initContainer
	SecretName string `json:"secretName,omitempty"`
	//StorageSecretName is used to mount existing repository bundles to the initContainer
	//these will be less than 1mb in size
	StorageSecretName string `json:"storageSecretName,omitempty"`
	Branch            string `json:"branch,omitempty"`
	Endpoint          string `json:"endpoint,omitempty"`
}

// PortalSyncStatus tracks the status of which portals are synced with a gateway.
type PortalSyncStatus struct {
	Name        string `json:"name,omitempty"`
	ApiCount    int    `json:"apiCount,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// L7ApiSyncStatus tracks the status of l7apis applied to a gateway.
type ApiSyncStatus struct {
	L7Apis      []ApiStatus `json:"apis,omitempty"`
	Count       int         `json:"count,omitempty"`
	LastUpdated string      `json:"lastUpdated,omitempty"`
}

type ApiStatus struct {
	Name string `json:"name,omitempty"`
	//Type is Gateway or Portal - defaults to Gateway
	Type        string `json:"type,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// Management defines configuration for Gateway Managment.
type Management struct {
	SecretName string   `json:"secretName,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Cluster    Cluster  `json:"cluster,omitempty"`
	Database   Database `json:"database,omitempty"`
	Restman    Restman  `json:"restman,omitempty"`
	Graphman   Graphman `json:"graphman,omitempty"`
	Service    Service  `json:"service,omitempty"`
}

// Restman is a Gateway Management interface that can be automatically provisioned.
type Restman struct {
	Enabled bool `json:"enabled,omitempty"`
}

// Graphman is a GraphQL Gateway Management interface that can be automatically provisioned.
// The initContainer image is required for bootstrapping graphman bundles defined by the repository controller.
type Graphman struct {
	Enabled            bool   `json:"enabled,omitempty"`
	DynamicSyncPort    int    `json:"dynamicSyncPort,omitempty"`
	InitContainerImage string `json:"initContainerImage,omitempty"`
}

// Bundle A Restman or Graphman bundle
type Bundle struct {
	Type   string `json:"type,omitempty"`
	Source string `json:"source,omitempty"`
	Name   string `json:"name,omitempty"`
	//Add secret...
	//Secret    Secret    `json:"secret,omitempty"`
	ConfigMap ConfigMap `json:"configMap,omitempty"`
	CSI       CSI       `json:"csi,omitempty"`
}

// ConfigMap
type ConfigMap struct {
	DefaultMode *int32 `json:"defaultMode,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
	Name        string `json:"name,omitempty"`
}

// type Secret struct {
// 	DefaultMode *int32 `json:"defaultMode,omitempty"`
// 	Optional    bool   `json:"optional,omitempty"`
// 	Name        string `json:"name,omitempty"`
// }

// CSI
type CSI struct {
	Driver           string `json:"driver,omitempty"`
	ReadOnly         bool   `json:"readOnly,omitempty"`
	VolumeAttributes `json:"volumeAttributes,omitempty"`
}

// VolumeAtttributes
type VolumeAttributes struct {
	SecretProviderClass string `json:"secretProviderClass,omitempty"`
}

// License is reference to a Kubernetes Secret Containing a Gateway v10/11.x license.
// license.accept must be set to true or the Gateway will not start.
type License struct {
	Accept     bool   `json:"accept,omitempty"`
	SecretName string `json:"secretName,omitempty"`
}

// Database configuration for the Gateway
type Database struct {
	Enabled  bool   `json:"enabled"`
	JDBCUrl  string `json:"jdbcUrl,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Image is the Gateway Image
type Image struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

// App contains Gateway specific deployment and application level configuration
type App struct {
	Annotations                        map[string]string                 `json:"annotations,omitempty"`
	PodAnnotations                     map[string]string                 `json:"podAnnotations,omitempty"`
	Labels                             map[string]string                 `json:"labels,omitempty"`
	ClusterProperties                  ClusterProperties                 `json:"cwp,omitempty"`
	Java                               Java                              `json:"java,omitempty"`
	Management                         Management                        `json:"management,omitempty"`
	System                             System                            `json:"system,omitempty"`
	UpdateStrategy                     UpdateStrategy                    `json:"updateStrategy,omitempty"`
	Image                              string                            `json:"image,omitempty"`
	ImagePullSecrets                   []corev1.LocalObjectReference     `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy                    corev1.PullPolicy                 `json:"imagePullPolicy,omitempty"`
	ListenPorts                        ListenPorts                       `json:"listenPorts,omitempty"`
	Replicas                           int32                             `json:"replicas,omitempty"`
	Service                            Service                           `json:"service,omitempty"`
	Bundle                             []Bundle                          `json:"bundle,omitempty"`
	SingletonExtraction                bool                              `json:"singletonExtraction,omitempty"`
	PortalReference                    PortalReference                   `json:"portalReference,omitempty"`
	PortalApiSyncIntervalSeconds       int                               `json:"portalApiSyncIntervalSeconds,omitempty"`
	RepositorySyncIntervalSeconds      int                               `json:"repositorySyncIntervalSeconds,omitempty"`
	ExternalSecretsSyncIntervalSeconds int                               `json:"externalSecretsSyncIntervalSeconds,omitempty"`
	ExternalKeysSyncIntervalSeconds    int                               `json:"externalKeysSyncIntervalSeconds,omitempty"`
	RepositoryReferences               []RepositoryReference             `json:"repositoryReferences,omitempty"`
	Ingress                            Ingress                           `json:"ingress,omitempty"`
	Sidecars                           []corev1.Container                `json:"sidecars,omitempty"`
	InitContainers                     []corev1.Container                `json:"initContainers,omitempty"`
	Resources                          PodResources                      `json:"resources,omitempty"`
	Autoscaling                        Autoscaling                       `json:"autoscaling,omitempty"`
	ServiceAccountName                 string                            `json:"serviceAccountName,omitempty"`
	Hazelcast                          Hazelcast                         `json:"hazelcast,omitempty"`
	Bootstrap                          Bootstrap                         `json:"bootstrap,omitempty"`
	Monitoring                         Monitoring                        `json:"monitoring,omitempty"`
	ContainerSecurityContext           corev1.SecurityContext            `json:"containerSecurityContext,omitempty"`
	PodSecurityContext                 corev1.PodSecurityContext         `json:"podSecurityContext,omitempty"`
	TopologySpreadConstraints          []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	Tolerations                        []corev1.Toleration               `json:"tolerations,omitempty"`
	Affinity                           corev1.Affinity                   `json:"affinity,omitempty"`
	PodDisruptionBudget                PodDisruptionBudgetSpec           `json:"pdb,omitempty"`
	NodeSelector                       map[string]string                 `json:"nodeSelector,omitempty"`
	ExternalSecrets                    []ExternalSecret                  `json:"externalSecrets,omitempty"`
	ExternalKeys                       []ExternalKey                     `json:"externalKeys,omitempty"`
}

// ClusterProperties are key value pairs of additional cluster-wide properties you wish to bootstrap to your Gateway.
type ClusterProperties struct {
	Enabled    bool       `json:"enabled,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

type SingletonExtraction struct {
	Enabled bool `json:"enabled,omitempty"`
	// ScheduledTasks bool `json:"scheduledTasks,omitempty"`
	// JmsListener    bool `json:"jmsListener,omitempty"`
}

// Property is a cluster-wide property k/v pair
type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ExternalSecret is a reference to an existing secret in Kubernetes
// The Layer7 Operator will attempt to convert this secret to a Graphman bundle that can be applied
// dynamically keeping any referenced secrets up-to-date.
// You can bring in external secrets using tools like the external secrets operator (external-secrets.io)
type ExternalSecret struct {
	Enabled              bool             `json:"enabled,omitempty"`
	Encryption           BundleEncryption `json:"encryption,omitempty"`
	Name                 string           `json:"name,omitempty"`
	Description          string           `json:"description,omitempty"`
	VariableReferencable bool             `json:"variableReferencable,omitempty"`
	Type                 string           `json:"type,omitempty"`
}

// ExternalKey is a reference to an existing TLS Secret in Kubernetes
// The Layer7 Operator will attempt to convert this secret to a Graphman bundle that can be applied
// dynamically keeping any referenced keys up-to-date.
// You can bring in external secrets using tools like cert-manager
type ExternalKey struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
	Port    string `json:"port,omitempty"`
}

// Monitoring - experimental feature that creates
// an OTEL Collector and Prometheus Service Monitor
// custom collectors and service monitors can be created separately
type Monitoring struct {
	Enabled        bool           `json:"enabled,omitempty"`
	Otel           Otel           `json:"otel,omitempty"`
	ServiceMonitor ServiceMonitor `json:"serviceMonitor,omitempty"`
}

// Otel
type Otel struct {
	Collector Collector `json:"collector,omitempty"`
}

// Collector is an OpenTelemetryCollector Configuration
type Collector struct {
	//Name   string
	Create bool `json:"create,omitempty"`
}

// ServiceMonitor is a Prom Service Monitor Configuration
type ServiceMonitor struct {
	//Name   string
	Create bool `json:"create,omitempty"`
}

// Bootstrap - optionally add a bootstrap script to the Gateway that migrates configuration from /opt/docker/custom to the correct Container Gateway locations for bootstrap
type Bootstrap struct {
	Script BootstrapScript `json:"script,omitempty"`
}

// BootstrapScript - enable/disable this functionality
type BootstrapScript struct {
	Enabled bool `json:"enabled,omitempty"`
	//Cleanup bool `json:"cleanup,omitempty"`
}

// ListenPorts The Layer7 Gateway instantiates the following HTTP(s) ports by default
// Harden applies the following changes, setting ports overrides this flag.
// - 8080 (HTTP)
//   - Disable
//   - Allow Published Service Message input only
//
// - 8443 (HTTPS)
//   - Remove Management Features (no Policy Manager Access)
//   - Enables TLSv1.2,TLS1.3 only
//   - Disables insecure Cipher Suites
//
// - 9443 (HTTPS)
//   - Enables TLSv1.2,TLS1.3 only
//   - Disables insecure Cipher Suites
//
// - 2124 (Inter-Node Communication)
//   - Not created
//   - if using an existing database 2124 will not be modified
type ListenPorts struct {
	Harden       bool             `json:"harden,omitempty"`
	CipherSuites []string         `json:"cipherSuites,omitempty"`
	TlsVersions  []string         `json:"tlsVersions,omitempty"`
	Custom       CustomListenPort `json:"custom,omitempty"`
	Ports        []ListenPort     `json:"ports,omitempty"`
}

// CustomListenPort - enable/disable custom listen ports
type CustomListenPort struct {
	Enabled bool `json:"enabled,omitempty"`
}

type ListenPort struct {
	Name               string     `json:"name,omitempty"`
	Enabled            bool       `json:"enabled,omitempty"`
	Protocol           string     `json:"protocol,omitempty"`
	Port               string     `json:"port,omitempty"`
	Tls                Tls        `json:"tls,omitempty"`
	ManagementFeatures []string   `json:"managementFeatures,omitempty"`
	Properties         []Property `json:"properties,omitempty"`
}

type Tls struct {
	Enabled              bool     `json:"enabled,omitempty"`
	PrivateKey           string   `json:"privateKey,omitempty"`
	ClientAuthentication string   `json:"clientAuthentication,omitempty"`
	Versions             []string `json:"versions,omitempty"`
	UseCipherSuitesOrder bool     `json:"useCipherSuitesOrder,omitempty"`
	CipherSuites         []string `json:"cipherSuites,omitempty"`
}

type Hazelcast struct {
	External bool   `json:"external,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type UpdateStrategy struct {
	Type          string                         `json:"type,omitempty"`
	RollingUpdate appsv1.RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

type Autoscaling struct {
	Enabled bool `json:"enabled,omitempty"`
	HPA     HPA  `json:"hpa,omitempty"`
}

type HPA struct {
	MinReplicas *int32                                        `json:"minReplicas,omitempty"`
	MaxReplicas int32                                         `json:"maxReplicas,omitempty"`
	Behavior    autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
	Metrics     []autoscalingv2.MetricSpec                    `json:"metrics,omitempty"`
}

type System struct {
	Properties string `json:"properties,omitempty"`
}

type RepositoryReference struct {
	Name         string           `json:"name,omitempty"`
	Enabled      bool             `json:"enabled"`
	Directories  []string         `json:"directories,omitempty"`
	Type         string           `json:"type,omitempty"`
	Encryption   BundleEncryption `json:"encryption,omitempty"`
	Singleton    bool             `json:"singleton,omitempty"`
	Notification Notification     `json:"notification,omitempty"`
}

type PortalReference struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
}

// BundleEncryption allows setting an encryption passphrase per repository or external secret/key reference
type BundleEncryption struct {
	// Passphrase - bundle encryption passphrase in plaintext
	Passphrase string `json:"passphrase,omitempty"`
	// ExistingSecret - reference to an existing secret
	ExistingSecret string `json:"existingSecret,omitempty"`
	// Key - the key in the kubernetes secret that the encryption passphrase is stored in.
	Key string `json:"key,omitempty"`
}

// This is currently configured for Slack
type Notification struct {
	Name    string  `json:"name,omitempty"`
	Enabled bool    `json:"enabled,omitempty"`
	Channel Channel `json:"channel,omitempty"`
}

type Channel struct {
	Webhook Webhook `json:"webhook,omitempty"`
}

type Webhook struct {
	Url                string            `json:"url,omitempty"`
	InsecureSkipVerify bool              `json:"insecureSkipVerify,omitempty"`
	Headers            map[string]string `json:"headers,omitempty"`
	Auth               WebhookAuth       `json:"auth,omitempty"`
}

type WebhookAuth struct {
	Type     string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type PodDisruptionBudgetSpec struct {
	Enabled        bool               `json:"enabled,omitempty"`
	MinAvailable   intstr.IntOrString `json:"minAvailable,omitempty"`
	MaxUnavailable intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

type PodResources struct {
	Requests corev1.ResourceList `json:"requests,omitempty"`
	Limits   corev1.ResourceList `json:"limits,omitempty"`
}

type ResourceList struct {
	Memory           resource.Quantity `json:"memory,omitempty"`
	CPU              resource.Quantity `json:"cpu,omitempty"`
	EphemeralStorage resource.Quantity `json:"ephemeral-storage,omitempty"`
}

type Cluster struct {
	Password string `json:"password,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

type Service struct {
	Enabled     bool               ` json:"enabled,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Type        corev1.ServiceType `json:"type,omitempty"`
	Ports       []Ports            `json:"ports,omitempty"`
}

type Ingress struct {
	Enabled          bool                       `json:"enabled,omitempty"`
	Annotations      map[string]string          `json:"annotations,omitempty"`
	IngressClassName string                     `json:"ingressClassName,omitempty"`
	TLS              []networkingv1.IngressTLS  `json:"tls,omitempty"`
	Rules            []networkingv1.IngressRule `json:"rules,omitempty"`
}

type Ports struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}

type Java struct {
	JVMHeap   JVMHeap  `json:"jvmHeap,omitempty"`
	ExtraArgs []string `json:"extraArgs,omitempty"`
}

type JVMHeap struct {
	Calculate  bool   `json:"calculate,omitempty"`
	Percentage int    `json:"percentage,omitempty"`
	Default    string `json:"default,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
