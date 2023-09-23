package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func syncExternalKeys(ctx context.Context, params Params) {
	cntr := 0
	for _, externalSecret := range params.Instance.Spec.App.ExternalSecrets {
		if externalSecret.Enabled {
			cntr++
		}
	}
	if cntr == 0 {
		_ = s.RemoveByTag(params.Instance.Name + "-sync-external-keys")
	}

	err := reconcileExternalKeys(ctx, params)
	if err != nil {
		params.Log.Info("failed to reconcile external keys", "Name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
	}

}

func reconcileExternalKeys(ctx context.Context, params Params) error {
	keySecretMap := []util.GraphmanKey{}
	bundleBytes := []byte{}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	for _, externalKey := range params.Instance.Spec.App.ExternalKeys {
		if externalKey.Enabled {

			secret, err := getGatewaySecret(ctx, params, externalKey.Name)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					params.Log.Info("secret not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "external key ref", externalKey.Name)
				} else {
					params.Log.Info("can't retrieve secret", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
				}
			}

			if secret.Type == corev1.SecretTypeTLS {
				keySecretMap = append(keySecretMap, util.GraphmanKey{
					Name: secret.Name,
					Crt:  string(secret.Data["tls.crt"]),
					Key:  string(secret.Data["tls.key"]),
				})
			}

		}
	}

	if len(keySecretMap) > 0 {
		bundleBytes, err = util.ConvertX509ToGraphmanBundle(keySecretMap)
		if err != nil {
			params.Log.Info("can't convert secrets to graphman bundle", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
		}
	} else {
		return nil
	}

	sort.Slice(keySecretMap, func(i, j int) bool {
		return keySecretMap[i].Name < keySecretMap[j].Name
	})

	keySecretMapBytes, err := json.Marshal(keySecretMap)

	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(keySecretMapBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-keys", sha1Sum)

	name := params.Instance.Name
	if params.Instance.Spec.App.Management.SecretName != "" {
		name = params.Instance.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	for i, pod := range podList.Items {
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		if ready && pod.Labels["security.brcmlabs.com/external-keys"] != sha1Sum {
			endpoint := pod.Status.PodIP + ":9443/graphman"

			params.Log.V(2).Info("applying latest key bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
			if err != nil {
				return err
			}

			params.Log.Info("applied latest key bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)

			if err := params.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "Failed to update pod label", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
				return err
			}
		}
	}
	return nil
}
