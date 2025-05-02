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

	if delete && err == nil {
		_ = deleteRepoRefStatus(ctx, params, *gwUpdReq.repository)
	} else {
		_ = updateRepoRefStatus(ctx, params, *gwUpdReq.repository, gwUpdReq.repositoryReference.Type, gwUpdReq.checksum, err)
	}
	gwUpdReq = nil
	if err != nil {
		return err
	}

	return nil
}
