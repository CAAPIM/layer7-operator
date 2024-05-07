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

// L7ApiSpec defines the desired state of L7Api

/// TODO: Include internal/templategen API definition

type L7ApiSpec struct {
	// ServiceUrl on the API Gateway
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ServiceUrl"
	ServiceUrl string `json:"serviceUrl,omitempty"`
	// PortalPublished
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="PortalPublished"
	PortalPublished bool `json:"portalPublished,omitempty"`
	// L7Portal is the L7Portal that this API is associated with when Portal Published is true
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="L7Portal"
	L7Portal string `json:"l7Portal,omitempty"`
	// PortalMeta is reserved for the API Developer Portal
	PortalMeta PortalMeta `json:"portalMeta,omitempty"`
	// GraphmanBundle associated with this API
	// currently limited to Service and Fragments
	// auto generated when PortalMeta is set and PortalPublished is true
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="GraphmanBundle"
	GraphmanBundle string `json:"graphmanBundle,omitempty"`
	// DeploymentTags target Gateway deployments that this API should be published to
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="DeploymentTags"
	DeploymentTags []string `json:"deploymentTags,omitempty"`
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
	Ready    bool                  `json:"ready,omitempty"`
	Checksum string                `json:"checksum,omitempty"`
	Gateways []LinkedGatewayStatus `json:"gateways,omitempty"`
}

// PortalMeta contains layer7 portal API Specific Metadata
type PortalMeta struct {
	TenantId        string           `json:"tenantId,omitempty"`
	Uuid            string           `json:"apiUuid,omitempty"`
	UuidStripped    string           `json:"apiId,omitempty"`
	ServiceId       string           `json:"serviceId,omitempty"`
	Name            string           `json:"name,omitempty"`
	ApiEnabled      bool             `json:"enabled,omitempty"`
	SsgUrl          string           `json:"ssgUrl,omitempty"`
	SsgUrlBase64    string           `json:"ssgUrlEncoded,omitempty"`
	LocationUrl     string           `json:"locationUrl,omitempty"`
	PublishedTs     int              `json:"publishedTs,omitempty"`
	CreateTs        int              `json:"createTs,omitempty"`
	ModifyTs        int              `json:"modifyTs,omitempty"`
	SsgServiceType  string           `json:"ssgServiceType,omitempty"`
	PolicyTemplates []PolicyTemplate `json:"policyEntities,omitempty"`
	CustomFields    []CustomField    `json:"customFieldValues,omitempty"`
	Checksum        string           `json:"checksum,omitempty"`
}

type PolicyTemplate struct {
	Uuid                       string              `json:"policyEntityUuid"`
	ApiPolicyTemplateArguments []PolicyTemplateArg `json:"policyTemplateArguments"`
}

type PolicyTemplateArg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LinkedGatewayStatus struct {
	Name string `json:"name,omitempty"`
	//Phase       corev1.PodPhase `json:"phase,omitempty"`
	// Reason: Success/Failed to sync because of x
	Deployment string `json:"deployment,omitempty"`
	//Ready       bool            `json:"ready,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
}

func init() {
	SchemeBuilder.Register(&L7Api{}, &L7ApiList{})
}
