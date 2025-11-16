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
package reconcile

import (
	"context"
	"strings"

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

		localRepoAvailable, err := checkLocalRepoOnFs(params, repository)
		if err != nil {
			return err
		}

		// If repository not available locally, check if storage secret exists with required bundles
		if !localRepoAvailable {
			if repository.Status.StorageSecretName != "" && repository.Status.StorageSecretName != "_" {
				storageSecret, err := getGatewaySecret(ctx, params, repository.Status.StorageSecretName)
				if err != nil {
					params.Log.V(2).Info("storage secret not found", "secret", repository.Status.StorageSecretName, "repository", repository.Name)
					return nil
				}

				// Check if the storage secret has bundles matching the requested directories
				hasBundles := false
				for _, dir := range repoRef.Directories {
					// Normalize directory name to match bundle key format
					normalizedDir := strings.TrimPrefix(strings.ReplaceAll(dir, "/", "-"), "-")
					if normalizedDir == "" || normalizedDir == "-" {
						normalizedDir = "bundle" // Root directory bundles might be stored as "bundle.gz"
					}

					// Check for both .gz and non-.gz variants
					bundleKey := normalizedDir + ".gz"
					if bundleData, exists := storageSecret.Data[bundleKey]; exists && len(bundleData) > 20 {
						hasBundles = true
						break
					}

					// Also check without .gz extension
					if bundleData, exists := storageSecret.Data[normalizedDir]; exists && len(bundleData) > 20 {
						hasBundles = true
						break
					}
				}

				if !hasBundles {
					params.Log.V(2).Info("storage secret exists but does not contain required bundles",
						"secret", repository.Status.StorageSecretName,
						"repository", repository.Name,
						"directories", repoRef.Directories)
					return nil
				}

				params.Log.Info("using storage secret for repository",
					"secret", repository.Status.StorageSecretName,
					"repository", repository.Name,
					"directories", repoRef.Directories)
				// Continue with bundle building - buildBundleFromCache will use the storage secret
			} else {
				return nil
			}
		}
	}

	commit := repository.Status.Commit

	// Only enable delete functionality if delete was requested (repository disabled/removed)
	// AND RepositoryReferenceDelete is enabled
	if delete && gateway.Spec.App.RepositoryReferenceDelete.Enabled {
		if gateway.Spec.App.RepositoryReferenceDelete.LimitToStateStore && repository.Spec.StateStoreReference == "" {
			delete = false
		}
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

	_ = updateRepoRefStatus(ctx, params, *gwUpdReq.repository, *gwUpdReq.repositoryReference, gwUpdReq.checksum, err, delete)
	gwUpdReq = nil
	if err != nil {
		return err
	}

	return nil
}
