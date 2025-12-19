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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package reconcile

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	v1 "github.com/caapim/layer7-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const repositoryFinalizer = "security.brcmlabs.com/layer7-operator"

func Finalizer(ctx context.Context, params Params) (err error) {
	isMarkedToBeDeleted := false
	if params.Instance.DeletionTimestamp != nil {
		isMarkedToBeDeleted = true
	}

	if isMarkedToBeDeleted {
		// Check if repository is still being actively used by any gateways
		inUse, gatewayNames, err := isRepositoryInUse(ctx, params)
		if err != nil {
			params.Log.Error(err, "failed to check if repository is in use", "name", params.Instance.Name)
			return err
		}

		if inUse {
			params.Log.Info("repository is still referenced by active gateways, skipping deletion",
				"repository", params.Instance.Name,
				"namespace", params.Instance.Namespace,
				"gateways", gatewayNames)
			return fmt.Errorf("repository %s is still referenced by gateways: %v", params.Instance.Name, gatewayNames)
		}

		// Clean up repository files and folders
		switch params.Instance.Spec.Type {
		case v1.RepositoryTypeGit:
			ref := params.Instance.Spec.Tag

			if params.Instance.Spec.Branch != "" {
				ref = params.Instance.Spec.Branch
			}

			err = os.RemoveAll("/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + ref)

			if err != nil {
				return err
			}
		case v1.RepositoryTypeHttp:
			fileURL, _ := url.Parse(params.Instance.Spec.Endpoint)
			path := fileURL.Path
			segments := strings.Split(path, "/")
			fileName := segments[len(segments)-1]
			ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]
			folderName := strings.ReplaceAll(fileName, "."+ext, "")
			if ext == "gz" && strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
				folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
			}

			fileName = "/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + fileName
			folderName = "/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + folderName

			err = os.RemoveAll(folderName)
			if err != nil {
				return err
			}

			err = os.Remove(fileName)
			if err != nil {
				return err
			}
		default:
			break
		}

		// Clean up cache directory
		cachePath := "/tmp/repo-cache/" + params.Instance.Name
		if err := os.RemoveAll(cachePath); err != nil {
			params.Log.V(2).Info("failed to remove cache directory", "path", cachePath, "error", err.Error())
		}

		// Clean up statestore cache if it exists
		stateStorePath := "/tmp/statestore/" + params.Instance.Name
		if err := os.RemoveAll(stateStorePath); err != nil {
			params.Log.V(2).Info("failed to remove statestore cache directory", "path", stateStorePath, "error", err.Error())
		}

		// todo
		// remove from statestore if push
		err = removeFinalizer(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// isRepositoryInUse checks if the repository is actively referenced by any gateways
func isRepositoryInUse(ctx context.Context, params Params) (bool, []string, error) {
	// List all gateways in the same namespace
	gatewayList := &v1.GatewayList{}
	if err := params.Client.List(ctx, gatewayList, client.InNamespace(params.Instance.Namespace)); err != nil {
		return false, nil, fmt.Errorf("failed to list gateways: %w", err)
	}

	var referencingGateways []string

	// Check each gateway for references to this repository
	for _, gateway := range gatewayList.Items {
		for _, repoRef := range gateway.Spec.App.RepositoryReferences {
			if repoRef.Name == params.Instance.Name && repoRef.Enabled {
				referencingGateways = append(referencingGateways, gateway.Name)
				break
			}
		}
	}

	if len(referencingGateways) > 0 {
		return true, referencingGateways, nil
	}

	return false, nil, nil
}

func removeFinalizer(ctx context.Context, params Params) error {
	params.Log.V(2).Info("removing finalizer", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	controllerutil.RemoveFinalizer(params.Instance, repositoryFinalizer)
	err := params.Client.Update(ctx, params.Instance)
	if err != nil {
		params.Log.V(2).Info("fail to remove finalizer", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
		return err
	}

	return nil
}
