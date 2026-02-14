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
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
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

	// Publish DMZ certs to Internal Gateway when DMZ key is updated
	if gateway.Spec.App.Otk.Type == securityv1.OtkTypeDMZ && gateway.Spec.App.Otk.DmzKeySecret != "" {
		err = publishDmzCertificatesToInternal(ctx, params, gateway)
		if err != nil {
			params.Log.V(2).Info("failed to publish DMZ certificates to Internal", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
		}
	}

	// Publish Internal certs to DMZ Gateway when Internal key is updated
	if gateway.Spec.App.Otk.Type == securityv1.OtkTypeInternal && gateway.Spec.App.Otk.InternalKeySecret != "" {
		err = publishInternalCertificatesToDmz(ctx, params, gateway)
		if err != nil {
			params.Log.V(2).Info("failed to publish Internal certificates to DMZ", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
		}
	}
}

// publishDmzCertificatesToInternal publishes DMZ certificates to Internal gateway when DMZ key is updated
// Handles ephemeral, DB-backed, and external gateways
func publishDmzCertificatesToInternal(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	// Check if Internal gateway reference is specified
	if gateway.Spec.App.Otk.InternalOtkGatewayReference == "" {
		return nil
	}

	// Get DMZ key secret
	dmzKeySecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.DmzKeySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("DMZ key secret not found, skipping cert publish", "secret", gateway.Spec.App.Otk.DmzKeySecret)
			return nil
		}
		return err
	}

	// Check if key was updated by comparing annotation
	annotation := "security.brcmlabs.com/otk-dmz-key"
	currentSha1Sum := ""

	if !gateway.Spec.App.Management.Database.Enabled {
		// Ephemeral gateway - check pod annotations
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		for _, pod := range podList.Items {
			if val, ok := pod.ObjectMeta.Annotations[annotation]; ok {
				currentSha1Sum = val
				break
			}
		}
	} else {
		// DB-backed gateway - check deployment annotations
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		currentSha1Sum = gatewayDeployment.ObjectMeta.Annotations[annotation]
	}

	// Calculate current key checksum
	certData := dmzKeySecret.Data["tls.crt"]
	keyData := dmzKeySecret.Data["tls.key"]
	if len(certData) == 0 || len(keyData) == 0 {
		return fmt.Errorf("DMZ key secret must contain tls.crt and tls.key")
	}

	keySecretMap := []struct {
		Name      string
		Crt       string
		Key       string
		Alias     string
		UsageType string
	}{
		{
			Name:      "dmz-key",
			Crt:       string(certData),
			Key:       string(keyData),
			Alias:     "otk-dmz-key",
			UsageType: "",
		},
	}

	dataBytes, _ := json.Marshal(&keySecretMap)
	h := sha1.New()
	h.Write(dataBytes)
	newSha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Only publish if key was updated
	if currentSha1Sum == newSha1Sum {
		params.Log.V(2).Info("DMZ key not updated, skipping cert publish", "gateway", gateway.Name)
		return nil
	}

	// Publish DMZ cert to Internal (handles ephemeral, DB-backed, and external gateways)
	return publishDmzCertToInternal(ctx, params, gateway, dmzKeySecret)
}

// publishInternalCertificatesToDmz publishes Internal certificates to DMZ gateway when Internal key is updated
// Handles ephemeral, DB-backed, and external gateways
func publishInternalCertificatesToDmz(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	// Check if DMZ gateway reference is specified
	if gateway.Spec.App.Otk.DmzOtkGatewayReference == "" {
		return nil
	}

	// Get Internal key secret
	internalKeySecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.InternalKeySecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("Internal key secret not found, skipping cert publish", "secret", gateway.Spec.App.Otk.InternalKeySecret)
			return nil
		}
		return err
	}

	// Check if key was updated by comparing annotation
	annotation := "security.brcmlabs.com/otk-internal-key"
	currentSha1Sum := ""

	if !gateway.Spec.App.Management.Database.Enabled {
		// Ephemeral gateway - check pod annotations
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		for _, pod := range podList.Items {
			if val, ok := pod.ObjectMeta.Annotations[annotation]; ok {
				currentSha1Sum = val
				break
			}
		}
	} else {
		// DB-backed gateway - check deployment annotations
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		currentSha1Sum = gatewayDeployment.ObjectMeta.Annotations[annotation]
	}

	// Calculate current key checksum
	certData := internalKeySecret.Data["tls.crt"]
	keyData := internalKeySecret.Data["tls.key"]
	if len(certData) == 0 || len(keyData) == 0 {
		return fmt.Errorf("Internal key secret must contain tls.crt and tls.key")
	}

	keySecretMap := []struct {
		Name      string
		Crt       string
		Key       string
		Alias     string
		UsageType string
	}{
		{
			Name:      "internal-key",
			Crt:       string(certData),
			Key:       string(keyData),
			Alias:     "otk-internal-key",
			UsageType: "",
		},
	}

	dataBytes, _ := json.Marshal(&keySecretMap)
	h := sha1.New()
	h.Write(dataBytes)
	newSha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	// Only publish if key was updated
	if currentSha1Sum == newSha1Sum {
		params.Log.V(2).Info("Internal key not updated, skipping cert publish", "gateway", gateway.Name)
		return nil
	}

	// Publish Internal cert to DMZ (handles ephemeral, DB-backed, and external gateways)
	return publishInternalCertToDmz(ctx, params, gateway, internalKeySecret)
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

	// Before adding new certs, remove existing ones if they were previously applied
	// Check if certs were previously applied by checking the annotation
	annotation := "security.brcmlabs.com/" + gateway.Name + "-dmz-certificates"
	thumbprintAnnotation := "security.brcmlabs.com/" + gateway.Name + "-dmz-certificates-thumbprints"
	previousCertChecksum := ""
	var oldThumbprints []string
	if !isExternalGateway {
		if !internalGateway.Spec.App.Management.Database.Enabled {
			podList, err := getGatewayPods(ctx, params)
			if err == nil {
				for _, pod := range podList.Items {
					if val, ok := pod.ObjectMeta.Annotations[annotation]; ok {
						previousCertChecksum = val
					}
					if val, ok := pod.ObjectMeta.Annotations[thumbprintAnnotation]; ok && val != "" {
						// Parse comma-separated thumbprints
						oldThumbprints = strings.Split(val, ",")
					}
					if previousCertChecksum != "" {
						break
					}
				}
			}
		} else {
			gatewayDeployment, err := getGatewayDeployment(ctx, params)
			if err == nil {
				previousCertChecksum = gatewayDeployment.ObjectMeta.Annotations[annotation]
				if val, ok := gatewayDeployment.ObjectMeta.Annotations[thumbprintAnnotation]; ok && val != "" {
					oldThumbprints = strings.Split(val, ",")
				}
			}
		}
	}

	bundle := graphman.Bundle{}

	// If we have old thumbprints, add deletion mappings before adding new certs
	if len(oldThumbprints) > 0 && previousCertChecksum != "" {
		if bundle.Properties == nil {
			bundle.Properties = &graphman.BundleProperties{
				Mappings: graphman.BundleMappings{},
			}
		}
		for _, thumbprint := range oldThumbprints {
			thumbprint = strings.TrimSpace(thumbprint)
			if thumbprint != "" {
				bundle.Properties.Mappings.TrustedCerts = append(bundle.Properties.Mappings.TrustedCerts, &graphman.MappingInstructionInput{
					Action: graphman.MappingActionDelete,
					Source: graphman.MappingSource{ThumbprintSha1: thumbprint},
				})
			}
		}
		// Also remove old FIP users with the same names
		// We'll identify them by the cert CommonName pattern
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
			// Remove FIP user by name (CommonName)
			bundle.Properties.Mappings.FipUsers = append(bundle.Properties.Mappings.FipUsers, &graphman.MappingInstructionInput{
				Action: graphman.MappingActionDelete,
				Source: graphman.MappingSource{Name: crtX509.Subject.CommonName},
			})
		}
	}

	// Calculate thumbprints for new certs and add to TrustedCerts
	var newThumbprints []string
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

		// Calculate thumbprint for this cert
		thumbprint, err := calculateCertThumbprint(crtX509.Raw)
		if err != nil {
			params.Log.V(2).Info("Failed to calculate cert thumbprint", "error", err, "cert", crtX509.Subject.CommonName)
			thumbprint = "" // Continue without thumbprint
		} else {
			newThumbprints = append(newThumbprints, thumbprint)
		}

		bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
			Name:                      crtX509.Subject.CommonName,
			CertBase64:                base64.StdEncoding.EncodeToString([]byte(certStr)),
			ThumbprintSha1:            thumbprint,
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

	// If gateway is external (not managed by operator), use specified port and auth secret
	if isExternalGateway {
		// Parse certificates to extract information for FIP user creation
		var certInfo []struct {
			commonName string
			subjectDn  string
			certRaw    []byte
		}

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

			// Extract full Subject DN from certificate
			subjectDn := extractSubjectDN(crtX509)

			certInfo = append(certInfo, struct {
				commonName string
				subjectDn  string
				certRaw    []byte
			}{
				commonName: crtX509.Subject.CommonName,
				subjectDn:  subjectDn,
				certRaw:    crtX509.Raw,
			})
		}

		return syncDmzCertToExternalInternalGateway(ctx, params, gateway, dmzKeySecret, certInfo)
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

	// Before adding new certs, remove existing ones if they were previously applied
	// Check if certs were previously applied by checking the annotation
	annotation := "security.brcmlabs.com/" + gateway.Name + "-internal-certificates"
	thumbprintAnnotation := "security.brcmlabs.com/" + gateway.Name + "-internal-certificates-thumbprints"
	previousCertChecksum := ""
	var oldThumbprints []string
	if !isExternalGateway {
		if !dmzGateway.Spec.App.Management.Database.Enabled {
			dmzParams := params
			dmzParams.Instance = dmzGateway
			podList, err := getGatewayPods(ctx, dmzParams)
			if err == nil {
				for _, pod := range podList.Items {
					if val, ok := pod.ObjectMeta.Annotations[annotation]; ok {
						previousCertChecksum = val
					}
					if val, ok := pod.ObjectMeta.Annotations[thumbprintAnnotation]; ok && val != "" {
						// Parse comma-separated thumbprints
						oldThumbprints = strings.Split(val, ",")
					}
					if previousCertChecksum != "" {
						break
					}
				}
			}
		} else {
			dmzParams := params
			dmzParams.Instance = dmzGateway
			gatewayDeployment, err := getGatewayDeployment(ctx, dmzParams)
			if err == nil {
				previousCertChecksum = gatewayDeployment.ObjectMeta.Annotations[annotation]
				if val, ok := gatewayDeployment.ObjectMeta.Annotations[thumbprintAnnotation]; ok && val != "" {
					oldThumbprints = strings.Split(val, ",")
				}
			}
		}
	}

	bundle := graphman.Bundle{}

	// If we have old thumbprints, add deletion mappings before adding new certs
	if len(oldThumbprints) > 0 && previousCertChecksum != "" {
		if bundle.Properties == nil {
			bundle.Properties = &graphman.BundleProperties{
				Mappings: graphman.BundleMappings{},
			}
		}
		for _, thumbprint := range oldThumbprints {
			thumbprint = strings.TrimSpace(thumbprint)
			if thumbprint != "" {
				bundle.Properties.Mappings.TrustedCerts = append(bundle.Properties.Mappings.TrustedCerts, &graphman.MappingInstructionInput{
					Action: graphman.MappingActionDelete,
					Source: graphman.MappingSource{ThumbprintSha1: thumbprint},
				})
			}
		}
	}

	// Calculate thumbprints for new certs and add to TrustedCerts
	var newThumbprints []string
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

		// Calculate thumbprint for this cert
		thumbprint, err := calculateCertThumbprint(crtX509.Raw)
		if err != nil {
			params.Log.V(2).Info("Failed to calculate cert thumbprint", "error", err, "cert", crtX509.Subject.CommonName)
			thumbprint = "" // Continue without thumbprint
		} else {
			newThumbprints = append(newThumbprints, thumbprint)
		}

		bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
			Name:                      crtX509.Subject.CommonName,
			CertBase64:                base64.StdEncoding.EncodeToString([]byte(certStr)),
			ThumbprintSha1:            thumbprint,
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

	// If gateway is external (not managed by operator), use specified port and auth secret
	if isExternalGateway {
		return syncInternalCertToExternalDmzGateway(ctx, params, gateway, bundleBytes, sha1Sum)
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

	// annotation is already declared above, reuse it
	// annotation := "security.brcmlabs.com/" + gateway.Name + "-internal-certificates"

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

// syncDmzCertToExternalInternalGateway syncs DMZ certificate to an external Internal gateway
// using graphman. First it adds the certificate as a trusted cert, then creates a FIP user.
func syncDmzCertToExternalInternalGateway(ctx context.Context, params Params, gateway *securityv1.Gateway, dmzKeySecret *corev1.Secret, certInfo []struct {
	commonName string
	subjectDn  string
	certRaw    []byte
}) error {
	// Get auth secret for external Internal gateway
	if gateway.Spec.App.Otk.InternalAuthSecret == "" {
		return fmt.Errorf("internalAuthSecret is required for external Internal gateway")
	}

	authSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.InternalAuthSecret)
	if err != nil {
		return fmt.Errorf("failed to get auth secret for external Internal gateway: %w", err)
	}

	// Parse username and password from auth secret
	username, password := parseGatewaySecret(authSecret)
	if username == "" || password == "" {
		return fmt.Errorf("could not retrieve gateway credentials from auth secret: %s", gateway.Spec.App.Otk.InternalAuthSecret)
	}

	// Build endpoint URL for external gateway
	// Format: <gateway-reference>:<port>/graphman
	// ApplyGraphmanBundle expects format: host:port/path (without https://)
	gatewayReference := gateway.Spec.App.Otk.InternalOtkGatewayReference
	port := gateway.Spec.App.Otk.InternalGatewayPort
	if port == 0 {
		port = 9443 // Default graphman port
	}

	// For external gateways, the reference might be a hostname or IP
	// If it's just a name without domain, we might need to construct a full hostname
	// For now, use the reference as-is (could be FQDN, hostname, or IP)
	endpoint := fmt.Sprintf("%s:%d/graphman", gatewayReference, port)

	// Step 1: Sync DMZ certificate as TrustedCert first
	certData := dmzKeySecret.Data["tls.crt"]
	if len(certData) == 0 {
		return fmt.Errorf("DMZ key secret must contain tls.crt")
	}

	crtStrings := strings.SplitAfter(string(certData), "-----END CERTIFICATE-----")
	trustedCertBundle := graphman.Bundle{}

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

		trustedCertBundle.TrustedCerts = append(trustedCertBundle.TrustedCerts, &graphman.TrustedCertInput{
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

	trustedCertBundleBytes, err := json.Marshal(trustedCertBundle)
	if err != nil {
		return fmt.Errorf("failed to marshal trusted cert bundle: %w", err)
	}

	params.Log.V(2).Info("Syncing DMZ certificate as TrustedCert to external Internal gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint)

	// Apply trusted cert bundle first
	err = util.ApplyGraphmanBundle(username, password, endpoint, "", trustedCertBundleBytes)
	if err != nil {
		return fmt.Errorf("failed to sync DMZ certificate as TrustedCert to external Internal gateway: %w", err)
	}

	params.Log.Info("Successfully synced DMZ certificate as TrustedCert to external Internal gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint)

	// Step 2: Create FIP user with DMZ certificate
	if len(certInfo) == 0 {
		params.Log.V(2).Info("No certificate info available, skipping FIP user creation")
		return nil
	}

	fipUserBundle := graphman.Bundle{}

	for _, info := range certInfo {
		// Use the extracted Subject DN (not just "cn=" + CommonName)
		// Since FIP identity provider doesn't have a default subject dn, we must provide it
		fipUserBundle.FipUsers = append(fipUserBundle.FipUsers, &graphman.FipUserInput{
			Name:         info.commonName,
			ProviderName: "otk-fips-provider",
			SubjectDn:    info.subjectDn, // Full Subject DN from certificate
			CertBase64:   base64.RawStdEncoding.EncodeToString(info.certRaw),
		})
	}

	fipUserBundleBytes, err := json.Marshal(fipUserBundle)
	if err != nil {
		return fmt.Errorf("failed to marshal FIP user bundle: %w", err)
	}

	params.Log.V(2).Info("Creating FIP user with DMZ certificate in external Internal gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint)

	// Apply FIP user bundle after certificate is synced
	err = util.ApplyGraphmanBundle(username, password, endpoint, "", fipUserBundleBytes)
	if err != nil {
		return fmt.Errorf("failed to create FIP user with DMZ certificate in external Internal gateway: %w", err)
	}

	params.Log.Info("Successfully created FIP user with DMZ certificate in external Internal gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint)

	return nil
}

// extractSubjectDN extracts the full Subject DN from an x509 certificate
// Format: CN=name,OU=org unit,O=org,C=country, etc.
func extractSubjectDN(cert *x509.Certificate) string {
	var parts []string

	// Add CommonName
	if cert.Subject.CommonName != "" {
		parts = append(parts, "CN="+cert.Subject.CommonName)
	}

	// Add Country
	for _, c := range cert.Subject.Country {
		if c != "" {
			parts = append(parts, "C="+c)
		}
	}

	// Add Organization
	for _, o := range cert.Subject.Organization {
		if o != "" {
			parts = append(parts, "O="+o)
		}
	}

	// Add Organizational Unit
	for _, ou := range cert.Subject.OrganizationalUnit {
		if ou != "" {
			parts = append(parts, "OU="+ou)
		}
	}

	// Add Locality
	for _, l := range cert.Subject.Locality {
		if l != "" {
			parts = append(parts, "L="+l)
		}
	}

	// Add Province/State
	for _, p := range cert.Subject.Province {
		if p != "" {
			parts = append(parts, "ST="+p)
		}
	}

	// Add Street Address
	for _, s := range cert.Subject.StreetAddress {
		if s != "" {
			parts = append(parts, "STREET="+s)
		}
	}

	// Add Postal Code
	for _, pc := range cert.Subject.PostalCode {
		if pc != "" {
			parts = append(parts, "POSTALCODE="+pc)
		}
	}

	// Add Serial Number
	if cert.Subject.SerialNumber != "" {
		parts = append(parts, "SERIALNUMBER="+cert.Subject.SerialNumber)
	}

	// Join all parts with comma
	if len(parts) == 0 {
		// Fallback to CN if nothing else is available
		if cert.Subject.CommonName != "" {
			return "CN=" + cert.Subject.CommonName
		}
		return ""
	}

	return strings.Join(parts, ",")
}

// syncInternalCertToExternalDmzGateway syncs Internal certificate to an external DMZ gateway
// using graphman. It adds the certificate as a trusted cert.
func syncInternalCertToExternalDmzGateway(ctx context.Context, params Params, gateway *securityv1.Gateway, bundleBytes []byte, sha1Sum string) error {
	// Get auth secret for external DMZ gateway
	if gateway.Spec.App.Otk.DmzAuthSecret == "" {
		return fmt.Errorf("dmzAuthSecret is required for external DMZ gateway")
	}

	authSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.DmzAuthSecret)
	if err != nil {
		return fmt.Errorf("failed to get auth secret for external DMZ gateway: %w", err)
	}

	// Parse username and password from auth secret
	username, password := parseGatewaySecret(authSecret)
	if username == "" || password == "" {
		return fmt.Errorf("could not retrieve gateway credentials from auth secret: %s", gateway.Spec.App.Otk.DmzAuthSecret)
	}

	// Build endpoint URL for external gateway
	// Format: <gateway-reference>:<port>/graphman
	// ApplyGraphmanBundle expects format: host:port/path (without https://)
	gatewayReference := gateway.Spec.App.Otk.DmzOtkGatewayReference
	port := gateway.Spec.App.Otk.DmzGatewayPort
	if port == 0 {
		port = 9443 // Default graphman port
	}

	// For external gateways, the reference might be a hostname or IP
	endpoint := fmt.Sprintf("%s:%d/graphman", gatewayReference, port)

	params.Log.V(2).Info("Syncing Internal certificate to external DMZ gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint,
		"sha1Sum", sha1Sum)

	// Apply bundle to external gateway using graphman
	err = util.ApplyGraphmanBundle(username, password, endpoint, "", bundleBytes)
	if err != nil {
		return fmt.Errorf("failed to sync Internal certificate to external DMZ gateway: %w", err)
	}

	params.Log.Info("Successfully synced Internal certificate to external DMZ gateway",
		"gateway", gatewayReference,
		"endpoint", endpoint,
		"sha1Sum", sha1Sum)

	return nil
}

// calculateCertThumbprint calculates the SHA1 thumbprint of a certificate in the format expected by Graphman
// Format: base64-encoded hex string of SHA1 fingerprint
func calculateCertThumbprint(rawCert []byte) (string, error) {
	fingerprint := sha1.Sum(rawCert)
	var buf bytes.Buffer
	for _, f := range fingerprint {
		fmt.Fprintf(&buf, "%02X", f)
	}
	hexDump, err := hex.DecodeString(buf.String())
	if err != nil {
		return "", err
	}
	buf.Reset()
	return base64.StdEncoding.EncodeToString(hexDump), nil
}
