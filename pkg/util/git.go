package util

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepository(url string, username string, token string, branch string, tag string, remoteName string, name string, vendor string) (string, error) {

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

	if vendor == "gitlab" {
		if !strings.Contains(url, ".git") {
			cloneOpts.URL = url + ".git"
		}
	}

	if tag != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(tag)
	}

	// this supercedes tag if set.
	if branch != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(branch)
	}

	if username != "" && token != "" {
		cloneOpts.Auth = &http.BasicAuth{Username: username, Password: token}
		pullOpts.Auth = &http.BasicAuth{Username: username, Password: token}
	}

	ext := cloneOpts.ReferenceName.String()

	r, err := git.PlainClone("/tmp/"+name+"-"+ext, false, &cloneOpts)
	if err == git.ErrRepositoryAlreadyExists {
		r, _ := git.PlainOpen("/tmp/" + name + "-" + ext)
		w, _ := r.Worktree()

		ref, _ := r.Head()
		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			return "", err
		}

		if ext == tag {
			return commit.Hash.String(), nil
		}

		_ = w.Pull(&pullOpts)
		// if err != nil {
		// 	return "", err
		// }

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
