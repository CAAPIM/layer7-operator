package reconcile

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"strings"

	v1 "github.com/caapim/layer7-operator/api/v1"
	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func StateStorage(ctx context.Context, params Params, statestore securityv1alpha1.L7StateStore) error {
	storageSecretName, repositoryPath, _, err := localRepoStorageInfo(params)
	if err != nil {
		return err
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

	rc := util.RedisClient(&statestore.Spec.Redis)
	version := params.Instance.Status.StateStoreVersion

	// this represents the current version
	bundleGzip, err := util.CompressGraphmanBundle(repositoryPath)

	if err != nil {
		return err
	}

	if version == 0 {
		rs := rc.Set(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest", bundleGzip, 0)
		if rs.Err() != nil {
			return fmt.Errorf("failed to reconcile state storage: %w", rs.Err())
		}
		return nil
	}

	// calculate delta
	currentVersion, err := util.GzipDecompress(bundleGzip)
	if err != nil {
		return fmt.Errorf("failed to decompress bundle: %w", err)
	}
	previousVersionGzip, err := rc.Get(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest").Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve previous version bundle from state store: %w", err)
	}
	previousVersion, err := util.GzipDecompress([]byte(previousVersionGzip))
	if err != nil {
		return fmt.Errorf("failed to decompress previous version bundle: %w", err)
	}
	deltaBytes, combinedBytes, err := graphman.SubtractBundle(previousVersion, currentVersion)
	if err != nil {
		return fmt.Errorf("failed to subtract current and previous version bundles: %w", err)
	}
	combinedGzip, err := util.GzipCompress(combinedBytes)
	if err != nil {
		return fmt.Errorf("failed to compress combined bundle: %w", err)
	}
	deltaGzip, err := util.GzipCompress(deltaBytes)
	if err != nil {
		return fmt.Errorf("failed to compress delta bundle: %w", err)
	}

	rs := rc.Set(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":latest", combinedGzip, 0)
	if rs.Err() != nil {
		return fmt.Errorf("failed to reconcile state storage: %w", rs.Err())
	}

	rs = rc.Set(ctx, statestore.Spec.Redis.GroupName+":"+statestore.Spec.Redis.StoreId+":"+"repository"+":"+storageSecretName+":delta", deltaGzip, 0)
	if rs.Err() != nil {
		return fmt.Errorf("failed to reconcile state storage: %w", rs.Err())

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

	rc := util.RedisClient(&statestore.Spec.Redis)

	bundle, err := rc.Get(ctx, params.Instance.Spec.StateStoreKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to bundle from state store: %w", err)
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
