package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func syncExternalSecrets(ctx context.Context, params Params) error {

	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "Name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	cntr := 0
	for _, externalSecret := range gateway.Spec.App.ExternalSecrets {
		if externalSecret.Enabled {
			cntr++
		}
	}
	if cntr == 0 {
		_ = s.RemoveByTag("sync-external-secrets")
		return nil
	}

	err = reconcileExternalSecrets(ctx, params, gateway)
	if err != nil {
		params.Log.Info("failed to reconcile external secrets", "Name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
	}

	return nil
}

func reconcileExternalSecrets(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	graphmanPort := 9443

	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	for _, es := range gateway.Spec.App.ExternalSecrets {
		opaqueSecretMap := []util.GraphmanSecret{}
		if es.Enabled {

			secret, err := getGatewaySecret(ctx, params, es.Name)
			if err != nil {
				return err
			}

			switch secret.Type {
			case corev1.SecretTypeOpaque:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 k,
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			case corev1.SecretTypeBasicAuth:
			case corev1.SecretTypeServiceAccountToken:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 es.Name + "-" + k,
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			case corev1.SecretTypeDockerConfigJson:
			case corev1.SecretTypeDockercfg:
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:                 es.Name + "-" + strings.Split(k, ".")[1],
						Secret:               string(v),
						Description:          es.Description,
						VariableReferencable: es.VariableReferencable,
					})
				}
			default:
				params.Log.V(2).Info("not a supported secret type", "secret name", es.Name, "secret type", secret.Type, "name", gateway.Name, "namespace", gateway.Namespace)
			}
		}

		if len(opaqueSecretMap) < 1 {
			continue
		}

		bundleBytes, err := util.ConvertOpaqueMapToGraphmanBundle(opaqueSecretMap)
		if err != nil {
			return err
		}
		sort.Slice(opaqueSecretMap, func(i, j int) bool {
			return opaqueSecretMap[i].Name < opaqueSecretMap[j].Name
		})

		opaqueSecretMapBytes, err := json.Marshal(opaqueSecretMap)
		if err != nil {
			return err
		}
		h := sha1.New()
		h.Write(opaqueSecretMapBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		graphmanEncryptionPassphrase := es.Encryption.Passphrase
		if es.Encryption.ExistingSecret != "" {
			graphmanEncryptionPassphrase, err = getGraphmanEncryptionPassphrase(ctx, params, es.Encryption.ExistingSecret, es.Encryption.Key)
			if err != nil {
				return err
			}
		}

		if !gateway.Spec.App.Management.Database.Enabled {
			err = applyExternalSecretEphemeral(ctx, params, gateway, *gwSecret, bundleBytes, es.Name, graphmanEncryptionPassphrase, sha1Sum, graphmanPort)
			if err != nil {
				params.Log.Info("failed to apply secret bundle", "sha1Sum", sha1Sum, "secret name", es.Name, "name", gateway.Name, "namespace", gateway.Namespace)
			}
		} else {
			err = applyExternalSecretDbBacked(ctx, params, gateway, *gwSecret, bundleBytes, es.Name, graphmanEncryptionPassphrase, sha1Sum, graphmanPort)
			if err != nil {
				params.Log.Info("failed to apply secret bundle", "sha1Sum", sha1Sum, "secret name", es.Name, "name", gateway.Name, "namespace", gateway.Namespace)
			}
		}
	}

	return nil
}

func applyExternalSecretEphemeral(ctx context.Context, params Params, gateway *securityv1.Gateway, gwSecret corev1.Secret, bundleBytes []byte, name string, graphmanEncryptionPassphrase string, sha1Sum string, graphmanPort int) error {

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-secret-"+name, sha1Sum)

	for i, pod := range podList.Items {
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		if pod.ObjectMeta.Annotations["security.brcmlabs.com/external-secret-"+name] == sha1Sum {
			return nil
		}

		if ready {
			requestCacheEntry := pod.Name + "-" + sha1Sum
			syncRequest, err := syncCache.Read(requestCacheEntry)
			if err != nil {
				params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync external secrets", "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
			}

			if syncRequest.Attempts > 0 {
				params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "SecretSha1Sum", sha1Sum, "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
				return nil
			}
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())

			endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"

			params.Log.V(2).Info("applying latest secret bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase, bundleBytes)
			if err != nil {
				return err
			}
			params.Log.Info("applied latest secret bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)

			if err := params.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "Failed to update pod label", "Namespace", gateway.Namespace, "Name", gateway.Name)
				return err
			}
		}
	}
	return nil

}

func applyExternalSecretDbBacked(ctx context.Context, params Params, gateway *securityv1.Gateway, gwSecret corev1.Secret, bundleBytes []byte, name string, graphmanEncryptionPassphrase string, sha1Sum string, graphmanPort int) error {

	gatewayDeployment, err := getGatewayDeployment(ctx, params)
	if err != nil {
		return err
	}

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-secret-"+name, sha1Sum)

	ready := false

	if gatewayDeployment.ObjectMeta.Annotations["security.brcmlabs.com/external-secret-"+name] == sha1Sum {
		return nil
	}

	if gatewayDeployment.Status.ReadyReplicas == gatewayDeployment.Status.Replicas {
		ready = true
	}

	if ready {
		requestCacheEntry := gatewayDeployment.Name + "-" + sha1Sum
		syncRequest, err := syncCache.Read(requestCacheEntry)
		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync external secrets", "Name", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "SecretSha1Sum", sha1Sum, "Name", gateway.Name, "Namespace", gateway.Namespace)
			return nil
		}
		syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())

		endpoint := gateway.Name + "." + gateway.Namespace + ".svc.cluster.local:" + strconv.Itoa(graphmanPort) + "/graphman"
		if gateway.Spec.App.Management.Service.Enabled {
			endpoint = gateway.Name + "-management-service." + gateway.Namespace + ".svc.cluster.local:" + strconv.Itoa(graphmanPort) + "/graphman"
		}
		params.Log.V(2).Info("applying latest secret bundle", "sha1Sum", sha1Sum, "name", gateway.Name, "namespace", gateway.Namespace)

		err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase, bundleBytes)
		if err != nil {
			return err
		}
		params.Log.Info("applied latest secret bundle", "sha1Sum", sha1Sum, "name", gateway.Name, "namespace", gateway.Namespace)

		if err := params.Client.Patch(context.Background(), &gatewayDeployment,
			client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
			params.Log.Error(err, "Failed to update deployment annotations", "Namespace", params.Instance.Namespace, "Name", params.Instance.Name)
			return err
		}
	}
	// }
	return nil

}
