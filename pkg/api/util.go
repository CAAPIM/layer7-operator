package api

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/internal/templategen"
	"strings"
)

func ConvertPortalPolicyXmlToGraphman(policyXml string, passwords []templategen.SecurePassword, passwordUndeploymentIds []string) ([]byte, string, error) {
	restmanBundle := Bundle{}
	graphmanBundle := graphman.Bundle{}
	err := xml.Unmarshal([]byte(policyXml), &restmanBundle)
	if err != nil {
		return nil, "", err
	}

	hasService := false
	hasPolicy := false
	/// convert items to graphman
	for _, item := range restmanBundle.References.Item {
		if item.Type == "POLICY" && item.Resource.Policy.PolicyDetail.PolicyType == "Include" {

			policyFragment := graphman.PolicyFragmentInput{
				FolderPath: "/Portal APIs",
				Name:       item.Name,
				Goid:       item.ID,
				Guid:       item.Resource.Policy.Guid,
				Policy: &graphman.PolicyInput{
					Xml: item.Resource.Policy.Resources.ResourceSet.Resource.Text,
				},
				Soap: false,
			}
			graphmanBundle.PolicyFragments = append(graphmanBundle.PolicyFragments, &policyFragment)
			hasPolicy = true
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
			var serviceId string
			serviceId = strings.ReplaceAll(item.ID, "-", "")
			l7Service := graphman.WebApiServiceInput{
				Name:           item.Name,
				Goid:           serviceId,
				FolderPath:     "/Portal APIs",
				ResolutionPath: item.Resource.Service.ServiceDetail.ServiceMappings.HttpMapping.UrlPattern,
				MethodsAllowed: methodsAllowed,
				Enabled:        item.Resource.Service.ServiceDetail.Enabled,
				Properties:     properties,
				Policy:         &graphman.PolicyInput{Xml: item.Resource.Service.Resources.ResourceSet.Resource.Text},
			}
			graphmanBundle.WebApiServices = append(graphmanBundle.WebApiServices, &l7Service)
			hasService = true
		}
	}

	for _, securePassword := range passwords {
		l7Secret := graphman.SecretInput{
			Name:                 securePassword.Name,
			Secret:               securePassword.Value,
			SecretType:           graphman.SecretTypePassword,
			Goid:                 securePassword.Id,
			VariableReferencable: true,
			Description:          securePassword.Description,
		}
		graphmanBundle.Secrets = append(graphmanBundle.Secrets, &l7Secret)
	}

	graphmanBundle.Properties = &graphman.BundleProperties{}

	for _, securePasswordIdsForUndeployment := range passwordUndeploymentIds {
		secretToDelete := "l7_secure_" + securePasswordIdsForUndeployment
		graphmanBundle.Secrets = append(graphmanBundle.Secrets, &graphman.SecretInput{
			Name:       secretToDelete,
			Secret:     "",
			SecretType: graphman.SecretTypePassword,
		})
		graphmanBundle.Properties.Mappings.Secrets = append(graphmanBundle.Properties.Mappings.Secrets,
			&graphman.MappingInstructionInput{
				Action: graphman.MappingActionDelete,
				Source: graphman.MappingSource{
					Name: secretToDelete,
				},
			})
	}

	if hasService {
		graphmanBundle.Properties.Mappings.WebApiServices = append(graphmanBundle.Properties.Mappings.WebApiServices,
			&graphman.MappingInstructionInput{
				Action:  graphman.MappingActionNewOrUpdate,
				Default: true,
			})
	}

	if hasPolicy {
		graphmanBundle.Properties.Mappings.PolicyFragments = append(graphmanBundle.Properties.Mappings.PolicyFragments,
			&graphman.MappingInstructionInput{
				Action:  graphman.MappingActionNewOrUpdate,
				Default: true,
			})
	}

	graphmanBundleBytes, _ := json.Marshal(graphmanBundle)

	h := sha1.New()
	h.Write(graphmanBundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	return graphmanBundleBytes, sha1Sum, nil
}
