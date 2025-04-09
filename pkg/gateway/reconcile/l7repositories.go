package reconcile

import (
	"context"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func ExternalRepository(ctx context.Context, params Params) error {
	gateway := params.Instance

	for _, repoRef := range gateway.Spec.App.RepositoryReferences {
		if repoRef.Enabled {
			err := reconcileDynamicRepository(ctx, params, repoRef, false)
			if err != nil {
				params.Log.Error(err, "failed to reconcile repository reference", "name", gateway.Name, "repository", repoRef.Name, "namespace", gateway.Namespace)
				return err
			}
		}
	}

	for _, repoStatus := range gateway.Status.RepositoryStatus {
		found := false
		disabled := false
		for _, repoRef := range gateway.Spec.App.RepositoryReferences {
			if repoStatus.Name == repoRef.Name {
				found = true
				if !repoRef.Enabled {
					disabled = true
				}
			}
		}
		if !found || disabled {
			repoRef := securityv1.RepositoryReference{Name: repoStatus.Name, Type: "dynamic", Encryption: securityv1.BundleEncryption{Passphrase: "delete"}}
			err := reconcileDynamicRepository(ctx, params, repoRef, true)
			if err != nil {
				params.Log.Error(err, "failed to remove repository reference", "name", gateway.Name, "repository", repoRef.Name, "namespace", gateway.Namespace)
				return err
			}
		}
	}

	return nil
}

func reconcileDynamicRepository(ctx context.Context, params Params, repoRef securityv1.RepositoryReference, delete bool) (err error) {
	gateway := params.Instance

	repository := &securityv1.Repository{}

	err = params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: gateway.Namespace}, repository)
	if err != nil && k8serrors.IsNotFound(err) {
		return err
	}

	if !repository.Status.Ready {
		params.Log.Info("repository not ready", "repository", repository.Name, "name", gateway.Name, "namespace", gateway.Namespace)
		return nil
	}

	commit := repository.Status.Commit
	// only support delete if a statestore is used
	if repository.Spec.StateStoreReference == "" {
		delete = false
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithChecksum(commit),
		WithDelete(delete),
		WithBundleType(BundleTypeRepository),
		WithRepositoryReference(repoRef),
		WithRepository(repository),
	)

	if err != nil {
		return err
	}

	if gwUpdReq == nil {
		return nil
	}

	err = SyncGateway(ctx, params, *gwUpdReq)
	if err != nil {
		gwUpdReq = nil
		return err
	}

	for _, sRepo := range gateway.Status.RepositoryStatus {
		if sRepo.Name == repoRef.Name {
			if sRepo.Commit != commit {
				_ = GatewayStatus(ctx, params)
			}
		}
	}
	gwUpdReq = nil
	return nil
}
