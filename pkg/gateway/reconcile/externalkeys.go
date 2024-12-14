/*
* Copyright (c) 2024 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */

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
	if len(gateway.Spec.App.ExternalKeys) == 0 && len(gateway.Status.LastAppliedExternalKeys) == 0 {
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

	notFound := []string{}
	keySecretMap := []util.GraphmanKey{}
	for _, k := range gateway.Status.LastAppliedExternalKeys {
		found := false
		for _, ek := range gateway.Spec.App.ExternalKeys {
			if k == ek.Alias && ek.Enabled {
				found = true
			}
		}
		if !found {
			notFound = append(notFound, k)
		}
	}

	var sha1Sum string
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
	}

	if len(keySecretMap) < 1 && len(notFound) < 1 {
		return nil
	}

	bundleBytes, err := util.ConvertX509ToGraphmanBundle(keySecretMap, notFound)
	if err != nil {
		return err
	}

	dataBytes, _ := json.Marshal(&keySecretMap)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum = fmt.Sprintf("%x", h.Sum(nil))

	annotation := "security.brcmlabs.com/external-keys"

	if !gateway.Spec.App.Management.Database.Enabled {
		err = ReconcileEphemeralGateway(ctx, params, "external keys", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, "", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		err = ReconcileDBGateway(ctx, params, "external keys", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, "", bundleBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
