package reconcile

import (
	"context"
	"fmt"
	"strconv"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func syncRepository(ctx context.Context, params Params) {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	cntr := 0
	for _, repoRef := range gateway.Spec.App.RepositoryReferences {
		if repoRef.Enabled && repoRef.Type == "dynamic" {
			cntr++
		}
	}
	if cntr == 0 {
		_ = removeJob("sync-repository-references")
	}

	err = params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	for _, repoRef := range gateway.Spec.App.RepositoryReferences {
		if repoRef.Enabled && repoRef.Type == "dynamic" {
			err := reconcileDynamicRepository(ctx, params, gateway, repoRef)
			if err != nil {
				params.Log.Error(err, "failed to reconcile repository reference", "name", gateway.Name, "repository", repoRef.Name, "namespace", gateway.Namespace)
			}
		}
	}
}

func reconcileDynamicRepository(ctx context.Context, params Params, gateway *securityv1.Gateway, repoRef securityv1.RepositoryReference) error {
	repository := &securityv1.Repository{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: gateway.Namespace}, repository)
	if err != nil && k8serrors.IsNotFound(err) {
		return err
	}

	commit := repository.Status.Commit

	switch repoRef.Type {
	case "dynamic":
		if !gateway.Spec.App.Management.Database.Enabled {
			err = applyEphemeral(ctx, params, gateway, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
			}
		} else {
			err = applyDbBacked(ctx, params, gateway, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
				return err
			}
		}
	}
	return nil
}

func applyEphemeral(ctx context.Context, params Params, gateway *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {
	graphmanPort := 9443

	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	graphmanEncryptionPassphrase, err := graphmanEncryptionPassphrase(ctx, params, repoRef)
	if err != nil {
		return err
	}

	for i, pod := range podList.Items {

		update := false
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		currentCommit := pod.ObjectMeta.Annotations["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]

		if currentCommit != commit || currentCommit == "" {
			update = true
		}

		if update && ready {
			endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"
			if len(repoRef.Directories) == 0 {
				repoRef.Directories = []string{"/"}
			}
			for d := range repoRef.Directories {
				gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[d]
				requestCacheEntry := pod.Name + "-" + repoRef.Name + "-" + commit
				syncRequest, err := syncCache.Read(requestCacheEntry)
				tryRequest := true
				if err != nil {
					params.Log.V(2).Info("request has not been attempted or cache was flushed", "repo", repoRef.Name, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
				}

				if syncRequest.Attempts > 0 {
					params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "repo", repoRef.Name, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					tryRequest = false
				}

				if tryRequest {
					syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
					name := gateway.Name
					if gateway.Spec.App.Management.SecretName != "" {
						name = gateway.Spec.App.Management.SecretName
					}
					gwSecret, err := getGatewaySecret(ctx, params, name)

					if err != nil {
						return err
					}

					params.Log.V(2).Info("applying latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)
					err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
					if err != nil {
						return err
					}

					params.Log.Info("applied latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)

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

func applyDbBacked(ctx context.Context, params Params, gateway *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {
	graphmanPort := 9443

	if params.Instance.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = params.Instance.Spec.App.Management.Graphman.DynamicSyncPort
	}

	patch := fmt.Sprintf("{\"spec\": { \"template\": {\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

	gatewayDeployment, err := getGatewayDeployment(ctx, params)

	graphmanEncryptionPassphrase, err := graphmanEncryptionPassphrase(ctx, params, repoRef)
	if err != nil {
		return err
	}

	if gatewayDeployment.Status.ReadyReplicas != gatewayDeployment.Status.Replicas {
		return nil
	}

	currentCommit := gatewayDeployment.Spec.Template.Annotations["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]
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
		gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[d]
		requestCacheEntry := gatewayDeployment.Name + "-" + repoRef.Name + "-" + commit
		syncRequest, err := syncCache.Read(requestCacheEntry)
		tryRequest := true
		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "repo", repoRef.Name, "deployment", gatewayDeployment.Name, "nam", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "Repo", repoRef.Name, "deployment", gatewayDeployment.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
			tryRequest = false
		}

		if tryRequest {
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())
			name := params.Instance.Name
			if params.Instance.Spec.App.Management.SecretName != "" {
				name = params.Instance.Spec.App.Management.SecretName
			}
			gwSecret, err := getGatewaySecret(ctx, params, name)

			if err != nil {
				return err
			}
			params.Log.Info("applying latest commit", "repo", repoRef.Name, "directory", repoRef.Directories[d], "commit", commit, "deployment", gatewayDeployment.Name, "name", gateway.Name, "namespace", gateway.Namespace)
			err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
			if err != nil {
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