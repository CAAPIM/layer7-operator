package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/caapim/layer7-operator/pkg/util"
)

func ExternalCerts(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalCerts) == 0 {
		return nil
	}
	certSecretMap := []util.GraphmanCert{}
	var bundleBytes []byte

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

	for _, externalCert := range gateway.Spec.App.ExternalCerts {
		if externalCert.Enabled {

			secret, err := getGatewaySecret(ctx, params, externalCert.Name)
			if err != nil {
				return err
			}

			for _, v := range secret.Data {
				if !strings.Contains(string(v), "-----BEGIN CERTIFICATE-----") {
					continue
				}

				trustedFor := []string{}

				for i := range externalCert.TrustedFor {
					trustedFor = append(trustedFor, string(externalCert.TrustedFor[i]))
				}

				graphmanCert := util.GraphmanCert{
					Name:                      externalCert.Name,
					Crt:                       string(v),
					VerifyHostname:            externalCert.VerifyHostname,
					RevocationCheckPolicyType: string(externalCert.RevocationCheckPolicyType),
					TrustedFor:                trustedFor,
				}

				certSecretMap = append(certSecretMap, graphmanCert)
			}
		}

		if len(certSecretMap) <= 0 {
			return nil
		}

		bundleBytes, err = util.ConvertCertsToGraphmanBundle(certSecretMap)
		if err != nil {
			return err
		}

		sort.Slice(certSecretMap, func(i, j int) bool {
			return certSecretMap[i].Name < certSecretMap[j].Name
		})

		keySecretMapBytes, err := json.Marshal(certSecretMap)

		if err != nil {
			return err
		}
		h := sha1.New()
		h.Write(keySecretMapBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

		annotation := "security.brcmlabs.com/external-certs-" + externalCert.Name

		if !gateway.Spec.App.Management.Database.Enabled {
			err = ReconcileEphemeralGateway(ctx, params, "external certs", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, externalCert.Name, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "external certs", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, externalCert.Name, bundleBytes)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
