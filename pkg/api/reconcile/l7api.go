// Copyright (c) 2025 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package reconcile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	v1 "github.com/caapim/layer7-operator/api/v1"
	v1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/caapim/layer7-operator/pkg/api"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const apiFinalizer = "security.brcmlabs.com/finalizer"
const L7API_REMOVED_ANNOTATION = "security.brcmlabs.com/l7api-removed"
const portalTempDirectory = "/tmp/portalapis/"
const GATEWAY_STATUS_CONDITIONS_MAX_LEN = 1
const (
	DEPLOY   = "DEPLOY"
	UNDEPLOY = "UNDEPLOY"
	SUCCESS  = "SUCCESS"
	FAILURE  = "FAILURE"
)

func Gateway(ctx context.Context, params Params) error {

	isMarkedToBeDeleted := false
	if params.Instance.DeletionTimestamp != nil {
		isMarkedToBeDeleted = true
	}

	if !params.Instance.Status.Ready {
		return fmt.Errorf("api %s not ready in namespace %s", params.Instance.Name, params.Instance.Namespace)
	}

	for _, tag := range params.Instance.Spec.DeploymentTags {
		gateway := &v1.Gateway{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: tag, Namespace: params.Instance.Namespace}, gateway)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("gateway not found", "name", tag, "namespace", params.Instance.Namespace)
			continue
		}

		// going to need a mechanism to throw an error if sync doesn't fully complete without interrupting other updates.
		updatedStatus, err := cloneL7ApiStatus(&params.Instance.Status)
		if err != nil {
			return err
		}

		if !isMarkedToBeDeleted {
			err = deployL7ApiToGateway(ctx, params, gateway, tag, updatedStatus)
		} else {
			if params.Instance.ObjectMeta.Annotations[L7API_REMOVED_ANNOTATION] == "true" {
				params.Log.V(2).Info("skip un-deployment since it has been done", "api", params.Instance.Name)
				return nil
			}
			err = undeployL7ApiToGateway(ctx, params, gateway, tag, updatedStatus)
		}
		if err != nil {
			return err
		}

		// persist the deployment status once per deployment tag
		if !reflect.DeepEqual(*updatedStatus, params.Instance.Status) {
			params.Instance.Status = *updatedStatus
			err := params.Client.Status().Update(ctx, params.Instance)
			if err != nil {
				params.Log.V(2).Info("failed to update api status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			}
		}
	}

	if isMarkedToBeDeleted {
		params.Instance.ObjectMeta.Annotations[L7API_REMOVED_ANNOTATION] = "true"
		_ = params.Client.Update(ctx, params.Instance)
		RemoveTempStorage(ctx, params)
		params.Log.V(2).Info("removed api from temp storage", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}
	// If it is portal published, g2c agent will remove the finalizer after call back to portal to update the api deployment status
	if !params.Instance.Spec.PortalPublished {
		err := finalizeL7Api(ctx, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func cloneL7ApiStatus(orig *v1alpha1.L7ApiStatus) (*v1alpha1.L7ApiStatus, error) {
	origJSON, err := json.Marshal(orig)
	if err != nil {
		return nil, err
	}

	clone := v1alpha1.L7ApiStatus{}
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
}

// deploy the L7Api in params to the gateway pods except the pods in which the api has been deployed or the pods aren't ready yet.
func deployL7ApiToGateway(ctx context.Context, params Params, gateway *v1.Gateway, tag string, updatedStatus *v1alpha1.L7ApiStatus) error {
	graphmanPort := 9443
	tryRequest := true
	checksum := params.Instance.Annotations["app.l7.traceId"]
	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params, gateway)
		if err != nil {
			params.Log.V(2).Info("error retrieving gateway pods", "api", params.Instance.Name, "gateway", gateway.Name, "namespace", params.Instance.Namespace)
			tryRequest = false
			return nil
		}

		var graphmanBundleBytes []byte

		if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
			portalMeta := templategen.PortalAPI{}
			portalMetaBytes, err := json.Marshal(params.Instance.Spec.PortalMeta)
			if err != nil {
				return err
			}
			err = json.Unmarshal(portalMetaBytes, &portalMeta)
			if err != nil {
				return err
			}

			portalMeta.LocationUrl = base64.StdEncoding.EncodeToString([]byte(portalMeta.LocationUrl))
			//trim wildcard char for usage in policy context var: serviceUrl
			serviceUrl, _ := strings.CutSuffix(portalMeta.SsgUrl, "*")
			portalMeta.SsgUrlBase64 = base64.StdEncoding.EncodeToString([]byte(serviceUrl))

			portalMeta.ApiEnabled = false
			if params.Instance.Spec.PortalMeta.ApiEnabled {
				portalMeta.ApiEnabled = true
			}

			policyXml := templategen.BuildTemplate(portalMeta)
			graphmanBundleBytes, _, err = api.ConvertPortalPolicyXmlToGraphman(policyXml, portalMeta.SecurePasswords, portalMeta.SecurePasswordIdsForUndeployment)
			if err != nil {
				return err
			}

		} else {
			graphmanBundleBytes, err = base64.StdEncoding.DecodeString(params.Instance.Spec.GraphmanBundle)
			if err != nil {
				return err
			}
		}

		for _, pod := range podList.Items {
			// if the checksum is in the pod condition already, it means the deployment has been done on the gateway pod.
			for _, us := range updatedStatus.Gateways {
				if us.Name == pod.Name && us.Deployment == tag {
					for _, condition := range us.Conditions {
						if condition.Checksum == checksum && condition.Action == DEPLOY {
							tryRequest = false
							continue
						}
					}
				}
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == "gateway" && !containerStatus.Ready {
					tryRequest = false
					continue
				}
			}

			if tryRequest {
				endpoint := "127.0.0.1" + ":" + strconv.Itoa(graphmanPort) + "/graphman"
				var errorMessage string
				status := SUCCESS
				name := gateway.Name
				if gateway.Spec.App.Management.SecretName != "" {
					name = gateway.Spec.App.Management.SecretName
				}
				gwSecret, err := getGatewaySecret(ctx, params, name)

				if err != nil {
					return err
				}

				params.Log.V(2).Info("applying api", "api", params.Instance.Name, "pod", pod.Name, "namespace", params.Instance.Namespace)
				err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "", graphmanBundleBytes)
				if err != nil {
					status = FAILURE
					errorMessage = err.Error()
					params.Log.Error(err, "applied api", "api", params.Instance.Name, "pod", pod.Name, "namespace", params.Instance.Namespace)
				} else {
					params.Log.Info("applied api", "api", params.Instance.Name, "pod", pod.Name, "namespace", params.Instance.Namespace)
				}
				updateL7ApiDeploymentStatusOnPod(tag, pod.Name, checksum, DEPLOY, status, errorMessage, updatedStatus)
			}
		}
	}
	return nil
}

// un-deploy the L7Api in params from the gateway pods except the pods in which the api has been un-deployed or the pods aren't ready yet.
func undeployL7ApiToGateway(ctx context.Context, params Params, gateway *v1.Gateway, tag string, updatedStatus *v1alpha1.L7ApiStatus) error {
	graphmanPort := 9443
	tryRequest := true
	checksum := params.Instance.Annotations["app.l7.traceId"]
	secretNames := []string{}
	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	for _, securePassword := range params.Instance.Spec.PortalMeta.SecurePasswords {
		secretNames = append(secretNames, securePassword.Name)
	}

	for _, securePasswordIdsForUndeployment := range params.Instance.Spec.PortalMeta.SecurePasswordIdsForUndeployment {
		secretToDelete := "l7_secure_" + securePasswordIdsForUndeployment
		secretNames = append(secretNames, secretToDelete)
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params, gateway)
		if err != nil {
			params.Log.V(2).Info("error retrieving gateway pods", "api", params.Instance.Name, "gateway", gateway.Name, "namespace", params.Instance.Namespace)
			return err
		}

		for _, pod := range podList.Items {
			// if the checksum is in the pod condition already, it means the un-deployment has been done on the gateway pod.
			for _, us := range updatedStatus.Gateways {
				if us.Name == pod.Name && us.Deployment == tag {
					for _, condition := range us.Conditions {
						if condition.Checksum == checksum && condition.Action == UNDEPLOY {
							tryRequest = false
							continue
						}
					}
				}
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == "gateway" && !containerStatus.Ready {
					tryRequest = false
					continue
				}
			}

			if tryRequest {
				endpoint := "127.0.0.1" + ":" + strconv.Itoa(graphmanPort) + "/graphman"
				status := SUCCESS
				name := gateway.Name
				if gateway.Spec.App.Management.SecretName != "" {
					name = gateway.Spec.App.Management.SecretName
				}
				gwSecret, err := getGatewaySecret(ctx, params, name)

				if err != nil {
					return err
				}

				params.Log.V(2).Info("removing api", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
				var errorMessage string
				err = util.RemoveL7API(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "/"+params.Instance.Spec.PortalMeta.SsgUrl, params.Instance.Spec.PortalMeta.Name+"-fragment", secretNames)
				if err != nil {
					status = FAILURE
					errorMessage = err.Error()
					params.Log.Error(err, "failed to remove api", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
				} else {
					params.Log.Info("removed api", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
				}
				updateL7ApiDeploymentStatusOnPod(tag, pod.Name, checksum, UNDEPLOY, status, errorMessage, updatedStatus)
			}
		}
	}
	return nil
}

// update updatedStatus instead of persisting status into k8s. the persisting status into k8s will happen after all pods are deployed.
// g2cagent expects only one k8s update event per deployment request from portal.
func updateL7ApiDeploymentStatusOnPod(tag string, podName string, checksum string, action string, status string, errorMessage string, updatedStatus *v1alpha1.L7ApiStatus) {
	condition := v1alpha1.GatewayPodDeploymentCondition{
		Action:     action,
		ActionTime: time.Now().UTC().Format(time.RFC3339),
		Checksum:   checksum,
		Reason:     errorMessage,
		Status:     status,
	}
	statusExists := false
	for i, ds := range updatedStatus.Gateways {
		if ds.Name == podName && ds.Deployment == tag {
			statusExists = true
			updatedStatus.Gateways[i].Conditions = append(updatedStatus.Gateways[i].Conditions, condition)
			// truncate the conditions if the size is too big.
			if len(updatedStatus.Gateways[i].Conditions) > GATEWAY_STATUS_CONDITIONS_MAX_LEN {
				updatedStatus.Gateways[i].Conditions = updatedStatus.Gateways[i].Conditions[len(updatedStatus.Gateways[i].Conditions)-GATEWAY_STATUS_CONDITIONS_MAX_LEN:]
			}
		}
	}

	if !statusExists {
		updatedStatus.Gateways = append(updatedStatus.Gateways, v1alpha1.LinkedGatewayStatus{
			Deployment: tag,
			Name:       podName,
			Conditions: []v1alpha1.GatewayPodDeploymentCondition{condition},
		})
	}
}

// TempStorage writes API Metadata to /tmp/portalapis/<l7PortalName>/apiname.json
// This does not track the deployment tag which will be resolved in a future update
// l7api and l7portal should have the same deployment tags, l7portal deployment tags are
// primarily used for bootstrapping portal apis to target container gateway deployments.
// this mechanism will be updated in the future.
func WriteTempStorage(ctx context.Context, params Params) error {
	apiPath := portalTempDirectory + params.Instance.Spec.L7Portal + "/"
	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
		portalMeta := templategen.PortalAPI{}
		portalMetaBytes, err := json.Marshal(params.Instance.Spec.PortalMeta)
		if err != nil {
			return err
		}
		err = json.Unmarshal(portalMetaBytes, &portalMeta)
		if err != nil {
			return err
		}

		if _, err := os.Stat(apiPath); os.IsNotExist(err) {
			err = os.MkdirAll(apiPath, 0760)
			if err != nil {
				return err
			}
		}

		err = os.WriteFile(apiPath+strings.ReplaceAll(params.Instance.Name, " ", "-")+".json", portalMetaBytes, 0660)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveTempStorage(ctx context.Context, params Params) error {
	apiPath := portalTempDirectory + params.Instance.Spec.L7Portal + "/"
	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
		err := os.Remove(apiPath + strings.ReplaceAll(params.Instance.Name, " ", "-") + ".json")
		if err != nil {
			return err
		}

	}
	return nil
}

// GetGatewayPods returns the pods in a Gateway Deployment
func getGatewayPods(ctx context.Context, params Params, gateway *v1.Gateway) (*corev1.PodList, error) {
	podList := &corev1.PodList{}

	listOpts := []client.ListOption{
		client.InNamespace(gateway.Namespace),
		client.MatchingLabels(util.DefaultLabels(gateway.Name, map[string]string{})),
	}
	if err := params.Client.List(ctx, podList, listOpts...); err != nil {
		return podList, err
	}
	return podList, nil
}

func getGatewaySecret(ctx context.Context, params Params, name string) (*corev1.Secret, error) {
	gwSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, gwSecret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return gwSecret, err
		}
	}
	return gwSecret, nil
}

func finalizeL7Api(ctx context.Context, params Params) error {
	if params.Instance.ObjectMeta.Annotations[L7API_REMOVED_ANNOTATION] == "true" {
		params.Log.V(2).Info("removing finalizer", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		controllerutil.RemoveFinalizer(params.Instance, apiFinalizer)
		err := params.Client.Update(ctx, params.Instance)
		if err != nil {
			params.Log.V(2).Info("fail to remove finalizer", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
	}
	return nil
}
