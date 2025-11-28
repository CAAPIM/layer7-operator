/*
* Copyright (c) 2025 Broadcom. All rights reserved.
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

package graphman

import (
	"encoding/json"
	"testing"
)

func TestCombineWithOverwrite_LatestWins(t *testing.T) {
	// Test that when same entity appears in both bundles, src (latest) wins
	bundle1 := Bundle{
		Services: []*L7ServiceInput{
			{
				Name:           "api1",
				ResolutionPath: "/api1",
				Goid:           "goid1",
			},
		},
	}

	bundle2 := Bundle{
		Services: []*L7ServiceInput{
			{
				Name:           "api1",
				ResolutionPath: "/api1-updated",
				Goid:           "goid1",
			},
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.Services) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(result.Services))
	}

	if result.Services[0].ResolutionPath != "/api1-updated" {
		t.Errorf("Expected updated resolution path '/api1-updated', got '%s'", result.Services[0].ResolutionPath)
	}
}

func TestCombineWithOverwrite_HttpConfigurationCompositeKey(t *testing.T) {
	// Test that HttpConfiguration uses composite key (Host + Port)
	bundle1 := Bundle{
		HttpConfigurations: []*HttpConfigurationInput{
			{Host: "localhost", Port: 8080, Path: "/old"},
			{Host: "localhost", Port: 8443, Path: "/ssl"},
		},
	}

	bundle2 := Bundle{
		HttpConfigurations: []*HttpConfigurationInput{
			{Host: "localhost", Port: 8080, Path: "/new"}, // Should overwrite same host:port
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.HttpConfigurations) != 2 {
		t.Errorf("Expected 2 http configs, got %d", len(result.HttpConfigurations))
	}

	// Find the 8080 config
	found := false
	for _, cfg := range result.HttpConfigurations {
		if cfg.Port == 8080 {
			if cfg.Path != "/new" {
				t.Errorf("Expected path '/new' for port 8080, got '%s'", cfg.Path)
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Could not find http config for port 8080")
	}
}

func TestCombineWithOverwrite_KeysByAlias(t *testing.T) {
	// Test that Keys are matched by Alias
	bundle1 := Bundle{
		Keys: []*KeyInput{
			{Alias: "ssl-key", KeystoreId: "00000000000000000000000000000002:keystore1"},
			{Alias: "signing-key", KeystoreId: "00000000000000000000000000000002:keystore1"},
		},
	}

	bundle2 := Bundle{
		Keys: []*KeyInput{
			{Alias: "ssl-key", KeystoreId: "00000000000000000000000000000002:keystore2"}, // Should overwrite
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.Keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(result.Keys))
	}

	// Find ssl-key
	found := false
	for _, key := range result.Keys {
		if key.Alias == "ssl-key" {
			if key.KeystoreId != "00000000000000000000000000000002:keystore2" {
				t.Errorf("Expected updated keystore for ssl-key")
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Could not find ssl-key")
	}
}

func TestCombineWithOverwrite_TrustedCertsByThumbprint(t *testing.T) {
	// Test that TrustedCerts are matched by ThumbprintSha1
	bundle1 := Bundle{
		TrustedCerts: []*TrustedCertInput{
			{Name: "cert1", ThumbprintSha1: "abc123"},
			{Name: "cert2", ThumbprintSha1: "def456"},
		},
	}

	bundle2 := Bundle{
		TrustedCerts: []*TrustedCertInput{
			{Name: "cert1-updated", ThumbprintSha1: "abc123"}, // Should overwrite based on thumbprint
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.TrustedCerts) != 2 {
		t.Errorf("Expected 2 trusted certs, got %d", len(result.TrustedCerts))
	}

	// Find cert with thumbprint abc123
	found := false
	for _, cert := range result.TrustedCerts {
		if cert.ThumbprintSha1 == "abc123" {
			if cert.Name != "cert1-updated" {
				t.Errorf("Expected updated name for cert with thumbprint abc123")
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Could not find cert with thumbprint abc123")
	}
}

func TestCombineWithOverwrite_SchemasBySystemId(t *testing.T) {
	// Test that Schemas are matched by SystemId
	bundle1 := Bundle{
		Schemas: []*SchemaInput{
			{SystemId: "http://example.com/schema1", Goid: "goid1"},
			{SystemId: "http://example.com/schema2", Goid: "goid2"},
		},
	}

	bundle2 := Bundle{
		Schemas: []*SchemaInput{
			{SystemId: "http://example.com/schema1", Goid: "goid3"}, // Should overwrite
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.Schemas) != 2 {
		t.Errorf("Expected 2 schemas, got %d", len(result.Schemas))
	}

	// Find schema1
	found := false
	for _, schema := range result.Schemas {
		if schema.SystemId == "http://example.com/schema1" {
			if schema.Goid != "goid3" {
				t.Errorf("Expected updated goid for schema1")
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Could not find schema1")
	}
}

func TestCombineWithOverwrite_DtdsBySystemId(t *testing.T) {
	// Test that DTDs are matched by SystemId
	bundle1 := Bundle{
		Dtds: []*DtdInput{
			{SystemId: "http://example.com/dtd1.dtd", Goid: "goid1"},
		},
	}

	bundle2 := Bundle{
		Dtds: []*DtdInput{
			{SystemId: "http://example.com/dtd1.dtd", Goid: "goid2"}, // Should overwrite
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.Dtds) != 1 {
		t.Errorf("Expected 1 dtd, got %d", len(result.Dtds))
	}

	if result.Dtds[0].Goid != "goid2" {
		t.Errorf("Expected updated goid for dtd")
	}
}

func TestCombineWithOverwrite_CustomKeyValuesByKey(t *testing.T) {
	// Test that CustomKeyValues are matched by Key
	bundle1 := Bundle{
		CustomKeyValues: []*CustomKeyValueInput{
			{Key: "custom.key1", Value: "value1"},
			{Key: "custom.key2", Value: "value2"},
		},
	}

	bundle2 := Bundle{
		CustomKeyValues: []*CustomKeyValueInput{
			{Key: "custom.key1", Value: "updated-value1"}, // Should overwrite
		},
	}

	result, err := CombineWithOverwrite(bundle2, bundle1)
	if err != nil {
		t.Fatalf("CombineWithOverwrite failed: %v", err)
	}

	if len(result.CustomKeyValues) != 2 {
		t.Errorf("Expected 2 custom key values, got %d", len(result.CustomKeyValues))
	}

	// Find custom.key1
	found := false
	for _, kv := range result.CustomKeyValues {
		if kv.Key == "custom.key1" {
			if kv.Value != "updated-value1" {
				t.Errorf("Expected updated value for custom.key1")
			}
			found = true
		}
	}
	if !found {
		t.Errorf("Could not find custom.key1")
	}
}

func TestCalculateDelta_DetectsNewEntities(t *testing.T) {
	// Test that new entities are detected
	current := Bundle{}
	desired := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1"},
		},
	}

	delta, combined, err := CalculateDelta(current, desired)
	if err != nil {
		t.Fatalf("CalculateDelta failed: %v", err)
	}

	if len(delta.Services) != 1 {
		t.Errorf("Expected 1 service in delta, got %d", len(delta.Services))
	}
	if len(combined.Services) != 1 {
		t.Errorf("Expected 1 service in combined, got %d", len(combined.Services))
	}
}

func TestCalculateDelta_DetectsDeletedEntities(t *testing.T) {
	// Test that deleted entities get delete mappings
	current := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1"},
		},
	}
	desired := Bundle{}

	delta, combined, err := CalculateDelta(current, desired)
	if err != nil {
		t.Fatalf("CalculateDelta failed: %v", err)
	}

	// Delta should include the entity with delete mapping
	if len(delta.Services) != 1 {
		t.Errorf("Expected 1 service in delta, got %d", len(delta.Services))
	}
	if delta.Properties == nil || len(delta.Properties.Mappings.Services) != 1 {
		t.Errorf("Expected 1 delete mapping in delta")
	}

	// Combined should also have entity with delete mapping
	if len(combined.Services) != 1 {
		t.Errorf("Expected 1 service in combined, got %d", len(combined.Services))
	}
	if combined.Properties == nil || len(combined.Properties.Mappings.Services) != 1 {
		t.Errorf("Expected 1 delete mapping in combined")
	}
	if combined.Properties.Mappings.Services[0].Action != MappingActionDelete {
		t.Errorf("Expected delete action in mapping")
	}
}

func TestResetMappings_RemovesEntitiesWithDeleteMappings(t *testing.T) {
	// Test that ResetMappings removes entities marked for deletion
	bundle := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1"},
			{Name: "api2", ResolutionPath: "/api2"},
		},
		Properties: &BundleProperties{
			Mappings: BundleMappings{
				Services: []*MappingInstructionInput{
					{
						Action: MappingActionDelete,
						Source: MappingSource{Name: "api1"},
					},
				},
			},
		},
	}

	err := ResetMappings(&bundle)
	if err != nil {
		t.Fatalf("ResetMappings failed: %v", err)
	}

	// api1 should be removed, api2 should remain
	if len(bundle.Services) != 1 {
		t.Errorf("Expected 1 service after reset, got %d", len(bundle.Services))
	}
	if bundle.Services[0].Name != "api2" {
		t.Errorf("Expected api2 to remain, got %s", bundle.Services[0].Name)
	}

	// Mappings should be cleared
	if bundle.Properties != nil && len(bundle.Properties.Mappings.Services) != 0 {
		t.Errorf("Expected mappings to be cleared")
	}
}

func TestRemoveDuplicates_RemovesDuplicatesByName(t *testing.T) {
	// Test that duplicates by name are removed (first wins)
	bundleBytes := []byte(`{"services":[{"name":"api1","resolutionPath":"/api1","goid":"goid1"},{"name":"api1","resolutionPath":"/api1-dup","goid":"goid2"}]}`)

	result, err := RemoveDuplicates(bundleBytes)
	if err != nil {
		t.Fatalf("RemoveDuplicates failed: %v", err)
	}

	var resultBundle Bundle
	if err := json.Unmarshal(result, &resultBundle); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(resultBundle.Services) != 1 {
		t.Errorf("Expected 1 service after deduplication, got %d", len(resultBundle.Services))
	}
	if resultBundle.Services[0].ResolutionPath != "/api1" {
		t.Errorf("Expected first occurrence to be kept, got '%s'", resultBundle.Services[0].ResolutionPath)
	}
}

func TestRemoveDuplicates_RemovesDuplicatesByGoid(t *testing.T) {
	// Test that duplicates by goid are removed
	bundleBytes := []byte(`{"services":[{"name":"api1","resolutionPath":"/api1","goid":"same-goid"},{"name":"api2","resolutionPath":"/api2","goid":"same-goid"}]}`)

	result, err := RemoveDuplicates(bundleBytes)
	if err != nil {
		t.Fatalf("RemoveDuplicates failed: %v", err)
	}

	var resultBundle Bundle
	if err := json.Unmarshal(result, &resultBundle); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(resultBundle.Services) != 1 {
		t.Errorf("Expected 1 service after deduplication by goid, got %d", len(resultBundle.Services))
	}
}
