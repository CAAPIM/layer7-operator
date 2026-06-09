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
 */
package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// ValidateRef rejects branch and tag values that contain path traversal
// sequences or separator characters to prevent os.RemoveAll / path construction
// from escaping /tmp.
func ValidateRef(ref string) error {
	if strings.Contains(ref, "..") || strings.ContainsAny(ref, "/\\") {
		return fmt.Errorf("branch/tag %q contains invalid characters", ref)
	}
	return nil
}

func CloneRepository(url string, username string, token string, privateKey []byte, privateKeyPass string, branch string, tag string, remoteName string, name string, vendor string, authType string, knownHosts []byte, namespace string, insecureSkipVerify bool) (string, error) {
	if branch != "" {
		if err := ValidateRef(branch); err != nil {
			return "", err
		}
	}
	if tag != "" {
		if err := ValidateRef(tag); err != nil {
			return "", err
		}
	}

	if remoteName == "" {
		remoteName = "origin"
	}

	cloneOpts := git.CloneOptions{
		URL:        url,
		RemoteName: remoteName,
	}

	pullOpts := git.PullOptions{
		RemoteName: remoteName,
	}

	if !strings.HasSuffix(url, ".git") {
		cloneOpts.URL = url + ".git"
	}

	if strings.ToLower(vendor) == "azure" {
		transport.UnsupportedCapabilities = []capability.Capability{
			capability.ThinPack,
		}
		cloneOpts.URL = url
	}

	if insecureSkipVerify || strings.Contains(strings.ToLower(vendor), "insecure") {
		cloneOpts.InsecureSkipTLS = true
		pullOpts.InsecureSkipTLS = true
	}

	if tag != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(tag)
		pullOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + tag)
	}

	// this supercedes tag if set.
	if branch != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(branch)
		pullOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + branch)
	}

	switch strings.ToLower(authType) {
	case "ssh":
		if strings.Contains(url, "https") {
			return "", fmt.Errorf("auth type %s is not valid for %s please use username,token instead", authType, url)
		}
		publicKeys, err := ssh.NewPublicKeys("git", privateKey, privateKeyPass)
		if err != nil {
			return "", err
		}
		cloneOpts.Auth = publicKeys
		pullOpts.Auth = publicKeys

		cb, err := knownHostsCallbackFromBytes(knownHosts)
		if err != nil {
			return "", fmt.Errorf("failed to parse known_hosts: %w", err)
		}
		publicKeys.HostKeyCallback = cb

	case "basic":
		if username != "" && token != "" {
			cloneOpts.Auth = &http.BasicAuth{Username: username, Password: token}
			pullOpts.Auth = &http.BasicAuth{Username: username, Password: token}
		}
	}

	ext := cloneOpts.ReferenceName.String()

	r, err := git.PlainClone("/tmp/"+name+"-"+namespace+"-"+ext, false, &cloneOpts)

	if err == git.ErrRepositoryAlreadyExists {
		r, _ := git.PlainOpen("/tmp/" + name + "-" + namespace + "-" + ext)
		w, _ := r.Worktree()

		ref, _ := r.Head()

		if ref == nil {
			_ = os.RemoveAll("/tmp/" + name + "-" + namespace + "-" + ext)
			return "", fmt.Errorf("ref is nil for %s", name)
		}
		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			return "", err
		}

		if ext == tag {
			return commit.Hash.String(), nil
		}

		gbytes, _ := os.ReadFile("/tmp/" + name + "-" + namespace + "-" + ext + "/.git/config")
		if !strings.Contains(string(gbytes), cloneOpts.URL) {
			err = os.RemoveAll("/tmp/" + name + "-" + namespace + "-" + ext)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("invalid git config for %s removing temp storage", name)
		}

		err = w.Pull(&pullOpts)
		if err != nil {
			if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
				return commit.Hash.String(), err
			}
			return "", err
		}

		// Re-read HEAD after a successful pull to return the actual new commit SHA.
		// The pre-pull `commit` captured above reflects the old HEAD; without this
		// re-read the caller receives the stale SHA and treats the new push as if
		// nothing changed, causing a one-cycle delay before any commit is detected.
		newRef, err := r.Head()
		if err != nil {
			return "", err
		}
		newCommit, err := r.CommitObject(newRef.Hash())
		if err != nil {
			return "", err
		}
		return newCommit.Hash.String(), nil
	}

	if err != nil {
		return "", err
	}

	ref, _ := r.Head()
	commit, err := r.CommitObject(ref.Hash())

	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

// knownHostsCallbackFromBytes builds a per-repository SSH HostKeyCallback from
// raw known_hosts bytes without touching any process-wide file or environment
// variable. It writes to a unique temp file, lets knownhosts.New parse it into
// an in-memory closure, then removes the file immediately.
func knownHostsCallbackFromBytes(data []byte) (gossh.HostKeyCallback, error) {
	f, err := os.CreateTemp("", "layer7-known-hosts-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(data); err != nil {
		f.Close()
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	return knownhosts.New(f.Name())
}
