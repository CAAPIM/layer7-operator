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
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	v1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func StateStorage(ctx context.Context, params Params, statestore securityv1alpha1.L7StateStore, commit string) error {
	storageSecretName, repositoryPath, _, err := localRepoStorageInfo(params)
	if err != nil {
		return err
	}

	tmpPath := "/tmp/statestore/" + params.Instance.Name
	commitTracker := commit + ".txt"
	fileName := "latest.json"
	_, dErr := os.Stat(tmpPath)
	if dErr != nil {
		_ = os.MkdirAll(tmpPath, 0755)
	}

	_, fErr := os.Stat(tmpPath + "/" + commitTracker)
	if fErr == nil {
		return nil
	}

	// Retrieve existing secret for Redis
	// this will need to be updated for multi-state store provider support
	if statestore.Spec.Redis.ExistingSecret != "" {
		stateStoreSecret, err := getStateStoreSecret(ctx, statestore.Spec.Redis.ExistingSecret, statestore, params)
		if err != nil {
			return err
		}
		statestore.Spec.Redis.Username = string(stateStoreSecret.Data["username"])
		statestore.Spec.Redis.MasterPassword = string(stateStoreSecret.Data["masterPassword"])
	}

	rc, err := util.RedisClient(&statestore.Spec.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to state store: %w", err)
	}

	projects, err := util.DetectGraphmanFolders(repositoryPath)
	if err != nil {
		return err
	}

	bundles := map[string][]byte{}
	compressedBundle := map[string][]byte{}
	for _, p := range projects {
		bundle, err := util.BuildAndValidateBundle(p, false)
		if err != nil {
			return err
		}
		bundleGzip, err := util.GzipCompress(bundle)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(bundleGzip)
		if err != nil {
			return err
		}

		if err := zw.Close(); err != nil {
			return err
		}

		keyName := strings.Replace(p, repositoryPath, "", 1)
		keyName = strings.Replace(strings.ReplaceAll(keyName, "/", "-"), "-", "", 1)
		bundles[keyName+".json"] = bundle
		compressedBundle[keyName+".gz"] = bundleGzip

		buf.Reset()
	}

	compressedBundleBytes, err := json.Marshal(compressedBundle)
	if err != nil {
		return err
	}

	// check for previous version - statestore may be empty
	stateStoreBundleMapString, err := rc.Get(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest").Result()
	if err != nil {
		// if the previous version can't be retrieved, write the current version
		rs := rc.Set(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest", compressedBundleBytes, 0)
		if rs.Err() != nil {
			return fmt.Errorf("failed to reconcile state storage: %w", rs.Err())
		}

		// then write that to file...
		err = os.WriteFile(tmpPath+"/"+fileName, compressedBundleBytes, 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(tmpPath+"/"+commitTracker, []byte{}, 0755)
		if err != nil {
			return err
		}

		return nil
	}

	stateStoreBundleMap := map[string][]byte{}
	err = json.Unmarshal([]byte(stateStoreBundleMapString), &stateStoreBundleMap)
	if err != nil {
		return err
	}

	for k1, b1bytes := range bundles {
		for k2, b2 := range stateStoreBundleMap {
			if strings.Split(k1, ".")[0] == strings.Split(k2, ".")[0] {
				b2Bytes, err := util.GzipDecompress(b2)
				if err != nil {
					//return fmt.Errorf("failed to decompress bundle 1: %s %w", k2, err)
					b2Bytes = b2
				}

				// reset current mappings from statestore
				// delete mappings should only persist for one version.
				srcBundle := graphman.Bundle{}
				destBundle := graphman.Bundle{}
				err = json.Unmarshal(b2Bytes, &srcBundle)
				if err != nil {
					return fmt.Errorf("failed to unmarshal state store bundle: %s %w", strings.Split(k2, ".")[0], err)
				}
				err = json.Unmarshal(b1bytes, &destBundle)
				if err != nil {
					return fmt.Errorf("failed to unmarshal local bundle: %s %w", strings.Split(k2, ".")[0], err)
				}

				// Reset mappings - removes entities marked for delete and clears all mappings
				// This represents the actual state after the last apply
				err = graphman.ResetMappings(&srcBundle)
				if err != nil {
					return fmt.Errorf("failed to reset mappings: %s %w", strings.Split(k2, ".")[0], err)
				}

				// Calculate delta: current=srcBundle (cleaned), desired=destBundle (filesystem)
				deltaBundle, combinedBundle, err := graphman.CalculateDelta(srcBundle, destBundle)
				if err != nil {
					return fmt.Errorf("failed to subtract current and previous version bundle: %s %w", strings.Split(k2, ".")[0], err)
				}

				deltaBundleBytes, err := json.Marshal(deltaBundle)
				if err != nil {
					return fmt.Errorf("failed to marshal delta bundle: %s %w", strings.Split(k1, ".")[0], err)
				}
				deltaGzip, err := util.GzipCompress(deltaBundleBytes)
				if err != nil {
					return fmt.Errorf("failed to compress delta bundle: %s %w", strings.Split(k1, ".")[0], err)
				}
				compressedBundle[strings.Split(k1, ".")[0]+"-delta.gz"] = deltaGzip

				// Store combined bundle (with delete mappings) for next iteration
				combinedBundleBytes, err := json.Marshal(combinedBundle)
				if err != nil {
					return fmt.Errorf("failed to marshal combined bundle: %s %w", strings.Split(k1, ".")[0], err)
				}
				combinedGzip, err := util.GzipCompress(combinedBundleBytes)
				if err != nil {
					return fmt.Errorf("failed to compress combined bundle: %w", err)
				}
				compressedBundle[strings.Split(k1, ".")[0]+".gz"] = combinedGzip
			}
		}
	}

	// prepare delete bundles for folders that have been removed
	for k2, b2 := range stateStoreBundleMap {
		found := false
		for k1 := range bundles {
			if strings.Split(k2, ".")[0] == strings.Split(k1, ".")[0] || strings.HasSuffix(k2, "-delta.gz") {
				found = true
			}
		}
		if !found {
			b2Bytes, err := util.GzipDecompress(b2)
			if err != nil {
				b2Bytes = b2
			}

			// reset current mappings from statestore
			// delete mappings should only persist for one version.
			srcBundle := graphman.Bundle{}
			err = json.Unmarshal(b2Bytes, &srcBundle)
			if err != nil {
				return fmt.Errorf("failed to unmarshal state store bundle: %s %w", strings.Split(k2, ".")[0], err)
			}

			// Reset mappings to get clean state (what's actually on gateway after last apply)
			err = graphman.ResetMappings(&srcBundle)
			if err != nil {
				return fmt.Errorf("failed to reset mappings: %s %w", strings.Split(k2, ".")[0], err)
			}

			// To delete all entities in srcBundle: current=srcBundle (cleaned), desired=empty
			// This generates delete mappings for everything.
			emptyBundle := graphman.Bundle{}
			deltaBundle, combinedBundle, err := graphman.CalculateDelta(srcBundle, emptyBundle)
			if err != nil {
				return fmt.Errorf("failed to subtract current and previous version bundle: %s %w", strings.Split(k2, ".")[0], err)
			}

			deltaBundleBytes, err := json.Marshal(deltaBundle)
			if err != nil {
				return fmt.Errorf("failed to marshal delta bundle: %s %w", strings.Split(k2, ".")[0], err)
			}

			deltaGzip, err := util.GzipCompress(deltaBundleBytes)
			if err != nil {
				return fmt.Errorf("failed to compress delta bundle: %s %w", strings.Split(k2, ".")[0], err)
			}
			compressedBundle[strings.Split(k2, ".")[0]+"-delta.gz"] = deltaGzip

			// Store combined bundle with delete mappings
			combinedBundleBytes, err := json.Marshal(combinedBundle)
			if err != nil {
				return fmt.Errorf("failed to marshal combined bundle: %s %w", strings.Split(k2, ".")[0], err)
			}
			combinedGzip, err := util.GzipCompress(combinedBundleBytes)
			if err != nil {
				return fmt.Errorf("failed to compress combined bundle: %w", err)
			}
			compressedBundle[strings.Split(k2, ".")[0]+".gz"] = combinedGzip
		}
	}

	compressedBundleBytes, err = json.Marshal(compressedBundle)
	if err != nil {
		return err
	}

	rs := rc.Set(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest", compressedBundleBytes, 0)
	if rs.Err() != nil {
		return fmt.Errorf("failed to reconcile state storage: %w", rs.Err())
	}

	err = os.WriteFile(tmpPath+"/"+fileName, compressedBundleBytes, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(tmpPath+"/"+commitTracker, []byte{}, 0755)
	if err != nil {
		return err
	}

	return nil
}

func GetStateStoreChecksum(ctx context.Context, params Params, statestore securityv1alpha1.L7StateStore) (commit string, err error) {
	if params.Instance.Spec.Type != v1.RepositoryTypeStateStore || params.Instance.Spec.StateStoreKey == "" {
		return "", fmt.Errorf("repository %s in namespace %s does not use a statestore or no statestore key defined, please check your repository configuration", params.Instance.Name, params.Instance.Namespace)
	}

	if statestore.Spec.Redis.ExistingSecret != "" {
		stateStoreSecret, err := getStateStoreSecret(ctx, statestore.Spec.Redis.ExistingSecret, statestore, params)
		if err != nil {
			return "", err
		}
		statestore.Spec.Redis.Username = string(stateStoreSecret.Data["username"])
		statestore.Spec.Redis.MasterPassword = string(stateStoreSecret.Data["masterPassword"])
	}

	rc, err := util.RedisClient(&statestore.Spec.Redis)
	if err != nil {
		return "", fmt.Errorf("failed to connect to state store: %w", err)
	}

	bundle, err := rc.Get(ctx, params.Instance.Spec.StateStoreKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve bundle from state store: %w", err)
	}

	h := sha1.New()
	h.Write([]byte(bundle))

	return fmt.Sprintf("%x", h.Sum(nil)), nil

}

func localRepoStorageInfo(params Params) (storageSecretName string, repositoryPath string, ext string, err error) {
	ext = params.Instance.Spec.Branch
	if ext == "" {
		ext = params.Instance.Spec.Tag
	}
	switch strings.ToLower(string(params.Instance.Spec.Type)) {
	case "http":
		fileURL, err := url.Parse(params.Instance.Spec.Endpoint)
		if err != nil {
			return "", "", "", err
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
		return storageSecretName, "/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + ext, ext, nil
	case "git":
		storageSecretName = params.Instance.Name + "-repository-" + ext
		return storageSecretName, "/tmp/" + params.Instance.Name + "-" + params.Instance.Namespace + "-" + ext, ext, nil
	default:
		params.Log.Info("repository type not set or unsupported", "name", params.Instance.Name, "namespace", params.Instance.Name)
		return "", "", "", errors.New("repository type not set or unsupported")
	}
}

func getStateStore(ctx context.Context, params Params) (securityv1alpha1.L7StateStore, error) {
	statestore := securityv1alpha1.L7StateStore{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Spec.StateStoreReference, Namespace: params.Instance.Namespace}, &statestore)
	if err != nil {
		return statestore, err
	}
	return statestore, nil
}

func getStateStoreSecret(ctx context.Context, name string, statestore securityv1alpha1.L7StateStore, params Params) (*corev1.Secret, error) {
	statestoreSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: statestore.Namespace}, statestoreSecret)
	if err != nil {
		return statestoreSecret, err
	}
	return statestoreSecret, nil
}
