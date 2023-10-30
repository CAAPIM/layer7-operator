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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GatewaySpec defines the desired state of Gateway
type GatewaySpec struct {
	License License `json:"license,omitempty"`
	App     App     `json:"app,omitempty"`
	// Version references the Gateway release that this Operator is intended to be used with
	// while all supported container gateway versions will work, some functionality will not be available
	Version string `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=gws;gw;l7gw;l7gws;l7gateway;l7gateways

// Gateway is the Schema for the Gateway Custom Resource
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GatewayList contains a list of Gateways
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

// GatewayStatus defines the observed state of Gateways
type GatewayStatus struct {
	//+operator-sdk:csv:customresourcedefinitions:type=status
	// Host is the Gateway Cluster Hostname
	Host       string                       `json:"host,omitempty"`
	Conditions []appsv1.DeploymentCondition `json:"conditions,omitempty"`
	Phase      corev1.PodPhase              `json:"phase,omitempty"`
	Gateway    []GatewayState               `json:"gateway,omitempty"`
	Ready      int32                        `json:"ready,omitempty"`
	State      corev1.PodConditionType      `json:"state,omitempty"`
	// Replicas is the number of Gateway Pods
	Replicas int32 `json:"replicas,omitempty"`
	// Version of the Gateway
	Version string `json:"version,omitempty"`
	// Image of the Gateway
	Image string `json:"image,omitempty"`
	// Management Pod is a Gateway with a special annotation is used as a selector for the
	// management service and applying singleton resources
	ManagementPod    string                    `json:"managementPod,omitempty"`
	RepositoryStatus []GatewayRepositoryStatus `json:"repositoryStatus,omitempty"`
	PortalSyncStatus PortalSyncStatus          `json:"PortalSyncStatus,omitempty"`
}

// GatewayState tracks the status of Gateway Resources
type GatewayState struct {
	// Name of the Gateway Pod
	Name  string          `json:"name,omitempty"`
	Phase corev1.PodPhase `json:"phase,omitempty"`
	// Ready is the state of the Gateway pod
	Ready bool `json:"ready"`
	// StartTime is when the Gateway pod was started
	StartTime string `json:"startTime,omitempty"`
}

// GatewayRepositoryStatus tracks the status of which Graphman repositories have been applied to the Gateway Resource.
type GatewayRepositoryStatus struct {
	// Enabled shows whether or not this repository reference is enabled
	Enabled bool `json:"enabled"`
	// Name of the Repository Reference
	Name string `json:"name,omitempty"`
	// Commit is the last commit that was applied
	Commit string `json:"commit,omitempty"`
	// Type is static or dynamic
	Type string `json:"type,omitempty"`
	//SecretName is used to mount the correct repository secret to the initContainer
	SecretName string `json:"secretName,omitempty"`
	//StorageSecretName is used to mount existing repository bundles to the initContainer
	//these will be less than 1mb in size
	StorageSecretName string `json:"storageSecretName,omitempty"`
	// Branch of the Git repo
	Branch string `json:"branch,omitempty"`
	// Tag is the git tag in the Git repo
	Tag string `json:"tag,omitempty"`
	// Endoint is the Git repo
	Endpoint string `json:"endpoint,omitempty"`
}

// PortalSyncStatus tracks the status of which portals are synced with a gateway.
type PortalSyncStatus struct {
	// Name of the L7Portal
	Name string `json:"name,omitempty"`
	// ApiCount is number of APIs that are related to the Referenced Portal
	ApiCount int `json:"apiCount,omitempty"`
	// LastUpdated is the last time this status was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// License is reference to a Kubernetes Secret Containing a Gateway v10/11.x license.
// license.accept must be set to true or the Gateway will not start.
type License struct {
	Accept bool `json:"accept,omitempty"`
	// SecretName is the Kubernetes Secret that contains the Gateway license
	// There must be a key called license.xml
	SecretName string `json:"secretName,omitempty"`
}

// App contains Gateway specific deployment and application level configuration
type App struct {
	// Annotations for the Gateway Deployment
	Annotations map[string]string `json:"annotations,omitempty"`
	// PodAnnotations for Gateway Pods
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
	// Labels for the Gateway Deployment
	Labels            map[string]string `json:"labels,omitempty"`
	ClusterProperties ClusterProperties `json:"cwp,omitempty"`
	Java              Java              `json:"java,omitempty"`
	Management        Management        `json:"management,omitempty"`
	Log               Log               `json:"log,omitempty"`
	System            System            `json:"system,omitempty"`
	// AutoMountServiceAccountToken optionally adds the Gateway Container's Kubernetes Service Account Token to Stored Passwords
	AutoMountServiceAccountToken bool           `json:"autoMountServiceAccountToken,omitempty"`
	UpdateStrategy               UpdateStrategy `json:"updateStrategy,omitempty"`
	// Image is the Gateway image
	Image            string                        `json:"image,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy  corev1.PullPolicy             `json:"imagePullPolicy,omitempty"`
	ListenPorts      ListenPorts                   `json:"listenPorts,omitempty"`
	// Replicas to deploy, overridden if autoscaling is enabled
	Replicas int32    `json:"replicas,omitempty"`
	Service  Service  `json:"service,omitempty"`
	Bundle   []Bundle `json:"bundle,omitempty"`
	// SingletonExtraction works with the Gateway in Ephemeral mode.
	// this enables scheduled tasks that are set to execute on a single node and jms destinations that are outbound
	// to be applied to one ephemeral gateway only.
	// This works inconjunction with repository references and only supports dynamic repository references.
	SingletonExtraction bool            `json:"singletonExtraction,omitempty"`
	PortalReference     PortalReference `json:"portalReference,omitempty"`
	// RepositorySyncIntervalSeconds is the period of time between attempts to apply repository references to gateways.
	RepositorySyncIntervalSeconds int `json:"repositorySyncIntervalSeconds,omitempty"`
	// ExternalSecretsSyncIntervalSeconds is the period of time between attempts to apply external secrets to gateways.
	ExternalSecretsSyncIntervalSeconds int `json:"externalSecretsSyncIntervalSeconds,omitempty"`
	// ExternalKeysSyncIntervalSeconds is the period of time between attempts to apply external keys to gateways.
	ExternalKeysSyncIntervalSeconds int                   `json:"externalKeysSyncIntervalSeconds,omitempty"`
	RepositoryReferences            []RepositoryReference `json:"repositoryReferences,omitempty"`
	Ingress                         Ingress               `json:"ingress,omitempty"`
	Sidecars                        []corev1.Container    `json:"sidecars,omitempty"`
	InitContainers                  []corev1.Container    `json:"initContainers,omitempty"`
	Resources                       PodResources          `json:"resources,omitempty"`
	Autoscaling                     Autoscaling           `json:"autoscaling,omitempty"`
	// ServiceAccount to use for the Gateway Deployment
	ServiceAccount            ServiceAccount                    `json:"serviceAccount,omitempty"`
	Hazelcast                 Hazelcast                         `json:"hazelcast,omitempty"`
	Bootstrap                 Bootstrap                         `json:"bootstrap,omitempty"`
	ContainerSecurityContext  corev1.SecurityContext            `json:"containerSecurityContext,omitempty"`
	PodSecurityContext        corev1.PodSecurityContext         `json:"podSecurityContext,omitempty"`
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	Tolerations               []corev1.Toleration               `json:"tolerations,omitempty"`
	Affinity                  corev1.Affinity                   `json:"affinity,omitempty"`
	PodDisruptionBudget       PodDisruptionBudgetSpec           `json:"pdb,omitempty"`
	NodeSelector              map[string]string                 `json:"nodeSelector,omitempty"`
	ExternalSecrets           []ExternalSecret                  `json:"externalSecrets,omitempty"`
	ExternalKeys              []ExternalKey                     `json:"externalKeys,omitempty"`
	LivenessProbe             corev1.Probe                      `json:"livenessProbe,omitempty"`
	ReadinessProbe            corev1.Probe                      `json:"readinessProbe,omitempty"`
	CustomConfig              CustomConfig                      `json:"customConfig,omitempty"`
	// TerminationGracePeriodSeconds is the time kubernetes will wait for the Gateway to shutdown before forceably removing it
	TerminationGracePeriodSeconds int64            `json:"terminationGracePeriodSeconds,omitempty"`
	LifecycleHooks                corev1.Lifecycle `json:"lifecycleHooks,omitempty"`
	PreStopScript                 PreStopScript    `json:"preStopScript,omitempty"`
	CustomHosts                   CustomHosts      `json:"customHosts,omitempty"`
}

type ServiceAccount struct {
	// Create a service account for the Gateway Deployment
	Create bool `json:"create,omitempty"`
	// Name of the service account
	Name string `json:"name,omitempty"`
}

type CustomHosts struct {
	// Enabled or disabled
	Enabled     bool               `json:"enabled,omitempty"`
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`
}

// Management defines configuration for Gateway Managment.
type Management struct {
	// SecretName is reference to an existing secret that contains
	// SSG_ADMIN_USERNAME, SSG_ADMIN_PASSWORD, SSG_CLUSTER_PASSPHRASE and optionally
	// SSG_DATABASE_USER and SSG_DATABASE_PASSWORD for mysql backed gateway clusters
	SecretName string `json:"secretName,omitempty"`
	// Username is the Gateway Admin username
	Username string `json:"username,omitempty"`
	// Password is the Gateway Admin password
	Password string   `json:"password,omitempty"`
	Cluster  Cluster  `json:"cluster,omitempty"`
	Database Database `json:"database,omitempty"`
	Restman  Restman  `json:"restman,omitempty"`
	Graphman Graphman `json:"graphman,omitempty"`
	// Service is the Gateway Management Service
	Service Service `json:"service,omitempty"`
}

type Log struct {
	// Override default log properties
	Override   bool   `json:"override,omitempty"`
	Properties string `json:"properties,omitempty"`
}

// Cluster is gateway cluster configuration
type Cluster struct {
	// Password is the Gateway Cluster Passphrase
	Password string `json:"password,omitempty"`
	// Hostname is the Gateway Cluster Hostname
	Hostname string `json:"hostname,omitempty"`
}

// Database configuration for the Gateway
type Database struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// JDBCUrl for the Gateway
	JDBCUrl string `json:"jdbcUrl,omitempty"`
	// Username MySQL - can be set in management.secretName
	Username string `json:"username,omitempty"`
	// Password MySQL - can be set in management.secretName
	Password string `json:"password,omitempty"`
}

// Restman is a Gateway Management interface that can be automatically provisioned.
type Restman struct {
	// Enabled optionally bootstrap the Restman Gateway Managment API
	Enabled bool `json:"enabled,omitempty"`
}

// Graphman is a GraphQL Gateway Management interface that can be automatically provisioned.
// The initContainer image is required for bootstrapping graphman bundles defined by the repository controller.
type Graphman struct {
	// Enabled optionally bootstrap the GraphQL Gateway Management Service
	Enabled bool `json:"enabled,omitempty"`
	// DynamicSyncPort is the Port the Gateway controller uses to apply dynamic repositories, external keys/secrets to the Gateway
	DynamicSyncPort int `json:"dynamicSyncPort,omitempty"`
	// InitContainerImage is the image used to bootstrap static repositories
	InitContainerImage string `json:"initContainerImage,omitempty"`
}

// Service
type Service struct {
	// Enabled or disabled
	Enabled bool ` json:"enabled,omitempty"`
	// Annotations for the service
	Annotations map[string]string `json:"annotations,omitempty"`
	// Type ClusterIP, NodePort, LoadBalancer
	Type corev1.ServiceType `json:"type,omitempty"`
	// Ports exposed by the Service
	// These are appended to the Gateway deployment containerPorts
	Ports                         []Ports                             `json:"ports,omitempty"`
	ClusterIP                     string                              `json:"clusterIP,omitempty"`
	ClusterIPs                    []string                            `json:"clusterIPs,omitempty"`
	ExternalIPs                   []string                            `json:"externalIPs,omitempty"`
	SessionAffinity               corev1.ServiceAffinity              `json:"sessionAffinity,omitempty"`
	LoadBalancerIP                string                              `json:"loadBalancerIP,omitempty"`
	LoadBalancerSourceRanges      []string                            `json:"loadBalancerSourceRanges,omitempty"`
	ExternalName                  string                              `json:"externalName,omitempty"`
	ExternalTrafficPolicy         corev1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`
	HealthCheckNodePort           int32                               `json:"healthCheckNodePort,omitempty"`
	SessionAffinityConfig         corev1.SessionAffinityConfig        `json:"sessionAffinityConfig,omitempty"`
	IPFamilies                    []corev1.IPFamily                   `json:"ipFamilies,omitempty"`
	IPFamilyPolicy                corev1.IPFamilyPolicy               `json:"ipFamilyPolicy,omitempty"`
	AllocateLoadBalancerNodePorts *bool                               `json:"allocateLoadBalancerNodePorts,omitempty"`
	LoadBalancerClass             string                              `json:"loadBalancerClass,omitempty"`
	InternalTrafficPolicy         corev1.ServiceInternalTrafficPolicy `json:"internalTrafficPolicy,omitempty"`
}

// Ports
type Ports struct {
	// Name of the Port
	Name string `json:"name,omitempty"`
	// Port number
	Port int32 `json:"port,omitempty"`
	// TargetPort on the Gateway Application
	TargetPort int32 `json:"targetPort,omitempty"`
	// Protocol
	Protocol string `json:"protocol,omitempty"`
}

// Bundle A Restman or Graphman bundle
type Bundle struct {
	// Type can be restman or graphman
	Type string `json:"type,omitempty"`
	// Source
	Source string `json:"source,omitempty"`
	Name   string `json:"name,omitempty"`
	// ConfigMap ConfigMap `json:"configMap,omitempty"`
	CSI CSI `json:"csi,omitempty"`
}

// CSI volume configuration
type CSI struct {
	// Driver is the secretstore csi driver
	Driver string `json:"driver,omitempty"`
	// ReadOnly
	ReadOnly         bool `json:"readOnly,omitempty"`
	VolumeAttributes `json:"volumeAttributes,omitempty"`
}

// VolumeAtttributes
type VolumeAttributes struct {
	//SecretProviderClass
	SecretProviderClass string `json:"secretProviderClass,omitempty"`
}

// ClusterProperties are key value pairs of additional cluster-wide properties you wish to bootstrap to your Gateway.
type ClusterProperties struct {
	// Enabled bootstraps clusterProperties to the Gateway
	Enabled bool `json:"enabled,omitempty"`
	// Properties are key/value pairs
	Properties []Property `json:"properties,omitempty"`
}

// PreStopScript During upgrades and other events where Gateway pods are replaced you may have APIs/Services that have long running connections open.
// This functionality delays Kubernetes sending a SIGTERM to the container gateway while connections remain open. This works in conjunction with terminationGracePeriodSeconds which should always be higher than preStopScript.timeoutSeconds. If preStopScript.timeoutSeconds is exceeded, the script will exit 0 and normal pod termination will resume.
// The preStop script will monitor connections to inbound (not outbound) Gateway Application TCP ports (i.e. inbound listener ports opened by the Gateway Application and not some other process) except those that are explicitly excluded.
// The following ports are excluded from monitoring by default.
// 8777 (Hazelcast) - Embedded Hazelcast.
// 2124 (Internode-Communication) - not utilised by the Container Gateway.
// If there are no open connections, the preStop script will exit immediately ignoring preStopScript.timeoutSeconds to avoid unnecessary resource utilisation (pod stuck in terminating state) during upgrades.
// While there aren't any explicit limits on preStopScript.timeoutSeconds and terminationGracePeriodSeconds running these for extended periods of time (i.e. more than 5 minutes) may be less reliable where other Kubernetes processes may remove the pod before terminationGracePeriodSeconds is reached. If you do run services like this we recommend testing before any real life implementation or better, creating a dedicated workload without autoscaling enabled (HPA) where you have more control over when/how pods are replaced.
type PreStopScript struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// PeriodSeconds between checks
	PeriodSeconds int `json:"periodSeconds,omitempty"`
	// TimeoutSeconds is the total time this script should run
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
	// ExcludedPorts is an array of port numbers, if not set the defaults are 8777 and 2124
	ExcludedPorts []int `json:"excludedPorts,omitempty"`
}

// CustomConfig Certain folders on the Container Gateway are not writeable by design. This configuration allows you to mount existing configMap/Secret keys to specific paths on the Gateway without the need for a root user or a custom/derived image.
type CustomConfig struct {
	// Enabled or disabled
	Enabled bool                `json:"enabled,omitempty"`
	Mounts  []CustomConfigMount `json:"mounts,omitempty"`
}

// CustomConfigMount
type CustomConfigMount struct {
	// Name is the mount name
	Name string `json:"name,omitempty"`
	// MountPath is the location on the container gateway this should go
	MountPath string `json:"mountPath,omitempty"`
	// SubPath is the file name
	SubPath   string    `json:"subPath,omitempty"`
	ConfigRef ConfigRef `json:"ref,omitempty"`
}

// ConfigRef configures the secret or configmap for a CustomConfigMount
type ConfigRef struct {
	// Name of the Secret or Configmap which already exists in Kubernetes
	Name string `json:"name,omitempty"`
	// Type is secret or configmap
	Type string        `json:"type,omitempty"`
	Item ConfigRefItem `json:"item,omitempty"`
}

// ConfigRefItem is the key in the secret or configmap to mount, path is where it should be created.
type ConfigRefItem struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}

// Property is a simple k/v pair
type Property struct {
	// Name
	Name string `json:"name,omitempty"`
	// Value
	Value string `json:"value,omitempty"`
}

// ExternalSecret is a reference to an existing secret in Kubernetes
// The Layer7 Operator will attempt to convert this secret to a Graphman bundle that can be applied
// dynamically keeping any referenced secrets up-to-date.
// You can bring in external secrets using tools like the external secrets operator (external-secrets.io)
type ExternalSecret struct {
	// Enabled or disabled
	Enabled    bool             `json:"enabled,omitempty"`
	Encryption BundleEncryption `json:"encryption,omitempty"`
	// Name of the Opaque/Generic Secret which already exists in Kubernetes
	Name string `json:"name,omitempty"`
	// Description given the Stored Password in the Gateway
	Description string `json:"description,omitempty"`
	// VariableReferencable permits/restricts use of the Stored Password in policy
	VariableReferencable bool `json:"variableReferencable,omitempty"`
}

// ExternalKey is a reference to an existing TLS Secret in Kubernetes
// The Layer7 Operator will attempt to convert this secret to a Graphman bundle that can be applied
// dynamically keeping any referenced keys up-to-date.
// You can bring in external secrets using tools like cert-manager
type ExternalKey struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// Name of the kubernetes.io/tls Secret which already exists in Kubernetes
	Name string `json:"name,omitempty"`
	// Port is reserved for future use
	Port string `json:"port,omitempty"`
}

// Bootstrap - optionally add a bootstrap script to the Gateway that migrates configuration from /opt/docker/custom to the correct Container Gateway locations for bootstrap
type Bootstrap struct {
	Script BootstrapScript `json:"script,omitempty"`
}

// BootstrapScript - enable/disable this functionality
type BootstrapScript struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
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
	// Harden
	Harden bool             `json:"harden,omitempty"`
	Custom CustomListenPort `json:"custom,omitempty"`
	Ports  []ListenPort     `json:"ports,omitempty"`
}

// CustomListenPort - enable/disable custom listen ports
type CustomListenPort struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
}

// ListenPort is translated into a Restman Bundle
type ListenPort struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// Name of the listen port
	Name string `json:"name,omitempty"`
	// Protocol
	Protocol string `json:"protocol,omitempty"`
	// Port
	Port string `json:"port,omitempty"`
	Tls  Tls    `json:"tls,omitempty"`
	// ManagementFeatures that should be available on this port
	// - Published service message input
	// - Administrative access
	// - Browser-based administration
	// - Built-in services
	ManagementFeatures []string   `json:"managementFeatures,omitempty"`
	Properties         []Property `json:"properties,omitempty"`
}

// Tls configuration for Gateway Ports
type Tls struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// PrivateKey the Port should use
	PrivateKey string `json:"privateKey,omitempty"`
	// ClientAuthentication MTLS for the Port
	// None, Optional, Required
	ClientAuthentication string `json:"clientAuthentication,omitempty"`
	// Versions of TLS
	// - TLS1.0 (not recommended)
	// - TLS1.1 (not recommended)
	// - TLS1.2
	// - TLS1.3
	Versions []string `json:"versions,omitempty"`
	// UseCipherSuitesOrder
	UseCipherSuitesOrder bool `json:"useCipherSuitesOrder,omitempty"`
	// CipherSuites
	// 	- TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
	// 	- TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
	// 	- TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384
	// 	- TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384
	// 	- TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
	// 	- TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
	// 	- TLS_DHE_RSA_WITH_AES_256_GCM_SHA384
	// 	- TLS_DHE_RSA_WITH_AES_256_CBC_SHA256
	// 	- TLS_DHE_RSA_WITH_AES_256_CBC_SHA
	// 	- TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	// 	- TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	// 	- TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
	// 	- TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
	// 	- TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
	// 	- TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
	// 	- TLS_DHE_RSA_WITH_AES_128_GCM_SHA256
	// 	- TLS_DHE_RSA_WITH_AES_128_CBC_SHA256
	// 	- TLS_DHE_RSA_WITH_AES_128_CBC_SHA
	// 	- TLS_AES_256_GCM_SHA384
	// 	- TLS_AES_128_GCM_SHA256
	//  - TLS_ECDH_RSA_WITH_AES_256_GCM_SHA384 (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_256_GCM_SHA384 (Disabled by Harden)
	//  - TLS_ECDH_RSA_WITH_AES_256_CBC_SHA384 (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA384 (Disabled by Harden)
	//  - TLS_ECDH_RSA_WITH_AES_256_CBC_SHA (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_256_GCM_SHA384 (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_256_CBC_SHA256 (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_256_CBC_SHA (Disabled by Harden)
	//  - TLS_ECDH_RSA_WITH_AES_128_GCM_SHA256 (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_128_GCM_SHA256 (Disabled by Harden)
	//  - TLS_ECDH_RSA_WITH_AES_128_CBC_SHA256 (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA256 (Disabled by Harden)
	//  - TLS_ECDH_RSA_WITH_AES_128_CBC_SHA (Disabled by Harden)
	//  - TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_128_GCM_SHA256 (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_128_CBC_SHA256 (Disabled by Harden)
	//  - TLS_RSA_WITH_AES_128_CBC_SHA (Disabled by Harden)
	CipherSuites []string `json:"cipherSuites,omitempty"`
}

type Hazelcast struct {
	// External set to true adds config for an external Hazelcast instance to the Gateway
	External bool `json:"external,omitempty"`
	// Endpoint is the hazelcast server and port
	// my.hazelcast:5701
	Endpoint string `json:"endpoint,omitempty"`
}

// UpdateStrategy for the Gateway Deployment
type UpdateStrategy struct {
	Type          string                         `json:"type,omitempty"`
	RollingUpdate appsv1.RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

// Autoscaling configuration for the Gateway
type Autoscaling struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	HPA     HPA  `json:"hpa,omitempty"`
}

type HPA struct {
	// MinReplicas
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	// MaxReplicas
	MaxReplicas int32                                         `json:"maxReplicas,omitempty"`
	Behavior    autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
	Metrics     []autoscalingv2.MetricSpec                    `json:"metrics,omitempty"`
}

// System
type System struct {
	// Properties for the Gateway
	Properties string `json:"properties,omitempty"`
}

// RepositoryReference is reference to a Git repository or HTTP endpoint that contains graphman bundles
type RepositoryReference struct {
	// Enabled or disabled
	Enabled bool `json:"enabled"`
	// Name of the existing repository
	Name string `json:"name,omitempty"`
	// Directories from the remote repository to sync with the Gateway
	// Limited to dynamic type
	Directories []string `json:"directories,omitempty"`
	// Type static or dynamic
	// static repositories are bootstrapped to the container gateway using an initContainer
	// it is recommended that these stay under 1mb in size when compressed
	// for larger static repositories it is recommended that you use a dedicated initContainer
	// dynamic repositories are applied directly to the gateway whenever the commit of a repository changes
	Type         string           `json:"type,omitempty"`
	Encryption   BundleEncryption `json:"encryption,omitempty"`
	Notification Notification     `json:"notification,omitempty"`
}

// PortalReference
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

// PodDisruptionBudgetSpec
type PodDisruptionBudgetSpec struct {
	// Enabled or disabled
	Enabled        bool               `json:"enabled,omitempty"`
	MinAvailable   intstr.IntOrString `json:"minAvailable,omitempty"`
	MaxUnavailable intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

// PodResources
type PodResources struct {
	Requests corev1.ResourceList `json:"requests,omitempty"`
	Limits   corev1.ResourceList `json:"limits,omitempty"`
}

type Ingress struct {
	// Enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
	// Annotations for the ingress resource
	Annotations map[string]string `json:"annotations,omitempty"`
	// IngressClassName
	IngressClassName string `json:"ingressClassName,omitempty"`
	// TLS
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`
	// Rules
	Rules []networkingv1.IngressRule `json:"rules,omitempty"`
}

// Java configuration for the Gateway
type Java struct {
	JVMHeap JVMHeap `json:"jvmHeap,omitempty"`
	// ExtraArgs java
	ExtraArgs []string `json:"extraArgs,omitempty"`
}

type JVMHeap struct {
	// Calculate the JVMHeap size based on resource requests and limits
	// if resources are left unset this will be ignored
	Calculate bool `json:"calculate,omitempty"`
	// Percentage of requests.limits.memory to allocate to the jvm
	// 50% is the default, should be no higher than 75%
	Percentage int `json:"percentage,omitempty"`
	// Default Heap Size to use if calculate is false or requests.limits.memory is not set
	// Set to 2g
	Default string `json:"default,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
