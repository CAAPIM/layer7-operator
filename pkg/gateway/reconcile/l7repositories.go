package reconcile

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ExternalRepository(ctx context.Context, params Params) error {
	gateway := params.Instance
	for _, repoRef := range gateway.Spec.App.RepositoryReferences {
		if !repoRef.Enabled {
			return nil
		}
	}
	for _, repoRef := range gateway.Spec.App.RepositoryReferences {
		if repoRef.Enabled {
			err := reconcileDynamicRepository(ctx, params, repoRef)
			if err != nil {
				params.Log.Error(err, "failed to reconcile repository reference", "name", gateway.Name, "repository", repoRef.Name, "namespace", gateway.Namespace)
				return err
			}
		}
	}

	return nil
}

func reconcileDynamicRepository(ctx context.Context, params Params, repoRef securityv1.RepositoryReference) error {
	gateway := params.Instance
	repository := &securityv1.Repository{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: gateway.Namespace}, repository)
	if err != nil && k8serrors.IsNotFound(err) {
		return err
	}

	if !repository.Status.Ready {
		params.Log.Info("repository not ready", "repository", repository.Name, "name", gateway.Name, "namespace", gateway.Namespace)
		return nil
	}

	commit := repository.Status.Commit

	switch repoRef.Type {
	case "dynamic":
		if !gateway.Spec.App.Management.Database.Enabled {
			err = applyEphemeral(ctx, params, repository, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
				return err
			}

		} else {
			err = applyDbBacked(ctx, params, repository, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
				return err
			}
		}
	}

	for _, sRepo := range gateway.Status.RepositoryStatus {
		if sRepo.Name == repoRef.Name {
			if sRepo.Commit != commit {
				_ = GatewayStatus(ctx, params)
			}
		}
	}

	return nil
}

func applyEphemeral(ctx context.Context, params Params, repository *securityv1.Repository, repoRef securityv1.RepositoryReference, commit string) error {
	gateway := params.Instance

	graphmanPort := 9443

	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	graphmanEncryptionPassphrase := repoRef.Encryption.Passphrase

	if repoRef.Encryption.ExistingSecret != "" {
		graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, repoRef.Encryption.ExistingSecret, repoRef.Encryption.Key)
		if err != nil {
			return err
		}
	}

	singleton := false
	if !gateway.Spec.App.SingletonExtraction {
		singleton = true
	}

	for i, pod := range podList.Items {

		update := false
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}
		latestCommit := commit
		currentCommit := pod.ObjectMeta.Annotations["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]

		if gateway.Spec.App.SingletonExtraction {
			if pod.ObjectMeta.Labels["management-access"] == "leader" {
				latestCommit = commit + "-leader"
				singleton = true
			}
		}

		patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, latestCommit)

		if currentCommit != latestCommit || currentCommit == "" {
			update = true
		}

		if update && ready {
			endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"
			if len(repoRef.Directories) == 0 {
				repoRef.Directories = []string{"/"}
			}
			for d := range repoRef.Directories {
				ext := repository.Spec.Branch
				if ext == "" {
					ext = repository.Spec.Tag
				}

				gitPath := "/tmp/" + repoRef.Name + "-" + gateway.Namespace + "-" + ext + "/" + repoRef.Directories[d]

				if strings.ToLower(repository.Spec.Type) == "http" {
					fileURL, err := url.Parse(repository.Spec.Endpoint)
					if err != nil {
						log.Fatal(err)
					}
					path := fileURL.Path
					segments := strings.Split(path, "/")
					fileName := segments[len(segments)-1]
					ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]
					folderName := strings.ReplaceAll(fileName, "."+ext, "")
					if ext == "gz" && strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
						folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
					}
					gitPath = "/tmp/" + repository.Name + "-" + gateway.Namespace + "-" + folderName

				}

				requestCacheEntry := pod.Name + "-" + repoRef.Name + "-" + latestCommit
				syncRequest, err := syncCache.Read(requestCacheEntry)
				tryRequest := true
				if err != nil {
					params.Log.V(2).Info("request has not been attempted or cache was flushed", "repo", repoRef.Name, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
				}

				if syncRequest.Attempts > 0 {
					params.Log.V(2).Info("request has been attempted in the last 3 seconds, backing off", "repo", repoRef.Name, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					tryRequest = false
					return errors.New("request has been attempted in the last 3 seconds, backing off")
				}

				if tryRequest {
					syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())
					name := gateway.Name
					if gateway.Spec.App.Management.SecretName != "" {
						name = gateway.Spec.App.Management.SecretName
					}
					gwSecret, err := getGatewaySecret(ctx, params, name)

					if err != nil {
						return err
					}
					start := time.Now()
					params.Log.V(2).Info("applying latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", latestCommit, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					err = util.ApplyToGraphmanTarget(gitPath, singleton, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
					if err != nil {
						params.Log.Info("failed to apply latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", latestCommit, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
						_ = captureGraphmanMetrics(ctx, params, start, pod.Name, "repository", repoRef.Name, latestCommit, true)
						return err
					}

					params.Log.Info("applied latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", latestCommit, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					_ = captureGraphmanMetrics(ctx, params, start, pod.Name, "repository", repoRef.Name, latestCommit, false)

					if err := params.Client.Patch(context.Background(), &podList.Items[i],
						client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
						params.Log.Error(err, "failed to update pod label", "Name", gateway.Name, "namespace", gateway.Namespace)
						return err
					}
				}
			}
		}
	}
	return nil
}

func applyDbBacked(ctx context.Context, params Params, repository *securityv1.Repository, repoRef securityv1.RepositoryReference, commit string) error {
	gateway := params.Instance

	graphmanPort := 9443

	if params.Instance.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = params.Instance.Spec.App.Management.Graphman.DynamicSyncPort
	}

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	graphmanEncryptionPassphrase := repoRef.Encryption.Passphrase

	if repoRef.Encryption.ExistingSecret != "" {
		graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, repoRef.Encryption.ExistingSecret, repoRef.Encryption.Key)
		if err != nil {
			return err
		}
	}

	if gatewayDeployment.Status.ReadyReplicas != gatewayDeployment.Status.Replicas {
		return nil
	}

	currentCommit := gatewayDeployment.ObjectMeta.Annotations["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]
	if currentCommit == commit {
		return nil
	}

	endpoint := params.Instance.Name + "." + params.Instance.Namespace + ".svc.cluster.local:" + strconv.Itoa(graphmanPort) + "/graphman"
	if params.Instance.Spec.App.Management.Service.Enabled {
		endpoint = params.Instance.Name + "-management-service." + params.Instance.Namespace + ".svc.cluster.local:9443/graphman"
	}

	if len(repoRef.Directories) == 0 {
		repoRef.Directories = []string{"/"}
	}
	for d := range repoRef.Directories {
		ext := repository.Spec.Branch

		if ext == "" {
			ext = repository.Spec.Tag
		}
		gitPath := "/tmp/" + repoRef.Name + "-" + gateway.Namespace + "-" + ext + "/" + repoRef.Directories[d]
		requestCacheEntry := gatewayDeployment.Name + "-" + repoRef.Name + "-" + commit
		syncRequest, err := syncCache.Read(requestCacheEntry)
		tryRequest := true
		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "repo", repoRef.Name, "deployment", gatewayDeployment.Name, "nam", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 3 seconds, backing off", "Repo", repoRef.Name, "deployment", gatewayDeployment.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
			tryRequest = false
		}

		if tryRequest {
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(3*time.Second).Unix())
			name := params.Instance.Name
			if params.Instance.Spec.App.Management.SecretName != "" {
				name = params.Instance.Spec.App.Management.SecretName
			}
			gwSecret, err := getGatewaySecret(ctx, params, name)

			if err != nil {
				return err
			}
			params.Log.V(2).Info("applying latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "deployment", gatewayDeployment.Name, "name", gateway.Name, "namespace", gateway.Namespace)
			err = util.ApplyToGraphmanTarget(gitPath, true, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
			if err != nil {
				params.Log.Info("failed to apply latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "deployment", gatewayDeployment.Name, "name", gateway.Name, "namespace", gateway.Namespace)
				return err
			}
			params.Log.Info("applied latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "deployment", gatewayDeployment.Name, "name", gateway.Name, "namespace", gateway.Namespace)

			if err := params.Client.Patch(context.Background(), &gatewayDeployment,
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "Failed to update deployment annotations", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
				return err
			}
		}
	}

	return nil
}
