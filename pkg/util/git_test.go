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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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
		prefix := filepath.Join(repositoryCheckoutBase, "namespace", "name") + "/"
		dir := RepoCheckoutDir("name", "namespace", ref)
		// A single flat segment under {base}/{namespace}/{name}/: no additional
		// "/" from the ref.
		if got := countSlashesAfter(dir, prefix); got != 0 {
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

// newSourceRepo creates a local git repository (path ending in ".git" so
// CloneRepository uses it verbatim) with one commit on master and a tag v1.
func newSourceRepo(t *testing.T, content string) (path string, head string) {
	t.Helper()
	path = filepath.Join(t.TempDir(), "src.git")
	r, err := git.PlainInit(path, false)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := os.WriteFile(filepath.Join(path, "file.txt"), []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}
	if _, err := w.Add("file.txt"); err != nil {
		t.Fatalf("add: %v", err)
	}
	h, err := w.Commit(content, &git.CommitOptions{Author: &object.Signature{Name: "t", Email: "t@example.com", When: time.Now()}})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if _, err := r.CreateTag("v1", h, nil); err != nil {
		t.Fatalf("tag: %v", err)
	}
	return path, h.String()
}

func withCheckoutBase(t *testing.T) {
	t.Helper()
	orig := repositoryCheckoutBase
	repositoryCheckoutBase = t.TempDir()
	t.Cleanup(func() { repositoryCheckoutBase = orig })
}

func TestRemoteURLMatches(t *testing.T) {
	const want = "https://github.com/victim/repo.git"
	cases := []struct {
		name string
		urls []string
		ok   bool
	}{
		{"exact", []string{want}, true},
		{"substring", []string{"https://evil.example/pwn.git?u=" + want}, false},
		{"prefix", []string{want + "x"}, false},
		{"multiple urls", []string{want, "https://evil.example/pwn.git"}, false},
	}
	for _, tc := range cases {
		r, err := git.PlainInit(t.TempDir(), false)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: tc.urls}); err != nil {
			t.Fatal(err)
		}
		if got := remoteURLMatches(r, "origin", want); got != tc.ok {
			t.Errorf("%s: remoteURLMatches = %v, want %v", tc.name, got, tc.ok)
		}
		if remoteURLMatches(r, "upstream", want) {
			t.Errorf("%s: missing remote matched", tc.name)
		}
	}
	if remoteURLMatches(nil, "origin", want) {
		t.Error("nil repository matched")
	}
}

// TestCloneRepositoryReusedDirRemoteMismatch pre-seeds the checkout directory
// with a clone of a different repository and asserts CloneRepository discards
// it and re-clones, for both the tag early-return and the branch/pull paths.
func TestCloneRepositoryReusedDirRemoteMismatch(t *testing.T) {
	cases := []struct {
		name   string
		branch string
		tag    string
		// remoteURL builds the seeded remote URL from (evil, good).
		remoteURL func(evil, good string) string
	}{
		{"tag path, foreign remote", "", "v1", func(evil, _ string) string { return evil }},
		{"branch path, foreign remote", "master", "", func(evil, _ string) string { return evil }},
		{"branch path, substring remote", "master", "", func(evil, good string) string { return evil + "#" + good }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withCheckoutBase(t)
			good, goodHead := newSourceRepo(t, "good")
			evil, evilHead := newSourceRepo(t, "evil")

			ref := tc.tag
			if tc.branch != "" {
				ref = tc.branch
			}
			dir := RepoCheckoutDir("repo", "ns", ref)
			seeded, err := git.PlainClone(dir, false, &git.CloneOptions{URL: evil})
			if err != nil {
				t.Fatalf("seed clone: %v", err)
			}
			if err := seeded.DeleteRemote("origin"); err != nil {
				t.Fatal(err)
			}
			if _, err := seeded.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{tc.remoteURL(evil, good)}}); err != nil {
				t.Fatal(err)
			}

			got, err := CloneRepository(good, "", "", nil, "", tc.branch, tc.tag, "", "repo", "", "", nil, "ns", false)
			if err != nil {
				t.Fatalf("CloneRepository: %v", err)
			}
			if got == evilHead {
				t.Fatal("returned the foreign checkout's HEAD")
			}
			if got != goodHead {
				t.Fatalf("CloneRepository = %s, want %s", got, goodHead)
			}
			b, err := os.ReadFile(filepath.Join(dir, "file.txt"))
			if err != nil || string(b) != "good" {
				t.Fatalf("checkout content = %q (%v), want re-cloned good content", b, err)
			}
		})
	}
}

// TestCloneRepositoryReusedDirMatching asserts a checkout whose remote matches
// exactly is reused rather than re-cloned.
func TestCloneRepositoryReusedDirMatching(t *testing.T) {
	withCheckoutBase(t)
	good, goodHead := newSourceRepo(t, "good")

	for _, ref := range []struct{ branch, tag string }{{"master", ""}, {"", "v1"}} {
		first, err := CloneRepository(good, "", "", nil, "", ref.branch, ref.tag, "", "repo", "", "", nil, "ns", false)
		if err != nil || first != goodHead {
			t.Fatalf("initial clone (%+v) = %s, %v", ref, first, err)
		}
		marker := filepath.Join(RepoCheckoutDir("repo", "ns", ref.branch+ref.tag), "marker")
		if err := os.WriteFile(marker, nil, 0644); err != nil {
			t.Fatal(err)
		}
		second, err := CloneRepository(good, "", "", nil, "", ref.branch, ref.tag, "", "repo", "", "", nil, "ns", false)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			t.Fatalf("reuse (%+v): %v", ref, err)
		}
		if second != goodHead {
			t.Fatalf("reuse (%+v) = %s, want %s", ref, second, goodHead)
		}
		if _, err := os.Stat(marker); err != nil {
			t.Fatalf("matching checkout (%+v) was re-cloned instead of reused", ref)
		}
	}
}
