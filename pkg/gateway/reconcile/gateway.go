package reconcile

import (
	"context"

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

func getGraphmanEncryptionPassphrase(ctx context.Context, params Params, existingSecretName string, existingSecretKey string) (string, error) {
	graphmanEncryptionSecret, err := getGatewaySecret(ctx, params, existingSecretName)
	if err != nil {
		return "", err
	}
	return string(graphmanEncryptionSecret.Data[existingSecretKey]), nil
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

// getPodNames returns the pod names of the array of pods passed in
// func getPodNames(pods []corev1.Pod) []string {
// 	var podNames []string
// 	for _, pod := range pods {
// 		podNames = append(podNames, pod.Name)
// 	}
// 	return podNames
// }

func GatewayLicense(ctx context.Context, params Params) error {
	gatewayLicense := &corev1.Secret{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Spec.License.SecretName, Namespace: params.Instance.Namespace}, gatewayLicense)
	if k8serrors.IsNotFound(err) {
		params.Log.Error(err, "license not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func ManagementPod(ctx context.Context, params Params) error {

	podList, err := getGatewayPods(ctx, params)

	if err != nil {
		return err
	}

	for p := range podList.Items {
		if podList.Items[p].Labels["management-access"] == "leader" {
			if podList.Items[p].DeletionTimestamp == nil {
				return nil
			}
		}
	}

	tagged := false
	for p := range podList.Items {
		if podList.Items[p].Status.Phase == "Running" && podList.Items[p].DeletionTimestamp == nil && !tagged {
			patch := []byte(`{"metadata":{"labels":{"management-access": "leader"}}}`)
			if err := params.Client.Patch(context.Background(), &podList.Items[p],
				client.RawPatch(types.StrategicMergePatchType, patch)); err != nil {
				params.Log.Error(err, "failed to update pod label", "namespace", params.Instance.Namespace, "name", params.Instance.Name)
				return err
			}

			params.Log.V(2).Info("new leader elected", "name", params.Instance.Name, "pod", podList.Items[p].Name, "namespace", params.Instance.Namespace)
			tagged = true
		}
	}
	return nil
}
