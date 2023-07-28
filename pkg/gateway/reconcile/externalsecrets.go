package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
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

	opaqueSecretMap := []util.GraphmanSecret{}
	bundleBytes := []byte{}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	for _, es := range gateway.Spec.App.ExternalSecrets {
		if es.Enabled {

			secret, err := getGatewaySecret(ctx, params, es.Name)
			if err != nil {
				return err
			}

			if secret.Type == corev1.SecretTypeOpaque {
				for k, v := range secret.Data {
					opaqueSecretMap = append(opaqueSecretMap, util.GraphmanSecret{
						Name:   k,
						Secret: string(v),
					})
				}
			}

		}
	}

	if len(opaqueSecretMap) > 0 {
		bundleBytes, err = util.ConvertOpaqueMapToGraphmanBundle(opaqueSecretMap)
		if err != nil {
			return err
		}
	} else {
		return nil
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

	patch := fmt.Sprintf("{\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-secrets", sha1Sum)

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
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

		if pod.ObjectMeta.Annotations["security.brcmlabs.com/external-secrets"] == sha1Sum {
			return nil
		}

		requestCacheEntry := pod.Name + "-" + sha1Sum
		syncRequest, err := syncCache.Read(requestCacheEntry)
		tryRequest := true

		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync external secrets", "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "SecretSha1Sum", sha1Sum, "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
			tryRequest = false
		}
		if tryRequest && ready {
			syncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())

			endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"

			params.Log.V(2).Info("applying latest secret bundle", "sha1Sum", sha1Sum, "pod", pod.Name, "name", gateway.Name, "namespace", gateway.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
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
