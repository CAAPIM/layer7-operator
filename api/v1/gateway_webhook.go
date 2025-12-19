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

* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
*/

package v1

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
//var gatewaylog = logf.Log.WithName("gateway-resource")

func (r *Gateway) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithDefaulter(r).
		WithValidator(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-security-brcmlabs-com-v1-gateway,mutating=true,failurePolicy=fail,sideEffects=None,groups=security.brcmlabs.com,resources=gateways,verbs=create;update,versions=v1,name=mgateway.kb.io,admissionReviewVersions=v1

var _ admission.CustomDefaulter = &Gateway{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Gateway) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}

//+kubebuilder:webhook:path=/validate-security-brcmlabs-com-v1-gateway,mutating=false,failurePolicy=fail,sideEffects=None,groups=security.brcmlabs.com,resources=gateways,verbs=create;update,versions=v1,name=vgateway.kb.io,admissionReviewVersions=v1

var _ admission.CustomValidator = &Gateway{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Gateway) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	gateway, ok := obj.(*Gateway)
	if !ok {
		return nil, fmt.Errorf("expected a Gateway, received %T", obj)
	}
	return validateGateway(gateway)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Gateway) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	_, ok := oldObj.(*Gateway)
	if !ok {
		return nil, fmt.Errorf("expected a Gateway for oldObj, received %T", oldObj)
	}
	gateway, ok := newObj.(*Gateway)
	if !ok {
		return nil, fmt.Errorf("expected a Gateway for newObj, received %T", newObj)
	}
	return validateGateway(gateway)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Gateway) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	//gatewaylog.Info("validate delete", "name", r.Name)
	return []string{}, nil
}

func validateGateway(r *Gateway) (admission.Warnings, error) {

	warnings := admission.Warnings{}

	if !r.Spec.License.Accept {
		return warnings, fmt.Errorf("please accept the gateway license")
	}

	if r.Spec.License.SecretName == "" {
		return warnings, fmt.Errorf("please create a kubernetes secret with a valid gateway license for the major version you are using. https://github.com/CAAPIM/layer7-operator#create-a-simple-gateway")
	}

	for i, b := range r.Spec.App.Bundle {
		if b.Name == "" {
			return warnings, fmt.Errorf("please specify a bundle name in your gateway configuration, this needs to be the name of an existing kubernetes configmap or secret. index: %d", i)
		}
		if b.Source != "configmap" && b.Source != "secret" {
			return warnings, fmt.Errorf("please specify a bundle source in your gateway configuration, valid sources are configmap or secret. bundle name: %s index: %d", b.Name, i)
		}
		if b.Type != "restman" && b.Type != "graphman" {
			return warnings, fmt.Errorf("please specify a bundle type in your gateway configuration, valid types are restman or graphman. bundle name: %s index: %d", b.Name, i)
		}
	}

	for i, rr := range r.Spec.App.RepositoryReferences {
		if rr.Enabled {
			if rr.Name == "" {
				return warnings, fmt.Errorf("please specify a repository reference name in your gateway configuration. index: %d", i)
			}
			if rr.Type == "" {
				return warnings, fmt.Errorf("please specify a repository reference type in your gateway configuration. reference name: %s index: %d", rr.Name, i)
			}
		}
	}

	if r.Spec.App.PodDisruptionBudget.Enabled {
		if r.Spec.App.PodDisruptionBudget.MinAvailable != (intstr.IntOrString{}) && r.Spec.App.PodDisruptionBudget.MaxUnavailable != (intstr.IntOrString{}) {
			return warnings, fmt.Errorf("PodDisruptionBudget minAvailable and maxUnavailable cannot be both set")
		}
	}

	if r.Spec.App.CustomConfig.Enabled {
		for i, cc := range r.Spec.App.CustomConfig.Mounts {
			if cc.Name == "" {
				return warnings, fmt.Errorf("please set the custom config mount name. index: %d", i)
			}
			if cc.MountPath == "" {
				return warnings, fmt.Errorf("please set the custom config mount path. name: %s index: %d", cc.Name, i)
			}
			if cc.SubPath == "" {
				return warnings, fmt.Errorf("please set the custom config sub path. name: %s index: %d", cc.Name, i)
			}
			if cc.ConfigRef.Name == "" {
				return warnings, fmt.Errorf("please set the custom config ref name. name: %s index: %d", cc.Name, i)
			}
			if strings.ToLower(cc.ConfigRef.Type) != "configmap" && strings.ToLower(cc.ConfigRef.Type) != "secret" {
				return warnings, fmt.Errorf("please specify a config type in your custom config configuration, valid types are configmap or secret. name: %s index: %d", cc.Name, i)
			}
			if cc.ConfigRef.Item.Key == "" || cc.ConfigRef.Item.Path == "" {
				return warnings, fmt.Errorf("please specify an item with key and path in your custom config configuration. name: %s index: %d", cc.Name, i)
			}
		}
	}

	for i, es := range r.Spec.App.ExternalSecrets {
		if es.Enabled {
			if es.Name == "" {
				return warnings, fmt.Errorf("please specify an external key name in your gateway configuration, this needs to be the name of an existing kubernetes secret. index: %d", i)
			}
		}
	}

	for i, ek := range r.Spec.App.ExternalKeys {
		if ek.Enabled {
			if ek.Name == "" {
				return warnings, fmt.Errorf("please specify an external key name in your gateway configuration, this needs to be the name of an existing kubernetes secret (kubernetes.io/tls). name: %s index: %d", ek.Name, i)
			}

			if ek.KeyUsageType != "" {
				if ek.KeyUsageType != KeyUsageTypeDefaultSSL && ek.KeyUsageType != KeyUsageTypeDefaultCA && ek.KeyUsageType != KeyUsageAuditViewer && ek.KeyUsageType != KeyUsageAuditSigning {
					return warnings, fmt.Errorf("please specify a valid key usage type, valid types are SSL,CA,AUDIT_VIEWER,AUDIT_SIGNING. name: %s index: %d", ek.Name, i)
				}
			}

		}
	}

	if r.Spec.App.Hazelcast.External {
		if r.Spec.App.Hazelcast.Endpoint == "" {
			return warnings, fmt.Errorf("please specify the endpoint for your external Hazelcast server")
		}
	}

	if r.Spec.App.Management.SecretName == "" {
		warnings = append(warnings, "using an existing secret for gateway credentials is strongly recommended")
		if r.Spec.App.Management.Username == "" || r.Spec.App.Management.Password == "" {
			return warnings, fmt.Errorf("please specify management username and password")
		}
		if r.Spec.App.Management.Cluster.Password == "" {
			return warnings, fmt.Errorf("please specify cluster password")
		}
		if r.Spec.App.Management.Database.Enabled {
			if r.Spec.App.Management.Database.Username == "" || r.Spec.App.Management.Database.Password == "" {
				return warnings, fmt.Errorf("please specify database username and password")
			}
		}
	}

	if r.Spec.App.Management.Cluster.Hostname == "" {
		return warnings, fmt.Errorf("please specify cluster hostname")
	}

	if r.Spec.App.Management.Database.Enabled {
		if r.Spec.App.Management.Database.JDBCUrl == "" {
			return warnings, fmt.Errorf("please specify a jdbcUrl for the gateway database")

		}
	}

	if r.Spec.App.Management.Service.Enabled {
		if r.Spec.App.Management.Service.Type != v1.ServiceTypeClusterIP && r.Spec.App.Management.Service.Type != v1.ServiceTypeLoadBalancer && r.Spec.App.Management.Service.Type != v1.ServiceTypeNodePort {
			return warnings, fmt.Errorf("please specify a valid management service type, valid types are LoadBalancer, ClusterIP and NodePort")
		}
	}

	if r.Spec.App.Service.Enabled {
		if r.Spec.App.Service.Type != v1.ServiceTypeClusterIP && r.Spec.App.Service.Type != v1.ServiceTypeLoadBalancer && r.Spec.App.Service.Type != v1.ServiceTypeNodePort {
			return warnings, fmt.Errorf("please specify a valid service type, valid types are LoadBalancer, ClusterIP and NodePort")
		}
	}

	if r.Spec.App.ListenPorts.Custom.Enabled {
		for i, lp := range r.Spec.App.ListenPorts.Ports {
			if lp.Name == "" {
				return warnings, fmt.Errorf("please specify a ListenPort name, valid features are Published service message input, Administrative access, Browser-based administration and Built-in services. index: %d", i)

			}
			for _, mf := range lp.ManagementFeatures {
				if mf != "Published service message input" && mf != "Administrative access" && mf != "Browser-based administration" && mf != "Built-in services" {
					return warnings, fmt.Errorf("please specify valid management features, valid features are Published service message input, Administrative access, Browser-based administration and Built-in services. port: %s", lp.Name)
				}
			}
		}
	}

	return warnings, nil
}
