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
)

func ExternalKeys(ctx context.Context, params Params) error {
	gateway := params.Instance
	keySecretMap := []util.GraphmanKey{}
	var bundleBytes []byte

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

	for _, externalKey := range gateway.Spec.App.ExternalKeys {
		if externalKey.Enabled {

			secret, err := getGatewaySecret(ctx, params, externalKey.Name)
			if err != nil {
				return err
			}

			usageType := ""
			switch strings.ToUpper(string(externalKey.KeyUsageType)) {
			case string(securityv1.KeyUsageTypeDefaultSSL), string(securityv1.KeyUsageTypeDefaultCA), string(securityv1.KeyUsageAuditSigning), string(securityv1.KeyUsageAuditViewer):
				usageType = strings.ToUpper(string(externalKey.KeyUsageType))
			}

			if secret.Type == corev1.SecretTypeTLS {
				keySecretMap = append(keySecretMap, util.GraphmanKey{
					Name:      secret.Name,
					Crt:       string(secret.Data["tls.crt"]),
					Key:       string(secret.Data["tls.key"]),
					Alias:     externalKey.Alias,
					UsageType: usageType,
				})
			}

		}

		if len(keySecretMap) <= 0 {
			return nil
		}

		bundleBytes, err = util.ConvertX509ToGraphmanBundle(keySecretMap)
		if err != nil {
			return err
		}

		sort.Slice(keySecretMap, func(i, j int) bool {
			return keySecretMap[i].Name < keySecretMap[j].Name
		})

		keySecretMapBytes, err := json.Marshal(keySecretMap)

		if err != nil {
			return err
		}
		h := sha1.New()
		h.Write(keySecretMapBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

		annotation := "security.brcmlabs.com/external-key-" + externalKey.Name

		if !gateway.Spec.App.Management.Database.Enabled {
			err = ReconcileEphemeralGateway(ctx, params, "external keys", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, externalKey.Name, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "external keys", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, bundleBytes)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
