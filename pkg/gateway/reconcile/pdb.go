package reconcile

import (
	"context"
	"fmt"
	"reflect"

	"github.com/caapim/layer7-operator/pkg/gateway"
	policyv1 "k8s.io/api/policy/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func PodDisruptionBudget(ctx context.Context, params Params) error {
	if !params.Instance.Spec.App.PodDisruptionBudget.Enabled {
		return nil
	}
	desiredPdb := gateway.NewPDB(params.Instance)
	currentPdb := &policyv1.PodDisruptionBudget{}

	if err := controllerutil.SetControllerReference(params.Instance, desiredPdb, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}

	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, currentPdb)

	if err != nil && k8serrors.IsNotFound(err) {
		if err = params.Client.Create(ctx, desiredPdb); err != nil {
			return fmt.Errorf("failed creating pod disruption budget: %w", err)
		}
		params.Log.Info("created pod disruption budget", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return nil
	}

	if err != nil {
		return err
	}

	if reflect.DeepEqual(currentPdb.Spec, desiredPdb.Spec) {
		params.Log.V(2).Info("no pod disruption budget updates needed", "name", desiredPdb.Name, "namespace", desiredPdb.Namespace)
		return nil
	}

	updated := currentPdb.DeepCopy()
	updated.Spec = desiredPdb.Spec

	updated.ObjectMeta.OwnerReferences = desiredPdb.ObjectMeta.OwnerReferences

	for k, v := range desiredPdb.ObjectMeta.Annotations {
		updated.ObjectMeta.Annotations[k] = v
	}
	for k, v := range desiredPdb.ObjectMeta.Labels {
		updated.ObjectMeta.Labels[k] = v
	}

	patch := client.MergeFrom(currentPdb)

	if err := params.Client.Patch(ctx, updated, patch); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	params.Log.Info("updated pod disruption budget", "name", desiredPdb.Name, "namespace", desiredPdb.Namespace)

	return nil

}
