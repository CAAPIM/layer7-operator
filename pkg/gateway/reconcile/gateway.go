package reconcile

import (
	"context"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetGatewayPods returns the pods in a Gateway Deployment
func getGatewayPods(ctx context.Context, params Params) (*corev1.PodList, error) {
	podList := &corev1.PodList{}

	listOpts := []client.ListOption{
		client.InNamespace(params.Instance.Namespace),
		client.MatchingLabels(util.DefaultLabels(params.Instance.Name, map[string]string{})),
	}
	if err := params.Client.List(ctx, podList, listOpts...); err != nil {
		return podList, err
	}
	return podList, nil
}

func getGatewayDeployment(ctx context.Context, params Params) (appsv1.Deployment, error) {
	gatewayDeployment := &appsv1.Deployment{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gatewayDeployment)
	if err != nil {
		return *gatewayDeployment, err
	}
	return *gatewayDeployment, nil
}

func graphmanEncryptionPassphrase(ctx context.Context, params Params, repoRef securityv1.RepositoryReference) (string, error) {
	var graphmanEncryptionPassphrase string
	if repoRef.Encryption.Passphrase != "" && repoRef.Encryption.ExistingSecret == "" {
		graphmanEncryptionPassphrase = repoRef.Encryption.Passphrase
	} else {
		graphmanEncryptionSecret, err := getGatewaySecret(ctx, params, repoRef.Encryption.ExistingSecret)
		if err != nil {
			return "", err
		}
		graphmanEncryptionPassphrase = string(graphmanEncryptionSecret.Data[repoRef.Encryption.Key])
	}
	return graphmanEncryptionPassphrase, nil
}

func getGatewaySecret(ctx context.Context, params Params, name string) (*corev1.Secret, error) {
	gwSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, gwSecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return gwSecret, err
			}
		}
	}
	return gwSecret, nil
}
