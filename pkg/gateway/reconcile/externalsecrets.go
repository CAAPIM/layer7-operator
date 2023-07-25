package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var externalSecretsScheduler = gocron.NewScheduler(time.Local).SingletonMode()

func ExternalSecrets(ctx context.Context, params Params) error {

	syncInterval := 5

	if params.Instance.Spec.App.RepositorySyncIntervalSeconds != 0 {
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
					return fmt.Errorf("already running: %w", errors.New("repository sync is already in progress"))
				}

				err := externalSecretsScheduler.RunByTag("sync-external-secrets")
				externalSecretsScheduler.StartAsync()
				if err != nil {
					return fmt.Errorf("failed to reconcile repository: %w", err)
				}

			}
		}
	}

	return nil
}

func registerExternalSecretJob(ctx context.Context, params Params, syncInterval int) error {
	if externalSecretsScheduler.Len() > 0 {
		_ = externalSecretsScheduler.RemoveByTag("sync-external-secrets")
	}
	externalSecretsScheduler.TagsUnique()
	_, err := externalSecretsScheduler.Every(syncInterval).Seconds().Tag("sync-external-secrets").Do(func() {
		cntr := 0
		for _, externalSecret := range params.Instance.Spec.App.ExternalSecrets {
			if externalSecret.Enabled {
				cntr++
			}
		}
		if cntr == 0 {
			_ = externalSecretsScheduler.RemoveByTag("sync-external-secrets")
		}

		err := reconcileExternalSecrets(ctx, params)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				params.Log.Info("Secret not found", "Name", params.Instance.Name, "Namespace", params.Instance.Namespace, "External Secret Ref", params.Instance.Name)
			} else {
				params.Log.Info("Can't retrieve secret", "Name", params.Instance.Name, "Namespace", params.Instance.Namespace, "Error", err.Error())
			}
		}

	})

	if err != nil {
		return err
	}
	return nil
}

func reconcileExternalSecrets(ctx context.Context, params Params) error {
	opaqueSecretMap := []util.GraphmanSecret{}
	bundleBytes := []byte{}

	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		return err
	}

	for _, es := range params.Instance.Spec.App.ExternalSecrets {
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

		if ready && pod.ObjectMeta.Annotations["security.brcmlabs.com/external-secrets"] != sha1Sum {
			endpoint := pod.Status.PodIP + ":9443/graphman"

			params.Log.Info("Applying Latest Secret Bundle", "Secret SHA", sha1Sum, "Pod", pod.Name, "Name", params.Instance.Name, "Namespace", params.Instance.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
			if err != nil {
				return err
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
