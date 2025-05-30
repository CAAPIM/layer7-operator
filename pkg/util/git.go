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
)

func CloneRepository(url string, username string, token string, privateKey []byte, privateKeyPass string, branch string, tag string, remoteName string, name string, vendor string, authType string, knownHosts []byte, namespace string) (string, error) {
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

		if os.Getenv("SSH_KNOWN_HOSTS") != "/tmp/known_hosts" {
			os.Setenv("SSH_KNOWN_HOSTS", "/tmp/known_hosts")
		}
		var newKnownHosts string
		currentKnownHosts, err := os.ReadFile("/tmp/known_hosts")
		if err != nil {
			err = os.WriteFile("/tmp/known_hosts", knownHosts, 0644)
			if err != nil {
				return "", err
			}
		} else {
			if len(currentKnownHosts) == 0 {
				newKnownHosts = string(knownHosts)
			}
			for _, c := range strings.Split(string(currentKnownHosts), "\n") {
				if !strings.Contains(newKnownHosts, c) {
					if newKnownHosts == "" {
						newKnownHosts = c
					} else {
						newKnownHosts = newKnownHosts + "\n" + c
					}

					for _, n := range strings.Split(string(knownHosts), "\n") {
						if !strings.Contains(newKnownHosts, n) {
							if newKnownHosts == "" {
								newKnownHosts = n
							} else {
								newKnownHosts = newKnownHosts + "\n" + n
							}
						}
					}
				}
			}

			err = os.WriteFile("/tmp/known_hosts", []byte(newKnownHosts), 0644)
			if err != nil {
				return "", err
			}
		}

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

		return commit.Hash.String(), nil
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
