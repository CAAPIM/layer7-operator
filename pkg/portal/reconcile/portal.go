package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type RawPortalAPISummary struct {
	APIs []templategen.PortalAPI `json:"results"`
}

// type PortalAPIShort struct {
// 	Name     string `json:"name"`
// 	Uuid     string `json:"uuid"`
// 	SsgUrl   string `json:"ssgUrl"`
// 	Checksum string `json:"checksum"`
// }

func syncPortal(ctx context.Context, params Params) {

	portal := &v1alpha1.L7Portal{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, portal)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Info("portal not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-sync-portal")
		return
	}

	if !portal.Spec.Enabled {
		_ = removeJob(params.Instance.Name + "-sync-portal")
		return
	}

	requestCacheEntry := portal.Name + "-access-token"
	syncRequest, _ := syncCache.Read(requestCacheEntry)

	if syncRequest.Attempts > 0 {
		params.Log.V(2).Info("request has failed in the last 30 seconds, backing off", "name", portal.Name, "namespace", portal.Namespace)
		return
	}

	params.Instance = portal

	token, err := util.GetPortalAccessToken(params.Instance.Spec.Name, params.Instance.Spec.Auth.Endpoint, params.Instance.Spec.Auth.PapiClientId, params.Instance.Spec.Auth.PapiClientSecret)
	if err != nil {
		params.Log.Info("failed to retrieve portal access token", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
		return
	}

	// TODO:
	// Refactor when deleted entities are available
	// Should only retrieve changes after a last modified date to reduce resource utilisation
	requestCacheEntry = portal.Name + "-api-summary"
	syncRequest, _ = syncCache.Read(requestCacheEntry)

	if syncRequest.Attempts > 0 {
		params.Log.V(2).Info("request has failed in the last 30 seconds, backing off", "name", portal.Name, "namespace", portal.Namespace)
		return
	}

	//modifyTs := time.UnixMilli(params.Instance.Status.LastUpdated)

	apiEndpoint := "https://" + params.Instance.Spec.Endpoint + ":443/" + params.Instance.Spec.Name + "/api-management/layer7-operator/0.1/apis?size=2000"

	// Get summary
	resp, err := util.RestCall("GET", apiEndpoint, true, map[string]string{"Authorization": "Bearer " + token}, "application/json;charset=utf-8", []byte{}, "", "")
	if err != nil {
		params.Log.V(2).Info("failed to retrieve portal api summary", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
		return
	}

	var portalAPISummary []templategen.PortalAPI

	err = json.Unmarshal(resp, &portalAPISummary)
	if err != nil {
		params.Log.Info("failed to unmarshal portal api summary", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		return
	}

	var portalAPIList []templategen.PortalAPI

	for _, api := range portalAPISummary {

		dataBytes, _ := json.Marshal(api)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		dataCheckSum := sha1Sum

		portalAPIList = append(portalAPIList, templategen.PortalAPI{
			Name:            api.Name,
			Uuid:            api.Uuid,
			SsgUrlBase64:    api.SsgUrlBase64,
			SsgUrl:          api.SsgUrl,
			ServiceId:       api.ServiceId,
			ApiEnabled:      api.ApiEnabled,
			LocationUrl:     base64.StdEncoding.EncodeToString([]byte(api.LocationUrl)),
			PolicyTemplates: api.PolicyTemplates,
			CustomFields:    api.CustomFields,
			Checksum:        dataCheckSum,
			SsgServiceType:  api.SsgServiceType,
			ModifyTs:        api.ModifyTs,
		})
	}
	var currentPortalAPIList []templategen.PortalAPI
	/// look up configmap and check if an API has been removed.. then schedule deletion
	currentSummary, err := getConfigmap(ctx, params, params.Instance.Name+"-api-summary")

	if err == nil {
		currentPortalApiSummaryBytes, err := base64.StdEncoding.DecodeString(currentSummary.Data["apis"])
		if err != nil {
			params.Log.Info("failed to decode portal api summary", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
			return
		}

		err = json.Unmarshal(currentPortalApiSummaryBytes, &currentPortalAPIList)
		if err != nil {
			params.Log.Info("failed to unmarshal portal api summary", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
			return
		}
	} else {
		params.Log.V(2).Info("failed to retrieve configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	// TODO:
	// Refactor when deleted entities are available
	apiRemovalList := []string{}
	if len(currentPortalAPIList) > len(portalAPIList) {

		for _, currentApi := range currentPortalAPIList {
			found := false
			for _, latestApi := range portalAPIList {
				if currentApi.Name == latestApi.Name {
					found = true
				}
			}
			if !found {
				apiRemovalList = append(apiRemovalList, currentApi.Name)
			}
		}
	}

	portalAPISummaryBytes, _ := json.Marshal(portalAPIList)
	err = ConfigMap(ctx, params, portalAPISummaryBytes)

	if err != nil {
		params.Log.V(2).Info("failed to reconcile configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return
	}

	for _, api := range apiRemovalList {
		l7Api := &v1alpha1.L7Api{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: strings.ToLower(strings.ReplaceAll(api, " ", "-")), Namespace: params.Instance.Namespace}, l7Api)
		if err != nil {
			params.Log.Info("failed to retrieve l7Api", "name", params.Instance.Name, "l7api", strings.ToLower(strings.ReplaceAll(api, " ", "-")), "namespace", params.Instance.Namespace)
			return
		}
		err = params.Client.Delete(ctx, l7Api)
		if err != nil {
			params.Log.Info("failed to remove l7Api", "name", params.Instance.Name, "l7api", strings.ToLower(strings.ReplaceAll(api, " ", "-")), "namespace", params.Instance.Namespace)
			return
		}
	}
}

func getConfigmap(ctx context.Context, params Params, name string) (*corev1.ConfigMap, error) {
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
