package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

func ExternalKeys(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalKeys) == 0 {
		for _, v := range gateway.Status.LastAppliedExternalKeys {
			if len(v) != 0 {
				continue
			}
			return nil
		}
	}

	//var bundleBytes []byte

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

	for k, v := range gateway.Status.LastAppliedExternalKeys {
		found := false
		notFound := []string{}

		for _, ek := range gateway.Spec.App.ExternalKeys {
			if k == ek.Name {
				found = true
			}
		}
		if !found {
			notFound = append(notFound, v...)
			bundleBytes, err := util.ConvertX509ToGraphmanBundle(nil, notFound)
			if err != nil {
				return err
			}

			annotation := "security.brcmlabs.com/external-secret-" + k
			if !gateway.Spec.App.Management.Database.Enabled {
				err = ReconcileEphemeralGateway(ctx, params, "external keys", *podList, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			} else {
				err = ReconcileDBGateway(ctx, params, "external keys", gatewayDeployment, gateway, gwSecret, "", annotation, "deleted", false, k, bundleBytes)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, externalKey := range gateway.Spec.App.ExternalKeys {
		var sha1Sum string
		keySecretMap := []util.GraphmanKey{}
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

			dataBytes, _ := json.Marshal(&secret.Data)
			h := sha1.New()
			h.Write(dataBytes)
			sha1Sum = fmt.Sprintf("%x", h.Sum(nil))
		}

		notFound := []string{}
		if gateway.Status.LastAppliedExternalKeys != nil && gateway.Status.LastAppliedExternalKeys[externalKey.Name] != nil {
			for _, appliedKey := range gateway.Status.LastAppliedExternalKeys[externalKey.Name] {
				found := false
				for _, desiredKey := range keySecretMap {
					if appliedKey == desiredKey.Alias {
						found = true
					}
				}
				if !found {
					notFound = append(notFound, appliedKey)
				}
			}
		}

		if len(keySecretMap) < 1 && len(notFound) < 1 {
			continue
		}

		bundleBytes, err := util.ConvertX509ToGraphmanBundle(keySecretMap, notFound)
		if err != nil {
			return err
		}

		if sha1Sum == "" {
			sha1Sum = "deleted"
		}

		annotation := "security.brcmlabs.com/external-key-" + externalKey.Name

		if !gateway.Spec.App.Management.Database.Enabled {
			err = ReconcileEphemeralGateway(ctx, params, "external keys", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, externalKey.Name, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "external keys", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, externalKey.Name, bundleBytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
