package util

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
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
	ClusterProperty *ClusterProperty `xml:"l7:ClusterProperty,omitempty"`
	ListenPort      *ListenPort      `xml:"l7:ListenPort,omitempty"`
	Policy          *Policy          `xml:"l7:Policy,omitempty"`
	Service         *Service         `xml:"l7:Service,omitempty"`
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

type ClusterProperty struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"l7:Name"`
	Value string `xml:"l7:Value"`
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

func randToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type ListenPort struct {
	ID              string          `xml:"id,attr"`
	Name            string          `xml:"l7:Name"`
	Enabled         string          `xml:"l7:Enabled"`
	Protocol        string          `xml:"l7:Protocol"`
	Port            string          `xml:"l7:Port"`
	EnabledFeatures EnabledFeatures `xml:"l7:EnabledFeatures"`
	TlsSettings     *TlsSettings    `xml:"l7:TlsSettings"`
	Properties      Properties      `xml:"l7:Properties"`
}

type TlsSettings struct {
	ClientAuthentication string              `xml:"l7:ClientAuthentication"`
	PrivateKeyReference  PrivateKeyReference `xml:"l7:PrivateKeyReference"`
	EnabledVersions      EnabledVersions     `xml:"l7:EnabledVersions"`
	EnabledCipherSuites  EnabledCipherSuites `xml:"l7:EnabledCipherSuites"`
	UseCipherSuitesOrder bool                `xml:"l7:UseCipherSuitesOrder"`
	Properties           Properties          `xml:"l7:Properties"`
}

type PrivateKeyReference struct {
	ID          string `xml:"id,attr"`
	ResourceURI string `xml:"resourceUri,attr"`
}

type EnabledVersions struct {
	StringValue []string `xml:"l7:StringValue"`
}

type EnabledCipherSuites struct {
	StringValue []string `xml:"l7:StringValue"`
}

type EnabledFeatures struct {
	StringValue []string `xml:"l7:StringValue"`
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
	refs := References{}
	items := []Item{}
	mapping := []Mapping{}
	cwpIds := []string{}

	for _, cwp := range cwps {
		randomId, err := randToken(16)
		cwpIds = append(cwpIds, randomId)

		if err != nil {
			return nil, "", err
		}

		resource := Resource{ClusterProperty: &ClusterProperty{
			ID:    randomId,
			Name:  cwp.Name,
			Value: cwp.Value,
		}}
		items = append(items, Item{Name: cwp.Name,
			ID:       randomId,
			Type:     "CLUSTER_PROPERTY",
			Resource: resource,
		})

		properties := []Property{{
			Key:         "MapBy",
			StringValue: "name",
		}, {
			Key:         "MapTo",
			StringValue: cwp.Name,
		},
		}

		mapping = append(mapping, Mapping{
			Action:     "NewOrUpdate",
			SrcId:      randomId,
			Type:       "CLUSTER_PROPERTY",
			Properties: Properties{Property: properties},
		})

		refs.Item = items
	}

	mappings := Mappings{Mapping: mapping}

	bundle := Bundle{
		XMLNS:      "http://ns.l7tech.com/2010/04/gateway-management",
		References: refs,
		Mappings:   mappings,
	}

	bundleBytes, err := xml.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	bundleString := string(bundleBytes)
	for _, cwpId := range cwpIds {
		bundleString = strings.ReplaceAll(bundleString, cwpId, "")
	}

	h := sha1.New()
	h.Write([]byte(bundleString))
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return bundleBytes, sha1Sum, nil

}

func BuildDefaultListenPortBundle() ([]byte, string, error) {
	trafficId, _ := randToken(16)
	managementId, _ := randToken(16)
	plaintextId, _ := randToken(16)
	portIds := []string{trafficId, managementId, plaintextId}

	cipherSuites := []string{
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
		"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256",
		"TLS_DHE_RSA_WITH_AES_256_CBC_SHA",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
		"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256",
		"TLS_DHE_RSA_WITH_AES_128_CBC_SHA",
		"TLS_AES_256_GCM_SHA384",
		"TLS_AES_128_GCM_SHA256"}

	tlsVersions := []string{"TLSv1.2", "TLSv1.3"}

	refs := References{}
	items := []Item{}
	mapping := []Mapping{}

	plantextPort := Item{
		Name: "Default HTTP (8080)",
		ID:   plaintextId,
		Type: "SSG_CONNECTOR",
		Resource: Resource{
			ListenPort: &ListenPort{
				ID:       plaintextId,
				Name:     "Default HTTP (8080)",
				Enabled:  "false",
				Protocol: "HTTP",
				Port:     "8080",
				EnabledFeatures: EnabledFeatures{
					StringValue: []string{
						"Published service message input",
					}},
			},
		},
	}

	items = append(items, plantextPort)

	managementPort := Item{
		Name: "Default HTTPS (9443)",
		ID:   trafficId,
		Type: "SSG_CONNECTOR",
		Resource: Resource{
			ListenPort: &ListenPort{
				ID:       managementId,
				Name:     "Default HTTPS (9443)",
				Enabled:  "true",
				Protocol: "HTTPS",
				Port:     "9443",
				EnabledFeatures: EnabledFeatures{
					StringValue: []string{
						"Published service message input",
						"Administrative access",
						"Browser-based administration",
						"Built-in services",
					}},
				TlsSettings: &TlsSettings{
					ClientAuthentication: "Optional",
					PrivateKeyReference:  PrivateKeyReference{ID: "00000000000000000000000000000002:ssl", ResourceURI: "http://ns.l7tech.com/2010/04/gateway-management/privateKeys"},
					EnabledVersions: EnabledVersions{
						StringValue: tlsVersions,
					},
					EnabledCipherSuites: EnabledCipherSuites{
						StringValue: cipherSuites,
					},
					UseCipherSuitesOrder: true,
					Properties: Properties{
						Property: []Property{
							{
								Key:          "usesTLS",
								BooleanValue: true,
							},
						},
					},
				},
			},
		},
	}

	items = append(items, managementPort)

	trafficPort := Item{
		Name: "Default HTTPS (8443)",
		ID:   trafficId,
		Type: "SSG_CONNECTOR",
		Resource: Resource{
			ListenPort: &ListenPort{
				ID:       trafficId,
				Name:     "Default HTTPS (8443)",
				Enabled:  "true",
				Protocol: "HTTPS",
				Port:     "8443",
				EnabledFeatures: EnabledFeatures{
					StringValue: []string{
						"Published service message input",
					}},
				TlsSettings: &TlsSettings{
					ClientAuthentication: "Optional",
					PrivateKeyReference:  PrivateKeyReference{ID: "00000000000000000000000000000002:ssl", ResourceURI: "http://ns.l7tech.com/2010/04/gateway-management/privateKeys"},
					EnabledVersions: EnabledVersions{
						StringValue: tlsVersions,
					},
					EnabledCipherSuites: EnabledCipherSuites{
						StringValue: cipherSuites,
					},
					UseCipherSuitesOrder: true,
					Properties: Properties{
						Property: []Property{
							{
								Key:          "usesTLS",
								BooleanValue: true,
							},
						},
					},
				},
			},
		},
	}

	items = append(items, trafficPort)

	refs.Item = items

	mapping = append(mapping, Mapping{
		Action: "NewOrUpdate",
		SrcId:  plaintextId,
		Type:   "SSG_CONNECTOR",
		Properties: Properties{Property: []Property{{
			Key:         "MapBy",
			StringValue: "name",
		}, {
			Key:         "MapTo",
			StringValue: "Default HTTP (8080)",
		},
		}},
	})

	mapping = append(mapping, Mapping{
		Action: "NewOrUpdate",
		SrcId:  managementId,
		Type:   "SSG_CONNECTOR",
		Properties: Properties{Property: []Property{{
			Key:         "MapBy",
			StringValue: "name",
		}, {
			Key:         "MapTo",
			StringValue: "Default HTTPS (9443)",
		},
		}},
	})

	mapping = append(mapping, Mapping{
		Action: "NewOrUpdate",
		SrcId:  trafficId,
		Type:   "SSG_CONNECTOR",
		Properties: Properties{Property: []Property{{
			Key:         "MapBy",
			StringValue: "name",
		}, {
			Key:         "MapTo",
			StringValue: "Default HTTPS (8443)",
		},
		}},
	})

	mappings := Mappings{Mapping: mapping}

	bundle := Bundle{
		XMLNS:      "http://ns.l7tech.com/2010/04/gateway-management",
		References: refs,
		Mappings:   mappings,
	}

	bundleBytes, err := xml.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	bundleString := string(bundleBytes)
	for _, portId := range portIds {
		bundleString = strings.ReplaceAll(bundleString, portId, "")
	}

	h := sha1.New()
	h.Write([]byte(bundleString))
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return bundleBytes, sha1Sum, nil
}

func BuildCustomListenPortBundle(gw *securityv1.Gateway) ([]byte, string, error) {
	refs := References{}
	items := []Item{}
	mapping := []Mapping{}
	portIds := []string{}

	for _, port := range gw.Spec.App.ListenPorts.Ports {
		portId, _ := randToken(16)
		portIds = append(portIds, portId)
		newPort := Item{
			Name: port.Name,
			ID:   portId,
			Type: "SSG_CONNECTOR",
			Resource: Resource{
				ListenPort: &ListenPort{
					ID:       portId,
					Name:     port.Name,
					Enabled:  strconv.FormatBool(port.Enabled),
					Protocol: port.Protocol,
					Port:     port.Port,
					EnabledFeatures: EnabledFeatures{
						StringValue: port.ManagementFeatures,
					},
				},
			},
		}

		if port.Tls.Enabled {

			privateKey := "00000000000000000000000000000002:ssl"

			if port.Tls.PrivateKey != "" {
				privateKey = port.Tls.PrivateKey
			}
			newPort.Resource.ListenPort.TlsSettings = &TlsSettings{
				ClientAuthentication: port.Tls.ClientAuthentication,
				PrivateKeyReference:  PrivateKeyReference{ID: privateKey, ResourceURI: "http://ns.l7tech.com/2010/04/gateway-management/privateKeys"},
				EnabledVersions:      EnabledVersions{StringValue: port.Tls.Versions},
				EnabledCipherSuites:  EnabledCipherSuites{StringValue: port.Tls.CipherSuites},
				UseCipherSuitesOrder: port.Tls.UseCipherSuitesOrder,
				Properties:           Properties{Property: []Property{{Key: "usesTLS", BooleanValue: port.Tls.Enabled}}},
			}
		}

		if len(port.Properties) > 0 {
			for _, property := range port.Properties {
				newProp := Property{Key: property.Name, StringValue: property.Value}
				newPort.Resource.ListenPort.Properties.Property = append(newPort.Resource.ListenPort.Properties.Property, newProp)
			}

		}
		items = append(items, newPort)

		mapping = append(mapping, Mapping{
			Action: "NewOrUpdate",
			SrcId:  portId,
			Type:   "SSG_CONNECTOR",
			Properties: Properties{Property: []Property{{
				Key:         "MapBy",
				StringValue: "name",
			}, {
				Key:         "MapTo",
				StringValue: port.Name,
			},
			}},
		})
	}

	refs.Item = items
	mappings := Mappings{Mapping: mapping}

	bundle := Bundle{
		XMLNS:      "http://ns.l7tech.com/2010/04/gateway-management",
		References: refs,
		Mappings:   mappings,
	}

	bundleBytes, err := xml.Marshal(bundle)
	if err != nil {
		return nil, "", err
	}

	bundleString := string(bundleBytes)
	for _, portId := range portIds {
		bundleString = strings.ReplaceAll(bundleString, portId, "")
	}

	h := sha1.New()
	h.Write([]byte(bundleString))
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return bundleBytes, sha1Sum, nil
}
