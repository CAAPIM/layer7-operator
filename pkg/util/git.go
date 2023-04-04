package util

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepository(url string, username string, token string, branch string, name string) (string, error) {
	r, err := git.PlainClone("/tmp/"+name, false, &git.CloneOptions{
		URL:               url,
		Auth:              &http.BasicAuth{Username: username, Password: token},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		RemoteName:        branch,
	})

	if err == git.ErrRepositoryAlreadyExists {
		r, _ := git.PlainOpen("/tmp/" + name)
		w, _ := r.Worktree()

		ref, _ := r.Head()
		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			return "", err
		}
		_ = w.Pull(&git.PullOptions{RemoteName: branch, Auth: &http.BasicAuth{Username: username, Password: token}})

		// if err == git.NoErrAlreadyUpToDate {
		// 	return commit.Hash.String(), nil
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
