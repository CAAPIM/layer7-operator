package reconcile

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var repositoryScheduler = gocron.NewScheduler(time.Local).SingletonMode()

func Repositories(ctx context.Context, params Params) error {

	syncInterval := 5

	if params.Instance.Spec.App.RepositorySyncIntervalSeconds != 0 {
		syncInterval = params.Instance.Spec.App.RepositorySyncIntervalSeconds
	}

	err := registerRepositoryJob(ctx, params, syncInterval)
	if err != nil {
		params.Log.V(2).Info("jobs already registered", "detail", err.Error())
	}

	for _, j := range repositoryScheduler.Jobs() {
		for _, t := range j.Tags() {

			if t == "sync-repository-references" {
				if j.IsRunning() {
					return fmt.Errorf("already running: %w", errors.New("repository sync is already in progress"))
				}

				err := repositoryScheduler.RunByTag("sync-repository-references")
				repositoryScheduler.StartAsync()
				if err != nil {
					return fmt.Errorf("failed to reconcile repository: %w", err)
				}

			}
		}
	}
	return nil
}

func registerRepositoryJob(ctx context.Context, params Params, syncInterval int) error {
	if repositoryScheduler.Len() > 0 {
		_ = stopAndDeregister()
	}
	repositoryScheduler.TagsUnique()
	_, err := repositoryScheduler.Every(syncInterval).Seconds().Tag("sync-repository-references").Do(func() {
		cntr := 0
		for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
			if repoRef.Enabled && repoRef.Type == "dynamic" {
				cntr++
			}
		}
		if cntr == 0 {
			_ = stopAndDeregister()
		}

		for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
			if repoRef.Enabled && repoRef.Type == "dynamic" {
				err := reconcileDynamicRepository(ctx, params, repoRef)
				if err != nil {
					if k8serrors.IsNotFound(err) {
						params.Log.Error(err, "repository not found", "Name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
					}
					if k8serrors.IsForbidden(err) {
						params.Log.Error(err, "permission denied when trying to get repository, please check the Layer7 Operators Role", "Name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
					}
				}
			}
		}
	})

	if err != nil {
		return err
	}
	return nil
}

func stopAndDeregister() error {
	err := repositoryScheduler.RemoveByTag("sync-repository-references")
	if err != nil {
		return err
	}
	return nil
}

func reconcileDynamicRepository(ctx context.Context, params Params, repoRef securityv1.RepositoryReference) error {
	repository := &securityv1.Repository{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
	if err != nil && k8serrors.IsNotFound(err) {
		return err
	}
	commit := repository.Status.Commit

	switch repoRef.Type {
	case "dynamic":
		if !params.Instance.Spec.App.Management.Database.Enabled {
			err = applyEphemeral(ctx, params, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
			}
		} else {
			err = applyDbBacked(ctx, params, repoRef, commit)
			if err != nil {
				params.Log.Info("failed to apply commit", "Name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
				return err
			}
		}
	}
	return nil
}

func applyEphemeral(ctx context.Context, params Params, repoRef securityv1.RepositoryReference, commit string) error {

	graphmanPort := 9443

	if params.Instance.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = params.Instance.Spec.App.Management.Graphman.DynamicSyncPort
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
			for i := range repoRef.Directories {
				gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
				params.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Pod", pod.Name, "Name", params.Instance.Name, "Namespace", params.Instance.Namespace)
				name := params.Instance.Name
				if params.Instance.Spec.App.Management.SecretName != "" {
					name = params.Instance.Spec.App.Management.SecretName
				}
				gwSecret, err := getGatewaySecret(ctx, params, name)

				if err != nil {
					return err
				}
				err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
				if err != nil {
					return err
				}
			}

			if err := params.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "Failed to update pod label", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
				return err
			}
		}
	}

	return nil
}

func applyDbBacked(ctx context.Context, params Params, repoRef securityv1.RepositoryReference, commit string) error {
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
	for i := range repoRef.Directories {
		gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
		params.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Name", params.Instance.Name, "Namespace", params.Instance.Namespace)
		name := params.Instance.Name
		if params.Instance.Spec.App.Management.SecretName != "" {
			name = params.Instance.Spec.App.Management.SecretName
		}
		gwSecret, err := getGatewaySecret(ctx, params, name)

		if err != nil {
			return err
		}
		err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
		if err != nil {
			return err
		}
	}

	if err := params.Client.Patch(context.Background(), &gatewayDeployment,
		client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
		params.Log.Error(err, "Failed to update deployment annotations", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
		return err
	}

	return nil
}
