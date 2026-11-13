// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package reconcile

import (
	"context"
	"fmt"

	"github.com/caapim/layer7-operator/pkg/gateway"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func KubernetesGateway(ctx context.Context, params Params) error {
	desired, err := gateway.NewKubernetesGatewayResources(params.Instance)
	if err != nil {
		return err
	}
	if err := reconcileKubernetesGateways(ctx, params, desired.Gateways); err != nil {
		return fmt.Errorf("failed to reconcile Kubernetes Gateways: %w", err)
	}
	if err := reconcileHTTPRoutes(ctx, params, desired.HTTPRoutes); err != nil {
		return fmt.Errorf("failed to reconcile HTTPRoutes: %w", err)
	}
	if err := reconcileTLSRoutes(ctx, params, desired.TLSRoutes); err != nil {
		return fmt.Errorf("failed to reconcile TLSRoutes: %w", err)
	}
	if err := reconcileBackendTLSPolicies(ctx, params, desired.BackendTLSPolicies); err != nil {
		return fmt.Errorf("failed to reconcile BackendTLSPolicies: %w", err)
	}
	return nil
}

func reconcileKubernetesGateways(ctx context.Context, params Params, desiredResources []*gatewayv1.Gateway) error {
	desiredNames := make(map[string]bool, len(desiredResources))
	for _, desired := range desiredResources {
		desiredNames[desired.Name] = true
		if err := controllerutil.SetControllerReference(params.Instance, desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		current := &gatewayv1.Gateway{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: params.Instance.Namespace}, current)
		if k8serrors.IsNotFound(err) {
			if err := params.Client.Create(ctx, desired); err != nil {
				return err
			}
			params.Log.Info("created Gateway", "name", desired.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}
		if !metav1.IsControlledBy(current, params.Instance) {
			return fmt.Errorf("refusing to modify Gateway %s/%s: not owned by this Layer7 Gateway", current.Namespace, current.Name)
		}
		if apiequality.Semantic.DeepEqual(current.Spec, desired.Spec) &&
			apiequality.Semantic.DeepEqual(current.OwnerReferences, desired.OwnerReferences) &&
			metadataKeysMatch(current.ObjectMeta, desired.ObjectMeta) {
			params.Log.V(2).Info("no Gateway updates needed", "name", desired.Name, "namespace", desired.Namespace)
			continue
		}
		updated := current.DeepCopy()
		updated.Spec = desired.Spec
		mergeGatewayAPIMetadata(&updated.ObjectMeta, desired.ObjectMeta)
		patch := client.MergeFrom(current)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return err
		}
		params.Log.Info("updated Gateway", "name", desired.Name, "namespace", desired.Namespace)
	}
	// List the namespace to find previously owned resources after rules are removed
	// or the feature is disabled, including resources whose labels have changed.
	currentResources := &gatewayv1.GatewayList{}
	if err := params.Client.List(ctx, currentResources, client.InNamespace(params.Instance.Namespace)); err != nil {
		// An absent, unused API has no resources to clean up. Other errors must
		// still be returned so failed cleanup is retried.
		if len(desiredResources) == 0 && (meta.IsNoMatchError(err) || k8serrors.IsNotFound(err)) {
			return nil
		}
		return err
	}
	for i := range currentResources.Items {
		current := &currentResources.Items[i]
		if desiredNames[current.Name] || !metav1.IsControlledBy(current, params.Instance) {
			continue
		}
		if err := deleteGatewayAPIResource(ctx, params, current, "Gateway"); err != nil {
			return err
		}
	}
	return nil
}

func reconcileHTTPRoutes(ctx context.Context, params Params, desiredResources []*gatewayv1.HTTPRoute) error {
	desiredNames := make(map[string]bool, len(desiredResources))
	for _, desired := range desiredResources {
		desiredNames[desired.Name] = true
		if err := controllerutil.SetControllerReference(params.Instance, desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		current := &gatewayv1.HTTPRoute{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: params.Instance.Namespace}, current)
		if k8serrors.IsNotFound(err) {
			if err := params.Client.Create(ctx, desired); err != nil {
				return err
			}
			params.Log.Info("created HTTPRoute", "name", desired.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}
		if !metav1.IsControlledBy(current, params.Instance) {
			return fmt.Errorf("refusing to modify HTTPRoute %s/%s: not owned by this Layer7 Gateway", current.Namespace, current.Name)
		}
		if apiequality.Semantic.DeepEqual(current.Spec, desired.Spec) &&
			apiequality.Semantic.DeepEqual(current.OwnerReferences, desired.OwnerReferences) &&
			metadataKeysMatch(current.ObjectMeta, desired.ObjectMeta) {
			params.Log.V(2).Info("no HTTPRoute updates needed", "name", desired.Name, "namespace", desired.Namespace)
			continue
		}
		updated := current.DeepCopy()
		updated.Spec = desired.Spec
		mergeGatewayAPIMetadata(&updated.ObjectMeta, desired.ObjectMeta)
		patch := client.MergeFrom(current)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return err
		}
		params.Log.Info("updated HTTPRoute", "name", desired.Name, "namespace", desired.Namespace)
	}
	// List the namespace to find previously owned resources after rules are removed
	// or the feature is disabled, including resources whose labels have changed.
	currentResources := &gatewayv1.HTTPRouteList{}
	if err := params.Client.List(ctx, currentResources, client.InNamespace(params.Instance.Namespace)); err != nil {
		// An absent, unused API has no resources to clean up. Other errors must
		// still be returned so failed cleanup is retried.
		if len(desiredResources) == 0 && (meta.IsNoMatchError(err) || k8serrors.IsNotFound(err)) {
			return nil
		}
		return err
	}
	for i := range currentResources.Items {
		current := &currentResources.Items[i]
		if desiredNames[current.Name] || !metav1.IsControlledBy(current, params.Instance) {
			continue
		}
		if err := deleteGatewayAPIResource(ctx, params, current, "HTTPRoute"); err != nil {
			return err
		}
	}
	return nil
}

func reconcileTLSRoutes(ctx context.Context, params Params, desiredResources []*gatewayv1.TLSRoute) error {
	desiredNames := make(map[string]bool, len(desiredResources))
	for _, desired := range desiredResources {
		desiredNames[desired.Name] = true
		if err := controllerutil.SetControllerReference(params.Instance, desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		current := &gatewayv1.TLSRoute{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: params.Instance.Namespace}, current)
		if k8serrors.IsNotFound(err) {
			if err := params.Client.Create(ctx, desired); err != nil {
				return err
			}
			params.Log.Info("created TLSRoute", "name", desired.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}
		if !metav1.IsControlledBy(current, params.Instance) {
			return fmt.Errorf("refusing to modify TLSRoute %s/%s: not owned by this Layer7 Gateway", current.Namespace, current.Name)
		}
		if apiequality.Semantic.DeepEqual(current.Spec, desired.Spec) &&
			apiequality.Semantic.DeepEqual(current.OwnerReferences, desired.OwnerReferences) &&
			metadataKeysMatch(current.ObjectMeta, desired.ObjectMeta) {
			params.Log.V(2).Info("no TLSRoute updates needed", "name", desired.Name, "namespace", desired.Namespace)
			continue
		}
		updated := current.DeepCopy()
		updated.Spec = desired.Spec
		mergeGatewayAPIMetadata(&updated.ObjectMeta, desired.ObjectMeta)
		patch := client.MergeFrom(current)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return err
		}
		params.Log.Info("updated TLSRoute", "name", desired.Name, "namespace", desired.Namespace)
	}
	// List the namespace to find previously owned resources after rules are removed
	// or the feature is disabled, including resources whose labels have changed.
	currentResources := &gatewayv1.TLSRouteList{}
	if err := params.Client.List(ctx, currentResources, client.InNamespace(params.Instance.Namespace)); err != nil {
		// An absent, unused API has no resources to clean up. Other errors must
		// still be returned so failed cleanup is retried.
		if len(desiredResources) == 0 && (meta.IsNoMatchError(err) || k8serrors.IsNotFound(err)) {
			return nil
		}
		return err
	}
	for i := range currentResources.Items {
		current := &currentResources.Items[i]
		if desiredNames[current.Name] || !metav1.IsControlledBy(current, params.Instance) {
			continue
		}
		if err := deleteGatewayAPIResource(ctx, params, current, "TLSRoute"); err != nil {
			return err
		}
	}
	return nil
}

func reconcileBackendTLSPolicies(ctx context.Context, params Params, desiredResources []*gatewayv1.BackendTLSPolicy) error {
	desiredNames := make(map[string]bool, len(desiredResources))
	for _, desired := range desiredResources {
		desiredNames[desired.Name] = true
		if err := controllerutil.SetControllerReference(params.Instance, desired, params.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		current := &gatewayv1.BackendTLSPolicy{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: params.Instance.Namespace}, current)
		if k8serrors.IsNotFound(err) {
			if err := params.Client.Create(ctx, desired); err != nil {
				return err
			}
			params.Log.Info("created BackendTLSPolicy", "name", desired.Name, "namespace", params.Instance.Namespace)
			continue
		}
		if err != nil {
			return err
		}
		if !metav1.IsControlledBy(current, params.Instance) {
			return fmt.Errorf("refusing to modify BackendTLSPolicy %s/%s: not owned by this Layer7 Gateway", current.Namespace, current.Name)
		}
		if apiequality.Semantic.DeepEqual(current.Spec, desired.Spec) &&
			apiequality.Semantic.DeepEqual(current.OwnerReferences, desired.OwnerReferences) &&
			metadataKeysMatch(current.ObjectMeta, desired.ObjectMeta) {
			params.Log.V(2).Info("no BackendTLSPolicy updates needed", "name", desired.Name, "namespace", desired.Namespace)
			continue
		}
		updated := current.DeepCopy()
		updated.Spec = desired.Spec
		mergeGatewayAPIMetadata(&updated.ObjectMeta, desired.ObjectMeta)
		patch := client.MergeFrom(current)
		if err := params.Client.Patch(ctx, updated, patch); err != nil {
			return err
		}
		params.Log.Info("updated BackendTLSPolicy", "name", desired.Name, "namespace", desired.Namespace)
	}
	// List the namespace to find previously owned resources after rules are removed
	// or the feature is disabled, including resources whose labels have changed.
	currentResources := &gatewayv1.BackendTLSPolicyList{}
	if err := params.Client.List(ctx, currentResources, client.InNamespace(params.Instance.Namespace)); err != nil {
		// An absent, unused API has no resources to clean up. Other errors must
		// still be returned so failed cleanup is retried.
		if len(desiredResources) == 0 && (meta.IsNoMatchError(err) || k8serrors.IsNotFound(err)) {
			return nil
		}
		return err
	}
	for i := range currentResources.Items {
		current := &currentResources.Items[i]
		if desiredNames[current.Name] || !metav1.IsControlledBy(current, params.Instance) {
			continue
		}
		if err := deleteGatewayAPIResource(ctx, params, current, "BackendTLSPolicy"); err != nil {
			return err
		}
	}
	return nil
}

// mergeGatewayAPIMetadata follows the Ingress metadata merge: retain extra keys
// written by other controllers and update the keys supplied by the Gateway CR.
func mergeGatewayAPIMetadata(current *metav1.ObjectMeta, desired metav1.ObjectMeta) {
	current.OwnerReferences = desired.OwnerReferences
	if current.Labels == nil {
		current.Labels = make(map[string]string)
	}
	for k, v := range desired.Labels {
		current.Labels[k] = v
	}
	if current.Annotations == nil {
		current.Annotations = make(map[string]string)
	}
	for k, v := range desired.Annotations {
		current.Annotations[k] = v
	}
}

func deleteGatewayAPIResource(ctx context.Context, params Params, current client.Object, kind string) error {
	// Guard against deleting a resource replaced or reparented since it was listed.
	uid, version := current.GetUID(), current.GetResourceVersion()
	if err := params.Client.Delete(ctx, current, &client.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid, ResourceVersion: &version}}); err != nil && !k8serrors.IsNotFound(err) {
		return err
	}
	params.Log.Info("removed Gateway API resource", "kind", kind, "name", current.GetName(), "namespace", current.GetNamespace())
	return nil
}
