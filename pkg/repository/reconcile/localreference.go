package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func LocalReference(ctx context.Context, params Params) error {
	repository, err := getRepository(ctx, params)
	if err != nil {
		return err
	}
	var commit string

	params.Instance = &repository

	repoStatus := repository.Status
	if !repository.Spec.Enabled {
		return nil
	}
	start := time.Now()

	patch := []byte(`[{"op": "replace", "path": "/status/ready", "value": false}]`)

	switch strings.ToLower(repository.Spec.Type) {
	case "local":
		commit, err = localReferenceShaSum(ctx, repository, params)
		if err != nil {
			err = setRepoStatus(ctx, params, patch)
			if err != nil {
				return err
			}
			return err
		}
	default:
		return nil
	}

	repoStatus.Commit = commit
	repoStatus.Name = repository.Name
	repoStatus.Vendor = repository.Spec.Auth.Vendor
	repoStatus.Ready = true

	repoStatus.StorageSecretName = repository.Spec.LocalReference.SecretName

	if !reflect.DeepEqual(repoStatus, repository.Status) {

		params.Log.Info("syncing repository", "name", repository.Name, "namespace", repository.Namespace)
		repoStatus.Updated = time.Now().String()
		repository.Status = repoStatus
		err = params.Client.Status().Update(ctx, &repository)
		if err != nil {
			_ = captureRepositorySyncMetrics(ctx, params, start, commit, true)
			params.Log.Info("failed to update repository status", "namespace", repository.Namespace, "name", repository.Name, "error", err.Error())
			return nil
		}
		params.Log.Info("reconciled", "name", repository.Name, "namespace", repository.Namespace, "commit", commit)
	}
	_ = captureRepositorySyncMetrics(ctx, params, start, commit, false)
	return nil
}

func localReferenceShaSum(ctx context.Context, repository securityv1.Repository, params Params) (string, error) {
	if repository.Spec.LocalReference.SecretName == "" {
		return "", fmt.Errorf("%s localReference secret name must be set", repository.Name)
	}

	localReference := &corev1.Secret{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.LocalReference.SecretName, Namespace: repository.Namespace}, localReference)
	if err != nil {
		return "", err
	}

	// bundleBytes, err := util.ConcatBundles(localReference.Data)
	// if err != nil {
	// 	return "", err
	// }
	dataBytes, _ := json.Marshal(&localReference.Data)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	return sha1Sum, nil
}
