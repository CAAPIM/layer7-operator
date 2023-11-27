package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncExternalSecrets(ctx context.Context, params Params) error {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "Name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = s.RemoveByTag(params.Instance.Name + "-sync-external-secrets")
		return nil
	}

	cntr := 0
	for _, externalSecret := range gateway.Spec.App.ExternalSecrets {
		if externalSecret.Enabled {
			cntr++
		}
	}
	if cntr == 0 {
		_ = s.RemoveByTag(params.Instance.Name + "-sync-external-secrets")
		return nil
	}

	err = applyExternalSecrets(ctx, params, gateway)
	if err != nil {
		params.Log.Info("failed to reconcile external secrets", "Name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
	}

	return nil
}

func applyExternalSecrets(ctx context.Context, params Params, gateway *securityv1.Gateway) error {

	name := gateway.Name
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
			err = ReconcileEphemeralGateway(ctx, params, "external secrets", *podList, gateway, gwSecret, graphmanEncryptionPassphrase, annotation, sha1Sum, true, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "otk policies", gatewayDeployment, gateway, gwSecret, graphmanEncryptionPassphrase, annotation, sha1Sum, false, bundleBytes)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
