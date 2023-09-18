package util

import (
	"encoding/json"
	"time"
)

type PortalPapiToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

var portalAuthSyncCache = NewSyncCache(3 * time.Second)

func GetPortalAccessToken(name string, authServer string, clientId string, clientSecret string) (string, error) {

	requestCacheEntry := name + "-" + clientId
	syncRequest, _ := portalAuthSyncCache.Read(requestCacheEntry)

	if syncRequest.CacheData != "" {
		token := syncRequest.CacheData
		return token, nil
	}

	tokenEndpoint := "https://" + authServer + "/auth/oauth/v2/token"
	formData := "grant_type=client_credentials&scope=OOB"

	var portalPapiToken PortalPapiToken
	resp, err := RestCall("POST", tokenEndpoint, true, map[string]string{}, "application/x-www-form-urlencoded", []byte(formData), clientId, clientSecret)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(resp, &portalPapiToken)

	if err != nil {
		return "", err
	}

	cacheData := portalPapiToken.AccessToken

	portalAuthSyncCache.Update(SyncRequest{RequestName: requestCacheEntry, CacheData: cacheData}, time.Now().Add(time.Duration(portalPapiToken.ExpiresIn-100)*time.Second).Unix())
	syncRequest, _ = portalAuthSyncCache.Read(requestCacheEntry)
	return portalPapiToken.AccessToken, nil
}
