package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type CloneRepositoryOpts struct {
	URL            string
	Username       string
	Token          string
	PrivateKey     []byte
	PrivateKeyPass string
	Branch         string
	Tag            string
	RemoteName     string
	Name           string
	Vendor         string
	AuthType       string
	KnownHosts     []byte
	Namespace      string
}

func CloneRepository(url string, opts *CloneRepositoryOpts) (string, error) {

	if opts.RemoteName == "" {
		opts.RemoteName = "origin"
	}

	cloneOpts := git.CloneOptions{
		URL:        url,
		RemoteName: opts.RemoteName,
	}

	pullOpts := git.PullOptions{
		RemoteName: opts.RemoteName,
	}

	if !strings.Contains(url, ".git") {
		cloneOpts.URL = url + ".git"
	}

	if opts.Tag != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(opts.Tag)
		pullOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + opts.Tag)
	}

	// this supercedes tag if set.
	if opts.Branch != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(opts.Branch)
		pullOpts.ReferenceName = plumbing.ReferenceName("refs/heads/" + opts.Branch)
	}

	switch strings.ToLower(opts.AuthType) {
	case "ssh":
		if strings.Contains(url, "https") {
			return "", fmt.Errorf("auth type %s is not valid for %s please use username,token instead", opts.AuthType, url)
		}
		publicKeys, err := ssh.NewPublicKeys("git", opts.PrivateKey, opts.PrivateKeyPass)
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
			err = os.WriteFile("/tmp/known_hosts", opts.KnownHosts, 0644)
			if err != nil {
				return "", err
			}
		} else {
			if len(currentKnownHosts) == 0 {
				newKnownHosts = string(opts.KnownHosts)
			}
			for _, c := range strings.Split(string(currentKnownHosts), "\n") {
				if !strings.Contains(newKnownHosts, c) {
					if newKnownHosts == "" {
						newKnownHosts = c
					} else {
						newKnownHosts = newKnownHosts + "\n" + c
					}

					for _, n := range strings.Split(string(opts.KnownHosts), "\n") {
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
		if opts.Username != "" && opts.Token != "" {
			cloneOpts.Auth = &http.BasicAuth{Username: opts.Username, Password: opts.Token}
			pullOpts.Auth = &http.BasicAuth{Username: opts.Username, Password: opts.Token}
		}
	}

	ext := cloneOpts.ReferenceName.String()

	r, err := git.PlainClone("/tmp/"+opts.Name+"-"+opts.Namespace+"-"+ext, false, &cloneOpts)

	if err == git.ErrRepositoryAlreadyExists {
		r, _ := git.PlainOpen("/tmp/" + opts.Name + "-" + opts.Namespace + "-" + ext)
		w, _ := r.Worktree()

		ref, _ := r.Head()

		if ref == nil {
			_ = os.RemoveAll("/tmp/" + opts.Name + "-" + opts.Namespace + "-" + ext)
			return "", fmt.Errorf("ref is nil for %s", opts.Name)
		}
		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			return "", err
		}

		if ext == opts.Tag {
			return commit.Hash.String(), nil
		}

		gbytes, _ := os.ReadFile("/tmp/" + opts.Name + "-" + opts.Namespace + "-" + ext + "/.git/config")
		if !strings.Contains(string(gbytes), cloneOpts.URL) {
			err = os.RemoveAll("/tmp/" + opts.Name + "-" + opts.Namespace + "-" + ext)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("invalid git config for %s removing temp storage", opts.Name)
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
