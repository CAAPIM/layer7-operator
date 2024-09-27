package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

func ExternalSecrets(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalSecrets) == 0 {
		for _, v := range gateway.Status.LastAppliedExternalSecrets {
			if len(v) != 0 {
				continue
			}
			return nil
		}
	}

	name := gateway.Name
	if gateway.Spec.App.Management.DisklessConfig.Disabled {
		name = gateway.Name + "-node-properties"
	}
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	for k, v := range gateway.Status.LastAppliedExternalSecrets {
		found := false
		notFound := []string{}

		for _, es := range gateway.Spec.App.ExternalSecrets {
			if k == es.Name {
				found = true
			}
		}
		if !found {
			notFound = append(notFound, v...)
			bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(nil, notFound)
			if err != nil {
				return err
			}
			annotation := "security.brcmlabs.com/external-secret-" + k
			if !gateway.Spec.App.Management.Database.Enabled {
				err = ReconcileEphemeralGateway(ctx, params, "external secrets", *podList, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			} else {
				err = ReconcileDBGateway(ctx, params, "external secrets", gatewayDeployment, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			}
		}

	}

	for _, es := range gateway.Spec.App.ExternalSecrets {
		var sha1Sum string
		opaqueSecretMap := []util.GraphmanSecret{}
		if es.Enabled {
			secret, err := getGatewaySecret(ctx, params, es.Name)
			if err != nil {
				return err
			}

			switch secret.Type {
			case corev1.SecretTypeOpaque:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 k,
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			case corev1.SecretTypeServiceAccountToken, corev1.SecretTypeBasicAuth:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 es.Name + "-" + k,
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			case corev1.SecretTypeDockercfg, corev1.SecretTypeDockerConfigJson:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 es.Name + "-" + strings.Split(k, ".")[1],
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			default:
				params.Log.V(2).Info("not a supported secret type", "secret name", es.Name, "secret type", secret.Type, "name", gateway.Name, "namespace", gateway.Namespace)
			}

			dataBytes, _ := json.Marshal(&secret.Data)
			h := sha1.New()
			h.Write(dataBytes)
			sha1Sum = fmt.Sprintf("%x", h.Sum(nil))
		}

		notFound := []string{}
		if gateway.Status.LastAppliedExternalSecrets != nil && gateway.Status.LastAppliedExternalSecrets[es.Name] != nil {

			for _, appliedSecret := range gateway.Status.LastAppliedExternalSecrets[es.Name] {
				found := false
				for _, desiredSecret := range opaqueSecretMap {
					if appliedSecret == desiredSecret.Name {
						found = true
					}
				}
				if !found {
					notFound = append(notFound, appliedSecret)
				}
			}
		}

		if len(opaqueSecretMap) < 1 && len(notFound) < 1 {
			continue
		}

		bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(opaqueSecretMap, notFound)
		if err != nil {
			return err
		}

		graphmanEncryptionPassphrase := es.Encryption.Passphrase
		if es.Encryption.ExistingSecret != "" {
			graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, es.Encryption.ExistingSecret, es.Encryption.Key)
			if err != nil {
				return err
			}
		}

		if sha1Sum == "" {
			sha1Sum = "deleted"
		}

		annotation := "security.brcmlabs.com/external-secret-" + es.Name

		if !gateway.Spec.App.Management.Database.Enabled {
			err = ReconcileEphemeralGateway(ctx, params, "external secrets", *podList, gateway, gwSecret, graphmanEncryptionPassphrase, annotation, sha1Sum, false, es.Name, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "external secrets", gatewayDeployment, gateway, gwSecret, graphmanEncryptionPassphrase, annotation, sha1Sum, false, es.Name, bundleBytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
