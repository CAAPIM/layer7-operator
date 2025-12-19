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

type StateStoreType string

const (
	StateStoreTypeRedis StateStoreType = "redis"
)

type RedisType string

const (
	RedisTypeStandalone RedisType = "standalone"
	RedisTypeSentinel   RedisType = "sentinel"
)

// L7StateStoreSpec defines the desired state of L7StateStore
type L7StateStoreSpec struct {
	// StateStoreType currently only supports Redis
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StateStoreType"
	StateStoreType StateStoreType `json:"type,omitempty"`
	// Redis state store configuration
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Redis"
	Redis Redis `json:"redis,omitempty"`
}

// L7StateStoreStatus defines the observed state of L7StateStore
type L7StateStoreStatus struct {
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// L7StateStore is the Schema for the l7statestores API
type L7StateStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   L7StateStoreSpec   `json:"spec,omitempty"`
	Status L7StateStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// L7StateStoreList contains a list of L7StateStore
type L7StateStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []L7StateStore `json:"items"`
}

type Redis struct {
	Type           RedisType       `json:"type,omitempty"`
	ExistingSecret string          `json:"existingSecret,omitempty"`
	Tls            RedisTls        `json:"tls,omitempty"`
	Username       string          `json:"username,omitempty"`
	MasterPassword string          `json:"masterPassword,omitempty"`
	GroupName      string          `json:"groupName,omitempty"`
	StoreId        string          `json:"storeId,omitempty"`
	Standalone     RedisStandalone `json:"standalone,omitempty"`
	Sentinel       RedisSentinel   `json:"sentinel,omitempty"`
	Database       int             `json:"database,omitempty"`
}

type RedisTls struct {
	Enabled    bool   `json:"enabled,omitempty"`
	RedisCrt   string `json:"redisCrt,omitempty"`
	VerifyPeer bool   `json:"verifyPeer,omitempty"`
}

type RedisSentinel struct {
	Master string              `json:"master,omitempty"`
	Nodes  []RedisSentinelNode `json:"nodes,omitempty"`
}

type RedisSentinelNode struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type RedisStandalone struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

func init() {
	SchemeBuilder.Register(&L7StateStore{}, &L7StateStoreList{})
}
