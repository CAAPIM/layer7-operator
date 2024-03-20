package api

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/caapim/layer7-operator/internal/graphman"
)

func ConvertPortalPolicyXmlToGraphman(policyXml string) ([]byte, string, error) {
	restmanBundle := Bundle{}
	graphmanBundle := graphman.Bundle{}
	err := xml.Unmarshal([]byte(policyXml), &restmanBundle)
	if err != nil {
		return nil, "", err
	}

	/// convert items to graphman
	for _, item := range restmanBundle.References.Item {
		if item.Type == "POLICY" && item.Resource.Policy.PolicyDetail.PolicyType == "Include" {

			policyFragment := graphman.PolicyFragmentInput{
				FolderPath: "/Portal APIs",
				Name:       item.Name,
				Guid:       item.Resource.Policy.Guid,
				Policy: &graphman.PolicyInput{
					Xml: item.Resource.Policy.Resources.ResourceSet.Resource.Text,
				},
				Soap: false,
			}
			graphmanBundle.PolicyFragments = append(graphmanBundle.PolicyFragments, &policyFragment)
		}
		if item.Type == "SERVICE" {

			methodsAllowed := []graphman.HttpMethod{}

			for _, verb := range item.Resource.Service.ServiceDetail.ServiceMappings.HttpMapping.Verbs.Verb {
				var method graphman.HttpMethod
				switch verb {
				case "GET":
					method = graphman.HttpMethodGet
				case "PUT":
					method = graphman.HttpMethodPut
				case "OPTIONS":
					method = graphman.HttpMethodOptions
				case "POST":
					method = graphman.HttpMethodPost
				case "DELETE":
					method = graphman.HttpMethodDelete
				case "HEAD":
					method = graphman.HttpMethodHead
				case "PATCH":
					method = graphman.HttpMethodPatch
				case "OTHER":
					method = graphman.HttpMethodOther
				}
				methodsAllowed = append(methodsAllowed, method)
			}

			properties := []*graphman.EntityPropertyInput{}

			for _, p := range item.Resource.Service.ServiceDetail.Properties.Property {
				graphmanEntityProperty := graphman.EntityPropertyInput{Name: p.Key}
				if p.BooleanValue != "" {
					graphmanEntityProperty.Value = p.BooleanValue
				}
				if p.StringValue != "" {
					graphmanEntityProperty.Value = p.StringValue
				}
				properties = append(properties, &graphmanEntityProperty)

			}

			l7Service := graphman.WebApiServiceInput{
				Name:           item.Name,
				FolderPath:     "/Portal APIs",
				ResolutionPath: item.Resource.Service.ServiceDetail.ServiceMappings.HttpMapping.UrlPattern,
				MethodsAllowed: methodsAllowed,
				Enabled:        item.Resource.Service.ServiceDetail.Enabled,
				Properties:     properties,
				Policy:         &graphman.PolicyInput{Xml: item.Resource.Service.Resources.ResourceSet.Resource.Text},
			}
			graphmanBundle.WebApiServices = append(graphmanBundle.WebApiServices, &l7Service)
		}
	}

	graphmanBundleBytes, _ := json.Marshal(graphmanBundle)

	h := sha1.New()
	h.Write(graphmanBundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	return graphmanBundleBytes, sha1Sum, nil

}
