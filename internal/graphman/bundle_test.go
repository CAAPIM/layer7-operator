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

func TestConcatBundle_EmptyAccumulator(t *testing.T) {
	// Test concatenating into an empty accumulator
	newBundle := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1"},
		},
	}

	newBytes, _ := json.Marshal(newBundle)
	result, err := ConcatBundle(newBytes, []byte("{}"))
	if err != nil {
		t.Fatalf("ConcatBundle failed: %v", err)
	}

	var resultBundle Bundle
	json.Unmarshal(result, &resultBundle)

	if len(resultBundle.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(resultBundle.Services))
	}
	if resultBundle.Services[0].Name != "api1" {
		t.Errorf("Expected service 'api1', got '%s'", resultBundle.Services[0].Name)
	}
}

func TestConcatBundle_LatestWins(t *testing.T) {
	// Test that latest bundle wins when same entity appears in multiple bundles
	// This is the key scenario we discussed - later bundles overwrite earlier ones
	acc := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1", Enabled: false},
		},
	}

	newBundle := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1-updated", Enabled: true},
		},
	}

	accBytes, _ := json.Marshal(acc)
	newBytes, _ := json.Marshal(newBundle)
	resultBytes, err := ConcatBundle(newBytes, accBytes)
	if err != nil {
		t.Fatalf("ConcatBundle failed: %v", err)
	}

	var result Bundle
	json.Unmarshal(resultBytes, &result)

	if len(result.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(result.Services))
	}
	if result.Services[0].ResolutionPath != "/api1-updated" {
		t.Errorf("Expected '/api1-updated', got '%s'", result.Services[0].ResolutionPath)
	}
	if !result.Services[0].Enabled {
		t.Errorf("Expected enabled=true from latest bundle")
	}
}

func TestConcatBundle_MultipleIterations(t *testing.T) {
	// Test multiple concatenations to ensure ordering is preserved
	// folder1/api1, folder2/api1 (should overwrite), folder3/api2

	// First folder - api1 with value A
	folder1 := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1", Goid: "valueA"},
		},
	}
	folder1Bytes, _ := json.Marshal(folder1)
	result1, err := ConcatBundle(folder1Bytes, []byte("{}"))
	if err != nil {
		t.Fatalf("ConcatBundle failed on iteration 1: %v", err)
	}

	// Second folder - api1 with value B (should overwrite)
	folder2 := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1", Goid: "valueB"},
		},
	}
	folder2Bytes, _ := json.Marshal(folder2)
	result2, err := ConcatBundle(folder2Bytes, result1)
	if err != nil {
		t.Fatalf("ConcatBundle failed on iteration 2: %v", err)
	}

	// Third folder - api2 (new entity)
	folder3 := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api2", ResolutionPath: "/api2", Goid: "valueC"},
		},
	}
	folder3Bytes, _ := json.Marshal(folder3)
	result3, err := ConcatBundle(folder3Bytes, result2)
	if err != nil {
		t.Fatalf("ConcatBundle failed on iteration 3: %v", err)
	}

	var finalResult Bundle
	json.Unmarshal(result3, &finalResult)

	// Final result should have: api1 with valueB (from folder2), api2 with valueC (from folder3)
	if len(finalResult.Services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(finalResult.Services))
	}

	// Find api1
	foundApi1 := false
	for _, svc := range finalResult.Services {
		if svc.Name == "api1" {
			if svc.Goid != "valueB" {
				t.Errorf("Expected api1 to have valueB (from folder2), got '%s'", svc.Goid)
			}
			foundApi1 = true
		}
	}
	if !foundApi1 {
		t.Errorf("api1 not found in result")
	}
}

func TestConcatBundle_PreservesUniqueEntities(t *testing.T) {
	// Test that unique entities from each bundle are preserved
	acc := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1"},
		},
		Policies: []*L7PolicyInput{
			{Name: "policy1"},
		},
	}

	newBundle := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api2", ResolutionPath: "/api2"},
		},
		ClusterProperties: []*ClusterPropertyInput{
			{Name: "cluster.hostname", Value: "test"},
		},
	}

	accBytes, _ := json.Marshal(acc)
	newBytes, _ := json.Marshal(newBundle)
	resultBytes, err := ConcatBundle(newBytes, accBytes)
	if err != nil {
		t.Fatalf("ConcatBundle failed: %v", err)
	}

	var result Bundle
	json.Unmarshal(resultBytes, &result)

	if len(result.Services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(result.Services))
	}
	if len(result.Policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(result.Policies))
	}
	if len(result.ClusterProperties) != 1 {
		t.Errorf("Expected 1 cluster property, got %d", len(result.ClusterProperties))
	}
}

func TestConcatBundle_DuplicatesAcrossFolders(t *testing.T) {
	// This is the key scenario we discussed - when same API appears in folder1 and folder2,
	// the one from folder2 (processed later) should win
	folder1 := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1", FolderPath: "/folder1"},
		},
	}

	folder2 := Bundle{
		Services: []*L7ServiceInput{
			{Name: "api1", ResolutionPath: "/api1", FolderPath: "/folder2"},
		},
	}

	// Simulate processing folder1 then folder2
	folder1Bytes, _ := json.Marshal(folder1)
	result, err := ConcatBundle(folder1Bytes, []byte("{}"))
	if err != nil {
		t.Fatalf("ConcatBundle failed on folder1: %v", err)
	}

	folder2Bytes, _ := json.Marshal(folder2)
	result, err = ConcatBundle(folder2Bytes, result)
	if err != nil {
		t.Fatalf("ConcatBundle failed on folder2: %v", err)
	}

	var finalResult Bundle
	json.Unmarshal(result, &finalResult)

	// Should have only 1 service with folder path from folder2
	if len(finalResult.Services) != 1 {
		t.Errorf("Expected 1 service after deduplication, got %d", len(finalResult.Services))
	}
	if finalResult.Services[0].FolderPath != "/folder2" {
		t.Errorf("Expected folder path '/folder2' (latest), got '%s'", finalResult.Services[0].FolderPath)
	}
}
