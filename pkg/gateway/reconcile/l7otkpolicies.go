package reconcile

import (
	"context"
	"strconv"
	"strings"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncOtkPolicies(ctx context.Context, params Params) {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-policies")
		return
	}

	if !gateway.Spec.App.Otk.Enabled {
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-policies")
		return
	}

	params.Instance = gateway

	err = applyOtkPolicies(ctx, params, gateway)
	if err != nil {
		params.Log.Error(err, "failed to reconcile otk policies", "name", gateway.Name, "namespace", gateway.Namespace)
	}

}

func applyOtkPolicies(ctx context.Context, params Params, gateway *securityv1.Gateway) error {
	internalGatewayPort := 9443
	defaultOtkPort := 8443
	if gateway.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		internalGatewayPort = gateway.Spec.App.Management.Graphman.DynamicSyncPort
	}

	if gateway.Spec.App.Otk.OTKPort != 0 {
		defaultOtkPort = gateway.Spec.App.Otk.OTKPort
	}

	var gatewayHost string
	switch gateway.Spec.App.Otk.Type {
	case securityv1.OtkTypeDMZ:
		// TODO: open this to internal gateways that are fully external or in a different namespace
		// This routes via 9443 or the management port by default
		if gateway.Spec.App.Otk.InternalGatewayPort != 0 {
			internalGatewayPort = gateway.Spec.App.Otk.InternalGatewayPort
		}
		gatewayHost = "https://" + gateway.Spec.App.Otk.InternalOtkGatewayReference + ":" + strconv.Itoa(internalGatewayPort)
	case securityv1.OtkTypeInternal:
		gatewayHost = "https://" + gateway.Name + ":" + strconv.Itoa(defaultOtkPort)
	}

	bundle, sha1Sum, err := util.BuildOtkOverrideBundle(strings.ToUpper(string(gateway.Spec.App.Otk.Type)), gatewayHost, defaultOtkPort)
	if err != nil {
		return err
	}

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	annotation := "security.brcmlabs.com/" + gateway.Name + "-" + string(gateway.Spec.App.Otk.Type) + "-policies"

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "otk policies", *podList, gateway, gwSecret, "", annotation, sha1Sum, false, "otk policies", bundle)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, params, "otk policies", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, bundle)
		if err != nil {
			return err
		}
	}

	return nil
}
