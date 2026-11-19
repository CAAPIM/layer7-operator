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
	"path/filepath"
	"strings"
	"testing"
)

// TestTmpPathsInjective asserts that (name, namespace) pairs which collided
// under the old hyphen-joined layout ("apis"/"team-a" vs "apis-team"/"a" both
// mapped to "apis-team-a") now resolve to distinct directories.
func TestTmpPathsInjective(t *testing.T) {
	builders := map[string]func(name, namespace string) string{
		"RepoCacheDir":         RepoCacheDir,
		"StateStoreCacheDir":   StateStoreCacheDir,
		"GatewayBundleWorkDir": GatewayBundleWorkDir,
		"RepoCheckoutDir":      func(name, namespace string) string { return RepoCheckoutDir(name, namespace, "main") },
	}
	for fn, build := range builders {
		victim := build("apis", "team-a")
		attacker := build("apis-team", "a")
		if victim == attacker {
			t.Errorf("%s: %q collides for apis/team-a and apis-team/a", fn, victim)
		}
		// Each directory must nest namespace then name.
		if !strings.Contains(victim, string(filepath.Separator)+filepath.Join("team-a", "apis")) {
			t.Errorf("%s: %q does not nest namespace/name", fn, victim)
		}
	}
}

func TestRepoCheckoutDirInjective(t *testing.T) {
	// name/ref boundary: ("a", ref "b-c") vs ("a-b", ref "c") in one namespace.
	if RepoCheckoutDir("a", "ns", "b-c") == RepoCheckoutDir("a-b", "ns", "c") {
		t.Error("name/ref boundary collides")
	}
	// Distinct bases: a namespace named like another base dir must not land in
	// that base (e.g. namespace "repo-cache" vs RepoCacheDir).
	if RepoCheckoutDir("victim-ns", "repo-cache", "victim") == RepoCacheDir("victim", "victim-ns") {
		t.Error("checkout dir collides with repo-cache dir")
	}
	if a, b := RepoCacheDir("x", "y"), StateStoreCacheDir("x", "y"); a == b {
		t.Errorf("cache bases collide: %q", a)
	}
}

func TestRepoCheckoutDirRefContained(t *testing.T) {
	base := filepath.Clean(repositoryCheckoutBase)
	for _, ref := range []string{"", ".", "..", "main", "feature/x", "release/1.3", "refs/tags/v1"} {
		dir := RepoCheckoutDir("repo", "ns", ref)
		want := filepath.Join(base, "ns", "repo")
		if filepath.Dir(dir) != want {
			t.Errorf("RepoCheckoutDir(ref=%q) = %q, want a single child of %q", ref, dir, want)
		}
	}
	if got := RepoCheckoutDir("repo", "ns", "feature/x"); filepath.Base(got) != "feature-x" {
		t.Errorf("ref not flattened with SafeRef: %q", got)
	}
}
