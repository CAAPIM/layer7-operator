package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncExternalKeys(ctx context.Context, params Params) error {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "Name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = s.RemoveByTag(params.Instance.Name + "-sync-external-keys")
		return nil
	}
	cntr := 0
	for _, externalKey := range gateway.Spec.App.ExternalKeys {
		if externalKey.Enabled {
			cntr++
		}
	}
	if cntr == 0 {
		_ = s.RemoveByTag(params.Instance.Name + "-sync-external-keys")
	}

	err = reconcileExternalKeys(ctx, params, gateway)
	if err != nil {
		params.Log.Info("failed to reconcile external keys", "Name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
	}
	return nil
}

func reconcileExternalKeys(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	keySecretMap := []util.GraphmanKey{}
	bundleBytes := []byte{}

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	for _, externalKey := range gateway.Spec.App.ExternalKeys {
		if externalKey.Enabled {

			secret, err := getGatewaySecret(ctx, params, externalKey.Name)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					params.Log.Info("secret not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "external key ref", externalKey.Name)
					continue
				} else {
					params.Log.Info("can't retrieve secret", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
					continue
				}
			}

			usageType := ""
			switch strings.ToUpper(string(externalKey.KeyUsageType)) {
			case string(securityv1.KeyUsageTypeDefaultSSL), string(securityv1.KeyUsageTypeDefaultCA), string(securityv1.KeyUsageAuditSigning), string(securityv1.KeyUsageAuditViewer):
				usageType = strings.ToUpper(string(externalKey.KeyUsageType))
			}

			if secret.Type == corev1.SecretTypeTLS {
				keySecretMap = append(keySecretMap, util.GraphmanKey{
					Name:      secret.Name,
					Crt:       string(secret.Data["tls.crt"]),
					Key:       string(secret.Data["tls.key"]),
					Alias:     externalKey.Alias,
					UsageType: usageType,
				})
			}

		}

		if len(keySecretMap) > 0 {
			bundleBytes, err = util.ConvertX509ToGraphmanBundle(keySecretMap)
			if err != nil {
				return err
				//params.Log.Info("can't convert secrets to graphman bundle", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "error", err.Error())
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

		//patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-keys", sha1Sum)

		annotation := "security.brcmlabs.com/external-key-" + externalKey.Name

		if !gateway.Spec.App.Management.Database.Enabled {
			err = ReconcileEphemeralGateway(ctx, params, "external keys", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, bundleBytes)
			if err != nil {
				return err
			}
		} else {
			err = ReconcileDBGateway(ctx, params, "external keys", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, bundleBytes)
			if err != nil {
				return err
			}
		}

	}

	// name := params.Instance.Name
	// if params.Instance.Spec.App.Management.SecretName != "" {
	// 	name = params.Instance.Spec.App.Management.SecretName
	// }
	// gwSecret, err := getGatewaySecret(ctx, params, name)

	// if err != nil {
	// 	return err
	// }

	// for i, pod := range podList.Items {
	// 	ready := false

	// 	for _, containerStatus := range pod.Status.ContainerStatuses {
	// 		if containerStatus.Name == "gateway" {
	// 			ready = containerStatus.Ready
	// 		}
	// 	}

	// 	if ready && pod.Labels["security.brcmlabs.com/external-keys"] != sha1Sum {
	// 		endpoint := pod.Status.PodIP + ":9443/graphman"

	// 		params.Log.V(2).Info("applying latest key bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)

	// 		err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		params.Log.Info("applied latest key bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)

	// 		if err := params.Client.Patch(context.Background(), &podList.Items[i],
	// 			client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
	// 			params.Log.Error(err, "Failed to update pod label", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
	// 			return err
	// 		}
	// 	}
	// }
	return nil
}
