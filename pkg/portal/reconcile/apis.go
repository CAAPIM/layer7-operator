package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/caapim/layer7-operator/pkg/portal"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func syncPortalApis(ctx context.Context, params Params) {

	l7Portal := &v1alpha1.L7Portal{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, l7Portal)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "portal not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-sync-portal-apis")
		return
	}

	if l7Portal.Spec.PortalManaged {
		return
	}
	portalApiSummaryConfigMap := corev1.ConfigMap{}

	err = params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name + "-api-summary", Namespace: params.Instance.Namespace}, &portalApiSummaryConfigMap)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				_ = removeJob(params.Instance.Name + "-sync-portal-apis")
				return
			}
		}
		return
	}

	portalApiSummaryBytes, err := base64.StdEncoding.DecodeString(portalApiSummaryConfigMap.Data["apis"])
	if err != nil {
		return
	}

	var portalApiSummary []templategen.PortalAPI

	err = json.Unmarshal(portalApiSummaryBytes, &portalApiSummary)
	if err != nil {
		return
	}

	for _, api := range portalApiSummary {
		policyXml := templategen.BuildTemplate(api)
		restmanBundle := portal.Bundle{}
		graphmanBundle := graphman.Bundle{}
		err = xml.Unmarshal([]byte(policyXml), &restmanBundle)
		if err != nil {
			return
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
		dataCheckSum := sha1Sum
		deletionGracePeriodSeconds := int64(5)
		desiredL7API := &v1alpha1.L7Api{
			ObjectMeta: metav1.ObjectMeta{
				Name:                       strings.ToLower(strings.ReplaceAll(api.Name, " ", "-")),
				Namespace:                  params.Instance.Namespace,
				DeletionGracePeriodSeconds: &deletionGracePeriodSeconds,
				Labels:                     util.DefaultLabels(strings.ToLower(strings.ReplaceAll(api.Name, " ", "-")), map[string]string{}),
				Annotations: map[string]string{
					"checksum/bundle": dataCheckSum,
				},
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1alpha1",
				Kind:       "L7Api",
			},
			Spec: v1alpha1.L7ApiSpec{
				ServiceUrl:      api.SsgUrl,
				PortalPublished: true,
				GraphmanBundle:  base64.StdEncoding.EncodeToString(graphmanBundleBytes),
				DeploymentTags:  params.Instance.Spec.DeploymentTags,
				PortalMeta: v1alpha1.PortalMeta{
					TenantId:       api.TenantId,
					Name:           api.Name,
					Uuid:           api.Uuid,
					UuidStripped:   api.UuidStripped,
					SsgUrlBase64:   api.SsgUrlBase64,
					SsgUrl:         api.SsgUrl,
					ServiceId:      api.ServiceId,
					ApiEnabled:     api.ApiEnabled,
					LocationUrl:    base64.StdEncoding.EncodeToString([]byte(api.LocationUrl)),
					Checksum:       dataCheckSum,
					SsgServiceType: api.SsgServiceType,
					ModifyTs:       api.ModifyTs,
				},
				L7Portal: params.Instance.Name,
			},
		}

		if err := controllerutil.SetControllerReference(params.Instance, desiredL7API, params.Scheme); err != nil {
			params.Log.Info("failed to set controller reference", "name", desiredL7API.Name, "namespace", params.Instance.Namespace, "error", err.Error())
			return
		}

		currentL7API := &v1alpha1.L7Api{}

		err = params.Client.Get(ctx, types.NamespacedName{Name: desiredL7API.Name, Namespace: params.Instance.Namespace}, currentL7API)
		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredL7API); err != nil {
				params.Log.V(2).Info("failed to create l7api", "name", desiredL7API.Name, "namespace", params.Instance.Namespace, "error", err.Error())
				return
			}
			params.Log.Info("created l7Api", "name", desiredL7API.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			params.Log.Info("failed to retrieve l7api", "name", desiredL7API.Name, "namespace", params.Instance.Namespace, "error", err.Error())
			return
		}

		updatedL7API := currentL7API.DeepCopy()
		updatedL7API.Spec = desiredL7API.Spec

		for k, v := range desiredL7API.ObjectMeta.Annotations {
			updatedL7API.ObjectMeta.Annotations[k] = v
		}

		for k, v := range desiredL7API.ObjectMeta.Labels {
			updatedL7API.ObjectMeta.Labels[k] = v
		}

		if desiredL7API.ObjectMeta.Annotations["checksum/bundle"] != currentL7API.ObjectMeta.Annotations["checksum/bundle"] || !reflect.DeepEqual(desiredL7API.Spec.DeploymentTags, currentL7API.Spec.DeploymentTags) {
			patch := client.MergeFrom(currentL7API)
			if err := params.Client.Patch(ctx, updatedL7API, patch); err != nil {
				params.Log.Info("failed to update l7Api", "name", desiredL7API.Name, "namespace", params.Instance.Namespace, "error", err.Error())
				return
			}
			params.Log.Info("l7Api updated", "name", desiredL7API.Name, "namespace", desiredL7API.Namespace)
		}
	}
}
