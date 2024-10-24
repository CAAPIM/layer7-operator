package reconcile

import (
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
)

func ExternalCerts(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalCerts) == 0 {
		for _, v := range gateway.Status.LastAppliedExternalCerts {
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

	for k, v := range gateway.Status.LastAppliedExternalCerts {
		found := false
		notFound := []string{}

		for _, ec := range gateway.Spec.App.ExternalCerts {
			if k == ec.Name {
				found = true
			}
		}
		if !found {
			notFound = append(notFound, v...)
			bundleBytes, err := util.ConvertCertsToGraphmanBundle(nil, notFound)
			if err != nil {
				return err
			}

			annotation := "security.brcmlabs.com/external-certs-" + k
			if !gateway.Spec.App.Management.Database.Enabled {
				err = ReconcileEphemeralGateway(ctx, params, "external certs", *podList, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			} else {
				err = ReconcileDBGateway(ctx, params, "external certs", gatewayDeployment, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, externalCert := range gateway.Spec.App.ExternalCerts {
		var sha1Sum string

		certSecretMap := []util.GraphmanCert{}
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

				crtStrings := strings.SplitAfter(string(v), "-----END CERTIFICATE-----")
				crtStrings = crtStrings[:len(crtStrings)-1]
				for crt := range crtStrings {
					b, _ := pem.Decode([]byte(crtStrings[crt]))
					crtX509, _ := x509.ParseCertificate(b.Bytes)

					revocationCheckPolicyType := string(graphman.PolicyUsageTypeUseDefault)
					if externalCert.RevocationCheckPolicyType == "" {
						revocationCheckPolicyType = string(graphman.PolicyUsageType(externalCert.RevocationCheckPolicyType))
					}

					gmanCert := util.GraphmanCert{
						Name:                      crtX509.Subject.CommonName,
						Crt:                       crtStrings[crt],
						VerifyHostname:            externalCert.VerifyHostname,
						TrustAnchor:               externalCert.TrustAnchor,
						TrustedFor:                trustedFor,
						RevocationCheckPolicyType: revocationCheckPolicyType,
						RevocationCheckPolicyName: externalCert.RevocationCheckPolicyName,
					}
					certSecretMap = append(certSecretMap, gmanCert)
				}

			}

			dataBytes, _ := json.Marshal(&secret.Data)
			h := sha1.New()
			h.Write(dataBytes)
			sha1Sum = fmt.Sprintf("%x", h.Sum(nil))
		}

		notFound := []string{}
		if gateway.Status.LastAppliedExternalCerts != nil && gateway.Status.LastAppliedExternalCerts[externalCert.Name] != nil {
			for _, appliedCert := range gateway.Status.LastAppliedExternalCerts[externalCert.Name] {
				found := false
				for _, desiredCert := range certSecretMap {
					if strings.Split(appliedCert, "-")[0] == desiredCert.Name {
						found = true
					}
				}
				if !found {
					notFound = append(notFound, appliedCert)
				}
			}
		}

		if len(certSecretMap) < 1 && len(notFound) < 1 {
			continue
		}

		bundleBytes, err := util.ConvertCertsToGraphmanBundle(certSecretMap, notFound)
		if err != nil {
			return err
		}

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
