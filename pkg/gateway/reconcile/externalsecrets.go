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
	"github.com/go-co-op/gocron"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var externalSecretsScheduler = gocron.NewScheduler(time.Local)
var secretSyncCache = util.NewSyncCache(3 * time.Second)

func ExternalSecrets(ctx context.Context, params Params) error {

	syncInterval := 5

	if params.Instance.Spec.App.ExternalSecretsSyncIntervalSeconds != 0 {
		syncInterval = params.Instance.Spec.App.ExternalSecretsSyncIntervalSeconds
	}

	err := registerExternalSecretJob(ctx, params, syncInterval)
	if err != nil {
		params.Log.V(2).Info("jobs already registered", "detail", err.Error())
	}

	for _, j := range externalSecretsScheduler.Jobs() {
		for _, t := range j.Tags() {

			if t == "sync-external-secrets" {
				if j.IsRunning() {
					//return fmt.Errorf("already running: %w", errors.New("external secret sync is already in progress"))
					params.Log.V(2).Info("external secret sync job is already in progress", "job", j.Tags(), "name", params.Instance.Name, "namespace", params.Instance.Namespace)
					return nil
				}

				err := externalSecretsScheduler.RunByTag("sync-external-secrets")

				if err != nil {
					return fmt.Errorf("failed to reconcile external secrets: %w", err)
				}

				externalSecretsScheduler.StartAsync()

			}
		}
	}

	return nil
}

func registerExternalSecretJob(ctx context.Context, params Params, syncInterval int) error {
	externalSecretsScheduler.TagsUnique()

	_, err := externalSecretsScheduler.Every(syncInterval).Seconds().Tag("sync-external-secrets").Do(func() {

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
			_ = externalSecretsScheduler.RemoveByTag("sync-external-secrets")
			return
		}

		err = reconcileExternalSecrets(ctx, params, gateway)
		if err != nil {
			params.Log.Info("failed to reconcile external secrets", "Name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
		}

	})

	if err != nil {
		return err
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
		syncRequest, err := secretSyncCache.Read(requestCacheEntry)
		tryRequest := true

		if err != nil {
			params.Log.V(2).Info("request has not been attempted or cache was flushed", "action", "sync external secrets", "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
		}

		if syncRequest.Attempts > 0 {
			params.Log.V(2).Info("request has been attempted in the last 30 seconds, backing off", "SecretSha1Sum", sha1Sum, "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)
			tryRequest = false
		}

		secretSyncCache.Update(util.SyncRequest{RequestName: requestCacheEntry, Attempts: 1}, time.Now().Add(30*time.Second).Unix())

		/// Change port to graphman port
		if tryRequest && ready {
			endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"

			params.Log.Info("Applying Latest Secret Bundle", "Secret SHA", sha1Sum, "Pod", pod.Name, "Name", gateway.Name, "Namespace", gateway.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
			if err != nil {
				return err
			}

			if err := params.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				params.Log.Error(err, "Failed to update pod label", "Namespace", gateway.Namespace, "Name", gateway.Name)
				return err
			}
		}
	}
	return nil
}
