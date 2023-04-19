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

// RepositorySpec defines the desired state of Repository
type RepositorySpec struct {
	//Labels - Custom Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Name Repository name
	Name string `json:"name,omitempty"`
	// Enabled - if enabled this repository will be synced
	Enabled bool `json:"enabled,omitempty"`
	// Endoint - Git repository endpoint
	Endpoint string `json:"endpoint,omitempty"`
	// LocalReference lets the Repository controller use a local Kubernetes Configmap/Secret as a repository source
	// This is not currently implemented
	LocalReference       LocalReference       `json:"localReference,omitempty"`
	RepositorySyncConfig RepositorySyncConfig `json:"sync,omitempty"`
	// Branch - specify which branch to check out
	Branch string `json:"branch,omitempty"`
	// Auth contains a reference to the credentials required to connect to your Git repository
	Auth RepositoryAuth `json:"auth,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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
	SecretName    string `json:"secretName,omitempty"`
	ConfigMapName string `json:"configMapName,omitempty"`
}

// RepositorySyncConfig
type RepositorySyncConfig struct {
	// Configure how frequently the remote is checked for new commits
	IntervalSeconds int64 `json:"interval,omitempty"`
}

// RepositoryAuth
type RepositoryAuth struct {
	// Vendor i.e. Github, Gitlab, BitBucket
	Vendor string `json:"vendor,omitempty"`
	// Username repository username
	Username string `json:"username,omitempty"`
	// Password repository Password
	// password or token are acceptable
	Password string `json:"password,omitempty"`
	// Token repository Access Token
	Token string `json:"token,omitempty"`
	// ExistingSecretName reference an existing secret
	ExistingSecretName string `json:"existingSecretName,omitempty"`
}

// RepositoryStatus defines the observed state of Repository
type RepositoryStatus struct {
	Name              string `json:"name,omitempty"`
	Commit            string `json:"commit,omitempty"`
	Updated           string `json:"updated,omitempty"`
	Summary           string `json:"summary,omitempty"`
	Vendor            string `json:"vendor,omitempty"`
	StorageSecretName string `json:"storageSecretName,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}
