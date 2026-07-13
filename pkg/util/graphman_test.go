/*
* Copyright (c) 2026 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/go-logr/logr"
)

func TestGraphmanBundleBytesHaveNoEntities(t *testing.T) {
	if !GraphmanBundleBytesHaveNoEntities([]byte(`{}`)) {
		t.Fatal("expected empty object to have no entities")
	}
	empty := graphman.Bundle{}
	b, err := json.Marshal(empty)
	if err != nil {
		t.Fatal(err)
	}
	if !GraphmanBundleBytesHaveNoEntities(b) {
		t.Fatal("expected marshaled empty Bundle to have no entities")
	}
	if !GraphmanBundleBytesHaveNoEntities([]byte(`{"goid":""}`)) {
		t.Fatal("expected tiny payload to be treated as empty")
	}
	withService := graphman.Bundle{
		Services: []*graphman.L7ServiceInput{{Name: "s"}},
	}
	b2, err := json.Marshal(withService)
	if err != nil {
		t.Fatal(err)
	}
	if GraphmanBundleBytesHaveNoEntities(b2) {
		t.Fatal("expected bundle with service to not be empty")
	}
	longInvalidJSON := make([]byte, 50)
	copy(longInvalidJSON, []byte(`not json`))
	if GraphmanBundleBytesHaveNoEntities(longInvalidJSON) {
		t.Fatal("invalid JSON should not be treated as empty")
	}
}

func TestBuildAndValidateBundle_LooseFullBundleStrictDecodeFails(t *testing.T) {
	d := t.TempDir()
	valid := `{
    "clusterProperties": [
        {
            "goid": "982cc1ee7369c6ca5a7ae1e4ad866070",
            "name": "cfg",
            "hiddenProperty": false,
            "value": "v"
        }
    ]
}`
	if err := os.WriteFile(filepath.Join(d, "base.json"), []byte(valid), 0644); err != nil {
		t.Fatal(err)
	}
	// Full bundle shape but unknown top-level key fails strict Decode (DisallowUnknownFields).
	bad := `{"services":[],"unexpectedTopLevel":true}`
	// Filename must not contain a substring that ParseEntityPath treats as an entity
	// (e.g. "bad-bundle.json" matches ".service" and would be skipped here).
	if err := os.WriteFile(filepath.Join(d, "invalid-full.json"), []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := BuildAndValidateBundle(d, false, logr.Discard())
	if err == nil {
		t.Fatal("expected error for invalid full bundle loose JSON")
	}
	if !strings.Contains(err.Error(), "invalid-full.json") {
		t.Fatalf("expected error to mention file path, got: %v", err)
	}
}

func TestBuildAndValidateBundle_UnrelatedJSONIgnored(t *testing.T) {
	d := t.TempDir()
	valid := `{"services":[{"name":"s","resolutionPath":"/x","enabled":true}]}`
	if err := os.WriteFile(filepath.Join(d, "svc.json"), []byte(valid), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "settings.json"), []byte(`{"app":{"foo":1}}`), 0644); err != nil {
		t.Fatal(err)
	}
	b, err := BuildAndValidateBundle(d, false, logr.Discard())
	if err != nil {
		t.Fatal(err)
	}
	if len(b) <= 40 {
		t.Fatalf("expected non-trivial bundle, len=%d", len(b))
	}
}

func TestBuildAndValidateBundle_DotServiceStringVersionFails(t *testing.T) {
	d := t.TempDir()
	// Basename must contain ".service" so parseEntity classifies as .service (readBundle case ".service").
	// L7ServiceInput.Version is int; a JSON string must surface an error with file path (not be dropped).
	bad := `{
  "goid": "6d2ca7e0a5c6c9f37dd0b124a6a26532",
  "guid": "0ed53d22-bed8-424b-8b4e-68e3e40a9311",
  "name": "Rest Api 4",
  "resolutionPath": "/api4",
  "serviceType": "WEB_API",
  "enabled": true,
  "version": "1"
}`
	path := filepath.Join(d, "api.service.json")
	if err := os.WriteFile(path, []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := BuildAndValidateBundle(d, false, logr.Discard())
	if err == nil {
		t.Fatal("expected error for invalid .service entity JSON (version must be number)")
	}
	if !strings.Contains(err.Error(), "api.service.json") {
		t.Fatalf("expected error to mention file path, got: %v", err)
	}
	if !strings.Contains(err.Error(), ".service") {
		t.Fatalf("expected error to mention entity kind, got: %v", err)
	}
}

func TestBuildAndValidateBundle_ServicesFolderStringVersionFails(t *testing.T) {
	d := t.TempDir()
	servicesDir := filepath.Join(d, "services")
	if err := os.MkdirAll(servicesDir, 0755); err != nil {
		t.Fatal(err)
	}
	bad := `{
  "goid": "6d2ca7e0a5c6c9f37dd0b124a6a26532",
  "guid": "0ed53d22-bed8-424b-8b4e-68e3e40a9311",
  "name": "svc",
  "resolutionPath": "/api4",
  "serviceType": "WEB_API",
  "enabled": true,
  "version": "1"
}`
	path := filepath.Join(servicesDir, "bad.json")
	if err := os.WriteFile(path, []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := BuildAndValidateBundle(d, false, logr.Discard())
	if err == nil {
		t.Fatal("expected error for invalid JSON under services/")
	}
	if !strings.Contains(err.Error(), "bad.json") {
		t.Fatalf("expected error to mention file path, got: %v", err)
	}
	if !strings.Contains(err.Error(), "services") {
		t.Fatalf("expected error to mention entity kind, got: %v", err)
	}
}

func TestConvertOpaqueMapToGraphmanBundle(t *testing.T) {
	secrets := []GraphmanSecret{{Name: "test1", Secret: "secret1"}, {Name: "test2", Secret: "secret2"}}
	bundleBytes, err := ConvertOpaqueMapToGraphmanBundle(secrets, []string{})
	if err != nil {
		t.Errorf("Error getting secret bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "secret1") {
		t.Errorf("bundle %s should contain %s", bundle, "secret1")
	}
	if !strings.Contains(bundle, "secret2") {
		t.Errorf("bundle %s should contain %s", bundle, "secret2")
	}

}

/*func TestConvertX509ToGraphmanBundle(t *testing.T) {
	key := GraphmanKey{
		Name: "test",
		Crt:  "-----BEGIN CERTIFICATE-----MIIC6jCCAdKgAwIBAgIGAYnTwy1CMA0GCSqGSIb3DQEBCwUAMDYxNDAyBgNVBAMMK2FNZTk4OXNSdk9YTlgyckgtYmZjVHUyZVB5NFJhVXhOOHpCenJsSzFpcFkwHhcNMjMwODA4MDYwODUxWhcNMjQwNjAzMDYwODUxWjA2MTQwMgYDVQQDDCthTWU5ODlzUnZPWE5YMnJILWJmY1R1MmVQeTRSYVV4Tjh6QnpybEsxaXBZMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0D5Q71trDwX4BnFR4jnxEBsohmD/R9CU19eWGN5iKRHVxQwsM9tR1569NZkUqAXRexo/6RHp3IT8fXdai2+227i3tpSt6hKoiVNMFznktnQ5WTZpBIa6D7iwdGZCDoY17juBrNr5ko/WWsIIVD06Z14BIs6wcqyM2QHPaKLQAgt+ZJfRps3vWCmoBHxLRhuQrcxpPhYnb/ZFsNW6fq5aJA7TG5fU7PKo69DVWUVga65ysTEEb79c7ytHHUrdEE5oR2dFmemN6yev36I92oSFqb5sBKkn2lim9VCTY6ZGZitF3XbUSqfJGkDxHIANLRi+trPdI71RKTWBtHMTkIsFhwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBJ7ShdGRSKeMVPnmb9NnQX9aZlU5Sphbb5UkgTCdd5y+8k/QpKgk+BG4u5P3wN359X1HgpDQGh3OfboMhZMJY2VnQ3qK7W0r8au6IQ5mFtlrUukBWjxAJc/1rbdBD2TlCHdEBpqgg2s7fEgeu+6NRIeJFYDLXOiQaZES01WMxL9CZDfxijwJSO6ZSSEMlDQ0K0UY3p/B0V0rSvXTrJIPE8boDzksL/0GiRBFOc0tQhqtq33h7pnKW70CjDiM7ib2fuZLLtLse+jrbZiJ79bINRmB+kd5HNJtI5xTTwXvf+sfs2v81Wdmpzdv3aKIVcnDk63+lVVh9+114QifWNNeuy-----END CERTIFICATE-----",
		Key:  "-----BEGIN PRIVATE KEY-----MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQPlDvW2sPBfgGcVHiOfEQGyiGYP9H0JTX15YY3mIpEdXFDCwz21HXnr01mRSoBdF7Gj/pEenchPx9d1qLb7bbuLe2lK3qEqiJU0wXOeS2dDlZNmkEhroPuLB0ZkIOhjXuO4Gs2vmSj9ZawghUPTpnXgEizrByrIzZAc9ootACC35kl9Gmze9YKagEfEtGG5CtzGk+Fidv9kWw1bp+rlokDtMbl9Ts8qjr0NVZRWBrrnKxMQRvv1zvK0cdSt0QTmhHZ0WZ6Y3rJ6/foj3ahIWpvmwEqSfaWKb1UJNjpkZmK0XddtRKp8kaQPEcgA0tGL62s90jvVEpNYG0cxOQiwWHAgMBAAECggEAbtF6yKXhpzEJ+IO9i6JCAswxGLHtqA3755E2sy1FF44CMMZ1j3Mbbp9vGWLJd1EBVX12nVWHGm863pnxeVqN+Qen3GXq1zHutoW5bHLGn8Hh8vPdlycLROqIHKl+ZbROZuUL8Szmu3QIImw3enzK489G03sisyPYIHOyKIDcKPl2OUEGeznXXVgzPI5LC5vXeazy7nI9ykzaBdlxf2bnykIR9RShRWxaJQ5/xqZ8hnywaGexXzv2Vpo88a3KRL9f9CyhCRR3Su0cCVua3rwj/Ijl3hx5oBPzXn0nZLeCZ8NhpC89BnQ5T0u15wKOVA9axKMFfE4PpZ/MXjaOJlOT4QKBgQD7foeM9lAhzUIW1k9qzgcR38q7BFLn8oK7FKRv1PIL0GGL+w6G00QYJ9xAQBJkZH2bzdij8dhflmbh/C3Of6ptxDuHFOoCGZArFWqbQSiMB3fb1AwGMc9hATeJ+d0Jx5/GvgWrzheVFCFXvAK9dkMIIKWYZ7xM0gA+lZrfCZuICwKBgQDT+WsGjzGIbrERirXhA1GhysYSrLP6fMGy0Ko0qlHgrwkAI1YF8eK2zqNdApqTpFoonxvwSAPICZNsh2rGup1tBjSaN3HCJC6PZonuzbApuPcfp+QZbuOBEakgp05HV2tG7cGoV48eW+FBrzuXXZgItZsZMly2VJHSxRhqVvJZ9QKBgQCwPfhyIX8AYR5qcJ9RArbToNgqfRo4b6uLvSiLMli5TLu/ZB3HADCdGPnxkLUS45Ve5T9njKkMO5M31QioyLC/oZ/xxwdCl3V/q898o4ntr6IgTJZslOV2Xmr0Z0SugNWIakwBHTlLgMLo/9mPulu5S1+g0TmVQClpsl/I46u6BwKBgQDSmK9bMfKdQJQNAImYhyqYGpRVQ14AU+hBVoxzjG+SUXQYvgKeH2YGByBIrOiUHKoyR3mDbJjNKa5dGeDcldUH1y11tfYAUuArOk15gsMtgIWM3smA9ylyNvCX74CW4mRDcL2BGZSoLdKK5qTGyobcyEjSbLWttDG4fHa4V6+p7QKBgGWLpT8TmN7qif0L07BiPTkqptqSvZRXP4lLh6wkAZxrOxzkJbSQjM+rW/cmuLWphyyMbru3xqrORzMjuc7t/3FKwbd97ZlYba8tvThNFPTA0cVAjplBIpEIcfnpBd0oPklxsj28fENz6dQ+3a1ZzSDF65kRE9/R/fgR3fw9LikM-----END PRIVATE KEY-----",
		Port: "8443"}
	keys := []GraphmanKey{key}

	bundleBytes, err := ConvertX509ToGraphmanBundle(keys)
	if err != nil {
		t.Errorf("Error getting key bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "test") {
		t.Errorf("bundle %s should contain key %s", bundle, "test")
	}
}*/
