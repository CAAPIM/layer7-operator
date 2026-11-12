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

import "testing"

func TestValidateRef(t *testing.T) {
	valid := []string{
		"", // empty is allowed (no branch/tag set)
		"main",
		"feature/x",     // slash-namespaced branch
		"release/1.3.0", // slash + dots
		"bugfix/JIRA-123",
		"v1.2.3",           // tag
		"team/sub/feature", // multiple slashes
	}
	for _, ref := range valid {
		if err := ValidateRef(ref); err != nil {
			t.Errorf("ValidateRef(%q) = %v, want nil", ref, err)
		}
	}

	invalid := []string{
		"../../etc",   // traversal
		"a/../b",      // embedded traversal
		"..",          // bare traversal
		"/abs",        // leading slash
		"trailing/",   // trailing slash
		"a//b",        // double slash
		"back\\slash", // backslash
		"tab\there",   // control character
		"null\x00byte",
	}
	for _, ref := range invalid {
		if err := ValidateRef(ref); err == nil {
			t.Errorf("ValidateRef(%q) = nil, want error", ref)
		}
	}
}

func TestSafeRef(t *testing.T) {
	cases := map[string]string{
		"main":          "main",
		"feature/x":     "feature-x",
		"release/1.3.0": "release-1.3.0",
		"team/sub/leaf": "team-sub-leaf",
		"back\\slash":   "back-slash",
		"":              "",
	}
	for in, want := range cases {
		if got := SafeRef(in); got != want {
			t.Errorf("SafeRef(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestSafeRefContainment asserts that a SafeRef'd ref, when used to build the
// repository /tmp working directory, cannot escape /tmp for any ref that passes
// ValidateRef.
func TestSafeRefContainment(t *testing.T) {
	refs := []string{"feature/x", "release/1.3", "a-b-c", "team/sub/leaf"}
	for _, ref := range refs {
		if err := ValidateRef(ref); err != nil {
			t.Fatalf("ValidateRef(%q) unexpectedly failed: %v", ref, err)
		}
		dir := "/tmp/name-namespace-" + SafeRef(ref)
		// A single flat segment under /tmp: exactly two path separators
		// ("/tmp/" prefix) and no additional "/" from the ref.
		if got := countSlashesAfter(dir, "/tmp/"); got != 0 {
			t.Errorf("SafeRef(%q) produced a nested dir %q (%d extra separators)", ref, dir, got)
		}
	}
}

func countSlashesAfter(s, prefix string) int {
	rest := s[len(prefix):]
	n := 0
	for _, r := range rest {
		if r == '/' || r == '\\' {
			n++
		}
	}
	return n
}
