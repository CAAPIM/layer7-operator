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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */

package reconcile

import (
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GatewayUpdateRequest struct {
	checksum                     string
	patchAnnotation              string
	delete                       bool
	graphmanPort                 int
	graphmanEncryptionPassphrase string
	bundle                       []byte
	bundleName                   string
	bundleType                   BundleType
	username                     string
	password                     string
	cacheEntry                   string
	stateStore                   bool
	repositoryReference          *securityv1.RepositoryReference
	repository                   *securityv1.Repository
	gateway                      *securityv1.Gateway
	ephemeral                    bool
	otkCerts                     bool
	podList                      *corev1.PodList
	deployment                   *appsv1.Deployment
	externalEntities             []ExternalEntity
}

type ExternalEntity struct {
	Name                 string
	Annotation           string
	Bundle               []byte
	Checksum             string
	CacheEntry           string
	EncryptionPassphrase string
}

type BundleType string

const (
	BundleTypeRepository             BundleType = "repository"
	BundleTypeExternalCert           BundleType = "external certs"
	BundleTypeExternalKey            BundleType = "external keys"
	BundleTypeExternalSecret         BundleType = "external secrets"
	BundleTypeClusterProp            BundleType = "cluster properties"
	BundleTypeListenPort             BundleType = "listen ports"
	BundleTypeOTKDatabaseMaintenance BundleType = "otk db maintenance"
)

type GatewayUpdateRequestOpt func(*GatewayUpdateRequest)

type MappingSource struct {
	Name           string `json:"name,omitempty"`
	Alias          string `json:"alias,omitempty"`
	KeystoreId     string `json:"keystoreId,omitempty"`
	ThumbprintSha1 string `json:"thumbprintSha1,omitempty"`
}

func NewGwUpdateRequest(ctx context.Context, gateway *securityv1.Gateway, params Params, opts ...GatewayUpdateRequestOpt) (*GatewayUpdateRequest, error) {
	graphmanPort := 9443
	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
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
		return nil, err
	}

	username, password := parseGatewaySecret(gwSecret)
	if username == "" || password == "" {
		return nil, fmt.Errorf("could not retrieve gateway credentials for %s", name)
	}

	gwUpdReq := &GatewayUpdateRequest{username: username, password: password, graphmanPort: graphmanPort, gateway: gateway}

	for _, opt := range opts {
		opt(gwUpdReq)
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		gwUpdReq.ephemeral = true
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return nil, err
	}
	gwUpdReq.podList = podList

	if !gwUpdReq.ephemeral {
		deployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return nil, err
		}
		gwUpdReq.deployment = deployment
	}

	switch gwUpdReq.bundleType {
	case BundleTypeRepository:

		if !gwUpdReq.delete {
			if (gwUpdReq.repository.Spec.StateStoreReference == "" && !gwUpdReq.repositoryReference.Enabled) || !gwUpdReq.repository.Spec.Enabled {
				return nil, nil
			}
		}

		gwUpdReq.patchAnnotation = "security.brcmlabs.com/" + gwUpdReq.repositoryReference.Name + "-" + string(gwUpdReq.repositoryReference.Type)
		graphmanEncryptionPassphrase := gwUpdReq.repositoryReference.Encryption.Passphrase

		// check for directory change
		directoryChange := false
		for _, repoStatus := range gwUpdReq.gateway.Status.RepositoryStatus {
			if gwUpdReq.repositoryReference.Name == repoStatus.Name {
				if !reflect.DeepEqual(gwUpdReq.repositoryReference.Directories, repoStatus.Directories) {
					directoryChange = true
				}
			}
		}

		/// if no pods are ready return nil
		if gwUpdReq.ephemeral {
			updCntr := 0
			ready := false
			for _, pod := range gwUpdReq.podList.Items {
				if (pod.ObjectMeta.Annotations[gwUpdReq.patchAnnotation] == gwUpdReq.checksum && pod.ObjectMeta.Labels["management-access"] != "leader") || pod.ObjectMeta.Annotations[gwUpdReq.patchAnnotation] == gwUpdReq.checksum+"-leader" {
					updCntr = updCntr + 1
				}
				for _, ps := range pod.Status.ContainerStatuses {
					if ps.Ready {
						ready = true
					}
				}
			}

			if updCntr == len(gwUpdReq.podList.Items) && !gwUpdReq.delete && !directoryChange {
				return nil, nil
			}

			// If pods aren't ready yet, only proceed if bootstrap is enabled OR if it's a delete
			if !ready && !gwUpdReq.delete {
				// With bootstrap enabled, we can patch pods that just started but aren't ready yet
				if !gwUpdReq.gateway.Spec.App.RepositoryReferenceBootstrap.Enabled {
					return nil, nil
				}
				// If bootstrap enabled and not ready, continue to patch with checksum
			}

		} else {
			if (gwUpdReq.deployment.Annotations[gwUpdReq.patchAnnotation] == gwUpdReq.checksum || gwUpdReq.repositoryReference.Type == securityv1.RepositoryReferenceTypeStatic) && !gwUpdReq.delete && !directoryChange {
				return nil, nil
			}
		}

		if gwUpdReq.repositoryReference.Encryption.ExistingSecret != "" {
			graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, gwUpdReq.repositoryReference.Encryption.ExistingSecret, gwUpdReq.repositoryReference.Encryption.Key)
			if err != nil {
				return nil, err
			}
		}

		if len(gwUpdReq.repositoryReference.Directories) == 0 {
			gwUpdReq.repositoryReference.Directories = []string{"/"}
		}

		if gwUpdReq.repository.Spec.Type == securityv1.RepositoryTypeLocal {
			gwUpdReq.bundle, err = readLocalReference(ctx, gwUpdReq.repository, params)
			if err != nil {
				return nil, err
			}
		} else {
			gwUpdReq.bundle, err = buildBundle(ctx, params, gwUpdReq.repositoryReference, gwUpdReq.repository, gwUpdReq.gateway, gwUpdReq.delete)
			if err != nil {
				return nil, err
			}
		}

		gwUpdReq.graphmanEncryptionPassphrase = graphmanEncryptionPassphrase

		gwUpdReq.cacheEntry = gwUpdReq.repositoryReference.Name + "-" + gwUpdReq.checksum
		gwUpdReq.bundleName = gwUpdReq.repositoryReference.Name
	case BundleTypeClusterProp:
		cm, err := getGatewayConfigMap(ctx, params, gateway.Name+"-cwp-bundle")
		if err != nil {
			return nil, err
		}

		bundle := graphman.Bundle{}
		err = json.Unmarshal([]byte(cm.Data["cwp.json"]), &bundle)
		if err != nil {
			return nil, err
		}

		notFound := []string{}

		if !gwUpdReq.delete {
			for _, sCwp := range gateway.Status.LastAppliedClusterProperties {
				found := false
				for _, cwp := range bundle.ClusterProperties {
					if cwp.Name == sCwp {
						found = true
					}
				}
				if !found {
					notFound = append(notFound, sCwp)
				}
			}
		} else {
			notFound = append(notFound, gateway.Status.LastAppliedClusterProperties...)
		}

		bundle.Properties = &graphman.BundleProperties{}
		for _, deletedCwp := range notFound {
			mappingSource := graphman.MappingSource{Name: deletedCwp}
			bundle.ClusterProperties = append(bundle.ClusterProperties, &graphman.ClusterPropertyInput{
				Name:  deletedCwp,
				Value: "to be deleted",
			})

			bundle.Properties.Mappings.ClusterProperties = append(bundle.Properties.Mappings.ClusterProperties, &graphman.MappingInstructionInput{
				Action: graphman.MappingActionDelete,
				Source: mappingSource,
			})
		}

		bundleBytes, err := json.Marshal(bundle)
		if err != nil {
			return nil, err
		}
		gwUpdReq.graphmanEncryptionPassphrase = ""
		gwUpdReq.bundle = bundleBytes
		gwUpdReq.patchAnnotation = "security.brcmlabs.com/" + gateway.Name + "-cwp-bundle"
		gwUpdReq.checksum = cm.ObjectMeta.Annotations["checksum/data"]
		gwUpdReq.cacheEntry = gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + gwUpdReq.checksum
		gwUpdReq.bundleName = string(gwUpdReq.bundleType)
	case BundleTypeListenPort:
		refreshOnKeyChanges := false
		checksum := ""
		var bundleBytes []byte
		if gateway.Spec.App.ListenPorts.Custom.Enabled {

			if gateway.Spec.App.ListenPorts.RefreshOnKeyChanges {
				refreshOnKeyChanges = true
			}
			if !gateway.Spec.App.ListenPorts.Custom.Enabled {
				bundleBytes, checksum, err = util.BuildDefaultListenPortBundle(refreshOnKeyChanges)
				if err != nil {
					return nil, err
				}
			} else {
				bundleBytes, checksum, err = util.BuildCustomListenPortBundle(gateway, refreshOnKeyChanges)
				if err != nil {
					return nil, err
				}
			}

			bundle := graphman.Bundle{}
			err = json.Unmarshal(bundleBytes, &bundle)
			if err != nil {
				return nil, err
			}

			notFound := []string{}
			if !gwUpdReq.delete {
				for _, slistenPort := range params.Instance.Status.LastAppliedListenPorts {
					found := false
					for _, listenPort := range bundle.ListenPorts {
						if listenPort.Name == slistenPort {
							found = true
						}
						// anti-lockout
						if listenPort.Name == slistenPort && listenPort.Port == gwUpdReq.graphmanPort {
							found = true
						}
					}
					if !found {
						notFound = append(notFound, slistenPort)
					}
				}
			} else {
				notFound = append(notFound, gateway.Status.LastAppliedListenPorts...)
			}

			bundle.Properties = &graphman.BundleProperties{}
			for _, deletedListenPort := range notFound {
				mappingSource := graphman.MappingSource{Name: deletedListenPort}
				bundle.ListenPorts = append(bundle.ListenPorts, &graphman.ListenPortInput{
					Name:     deletedListenPort,
					Port:     1,
					Enabled:  false,
					Protocol: "HTTP",
					EnabledFeatures: []graphman.ListenPortFeature{
						graphman.ListenPortFeaturePublishedServiceMessageInput,
					},
				})

				bundle.Properties.Mappings.ListenPorts = append(bundle.Properties.Mappings.ListenPorts, &graphman.MappingInstructionInput{
					Action: graphman.MappingActionDelete,
					Source: mappingSource,
				})
			}

			gwUpdReq.bundle, err = json.Marshal(bundle)
			if err != nil {
				return nil, err
			}

			gwUpdReq.graphmanEncryptionPassphrase = ""
			gwUpdReq.patchAnnotation = "security.brcmlabs.com/" + params.Instance.Name + "-listen-port-bundle"
			gwUpdReq.checksum = checksum
			gwUpdReq.cacheEntry = gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + gwUpdReq.checksum
			gwUpdReq.bundleName = string(gwUpdReq.bundleType)
		}
	case BundleTypeExternalCert:
		externalCerts := []ExternalEntity{}

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
					return nil, err
				}

				annotation := "security.brcmlabs.com/external-certs-" + k

				externalCerts = append(externalCerts, ExternalEntity{
					Name:       k,
					Annotation: annotation,
					Bundle:     bundleBytes,
					Checksum:   "deleted",
					CacheEntry: gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + k + "-deleted",
				})
			}
		}

		for _, externalCert := range gateway.Spec.App.ExternalCerts {
			var sha1Sum string

			certSecretMap := []util.GraphmanCert{}
			if externalCert.Enabled {
				secret, err := getGatewaySecret(ctx, params, externalCert.Name)
				if err != nil {
					return nil, err
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
				return nil, err
			}

			if sha1Sum == "" {
				sha1Sum = "deleted"
			}

			annotation := "security.brcmlabs.com/external-certs-" + externalCert.Name

			externalCerts = append(externalCerts, ExternalEntity{
				Name:       externalCert.Name,
				Annotation: annotation,
				Bundle:     bundleBytes,
				Checksum:   sha1Sum,
				CacheEntry: gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + externalCert.Name + "-" + sha1Sum,
			})
		}

		gwUpdReq.graphmanEncryptionPassphrase = ""
		gwUpdReq.bundleName = string(gwUpdReq.bundleType)
		gwUpdReq.externalEntities = externalCerts
	case BundleTypeExternalKey:
		externalKeys := []ExternalEntity{}

		notFound := []string{}
		keySecretMap := []util.GraphmanKey{}
		for _, k := range gateway.Status.LastAppliedExternalKeys {
			found := false
			for _, ek := range gateway.Spec.App.ExternalKeys {
				if k == ek.Alias && ek.Enabled && !ek.Otk {
					// Only process non-OTK keys in regular external keys flow
					found = true
				}
			}
			if !found {
				notFound = append(notFound, k)
			}
		}

		var sha1Sum string
		for _, externalKey := range gateway.Spec.App.ExternalKeys {

			if externalKey.Enabled && !externalKey.Otk {
				// Skip keys with otk: true - they are handled separately by OTK reconciliation
				secret, err := getGatewaySecret(ctx, params, externalKey.Name)
				if err != nil {
					return nil, err
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
			return nil, nil
		}

		bundleBytes, err := util.ConvertX509ToGraphmanBundle(keySecretMap, notFound)
		if err != nil {
			return nil, err
		}

		dataBytes, _ := json.Marshal(&keySecretMap)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum = fmt.Sprintf("%x", h.Sum(nil))

		annotation := "security.brcmlabs.com/external-keys"

		externalKeys = append(externalKeys, ExternalEntity{
			Name:       "",
			Annotation: annotation,
			Bundle:     bundleBytes,
			Checksum:   sha1Sum,
			CacheEntry: gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + sha1Sum,
		})
		gwUpdReq.graphmanEncryptionPassphrase = ""
		gwUpdReq.bundleName = string(gwUpdReq.bundleType)
		gwUpdReq.externalEntities = externalKeys
	case BundleTypeExternalSecret:
		externalSecrets := []ExternalEntity{}

		for k, v := range gateway.Status.LastAppliedExternalSecrets {
			found := false
			notFound := []string{}

			for _, es := range gateway.Spec.App.ExternalSecrets {
				if k == es.Name {
					found = true
				}
			}
			if !found {
				notFound = append(notFound, v...)
				bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(nil, notFound)
				if err != nil {
					return nil, err
				}

				annotation := "security.brcmlabs.com/external-secret-" + k

				externalSecrets = append(externalSecrets, ExternalEntity{
					Name:                 k,
					Annotation:           annotation,
					Bundle:               bundleBytes,
					Checksum:             "deleted",
					CacheEntry:           gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + k + "-deleted",
					EncryptionPassphrase: "",
				})
			}

		}

		for _, es := range gateway.Spec.App.ExternalSecrets {
			var sha1Sum string
			opaqueSecretMap := []util.GraphmanSecret{}
			if es.Enabled {
				secret, err := getGatewaySecret(ctx, params, es.Name)
				if err != nil {
					return nil, err
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
				case corev1.SecretTypeServiceAccountToken, corev1.SecretTypeBasicAuth:
					for k, v := range secret.Data {
						opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
							Name:                 es.Name + "-" + k,
							Secret:               string(v),
							Description:          es.Description,
							VariableReferencable: es.VariableReferencable,
						})
					}
				case corev1.SecretTypeDockercfg, corev1.SecretTypeDockerConfigJson:
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

				dataBytes, _ := json.Marshal(&secret.Data)
				h := sha1.New()
				h.Write(dataBytes)
				sha1Sum = fmt.Sprintf("%x", h.Sum(nil))
			}

			notFound := []string{}
			if gateway.Status.LastAppliedExternalSecrets != nil && gateway.Status.LastAppliedExternalSecrets[es.Name] != nil {

				for _, appliedSecret := range gateway.Status.LastAppliedExternalSecrets[es.Name] {
					found := false
					for _, desiredSecret := range opaqueSecretMap {
						if appliedSecret == desiredSecret.Name {
							found = true
						}
					}
					if !found {
						notFound = append(notFound, appliedSecret)
					}
				}
			}

			if len(opaqueSecretMap) < 1 && len(notFound) < 1 {
				continue
			}

			bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(opaqueSecretMap, notFound)
			if err != nil {
				return nil, err
			}

			graphmanEncryptionPassphrase := es.Encryption.Passphrase
			if es.Encryption.ExistingSecret != "" {
				graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, es.Encryption.ExistingSecret, es.Encryption.Key)
				if err != nil {
					return nil, err
				}
			}

			if sha1Sum == "" {
				sha1Sum = "deleted"
			}

			annotation := "security.brcmlabs.com/external-secret-" + es.Name

			externalSecrets = append(externalSecrets, ExternalEntity{
				Name:                 es.Name,
				Annotation:           annotation,
				Bundle:               bundleBytes,
				Checksum:             sha1Sum,
				CacheEntry:           gateway.Name + "-" + string(gwUpdReq.bundleType) + "-" + es.Name + "-" + sha1Sum,
				EncryptionPassphrase: graphmanEncryptionPassphrase,
			})
		}

		gwUpdReq.graphmanEncryptionPassphrase = ""
		gwUpdReq.bundleName = string(gwUpdReq.bundleType)
		gwUpdReq.externalEntities = externalSecrets
	}

	return gwUpdReq, nil
}

func WithBundleType(bundleType BundleType) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.bundleType = bundleType
	}
}

func WithRepositoryReference(repositoryReference securityv1.RepositoryReference) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.repositoryReference = &repositoryReference
	}
}

func WithRepository(repository *securityv1.Repository) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.repository = repository
	}
}

func WithBundle(bundle []byte) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.bundle = bundle
	}
}

func WithBundleName(bundleName string) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.bundleName = bundleName
	}
}

func WithDelete(delete bool) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.delete = delete
	}
}

func WithChecksum(checksum string) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.checksum = checksum
	}
}

func WithPatchAnnotation(patchAnnotation string) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.patchAnnotation = patchAnnotation
	}
}

func WithCacheEntry(cacheEntry string) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.cacheEntry = cacheEntry
	}
}

func WithOTKCerts(otkCerts bool) GatewayUpdateRequestOpt {
	return func(gUpdReq *GatewayUpdateRequest) {
		gUpdReq.otkCerts = otkCerts
	}
}

func SyncGateway(ctx context.Context, params Params, gwUpdReq GatewayUpdateRequest) (err error) {
	switch gwUpdReq.ephemeral {
	case true:
		err = updateGatewayPods(ctx, params, &gwUpdReq)
		if err != nil {
			return err
		}
	default:
		err = updateGatewayDeployment(ctx, params, &gwUpdReq)
		if err != nil {
			return err
		}
	}
	return nil
}

// buildDeleteBundle handles repository deletion by creating a bundle with delete mappings
func buildDeleteBundle(repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, gateway *securityv1.Gateway, params Params) ([]byte, error) {
	// Check if delete is enabled
	if !gateway.Spec.App.RepositoryReferenceDelete.Enabled {
		params.Log.V(2).Info("repository delete skipped - RepositoryReferenceDelete.Enabled is false", "repository", repoRef.Name)
		return []byte("{}"), nil
	}

	// Check if we should skip delete for non-statestore repos
	if repository.Spec.StateStoreReference == "" && !gateway.Spec.App.RepositoryReferenceDelete.IncludeEfs {
		params.Log.V(2).Info("repository delete skipped - non-statestore repo and IncludeEfs is false", "repository", repoRef.Name)
		return []byte("{}"), nil
	}

	// Determine cache location
	cachePath, cacheFileName := determineCacheLocation(repository, gateway)

	// Build bundle from cache
	bundleBytes, err := buildBundleFromCache(repository, repoRef, cachePath, cacheFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to build delete bundle from cache: %w", err)
	}

	// Set default action to delete
	var bundle graphman.Bundle
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle for delete: %w", err)
	}

	bundle.Properties = &graphman.BundleProperties{DefaultAction: graphman.MappingActionDelete}

	bundleBytes, err = json.Marshal(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bundle with delete mapping: %w", err)
	}

	params.Log.V(2).Info("applied delete mapping to bundle", "repository", repoRef.Name)
	return bundleBytes, nil
}

// checkRetryScenario checks if we should retry the last applied bundle due to previous failure
func checkRetryScenario(gateway *securityv1.Gateway, repoRefName string, currentCommit string, tmpPath string, params Params) (shouldRetry bool, bundle []byte) {
	// Check gateway status for apply failures for this repository
	for _, repoStatus := range gateway.Status.RepositoryStatus {
		if repoStatus.Name == repoRefName {
			// Check if there's a failure condition with "failed" or "error" in reason
			for _, condition := range repoStatus.Conditions {
				// Look for failures in the reason field
				reasonLower := strings.ToLower(condition.Reason)
				if (strings.Contains(reasonLower, "failed") || strings.Contains(reasonLower, "error")) && condition.Status != "success" {
					// Found a failure - check if commit has changed
					if repoStatus.Commit == currentCommit {
						// Commit unchanged, this is a retry scenario
						lastAppliedFile := tmpPath + "/last_applied_" + repoRefName + ".json"
						if bundleBytes, err := os.ReadFile(lastAppliedFile); err == nil {
							params.Log.V(2).Info("retry scenario detected - using last applied bundle",
								"repository", repoRefName,
								"commit", currentCommit,
								"failureReason", condition.Reason)
							return true, bundleBytes
						} else {
							params.Log.V(2).Info("retry scenario but last applied bundle not found, will rebuild",
								"repository", repoRefName,
								"error", err)
							return true, nil
						}
					}
				}
			}
		}
	}
	return false, nil
}

// determineCacheLocation returns the cache path and filename based on repository and gateway configuration
func determineCacheLocation(repository *securityv1.Repository, gateway *securityv1.Gateway) (cachePath string, cacheFileName string) {
	// Local and HTTP repos always use vanilla bundles (no operator-generated delete mappings)
	if repository.Spec.Type == securityv1.RepositoryTypeLocal || repository.Spec.Type == securityv1.RepositoryTypeHttp {
		cachePath = "/tmp/repo-cache/" + repository.Name
		cacheFileName = repository.Status.Commit + ".json"
		return cachePath, cacheFileName
	}

	if repository.Spec.StateStoreReference != "" {
		// Statestore-backed repository
		cachePath = "/tmp/statestore/" + repository.Name
		cacheFileName = "latest.json"
	} else if gateway.Spec.App.RepositoryReferenceDelete.Enabled && gateway.Spec.App.RepositoryReferenceDelete.IncludeEfs {
		// Non-statestore with delete enabled
		cachePath = "/tmp/repo-cache/" + repository.Name
		cacheFileName = "combined.json"
	} else {
		// Non-statestore with delete disabled (vanilla)
		cachePath = "/tmp/repo-cache/" + repository.Name
		cacheFileName = repository.Status.Commit + ".json"
	}
	return cachePath, cacheFileName
}

// shouldSkipDeltaComparison determines if we should skip delta comparison and build a vanilla bundle
func shouldSkipDeltaComparison(gateway *securityv1.Gateway, repository *securityv1.Repository) bool {
	// Local and HTTP repos are ALWAYS vanilla (no delta comparison)
	if repository.Spec.Type == securityv1.RepositoryTypeLocal || repository.Spec.Type == securityv1.RepositoryTypeHttp {
		return true
	}

	// If master flag is disabled, always skip delta comparison
	if !gateway.Spec.App.RepositoryReferenceDelete.Enabled {
		return true
	}

	// If non-statestore and includeEfs is false, skip delta comparison
	if repository.Spec.StateStoreReference == "" && !gateway.Spec.App.RepositoryReferenceDelete.IncludeEfs {
		return true
	}

	return false
}

// buildVanillaBundleAndCache builds a vanilla bundle from cache and stores it
func buildVanillaBundleAndCache(repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, gateway *securityv1.Gateway, cachePath string, cacheFileName string, tmpPath string, fileName string, params Params) ([]byte, error) {
	// Build bundle from cache
	bundleBytes, err := buildBundleFromCache(repository, repoRef, cachePath, cacheFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to build vanilla bundle: %w", err)
	}

	// Write to cache and last_applied
	if err := writeBundlesToDisk(repository, repoRef, gateway, bundleBytes, nil, tmpPath, fileName, cachePath, params); err != nil {
		return nil, err
	}

	return bundleBytes, nil
}

// handleDirectoryChange handles directory changes with delta calculation
func handleDirectoryChange(ctx context.Context, params Params, repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, gateway *securityv1.Gateway, cachePath string, cacheFileName string, tmpPath string, fileName string, previousDirectories []string) ([]byte, error) {
	// Step 1: Build new bundle from current directories
	newBundleBytes, err := buildBundleFromCache(repository, repoRef, cachePath, cacheFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to build new bundle for directory change: %w", err)
	}

	// For combined.json and latest.json (StateStore) with directory additions/reordering (not removals), repository mappings are correct
	// But if directories were removed, we need to calculate delta to generate DELETE mappings
	// Check if any previous directories are missing from current directories
	directoriesRemoved := false
	for _, prevDir := range previousDirectories {
		found := false
		for _, currDir := range repoRef.Directories {
			if prevDir == currDir {
				found = true
				break
			}
		}
		if !found {
			directoriesRemoved = true
			params.Log.Info("directory removed, calculating delta",
				"repository", repoRef.Name,
				"removedDirectory", prevDir)
			break
		}
	}

	// If no directories were removed, just use repository bundle as-is for repository-controlled mappings
	if (cacheFileName == "combined.json" || cacheFileName == "latest.json") && !directoriesRemoved {
		sourceType := cacheFileName
		params.Log.V(2).Info("directory change with repository-controlled mappings - using repository mappings",
			"repository", repoRef.Name,
			"sourceType", sourceType,
			"previousDirs", previousDirectories,
			"currentDirs", repoRef.Directories)

		if err := writeBundlesToDisk(repository, repoRef, gateway, newBundleBytes, nil, tmpPath, fileName, cachePath, params); err != nil {
			return nil, err
		}
		return newBundleBytes, nil
	}

	// Step 2: Look up previous bundle
	var previousBundleBytes []byte

	if len(previousDirectories) == 0 {
		// No previous directories tracked - treat as first deployment
		params.Log.Info("no previous directories tracked, treating as first deployment",
			"repository", repoRef.Name)

		// Just write and return the new bundle
		if err := writeBundlesToDisk(repository, repoRef, gateway, newBundleBytes, nil, tmpPath, fileName, cachePath, params); err != nil {
			return nil, err
		}
		return newBundleBytes, nil
	}

	// Known previous directories - try to read the bundle
	previousFileName := calculateBundleFileName(params.Instance, repoRef.Name, previousDirectories)

	// Try cachePath first (persistent for StateStore), then tmpPath (ephemeral)
	previousBundleBytes, err = os.ReadFile(cachePath + "/" + previousFileName)
	if err != nil {
		// Try tmpPath as fallback
		previousBundleBytes, err = os.ReadFile(tmpPath + "/" + previousFileName)
		if err != nil {
			// Previous bundle file not found (e.g. after operator restart)
			// But we know what directories were applied from status
			// Reconstruct the previous bundle from cache
			params.Log.Info("previous bundle file not found, reconstructing from cache",
				"repository", repoRef.Name,
				"previousDirectories", previousDirectories,
				"previousFileName", previousFileName)

			// Create a temporary repoRef with previous directories to reconstruct
			tempRepoRef := &securityv1.RepositoryReference{
				Name:        repoRef.Name,
				Enabled:     repoRef.Enabled,
				Type:        repoRef.Type,
				Directories: previousDirectories,
			}

			previousBundleBytes, err = buildBundleFromCache(repository, tempRepoRef, cachePath, cacheFileName)
			if err != nil {
				return nil, fmt.Errorf("failed to reconstruct previous bundle from cache: %w", err)
			}

			params.Log.V(2).Info("successfully reconstructed previous bundle from cache",
				"repository", repoRef.Name,
				"previousDirectories", previousDirectories)
		} else {
			params.Log.V(2).Info("found previous bundle in tmpPath",
				"repository", repoRef.Name,
				"previousFileName", previousFileName)
		}
	} else {
		params.Log.V(2).Info("found previous bundle in cachePath",
			"repository", repoRef.Name,
			"previousFileName", previousFileName)
	}

	// Step 3: Parse bundles
	var previousBundle, newBundle graphman.Bundle
	if err := json.Unmarshal(previousBundleBytes, &previousBundle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal previous bundle: %w", err)
	}
	if err := json.Unmarshal(newBundleBytes, &newBundle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal new bundle: %w", err)
	}

	// Step 4: For repository-controlled mappings (combined.json/latest.json), preserve existing DELETE mappings
	// The repository controller has already determined what should be deleted based on commit changes
	// We only need to ADD delete mappings for entities in removed directories
	preserveRepoMappings := (cacheFileName == "combined.json" || cacheFileName == "latest.json")

	// Save the repository-generated DELETE mappings before any processing
	var repoDeleteMappings graphman.BundleMappings
	if preserveRepoMappings && newBundle.Properties != nil {
		repoDeleteMappings = newBundle.Properties.Mappings
	}

	// Step 5: Reset mappings on previous bundle
	if err := graphman.ResetMappings(&previousBundle); err != nil {
		params.Log.V(2).Info("failed to reset mappings on previous bundle, continuing anyway", "error", err)
	}

	// Step 6: Calculate delta for removed directories
	_, combinedBundle, err := graphman.CalculateDelta(previousBundle, newBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate delta for directory change: %w", err)
	}

	// Step 7: For repository-controlled mappings, merge back the repository DELETE mappings
	// This ensures DELETE mappings from commit changes are preserved alongside directory-change DELETEs
	if preserveRepoMappings {
		if combinedBundle.Properties == nil {
			combinedBundle.Properties = &graphman.BundleProperties{}
		}
		// Merge repository DELETE mappings with calculated delta DELETE mappings
		// We need to preserve DELETE mappings from the repository (commit-based deletes)
		// while adding new DELETE mappings from directory changes
		combinedBundle.Properties.Mappings.Services = append(combinedBundle.Properties.Mappings.Services, repoDeleteMappings.Services...)
		combinedBundle.Properties.Mappings.WebApiServices = append(combinedBundle.Properties.Mappings.WebApiServices, repoDeleteMappings.WebApiServices...)
		combinedBundle.Properties.Mappings.InternalWebApiServices = append(combinedBundle.Properties.Mappings.InternalWebApiServices, repoDeleteMappings.InternalWebApiServices...)
		combinedBundle.Properties.Mappings.SoapServices = append(combinedBundle.Properties.Mappings.SoapServices, repoDeleteMappings.SoapServices...)
		combinedBundle.Properties.Mappings.InternalSoapServices = append(combinedBundle.Properties.Mappings.InternalSoapServices, repoDeleteMappings.InternalSoapServices...)
		combinedBundle.Properties.Mappings.Policies = append(combinedBundle.Properties.Mappings.Policies, repoDeleteMappings.Policies...)
		combinedBundle.Properties.Mappings.PolicyFragments = append(combinedBundle.Properties.Mappings.PolicyFragments, repoDeleteMappings.PolicyFragments...)
		combinedBundle.Properties.Mappings.EncassConfigs = append(combinedBundle.Properties.Mappings.EncassConfigs, repoDeleteMappings.EncassConfigs...)
		combinedBundle.Properties.Mappings.HttpConfigurations = append(combinedBundle.Properties.Mappings.HttpConfigurations, repoDeleteMappings.HttpConfigurations...)
		combinedBundle.Properties.Mappings.Keys = append(combinedBundle.Properties.Mappings.Keys, repoDeleteMappings.Keys...)
		combinedBundle.Properties.Mappings.TrustedCerts = append(combinedBundle.Properties.Mappings.TrustedCerts, repoDeleteMappings.TrustedCerts...)
		combinedBundle.Properties.Mappings.Schemas = append(combinedBundle.Properties.Mappings.Schemas, repoDeleteMappings.Schemas...)
		combinedBundle.Properties.Mappings.Dtds = append(combinedBundle.Properties.Mappings.Dtds, repoDeleteMappings.Dtds...)
		combinedBundle.Properties.Mappings.CustomKeyValues = append(combinedBundle.Properties.Mappings.CustomKeyValues, repoDeleteMappings.CustomKeyValues...)
		combinedBundle.Properties.Mappings.ClusterProperties = append(combinedBundle.Properties.Mappings.ClusterProperties, repoDeleteMappings.ClusterProperties...)
		combinedBundle.Properties.Mappings.JdbcConnections = append(combinedBundle.Properties.Mappings.JdbcConnections, repoDeleteMappings.JdbcConnections...)
		combinedBundle.Properties.Mappings.CassandraConnections = append(combinedBundle.Properties.Mappings.CassandraConnections, repoDeleteMappings.CassandraConnections...)
		combinedBundle.Properties.Mappings.JmsDestinations = append(combinedBundle.Properties.Mappings.JmsDestinations, repoDeleteMappings.JmsDestinations...)
		combinedBundle.Properties.Mappings.Secrets = append(combinedBundle.Properties.Mappings.Secrets, repoDeleteMappings.Secrets...)
		combinedBundle.Properties.Mappings.Fips = append(combinedBundle.Properties.Mappings.Fips, repoDeleteMappings.Fips...)
		combinedBundle.Properties.Mappings.Ldaps = append(combinedBundle.Properties.Mappings.Ldaps, repoDeleteMappings.Ldaps...)
		combinedBundle.Properties.Mappings.InternalGroups = append(combinedBundle.Properties.Mappings.InternalGroups, repoDeleteMappings.InternalGroups...)
		combinedBundle.Properties.Mappings.FipGroups = append(combinedBundle.Properties.Mappings.FipGroups, repoDeleteMappings.FipGroups...)
		combinedBundle.Properties.Mappings.InternalUsers = append(combinedBundle.Properties.Mappings.InternalUsers, repoDeleteMappings.InternalUsers...)
		combinedBundle.Properties.Mappings.FipUsers = append(combinedBundle.Properties.Mappings.FipUsers, repoDeleteMappings.FipUsers...)
		combinedBundle.Properties.Mappings.GlobalPolicies = append(combinedBundle.Properties.Mappings.GlobalPolicies, repoDeleteMappings.GlobalPolicies...)
		combinedBundle.Properties.Mappings.BackgroundTasks = append(combinedBundle.Properties.Mappings.BackgroundTasks, repoDeleteMappings.BackgroundTasks...)
		combinedBundle.Properties.Mappings.ScheduledTasks = append(combinedBundle.Properties.Mappings.ScheduledTasks, repoDeleteMappings.ScheduledTasks...)
		combinedBundle.Properties.Mappings.ServerModuleFiles = append(combinedBundle.Properties.Mappings.ServerModuleFiles, repoDeleteMappings.ServerModuleFiles...)
		combinedBundle.Properties.Mappings.SiteMinderConfigs = append(combinedBundle.Properties.Mappings.SiteMinderConfigs, repoDeleteMappings.SiteMinderConfigs...)
		combinedBundle.Properties.Mappings.ActiveConnectors = append(combinedBundle.Properties.Mappings.ActiveConnectors, repoDeleteMappings.ActiveConnectors...)
		combinedBundle.Properties.Mappings.EmailListeners = append(combinedBundle.Properties.Mappings.EmailListeners, repoDeleteMappings.EmailListeners...)
		combinedBundle.Properties.Mappings.ListenPorts = append(combinedBundle.Properties.Mappings.ListenPorts, repoDeleteMappings.ListenPorts...)
		combinedBundle.Properties.Mappings.AdministrativeUserAccountProperties = append(combinedBundle.Properties.Mappings.AdministrativeUserAccountProperties, repoDeleteMappings.AdministrativeUserAccountProperties...)
		combinedBundle.Properties.Mappings.PasswordPolicies = append(combinedBundle.Properties.Mappings.PasswordPolicies, repoDeleteMappings.PasswordPolicies...)
		combinedBundle.Properties.Mappings.RevocationCheckPolicies = append(combinedBundle.Properties.Mappings.RevocationCheckPolicies, repoDeleteMappings.RevocationCheckPolicies...)
		combinedBundle.Properties.Mappings.LogSinks = append(combinedBundle.Properties.Mappings.LogSinks, repoDeleteMappings.LogSinks...)
		combinedBundle.Properties.Mappings.ServiceResolutionConfigs = append(combinedBundle.Properties.Mappings.ServiceResolutionConfigs, repoDeleteMappings.ServiceResolutionConfigs...)
		combinedBundle.Properties.Mappings.Folders = append(combinedBundle.Properties.Mappings.Folders, repoDeleteMappings.Folders...)
		combinedBundle.Properties.Mappings.FederatedIdps = append(combinedBundle.Properties.Mappings.FederatedIdps, repoDeleteMappings.FederatedIdps...)
		combinedBundle.Properties.Mappings.FederatedGroups = append(combinedBundle.Properties.Mappings.FederatedGroups, repoDeleteMappings.FederatedGroups...)
		combinedBundle.Properties.Mappings.FederatedUsers = append(combinedBundle.Properties.Mappings.FederatedUsers, repoDeleteMappings.FederatedUsers...)
		combinedBundle.Properties.Mappings.InternalIdps = append(combinedBundle.Properties.Mappings.InternalIdps, repoDeleteMappings.InternalIdps...)
		combinedBundle.Properties.Mappings.LdapIdps = append(combinedBundle.Properties.Mappings.LdapIdps, repoDeleteMappings.LdapIdps...)
		combinedBundle.Properties.Mappings.SimpleLdapIdps = append(combinedBundle.Properties.Mappings.SimpleLdapIdps, repoDeleteMappings.SimpleLdapIdps...)
		combinedBundle.Properties.Mappings.PolicyBackedIdps = append(combinedBundle.Properties.Mappings.PolicyBackedIdps, repoDeleteMappings.PolicyBackedIdps...)
		combinedBundle.Properties.Mappings.Roles = append(combinedBundle.Properties.Mappings.Roles, repoDeleteMappings.Roles...)
		combinedBundle.Properties.Mappings.GenericEntities = append(combinedBundle.Properties.Mappings.GenericEntities, repoDeleteMappings.GenericEntities...)
		combinedBundle.Properties.Mappings.AuditConfigurations = append(combinedBundle.Properties.Mappings.AuditConfigurations, repoDeleteMappings.AuditConfigurations...)

		params.Log.V(2).Info("merged repository DELETE mappings with directory delta",
			"repository", repoRef.Name,
			"sourceType", cacheFileName,
			"repoServiceMappings", len(repoDeleteMappings.Services),
			"deltaServiceMappings", len(combinedBundle.Properties.Mappings.Services)-len(repoDeleteMappings.Services))
	}

	// Step 8: Marshal the combined bundle (with delete mappings)
	bundleWithMappings, err := json.Marshal(combinedBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal combined bundle: %w", err)
	}

	// Step 9: Create clean version (reset mappings)
	cleanBundle := combinedBundle
	if err := graphman.ResetMappings(&cleanBundle); err != nil {
		params.Log.V(2).Info("failed to reset mappings for clean bundle", "error", err)
	}
	cleanBundleBytes, err := json.Marshal(cleanBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal clean bundle: %w", err)
	}

	// Step 10: Write bundles to disk
	if err := writeBundlesToDisk(repository, repoRef, gateway, bundleWithMappings, cleanBundleBytes, tmpPath, fileName, cachePath, params); err != nil {
		return nil, err
	}

	params.Log.V(2).Info("directory change handled with delta calculation",
		"repository", repoRef.Name,
		"previousDirs", previousDirectories,
		"currentDirs", repoRef.Directories)

	return bundleWithMappings, nil
}

// handleCommitChange handles commit changes with optional delta calculation
func handleCommitChange(ctx context.Context, params Params, repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, gateway *securityv1.Gateway, cachePath string, cacheFileName string, tmpPath string, fileName string) ([]byte, error) {
	// Step 1: Build new bundle from latest commit
	newBundleBytes, err := buildBundleFromCache(repository, repoRef, cachePath, cacheFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to build bundle for new commit: %w", err)
	}

	// Step 2: For combined.json and statestore, the repository controller already calculated deltas
	// Just apply the bundle with its mappings as-is
	if cacheFileName == "combined.json" || repository.Spec.StateStoreReference != "" {
		sourceType := "combined.json"
		if repository.Spec.StateStoreReference != "" {
			sourceType = "statestore"
		}

		params.Log.V(2).Info("commit change - using repository-calculated mappings",
			"repository", repoRef.Name,
			"commit", repository.Status.Commit,
			"sourceType", sourceType)

		// Log bundle details to verify mappings are present
		var newBundle graphman.Bundle
		if err := json.Unmarshal(newBundleBytes, &newBundle); err == nil {
			serviceMappingsCount := 0
			if newBundle.Properties != nil {
				serviceMappingsCount = len(newBundle.Properties.Mappings.Services)
			}
			params.Log.Info("bundle from repository",
				"repository", repoRef.Name,
				"services", len(newBundle.Services),
				"serviceMappings", serviceMappingsCount)
		}

		// Write and return the bundle (repository controller already added delete mappings)
		if err := writeBundlesToDisk(repository, repoRef, gateway, newBundleBytes, nil, tmpPath, fileName, cachePath, params); err != nil {
			return nil, err
		}
		return newBundleBytes, nil
	}

	// Step 3: For {commit}.json (vanilla bundles), just write and return
	// User controls all mappings explicitly - no operator-generated deletes
	params.Log.V(2).Info("commit change with vanilla bundle - no delta calculation",
		"repository", repoRef.Name,
		"commit", repository.Status.Commit)

	// Write and return the bundle as-is (user-defined mappings only)
	if err := writeBundlesToDisk(repository, repoRef, gateway, newBundleBytes, nil, tmpPath, fileName, cachePath, params); err != nil {
		return nil, err
	}
	return newBundleBytes, nil
}

// writeBundlesToDisk writes bundle versions to disk for caching and retry
func writeBundlesToDisk(repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, gateway *securityv1.Gateway, bundleWithMappings []byte, cleanBundle []byte, tmpPath string, fileName string, cachePath string, params Params) error {
	// Clean up old bundles
	cleanupOldBundles(tmpPath)

	// Write commit marker
	commitMarkerPath := tmpPath + "/" + repository.Status.Commit + ".txt"
	if err := os.WriteFile(commitMarkerPath, []byte{}, 0755); err != nil {
		return fmt.Errorf("failed to write commit marker: %w", err)
	}

	// Write clean bundle (if provided)
	if cleanBundle != nil {
		// Write to tmpPath for immediate access
		cleanBundlePath := tmpPath + "/" + fileName
		if err := os.WriteFile(cleanBundlePath, cleanBundle, 0755); err != nil {
			return fmt.Errorf("failed to write clean bundle to tmp: %w", err)
		}
		params.Log.V(5).Info("wrote clean bundle to tmp", "path", cleanBundlePath)

		// Also write to cachePath for persistence (for directory change comparisons)
		cleanBundleCachePath := cachePath + "/" + fileName
		if err := os.WriteFile(cleanBundleCachePath, cleanBundle, 0755); err != nil {
			params.Log.V(2).Info("failed to write clean bundle to cache, continuing", "error", err, "path", cleanBundleCachePath)
		} else {
			params.Log.V(5).Info("wrote clean bundle to cache", "path", cleanBundleCachePath)
		}
	} else {
		// If no separate clean bundle, write the main bundle as clean
		cleanBundlePath := tmpPath + "/" + fileName
		if err := os.WriteFile(cleanBundlePath, bundleWithMappings, 0755); err != nil {
			return fmt.Errorf("failed to write bundle to tmp: %w", err)
		}
		params.Log.V(5).Info("wrote bundle to tmp", "path", cleanBundlePath)

		// Also write to cachePath for persistence
		cleanBundleCachePath := cachePath + "/" + fileName
		if err := os.WriteFile(cleanBundleCachePath, bundleWithMappings, 0755); err != nil {
			params.Log.V(2).Info("failed to write bundle to cache, continuing", "error", err, "path", cleanBundleCachePath)
		} else {
			params.Log.V(5).Info("wrote bundle to cache", "path", cleanBundleCachePath)
		}
	}

	// Write bundle with mappings (if different from clean)
	if cleanBundle != nil {
		bundleWithMappingsPath := tmpPath + "/" + fileName + "_with_mappings"
		if err := os.WriteFile(bundleWithMappingsPath, bundleWithMappings, 0755); err != nil {
			return fmt.Errorf("failed to write bundle with mappings: %w", err)
		}
		params.Log.V(5).Info("wrote bundle with mappings", "path", bundleWithMappingsPath)
	}

	// Write last applied bundle for retry mechanism (to tmpPath for immediate access)
	lastAppliedTmpPath := tmpPath + "/last_applied_" + repoRef.Name + ".json"
	if err := os.WriteFile(lastAppliedTmpPath, bundleWithMappings, 0755); err != nil {
		return fmt.Errorf("failed to write last applied bundle to tmp: %w", err)
	}
	params.Log.V(5).Info("wrote last applied bundle to tmp", "path", lastAppliedTmpPath)

	// Also write to cachePath for persistence across pod restarts
	// Use repository name + reference name for consistency across directory changes
	lastAppliedCachePath := cachePath + "/last_applied_" + repoRef.Name + ".json"
	if err := os.WriteFile(lastAppliedCachePath, bundleWithMappings, 0755); err != nil {
		return fmt.Errorf("failed to write last applied bundle to cache: %w", err)
	}
	params.Log.V(5).Info("wrote last applied bundle to cache", "path", lastAppliedCachePath)

	return nil
}

func buildBundle(ctx context.Context, params Params, repoRef *securityv1.RepositoryReference, repository *securityv1.Repository, gateway *securityv1.Gateway, delete bool) (bundleBytes []byte, err error) {
	tmpPath := "/tmp/bundles/" + repository.Name
	fileName := calculateBundleFileName(params.Instance, repoRef.Name, repoRef.Directories)

	// Ensure temp directory exists
	if err := os.MkdirAll(tmpPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Step 1: Handle delete operations
	if delete {
		return buildDeleteBundle(repository, repoRef, gateway, params)
	}

	// Step 2: Check for retry scenario
	if shouldRetry, retryBundle := checkRetryScenario(gateway, repoRef.Name, repository.Status.Commit, tmpPath, params); shouldRetry {
		if retryBundle != nil {
			params.Log.V(2).Info("retrying last applied bundle due to previous failure", "repository", repoRef.Name)
			return retryBundle, nil
		}
	}

	// Step 3: Determine cache location and file
	cachePath, cacheFileName := determineCacheLocation(repository, gateway)
	params.Log.V(5).Info("using cache", "cachePath", cachePath, "cacheFileName", cacheFileName, "repository", repoRef.Name)

	// Step 4: Check if we should skip delta comparison (vanilla bundle)
	if shouldSkipDeltaComparison(gateway, repository) {
		params.Log.V(5).Info("skipping delta comparison, building vanilla bundle", "repository", repoRef.Name)
		return buildVanillaBundleAndCache(repository, repoRef, gateway, cachePath, cacheFileName, tmpPath, fileName, params)
	}

	// Step 5: Check if commit changed
	newCommit := false
	if _, err := os.Stat(tmpPath + "/" + repository.Status.Commit + ".txt"); err != nil {
		newCommit = true
		params.Log.V(5).Info("new commit detected", "repository", repoRef.Name, "commit", repository.Status.Commit)
	}

	// Step 6: Check if directories changed
	directoryChanged := false
	var previousDirectories []string
	for _, repoStatus := range gateway.Status.RepositoryStatus {
		if repoStatus.Name == repoRef.Name {
			previousDirectories = repoStatus.Directories
			if len(previousDirectories) > 0 && !reflect.DeepEqual(repoRef.Directories, repoStatus.Directories) {
				directoryChanged = true
				params.Log.Info("directory change detected", "repository", repoRef.Name, "previous", previousDirectories, "current", repoRef.Directories)
			}
			break
		}
	}

	// Step 7: Return cached bundle if nothing changed
	if !newCommit && !directoryChanged {
		// If delete/reconcile is enabled, try to use the bundle with mappings first
		// This ensures we continue applying delete mappings if needed
		if !shouldSkipDeltaComparison(gateway, repository) {
			if bundleWithMappings, err := os.ReadFile(tmpPath + "/" + fileName + "_with_mappings"); err == nil {
				params.Log.V(5).Info("returning cached bundle with mappings (no changes)", "repository", repoRef.Name)
				return bundleWithMappings, nil
			}
		}

		// Otherwise use the clean bundle
		if cachedBundle, err := os.ReadFile(tmpPath + "/" + fileName); err == nil {
			params.Log.V(5).Info("returning cached bundle (no changes)", "repository", repoRef.Name)
			return cachedBundle, nil
		}
	}

	// Step 8: Route to appropriate handler
	if directoryChanged && gateway.Spec.App.RepositoryReferenceDelete.ReconcileDirectoryChanges {
		params.Log.V(2).Info("handling directory change with reconciliation", "repository", repoRef.Name)
		return handleDirectoryChange(ctx, params, repository, repoRef, gateway, cachePath, cacheFileName, tmpPath, fileName, previousDirectories)
	} else if directoryChanged && !gateway.Spec.App.RepositoryReferenceDelete.ReconcileDirectoryChanges {
		params.Log.V(5).Info("directory changed but reconciliation disabled, building vanilla bundle", "repository", repoRef.Name)
		return buildVanillaBundleAndCache(repository, repoRef, gateway, cachePath, cacheFileName, tmpPath, fileName, params)
	} else if newCommit {
		params.Log.V(2).Info("handling commit change", "repository", repoRef.Name, "commit", repository.Status.Commit)
		return handleCommitChange(ctx, params, repository, repoRef, gateway, cachePath, cacheFileName, tmpPath, fileName)
	}

	// Fallback: build vanilla bundle
	params.Log.V(5).Info("building vanilla bundle (fallback)", "repository", repoRef.Name)
	return buildVanillaBundleAndCache(repository, repoRef, gateway, cachePath, cacheFileName, tmpPath, fileName, params)
}

// calculateBundleFileName generates a unique filename based on directories and commit
func calculateBundleFileName(gateway *securityv1.Gateway, referenceName string, directories []string) string {
	dirChecksum := ""
	for _, d := range directories {
		h := sha1.New()
		h.Write([]byte(d))
		dirChecksum += fmt.Sprintf("%x", h.Sum(nil))
	}

	h := sha1.New()
	h.Write([]byte(gateway.Name + "-" + referenceName + "-" + dirChecksum))
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	return sha1Sum[30:] + ".json"
}

// buildBundleFromCache loads bundles from cached directory structure or storage secret
func buildBundleFromCache(repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, cachePath string, fileName string) ([]byte, error) {
	//fileName := repository.Status.Commit + ".json"

	bundleMapBytes, err := os.ReadFile(cachePath + "/" + fileName)
	if err != nil {
		// Cache not available - this might happen if repository isn't ready yet
		// In this case, we cannot build the bundle
		return nil, fmt.Errorf("failed to read cached bundle: %w", err)
	}

	bundleMap := map[string][]byte{}
	if err := json.Unmarshal(bundleMapBytes, &bundleMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bundle map: %w", err)
	}

	// Determine if we should preserve repository mappings
	// For combined.json and latest.json (StateStore), preserve DELETE mappings from repository controller
	// For commit.json (user-controlled), clean DELETE mappings for re-added entities
	preserveRepoMappings := (fileName == "combined.json" || fileName == "latest.json")

	// Local and HTTP repos ALWAYS use all directories ["/"] regardless of what's specified in Gateway CR
	// This gives users full control and simplifies the model for development/testing repos
	isLocalOrHttp := (repository.Spec.Type == securityv1.RepositoryTypeLocal || repository.Spec.Type == securityv1.RepositoryTypeHttp)

	// If requesting all directories OR if it's a local/http repo, concatenate everything, no ordering
	if isLocalOrHttp || (len(repoRef.Directories) == 1 && repoRef.Directories[0] == "/") {
		if preserveRepoMappings {
			return util.ConcatBundlesPreservingMappings(bundleMap)
		}
		return util.ConcatBundles(bundleMap)
	}

	// Otherwise, filter by specific directories (git repos only)
	return buildBundleFromDirectories(repoRef.Directories, bundleMap, preserveRepoMappings)
}

// buildBundleFromDirectories combines bundles from specific directories
// Processes directories in order, with later directories overwriting earlier ones
// preserveRepoMappings: if true, preserve DELETE mappings from repository controller (for combined.json)
//
//	if false, clean DELETE mappings for re-added entities (for commit.json)
func buildBundleFromDirectories(directories []string, bundleMap map[string][]byte, preserveRepoMappings bool) ([]byte, error) {
	srcBundle := graphman.Bundle{}
	bundleBytes, err := json.Marshal(srcBundle)
	if err != nil {
		return nil, err
	}

	// Process directories in the order specified
	for _, d := range directories {
		keyName := strings.TrimPrefix(strings.ReplaceAll(d, "/", "-"), "-")

		// Use full bundles (not deltas) for concatenation
		bundleKey := keyName + ".gz"

		if bundleGz, exists := bundleMap[bundleKey]; exists {
			decompressedBytes, err := util.GzipDecompress(bundleGz)
			if err != nil {
				return nil, fmt.Errorf("failed to decompress bundle %s: %w", bundleKey, err)
			}

			var bundle graphman.Bundle
			if err := json.Unmarshal(decompressedBytes, &bundle); err != nil {
				return nil, fmt.Errorf("failed to unmarshal bundle from %s: %w", bundleKey, err)
			}

			cleanBytes, err := json.Marshal(bundle)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal clean bundle from %s: %w", bundleKey, err)
			}

			// Concatenate bundles
			if preserveRepoMappings {
				// For repository-controlled mappings (combined.json, latest.json),
				// preserve DELETE mappings as-is - don't clean them
				bundleBytes, err = graphman.ConcatBundlePreservingMappings(cleanBytes, bundleBytes)
			} else {
				// For user-controlled mappings (commit.json),
				// clean DELETE mappings for re-added entities
				bundleBytes, err = graphman.ConcatBundle(cleanBytes, bundleBytes)
			}
			if err != nil {
				return nil, fmt.Errorf("failed to concat bundle from %s: %w", bundleKey, err)
			}
		}
	}

	return bundleBytes, nil
}

// buildBundleFromStorageSecret reads bundles from storage secret when cache is not available
func buildBundleFromStorageSecret(ctx context.Context, repository *securityv1.Repository, repoRef *securityv1.RepositoryReference, params Params) ([]byte, error) {
	storageSecret, err := getGatewaySecret(ctx, params, repository.Status.StorageSecretName)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage secret: %w", err)
	}

	// Storage secret Data contains the same bundleMap structure as the cache
	// If requesting all directories, concatenate everything
	if len(repoRef.Directories) == 1 && repoRef.Directories[0] == "/" {
		return util.ConcatBundles(storageSecret.Data)
	}

	// Otherwise, filter by specific directories
	// For storage secret, we preserve user-defined mappings (not repository-controlled)
	return buildBundleFromDirectories(repoRef.Directories, storageSecret.Data, false)
}

// cleanupOldBundles removes bundles older than 10 days
func cleanupOldBundles(tmpPath string) {
	existingBundles, err := os.ReadDir(tmpPath)
	if err != nil {
		return
	}

	for _, f := range existingBundles {
		fInfo, err := f.Info()
		if err != nil {
			continue
		}
		if time.Since(fInfo.ModTime()) > 240*time.Hour {
			os.Remove(tmpPath + "/" + f.Name())
		}
	}
}

func checkLocalRepoOnFs(params Params, repository *securityv1.Repository) (bool, error) {

	// Check if pre-built bundle cache exists
	var cachePath string
	if repository.Spec.StateStoreReference != "" {
		cachePath = "/tmp/statestore/" + repository.Name
		fileName := "latest.json"
		if _, err := os.Stat(cachePath + "/" + fileName); err == nil {
			return true, nil
		}
	} else {
		cachePath = "/tmp/repo-cache/" + repository.Name
		fileName := repository.Status.Commit + ".json"
		if _, err := os.Stat(cachePath + "/" + fileName); err == nil {
			// Pre-built bundle cache exists, can use it
			return true, nil
		}
	}

	// If no cache, check if raw repository exists on filesystem
	if repository.Spec.StateStoreReference != "" {
		// For state store repos, check if state store path exists
		stateStorePath := "/tmp/statestore/" + repository.Name
		if _, err := os.Stat(stateStorePath); err == nil {
			return true, nil
		}
	}

	return true, nil
}

func updateGatewayDeployment(ctx context.Context, params Params, gwUpdReq *GatewayUpdateRequest) (err error) {
	update := false
	ready := false
	endpoint := ""

	leaderAvailable := false
	for _, pod := range gwUpdReq.podList.Items {
		if pod.ObjectMeta.Labels["management-access"] == "leader" {
			endpoint = podIP(pod.Status.PodIP) + ":" + strconv.Itoa(gwUpdReq.graphmanPort) + "/graphman"
			leaderAvailable = true
		}
	}

	if !leaderAvailable {
		return nil
	}

	currentChecksum := gwUpdReq.deployment.ObjectMeta.Annotations[gwUpdReq.patchAnnotation]
	// Skip if it already has the correct checksum and no directory change
	if currentChecksum == gwUpdReq.checksum && !gwUpdReq.delete {
		// Check if there's a directory change for repositories
		if gwUpdReq.bundleType == BundleTypeRepository {
			directoryChangeForPod := false
			for _, repoStatus := range gwUpdReq.gateway.Status.RepositoryStatus {
				if repoStatus.Name == gwUpdReq.repositoryReference.Name {
					if !reflect.DeepEqual(gwUpdReq.repositoryReference.Directories, repoStatus.Directories) {
						directoryChangeForPod = true
					}
					break
				}
			}
			if !directoryChangeForPod {
				return nil // skip, it's already up to date
			}
		} else {
			return nil //skip, it's already up to date
		}
	}

	if gwUpdReq.bundleType == BundleTypeRepository {
		if (currentChecksum == "deleted" && !gwUpdReq.repositoryReference.Enabled) || (currentChecksum == "" && (gwUpdReq.delete || !gwUpdReq.repositoryReference.Enabled)) {
			return nil
		}
		for _, repoStatus := range gwUpdReq.gateway.Status.RepositoryStatus {
			if repoStatus.Name == gwUpdReq.repositoryReference.Name {
				if !reflect.DeepEqual(gwUpdReq.repositoryReference.Directories, repoStatus.Directories) {
					update = true
				}
				break
			}
		}
	}

	if currentChecksum != gwUpdReq.checksum || currentChecksum == "" || gwUpdReq.delete {
		update = true
	}

	if gwUpdReq.deployment.Status.ReadyReplicas != 0 {
		ready = true
	}

	// Build patch with annotations
	annotations := make(map[string]string)

	if gwUpdReq.delete {
		annotations[gwUpdReq.patchAnnotation] = "deleted"

		// If ReconcileReferences is enabled, also clear other repository annotations to force reapply
		if gwUpdReq.gateway.Spec.App.RepositoryReferenceDelete.ReconcileReferences {
			for _, repoRef := range gwUpdReq.gateway.Spec.App.RepositoryReferences {
				if repoRef.Name == gwUpdReq.repositoryReference.Name {
					continue
				}
				// Skip static type repositories - there are no singleton configs with database backed gateways
				if repoRef.Type == securityv1.RepositoryReferenceTypeStatic {
					continue
				}
				annotationKey := "security.brcmlabs.com/" + repoRef.Name + "-" + string(repoRef.Type)
				annotations[annotationKey] = ""
			}
		}
	} else {
		annotations[gwUpdReq.patchAnnotation] = gwUpdReq.checksum
	}

	patchData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": annotations,
		},
	}
	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		return err
	}
	patch := string(patchBytes)

	if ready && update {
		requestCacheEntry := gwUpdReq.deployment.Name + "-" + gwUpdReq.cacheEntry
		syncRequest, err := syncCache.Read(requestCacheEntry)
		if err != nil {
			params.Log.V(5).Info("request has not been attempted or cache was flushed", "type", string(gwUpdReq.bundleType), "bundle", gwUpdReq.bundleName, "deployment", gwUpdReq.deployment.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(5).Info("request has been attempted in the last 3 seconds, backing off", "type", string(gwUpdReq.bundleType), "bundle", gwUpdReq.bundleName, "deployment", gwUpdReq.deployment.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
			return errors.New("request has been attempted in the last 3 seconds, backing off")
		}

		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())
		start := time.Now()

		logAction := "applying latest"
		if gwUpdReq.delete {
			logAction = "removing"
		}

		params.Log.V(5).Info(logAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "checksum", gwUpdReq.checksum, "deployment", gwUpdReq.deployment.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
		err = util.ApplyToGraphmanTarget(gwUpdReq.bundle, true, gwUpdReq.username, gwUpdReq.password, endpoint, gwUpdReq.graphmanEncryptionPassphrase, gwUpdReq.delete)
		if err != nil {
			failedAction := "failed to apply"
			if gwUpdReq.delete {
				failedAction = "failed to remove"
			}
			params.Log.Info(failedAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "checksum", gwUpdReq.checksum, "deployment", gwUpdReq.deployment.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
			_ = captureGraphmanMetrics(ctx, params, start, gwUpdReq.deployment.Name, string(gwUpdReq.bundleType), gwUpdReq.bundleName, gwUpdReq.checksum, true)
			return err
		}

		successAction := "applied latest"
		if gwUpdReq.delete {
			successAction = "removed"
		}
		params.Log.Info(successAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "hash", gwUpdReq.checksum, "deployment", gwUpdReq.deployment.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
		_ = captureGraphmanMetrics(ctx, params, start, gwUpdReq.deployment.Name, string(gwUpdReq.bundleType), gwUpdReq.bundleName, gwUpdReq.checksum, false)

		err = updateEntityStatus(ctx, string(gwUpdReq.bundleType), gwUpdReq.bundleName, gwUpdReq.bundle, params)
		if err != nil {
			return err
		}

		if err := params.Client.Patch(ctx, gwUpdReq.deployment,
			client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
			params.Log.Error(err, "failed to update deployment annotations", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
			return err
		}
	} else {
		if (!ready && gwUpdReq.bundleType == BundleTypeClusterProp) || (!ready && gwUpdReq.bundleType == BundleTypeListenPort) {
			if err := params.Client.Patch(ctx, gwUpdReq.deployment,
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "failed to update deployment annotations", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
				return err
			}
		}
	}
	return nil
}

func updateGatewayPods(ctx context.Context, params Params, gwUpdReq *GatewayUpdateRequest) (err error) {
	updateStatus := false
	for i, pod := range gwUpdReq.podList.Items {

		singleton := false
		checksum := gwUpdReq.checksum
		update := false
		ready := false

		if pod.DeletionTimestamp != nil {
			continue
		}

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		if gwUpdReq.otkCerts {
			if pod.ObjectMeta.Annotations["security.brcmlabs.com/"+gwUpdReq.gateway.Name+"-"+string(gwUpdReq.gateway.Spec.App.Otk.Type)+"-policies"] == "" {
				ready = false
			}
		}

		currentChecksum := pod.ObjectMeta.Annotations[gwUpdReq.patchAnnotation]

		// For repositories with SingletonExtraction, adjust checksum for leader before comparison
		if gwUpdReq.bundleType == BundleTypeRepository && gwUpdReq.gateway.Spec.App.SingletonExtraction && pod.ObjectMeta.Labels["management-access"] == "leader" {
			checksum = gwUpdReq.checksum + "-leader"
		}

		// Skip this pod if it already has the correct checksum and no directory change
		if currentChecksum == checksum && !gwUpdReq.delete {
			// Check if there's a directory change for repositories
			if gwUpdReq.bundleType == BundleTypeRepository {
				directoryChangeForPod := false
				for _, repoStatus := range gwUpdReq.gateway.Status.RepositoryStatus {
					if repoStatus.Name == gwUpdReq.repositoryReference.Name {
						if !reflect.DeepEqual(gwUpdReq.repositoryReference.Directories, repoStatus.Directories) {
							directoryChangeForPod = true
						}
						break
					}
				}
				if !directoryChangeForPod {
					continue // Skip this pod, it's already up to date
				}
			} else {
				continue // Skip this pod, it's already up to date
			}
		}

		if gwUpdReq.bundleType == BundleTypeOTKDatabaseMaintenance {
			if pod.ObjectMeta.Labels["management-access"] == "leader" {
				checksum = gwUpdReq.checksum + "-leader"
				singleton = true
			} else {
				continue
			}
		}

		if gwUpdReq.bundleType == BundleTypeRepository {
			// Skip if already deleted (annotation = "deleted" and repo still disabled)
			if currentChecksum == "deleted" && !gwUpdReq.repositoryReference.Enabled {
				return nil
			}

			for _, repoStatus := range gwUpdReq.gateway.Status.RepositoryStatus {
				if repoStatus.Name == gwUpdReq.repositoryReference.Name {
					if !reflect.DeepEqual(gwUpdReq.repositoryReference.Directories, repoStatus.Directories) {
						update = true
					}
					break
				}
			}

			if gwUpdReq.gateway.Spec.App.SingletonExtraction && pod.ObjectMeta.Labels["management-access"] == "leader" {
				checksum = gwUpdReq.checksum + "-leader"
				singleton = true
			}

			// Handle static repositories - only apply to leader pod
			if gwUpdReq.repositoryReference.Type == securityv1.RepositoryReferenceTypeStatic {
				if pod.ObjectMeta.Labels["management-access"] != "leader" {
					continue // Skip non-leader pods for static repos
				}

				// Extract singleton entities for static repos
				bundle := graphman.Bundle{}
				singletonBundle := graphman.Bundle{}
				err = json.Unmarshal(gwUpdReq.bundle, &bundle)
				if err != nil {
					return err
				}

				for _, scheduledTask := range bundle.ScheduledTasks {
					if scheduledTask.ExecuteOnSingleNode {
						singletonBundle.ScheduledTasks = append(singletonBundle.ScheduledTasks, scheduledTask)
					}
				}
				for _, jmsDestination := range bundle.JmsDestinations {
					if jmsDestination.Direction == "OUTBOUND" {
						singletonBundle.JmsDestinations = append(singletonBundle.JmsDestinations, jmsDestination)
					}
				}

				if len(singletonBundle.ScheduledTasks) > 0 || len(singletonBundle.JmsDestinations) > 0 {
					gwUpdReq.bundle, err = json.Marshal(singletonBundle)
					if err != nil {
						return err
					}
				} else {
					continue // No singleton entities to apply for static repo
				}
			}
		}

		// Build patch with annotations
		annotations := make(map[string]string)

		if gwUpdReq.delete {
			annotations[gwUpdReq.patchAnnotation] = "deleted"

			// If ReconcileReferences is enabled, also clear other repository annotations to force reapply
			if gwUpdReq.gateway.Spec.App.RepositoryReferenceDelete.ReconcileReferences {
				for _, repoRef := range gwUpdReq.gateway.Spec.App.RepositoryReferences {
					if repoRef.Name == gwUpdReq.repositoryReference.Name {
						continue
					}
					annotationKey := "security.brcmlabs.com/" + repoRef.Name + "-" + string(repoRef.Type)
					annotations[annotationKey] = ""
				}
			}
		} else {
			annotations[gwUpdReq.patchAnnotation] = checksum
		}

		patchData := map[string]interface{}{
			"metadata": map[string]interface{}{
				"annotations": annotations,
			},
		}
		patchBytes, err := json.Marshal(patchData)
		if err != nil {
			return err
		}
		patch := string(patchBytes)

		// Set update=true if checksums don't match, checksum is empty, or it's a delete
		// (update may already be true from directory change check above)
		if currentChecksum != checksum || currentChecksum == "" || gwUpdReq.delete {
			update = true
		}

		if update && ready {
			updateStatus = true
			endpoint := podIP(pod.Status.PodIP) + ":" + strconv.Itoa(gwUpdReq.graphmanPort) + "/graphman"
			requestCacheEntry := pod.Name + "-" + gwUpdReq.cacheEntry
			syncRequest, err := syncCache.Read(requestCacheEntry)
			tryRequest := true
			if err != nil {
				params.Log.V(5).Info("request has not been attempted or cache was flushed", "type", gwUpdReq.bundleType, "name", gwUpdReq.bundleName, "pod", pod.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
			}

			if syncRequest.Attempts > 0 {
				params.Log.V(5).Info("request has been attempted in the last 3 seconds, backing off", "type", gwUpdReq.bundleType, "name", gwUpdReq.bundleName, "pod", pod.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
				tryRequest = false
				return errors.New("request has been attempted in the last 3 seconds, backing off")
			}

			if tryRequest {
				syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())
				start := time.Now()

				logAction := "applying latest"
				if gwUpdReq.delete {
					logAction = "removing"
				}

				params.Log.V(5).Info(logAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "checksum", checksum, "pod", pod.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
				err = util.ApplyToGraphmanTarget(gwUpdReq.bundle, singleton, gwUpdReq.username, gwUpdReq.password, endpoint, gwUpdReq.graphmanEncryptionPassphrase, gwUpdReq.delete)
				if err != nil {
					failedAction := "failed to apply"
					if gwUpdReq.delete {
						failedAction = "failed to remove"
					}
					params.Log.Info(failedAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "checksum", checksum, "pod", pod.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
					_ = captureGraphmanMetrics(ctx, params, start, pod.Name, string(gwUpdReq.bundleType), gwUpdReq.bundleName, checksum, true)
					return err
				}

				successAction := "applied latest"
				if gwUpdReq.delete {
					successAction = "removed"
				}
				params.Log.Info(successAction+" "+string(gwUpdReq.bundleType)+" "+gwUpdReq.bundleName, "hash", checksum, "pod", pod.Name, "name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
				_ = captureGraphmanMetrics(ctx, params, start, pod.Name, string(gwUpdReq.bundleType), gwUpdReq.bundleName, checksum, false)

				if err := params.Client.Patch(ctx, &gwUpdReq.podList.Items[i],
					client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
					params.Log.Error(err, "failed to update pod label", "Name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
					return err
				}
			}
		} else {
			// Patch annotation for non-ready pods
			if (!ready && gwUpdReq.bundleType == BundleTypeClusterProp) ||
				(!ready && gwUpdReq.bundleType == BundleTypeListenPort) ||
				(!ready && gwUpdReq.bundleType == BundleTypeRepository && gwUpdReq.gateway.Spec.App.RepositoryReferenceBootstrap.Enabled && pod.ObjectMeta.Labels["management-access"] != "leader" && !singleton) {
				if err := params.Client.Patch(ctx, &gwUpdReq.podList.Items[i],
					client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
					params.Log.Error(err, "failed to update pod label", "Name", gwUpdReq.gateway.Name, "namespace", gwUpdReq.gateway.Namespace)
					return err
				}
			}
		}
	}

	if updateStatus || (!updateStatus && gwUpdReq.bundleType == BundleTypeClusterProp) || (!updateStatus && gwUpdReq.bundleType == BundleTypeListenPort) {
		err := updateEntityStatus(ctx, string(gwUpdReq.bundleType), gwUpdReq.bundleName, gwUpdReq.bundle, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func readLocalReference(ctx context.Context, repository *securityv1.Repository, params Params) ([]byte, error) {
	if repository.Spec.LocalReference.SecretName == "" {
		return nil, fmt.Errorf("%s localReference secret name must be set", repository.Name)
	}

	localReference := &corev1.Secret{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.LocalReference.SecretName, Namespace: repository.Namespace}, localReference)
	if err != nil {
		return nil, err
	}

	bundleBytes, err := util.ConcatBundles(localReference.Data)
	if err != nil {
		return nil, err
	}

	return bundleBytes, nil
}

func readStorageSecret(ctx context.Context, repository *securityv1.Repository, params Params) ([]byte, error) {
	if repository.Status.StorageSecretName == "_" {
		return nil, fmt.Errorf("%s storage secret does not exist", repository.Name)
	}

	storageSecret := &corev1.Secret{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Status.StorageSecretName, Namespace: repository.Namespace}, storageSecret)
	if err != nil {
		return nil, err
	}

	bundleBytes, err := util.ConcatBundles(storageSecret.Data)
	if err != nil {
		return nil, err
	}

	return bundleBytes, nil
}

// GetGatewayPods returns the pods in a Gateway Deployment
func getGatewayPods(ctx context.Context, params Params) (*corev1.PodList, error) {
	podList := &corev1.PodList{}

	listOpts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(util.DefaultLabels(params.Instance.Name, map[string]string{})),
	}
	if err := params.Client.List(ctx, podList, listOpts...); err != nil {
		return podList, err
	}
	return podList, nil
}

func getGatewayDeployment(ctx context.Context, params Params) (*appsv1.Deployment, error) {
	gatewayDeployment := &appsv1.Deployment{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gatewayDeployment)
	if err != nil {
		return gatewayDeployment, err
	}
	return gatewayDeployment, nil
}

func getGraphmanEncryptionPassphrase(ctx context.Context, params Params, existingSecretName string, existingSecretKey string) (string, error) {
	graphmanEncryptionSecret, err := getGatewaySecret(ctx, params, existingSecretName)
	if err != nil {
		return "", err
	}
	return string(graphmanEncryptionSecret.Data[existingSecretKey]), nil
}

func getGatewaySecret(ctx context.Context, params Params, name string) (*corev1.Secret, error) {
	gwSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, gwSecret)
	if err != nil {
		return gwSecret, err
	}
	return gwSecret, nil
}

func getGatewayConfigMap(ctx context.Context, params Params, name string) (*corev1.ConfigMap, error) {
	gwConfigmap := &corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, gwConfigmap)
	if err != nil {
		return gwConfigmap, err
	}
	return gwConfigmap, nil
}

func parseGatewaySecret(gwSecret *corev1.Secret) (string, string) {
	var username string
	var password string
	if string(gwSecret.Data["node.properties"]) != "" {
		usernameRe := regexp.MustCompile(`(?m)(admin.user=)(.*)`)
		passwordRe := regexp.MustCompile(`(?m)(admin.pass=)(.*)`)
		username = usernameRe.FindStringSubmatch(string(gwSecret.Data["node.properties"]))[2]
		password = passwordRe.FindStringSubmatch(string(gwSecret.Data["node.properties"]))[2]
	} else {
		username = string(gwSecret.Data["SSG_ADMIN_USERNAME"])
		password = string(gwSecret.Data["SSG_ADMIN_PASSWORD"])
	}
	return username, password
}

func getStateStoreSecret(ctx context.Context, name string, statestore securityv1alpha1.L7StateStore, params Params) (*corev1.Secret, error) {
	statestoreSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: statestore.Namespace}, statestoreSecret)
	if err != nil {
		return statestoreSecret, err
	}
	return statestoreSecret, nil
}

// HardenGraphmanService adds required mutual TLS to the Gateway's GraphQL Management API (Graphman)
// This process also creates a user (PKI) and restricts Graphman to that user effectively locking remote Gateway management to
// the Layer7 Operator only.
// This feature is intended for Ephemeral Gateways, while it will work for MySQL backed Gateways we strongly recommend you supply your own
// PKI Pair as losing this means you will need to update the user in Policy Manager as no remote interaction will be available.
func HardenGraphmanService(ctx context.Context, params Params) error {
	// potentially bootstrap this...
	return nil

}

func GatewayLicense(ctx context.Context, params Params) error {
	gatewayLicense := &corev1.Secret{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Spec.License.SecretName, Namespace: params.Instance.Namespace}, gatewayLicense)
	if k8serrors.IsNotFound(err) {
		params.Log.Error(err, "license not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func ManagementPod(ctx context.Context, params Params) error {
	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	for p := range podList.Items {
		if podList.Items[p].Labels["management-access"] == "leader" {
			if podList.Items[p].DeletionTimestamp == nil {
				return nil
			}
		}
	}
	tagged := false
	for p := range podList.Items {
		if podList.Items[p].Status.Phase == "Running" && podList.Items[p].DeletionTimestamp == nil && !tagged {
			patch := []byte(`{"metadata":{"labels":{"management-access": "leader"}}}`)
			if err := params.Client.Patch(ctx, &podList.Items[p],
				client.RawPatch(types.StrategicMergePatchType, patch)); err != nil {
				params.Log.Error(err, "failed to update pod label", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
				return err
			}
			params.Log.V(2).Info("new leader elected", "name", params.Instance.Name, "pod", podList.Items[p].Name, "namespace", params.Instance.Namespace)
			tagged = true
		}
	}
	return nil
}

func ReconcileEphemeralGateway(ctx context.Context, params Params, kind string, podList corev1.PodList, gateway *securityv1.Gateway, gwSecret *corev1.Secret, graphmanEncryptionPassphrase string, annotation string, sha1Sum string, otkCerts bool, name string, bundle []byte) error {
	graphmanPort := 9443

	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	username, password := parseGatewaySecret(gwSecret)

	if username == "" || password == "" {
		return fmt.Errorf("could not retrieve gateway credentials for %s", name)
	}

	updateStatus := false

	for i, pod := range podList.Items {
		currentSha1Sum := pod.ObjectMeta.Annotations[annotation]

		update := false
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		if otkCerts {
			if pod.ObjectMeta.Annotations["security.brcmlabs.com/"+gateway.Name+"-"+string(gateway.Spec.App.Otk.Type)+"-policies"] == "" {
				ready = false
			}
		}

		patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", annotation, sha1Sum)

		if currentSha1Sum != sha1Sum || currentSha1Sum == "" {
			update = true
		}

		if update && ready {
			updateStatus = true
			endpoint := podIP(pod.Status.PodIP) + ":" + strconv.Itoa(graphmanPort) + "/graphman"

			requestCacheEntry := pod.Name + "-" + gateway.Name + "-" + name + "-" + sha1Sum
			syncRequest, err := syncCache.Read(requestCacheEntry)
			tryRequest := true
			if err != nil {
				params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync "+kind, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
			}

			if syncRequest.Attempts > 0 {
				params.Log.V(2).Info("request has been attempted in the last 3 seconds, backing off", "hash", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
				tryRequest = false
			}

			if tryRequest {
				syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())
				start := time.Now()
				params.Log.V(2).Info("applying latest "+kind, "hash", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
				err = util.ApplyGraphmanBundle(username, password, endpoint, graphmanEncryptionPassphrase, bundle)
				if err != nil {
					params.Log.Info("failed to apply "+kind, "hash", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					_ = captureGraphmanMetrics(ctx, params, start, pod.Name, kind, name, sha1Sum, true)
					return err
				}
				_ = captureGraphmanMetrics(ctx, params, start, pod.Name, kind, name, sha1Sum, false)
				params.Log.Info("applied latest "+kind, "hash", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)

				if err := params.Client.Patch(ctx, &podList.Items[i],
					client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
					params.Log.Error(err, "failed to update pod label", "Name", gateway.Name, "namespace", gateway.Namespace)
					return err
				}

			}
		}

		// if the Gateway is not ready then cluster properties and listenPorts have already been applied via bootsrap
		if (!ready && kind == "cluster properties") || (!ready && kind == "listen ports") {
			if err := params.Client.Patch(ctx, &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "failed to update pod label", "Name", gateway.Name, "namespace", gateway.Namespace)
				return err
			}
		}
	}

	if updateStatus || (!updateStatus && kind == "cluster properties") || (!updateStatus && kind == "listen ports") {
		err := updateEntityStatus(ctx, kind, name, bundle, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReconcileDBGateway(ctx context.Context, params Params, kind string, gatewayDeployment appsv1.Deployment, gateway *securityv1.Gateway, gwSecret *corev1.Secret, graphmanEncryptionPassphrase string, annotation string, sha1Sum string, otkCerts bool, name string, bundle []byte) error {
	graphmanPort := 9443

	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	username, password := parseGatewaySecret(gwSecret)
	if username == "" || password == "" {
		return fmt.Errorf("could not retrieve gateway credentials for %s", name)
	}

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", annotation, sha1Sum)

	ready := false

	if gatewayDeployment.ObjectMeta.Annotations[annotation] == sha1Sum {
		return nil
	}

	if gatewayDeployment.Status.ReadyReplicas == gatewayDeployment.Status.Replicas {
		ready = true
	}

	if ready {
		requestCacheEntry := gatewayDeployment.Name + "-" + name + "-" + sha1Sum
		syncRequest, err := syncCache.Read(requestCacheEntry)
		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync "+kind, "Name", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 3 seconds, backing off", "hash", sha1Sum, "Name", gateway.Name, "Namespace", gateway.Namespace)
			return errors.New("request has been attempted in the last 3 seconds, backing off")

		}
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())

		endpoint := gateway.Name + "." + gateway.Namespace + ".svc.cluster.local:" + strconv.Itoa(graphmanPort) + "/graphman"
		if gateway.Spec.App.Management.Service.Enabled {
			endpoint = gateway.Name + "-management-service." + gateway.Namespace + ".svc.cluster.local:" + strconv.Itoa(graphmanPort) + "/graphman"
		}
		start := time.Now()
		params.Log.V(2).Info("applying latest "+kind, "sha1Sum", sha1Sum, "name", gateway.Name, "namespace", gateway.Namespace)

		err = util.ApplyGraphmanBundle(username, password, endpoint, graphmanEncryptionPassphrase, bundle)
		if err != nil {
			params.Log.Info("failed to apply "+kind, "sha1Sum", sha1Sum, "name", gateway.Name, "namespace", gateway.Namespace)
			_ = captureGraphmanMetrics(ctx, params, start, gateway.Name, kind, name, sha1Sum, true)
			return err
		}

		params.Log.Info("applied latest "+kind, "sha1Sum", sha1Sum, "name", gateway.Name, "namespace", gateway.Namespace)
		_ = captureGraphmanMetrics(ctx, params, start, gateway.Name, kind, name, sha1Sum, false)

		err = updateEntityStatus(ctx, kind, name, bundle, params)
		if err != nil {
			return err
		}

		if err := params.Client.Patch(ctx, &gatewayDeployment,
			client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
			params.Log.Error(err, "Failed to update deployment annotations", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
			return err
		}
	}
	return nil
}

func updateRepoRefStatus(ctx context.Context, params Params, repository securityv1.Repository, repoRef securityv1.RepositoryReference, commit string, applyError error, delete bool) (err error) {
	gatewayStatus := params.Instance.Status

	// If delete was successful, remove the repository from status
	if delete && applyError == nil {
		var updatedStatus []securityv1.GatewayRepositoryStatus
		for _, rs := range gatewayStatus.RepositoryStatus {
			if rs.Name != repoRef.Name {
				updatedStatus = append(updatedStatus, rs)
			}
		}
		gatewayStatus.RepositoryStatus = updatedStatus

		params.Instance.Status = gatewayStatus
		err = params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.V(2).Info("failed to update gateway status after delete", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
		params.Log.Info("removed repository from gateway status", "repository", repoRef.Name, "gateway", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	var conditions []securityv1.RepositoryCondition
	secretName := repository.Name
	if repository.Spec.Auth.ExistingSecretName != "" {
		secretName = repository.Spec.Auth.ExistingSecretName
	}

	if repository.Spec.Auth == (securityv1.RepositoryAuth{}) {
		secretName = ""
	}

	nrs := securityv1.GatewayRepositoryStatus{
		Commit:            commit,
		Enabled:           !delete,
		Name:              repoRef.Name,
		RepoType:          string(repository.Spec.Type),
		Vendor:            repository.Spec.Auth.Vendor,
		AuthType:          string(repository.Spec.Auth.Type),
		Type:              string(repoRef.Type),
		SecretName:        secretName,
		StorageSecretName: repository.Status.StorageSecretName,
		Endpoint:          repository.Spec.Endpoint,
		Directories:       repoRef.Directories,
	}

	if repository.Spec.Tag != "" && repository.Spec.Branch == "" {
		nrs.Tag = repository.Spec.Tag
	}

	if repository.Spec.Branch != "" {
		nrs.Branch = repository.Spec.Branch
	}

	nrs.RemoteName = "origin"
	if repository.Spec.RemoteName != "" {
		nrs.RemoteName = repository.Spec.RemoteName
	}

	// cleanup old conditions
	for _, ors := range gatewayStatus.RepositoryStatus {
		if ors.Name == repository.Name {
			nrs.Conditions = ors.Conditions
		}
	}

	for _, condition := range nrs.Conditions {
		t, err := time.Parse(time.RFC3339, condition.Time)

		if err != nil {
			return err
		}
		// if condition is older than 5 minutes, clean up
		if t.Add(5 * time.Minute).Before(time.Now()) {
			continue
		}
		conditions = append(conditions, condition)
	}

	if applyError != nil {
		errorMsg := applyError.Error()

		if len(errorMsg) > 200 {
			errorMsg = "gateway failed to apply repository"
		}

		conditions = append(conditions, securityv1.RepositoryCondition{
			Time:   time.Now().Format(time.RFC3339),
			Status: "FAILURE",
			Reason: errorMsg,
		})
	} else {
		conditions = []securityv1.RepositoryCondition{}
		conditions = append(conditions, securityv1.RepositoryCondition{
			Time:   time.Now().Format(time.RFC3339),
			Status: "SUCCESS",
			Reason: "",
		})
	}

	nrs.Conditions = conditions

	if repository.Spec.StateStoreReference != "" {
		ext := repository.Spec.Branch
		if ext == "" {
			ext = repository.Spec.Tag
		}
		stateStoreKey := repository.Name + "-repository-" + ext
		nrs.StateStoreReference = repository.Spec.StateStoreReference
		statestore := &securityv1alpha1.L7StateStore{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.StateStoreReference, Namespace: params.Instance.Namespace}, statestore)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("state store not found", "name", repository.Spec.StateStoreReference, "repository", repository.Name, "namespace", params.Instance.Namespace)
			return err
		}
		nrs.StateStoreKey = statestore.Spec.Redis.GroupName + ":" + statestore.Spec.Redis.StoreId + ":" + "repository" + ":" + stateStoreKey + ":latest"
		if repository.Spec.StateStoreKey != "" {
			nrs.StateStoreKey = repository.Spec.StateStoreKey
		}
	}

	found := false
	for i, rs := range gatewayStatus.RepositoryStatus {
		if rs.Name == nrs.Name {
			gatewayStatus.RepositoryStatus[i] = nrs
			found = true
		}
	}

	if !found {
		gatewayStatus.RepositoryStatus = append(gatewayStatus.RepositoryStatus, nrs)
	}

	params.Instance.Status = gatewayStatus
	err = params.Client.Status().Update(ctx, params.Instance)
	if err != nil {
		params.Log.V(2).Info("failed to update gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
		return err
	}
	params.Log.V(2).Info("updated gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	return nil
}

func updateEntityStatus(ctx context.Context, kind string, name string, bundleBytes []byte, params Params) error {
	switch kind {
	case "cluster properties":
		bundle := graphman.Bundle{}
		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}
		clusterProps := []string{}
		if params.Instance.Status.LastAppliedClusterProperties == nil {
			for _, cwp := range params.Instance.Spec.App.ClusterProperties.Properties {
				clusterProps = append(clusterProps, cwp.Name)
			}
		} else {
			for _, appliedCwp := range bundle.ClusterProperties {
				mappingSource := MappingSource{}
				found := false
				for _, cwp := range params.Instance.Status.LastAppliedClusterProperties {
					if cwp == appliedCwp.Name {
						for _, mapping := range bundle.Properties.Mappings.ClusterProperties {
							sourceBytes, err := json.Marshal(mapping.Source)
							if err != nil {
								return err
							}
							err = json.Unmarshal(sourceBytes, &mappingSource)
							if err != nil {
								return err
							}
							if appliedCwp.Name == mappingSource.Name && mapping.Action == graphman.MappingActionDelete {
								found = true
							}
						}
					}
				}
				if !found {
					clusterProps = append(clusterProps, appliedCwp.Name)
				}
			}
		}
		params.Instance.Status.LastAppliedClusterProperties = clusterProps
		if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
			return fmt.Errorf("failed to update cluster properties status: %w", err)
		}
	case "listen ports":
		bundle := graphman.Bundle{}
		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}
		listenPorts := []string{}
		if params.Instance.Status.LastAppliedListenPorts == nil {
			for _, listenPort := range params.Instance.Spec.App.ListenPorts.Ports {
				listenPorts = append(listenPorts, listenPort.Name)
			}
		} else {
			for _, appliedListenPort := range bundle.ListenPorts {
				mappingSource := MappingSource{}
				found := false
				for _, lp := range params.Instance.Status.LastAppliedListenPorts {
					if lp == appliedListenPort.Name {
						for _, mapping := range bundle.Properties.Mappings.ListenPorts {
							sourceBytes, err := json.Marshal(mapping.Source)
							if err != nil {
								return err
							}
							err = json.Unmarshal(sourceBytes, &mappingSource)
							if err != nil {
								return err
							}
							if appliedListenPort.Name == mappingSource.Name && mapping.Action == graphman.MappingActionDelete {
								found = true
							}
						}
					}
				}
				if !found {
					listenPorts = append(listenPorts, appliedListenPort.Name)
				}
			}
		}
		params.Instance.Status.LastAppliedListenPorts = listenPorts
		if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
			return fmt.Errorf("failed to update listenPort status: %w", err)
		}
	case "external secrets":
		bundle := graphman.Bundle{}
		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}
		secrets := []string{}
		if params.Instance.Status.LastAppliedExternalSecrets == nil {
			for _, secret := range bundle.Secrets {
				secrets = append(secrets, secret.Name)
			}
		} else {
			for _, appliedSecret := range bundle.Secrets {
				mappingSource := MappingSource{}
				found := false
				for _, secret := range params.Instance.Status.LastAppliedExternalSecrets[name] {
					if bundle.Properties != nil && secret == appliedSecret.Name {
						for _, mapping := range bundle.Properties.Mappings.Secrets {
							sourceBytes, err := json.Marshal(mapping.Source)
							if err != nil {
								return err
							}
							err = json.Unmarshal(sourceBytes, &mappingSource)
							if err != nil {
								return err
							}
							if appliedSecret.Name == mappingSource.Name && mapping.Action == graphman.MappingActionDelete {
								found = true
							}
						}
					}
				}
				if !found {
					secrets = append(secrets, appliedSecret.Name)
				}
			}
		}
		if params.Instance.Status.LastAppliedExternalSecrets == nil {
			params.Instance.Status.LastAppliedExternalSecrets = map[string][]string{}
		}

		params.Instance.Status.LastAppliedExternalSecrets[name] = secrets
		if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
			return fmt.Errorf("failed to update external secret status: %w", err)
		}
	case "external keys":
		bundle := graphman.Bundle{}

		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}
		keys := []string{}
		if params.Instance.Status.LastAppliedExternalKeys == nil {
			for _, key := range bundle.Keys {
				keys = append(keys, key.Alias)
			}
		} else {
			for _, appliedKey := range bundle.Keys {
				mappingSource := MappingSource{}
				found := false
				for _, key := range params.Instance.Status.LastAppliedExternalKeys {
					if bundle.Properties != nil && key == appliedKey.Alias {
						for _, mapping := range bundle.Properties.Mappings.Keys {
							sourceBytes, err := json.Marshal(mapping.Source)
							if err != nil {
								return err
							}
							err = json.Unmarshal(sourceBytes, &mappingSource)
							if err != nil {
								return err
							}
							if appliedKey.Alias == mappingSource.Alias && mapping.Action == graphman.MappingActionDelete {
								found = true
							}
						}
					}
				}
				if !found {
					keys = append(keys, appliedKey.Alias)
				}
			}
		}
		if params.Instance.Status.LastAppliedExternalKeys == nil {
			params.Instance.Status.LastAppliedExternalKeys = []string{}
		}

		params.Instance.Status.LastAppliedExternalKeys = keys
		if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
			return fmt.Errorf("failed to update external key status: %w", err)
		}
	case "external certs":
		bundle := graphman.Bundle{}

		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}
		certs := []string{}
		if params.Instance.Status.LastAppliedExternalCerts == nil {
			for _, cert := range bundle.TrustedCerts {
				certs = append(certs, cert.Name+"-"+cert.ThumbprintSha1)
			}
		} else {
			for _, appliedCert := range bundle.TrustedCerts {
				mappingSource := MappingSource{}
				found := false
				for _, cert := range params.Instance.Status.LastAppliedExternalCerts[name] {
					if bundle.Properties != nil && strings.Split(cert, "-")[0] == appliedCert.Name {
						for _, mapping := range bundle.Properties.Mappings.TrustedCerts {
							sourceBytes, err := json.Marshal(mapping.Source)
							if err != nil {
								return err
							}
							err = json.Unmarshal(sourceBytes, &mappingSource)
							if err != nil {
								return err
							}
							if appliedCert.ThumbprintSha1 == mappingSource.ThumbprintSha1 && mapping.Action == graphman.MappingActionDelete {
								found = true
							}
						}
					}
				}
				if !found {
					certs = append(certs, appliedCert.Name+"-"+appliedCert.ThumbprintSha1)
				}
			}
		}
		if params.Instance.Status.LastAppliedExternalCerts == nil {
			params.Instance.Status.LastAppliedExternalCerts = map[string][]string{}
		}

		params.Instance.Status.LastAppliedExternalCerts[name] = certs
		if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
			return fmt.Errorf("failed to update external cert status: %w", err)
		}
	}

	return nil
}

func getStateStore(ctx context.Context, params Params, stateStoreName string) (securityv1alpha1.L7StateStore, error) {
	statestore := securityv1alpha1.L7StateStore{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: stateStoreName, Namespace: params.Instance.Namespace}, &statestore)
	if err != nil {
		return statestore, err
	}
	return statestore, nil
}

func isIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

func podIP(podIp string) string {
	if isIPv6(podIp) {
		return "[" + podIp + "]"
	}
	return podIp
}
