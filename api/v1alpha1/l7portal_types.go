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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// L7PortalSpec defines the desired state of L7Portal
type L7PortalSpec struct {
	// PortalTenant is the tenantId of the API Developer Portal
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="PortalTenant"
	PortalTenant string `json:"portalTenant,omitempty"`
	//Labels - Custom Labels
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Labels"
	Labels map[string]string `json:"labels,omitempty"`
	// Enabled - if enabled this Portal and its APIs will be synced
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled"
	Enabled bool `json:"enabled,omitempty"`
	// Endoint - Portal endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Endpoint"
	Endpoint string `json:"endpoint,omitempty"`
	// Mode determines how or if the Portal is contacted
	// defaults to auto, options are auto, local. Local requires
	// enrollmentBundle to be set.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Mode"
	Mode string `json:"mode,omitempty"`
	// EnrollmentBundle - allows a custom enrollment bundle to be set in the Portal CR
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="EnrollmentBundle"
	EnrollmentBundle string `json:"enrollmentBundle,omitempty"`
	// Deployment Tags - determines which Gateway deployments these APIs will be applied to
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="DeploymentTags"
	DeploymentTags []string `json:"deploymentTags,omitempty"`
	// Auth - Portal credentials
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auth"
	Auth PortalAuth `json:"auth,omitempty"`
	// SyncIntervalSeconds how often the Portal CR is reconciled. Default is 10 seconds
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SyncIntervalSeconds"
	SyncIntervalSeconds int `json:"syncIntervalSeconds,omitempty"`
	// PortalManaged if PortalManaged is true the portal controller will not manage APIs and will be only be responsible for maintaining a list of L7Api Portal published Metadata.
	PortalManaged bool `json:"portalManaged,omitempty"`
}

// L7PortalStatus defines the observed state of L7Portal
type L7PortalStatus struct {
	Ready               bool             `json:"ready,omitempty"`
	GatewayProxies      []GatewayProxy   `json:"proxies,omitempty"`
	LastUpdated         int64            `json:"lastUpdated,omitempty"`
	EnrollmentBundle    EnrollmentBundle `json:"enrollmentBundle,omitempty"`
	ApiSummaryConfigMap string           `json:"apiSummaryConfigMap,omitempty"`
	ApiCount            int              `json:"apiCount,omitempty"`
	Checksum            string           `json:"checksum,omitempty"`
}

//+kubebuilder:object:root=true
// +operator-sdk:csv:customresourcedefinitions:resources={{ConfigMaps,v1}}
//+kubebuilder:subresource:status

// L7Portal is the Schema for the l7portals API
type L7Portal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   L7PortalSpec   `json:"spec,omitempty"`
	Status L7PortalStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// L7PortalList contains a list of L7Portal
type L7PortalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []L7Portal `json:"items"`
}

// PortalAuth
type PortalAuth struct {
	Endpoint           string `json:"endpoint,omitempty"`
	PapiClientId       string `json:"clientId,omitempty"`
	PapiClientSecret   string `json:"clientSecret,omitempty"`
	ExistingSecretName string `json:"existingSecretName,omitempty"`
}

// GatewayProxy
type GatewayProxy struct {
	Name string `json:"name,omitempty"`
	// Type - Ephemeral or DbBacked
	Type     string         `json:"type,omitempty"`
	Gateways []ProxyGateway `json:"gateways,omitempty"`
}

type ProxyGateway struct {
	Name         string `json:"name,omitempty"`
	Synchronised bool   `json:"synchronised,omitempty"`
	LastUpdated  string `json:"lastUpdated,omitempty"`
}

// EnrollmentBundle
type EnrollmentBundle struct {
	SecretName  string `json:"secretName,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

func init() {
	SchemeBuilder.Register(&L7Portal{}, &L7PortalList{})
}
