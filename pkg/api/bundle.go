package api

import (
	"encoding/xml"
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
	Policy  *Policy  `xml:"Policy,omitempty"`
	Service *Service `xml:"Service,omitempty"`
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
