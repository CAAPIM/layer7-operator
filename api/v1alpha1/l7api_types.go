// Copyright (c) 2025 Broadcom Inc. and its subsidiaries. All Rights Reserved.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// L7ApiSpec defines the desired state of L7Api
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
// +kubebuilder:resource:shortName=api;apis

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
	TenantId                         string           `json:"tenantId,omitempty"`
	Uuid                             string           `json:"apiUuid,omitempty"`
	UuidStripped                     string           `json:"apiId,omitempty"`
	ServiceId                        string           `json:"serviceId,omitempty"`
	Name                             string           `json:"name,omitempty"`
	ApiEnabled                       bool             `json:"enabled,omitempty"`
	SsgUrl                           string           `json:"ssgUrl,omitempty"`
	SsgUrlBase64                     string           `json:"ssgUrlEncoded,omitempty"`
	LocationUrl                      string           `json:"locationUrl,omitempty"`
	PublishedTs                      int              `json:"publishedTs,omitempty"`
	CreateTs                         int              `json:"createTs,omitempty"`
	ModifyTs                         int              `json:"modifyTs,omitempty"`
	SsgServiceType                   string           `json:"ssgServiceType,omitempty"`
	PolicyTemplates                  []PolicyTemplate `json:"policyEntities,omitempty"`
	CustomFields                     []CustomField    `json:"customFieldValues,omitempty"`
	SecurePasswords                  []SecurePassword `json:"securePasswords,omitempty"`
	SecurePasswordIdsForUndeployment []string         `json:"securePasswordIdsForUndeployment,omitempty"`
	Checksum                         string           `json:"checksum,omitempty"`
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

type SecurePassword struct {
	Id          string `json:"id"`
	Value       string `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GatewayPodDeploymentCondition struct {
	Action     string `json:"action,omitempty"`
	ActionTime string `json:"actionTime,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
	Status     string `json:"status,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

type LinkedGatewayStatus struct {
	Name       string                          `json:"name,omitempty"`
	Deployment string                          `json:"deployment,omitempty"`
	Conditions []GatewayPodDeploymentCondition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&L7Api{}, &L7ApiList{})
}
