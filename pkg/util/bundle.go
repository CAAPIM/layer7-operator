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
package util

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
)

type Bundle struct {
	XMLName    xml.Name   `xml:"l7:Bundle"`
	XMLNS      string     `xml:"xmlns:l7,attr"`
	References References `xml:"l7:References"`
	Mappings   Mappings   `xml:"l7:Mappings"`
}
type PortalBundle struct {
	XMLName    xml.Name   `xml:"Bundle"`
	XMLNS      string     `xml:"xmlns:l7,attr"`
	References References `xml:"l7:References"`
	Mappings   Mappings   `xml:"l7:Mappings"`
}

type References struct {
	Item []Item `xml:"l7:Item"`
}

type Item struct {
	Name     string   `xml:"l7:Name"`
	ID       string   `xml:"l7:Id"`
	Type     string   `xml:"l7:Type"`
	Resource Resource `xml:"l7:Resource"`
}

type Resource struct {
	Policy  *Policy  `xml:"l7:Policy,omitempty"`
	Service *Service `xml:"l7:Service,omitempty"`
}

type Policy struct {
	Guid         string       `xml:"guid,attr"`
	ID           string       `xml:"id,attr"`
	Version      string       `xml:"version,attr"`
	PolicyDetail PolicyDetail `xml:"l7:PolicyDetail"`

	Resources PolicyResources
}

type PolicyResources struct {
	ResourceSet PolicyResourceSet `xml:"l7:ResourceSet"`
}

type PolicyResourceSet struct {
	Tag      string         `xml:"tag,attr"`
	Resource PolicyResource `xml:"l7:Resource"`
}

type PolicyResource struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type PolicyDetail struct {
	FolderId   string                 `xml:"folderId,attr"`
	Guid       string                 `xml:"guid,attr"`
	ID         string                 `xml:"id,attr"`
	Name       string                 `xml:"l7:Name"`
	PolicyType string                 `xml:"l7:PolicyType"`
	Properties PolicyDetailProperties `xml:"l7:Properties"`
}

type PolicyDetailProperties struct {
	Text     string `xml:",chardata"`
	Property []PolicyProperty
}

type PolicyProperty struct {
	Text         string `xml:",chardata"`
	Key          string `xml:"key,attr"`
	LongValue    string `xml:"LongValue"`
	BooleanValue string `xml:"BooleanValue"`
}

type Service struct {
	ID            string           `xml:"id,attr"`
	L7            string           `xml:"l7,attr"`
	ServiceDetail ServiceDetail    `xml:"l7:ServiceDetail"`
	Properties    ServiceProperty  `xml:"l7:Properties"`
	Resources     ServiceResources `xml:"l7:Resources"`
}

type ServiceDetail struct {
	Text            string            `xml:",chardata"`
	FolderId        string            `xml:"folderId,attr"`
	ID              string            `xml:"id,attr"`
	Name            string            `xml:"l7:Name"`
	Enabled         bool              `xml:"l7:Enabled"`
	ServiceMappings ServiceMappings   `xml:"l7:ServiceMappings"`
	Properties      ServiceProperties `xml:"l7:Properties"`
}

type ServiceMappings struct {
	HttpMapping HttpMapping
}

type HttpMapping struct {
	UrlPattern string `xml:"l7:UrlPattern"`
	Verbs      Verbs  `xml:"l7:Verbs"`
}

type Verbs struct {
	Verb []string `xml:"l7:Verb"`
}

type ServiceProperties struct {
	Property []ServiceProperty `xml:"l7:Property"`
}

type ServiceProperty struct {
	Key          string `xml:"key,attr"`
	BooleanValue string `xml:"l7:BooleanValue"`
	StringValue  string `xml:"l7:StringValue"`
}

type ServiceResources struct {
	ResourceSet ServiceResourceSet `xml:"l7:ResourceSet"`
}

type ServiceResourceSet struct {
	Tag      string          `xml:"tag,attr"`
	Resource ServiceResource `xml:"l7:Resource"`
}

type ServiceResource struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type Mappings struct {
	Mapping []Mapping `xml:"l7:Mapping"`
}
type Mapping struct {
	Action     string     `xml:"action,attr"`
	SrcId      string     `xml:"srcId,attr"`
	Type       string     `xml:"type,attr"`
	Properties Properties `xml:"l7:Properties"`
}

type Properties struct {
	Property []Property `xml:"l7:Property"`
}

type Property struct {
	Key          string `xml:"key,attr"`
	StringValue  string `xml:"l7:StringValue,omitempty"`
	BooleanValue bool   `xml:"l7:BooleanValue,omitempty"`
}

type PolicyXml struct {
	XMLName  xml.Name   `xml:"wsp:Policy"`
	XMLNSL7p string     `xml:"xmlns:L7p,attr"`
	XMLNSWsp string     `xml:"xmlns:wsp,attr"`
	All      PolicyBody `xml:"wsp:All"`
}

type PolicyBody struct {
	Usage            string                  `xml:"wsp:Usage,attr"`
	CommentAssertion *CommentAssertion       `xml:"L7p:CommentAssertion"`
	SetVariable      *[]SetVariableAssertion `xml:"L7p:SetVariable"`
	OneOrMore        *OneOrMore              `xml:"wsp:OneOrMore"`
}

type CommentAssertion struct {
	Comment CommentAssertionComment `xml:"L7p:Comment"`
}

type CommentAssertionComment struct {
	StringValue string `xml:"stringValue,attr"`
}

type SetVariableAssertion struct {
	AssertionComment AssertionComment  `xml:"L7p:AssertionComment"`
	Base64Expression PolicyStringValue `xml:"L7p:Base64Expression"`
	VariableToSet    PolicyStringValue `xml:"L7p:VariableToSet"`
}

type AssertionComment struct {
	AssertionComment string            `xml:"assertionComment,attr"`
	Properties       MappingProperties `xml:"L7p:Properties"`
}

type MappingProperties struct {
	MapValue string         `xml:"mapValue,attr"`
	Entry    []MappingEntry `xml:"L7p:entry"`
}

type MappingEntry struct {
	Key   PolicyStringValue `xml:"L7p:key"`
	Value PolicyStringValue `xml:"L7p:value"`
}

type PolicyStringValue struct {
	StringValue string `xml:"stringValue,attr"`
}

type GPolicy struct {
	Xml string `json:"xml"`
}

type OneOrMore struct {
	Text             string                          `xml:",chardata"`
	Usage            string                          `xml:"Usage,attr"`
	CommentAssertion *CommentAssertion               `xml:"L7p:CommentAssertion"`
	Authentication   AuthenticateAgainstIdpAssertion `xml:"L7p:Authentication"`
	Encapsulated     EncapsulatedAssertion           `xml:"L7p:Encapsulated"`
}

type AuthenticateAgainstIdpAssertion struct {
	IdentityProviderOid IdentityProviderOid `xml:"L7p:IdentityProviderOid"`
}

type IdentityProviderOid struct {
	GoidValue string `xml:"goidValue,attr"`
}

type EncapsulatedAssertion struct {
	AssertionComment                *AssertionComment `xml:"L7p:AssertionComment"`
	EncapsulatedAssertionConfigGuid PolicyStringValue `xml:"L7p:EncapsulatedAssertionConfigGuid"`
	EncapsulatedAssertionConfigName PolicyStringValue `xml:"L7p:EncapsulatedAssertionConfigName"`
	Parameters                      MappingProperties `xml:"L7p:Parameters"`
}

func BuildLayer7PolicyXml(name string, gatewayHost string, fipId string) ([]byte, error) {
	switch name {
	case "#OTK Client Context Variables":
		policy := PolicyXml{
			XMLNSL7p: "http://www.layer7tech.com/ws/policy",
			XMLNSWsp: "http://schemas.xmlsoap.org/ws/2002/12/policy",
			All: PolicyBody{
				Usage: "Required",
				SetVariable: &[]SetVariableAssertion{
					{
						VariableToSet: PolicyStringValue{
							StringValue: "host_oauth2_auth_server",
						},
						Base64Expression: PolicyStringValue{
							StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
						},
					},
					{
						VariableToSet: PolicyStringValue{
							StringValue: "audience_recipient_restriction",
						},
						Base64Expression: PolicyStringValue{
							StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
						},
					},
				},
			},
		}
		policyBytes, err := xml.Marshal(policy)
		if err != nil {
			return nil, err
		}

		return policyBytes, nil
	case "OTK FIP Client Authentication Extension":
		policy := PolicyXml{
			XMLNSL7p: "http://www.layer7tech.com/ws/policy",
			XMLNSWsp: "http://schemas.xmlsoap.org/ws/2002/12/policy",
			All: PolicyBody{
				Usage: "Required",
				OneOrMore: &OneOrMore{
					Authentication: AuthenticateAgainstIdpAssertion{
						IdentityProviderOid: IdentityProviderOid{
							GoidValue: "41e5cacd15f86758f03ff2952616d4f3",
						},
					},
					Encapsulated: EncapsulatedAssertion{
						EncapsulatedAssertionConfigGuid: PolicyStringValue{
							StringValue: "56bd8147-3ab4-4d09-9460-8b2de02b7a9e",
						},
						EncapsulatedAssertionConfigName: PolicyStringValue{
							StringValue: "OTK Fail with error message",
						},
						Parameters: MappingProperties{
							MapValue: "included",
							Entry: []MappingEntry{
								{
									Key: PolicyStringValue{
										StringValue: "apiPrefix",
									},
									Value: PolicyStringValue{
										StringValue: "${apiPrefix}",
									},
								},
								{
									Key: PolicyStringValue{
										StringValue: "givenErrorCode",
									},
									Value: PolicyStringValue{
										StringValue: "205",
									},
								},
							},
						},
					},
					Usage: "Required",
				},
			},
		}
		policyBytes, err := xml.Marshal(policy)
		if err != nil {
			return nil, err
		}

		return policyBytes, nil

	case "#OTK OVP Configuration":
		policy := PolicyXml{
			XMLNSL7p: "http://www.layer7tech.com/ws/policy",
			XMLNSWsp: "http://schemas.xmlsoap.org/ws/2002/12/policy",
			All: PolicyBody{
				Usage: "Required",
				SetVariable: &[]SetVariableAssertion{{
					VariableToSet: PolicyStringValue{
						StringValue: "host_oauth_ovp_server",
					},
					Base64Expression: PolicyStringValue{
						StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
					}},
				},
			},
		}
		policyBytes, err := xml.Marshal(policy)
		if err != nil {
			return nil, err
		}

		return policyBytes, nil

	case "#OTK Storage Configuration":
		policy := PolicyXml{
			XMLNSL7p: "http://www.layer7tech.com/ws/policy",
			XMLNSWsp: "http://schemas.xmlsoap.org/ws/2002/12/policy",
			All: PolicyBody{
				Usage: "Required",
				SetVariable: &[]SetVariableAssertion{
					{
						VariableToSet: PolicyStringValue{
							StringValue: "host_oauth_tokenstore_server",
						},
						Base64Expression: PolicyStringValue{
							StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
						},
					},
					{
						VariableToSet: PolicyStringValue{
							StringValue: "host_oauth_clientstore_server",
						},
						Base64Expression: PolicyStringValue{
							StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
						},
					},
					{
						VariableToSet: PolicyStringValue{
							StringValue: "host_oauth_session_server",
						},
						Base64Expression: PolicyStringValue{
							StringValue: base64.StdEncoding.EncodeToString([]byte(gatewayHost)),
						},
					},
				},
			},
		}
		policyBytes, err := xml.Marshal(policy)
		if err != nil {
			return nil, err
		}

		return policyBytes, nil
	}
	return nil, nil
}

func BuildCWPBundle(cwps []securityv1.Property) ([]byte, string, error) {
	bundle := graphman.Bundle{}

	for _, cwp := range cwps {
		bundle.ClusterProperties = append(bundle.ClusterProperties, &graphman.ClusterPropertyInput{
			Name:  cwp.Name,
			Value: cwp.Value,
		})
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	return bundleBytes, sha1Sum, nil
}

func BuildDefaultListenPortBundle(refreshOnKeyChanges bool) ([]byte, string, error) {
	bundle := graphman.Bundle{}

	cipherSuites := []string{
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
		"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
		"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_AES_256_GCM_SHA384",
		"TLS_AES_128_GCM_SHA256",
	}

	tlsVersions := []string{"TLSv1.2", "TLSv1.3"}

	httpPort := &graphman.ListenPortInput{
		Name:     "Default HTTP (8080)",
		Enabled:  false,
		Protocol: "HTTP",
		Port:     8080,
		EnabledFeatures: []graphman.ListenPortFeature{
			"PUBLISHED_SERVICE_MESSAGE_INPUT",
		},
	}

	httpsPort := &graphman.ListenPortInput{
		Name:     "Default HTTPS (8443)",
		Enabled:  true,
		Protocol: "HTTPS",
		Port:     8443,
		EnabledFeatures: []graphman.ListenPortFeature{
			"PUBLISHED_SERVICE_MESSAGE_INPUT",
		},
		TlsSettings: &graphman.ListenPortTlsSettingsInput{
			ClientAuthentication: graphman.ListenPortClientAuthOptional,
			TlsVersions:          tlsVersions,
			CipherSuites:         cipherSuites,
			UseCipherSuitesOrder: true,
		},
	}

	managementPort := &graphman.ListenPortInput{
		Name:     "Default HTTPS (9443)",
		Enabled:  true,
		Protocol: "HTTPS",
		Port:     9443,
		EnabledFeatures: []graphman.ListenPortFeature{
			"PUBLISHED_SERVICE_MESSAGE_INPUT",
			"ADMINISTRATIVE_ACCESS",
			"BROWSER_BASED_ADMINISTRATION",
			"BUILT_IN_SERVICES",
		},
		TlsSettings: &graphman.ListenPortTlsSettingsInput{
			ClientAuthentication: graphman.ListenPortClientAuthOptional,
			TlsVersions:          tlsVersions,
			CipherSuites:         cipherSuites,
			UseCipherSuitesOrder: true,
		},
	}

	if refreshOnKeyChanges {
		refreshOnKeyChangesProp := &graphman.EntityPropertyInput{
			Name:  "refreshOnKeyChanges",
			Value: "true",
		}
		httpsPort.Properties = append(httpsPort.Properties, refreshOnKeyChangesProp)
		managementPort.Properties = append(managementPort.Properties, refreshOnKeyChangesProp)
	}

	bundle.ListenPorts = append(bundle.ListenPorts, httpPort)
	bundle.ListenPorts = append(bundle.ListenPorts, httpsPort)
	bundle.ListenPorts = append(bundle.ListenPorts, managementPort)

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return bundleBytes, sha1Sum, nil
}

func BuildCustomListenPortBundle(gw *securityv1.Gateway, refreshOnKeyChanges bool) ([]byte, string, error) {
	bundle := graphman.Bundle{}
	//privateKey := "00000000000000000000000000000002:ssl"
	clientAuthentication := graphman.ListenPortClientAuthOptional
	for _, port := range gw.Spec.App.ListenPorts.Ports {
		enabledFeatures := []graphman.ListenPortFeature{}
		for i := range port.ManagementFeatures {
			managementFeature := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(port.ManagementFeatures[i], " ", "_"), "-", "_"))
			enabledFeatures = append(enabledFeatures, graphman.ListenPortFeature(managementFeature))
		}

		newPort := graphman.ListenPortInput{
			Name:            port.Name,
			Enabled:         port.Enabled,
			Protocol:        port.Protocol,
			Port:            port.Port,
			EnabledFeatures: enabledFeatures,
		}

		if port.Tls.Enabled {
			if port.Tls.ClientAuthentication != "" {
				clientAuthentication = graphman.ListenPortClientAuth(strings.ToUpper(port.Tls.ClientAuthentication))
			}

			newPort.TlsSettings = &graphman.ListenPortTlsSettingsInput{
				ClientAuthentication: clientAuthentication,
				CipherSuites:         port.Tls.CipherSuites,
				UseCipherSuitesOrder: port.Tls.UseCipherSuitesOrder,
				TlsVersions:          port.Tls.Versions,
				// KeystoreId:           strings.Split(privateKey, ":")[0],
				// KeyAlias:             strings.Split(privateKey, ":")[1],
			}

			hasRefreshOnKeyChangeProp := false
			for _, prop := range port.Properties {
				if prop.Name == "refreshOnKeyChanges" {
					hasRefreshOnKeyChangeProp = true
				}
				newPort.Properties = append(newPort.Properties, &graphman.EntityPropertyInput{
					Name:  prop.Name,
					Value: prop.Value,
				})
			}

			if refreshOnKeyChanges && !hasRefreshOnKeyChangeProp {
				refreshOnKeyChangesProp := &graphman.EntityPropertyInput{
					Name:  "refreshOnKeyChanges",
					Value: "true",
				}
				newPort.Properties = append(newPort.Properties, refreshOnKeyChangesProp)
			}

			if port.Tls.PrivateKey != "" {
				newPort.TlsSettings.KeystoreId = strings.Split(port.Tls.PrivateKey, ":")[0]
				newPort.TlsSettings.KeyAlias = strings.Split(port.Tls.PrivateKey, ":")[1]
			}
		}
		bundle.ListenPorts = append(bundle.ListenPorts, &newPort)
	}
	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return bundleBytes, sha1Sum, nil
}
