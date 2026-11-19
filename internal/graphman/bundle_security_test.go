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
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupRefRepo creates <tmp>/repo/<entityDir> containing an in-repo file and a symlink that escapes the
// repo to <tmp>/secret.txt. It returns the entity directory.
func setupRefRepo(t *testing.T, entityDir string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "repo", entityDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "secret.txt"), []byte("top-secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "content.txt"), []byte("in-repo"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "secret.txt"), filepath.Join(dir, "escape.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("content.txt", filepath.Join(dir, "inside.txt")); err != nil {
		t.Fatal(err)
	}
	return dir
}

func expectEscapeError(t *testing.T, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), "escapes repository directory") {
		t.Fatalf("expected escape error, got %v", err)
	}
}

func TestParseCacString_RejectsSymlinkEscape(t *testing.T) {
	dir := setupRefRepo(t, "trustedCerts")
	_, err := parseCacString("trustedCerts", filepath.Join(dir, "cert.json"), "{{escape.txt}}")
	expectEscapeError(t, err)
}

func TestParseCacString_RejectsDotDot(t *testing.T) {
	dir := setupRefRepo(t, "trustedCerts")
	_, err := parseCacString("trustedCerts", filepath.Join(dir, "cert.json"), "{{../../secret.txt}}")
	expectEscapeError(t, err)
}

func TestParseCacString_ReadsInRepoReference(t *testing.T) {
	dir := setupRefRepo(t, "trustedCerts")
	want := base64.StdEncoding.EncodeToString([]byte("in-repo"))
	for _, ref := range []string{"{{content.txt}}", "{{inside.txt}}"} {
		got, err := parseCacString("trustedCerts", filepath.Join(dir, "cert.json"), ref)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", ref, err)
		}
		if got != want {
			t.Errorf("%s: expected %q, got %q", ref, want, got)
		}
	}
}

func TestParsePacCode_RejectsSymlinkEscape(t *testing.T) {
	dir := setupRefRepo(t, "tree")
	_, err := parsePacCode(&PolicyInput{Xml: "{{escape.txt}}"}, filepath.Join(dir, "api.service.json"))
	expectEscapeError(t, err)
}

func TestParsePacCode_RejectsDotDot(t *testing.T) {
	dir := setupRefRepo(t, "tree")
	_, err := parsePacCode(&PolicyInput{Json: "{../../secret.txt}"}, filepath.Join(dir, "api.policy.json"))
	expectEscapeError(t, err)
}

func TestParsePacCode_ReadsInRepoReference(t *testing.T) {
	dir := setupRefRepo(t, "tree")
	for _, ref := range []string{"{{content.txt}}", "{{inside.txt}}"} {
		p, err := parsePacCode(&PolicyInput{Xml: ref}, filepath.Join(dir, "api.policy.json"))
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", ref, err)
		}
		if p.Xml != "in-repo" {
			t.Errorf("%s: expected %q, got %q", ref, "in-repo", p.Xml)
		}
	}
}

func TestConcatBundle_NullEntityReturnsError(t *testing.T) {
	acc, _ := json.Marshal(Bundle{Services: []*L7ServiceInput{{Name: "api1"}}})
	src := []byte(`{"services":[null]}`)

	for name, concat := range map[string]func([]byte, []byte) ([]byte, error){
		"ConcatBundle":                   ConcatBundle,
		"ConcatBundlePreservingMappings": ConcatBundlePreservingMappings,
	} {
		result, err := concat(src, acc)
		if err == nil {
			t.Errorf("%s: expected error for null entity, got result %s", name, result)
		}
		if result != nil {
			t.Errorf("%s: expected nil result on error, got %s", name, result)
		}
	}
}

func TestCombineWithOverwrite_NullEntityReturnsError(t *testing.T) {
	src := Bundle{Services: []*L7ServiceInput{nil}}
	dest := Bundle{Services: []*L7ServiceInput{{Name: "api1"}}}
	if _, err := CombineWithOverwrite(src, dest); err == nil {
		t.Error("CombineWithOverwrite: expected error for null src entity")
	}
	if _, err := CombineWithOverwritePreservingDeleteMappings(dest, src); err == nil {
		t.Error("CombineWithOverwritePreservingDeleteMappings: expected error for null dest entity")
	}
}

func TestRemoveDuplicates_NullEntityReturnsError(t *testing.T) {
	result, err := RemoveDuplicates([]byte(`{"services":[{"name":"api1"},null]}`))
	if err == nil {
		t.Fatalf("expected error for null entity, got result %s", result)
	}
}

func TestResetMappings_SkipsNullEntriesWithoutError(t *testing.T) {
	var bundle Bundle
	data := `{"services":[null,{"name":"api1"},{"name":"api2"}],` +
		`"properties":{"mappings":{"services":[null,{"action":"DELETE","source":{"name":"api1"}}]}}}`
	if err := json.Unmarshal([]byte(data), &bundle); err != nil {
		t.Fatal(err)
	}
	if err := ResetMappings(&bundle); err != nil {
		t.Fatalf("ResetMappings failed: %v", err)
	}
	for _, s := range bundle.Services {
		if s != nil && s.Name == "api1" {
			t.Error("expected api1 to be removed by delete mapping")
		}
	}
}
