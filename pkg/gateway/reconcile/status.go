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
	"fmt"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func GatewayStatus(ctx context.Context, params Params) error {
	gatewayStatus := params.Instance.Status
	gatewayStatus.Host = params.Instance.Spec.App.Management.Cluster.Hostname
	gatewayStatus.Image = params.Instance.Spec.App.Image
	gatewayStatus.Version = params.Instance.Spec.Version
	gatewayStatus.Gateway = []securityv1.GatewayState{}

	dep, err := getGatewayDeployment(ctx, params)
	if err != nil {
		params.Log.V(2).Info("deployment hasn't been created yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	} else {
		gatewayStatus.Replicas = dep.Status.Replicas
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		params.Log.V(2).Info("pods aren't available yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	if len(params.Instance.Spec.App.RepositoryReferences) < len(gatewayStatus.RepositoryStatus) {
		for i, repoStatus := range gatewayStatus.RepositoryStatus {
			found := false
			for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
				if repoStatus.Name == repoRef.Name {
					gatewayStatus.RepositoryStatus[i].Enabled = repoRef.Enabled
					found = true
				}
			}
			if !found {
				gatewayStatus.RepositoryStatus[i].Enabled = false
				gatewayStatus.RepositoryStatus[i].Commit = ""
			}
		}
	}

	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("repository not found", "name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
			return err
		}

		if repository.Status.Commit == "" {
			return fmt.Errorf("repository %s is not ready yet", repository.Name)
		}

		found := false
		for i, repoStatus := range gatewayStatus.RepositoryStatus {
			if repoStatus.Name == repository.Name {
				rs, err := buildRepoStatus(ctx, params, *repository, repoRef)
				if err != nil {
					params.Log.V(2).Info("failed to build repository status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
					return err
				}
				if len(gatewayStatus.RepositoryStatus[i].Conditions) > 0 {
					rs.Conditions = gatewayStatus.RepositoryStatus[i].Conditions
				}
				rs.Directories = gatewayStatus.RepositoryStatus[i].Directories
				gatewayStatus.RepositoryStatus[i] = rs
				found = true
			}
		}

		if !found && repoRef.Enabled {
			rs, err := buildRepoStatus(ctx, params, *repository, repoRef)
			if err != nil {
				params.Log.V(2).Info("failed to build repository status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
				return err
			}
			gatewayStatus.RepositoryStatus = append(gatewayStatus.RepositoryStatus, rs)
		}
	}

	if podList != nil {
		for _, p := range podList.Items {
			if p.ObjectMeta.Labels["management-access"] == "leader" {
				gatewayStatus.ManagementPod = p.Name
			}
		}
	}

	if !reflect.DeepEqual(gatewayStatus, params.Instance.Status) {
		params.Instance.Status = gatewayStatus
		err = params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.V(2).Info("failed to update gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
		params.Log.V(2).Info("updated gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}
	return nil
}

func buildRepoStatus(ctx context.Context, params Params, repository securityv1.Repository, repoRef securityv1.RepositoryReference) (repoStatus securityv1.GatewayRepositoryStatus, err error) {
	secretName := repository.Name
	if repository.Spec.Auth.ExistingSecretName != "" {
		secretName = repository.Spec.Auth.ExistingSecretName
	}

	if repository.Spec.Auth == (securityv1.RepositoryAuth{}) {
		secretName = ""
	}

	rs := securityv1.GatewayRepositoryStatus{
		Commit:            repository.Status.Commit,
		Enabled:           repoRef.Enabled,
		Name:              repoRef.Name,
		RepoType:          string(repository.Spec.Type),
		Vendor:            repository.Spec.Auth.Vendor,
		AuthType:          string(repository.Spec.Auth.Type),
		Type:              string(repoRef.Type),
		SecretName:        secretName,
		StorageSecretName: repository.Status.StorageSecretName,
		Endpoint:          repository.Spec.Endpoint,
		//Directories:       repoRef.Directories,
	}

	if repository.Spec.Tag != "" && repository.Spec.Branch == "" {
		rs.Tag = repository.Spec.Tag
	}

	if repository.Spec.Branch != "" {
		rs.Branch = repository.Spec.Branch
	}

	rs.RemoteName = "origin"
	if repository.Spec.RemoteName != "" {
		rs.RemoteName = repository.Spec.RemoteName
	}

	if repository.Spec.StateStoreReference != "" {
		ext := repository.Spec.Branch
		if ext == "" {
			ext = repository.Spec.Tag
		}
		stateStoreKey := repository.Name + "-repository-" + ext
		rs.StateStoreReference = repository.Spec.StateStoreReference
		statestore := &securityv1alpha1.L7StateStore{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: repository.Spec.StateStoreReference, Namespace: params.Instance.Namespace}, statestore)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("state store not found", "name", repository.Spec.StateStoreReference, "repository", repository.Name, "namespace", params.Instance.Namespace)
			return securityv1.GatewayRepositoryStatus{}, err
		}
		rs.StateStoreKey = statestore.Spec.Redis.GroupName + ":" + statestore.Spec.Redis.StoreId + ":" + "repository" + ":" + stateStoreKey + ":latest"
		if repository.Spec.StateStoreKey != "" {
			rs.StateStoreKey = repository.Spec.StateStoreKey
		}
	}
	return rs, nil
}
