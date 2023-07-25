/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gateway

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway/reconcile"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	Scheme   *runtime.Scheme
	muTasks  sync.Mutex
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("gateway", req.NamespacedName)
	gw, err := getGateway(r, ctx, req.NamespacedName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	params := reconcile.Params{
		Client:   r.Client,
		Recorder: r.Recorder,
		Scheme:   r.Scheme,
		Log:      log,
		Instance: gw,
	}

	err = checkGatewayLicense(r, ctx, gw)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	err = reconcile.Secret(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	if gw.Spec.App.Management.Service.Enabled {
		err = tagManagementPod(r, ctx, gw)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconcile.Services(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Ingress(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.HorizontalPodAutoscaler(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.PodDisruptionBudget(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = updateGatewayStatus(r, ctx, gw)
	if err != nil {
		r.Log.Error(err, "status update err")
		return ctrl.Result{}, err
	}

	err = reconcile.ConfigMaps(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Deployment(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.Repositories(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcile.ExternalSecrets(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func reconcileExternalKeys(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	keySecretMap := []util.GraphmanKey{}
	bundleBytes := []byte{}

	podList, err := getGatewayPods(r, ctx, gw)
	if err != nil {
		return err
	}

	for _, externalKey := range gw.Spec.App.ExternalKeys {
		if externalKey.Enabled {

			secret, err := getSecret(r, ctx, gw, externalKey.Name)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					r.Log.Info("Secret not found", "Name", gw.Name, "Namespace", gw.Namespace, "External Key Ref", externalKey.Name)
				} else {
					r.Log.Info("Can't retrieve secret", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
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
			r.Log.Info("Can't convert secrets to Graphman bundle", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
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

	name := gw.Name
	if gw.Spec.App.Management.SecretName != "" {
		name = gw.Spec.App.Management.SecretName
	}
	gwSecret, err := getSecret(r, ctx, gw, name)

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

			r.Log.Info("Applying Latest Secret Bundle", "Secret SHA", sha1Sum, "Pod", pod.Name, "Name", gw.Name, "Namespace", gw.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
			if err != nil {
				return err
			}

			if err := r.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
				return err
			}
		}
	}

	return nil
}

func reconcileExternalSecrets(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	opaqueSecretMap := []util.GraphmanSecret{}
	bundleBytes := []byte{}
	// TODO Confirm checksum changed before proceeding to get every referenced secret..

	//extSecretBytes, err := json.Marshal(gw.Spec.App.ExternalSecrets)

	podList, err := getGatewayPods(r, ctx, gw)
	if err != nil {
		return err
	}

	for _, externalSecret := range gw.Spec.App.ExternalSecrets {
		if externalSecret.Enabled {

			secret, err := getSecret(r, ctx, gw, externalSecret.Name)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					r.Log.Info("Secret not found", "Name", gw.Name, "Namespace", gw.Namespace, "External Secret Ref", externalSecret.Name)
				} else {
					r.Log.Info("Can't retrieve secret", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
				}
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
			r.Log.Info("Can't convert secrets to Graphman bundle", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
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

	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-secrets", sha1Sum)

	name := gw.Name
	if gw.Spec.App.Management.SecretName != "" {
		name = gw.Spec.App.Management.SecretName
	}
	gwSecret, err := getSecret(r, ctx, gw, name)

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

		if ready && pod.Labels["security.brcmlabs.com/external-secrets"] != sha1Sum {
			endpoint := pod.Status.PodIP + ":9443/graphman"

			r.Log.Info("Applying Latest Secret Bundle", "Secret SHA", sha1Sum, "Pod", pod.Name, "Name", gw.Name, "Namespace", gw.Namespace)

			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
			if err != nil {
				return err
			}

			if err := r.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
				return err
			}
		}
	}

	return nil
}

func updateGatewayStatus(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	gatewayStatus := gw.Status
	gatewayStatus.RepositoryStatus = []securityv1.GatewayRepositoryStatus{}
	gatewayStatus.Host = gw.Spec.App.Management.Cluster.Hostname
	gatewayStatus.Image = gw.Spec.App.Image
	gatewayStatus.Version = gw.Spec.Version
	gatewayStatus.Gateway = []securityv1.GatewayState{}

	dep, err := getGatewayDeployment(r, ctx, gw)
	if err != nil || k8serrors.IsNotFound(err) {
		//return nil
		r.Log.Info("Deployment hasn't been created yet", "Name", gw.Name, "Namespace", gw.Namespace)

	} else {
		gatewayStatus.Replicas = dep.Status.Replicas
		gatewayStatus.Ready = dep.Status.ReadyReplicas
		gatewayStatus.State = corev1.PodInitialized
	}

	if dep.Status.ReadyReplicas == dep.Status.Replicas {
		gatewayStatus.State = corev1.PodReady
	}

	for _, repoRef := range gw.Spec.App.RepositoryReferences {
		repository := &securityv1.Repository{}

		err := r.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: gw.Namespace}, repository)
		if err != nil && k8serrors.IsNotFound(err) {
			r.Log.Info("Repository not found", "Name", gw.Name, "Repository", repoRef.Name, "Namespace", gw.Namespace)
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
	podList, err := getGatewayPods(r, ctx, gw)

	if err != nil {
		return err
	}

	ready := false
	for _, p := range podList.Items {
		for _, cs := range p.Status.ContainerStatuses {
			if cs.Image == gw.Spec.App.Image {
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

	if !reflect.DeepEqual(gatewayStatus, gw.Status) {
		gw.Status = gatewayStatus
		err = r.Client.Status().Update(ctx, gw)
		if err != nil {
			r.Log.Info("Failed to update Gateway status", "Name", gw.Name, "Namespace", gw.Namespace, "Message", err.Error())
		}
	}
	return nil
}

func tagManagementPod(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	podList, err := getGatewayPods(r, ctx, gw)

	if err != nil {
		return err
	}

	podNames := getPodNames(podList.Items)
	if gw.Status.ManagementPod != "" {
		if util.Contains(podNames, gw.Status.ManagementPod) {
			return nil
		}
	}
	for p := range podList.Items {
		if p == 0 {
			patch := []byte(`{"metadata":{"labels":{"management-access": "leader"}}}`)
			if err := r.Client.Patch(context.Background(), &podList.Items[p],
				client.RawPatch(types.StrategicMergePatchType, patch)); err != nil {
				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
				return err
			}

			gw.Status.ManagementPod = podList.Items[0].Name
			if err := r.Client.Status().Update(ctx, gw); err != nil {
				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
				return err
			}
		}
	}
	return nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func getGatewayDeployment(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) (appsv1.Deployment, error) {
	gatewayDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, gatewayDeployment)
	if err != nil {
		return *gatewayDeployment, err
	}
	return *gatewayDeployment, nil
}

func getGatewayPods(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(gw.Namespace),
		client.MatchingLabels(util.DefaultLabels(gw.Name, map[string]string{})),
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		return podList, err
	}
	return podList, nil
}

// func getDeploymentStatus(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
// 	dep := appsv1.Deployment{}

// 	if err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, &dep); err != nil {
// 		return err
// 	}

// 	if dep.Status.AvailableReplicas == dep.Status.Replicas {
// 		return nil
// 	}

// 	return errors.New("Deployment not yet ready")
// }

func getSecret(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, name string) (*corev1.Secret, error) {
	gwSecret := &corev1.Secret{}

	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: gw.Namespace}, gwSecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			if err != nil {
				return gwSecret, err
			}
		}
	}
	return gwSecret, nil
}

func getGateway(r *GatewayReconciler, ctx context.Context, namespace types.NamespacedName) (*securityv1.Gateway, error) {
	gw := &securityv1.Gateway{}

	err := r.Get(ctx, namespace, gw)
	if err != nil {
		return gw, err
	}
	return gw, nil
}

// func applyGraphmanBundleEphemeral(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {

// 	notify := false
// 	notificationMessage := map[string]string{}
// 	applySuccess := []string{}
// 	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

// 	podList, err := getGatewayPods(r, ctx, gw)
// 	if err != nil {
// 		return err
// 	}

// 	graphmanEncryptionPassphrase, err := getGraphmanEncryptionPassphrase(r, ctx, gw, repoRef)
// 	if err != nil {
// 		return err
// 	}

// 	for i, pod := range podList.Items {
// 		update := false
// 		ready := false

// 		for _, containerStatus := range pod.Status.ContainerStatuses {
// 			if containerStatus.Name == "gateway" {
// 				ready = containerStatus.Ready
// 			}
// 		}

// 		currentCommit := pod.Labels["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]

// 		if currentCommit != commit || currentCommit == "" {
// 			update = true
// 		}

// 		if update && ready {
// 			notify = repoRef.Notification.Enabled
// 			endpoint := pod.Status.PodIP + ":9443/graphman"
// 			if len(repoRef.Directories) == 0 {
// 				repoRef.Directories = []string{"/"}
// 			}
// 			for i := range repoRef.Directories {
// 				gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
// 				r.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Pod", pod.Name, "Name", gw.Name, "Namespace", gw.Namespace)
// 				name := gw.Name
// 				if gw.Spec.App.Management.SecretName != "" {
// 					name = gw.Spec.App.Management.SecretName
// 				}
// 				gwSecret, err := getSecret(r, ctx, gw, name)

// 				if err != nil {
// 					return err
// 				}
// 				err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
// 				if err != nil {
// 					return err
// 				}
// 			}

// 			notificationMessage[pod.Name] = "successfully applied commit"
// 			applySuccess = append(applySuccess, pod.Name)

// 			if err := r.Client.Patch(context.Background(), &podList.Items[i],
// 				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
// 				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
// 				return err
// 			}
// 		}
// 	}

// 	// This is currently limited to WebHooks
// 	// When additional channels are added this will be refactored.
// 	if notify {
// 		notificationMessage["repository"] = repoRef.Name
// 		notificationMessage["commit"] = commit
// 		notificationMessage["deployment"] = gw.Name
// 		notificationMessage["text"] = "Repository: " + repoRef.Name + "\nCommit: " + commit + "\nSuccess: " + strings.Join(applySuccess, ",")
// 		notificationBytes, _ := json.Marshal(notificationMessage)
// 		_, err := util.RestCall("POST", repoRef.Notification.Channel.Webhook.Url, repoRef.Notification.Channel.Webhook.InsecureSkipVerify, repoRef.Notification.Channel.Webhook.Headers, "application/json", notificationBytes, repoRef.Notification.Channel.Webhook.Auth.Username, repoRef.Notification.Channel.Webhook.Auth.Password)

// 		if err != nil {
// 			r.Log.Info("Failed to send notification", "error", err.Error())
// 		}
// 	}
// 	return nil
// }

// func applyGraphmanBundleDbBacked(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {
// 	notify := false
// 	notificationMessage := map[string]string{}
// 	patch := fmt.Sprintf("{\"spec\": { \"template\": {\"metadata\": {\"annotations\": {\"%s\": \"%s\"}}}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

// 	gatewayDeployment, err := getGatewayDeployment(r, ctx, gw)

// 	graphmanEncryptionPassphrase, err := getGraphmanEncryptionPassphrase(r, ctx, gw, repoRef)
// 	if err != nil {
// 		return err
// 	}

// 	if gatewayDeployment.Status.ReadyReplicas != gatewayDeployment.Status.Replicas {
// 		return nil
// 	}

// 	currentCommit := gatewayDeployment.Spec.Template.Annotations["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]
// 	if currentCommit == commit {
// 		return nil
// 	}

// 	endpoint := gw.Name + "." + gw.Namespace + ".svc.cluster.local:8443/graphman"
// 	if gw.Spec.App.Management.Service.Enabled {
// 		endpoint = gw.Name + "-management-service." + gw.Namespace + ".svc.cluster.local:9443/graphman"
// 	}

// 	notify = repoRef.Notification.Enabled
// 	if len(repoRef.Directories) == 0 {
// 		repoRef.Directories = []string{"/"}
// 	}
// 	for i := range repoRef.Directories {
// 		gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
// 		r.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Name", gw.Name, "Namespace", gw.Namespace)
// 		name := gw.Name
// 		if gw.Spec.App.Management.SecretName != "" {
// 			name = gw.Spec.App.Management.SecretName
// 		}
// 		gwSecret, err := getSecret(r, ctx, gw, name)

// 		if err != nil {
// 			return err
// 		}
// 		err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
// 		if err != nil {
// 			return err
// 		}
// 		notificationMessage[gw.Name] = "successfully applied commit"
// 	}

// 	if err := r.Client.Patch(context.Background(), &gatewayDeployment,
// 		client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
// 		r.Log.Error(err, "Failed to update deployment annotations", "Namespace", gw.Namespace, "Name", gw.Name)
// 		return err
// 	}

// 	// This is currently limited to WebHooks
// 	// When additional channels are added this will be refactored.
// 	if notify {
// 		notificationMessage["repository"] = repoRef.Name
// 		notificationMessage["commit"] = commit
// 		notificationMessage["deployment"] = gw.Name
// 		notificationMessage["text"] = "Repository: " + repoRef.Name + "\nCommit: " + commit + "\nDeployment updated"
// 		notificationBytes, _ := json.Marshal(notificationMessage)
// 		_, err := util.RestCall("POST", repoRef.Notification.Channel.Webhook.Url, repoRef.Notification.Channel.Webhook.InsecureSkipVerify, repoRef.Notification.Channel.Webhook.Headers, "application/json", notificationBytes, repoRef.Notification.Channel.Webhook.Auth.Username, repoRef.Notification.Channel.Webhook.Auth.Password)

// 		if err != nil {
// 			r.Log.Info("Failed to send notification", "error", err.Error())
// 		}
// 	}
// 	return nil
// }

// func getGraphmanEncryptionPassphrase(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference) (string, error) {
// 	var graphmanEncryptionPassphrase string
// 	if repoRef.Encryption.Passphrase != "" && repoRef.Encryption.ExistingSecret == "" {
// 		graphmanEncryptionPassphrase = repoRef.Encryption.Passphrase
// 	} else {
// 		graphmanEncryptionSecret, err := getSecret(r, ctx, gw, repoRef.Encryption.ExistingSecret)
// 		if err != nil {
// 			return "", err
// 		}
// 		graphmanEncryptionPassphrase = string(graphmanEncryptionSecret.Data[repoRef.Encryption.Key])
// 	}
// 	return graphmanEncryptionPassphrase, nil
// }

func checkGatewayLicense(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	gatewayLicense := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: gw.Spec.License.SecretName, Namespace: gw.Namespace}, gatewayLicense)
	if k8serrors.IsNotFound(err) {
		r.Log.Error(err, "License not found", "Name", gw.Name, "Namespace", gw.Namespace)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.Gateway{}).
		Complete(r)
}
