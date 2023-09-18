package portal

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
)

type Bundle struct {
	XMLName    xml.Name   `xml:"Bundle"`
	XMLNS      string     `xml:"xmlns:l7,attr"`
	References References `xml:"References"`
	Mappings   Mappings   `xml:"Mappings"`
}

type References struct {
	Item []Item `xml:"Item"`
}

type Item struct {
	Name     string   `xml:"Name"`
	ID       string   `xml:"Id"`
	Type     string   `xml:"Type"`
	Resource Resource `xml:"Resource"`
}

type Resource struct {
	ClusterProperty *ClusterProperty `xml:"ClusterProperty,omitempty"`
	ListenPort      *ListenPort      `xml:"ListenPort,omitempty"`
	Policy          *Policy          `xml:"Policy,omitempty"`
	Service         *Service         `xml:"Service,omitempty"`
}

type Policy struct {
	Guid         string       `xml:"guid,attr"`
	ID           string       `xml:"id,attr"`
	Version      string       `xml:"version,attr"`
	PolicyDetail PolicyDetail `xml:"PolicyDetail"`

	Resources PolicyResources
}

type PolicyResources struct {
	ResourceSet PolicyResourceSet `xml:"ResourceSet"`
}

type PolicyResourceSet struct {
	Tag      string         `xml:"tag,attr"`
	Resource PolicyResource `xml:"Resource"`
}

type PolicyResource struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type PolicyDetail struct {
	FolderId   string                 `xml:"folderId,attr"`
	Guid       string                 `xml:"guid,attr"`
	ID         string                 `xml:"id,attr"`
	Name       string                 `xml:"Name"`
	PolicyType string                 `xml:"PolicyType"`
	Properties PolicyDetailProperties `xml:"Properties"`
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
	ServiceDetail ServiceDetail    `xml:"ServiceDetail"`
	Properties    ServiceProperty  `xml:"Properties"`
	Resources     ServiceResources `xml:"Resources"`
}

type ServiceDetail struct {
	Text            string            `xml:",chardata"`
	FolderId        string            `xml:"folderId,attr"`
	ID              string            `xml:"id,attr"`
	Name            string            `xml:"Name"`
	Enabled         bool              `xml:"Enabled"`
	ServiceMappings ServiceMappings   `xml:"ServiceMappings"`
	Properties      ServiceProperties `xml:"Properties"`
}

type ServiceMappings struct {
	HttpMapping HttpMapping
}

type HttpMapping struct {
	UrlPattern string `xml:"UrlPattern"`
	Verbs      Verbs  `xml:"Verbs"`
}

type Verbs struct {
	Verb []string `xml:"Verb"`
}

type ServiceProperties struct {
	Property []ServiceProperty `xml:"Property"`
}

type ServiceProperty struct {
	Key          string `xml:"key,attr"`
	BooleanValue string `xml:"BooleanValue"`
	StringValue  string `xml:"StringValue"`
}

type ServiceResources struct {
	ResourceSet ServiceResourceSet `xml:"ResourceSet"`
}

type ServiceResourceSet struct {
	Tag      string          `xml:"tag,attr"`
	Resource ServiceResource `xml:"Resource"`
}

type ServiceResource struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

type ClusterProperty struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}

type Mappings struct {
	Mapping []Mapping `xml:"Mapping"`
}
type Mapping struct {
	Action     string     `xml:"action,attr"`
	SrcId      string     `xml:"srcId,attr"`
	Type       string     `xml:"type,attr"`
	Properties Properties `xml:"Properties"`
}

type Properties struct {
	Property []Property `xml:"Property"`
}

type Property struct {
	Key          string `xml:"key,attr"`
	StringValue  string `xml:"StringValue,omitempty"`
	BooleanValue bool   `xml:"BooleanValue,omitempty"`
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
	Name            string          `xml:"Name"`
	Enabled         string          `xml:"Enabled"`
	Protocol        string          `xml:"Protocol"`
	Port            string          `xml:"Port"`
	EnabledFeatures EnabledFeatures `xml:"EnabledFeatures"`
	TlsSettings     *TlsSettings    `xml:"TlsSettings"`
	Properties      Properties      `xml:"Properties"`
}

type TlsSettings struct {
	ClientAuthentication string              `xml:"ClientAuthentication"`
	PrivateKeyReference  PrivateKeyReference `xml:"PrivateKeyReference"`
	EnabledVersions      EnabledVersions     `xml:"EnabledVersions"`
	EnabledCipherSuites  EnabledCipherSuites `xml:"EnabledCipherSuites"`
	UseCipherSuitesOrder bool                `xml:"UseCipherSuitesOrder"`
	Properties           Properties          `xml:"Properties"`
}

type PrivateKeyReference struct {
	ID          string `xml:"id,attr"`
	ResourceURI string `xml:"resourceUri,attr"`
}

type EnabledVersions struct {
	StringValue []string `xml:"StringValue"`
}

type EnabledCipherSuites struct {
	StringValue []string `xml:"StringValue"`
}

type EnabledFeatures struct {
	StringValue []string `xml:"StringValue"`
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
