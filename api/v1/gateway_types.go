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
	Version string  `json:"version,omitempty"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	Host               string                       `json:"host,omitempty"`
	Conditions         []appsv1.DeploymentCondition `json:"conditions,omitempty"`
	Phase              corev1.PodPhase              `json:"phase,omitempty"`
	Gateway            []GatewayState               `json:"gateway,omitempty"`
	ObservedGeneration int64                        `json:"observedGeneration,omitempty"`
	Ready              int32                        `json:"ready,omitempty"`
	State              corev1.PodConditionType      `json:"state,omitempty"`
	Replicas           int32                        `json:"replicas,omitempty"`
	Version            string                       `json:"version,omitempty"`
	Image              string                       `json:"image,omitempty"`
	LabelSelectorPath  string                       `json:"labelSelectorPath,omitempty"`
	ManagementPod      string                       `json:"managementPod,omitempty"`
	RepositoryStatus   []GatewayRepositoryStatus    `json:"repositoryStatus,omitempty"`
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

type Restman struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Graphman struct {
	Enabled            bool   `json:"enabled,omitempty"`
	InitContainerImage string `json:"initContainerImage,omitempty"`
}

type Bundle struct {
	Type      string    `json:"type,omitempty"`
	Source    string    `json:"source,omitempty"`
	Name      string    `json:"name,omitempty"`
	ConfigMap ConfigMap `json:"configMap,omitempty"`
	CSI       CSI       `json:"csi,omitempty"`
}

type ConfigMap struct {
	DefaultMode *int32 `json:"defaultMode,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
	Name        string `json:"name,omitempty"`
}

type CSI struct {
	Driver           string `json:"driver,omitempty"`
	ReadOnly         bool   `json:"readOnly,omitempty"`
	VolumeAttributes `json:"volumeAttributes,omitempty"`
}

type VolumeAttributes struct {
	SecretProviderClass string `json:"secretProviderClass,omitempty"`
}

type License struct {
	Accept     bool   `json:"accept,omitempty"`
	SecretName string `json:"secretName,omitempty"`
}

type Database struct {
	Enabled  bool   `json:"enabled"`
	JDBCUrl  string `json:"jdbcUrl,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Image struct {
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type App struct {
	Annotations               map[string]string                 `json:"annotations,omitempty"`
	Labels                    map[string]string                 `json:"labels,omitempty"`
	ClusterProperties         ClusterProperties                 `json:"cwp,omitempty"`
	Java                      Java                              `json:"java,omitempty"`
	Management                Management                        `json:"management,omitempty"`
	System                    System                            `json:"system,omitempty"`
	UpdateStrategy            UpdateStrategy                    `json:"updateStrategy,omitempty"`
	Image                     string                            `json:"image,omitempty"`
	ImagePullSecrets          []corev1.LocalObjectReference     `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy           corev1.PullPolicy                 `json:"imagePullPolicy,omitempty"`
	ListenPorts               ListenPorts                       `json:"listenPorts,omitempty"`
	Replicas                  int32                             `json:"replicas,omitempty"`
	Service                   Service                           `json:"service,omitempty"`
	Bundle                    []Bundle                          `json:"bundle,omitempty"`
	RepositoryReferences      []RepositoryReference             `json:"repositoryReferences,omitempty"`
	Ingress                   Ingress                           `json:"ingress,omitempty"`
	Sidecars                  []corev1.Container                `json:"sidecars,omitempty"`
	InitContainers            []corev1.Container                `json:"initContainers,omitempty"`
	Resources                 PodResources                      `json:"resources,omitempty"`
	Autoscaling               Autoscaling                       `json:"autoscaling,omitempty"`
	ServiceAccountName        string                            `json:"serviceAccountName,omitempty"`
	Hazelcast                 Hazelcast                         `json:"hazelcast,omitempty"`
	Bootstrap                 Bootstrap                         `json:"bootstrap,omitempty"`
	Monitoring                Monitoring                        `json:"monitoring,omitempty"`
	ContainerSecurityContext  corev1.SecurityContext            `json:"containerSecurityContext,omitempty"`
	PodSecurityContext        corev1.PodSecurityContext         `json:"podSecurityContext,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	Tolerations               []corev1.Toleration               `json:"tolerations,omitempty"`
	Affinity                  corev1.Affinity                   `json:"affinity,omitempty"`
	PodDisruptionBudget       PodDisruptionBudgetSpec           `json:"pdb,omitempty"`
	NodeSelector              map[string]string                 `json:"nodeSelector,omitempty"`
}

type ClusterProperties struct {
	Enabled    bool       `json:"enabled,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Monitoring - experimental feature that creates
// an OTEL Collector and Prometheus Service Monitor
// custom collectors and service monitors can be created separately
type Monitoring struct {
	Enabled        bool           `json:"enabled,omitempty"`
	Otel           Otel           `json:"otel,omitempty"`
	ServiceMonitor ServiceMonitor `json:"serviceMonitor,omitempty"`
}

type Otel struct {
	Collector Collector `json:"collector,omitempty"`
}

type Collector struct {
	//Name   string
	Create bool `json:"create,omitempty"`
}

type ServiceMonitor struct {
	//Name   string
	Create bool `json:"create,omitempty"`
}

type Bootstrap struct {
	Script BootstrapScript `json:"script,omitempty"`
}

type BootstrapScript struct {
	Enabled bool `json:"enabled,omitempty"`
	//Cleanup bool `json:"cleanup,omitempty"`
}

// Layer7 Gateway instantiates the following HTTP(s) ports by default
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
	Name         string                   `json:"name,omitempty"`
	Enabled      bool                     `json:"enabled"`
	Directories  []string                 `json:"directories,omitempty"`
	Type         string                   `json:"type,omitempty"`
	Encryption   GraphmanBundleEncryption `json:"encryption,omitempty"`
	Notification Notification             `json:"notification,omitempty"`
}

// GraphmanBundleEncryption allows setting an encryption passphrase per repository reference
type GraphmanBundleEncryption struct {
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
