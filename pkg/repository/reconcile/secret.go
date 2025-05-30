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
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	v1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/repository"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Secret(ctx context.Context, params Params) error {

	if params.Instance.Spec.Auth.ExistingSecretName != "" || params.Instance.Spec.Auth == (v1.RepositoryAuth{}) {
		return nil
	}

	data := map[string][]byte{}

	switch params.Instance.Spec.Auth.Type {
	case v1.RepositoryAuthTypeBasic:
		data["USERNAME"] = []byte(params.Instance.Spec.Auth.Username)
		data["PASSWORD"] = []byte(params.Instance.Spec.Auth.Password)
		data["TOKEN"] = []byte(params.Instance.Spec.Auth.Token)
	case v1.RepositoryAuthTypeSSH:
		data["SSH_KEY"] = []byte(params.Instance.Spec.Auth.SSHKey)
		data["SSH_KEY_PASS"] = []byte(params.Instance.Spec.Auth.SSHKeyPass)
		data["KNOWN_HOSTS"] = []byte(params.Instance.Spec.Auth.KnownHosts)
	case v1.RepositoryAuthTypeNone:
		return nil
	default:
		return fmt.Errorf("failed to reconcile secret: %s please set auth type to basic, ssh or none", params.Instance.Name)
	}

	desiredSecret := repository.NewSecret(params.Instance, params.Instance.Name, data)

	if err := reconcileSecret(ctx, params, desiredSecret); err != nil {
		return fmt.Errorf("failed to reconcile secrets: %w", err)
	}

	return nil
}

func StorageSecret(ctx context.Context, params Params) error {
	var storageSecretName string
	ext := params.Instance.Spec.Branch
	if ext == "" {
		ext = params.Instance.Spec.Tag
	}

	if params.Instance.Status.StorageSecretName == "_" {
		return nil
	}

	switch strings.ToLower(string(params.Instance.Spec.Type)) {
	case "http":
		fileURL, err := url.Parse(params.Instance.Spec.Endpoint)
		if err != nil {
			return err
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]
		ext = strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		if ext == "gz" && strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
			folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
		}
		storageSecretName = params.Instance.Name + "-repository-" + folderName
		ext = folderName
	case "git":
		storageSecretName = params.Instance.Name + "-repository-" + ext
	default:
		params.Log.Info("repository type not set", "name", params.Instance.Name, "namespace", params.Instance.Name)
		return nil
	}

	dirSize, err := getDirSize("/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + ext)
	if err != nil {
		return err
	}
	if dirSize/1000.0/1000.0/3.5 > 0.9 {
		return errors.New("exceededMaxSize")
	}

	bundleGzip, err := util.CompressGraphmanBundle("/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + ext)
	if err != nil {
		return err
	}

	data := map[string][]byte{
		params.Instance.Name + ".gz": bundleGzip,
	}

	desiredSecret := repository.NewSecret(params.Instance, storageSecretName, data)

	if err := reconcileSecret(ctx, params, desiredSecret); err != nil {
		return fmt.Errorf("failed to reconcile secrets: %w", err)
	}

	return nil
}

func getDirSize(path string) (sizef float64, err error) {
	var size int64
	var mu sync.Mutex

	var calculateSize func(string) error
	calculateSize = func(p string) error {
		fileInfo, err := os.Lstat(p)
		if err != nil {
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if fileInfo.IsDir() {
			entries, err := os.ReadDir(p)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if err := calculateSize(filepath.Join(p, entry.Name())); err != nil {
					return err
				}
			}
		} else {
			mu.Lock()
			size += fileInfo.Size()
			mu.Unlock()
		}
		return nil
	}

	if err := calculateSize(path); err != nil {
		return 0, err
	}

	return float64(size), nil
}

func reconcileSecret(ctx context.Context, params Params, desiredSecret *corev1.Secret) error {

	if err := controllerutil.SetControllerReference(params.Instance, desiredSecret, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	currentSecret := corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: desiredSecret.Name, Namespace: params.Instance.Namespace}, &currentSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredSecret); err != nil {
			return err
		}
		params.Log.Info("created secret", "name", desiredSecret.Name, "namespace", params.Instance.Namespace)
	}
	if err != nil {
		return err
	}

	if desiredSecret.ObjectMeta.Annotations["checksum/data"] != currentSecret.ObjectMeta.Annotations["checksum/data"] {
		patch := client.MergeFrom(&currentSecret)
		if err := params.Client.Patch(ctx, desiredSecret, patch); err != nil {
			return err
		}
		params.Log.V(2).Info("secret updated", "name", desiredSecret.Name, "namespace", desiredSecret.Namespace)
	}

	return nil
}
