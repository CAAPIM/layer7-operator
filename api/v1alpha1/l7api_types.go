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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// L7ApiSpec defines the desired state of L7Api
type L7ApiSpec struct {
	// Name of the API
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name"
	Name string `json:"name,omitempty"`
	// ServiceUrl on the API Gateway
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ServiceUrl"
	ServiceUrl string `json:"serviceUrl,omitempty"`
	// PortalPublished
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="PortalPublished"
	PortalPublished bool `json:"portalPublished,omitempty"`
	// GraphmanBundle associated with this API
	// currently limited to Service and Fragments
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="GraphmanBundle"
	GraphmanBundle string `json:"graphmanBundle,omitempty"`
	// DeploymentTags target Gateway deployments that this API should be published to
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="DeploymentTags"
	DeploymentTags []string `json:"deploymentTags,omitempty"`
	// L7Portal is the L7Portal that this API is associated with when Portal Published
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="L7Portal"
	L7Portal string `json:"l7Portal,omitempty"`
}

//+kubebuilder:object:root=true
// +operator-sdk:csv:customresourcedefinitions:resources={{ConfigMaps,v1},{Secrets,v1}}
//+kubebuilder:subresource:status

// L7Api is the Schema for the l7apis API
type L7Api struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   L7ApiSpec   `json:"spec,omitempty"`
	Status L7ApiStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// L7ApiList contains a list of L7Api
type L7ApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []L7Api `json:"items"`
}

// L7ApiStatus defines the observed state of L7Api
type L7ApiStatus struct {
	Gateways []LinkedGatewayStatus `json:"gateways,omitempty"`
}

type LinkedGatewayStatus struct {
	Name        string          `json:"name,omitempty"`
	Phase       corev1.PodPhase `json:"phase,omitempty"`
	Deployment  string          `json:"deployment,omitempty"`
	Ready       bool            `json:"ready,omitempty"`
	LastUpdated string          `json:"lastUpdated,omitempty"`
	Checksum    string          `json:"checksum,omitempty"`
}

func init() {
	SchemeBuilder.Register(&L7Api{}, &L7ApiList{})
}
