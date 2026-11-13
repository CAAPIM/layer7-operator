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
	"encoding/pem"
	"strings"
	"testing"

	"github.com/caapim/layer7-operator/internal/graphman"
)

// malformedCertPEM is a well-formed PEM CERTIFICATE block whose body is not valid DER.
var malformedCertPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("not a der certificate")}))

func TestSplitAppliedCertKey(t *testing.T) {
	tests := []struct {
		entry, name, thumbprint string
		ok                      bool
	}{
		{entry: "simple-AbC+/12=", name: "simple", thumbprint: "AbC+/12=", ok: true},
		{entry: "my-app.example.com-AbC+/12=", name: "my-app.example.com", thumbprint: "AbC+/12=", ok: true},
		{entry: "-AbC+/12=", name: "", thumbprint: "AbC+/12=", ok: true},
		{entry: "noseparator", ok: false},
		{entry: "trailing-", ok: false},
		{entry: "", ok: false},
	}
	for _, tt := range tests {
		name, thumbprint, ok := SplitAppliedCertKey(tt.entry)
		if ok != tt.ok || name != tt.name || thumbprint != tt.thumbprint {
			t.Errorf("SplitAppliedCertKey(%q) = (%q, %q, %v), want (%q, %q, %v)", tt.entry, name, thumbprint, ok, tt.name, tt.thumbprint, tt.ok)
		}
	}
}

func TestConvertCertsToGraphmanBundle_DashedNameDeleteMapping(t *testing.T) {
	_, crt, err := createDummyKey()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := pem.Decode([]byte(crt))
	thumbprint, err := getSha1Thumbprint(b.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(thumbprint, "-") {
		t.Fatalf("thumbprint %q contains '-'", thumbprint)
	}

	// same format the gateway status uses for applied external certs
	bundleBytes, err := ConvertCertsToGraphmanBundle(nil, []string{"my-app.example.com-" + thumbprint, "malformed"})
	if err != nil {
		t.Fatal(err)
	}
	bundle := graphman.Bundle{}
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		t.Fatal(err)
	}
	if len(bundle.TrustedCerts) != 1 {
		t.Fatalf("expected 1 trusted cert, got %d", len(bundle.TrustedCerts))
	}
	if bundle.TrustedCerts[0].Name != "my-app.example.com" || bundle.TrustedCerts[0].ThumbprintSha1 != thumbprint {
		t.Errorf("unexpected trusted cert name=%q thumbprint=%q", bundle.TrustedCerts[0].Name, bundle.TrustedCerts[0].ThumbprintSha1)
	}
	if bundle.Properties == nil || len(bundle.Properties.Mappings.TrustedCerts) != 1 {
		t.Fatalf("expected 1 delete mapping")
	}
	mapping := bundle.Properties.Mappings.TrustedCerts[0]
	sourceBytes, _ := json.Marshal(mapping.Source)
	source := graphman.MappingSource{}
	if err := json.Unmarshal(sourceBytes, &source); err != nil {
		t.Fatal(err)
	}
	if mapping.Action != graphman.MappingActionDelete || source.ThumbprintSha1 != thumbprint {
		t.Errorf("unexpected mapping action=%q thumbprint=%q", mapping.Action, source.ThumbprintSha1)
	}
}

func TestConvertCertsToGraphmanBundle_MalformedDERSkipped(t *testing.T) {
	_, crt, err := createDummyKey()
	if err != nil {
		t.Fatal(err)
	}
	bundleBytes, err := ConvertCertsToGraphmanBundle([]GraphmanCert{
		{Name: "bad", Crt: malformedCertPEM},
		{Name: "good", Crt: crt},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	bundle := graphman.Bundle{}
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		t.Fatal(err)
	}
	if len(bundle.TrustedCerts) != 1 || bundle.TrustedCerts[0].Name != "good" {
		t.Fatalf("expected only the valid cert in the bundle, got %d certs", len(bundle.TrustedCerts))
	}
}

func TestConvertX509ToGraphmanBundle_MalformedDERSkipped(t *testing.T) {
	key, crt, err := createDummyKey()
	if err != nil {
		t.Fatal(err)
	}
	bundleBytes, err := ConvertX509ToGraphmanBundle([]GraphmanKey{
		{Name: "only-bad", Crt: malformedCertPEM, Key: key},
		{Name: "mixed", Crt: malformedCertPEM + crt, Key: key},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	bundle := graphman.Bundle{}
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		t.Fatal(err)
	}
	if len(bundle.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(bundle.Keys))
	}
	if chainLen(bundle.Keys[0].CertChain) != 0 {
		t.Errorf("expected malformed cert to be dropped from chain, got %d certs", chainLen(bundle.Keys[0].CertChain))
	}
	if chainLen(bundle.Keys[1].CertChain) != 1 {
		t.Errorf("expected only the valid cert in chain, got %d certs", chainLen(bundle.Keys[1].CertChain))
	}
}

// chainLen returns the length of a decoded KeyInput.CertChain.
func chainLen(chain interface{}) int {
	c, _ := chain.([]interface{})
	return len(c)
}
