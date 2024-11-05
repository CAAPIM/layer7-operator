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
const portalTempDirectory = "/tmp/portalapis/"

func Gateway(ctx context.Context, params Params) error {
	graphmanPort := 9443
	checksum := params.Instance.Annotations["app.l7.traceId"]
	tryRequest := true
	isMarkedToBeDeleted := false
	if params.Instance.DeletionTimestamp != nil {
		isMarkedToBeDeleted = true
	}
	// going to need a mechanism to throw an error if sync doesn't fully complete without interrupting other updates.
	updatedStatus := v1alpha1.L7ApiStatus{}

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
				portalMeta.SsgUrlBase64 = base64.StdEncoding.EncodeToString([]byte(portalMeta.SsgUrl))

				portalMeta.ApiEnabled = false
				if params.Instance.Spec.PortalMeta.ApiEnabled {
					portalMeta.ApiEnabled = true
				}

				policyXml := templategen.BuildTemplate(portalMeta)
				graphmanBundleBytes, _, err = api.ConvertPortalPolicyXmlToGraphman(policyXml)
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

				if isMarkedToBeDeleted {
					tryRequest = true
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
									}
								}
							}
						}

						if !statusExists {
							updatedStatus.Gateways = append(updatedStatus.Gateways, v1alpha1.LinkedGatewayStatus{Checksum: checksum, Deployment: tag, Name: pod.Name})
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
							err = util.RemoveL7API(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "/"+params.Instance.Spec.PortalMeta.SsgUrl+"*", params.Instance.Spec.PortalMeta.Name+"-fragment")
							if err != nil {
								params.Log.Info("failed to remove api", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
							}
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
			RemoveTempStorage(ctx, params)
			params.Log.V(2).Info("removed api from temp storage", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		}
	}

	_ = finalizeL7Api(ctx, params)
	return nil
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
