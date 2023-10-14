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
	//Labels - Custom Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Name Portal name
	Name string `json:"name,omitempty"`
	// Enabled - if enabled this Portal and its APIs will be synced
	Enabled bool `json:"enabled,omitempty"`
	// Endoint - Portal endpoint
	Endpoint string `json:"endpoint,omitempty"`
	// Mode determines how or if the Portal is contacted
	// defaults to auto, options are auto, local. Local requires
	// enrollmentBundle to be set.
	Mode string `json:"mode,omitempty"`
	// EnrollmentBundle - allows a custom enrollment bundle to be set in the Portal CR
	EnrollmentBundle string `json:"enrollmentBundle,omitempty"`
	// Deployment Tags - determines which Gateway deployments these APIs will be applied to
	DeploymentTags []string `json:"deploymentTags,omitempty"`
	// Auth - Portal credentials
	Auth PortalAuth `json:"auth,omitempty"`
}

// L7PortalStatus defines the observed state of L7Portal
type L7PortalStatus struct {
	Ready               bool             `json:"ready,omitempty"`
	GatewayProxies      []GatewayProxy   `json:"proxies,omitempty"`
	Updated             string           `json:"updated,omitempty"`
	EnrollmentBundle    EnrollmentBundle `json:"enrollmentBundle,omitempty"`
	ApiSummaryConfigMap string           `json:"apiSummaryConfigMap,omitempty"`
	ApiCount            int              `json:"apiCount,omitempty"`
}

//+kubebuilder:object:root=true
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