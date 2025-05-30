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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *Repository) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-security-brcmlabs-com-v1-repository,mutating=true,failurePolicy=fail,sideEffects=None,groups=security.brcmlabs.com,resources=repositories,verbs=create;update,versions=v1,name=mrepository.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Repository{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Repository) Default() {
	//repositorylog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-security-brcmlabs-com-v1-repository,mutating=false,failurePolicy=fail,sideEffects=None,groups=security.brcmlabs.com,resources=repositories,verbs=create;update,versions=v1,name=vrepository.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Repository{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Repository) ValidateCreate() (admission.Warnings, error) {
	return validateRepository(r)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Repository) ValidateUpdate(obj runtime.Object) (admission.Warnings, error) {
	_, ok := obj.(*Repository)
	if !ok {
		return nil, fmt.Errorf("expected a Repository, received %T", obj)
	}
	return validateRepository(r)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Repository) ValidateDelete() (admission.Warnings, error) {
	// Could extend to checking which gateways reference this before deletion.
	return []string{}, nil
}

func validateRepository(r *Repository) (admission.Warnings, error) {
	// Could extend to checking the remote before the resource is created/updated
	warnings := admission.Warnings{}

	if r.Spec.Enabled {

		if r.Spec.Branch == "" && r.Spec.Tag == "" && string(r.Spec.Type) == "git" {
			return warnings, fmt.Errorf("please set a repository branch or tag. name: %s ", r.Name)
		}

		if r.Spec.Branch != "" && r.Spec.Tag != "" {
			warnings = append(warnings, "repository tag and branch are both set. branch supersedes tag.")
		}

		if r.Spec.Endpoint == "" {
			if strings.ToLower(string(r.Spec.Type)) != "local" && strings.ToLower(string(r.Spec.Type)) != "statestore" {
				return warnings, fmt.Errorf("please set a repository endpoint. name: %s ", r.Name)
			}
		}

		switch strings.ToLower(string(r.Spec.Type)) {
		case "git":
			if !strings.HasPrefix(r.Spec.Endpoint, "https://") && !strings.HasPrefix(r.Spec.Endpoint, "ssh://") {
				return warnings, fmt.Errorf("repository endpoint must start with https:// or ssh://. name: %s ", r.Name)
			}
			if r.Spec.Auth != (RepositoryAuth{}) {
				if r.Spec.Auth.Type != RepositoryAuthTypeNone && r.Spec.Auth.Type != RepositoryAuthTypeBasic && r.Spec.Auth.Type != RepositoryAuthTypeSSH {
					return warnings, fmt.Errorf("please set a valid auth type, valid options for Git are none, basic and ssh. name: %s ", r.Name)
				}
			}
		case "http":
			if !strings.HasPrefix(r.Spec.Endpoint, "https://") && !strings.HasPrefix(r.Spec.Endpoint, "ssh://") {
				return warnings, fmt.Errorf("repository endpoint must start with https://. name: %s ", r.Name)
			}
			if r.Spec.Auth != (RepositoryAuth{}) {
				if r.Spec.Auth.Type != RepositoryAuthTypeNone && r.Spec.Auth.Type != RepositoryAuthTypeBasic {
					return warnings, fmt.Errorf("please set a valid auth type, valid options for HTTP refs are none and basic. name: %s ", r.Name)
				}
			}
		case "local":
			if r.Spec.LocalReference.SecretName == "" {
				return warnings, fmt.Errorf("local repository type must reference an existing kubernetes secret. name: %s ", r.Name)
			}
		case "statestore":
			if r.Spec.StateStoreReference == "" {
				return warnings, fmt.Errorf("statestore repository type must reference an existing L7StateStore. name: %s ", r.Name)
			}
			if r.Spec.StateStoreKey == "" {
				return warnings, fmt.Errorf("statestore repository type must reference an existing key in the L7StateStore. name: %s ", r.Name)
			}
		default:
			return warnings, fmt.Errorf("please set a repository type, valid types are git, http, local and statestore. name: %s ", r.Name)
		}

		if r.Spec.Auth != (RepositoryAuth{}) {
			switch strings.ToLower(string(r.Spec.Auth.Type)) {
			case string(RepositoryAuthTypeNone):
				if strings.ToLower(string(r.Spec.Type)) != "local" {
					warnings = append(warnings, "it is strongly recommend using authentication for your remote repository "+r.Name)
				}
			case string(RepositoryAuthTypeBasic):
				if r.Spec.Auth.ExistingSecretName == "" {
					secret := r.Spec.Auth.Token
					if secret == "" {
						secret = r.Spec.Auth.Password
					}

					if r.Spec.Auth.Username == "" || secret == "" {
						return warnings, fmt.Errorf("please set a repository auth username and password or use an existingSecret. name: %s ", r.Name)
					}
				}
			case string(RepositoryAuthTypeSSH):
				if r.Spec.Auth.ExistingSecretName == "" {
					if r.Spec.Auth.KnownHosts == "" {
						return warnings, fmt.Errorf("please set knownHosts or use an existingSecret. name: %s ", r.Name)
					}
					if r.Spec.Auth.SSHKey == "" {
						return warnings, fmt.Errorf("please set an SSH Key or use an existingSecret. name: %s ", r.Name)
					}
					if r.Spec.Auth.SSHKeyPass == "" {
						warnings = append(warnings, "SSH Key used for repository "+r.Name+" is plaintext, consider encrypting it.")
					}
				}
			}
		} else {
			if r.Spec.StateStoreKey == "" {
				warnings = append(warnings, "using authentication for your remote repository is strongly recommended. repository: "+r.Name)
			}

		}

	}

	return warnings, nil
}
