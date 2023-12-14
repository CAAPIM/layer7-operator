package reconcile

import (
	"context"
	"encoding/base64"
	"reflect"
	"strconv"

	v1 "github.com/caapim/layer7-operator/api/v1"
	v1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const apiFinalizer = "security.brcmlabs.com/finalizer"

func Gateway(ctx context.Context, params Params) error {
	graphmanPort := 9443
	checksum := params.Instance.Annotations["checksum/bundle"]
	tryRequest := true
	isMarkedToBeDeleted := params.Instance.DeletionTimestamp != nil
	// going to need a mechanism to throw an error if sync doesn't fully complete without interrupting other updates.
	updatedStatus := v1alpha1.L7ApiStatus{}
	for _, tag := range params.Instance.Spec.DeploymentTags {
		gateway := &v1.Gateway{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: tag, Namespace: params.Instance.Namespace}, gateway)
		if err != nil && k8serrors.IsNotFound(err) {
			params.Log.V(2).Info("gateway not found", "name", tag, "namespace", params.Instance.Namespace)
			continue
		}

		if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
			graphmanPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
		}

		if !gateway.Spec.App.Management.Database.Enabled {
			podList, err := getGatewayPods(ctx, params, gateway)
			if err != nil {
				params.Log.V(2).Info("error retrieving gateway pods", "api", params.Instance.Name, "gateway", gateway.Name, "namespace", params.Instance.Namespace)
				tryRequest = false
				continue
			}

			for _, pod := range podList.Items {

				for _, ds := range params.Instance.Status.Gateways {
					if ds.Name == pod.Name && ds.Deployment == tag {
						if ds.Checksum == checksum && !isMarkedToBeDeleted {
							tryRequest = false
							continue
						}
					}
				}

				for _, us := range updatedStatus.Gateways {
					if us.Name == pod.Name && us.Deployment == tag {
						if us.Checksum == checksum && !isMarkedToBeDeleted {
							tryRequest = false
							continue
						}
					}
				}

				for _, containerStatus := range pod.Status.ContainerStatuses {
					if containerStatus.Name == "gateway" && !containerStatus.Ready {
						tryRequest = false
						continue
					}
				}

				endpoint := pod.Status.PodIP + ":" + strconv.Itoa(graphmanPort) + "/graphman"

				if tryRequest {
					name := gateway.Name
					if gateway.Spec.App.Management.SecretName != "" {
						name = gateway.Spec.App.Management.SecretName
					}
					gwSecret, err := getGatewaySecret(ctx, params, name)

					if err != nil {
						return err
					}

					if !isMarkedToBeDeleted {
						graphmanBundleBytes, err := base64.StdEncoding.DecodeString(params.Instance.Spec.GraphmanBundle)
						if err != nil {
							return err
						}
						params.Log.V(2).Info("applying api", "api", params.Instance.Name, "pod", pod.Name, "namespace", params.Instance.Namespace)
						err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "", graphmanBundleBytes)
						if err != nil {
							return err
						}
						params.Log.Info("applied api", "api", params.Instance.Name, "pod", pod.Name, "namespace", params.Instance.Namespace)
						statusExists := false

						///TODO: Use annotations and update status separately.
						for i, ds := range params.Instance.Status.Gateways {
							if ds.Name == pod.Name && ds.Deployment == tag {
								for _, us := range updatedStatus.Gateways {
									if ds.Name == us.Name && ds.Deployment == us.Deployment {
										statusExists = true
										updatedStatus.Gateways[i].Checksum = checksum
										updatedStatus.Gateways[i].Phase = pod.Status.Phase
									}
								}
							}
						}

						if !statusExists {
							updatedStatus.Gateways = append(updatedStatus.Gateways, v1alpha1.LinkedGatewayStatus{Checksum: checksum, Deployment: tag, Name: pod.Name, Phase: pod.Status.Phase, Ready: true})
						}

						if !reflect.DeepEqual(updatedStatus, params.Instance.Status) && !isMarkedToBeDeleted {
							params.Instance.Status = updatedStatus
							err := params.Client.Status().Update(ctx, params.Instance)
							if err != nil {
								params.Log.V(2).Info("failed to update api status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
							}
						}

					} else {
						if isMarkedToBeDeleted {
							params.Log.V(2).Info("removing api", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
							err = util.RemoveL7API(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "/"+params.Instance.Spec.ServiceUrl+"*", params.Instance.Spec.Name+"-fragment")
							if err != nil {
								params.Log.Info("failed to remove api", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
							}
							params.Log.Info("removed api", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
						}
					}

				}
			}
		}
	}

	if isMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(params.Instance, apiFinalizer) {
			params.Instance.ObjectMeta.Annotations["security.brcmlabs.com/l7api-removed"] = "true"
			_ = params.Client.Update(ctx, params.Instance)
		}
	}

	_ = finalizeL7Api(ctx, params)

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
			if err != nil {
				return gwSecret, err
			}
		}
	}
	return gwSecret, nil
}

func finalizeL7Api(ctx context.Context, params Params) error {
	if params.Instance.ObjectMeta.Annotations["security.brcmlabs.com/l7api-removed"] == "true" {
		params.Log.V(2).Info("removing finalizer", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		controllerutil.RemoveFinalizer(params.Instance, apiFinalizer)
		err := params.Client.Update(ctx, params.Instance)
		if err != nil {
			return err
		}
	}

	return nil
}
