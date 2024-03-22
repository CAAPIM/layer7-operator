package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

//////
//////  REFACTOR
//////

const tempDirectoryBase = "/tmp/portalapis/"

type RawPortalAPISummary struct {
	APIs []templategen.PortalAPI `json:"results"`
}

func syncPortal(ctx context.Context, params Params) {

	l7Portal := &v1alpha1.L7Portal{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, l7Portal)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Info("portal not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-sync-portal")
		return
	}

	if !l7Portal.Spec.Enabled {
		_ = removeJob(params.Instance.Name + "-sync-portal")
		return
	}

	if l7Portal.Spec.PortalManaged {
		externalManaged(params, ctx, l7Portal)
	} else {
		locallyManaged(params, ctx, l7Portal)
	}

}

func getConfigmap(ctx context.Context, params Params, name string) (*corev1.ConfigMap, error) {
	shortSummary := &corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, shortSummary)
	if err != nil {
		return shortSummary, err
	}
	return shortSummary, nil
}

// Deprecate the local approach in the future
// should only be managed by Portal
func locallyManaged(params Params, ctx context.Context, l7Portal *v1alpha1.L7Portal) {
	requestCacheEntry := l7Portal.Name + "-access-token"
	syncRequest, _ := syncCache.Read(requestCacheEntry)

	if syncRequest.Attempts > 0 {
		params.Log.V(2).Info("request has failed in the last 30 seconds, backing off", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	token, err := util.GetPortalAccessToken(l7Portal.Spec.PortalTenant, l7Portal.Spec.Auth.Endpoint, l7Portal.Spec.Auth.PapiClientId, l7Portal.Spec.Auth.PapiClientSecret)
	if err != nil {
		params.Log.Info("failed to retrieve portal access token", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
		return
	}

	// TODO:
	// Refactor when deleted entities are available
	// Should only retrieve changes after a last modified date to reduce resource utilisation
	requestCacheEntry = l7Portal.Name + "-api-summary"
	syncRequest, _ = syncCache.Read(requestCacheEntry)

	if syncRequest.Attempts > 0 {
		params.Log.V(2).Info("request has failed in the last 30 seconds, backing off", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	//modifyTs := time.UnixMilli(params.Instance.Status.LastUpdated)

	var apiEndpoint string
	if strings.Contains(l7Portal.Spec.Endpoint, ":") {
		apiEndpoint = "https://" + l7Portal.Spec.Endpoint + "/" + l7Portal.Spec.PortalTenant + "/api-management/layer7-operator/0.1/apis?size=2000"
	} else {
		apiEndpoint = "https://" + l7Portal.Spec.Endpoint + ":443/" + l7Portal.Spec.PortalTenant + "/api-management/layer7-operator/0.1/apis?size=2000"
	}

	// @Todo - Use modifyTs query param in above query so we only get the changes since the last reconciliation attempt.
	// 		   This will allow us to skip the costly step of iterating over the full set of APIs in order to compare checksums.
	//

	// Get summary
	resp, err := util.RestCall("GET", apiEndpoint, true, map[string]string{"Authorization": "Bearer " + token}, "application/json;charset=utf-8", []byte{}, "", "")
	if err != nil {
		params.Log.Info("Failed to retrieve portal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace, "endpoint", apiEndpoint)
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
		return
	}
	params.Log.V(2).Info("Successfully retrieved portal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace, "endpoint", apiEndpoint)

	var portalAPISummary []templategen.PortalAPI

	err = json.Unmarshal(resp, &portalAPISummary)
	if err != nil {
		params.Log.Info("failed to unmarshal portal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace, "error", err.Error())
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
			TenantId:        api.TenantId,
			Name:            api.Name,
			Uuid:            api.Uuid,
			UuidStripped:    api.UuidStripped,
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

	params.Log.Info("Creating config map")
	// DMUN : ConfigMap lookup was failing as it did not exist,  creating at
	portalAPISummaryBytes, _ := json.Marshal(portalAPIList)
	err = ConfigMap(ctx, params, portalAPISummaryBytes)

	currentSummary, err := getConfigmap(ctx, params, l7Portal.Name+"-api-summary")

	if err == nil {
		currentPortalApiSummaryBytes, err := base64.StdEncoding.DecodeString(currentSummary.Data["apis"])
		if err != nil {
			params.Log.Info("failed to decode portal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
			return
		}

		err = json.Unmarshal(currentPortalApiSummaryBytes, &currentPortalAPIList)
		if err != nil {
			params.Log.Info("failed to unmarshal portal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
			return
		}
	} else {
		params.Log.V(2).Info("failed to retrieve configmap", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
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

	// DMUn config map
	//portalAPISummaryBytes, _ := json.Marshal(portalAPIList)
	//err = ConfigMap(ctx, params, portalAPISummaryBytes)

	if err != nil {
		params.Log.V(2).Info("failed to reconcile configmap", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	for _, api := range apiRemovalList {
		l7Api := &v1alpha1.L7Api{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: strings.ToLower(strings.ReplaceAll(api, " ", "-")), Namespace: l7Portal.Namespace}, l7Api)
		if err != nil {
			params.Log.Info("failed to retrieve l7Api", "name", l7Portal.Name, "l7api", strings.ToLower(strings.ReplaceAll(api, " ", "-")), "namespace", l7Portal.Namespace)
			return
		}
		err = params.Client.Delete(ctx, l7Api)
		if err != nil {
			params.Log.Info("failed to remove l7Api", "name", l7Portal.Name, "l7api", strings.ToLower(strings.ReplaceAll(api, " ", "-")), "namespace", l7Portal.Namespace)
			return
		}
	}
}

func externalManaged(params Params, ctx context.Context, l7Portal *v1alpha1.L7Portal) {
	portalAPIs := []templategen.PortalAPI{}
	portalTempDirectory := tempDirectoryBase + l7Portal.Name

	folderInfo, err := os.Stat(portalTempDirectory)
	if err != nil {
		params.Log.Info("failed to scan temp storage", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	_, err = getConfigmap(ctx, params, l7Portal.Name+"-api-summary")

	if folderInfo.ModTime().Add(30*time.Second).Before(time.Now()) && err == nil {
		return
	}

	dInfo, err := os.ReadDir(portalTempDirectory)
	if err != nil {
		params.Log.V(2).Info("failed to read temp storage directory", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	for _, f := range dInfo {

		fBytes, err := os.ReadFile(portalTempDirectory + "/" + f.Name())
		if err != nil {
			params.Log.V(2).Info("failed to read portal api from temp storage", "name", l7Portal.Name, "summary file", f.Name(), "namespace", l7Portal.Namespace)
		}

		portalAPI := templategen.PortalAPI{}
		err = json.Unmarshal(fBytes, &portalAPI)
		if err != nil {
			params.Log.Info("failed to unmarshal portal api summary", "name", l7Portal.Name, "summary file", f.Name(), "namespace", l7Portal.Namespace)
			continue
		}
		portalAPIs = append(portalAPIs, portalAPI)
	}

	portalApiBytes, err := json.Marshal(portalAPIs)
	if err != nil {
		params.Log.Info("failed to read api from temp storage directory", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}
	err = ConfigMap(ctx, params, portalApiBytes)
	if err != nil {
		params.Log.Info("failed to reconcile configmap", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
	}
}
