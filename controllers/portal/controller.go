/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package portal

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/portal"
	"github.com/caapim/layer7-operator/pkg/portal/reconcile"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

/// TODO: refactor

// L7PortalReconciler reconciles a Gateway object
type L7PortalReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

type RawPortalAPISummary struct {
	APIs []PortalAPI `json:"results"`
}

type PortalAPI struct {
	Uuid              string   `json:"uuid"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Type              string   `json:"type"`
	PortalStatus      string   `json:"portalStatus"`
	AccessStatus      string   `json:"accessStatus"`
	SsgUrl            string   `json:"ssgUrl"`
	Version           string   `json:"version"`
	ApiEulaUuid       string   `json:"apiEulaUuid"`
	CreateTs          int      `json:"createTs"`
	ModifyTs          int      `json:"modifyTs"`
	SsgServiceType    string   `json:"ssgServiceType"`
	ApplicationUsage  int      `json:"applicationUsage"`
	Tags              []string `json:"tags"`
	PublishedByPortal bool     `json:"publishedByPortal"`
}

type PortalAPIShort struct {
	Name     string `json:"name"`
	Uuid     string `json:"uuid"`
	SsgUrl   string `json:"ssgUrl"`
	Checksum string `json:"checksum"`
}

const apiFinalizer = "security.brcmlabs.com/finalizer"

func (r *L7PortalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("portal", req.NamespacedName)

	portal := &securityv1alpha1.L7Portal{}
	err := r.Get(ctx, req.NamespacedName, portal)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: portal,
	}

	err = reconcilePortal(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcilePortalAPIs(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func reconcilePortal(ctx context.Context, params reconcile.Params) error {
	token, err := util.GetPortalAccessToken(params.Instance.Spec.Name, params.Instance.Spec.Auth.Endpoint, params.Instance.Spec.Auth.PapiClientId, params.Instance.Spec.Auth.PapiClientSecret)
	if err != nil {
		return err
	}

	apiEndpoint := "https://" + params.Instance.Spec.Endpoint + ":443/" + params.Instance.Spec.Name + "/api-management/1.0/apis?size=2000"

	// Get summary
	resp, err := util.RestCall("GET", apiEndpoint, true, map[string]string{"Authorization": "Bearer " + token}, "application/json;charset=utf-8", []byte{}, "", "")
	if err != nil {
		return err
	}

	var portalAPISummary RawPortalAPISummary

	err = json.Unmarshal(resp, &portalAPISummary)
	if err != nil {
		return err
	}

	var portalAPISummaryShort []PortalAPIShort

	for _, api := range portalAPISummary.APIs {

		dataBytes, _ := json.Marshal(api)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		dataCheckSum := sha1Sum

		portalAPISummaryShort = append(portalAPISummaryShort, PortalAPIShort{Name: api.Name, Checksum: dataCheckSum, SsgUrl: api.SsgUrl, Uuid: api.Uuid})
	}
	currentPortalAPISummaryShort := []PortalAPIShort{}
	/// look up configmap and check if an API has been removed.. then schedule deletion
	currentShortSummary, err := getShortSummaryCm(ctx, params, params.Instance.Name+"-api-summary")

	if err == nil {
		currentPortalApiSummaryBytes, err := base64.StdEncoding.DecodeString(currentShortSummary.Data["apiSummary"])
		if err != nil {
			return err
		}

		err = json.Unmarshal(currentPortalApiSummaryBytes, &currentPortalAPISummaryShort)
		if err != nil {
			return err
		}
	} else {
		params.Log.Info("error retrieving configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	apiRemovalList := []string{}
	if len(currentPortalAPISummaryShort) > len(portalAPISummaryShort) {

		for _, currentApi := range currentPortalAPISummaryShort {
			found := false
			for _, latestApi := range portalAPISummaryShort {
				if currentApi.Name == latestApi.Name {
					found = true
				}
			}
			if !found {
				apiRemovalList = append(apiRemovalList, currentApi.Name)
			}
		}
	}

	portalAPISummaryShortBytes, _ := json.Marshal(portalAPISummaryShort)

	updated, err := reconcile.ConfigMap(ctx, params, portalAPISummaryShortBytes)

	if err != nil {
		return err
	}

	for _, api := range apiRemovalList {
		l7Api := &securityv1alpha1.L7Api{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: strings.ToLower(strings.ReplaceAll(api, " ", "-")), Namespace: params.Instance.Namespace}, l7Api)
		if err != nil {
			return err
		}
		err = params.Client.Delete(ctx, l7Api)
		if err != nil {
			return err
		}
	}

	if updated {
		err = reconcilePortalStatus(ctx, params, len(portalAPISummaryShort))
		if err != nil {
			return err
		}
	}

	return nil
}

func reconcilePortalAPIs(ctx context.Context, params reconcile.Params) error {
	portalApiSummaryConfigMap := corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name + "-api-summary", Namespace: params.Instance.Namespace}, &portalApiSummaryConfigMap)
	if err != nil {
		return err
	}

	portalApiSummaryBytes, err := base64.StdEncoding.DecodeString(portalApiSummaryConfigMap.Data["apiSummary"])
	if err != nil {
		return err
	}

	var portalAPIShort []PortalAPIShort

	err = json.Unmarshal(portalApiSummaryBytes, &portalAPIShort)
	if err != nil {
		return err
	}

	token, err := util.GetPortalAccessToken(params.Instance.Spec.Name, params.Instance.Spec.Auth.Endpoint, params.Instance.Spec.Auth.PapiClientId, params.Instance.Spec.Auth.PapiClientSecret)
	if err != nil {
		return err
	}

	for _, api := range portalAPIShort {

		apiEndpoint := "https://" + params.Instance.Spec.Endpoint + "/" + params.Instance.Spec.Name + "/api-management/1.0/apis/" + api.Uuid + "/bundle"
		var resp []byte

		resp, err = util.RestCall("GET", apiEndpoint, true, map[string]string{"Authorization": "Bearer " + token}, "application/json;charset=utf-8", []byte{}, "", "")
		if err != nil {
			return err
		}

		restmanBundle := portal.Bundle{}
		graphmanBundle := graphman.Bundle{}
		err = xml.Unmarshal(resp, &restmanBundle)
		if err != nil {
			return err
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
		desiredL7API := &securityv1alpha1.L7Api{
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
			Spec: securityv1alpha1.L7ApiSpec{
				Name:            api.Name,
				ServiceUrl:      api.SsgUrl,
				PortalPublished: true,
				GraphmanBundle:  base64.StdEncoding.EncodeToString(graphmanBundleBytes),
				DeploymentTags:  params.Instance.Spec.DeploymentTags,
				L7Portal:        params.Instance.Name,
			},
		}

		if err := controllerutil.SetControllerReference(params.Instance, desiredL7API, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}

		if !controllerutil.ContainsFinalizer(desiredL7API, apiFinalizer) {
			controllerutil.AddFinalizer(desiredL7API, apiFinalizer)
		}

		currentL7API := &securityv1alpha1.L7Api{}

		err = params.Client.Get(ctx, types.NamespacedName{Name: desiredL7API.Name, Namespace: params.Instance.Namespace}, currentL7API)
		if err != nil && k8serrors.IsNotFound(err) {
			if err = params.Client.Create(ctx, desiredL7API); err != nil {
				return err
			}
			params.Log.Info("created l7Api", "name", desiredL7API.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
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
				return err
			}
			params.Log.Info("l7Api updated", "name", desiredL7API.Name, "namespace", desiredL7API.Namespace)
		}
	}

	return nil
}

func reconcilePortalStatus(ctx context.Context, params reconcile.Params, apiCount int) error {
	portalStatus := params.Instance.Status
	portalStatus.ApiSummaryConfigMap = params.Instance.Name + "-api-summary"
	portalStatus.Updated = time.Now().String()
	portalStatus.ApiCount = apiCount

	if !reflect.DeepEqual(portalStatus, params.Instance.Status) {
		params.Instance.Status = portalStatus
		err := params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.Info("failed to update portal status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
	}
	return nil
}

func getShortSummaryCm(ctx context.Context, params reconcile.Params, name string) (*corev1.ConfigMap, error) {
	shortSummary := &corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, shortSummary)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return shortSummary, err
			}
		}
	}
	return shortSummary, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *L7PortalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.L7Portal{}).
		Complete(r)
}
