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
package graphman

import (
	"crypto/sha1"
	"encoding/base64"
	"testing"
)

func certBundle(thumb string) Bundle {
	der := []byte{1, 2, 3, 4, 5}
	return Bundle{TrustedCerts: []*TrustedCertInput{{
		Name:           "apim-ssg.example",
		CertBase64:     base64.StdEncoding.EncodeToString(der),
		ThumbprintSha1: thumb,
	}}}
}

func wantThumb() string {
	sum := sha1.Sum([]byte{1, 2, 3, 4, 5})
	return base64.StdEncoding.EncodeToString(sum[:])
}

func TestTrustedCertWithoutThumbprintSurvivesMerge(t *testing.T) {
	merges := map[string]func(Bundle, Bundle) (Bundle, error){
		"CombineWithOverwrite":                     CombineWithOverwrite,
		"CombineWithOverwritePreservingDeleteMaps": CombineWithOverwritePreservingDeleteMappings,
	}
	for name, merge := range merges {
		out, err := merge(certBundle(""), Bundle{})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if len(out.TrustedCerts) != 1 {
			t.Fatalf("%s: trusted cert was dropped (got %d)", name, len(out.TrustedCerts))
		}
		if got := out.TrustedCerts[0].ThumbprintSha1; got != wantThumb() {
			t.Fatalf("%s: thumbprint = %q, want %q", name, got, wantThumb())
		}
	}
}

func TestTrustedCertExistingThumbprintUnchanged(t *testing.T) {
	out, err := CombineWithOverwrite(certBundle("keep-me"), Bundle{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.TrustedCerts) != 1 || out.TrustedCerts[0].ThumbprintSha1 != "keep-me" {
		t.Fatalf("existing thumbprint was changed: %+v", out.TrustedCerts)
	}
}

func TestTrustedCertBadBase64LeftAlone(t *testing.T) {
	b := Bundle{TrustedCerts: []*TrustedCertInput{{Name: "x", CertBase64: "%%%not-base64"}}}
	normalizeTrustedCertThumbprints(&b)
	if b.TrustedCerts[0].ThumbprintSha1 != "" {
		t.Fatalf("thumbprint set for undecodable cert: %q", b.TrustedCerts[0].ThumbprintSha1)
	}
}
