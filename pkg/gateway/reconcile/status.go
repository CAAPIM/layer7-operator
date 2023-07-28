package reconcile

import (
	"context"
	"reflect"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func GatewayStatus(ctx context.Context, params Params) error {
	gatewayStatus := params.Instance.Status
	gatewayStatus.RepositoryStatus = []securityv1.GatewayRepositoryStatus{}
	gatewayStatus.Host = params.Instance.Spec.App.Management.Cluster.Hostname
	gatewayStatus.Image = params.Instance.Spec.App.Image
	gatewayStatus.Version = params.Instance.Spec.Version
	gatewayStatus.Gateway = []securityv1.GatewayState{}

	dep, err := getGatewayDeployment(ctx, params)
	if err != nil || k8serrors.IsNotFound(err) {
		params.Log.V(2).Info("deployment hasn't been created yet", "name", params.Instance.Name, "namespace", params.Instance.Namespace)

	} else {
		gatewayStatus.Replicas = dep.Status.Replicas
		gatewayStatus.Ready = dep.Status.ReadyReplicas
		gatewayStatus.State = corev1.PodInitialized
	}

	if dep.Status.ReadyReplicas == dep.Status.Replicas {
		gatewayStatus.State = corev1.PodReady
	}

	for _, repoRef := range params.Instance.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := params.Client.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: params.Instance.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.Info("repository not found", "name", params.Instance.Name, "repository", repoRef.Name, "namespace", params.Instance.Namespace)
			return err
		}

		secretName := repository.Name
		if repository.Spec.Auth.ExistingSecretName != "" {
			secretName = repository.Spec.Auth.ExistingSecretName
		}

		commit := repository.Status.Commit

		gatewayStatus.RepositoryStatus = append(gatewayStatus.RepositoryStatus, securityv1.GatewayRepositoryStatus{
			Commit:            commit,
			Enabled:           repoRef.Enabled,
			Name:              repoRef.Name,
			Type:              repoRef.Type,
			SecretName:        secretName,
			StorageSecretName: repository.Status.StorageSecretName,
			Branch:            repository.Spec.Branch,
			Endpoint:          repository.Spec.Endpoint,
		})
	}

	gatewayStatus.Conditions = dep.Status.Conditions
	podList, err := getGatewayPods(ctx, params)

	if err != nil {
		return err
	}

	ready := false
	for _, p := range podList.Items {
		for _, cs := range p.Status.ContainerStatuses {
			if cs.Image == params.Instance.Spec.App.Image {
				ready = cs.Ready
			}
		}

		gatewayState := securityv1.GatewayState{
			Name:  p.Name,
			Phase: p.Status.Phase,
			Ready: ready,
		}

		if p.Status.Phase == corev1.PodRunning {
			gatewayState.StartTime = p.Status.StartTime.String()
		}
		gatewayStatus.Gateway = append(gatewayStatus.Gateway, gatewayState)
	}

	if !reflect.DeepEqual(gatewayStatus, params.Instance.Status) {
		params.Instance.Status = gatewayStatus
		err = params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.Info("failed to update gateway status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
		}
	}
	return nil
}
