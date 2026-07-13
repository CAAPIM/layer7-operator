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

package reconcile

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
)

func TestParseStateStoreKeyPayload_rawJSON(t *testing.T) {
	blob := []byte(`{"services":[{"name":"x"}]}`)
	got, err := parseStateStoreKeyPayload(blob)
	if err != nil {
		t.Fatal(err)
	}
	var b graphman.Bundle
	if err := json.Unmarshal(got, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.Services) != 1 || b.Services[0].Name != "x" {
		t.Fatalf("bundle: %+v", b)
	}
}

func TestParseStateStoreKeyPayload_gzipJSON(t *testing.T) {
	raw := []byte(`{"services":[]}`)
	gz, err := util.GzipCompress(raw)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseStateStoreKeyPayload(gz)
	if err != nil {
		t.Fatal(err)
	}
	var b graphman.Bundle
	if err := json.Unmarshal(got, &b); err != nil {
		t.Fatal(err)
	}
}

func TestParseStateStoreKeyPayload_emptyError(t *testing.T) {
	_, err := parseStateStoreKeyPayload([]byte("   "))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseStateStoreKeyPayload_rejectsTopLevelArray(t *testing.T) {
	_, err := parseStateStoreKeyPayload([]byte(`[{"name":"x"}]`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "top-level JSON object") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestStateStorePayloadToLatestJSON_shapeAndRoundTrip(t *testing.T) {
	repoName := "my-repo"
	raw := []byte(`{"services":[{"name":"svc"}]}`)
	jsonBytes, err := stateStorePayloadToLatestJSON(repoName, raw)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string][]byte
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		t.Fatal(err)
	}
	key := repoName + ".gz"
	gz, ok := m[key]
	if !ok || len(m) != 1 {
		t.Fatalf("map keys: %#v", m)
	}
	bundleJSON, err := util.GzipDecompress(gz)
	if err != nil {
		t.Fatal(err)
	}
	var b graphman.Bundle
	if err := json.Unmarshal(bundleJSON, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.Services) != 1 || b.Services[0].Name != "svc" {
		t.Fatalf("services: %+v", b.Services)
	}
}
