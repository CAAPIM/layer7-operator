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
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ExternalKeys(ctx context.Context, params Params) error {
	gateway := params.Instance

	// Handle OTK keys if OTK is enabled
	if gateway.Spec.App.Otk.Enabled {
		err := handleOtkKeys(ctx, params, gateway)
		if err != nil {
			params.Log.Error(err, "failed to handle OTK keys", "name", gateway.Name, "namespace", gateway.Namespace)
			return err
		}
	}

	// Handle regular external keys
	if len(gateway.Spec.App.ExternalKeys) == 0 && len(gateway.Status.LastAppliedExternalKeys) == 0 {
		return nil
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeExternalKey),
	)

	if err != nil {
		return err
	}

	if gwUpdReq == nil {
		return nil
	}

	for _, extKey := range gwUpdReq.externalEntities {
		extKeyUpdReq := gwUpdReq
		extKeyUpdReq.bundle = extKey.Bundle
		extKeyUpdReq.bundleName = extKey.Name
		extKeyUpdReq.checksum = extKey.Checksum
		extKeyUpdReq.cacheEntry = extKey.CacheEntry
		extKeyUpdReq.patchAnnotation = extKey.Annotation
		err = SyncGateway(ctx, params, *extKeyUpdReq)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleOtkKeys(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	// Handle DMZ key updates only if there's an externalKey with otk: true referencing the same secret
	if gateway.Spec.App.Otk.DmzKeySecret != "" && gateway.Spec.App.Otk.Type == securityv1.OtkTypeDMZ {
		// Check if there's an externalKey with otk: true that references this secret
		hasOtkExternalKey := false
		for _, ek := range gateway.Spec.App.ExternalKeys {
			if ek.Enabled && ek.Otk && ek.Name == gateway.Spec.App.Otk.DmzKeySecret {
				hasOtkExternalKey = true
				break
			}
		}

		// Only process if there's an externalKey with otk: true
		if hasOtkExternalKey {
			err := handleDmzKeyUpdate(ctx, params, gateway)
			if err != nil {
				params.Log.Error(err, "failed to handle DMZ key update", "name", gateway.Name, "namespace", gateway.Namespace)
				return err
			}
		} else {
			params.Log.V(2).Info("Skipping DMZ key update - no externalKey with otk: true found", "secret", gateway.Spec.App.Otk.DmzKeySecret)
		}
	}

	// Handle Internal key updates only if there's an externalKey with otk: true referencing the same secret
	if gateway.Spec.App.Otk.InternalKeySecret != "" && gateway.Spec.App.Otk.Type == securityv1.OtkTypeInternal {
		// Check if there's an externalKey with otk: true that references this secret
		hasOtkExternalKey := false
		for _, ek := range gateway.Spec.App.ExternalKeys {
			if ek.Enabled && ek.Otk && ek.Name == gateway.Spec.App.Otk.InternalKeySecret {
				hasOtkExternalKey = true
				break
			}
		}

		// Only process if there's an externalKey with otk: true
		if hasOtkExternalKey {
			err := handleInternalKeyUpdate(ctx, params, gateway)
			if err != nil {
				params.Log.Error(err, "failed to handle Internal key update", "name", gateway.Name, "namespace", gateway.Namespace)
				return err
			}
		} else {
			params.Log.V(2).Info("Skipping Internal key update - no externalKey with otk: true found", "secret", gateway.Spec.App.Otk.InternalKeySecret)
		}
	}

	return nil
}

func handleDmzKeyUpdate(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	// Get DMZ key secret
	dmzKeySecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.DmzKeySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("DMZ key secret not found, skipping", "secret", gateway.Spec.App.Otk.DmzKeySecret)
			return nil
		}
		return err
	}

	// Check if operator managed (ephemeral mode)
	isOperatorManaged := !gateway.Spec.App.Management.Database.Enabled

	if isOperatorManaged {
		// Update DMZ with the new key
		err = updateDmzWithKey(ctx, params, gateway, dmzKeySecret)
		if err != nil {
			return fmt.Errorf("failed to update DMZ with key: %w", err)
		}
	}

	// Publish DMZ cert to Internal for Truststore & FIP User
	if gateway.Spec.App.Otk.InternalOtkGatewayReference != "" {
		err = publishDmzCertToInternal(ctx, params, gateway, dmzKeySecret)
		if err != nil {
			return fmt.Errorf("failed to publish DMZ cert to Internal: %w", err)
		}
	}

	return nil
}

func handleInternalKeyUpdate(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	// Get Internal key secret
	internalKeySecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.InternalKeySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("Internal key secret not found, skipping", "secret", gateway.Spec.App.Otk.InternalKeySecret)
			return nil
		}
		return err
	}

	// Check if operator managed (ephemeral mode)
	isOperatorManaged := !gateway.Spec.App.Management.Database.Enabled

	if isOperatorManaged {
		// Update Internal with the new key
		err = updateInternalWithKey(ctx, params, gateway, internalKeySecret)
		if err != nil {
			return fmt.Errorf("failed to update Internal with key: %w", err)
		}
	}

	// Publish Internal cert to DMZ for Truststore
	if gateway.Spec.App.Otk.DmzOtkGatewayReference != "" {
		err = publishInternalCertToDmz(ctx, params, gateway, internalKeySecret)
		if err != nil {
			return fmt.Errorf("failed to publish Internal cert to DMZ: %w", err)
		}
	}

	return nil
}

func updateDmzWithKey(ctx context.Context, params Params, gateway *securityv1.Gateway, keySecret *corev1.Secret) error {
	if keySecret.Type != corev1.SecretTypeTLS {
		return fmt.Errorf("DMZ key secret must be of type kubernetes.io/tls")
	}

	certData := keySecret.Data["tls.crt"]
	keyData := keySecret.Data["tls.key"]

	if len(certData) == 0 || len(keyData) == 0 {
		return fmt.Errorf("DMZ key secret must contain tls.crt and tls.key")
	}

	// Extract certificate from chain
	crtStrings := strings.SplitAfter(string(certData), "-----END CERTIFICATE-----")
	if len(crtStrings) == 0 {
		return fmt.Errorf("invalid certificate format in DMZ key secret")
	}

	// Use first certificate in chain
	firstCert := crtStrings[0]
	b, _ := pem.Decode([]byte(firstCert))
	if b == nil {
		return fmt.Errorf("failed to decode certificate")
	}
	crtX509, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Create Graphman key bundle
	keySecretMap := []util.GraphmanKey{
		{
			Name:      crtX509.Subject.CommonName,
			Crt:       string(certData),
			Key:       string(keyData),
			Alias:     "otk-dmz-key",
			UsageType: "SSL",
		},
	}

	bundleBytes, err := util.ConvertX509ToGraphmanBundle(keySecretMap, []string{})
	if err != nil {
		return fmt.Errorf("failed to convert key to bundle: %w", err)
	}

	// Calculate checksum
	dataBytes, _ := json.Marshal(&keySecretMap)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Get gateway secret for authentication
	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)
	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/otk-dmz-key"

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "otk dmz key", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, "otk dmz key", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, params, "otk dmz key", *gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, "otk dmz key", bundleBytes)
		if err != nil {
			return err
		}
	}

	// Update cluster property otk.dmz.private_key.name after DMZ key is updated
	//if err := updateDmzPrivateKeyClusterProperty(ctx, params, gateway, "otk-dmz-key"); err != nil {
	//	params.Log.V(2).Info("Failed to update DMZ private key cluster property", "error", err, "gateway", gateway.Name)
	//	// Don't fail the entire operation if cluster property update fails
	//}

	return nil
}

func updateInternalWithKey(ctx context.Context, params Params, gateway *securityv1.Gateway, keySecret *corev1.Secret) error {
	if keySecret.Type != corev1.SecretTypeTLS {
		return fmt.Errorf("Internal key secret must be of type kubernetes.io/tls")
	}

	certData := keySecret.Data["tls.crt"]
	keyData := keySecret.Data["tls.key"]

	if len(certData) == 0 || len(keyData) == 0 {
		return fmt.Errorf("Internal key secret must contain tls.crt and tls.key")
	}

	// Extract certificate from chain
	crtStrings := strings.SplitAfter(string(certData), "-----END CERTIFICATE-----")
	if len(crtStrings) == 0 {
		return fmt.Errorf("invalid certificate format in Internal key secret")
	}

	// Use first certificate in chain
	firstCert := crtStrings[0]
	b, _ := pem.Decode([]byte(firstCert))
	if b == nil {
		return fmt.Errorf("failed to decode certificate")
	}
	crtX509, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Create Graphman key bundle
	keySecretMap := []util.GraphmanKey{
		{
			Name:      crtX509.Subject.CommonName,
			Crt:       string(certData),
			Key:       string(keyData),
			Alias:     "otk-internal-key",
			UsageType: "SSL",
		},
	}

	bundleBytes, err := util.ConvertX509ToGraphmanBundle(keySecretMap, []string{})
	if err != nil {
		return fmt.Errorf("failed to convert key to bundle: %w", err)
	}

	// Calculate checksum
	dataBytes, _ := json.Marshal(&keySecretMap)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Get gateway secret for authentication
	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)
	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/otk-internal-key"

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "otk internal key", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, "otk internal key", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, params, "otk internal key", *gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, "otk internal key", bundleBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func publishDmzCertToInternal(ctx context.Context, params Params, gateway *securityv1.Gateway, dmzKeySecret *corev1.Secret) error {
	// Get Internal gateway
	internalGateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{
		Name:      gateway.Spec.App.Otk.InternalOtkGatewayReference,
		Namespace: gateway.Namespace,
	}, internalGateway)

	isExternalGateway := false
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Gateway not found - check if it's external (port specified)
			if gateway.Spec.App.Otk.InternalGatewayPort != 0 {
				params.Log.V(2).Info("Internal gateway not found but port specified, treating as external",
					"gateway", gateway.Spec.App.Otk.InternalOtkGatewayReference,
					"port", gateway.Spec.App.Otk.InternalGatewayPort)
				isExternalGateway = true
			} else {
				params.Log.V(2).Info("Internal gateway not found and no port specified, skipping cert publish",
					"gateway", gateway.Spec.App.Otk.InternalOtkGatewayReference)
				return nil
			}
		} else {
			return err
		}
	}

	certData := dmzKeySecret.Data["tls.crt"]
	if len(certData) == 0 {
		return fmt.Errorf("DMZ key secret must contain tls.crt")
	}

	// Parse certificate
	crtStrings := strings.SplitAfter(string(certData), "-----END CERTIFICATE-----")
	if len(crtStrings) == 0 {
		return fmt.Errorf("invalid certificate format")
	}

	bundle := graphman.Bundle{}

	// Add to TrustedCerts
	for _, certStr := range crtStrings {
		if certStr == "" {
			continue
		}
		b, _ := pem.Decode([]byte(certStr))
		if b == nil {
			continue
		}
		crtX509, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			continue
		}

		bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
			Name:                      crtX509.Subject.CommonName,
			CertBase64:                base64.StdEncoding.EncodeToString([]byte(certStr)),
			TrustAnchor:               true,
			VerifyHostname:            false,
			RevocationCheckPolicyType: "USE_DEFAULT",
			TrustedFor: []graphman.TrustedForType{
				"SSL",
				"SIGNING_SERVER_CERTS",
			},
		})

		// Add to FIP Users
		bundle.FipUsers = append(bundle.FipUsers, &graphman.FipUserInput{
			Name:         crtX509.Subject.CommonName,
			ProviderName: "otk-fips-provider",
			SubjectDn:    "cn=" + crtX509.Subject.CommonName,
			CertBase64:   base64.RawStdEncoding.EncodeToString(crtX509.Raw),
		})
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	// Calculate checksum
	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// If gateway is external (not managed by operator), use specified port
	if isExternalGateway {
		// For external gateways, we can't use the standard reconciliation
		// Log that external gateway support requires additional configuration
		params.Log.V(2).Info("External Internal gateway detected, port will be used for connection",
			"gateway", gateway.Spec.App.Otk.InternalOtkGatewayReference,
			"port", gateway.Spec.App.Otk.InternalGatewayPort)
		// Note: External gateway connection would require additional implementation
		// For now, we skip reconciliation for external gateways
		return nil
	}

	// Get Internal gateway secret
	name := internalGateway.Name
	if internalGateway.Spec.App.Management.SecretName != "" {
		name = internalGateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)
	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/" + gateway.Name + "-dmz-certificates"

	internalParams := params
	internalParams.Instance = internalGateway

	// Note: InternalGatewayPort is used when the gateway is external (not found in cluster)
	// For operator-managed gateways, the gateway's own graphman port configuration is used

	if !internalGateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, internalParams)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, internalParams, "otk certificates", *podList, internalGateway, gwSecret, "", annotation, sha1Sum, true, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, internalParams)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, internalParams, "otk certificates", *gatewayDeployment, internalGateway, gwSecret, "", annotation, sha1Sum, false, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func publishInternalCertToDmz(ctx context.Context, params Params, gateway *securityv1.Gateway, internalKeySecret *corev1.Secret) error {
	// Get DMZ gateway
	dmzGateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{
		Name:      gateway.Spec.App.Otk.DmzOtkGatewayReference,
		Namespace: gateway.Namespace,
	}, dmzGateway)

	isExternalGateway := false
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Gateway not found - check if it's external (port specified)
			if gateway.Spec.App.Otk.DmzGatewayPort != 0 {
				params.Log.V(2).Info("DMZ gateway not found but port specified, treating as external",
					"gateway", gateway.Spec.App.Otk.DmzOtkGatewayReference,
					"port", gateway.Spec.App.Otk.DmzGatewayPort)
				isExternalGateway = true
			} else {
				params.Log.V(2).Info("DMZ gateway not found and no port specified, skipping cert publish",
					"gateway", gateway.Spec.App.Otk.DmzOtkGatewayReference)
				return nil
			}
		} else {
			return err
		}
	}

	certData := internalKeySecret.Data["tls.crt"]
	if len(certData) == 0 {
		return fmt.Errorf("Internal key secret must contain tls.crt")
	}

	// Parse certificate
	crtStrings := strings.SplitAfter(string(certData), "-----END CERTIFICATE-----")
	if len(crtStrings) == 0 {
		return fmt.Errorf("invalid certificate format")
	}

	bundle := graphman.Bundle{}

	// Add to TrustedCerts
	for _, certStr := range crtStrings {
		if certStr == "" {
			continue
		}
		b, _ := pem.Decode([]byte(certStr))
		if b == nil {
			continue
		}
		crtX509, err := x509.ParseCertificate(b.Bytes)
		if err != nil {
			continue
		}

		bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
			Name:                      crtX509.Subject.CommonName,
			CertBase64:                base64.StdEncoding.EncodeToString([]byte(certStr)),
			TrustAnchor:               true,
			VerifyHostname:            false,
			RevocationCheckPolicyType: "USE_DEFAULT",
			TrustedFor: []graphman.TrustedForType{
				"SSL",
				"SIGNING_SERVER_CERTS",
			},
		})
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	// Calculate checksum
	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// If gateway is external (not managed by operator), use specified port
	if isExternalGateway {
		// For external gateways, we can't use the standard reconciliation
		// Log that external gateway support requires additional configuration
		params.Log.V(2).Info("External DMZ gateway detected, port will be used for connection",
			"gateway", gateway.Spec.App.Otk.DmzOtkGatewayReference,
			"port", gateway.Spec.App.Otk.DmzGatewayPort)
		// Note: External gateway connection would require additional implementation
		// For now, we skip reconciliation for external gateways
		return nil
	}

	// Get DMZ gateway secret
	name := dmzGateway.Name
	if dmzGateway.Spec.App.Management.SecretName != "" {
		name = dmzGateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)
	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/" + gateway.Name + "-internal-certificates"

	dmzParams := params
	dmzParams.Instance = dmzGateway

	// Note: DmzGatewayPort is used when the gateway is external (not found in cluster)
	// For operator-managed gateways, the gateway's own graphman port configuration is used

	if !dmzGateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, dmzParams)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, dmzParams, "otk certificates", *podList, dmzGateway, gwSecret, "", annotation, sha1Sum, true, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, dmzParams)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, dmzParams, "otk certificates", *gatewayDeployment, dmzGateway, gwSecret, "", annotation, sha1Sum, false, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateDmzPrivateKeyClusterProperty updates the cluster property otk.dmz.private_key.name
// with the DMZ private key name. This is called after the DMZ key is successfully updated.
func updateDmzPrivateKeyClusterProperty(ctx context.Context, params Params, gateway *securityv1.Gateway, keyName string) error {
	// Only update cluster property for DMZ gateway type
	if gateway.Spec.App.Otk.Type != securityv1.OtkTypeDMZ {
		return nil
	}

	// Get or create the cluster properties ConfigMap
	cmName := gateway.Name + "-cwp-bundle"
	cm, err := getGatewayConfigMap(ctx, params, cmName)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to get cluster properties ConfigMap: %w", err)
		}
		// ConfigMap doesn't exist, create it with the property
		return createDmzPrivateKeyClusterProperty(ctx, params, gateway, keyName, cmName)
	}

	// Parse existing bundle
	bundle := graphman.Bundle{}
	bundleJSON := cm.Data["cwp.json"]
	if bundleJSON == "" {
		// Empty bundle, create new one
		return createDmzPrivateKeyClusterProperty(ctx, params, gateway, keyName, cmName)
	}

	err = json.Unmarshal([]byte(bundleJSON), &bundle)
	if err != nil {
		return fmt.Errorf("failed to parse cluster properties bundle: %w", err)
	}

	// Initialize bundle properties if nil
	if bundle.Properties == nil {
		bundle.Properties = &graphman.BundleProperties{
			Mappings: graphman.BundleMappings{},
		}
	}

	// Check if property already exists and update it, or add new one
	propertyName := "otk.dmz.private_key.name"
	found := false
	for _, cwp := range bundle.ClusterProperties {
		if cwp.Name == propertyName {
			cwp.Value = keyName
			found = true
			break
		}
	}

	if !found {
		// Add new cluster property
		bundle.ClusterProperties = append(bundle.ClusterProperties, &graphman.ClusterPropertyInput{
			Name:  propertyName,
			Value: keyName,
		})
	}

	// Marshal bundle back to JSON
	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster properties bundle: %w", err)
	}

	// Calculate checksum
	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Update ConfigMap
	cm.Data["cwp.json"] = string(bundleBytes)
	if cm.ObjectMeta.Annotations == nil {
		cm.ObjectMeta.Annotations = make(map[string]string)
	}
	cm.ObjectMeta.Annotations["checksum/data"] = sha1Sum

	err = params.Client.Update(ctx, cm)
	if err != nil {
		return fmt.Errorf("failed to update cluster properties ConfigMap: %w", err)
	}

	params.Log.V(2).Info("Updated cluster property ConfigMap", "property", propertyName, "value", keyName, "gateway", gateway.Name)

	// Apply the cluster property using the existing mechanism
	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeClusterProp),
	)
	if err != nil {
		return fmt.Errorf("failed to create gateway update request: %w", err)
	}

	err = SyncGateway(ctx, params, *gwUpdReq)
	if err != nil {
		return fmt.Errorf("failed to sync cluster property: %w", err)
	}

	params.Log.V(2).Info("Applied cluster property", "property", propertyName, "value", keyName, "gateway", gateway.Name)

	return nil
}

// createDmzPrivateKeyClusterProperty creates a new cluster properties ConfigMap with the DMZ private key property
func createDmzPrivateKeyClusterProperty(ctx context.Context, params Params, gateway *securityv1.Gateway, keyName string, cmName string) error {
	// Create new bundle with the property
	bundle := graphman.Bundle{
		ClusterProperties: []*graphman.ClusterPropertyInput{
			{
				Name:  "otk.dmz.private_key.name",
				Value: keyName,
			},
		},
		Properties: &graphman.BundleProperties{
			Mappings: graphman.BundleMappings{},
		},
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster properties bundle: %w", err)
	}

	// Calculate checksum
	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Create ConfigMap
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: gateway.Namespace,
			Annotations: map[string]string{
				"checksum/data": sha1Sum,
			},
		},
		Data: map[string]string{
			"cwp.json": string(bundleBytes),
		},
	}

	// Set controller reference
	if err := controllerutil.SetControllerReference(gateway, cm, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err = params.Client.Create(ctx, cm)
	if err != nil {
		return fmt.Errorf("failed to create cluster properties ConfigMap: %w", err)
	}

	params.Log.V(2).Info("Created cluster property ConfigMap", "property", "otk.dmz.private_key.name", "value", keyName, "gateway", gateway.Name)

	// Apply the cluster property using the existing mechanism
	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeClusterProp),
	)
	if err != nil {
		return fmt.Errorf("failed to create gateway update request: %w", err)
	}

	err = SyncGateway(ctx, params, *gwUpdReq)
	if err != nil {
		return fmt.Errorf("failed to sync cluster property: %w", err)
	}

	params.Log.V(2).Info("Applied cluster property", "property", "otk.dmz.private_key.name", "value", keyName, "gateway", gateway.Name)

	return nil
}
