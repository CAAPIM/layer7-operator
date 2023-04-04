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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/gateway"
	"github.com/caapim/layer7-operator/pkg/gateway/config"
	"github.com/caapim/layer7-operator/pkg/gateway/hpa"
	"github.com/caapim/layer7-operator/pkg/gateway/ingress"
	"github.com/caapim/layer7-operator/pkg/gateway/monitoring"
	"github.com/caapim/layer7-operator/pkg/gateway/secrets"
	"github.com/caapim/layer7-operator/pkg/gateway/service"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	otelv1alpha1 "github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = r.Log.WithValues("gateway", req.NamespacedName)

	gw, err := getGateway(r, ctx, req.NamespacedName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = checkGatewayLicense(r, ctx, gw)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = reconcileServiceMonitor(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileOtelCollector(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileConfigMap(r, gw.Name, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileConfigMap(r, gw.Name+"-system", ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	if gw.Spec.App.ClusterProperties.Enabled {
		err = reconcileConfigMap(r, gw.Name+"-cwp-bundle", ctx, gw)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if gw.Spec.App.ListenPorts.Harden || gw.Spec.App.ListenPorts.Custom.Enabled {
		err = reconcileConfigMap(r, gw.Name+"-listen-port-bundle", ctx, gw)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconcileSecret(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileService(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	if gw.Spec.App.Management.Service.Enabled {
		err = reconcileManagementService(r, ctx, gw)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = tagManagementPod(r, ctx, gw)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconcileIngress(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileHPA(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = updateGatewayStatus(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileConfigMap(r, gw.Name+"-repository-init-config", ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileDeployment(r, ctx, gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, repository := range gw.Spec.App.RepositoryReferences {
		if repository.Enabled {
			err = reconcileDynamicRepository(r, ctx, gw, repository)
			if err != nil {
				if k8serrors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

func reconcileDynamicRepository(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference) error {
	repository := &securityv1.Repository{}

	err := r.Get(ctx, types.NamespacedName{Name: repoRef.Name, Namespace: gw.Namespace}, repository)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Repository not found", "Name", gw.Name, "Repository", repoRef.Name, "Namespace", gw.Namespace)
		return err
	}
	commit := repository.Status.Commit

	switch repoRef.Type {
	case "dynamic":
		if !gw.Spec.App.Management.Database.Enabled {
			err = applyGraphmanBundleEphemeral(r, ctx, gw, repoRef, commit)
			if err != nil {
				r.Log.Info("Failed to apply commit", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
			}
		} else {
			err = applyGraphmanBundleDbBacked(r, ctx, gw, repoRef, commit)
			if err != nil {
				r.Log.Info("Failed to apply commit", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
				return err
			}
		}

	}
	return nil
}

func reconcileHPA(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	if !gw.Spec.App.Autoscaling.Enabled {
		return nil
	}
	currHPA := &autoscalingv2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currHPA)
	newHpa := hpa.NewHPA(gw)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating HPA", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, newHpa, r.Scheme)
		err = r.Create(ctx, newHpa)
		if err != nil {
			r.Log.Error(err, "Failed creating HPA", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(currHPA, newHpa) {
		ctrl.SetControllerReference(gw, newHpa, r.Scheme)
		return r.Update(ctx, newHpa)
	}
	return nil
}

func reconcileConfigMap(r *GatewayReconciler, name string, ctx context.Context, gw *securityv1.Gateway) error {
	currMap := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: gw.Namespace}, currMap)
	cm := config.NewConfigMap(gw, name)

	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating ConfigMap", "Name", name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, cm, r.Scheme)
		err = r.Create(ctx, cm)
		if err != nil {
			r.Log.Error(err, "Failed creating ConfigMap", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(currMap.Data, cm.Data) {
		ctrl.SetControllerReference(gw, cm, r.Scheme)
		return r.Update(ctx, cm)
	}
	return nil
}

func reconcileSecret(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {

	if gw.Spec.App.Management.SecretName != "" {
		_, err := getSecret(r, ctx, gw, gw.Spec.App.Management.SecretName)
		if err != nil {
			r.Log.Error(err, "Secret not found", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	currSecret := &corev1.Secret{}
	secret := secrets.NewSecret(gw)
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currSecret)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Secret", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, secret, r.Scheme)
		err = r.Create(ctx, secret)
		if err != nil {
			r.Log.Error(err, "Failed creating Secret", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(currSecret.Data, secret.Data) {
		ctrl.SetControllerReference(gw, secret, r.Scheme)
		return r.Update(ctx, secret)
	}
	return nil
}

func reconcileService(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	currService := &corev1.Service{}
	svc := service.NewService(gw)
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currService)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Service", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, svc, r.Scheme)
		err = r.Create(ctx, svc)
		if err != nil {
			r.Log.Error(err, "Failed creating Service", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}
	return nil
}

func reconcileManagementService(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	currService := &corev1.Service{}
	svc := service.NewManagementService(gw)
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name + "-management-service", Namespace: gw.Namespace}, currService)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Management Service", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, svc, r.Scheme)
		err = r.Create(ctx, svc)
		if err != nil {
			r.Log.Error(err, "Failed creating Management Service", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}
	return nil
}

func reconcileIngress(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	if !gw.Spec.App.Ingress.Enabled {
		return nil
	}
	currIngress := &networkingv1.Ingress{}
	ingress := ingress.NewIngress(gw)
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currIngress)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Ingress", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, ingress, r.Scheme)
		err = r.Create(ctx, ingress)
		if err != nil {
			r.Log.Error(err, "Failed creating Ingress", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if !reflect.DeepEqual(currIngress.Spec, ingress.Spec) {
		ctrl.SetControllerReference(gw, ingress, r.Scheme)
		return r.Update(ctx, ingress)
	}
	return nil
}

func reconcileDeployment(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	currDeployment := &appsv1.Deployment{}
	dep := gateway.NewDeployment(gw)
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currDeployment)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Deployment", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, dep, r.Scheme)
		err = r.Create(ctx, dep)
		if err != nil {
			r.Log.Error(err, "Failed creating Deployment", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if gw.Spec.App.Autoscaling.Enabled {
		dep.Spec.Replicas = currDeployment.Spec.Replicas
	}

	//TODO: Refactor and simplify.
	update := false

	var cGatewayContainer corev1.Container
	var nGatewayContainer corev1.Container

	for _, c := range currDeployment.Spec.Template.Spec.Containers {
		if c.Name == "gateway" {
			cGatewayContainer = c
		}
	}

	for _, c := range dep.Spec.Template.Spec.Containers {
		if c.Name == "gateway" {
			nGatewayContainer = c
		}
	}

	currContainerBytes, _ := json.Marshal(cGatewayContainer)
	newContainerBytes, _ := json.Marshal(nGatewayContainer)

	if !reflect.DeepEqual(currContainerBytes, newContainerBytes) {
		update = true
	}

	if *dep.Spec.Replicas != *currDeployment.Spec.Replicas {
		update = true
	}

	if currDeployment.Spec.Template.Spec.ServiceAccountName != dep.Spec.Template.Spec.ServiceAccountName {
		update = true
	}

	if !reflect.DeepEqual(currDeployment.Spec.Template.Spec.ImagePullSecrets, dep.Spec.Template.Spec.ImagePullSecrets) {
		update = true
	}

	if !reflect.DeepEqual(currDeployment.Spec.Template.Spec.InitContainers, dep.Spec.Template.Spec.InitContainers) {
		update = true
	}

	if !reflect.DeepEqual(currDeployment.Spec.Strategy, dep.Spec.Strategy) {
		update = true
	}

	if update {
		ctrl.SetControllerReference(gw, dep, r.Scheme)
		return r.Update(ctx, dep)
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

func getDeploymentStatus(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	dep := appsv1.Deployment{}

	if err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, &dep); err != nil {
		return err
	}

	if dep.Status.AvailableReplicas == dep.Status.Replicas {
		return nil
	}

	return errors.New("Deployment not yet ready")
}

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

func applyGraphmanBundleEphemeral(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {

	notify := false
	notificationMessage := map[string]string{}
	applySuccess := []string{}
	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

	podList, err := getGatewayPods(r, ctx, gw)
	if err != nil {
		return err
	}

	graphmanEncryptionPassphrase, err := getGraphmanEncryptionPassphrase(r, ctx, gw, repoRef)
	if err != nil {
		return err
	}

	for i, pod := range podList.Items {
		update := false
		ready := false

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				ready = containerStatus.Ready
			}
		}

		currentCommit := pod.Labels["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]

		if currentCommit != commit || currentCommit == "" {
			update = true
		}

		if update && ready {
			notify = repoRef.Notification.Enabled
			endpoint := pod.Status.PodIP + ":8443/graphman"
			if len(repoRef.Directories) == 0 {
				repoRef.Directories = []string{"/"}
			}
			for i := range repoRef.Directories {
				gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
				r.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Pod", pod.Name, "Name", gw.Name, "Namespace", gw.Namespace)
				name := gw.Name
				if gw.Spec.App.Management.SecretName != "" {
					name = gw.Spec.App.Management.SecretName
				}
				gwSecret, err := getSecret(r, ctx, gw, name)

				if err != nil {
					return err
				}
				err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
				if err != nil {
					return err
				}
			}

			notificationMessage[pod.Name] = "successfully applied commit"
			applySuccess = append(applySuccess, pod.Name)

			if err := r.Client.Patch(context.Background(), &podList.Items[i],
				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
				return err
			}
		}
	}

	// This is currently limited to WebHooks
	// When additional channels are added this will be refactored.
	if notify {
		notificationMessage["repository"] = repoRef.Name
		notificationMessage["commit"] = commit
		notificationMessage["deployment"] = gw.Name
		notificationMessage["text"] = "Repository: " + repoRef.Name + "\nCommit: " + commit + "\nSuccess: " + strings.Join(applySuccess, ",")
		notificationBytes, _ := json.Marshal(notificationMessage)
		_, err := util.RestCall("POST", repoRef.Notification.Channel.Webhook.Url, repoRef.Notification.Channel.Webhook.InsecureSkipVerify, repoRef.Notification.Channel.Webhook.Headers, "application/json", notificationBytes, repoRef.Notification.Channel.Webhook.Auth.Username, repoRef.Notification.Channel.Webhook.Auth.Password)

		if err != nil {
			r.Log.Info("Failed to send notification", "error", err.Error())
		}
	}
	return nil
}

func applyGraphmanBundleDbBacked(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference, commit string) error {
	notify := false
	notificationMessage := map[string]string{}
	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type, commit)

	gatewayDeployment, err := getGatewayDeployment(r, ctx, gw)

	graphmanEncryptionPassphrase, err := getGraphmanEncryptionPassphrase(r, ctx, gw, repoRef)
	if err != nil {
		return err
	}

	if gatewayDeployment.Status.ReadyReplicas != gatewayDeployment.Status.Replicas {
		return nil
	}

	currentCommit := gatewayDeployment.Labels["security.brcmlabs.com/"+repoRef.Name+"-"+repoRef.Type]
	if currentCommit == commit {
		return nil
	}

	endpoint := gw.Name + "." + gw.Namespace + ".svc.cluster.local:8443/graphman"
	if gw.Spec.App.Management.Service.Enabled {
		endpoint = gw.Name + "-management-service." + gw.Namespace + ".svc.cluster.local:9443/graphman"
	}

	notify = repoRef.Notification.Enabled
	if len(repoRef.Directories) == 0 {
		repoRef.Directories = []string{"/"}
	}
	for i := range repoRef.Directories {
		gitPath := "/tmp/" + repoRef.Name + "/" + repoRef.Directories[i]
		r.Log.Info("Applying Latest Commit", "Repo", repoRef.Name, "Directory", repoRef.Directories[i], "Commit", commit, "Name", gw.Name, "Namespace", gw.Namespace)
		name := gw.Name
		if gw.Spec.App.Management.SecretName != "" {
			name = gw.Spec.App.Management.SecretName
		}
		gwSecret, err := getSecret(r, ctx, gw, name)

		if err != nil {
			return err
		}
		err = util.ApplyToGraphmanTarget(gitPath, string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, graphmanEncryptionPassphrase)
		if err != nil {
			return err
		}
		notificationMessage[gw.Name] = "successfully applied commit"
	}

	if err := r.Client.Patch(context.Background(), &gatewayDeployment,
		client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
		r.Log.Error(err, "Failed to update deployment label", "Namespace", gw.Namespace, "Name", gw.Name)
		return err
	}

	// This is currently limited to WebHooks
	// When additional channels are added this will be refactored.
	if notify {
		notificationMessage["repository"] = repoRef.Name
		notificationMessage["commit"] = commit
		notificationMessage["deployment"] = gw.Name
		notificationMessage["text"] = "Repository: " + repoRef.Name + "\nCommit: " + commit + "\nDeployment updated"
		notificationBytes, _ := json.Marshal(notificationMessage)
		_, err := util.RestCall("POST", repoRef.Notification.Channel.Webhook.Url, repoRef.Notification.Channel.Webhook.InsecureSkipVerify, repoRef.Notification.Channel.Webhook.Headers, "application/json", notificationBytes, repoRef.Notification.Channel.Webhook.Auth.Username, repoRef.Notification.Channel.Webhook.Auth.Password)

		if err != nil {
			r.Log.Info("Failed to send notification", "error", err.Error())
		}
	}
	return nil
}

func getGraphmanEncryptionPassphrase(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, repoRef securityv1.RepositoryReference) (string, error) {
	var graphmanEncryptionPassphrase string
	if repoRef.Encryption.Passphrase != "" && repoRef.Encryption.ExistingSecret == "" {
		graphmanEncryptionPassphrase = repoRef.Encryption.Passphrase
	} else {
		graphmanEncryptionSecret, err := getSecret(r, ctx, gw, repoRef.Encryption.ExistingSecret)
		if err != nil {
			return "", err
		}
		graphmanEncryptionPassphrase = string(graphmanEncryptionSecret.Data[repoRef.Encryption.Key])
	}
	return graphmanEncryptionPassphrase, nil
}

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

// This functionality is experimental and will likely be removed
// in favour of an OtelCollector being created before the Gateway
// resource so that there are no race conditions.
func reconcileOtelCollector(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	if !gw.Spec.App.Monitoring.Otel.Collector.Create {
		return nil
	}
	currOtelCollector := &otelv1alpha1.OpenTelemetryCollector{}
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currOtelCollector)

	newOtelCollector := monitoring.NewOtelCollector(gw)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating OTel Collector", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, newOtelCollector, r.Scheme)
		err = r.Create(ctx, newOtelCollector)
		if err != nil {
			r.Log.Error(err, "Failed creating OTel Collector", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if err != nil {
		r.Log.Info("OtelCollector Error", "Name", gw.Name, "Namespace", gw.Namespace, "error", err.Error())
	}
	return nil
}
func reconcileServiceMonitor(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
	if !gw.Spec.App.Monitoring.ServiceMonitor.Create {
		return nil
	}
	currServiceMonitor := &monitoringv1.ServiceMonitor{}
	err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currServiceMonitor)
	newServiceMonitor := monitoring.NewServiceMonitor(gw)
	if err != nil && k8serrors.IsNotFound(err) {
		r.Log.Info("Creating Service Monitor", "Name", gw.Name, "Namespace", gw.Namespace)
		ctrl.SetControllerReference(gw, newServiceMonitor, r.Scheme)
		err = r.Create(ctx, newServiceMonitor)
		if err != nil {
			r.Log.Error(err, "Failed creating Service Monitor", "Name", gw.Name, "Namespace", gw.Namespace)
			return err
		}
		return nil
	}

	if err != nil {
		r.Log.Info("Service Monitor Error", "Name", gw.Name, "Namespace", gw.Namespace, "error", err.Error())
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.Gateway{}).
		Complete(r)
}
