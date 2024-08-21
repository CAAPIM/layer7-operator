package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

func ExternalSecrets(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalSecrets) == 0 {
		return nil
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

	for _, es := range gateway.Spec.App.ExternalSecrets {
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
			case corev1.SecretTypeBasicAuth:
			case corev1.SecretTypeServiceAccountToken:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 es.Name + "-" + k,
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			case corev1.SecretTypeDockerConfigJson:
			case corev1.SecretTypeDockercfg:
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
		}

		if len(opaqueSecretMap) < 1 {
			continue
		}

		bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(opaqueSecretMap)
		if err != nil {
			return err
		}
		sort.Slice(opaqueSecretMap, func(i, j int) bool {
			return opaqueSecretMap[i].Name < opaqueSecretMap[j].Name
		})

		opaqueSecretMapBytes, err := json.Marshal(opaqueSecretMap)
		if err != nil {
			return err
		}
		h := sha1.New()
		h.Write(opaqueSecretMapBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		graphmanEncryptionPassphrase := es.Encryption.Passphrase
		if es.Encryption.ExistingSecret != "" {
			graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, es.Encryption.ExistingSecret, es.Encryption.Key)
			if err != nil {
				return err
			}
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
