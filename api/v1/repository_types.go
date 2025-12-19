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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RepositoryType string

const (
	RepositoryTypeGit        RepositoryType = "git"
	RepositoryTypeLocal      RepositoryType = "local"
	RepositoryTypeHttp       RepositoryType = "http"
	RepositoryTypeStateStore RepositoryType = "statestore"
)

// RepositorySpec defines the desired state of Repository
type RepositorySpec struct {
	//Labels - Custom Labels
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Labels"
	Labels map[string]string `json:"labels,omitempty"`
	//Annotations - Custom Annotations
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Annotations"
	Annotations map[string]string `json:"annotations,omitempty"`
	// Enabled - if enabled this repository will be synced
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled"
	Enabled bool `json:"enabled,omitempty"`
	// Endoint - Git repository endpoint
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Endpoint"
	Endpoint string `json:"endpoint,omitempty"`
	// Type of Repository - git, http, local, statestore
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Type"
	Type RepositoryType `json:"type,omitempty"`
	// StateStoreReference which L7StateStore connection should be used to store or retrieve this key
	// if type is statestore this reference will read everything from the state store
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StateStoreReference"
	StateStoreReference string `json:"stateStoreReference,omitempty"`
	// StateStoreKey where the repository is stored in the L7StateStore
	// this only takes effect if type is statestore
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StateStoreReference"
	StateStoreKey string `json:"stateStoreKey,omitempty"`
	// LocalReference lets the Repository controller use a local Kubernetes Secret as a repository source
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="LocalReference"
	LocalReference LocalReference `json:"localReference,omitempty"`
	// RepositorySyncConfig defines how often this repository is synced
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="RepositorySyncConfig"
	RepositorySyncConfig RepositorySyncConfig `json:"sync,omitempty"`
	// Remote Name - defaults to "origin"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="RemoteName"
	RemoteName string `json:"remoteName,omitempty"`
	// Branch - specify which branch to clone
	// if branch and tag are both specified branch will take precedence and tag will be ignored
	// if branch and tag are both missing the entire repository will be cloned
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Branch"
	Branch string `json:"branch,omitempty"`
	// Tag - clone a specific tag.
	// tags do not change, once cloned this will not be checked for updates
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tag"
	Tag string `json:"tag,omitempty"`
	// Auth contains a reference to the credentials required to connect to your Git repository
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auth"
	Auth RepositoryAuth `json:"auth,omitempty"`
}

//+kubebuilder:object:root=true
// +operator-sdk:csv:customresourcedefinitions:resources={{ConfigMaps,v1},{Secrets,v1}}
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=repo;repos;l7repo;l7repos;l7repository;l7repositories

// Repository is the Schema for the repositories API
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec - Repository Spec
	Spec RepositorySpec `json:"spec,omitempty"`
	// Status - Repository Status
	Status RepositoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

type LocalReference struct {
	SecretName string `json:"secretName,omitempty"`
}

// RepositorySyncConfig
type RepositorySyncConfig struct {
	// Configure how frequently the remote is checked for new commits
	IntervalSeconds int `json:"interval,omitempty"`
}

// RepositoryAuth
type RepositoryAuth struct {
	// Vendor i.e. Github, Gitlab, BitBucket, Azure
	Vendor string `json:"vendor,omitempty"`
	// Auth Type defaults to basic, possible options are
	// none, basic or ssh
	Type RepositoryAuthType `json:"type,omitempty"`
	// Username repository username
	Username string `json:"username,omitempty"`
	// Password repository Password
	// password or token are acceptable
	Password string `json:"password,omitempty"`
	// Token repository Access Token
	Token string `json:"token,omitempty"`
	// SSHKey for Git SSH Authentication
	SSHKey string `json:"sshKey,omitempty"`
	// SSHKeyPass
	SSHKeyPass string `json:"sshKeyPass,omitempty"`
	// KnownHosts is required for SSH Auth
	KnownHosts string `json:"knownHosts,omitempty"`
	// ExistingSecretName reference an existing secret
	ExistingSecretName string `json:"existingSecretName,omitempty"`
}

type RepositoryAuthType string

const (
	RepositoryAuthTypeBasic RepositoryAuthType = "basic"
	RepositoryAuthTypeSSH   RepositoryAuthType = "ssh"
	RepositoryAuthTypeNone  RepositoryAuthType = "none"
)

// RepositoryStatus defines the observed state of Repository
type RepositoryStatus struct {
	// Name of the Repository
	Name string `json:"name,omitempty"`
	// Ready to apply to Gateway Deployments
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:displayName="Ready"
	Ready bool `json:"ready,omitempty"`
	// Commit is either current git commit that has been synced or a sha1sum of the http repository contents
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:displayName="Commit"
	Commit string `json:"commit,omitempty"`
	// Updated the last time this repository was successfully updated
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:displayName="Updated"
	Updated            string `json:"updated,omitempty"`
	Summary            string `json:"summary,omitempty"`
	LastAppliedSummary string `json:"lastAppliedSummary,omitempty"`
	Vendor             string `json:"vendor,omitempty"`
	// StorageSecretName is the Kubernetes Secret that this repository is stored in
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:displayName="StorageSecretName"
	StorageSecretName string `json:"storageSecretName,omitempty"`
	// StateStoreSynced whether or not the state store has been written to correctly
	// +operator-sdk:csv:customresourcedefinitions:type=status
	// +operator-sdk:csv:customresourcedefinitions:displayName="StateStoreVersion"
	StateStoreSynced bool `json:"stateStoreSynced"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}
