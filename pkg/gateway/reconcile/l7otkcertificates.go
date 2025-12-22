/*
* Copyright (c) 2025 Broadcom. All rights reserved.
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
	"encoding/base64"
	"encoding/json"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncOtkCertificates(ctx context.Context, params Params) {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-certificates")
		return
	}

	if !gateway.Spec.App.Otk.Enabled {
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-certificates")
		return
	}

	err = applyOtkCertificates(ctx, params, gateway)
	if err != nil {
		params.Log.Info("failed to reconcile otk certificates", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
	}
}

func applyOtkCertificates(ctx context.Context, params Params, gateway *securityv1.Gateway) error {

	bundle := graphman.Bundle{}
	annotation := ""
	sha1Sum := ""

	switch gateway.Spec.App.Otk.Type {
	case securityv1.OtkTypeDMZ:
		internalSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.InternalOtkGatewayReference+"-otk-internal-certificates")
		sha1Sum = internalSecret.ObjectMeta.Annotations["checksum/data"]

		if err != nil {
			return err
		}
		annotation = "security.brcmlabs.com/" + gateway.Name + "-" + string(gateway.Spec.App.Otk.Type) + "-certificates"
		for k, v := range internalSecret.Data {
			bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
				Name:                      k,
				CertBase64:                base64.StdEncoding.EncodeToString(v),
				TrustAnchor:               true,
				VerifyHostname:            false,
				RevocationCheckPolicyType: "USE_DEFAULT",
				TrustedFor: []graphman.TrustedForType{
					"SSL",
					"SIGNING_SERVER_CERTS",
				},
			})
		}

	case securityv1.OtkTypeInternal:
		dmzSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.DmzOtkGatewayReference+"-otk-dmz-certificates")
		sha1Sum = dmzSecret.ObjectMeta.Annotations["checksum/data"]
		if err != nil {
			return err
		}
		annotation = "security.brcmlabs.com/" + gateway.Name + "-" + string(gateway.Spec.App.Otk.Type) + "-fips-users"
		for k, v := range dmzSecret.Data {
			bundle.FipUsers = append(bundle.FipUsers, &graphman.FipUserInput{
				Name:         k,
				ProviderName: "otk-fips-provider",
				SubjectDn:    "cn=" + k,
				CertBase64:   base64.RawStdEncoding.EncodeToString(v),
			})
		}
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "otk certificates", *podList, gateway, gwSecret, "", annotation, sha1Sum, true, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, params, "otk certificates", *gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
