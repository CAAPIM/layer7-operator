// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package reconcile

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// gatewayAPIManagedMetadataAnnotation records the label and annotation keys the
// operator last applied to a generated Gateway API resource, so keys removed from
// the Gateway CR can be deleted without touching keys written by other controllers.
const gatewayAPIManagedMetadataAnnotation = "security.brcmlabs.com/managed-metadata"

// gatewayAPIApplyFailedReason is the Event reason used when a Gateway API resource
// can't be applied, so users see the API server's message on the Gateway CR.
const gatewayAPIApplyFailedReason = "GatewayAPIApplyFailed"

// gatewayAPIKind adapts one Gateway API resource type to reconcileGatewayAPIResources.
type gatewayAPIKind[T client.Object] struct {
	kind      string
	newObj    func() T
	list      func(ctx context.Context, c client.Client, namespace string) ([]T, error)
	specEqual func(current, desired T) bool
	copySpec  func(dst, src T)
}

func itemPointers[E any](items []E) []*E {
	ptrs := make([]*E, len(items))
	for i := range items {
		ptrs[i] = &items[i]
	}
	return ptrs
}

var kubernetesGatewayKind = gatewayAPIKind[*gatewayv1.Gateway]{
	kind:   "Gateway",
	newObj: func() *gatewayv1.Gateway { return &gatewayv1.Gateway{} },
	list: func(ctx context.Context, c client.Client, namespace string) ([]*gatewayv1.Gateway, error) {
		l := &gatewayv1.GatewayList{}
		if err := c.List(ctx, l, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		return itemPointers(l.Items), nil
	},
	specEqual: func(current, desired *gatewayv1.Gateway) bool {
		return apiequality.Semantic.DeepEqual(current.Spec, desired.Spec)
	},
	copySpec: func(dst, src *gatewayv1.Gateway) { dst.Spec = src.Spec },
}

var httpRouteKind = gatewayAPIKind[*gatewayv1.HTTPRoute]{
	kind:   "HTTPRoute",
	newObj: func() *gatewayv1.HTTPRoute { return &gatewayv1.HTTPRoute{} },
	list: func(ctx context.Context, c client.Client, namespace string) ([]*gatewayv1.HTTPRoute, error) {
		l := &gatewayv1.HTTPRouteList{}
		if err := c.List(ctx, l, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		return itemPointers(l.Items), nil
	},
	specEqual: func(current, desired *gatewayv1.HTTPRoute) bool {
		return apiequality.Semantic.DeepEqual(current.Spec, desired.Spec)
	},
	copySpec: func(dst, src *gatewayv1.HTTPRoute) { dst.Spec = src.Spec },
}

var tlsRouteKind = gatewayAPIKind[*gatewayv1.TLSRoute]{
	kind:   "TLSRoute",
	newObj: func() *gatewayv1.TLSRoute { return &gatewayv1.TLSRoute{} },
	list: func(ctx context.Context, c client.Client, namespace string) ([]*gatewayv1.TLSRoute, error) {
		l := &gatewayv1.TLSRouteList{}
		if err := c.List(ctx, l, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		return itemPointers(l.Items), nil
	},
	specEqual: func(current, desired *gatewayv1.TLSRoute) bool {
		return apiequality.Semantic.DeepEqual(current.Spec, desired.Spec)
	},
	copySpec: func(dst, src *gatewayv1.TLSRoute) { dst.Spec = src.Spec },
}

var backendTLSPolicyKind = gatewayAPIKind[*gatewayv1.BackendTLSPolicy]{
	kind:   "BackendTLSPolicy",
	newObj: func() *gatewayv1.BackendTLSPolicy { return &gatewayv1.BackendTLSPolicy{} },
	list: func(ctx context.Context, c client.Client, namespace string) ([]*gatewayv1.BackendTLSPolicy, error) {
		l := &gatewayv1.BackendTLSPolicyList{}
		if err := c.List(ctx, l, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		return itemPointers(l.Items), nil
	},
	specEqual: func(current, desired *gatewayv1.BackendTLSPolicy) bool {
		return apiequality.Semantic.DeepEqual(current.Spec, desired.Spec)
	},
	copySpec: func(dst, src *gatewayv1.BackendTLSPolicy) { dst.Spec = src.Spec },
}

func KubernetesGateway(ctx context.Context, params Params) error {
	desired, err := gateway.NewKubernetesGatewayResources(params.Instance)
	if err != nil {
		return err
	}
	if err := reconcileGatewayAPIResources(ctx, params, kubernetesGatewayKind, desired.Gateways); err != nil {
		return fmt.Errorf("failed to reconcile Kubernetes Gateways: %w", err)
	}
	if err := reconcileGatewayAPIResources(ctx, params, httpRouteKind, desired.HTTPRoutes); err != nil {
		return fmt.Errorf("failed to reconcile HTTPRoutes: %w", err)
	}
	if err := reconcileGatewayAPIResources(ctx, params, tlsRouteKind, desired.TLSRoutes); err != nil {
		return fmt.Errorf("failed to reconcile TLSRoutes: %w", err)
	}
	if err := reconcileGatewayAPIResources(ctx, params, backendTLSPolicyKind, desired.BackendTLSPolicies); err != nil {
		return fmt.Errorf("failed to reconcile BackendTLSPolicies: %w", err)
	}
	return nil
}

// reconcileGatewayAPIResources creates or updates every desired resource of one kind,
// then deletes resources of that kind previously owned by this Gateway that are no
// longer desired.
func reconcileGatewayAPIResources[T client.Object](ctx context.Context, params Params, k gatewayAPIKind[T], desiredResources []T) error {
	desiredNames := make(map[string]bool, len(desiredResources))
	for _, desired := range desiredResources {
		desiredNames[desired.GetName()] = true
		if err := applyGatewayAPIResource(ctx, params, k, desired); err != nil {
			recordGatewayAPIFailure(params, k.kind, desired.GetName(), err)
			return err
		}
	}
	// List the namespace to find previously owned resources after rules are removed
	// or the feature is disabled, including resources whose labels have changed.
	currentResources, err := k.list(ctx, params.Client, params.Instance.Namespace)
	if err != nil {
		// An absent, unused API has no resources to clean up. Other errors must
		// still be returned so failed cleanup is retried.
		if len(desiredResources) == 0 && (meta.IsNoMatchError(err) || k8serrors.IsNotFound(err)) {
			return nil
		}
		return err
	}
	for _, current := range currentResources {
		if desiredNames[current.GetName()] || !metav1.IsControlledBy(current, params.Instance) {
			continue
		}
		if err := deleteGatewayAPIResource(ctx, params, current, k.kind); err != nil {
			recordGatewayAPIFailure(params, k.kind, current.GetName(), err)
			return err
		}
	}
	return nil
}

func applyGatewayAPIResource[T client.Object](ctx context.Context, params Params, k gatewayAPIKind[T], desired T) error {
	if err := controllerutil.SetControllerReference(params.Instance, desired, params.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference: %w", err)
	}
	setManagedMetadataAnnotation(desired)
	current := k.newObj()
	err := params.Client.Get(ctx, types.NamespacedName{Name: desired.GetName(), Namespace: params.Instance.Namespace}, current)
	if k8serrors.IsNotFound(err) {
		if err := params.Client.Create(ctx, desired); err != nil {
			return err
		}
		params.Log.Info("created "+k.kind, "name", desired.GetName(), "namespace", params.Instance.Namespace)
		return nil
	}
	if err != nil {
		return err
	}
	if !metav1.IsControlledBy(current, params.Instance) {
		return fmt.Errorf("refusing to modify %s %s/%s: not owned by this Layer7 Gateway", k.kind, current.GetNamespace(), current.GetName())
	}
	if k.specEqual(current, desired) &&
		apiequality.Semantic.DeepEqual(current.GetOwnerReferences(), desired.GetOwnerReferences()) &&
		gatewayAPIMetadataMatches(current, desired) {
		params.Log.V(2).Info("no "+k.kind+" updates needed", "name", desired.GetName(), "namespace", desired.GetNamespace())
		return nil
	}
	updated := current.DeepCopyObject().(T)
	k.copySpec(updated, desired)
	mergeGatewayAPIMetadata(updated, desired)
	if err := params.Client.Patch(ctx, updated, client.MergeFrom(current)); err != nil {
		return err
	}
	params.Log.Info("updated "+k.kind, "name", desired.GetName(), "namespace", desired.GetNamespace())
	return nil
}

// recordGatewayAPIFailure surfaces the API server's error on the Gateway CR, so a
// rejected route (for example an invalid schemaless match) is visible with
// kubectl describe rather than only in operator logs.
func recordGatewayAPIFailure(params Params, kind, name string, err error) {
	if params.Recorder == nil {
		return
	}
	params.Recorder.Eventf(params.Instance, corev1.EventTypeWarning, gatewayAPIApplyFailedReason, "%s %s: %v", kind, name, err)
}

type managedMetadata struct {
	Labels      []string `json:"labels,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
}

// setManagedMetadataAnnotation records the desired label and annotation keys on
// the desired object itself, so it is written on create and compared like any
// other desired annotation on update.
func setManagedMetadataAnnotation(desired client.Object) {
	annotations := desired.GetAnnotations()
	keys := managedMetadata{
		Labels:      slices.Sorted(maps.Keys(desired.GetLabels())),
		Annotations: slices.Sorted(maps.Keys(annotations)),
	}
	keys.Annotations = slices.DeleteFunc(keys.Annotations, func(k string) bool { return k == gatewayAPIManagedMetadataAnnotation })
	raw, _ := json.Marshal(keys)
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[gatewayAPIManagedMetadataAnnotation] = string(raw)
	desired.SetAnnotations(annotations)
}

// previouslyManagedMetadata returns the keys recorded on current. A missing or
// unreadable record is treated as empty, so no keys are deleted on the first
// reconcile of a resource created before tracking existed.
func previouslyManagedMetadata(current client.Object) managedMetadata {
	keys := managedMetadata{}
	if raw, ok := current.GetAnnotations()[gatewayAPIManagedMetadataAnnotation]; ok {
		if err := json.Unmarshal([]byte(raw), &keys); err != nil {
			return managedMetadata{}
		}
	}
	return keys
}

// staleKeys returns keys the operator managed previously that are still present
// on current but no longer desired.
func staleKeys(previous []string, current, desired map[string]string) []string {
	stale := []string{}
	for _, k := range previous {
		if _, wanted := desired[k]; wanted {
			continue
		}
		if _, present := current[k]; present {
			stale = append(stale, k)
		}
	}
	return stale
}

// gatewayAPIMetadataMatches reports whether current already carries every desired
// label and annotation (including the managed-metadata record) and none of the
// keys the operator managed previously but no longer wants.
func gatewayAPIMetadataMatches(current, desired client.Object) bool {
	currentMeta := metav1.ObjectMeta{Labels: current.GetLabels(), Annotations: current.GetAnnotations()}
	desiredMeta := metav1.ObjectMeta{Labels: desired.GetLabels(), Annotations: desired.GetAnnotations()}
	if !metadataKeysMatch(currentMeta, desiredMeta) {
		return false
	}
	previous := previouslyManagedMetadata(current)
	return len(staleKeys(previous.Labels, current.GetLabels(), desired.GetLabels())) == 0 &&
		len(staleKeys(previous.Annotations, current.GetAnnotations(), desired.GetAnnotations())) == 0
}

// mergeGatewayAPIMetadata updates the keys supplied by the Gateway CR, deletes keys
// the operator managed previously but no longer wants, and keeps keys written by
// other controllers.
func mergeGatewayAPIMetadata(current, desired client.Object) {
	current.SetOwnerReferences(desired.GetOwnerReferences())
	previous := previouslyManagedMetadata(current)

	labels := maps.Clone(current.GetLabels())
	if labels == nil {
		labels = map[string]string{}
	}
	for _, k := range staleKeys(previous.Labels, labels, desired.GetLabels()) {
		delete(labels, k)
	}
	maps.Copy(labels, desired.GetLabels())
	current.SetLabels(labels)

	annotations := maps.Clone(current.GetAnnotations())
	if annotations == nil {
		annotations = map[string]string{}
	}
	for _, k := range staleKeys(previous.Annotations, annotations, desired.GetAnnotations()) {
		delete(annotations, k)
	}
	maps.Copy(annotations, desired.GetAnnotations())
	current.SetAnnotations(annotations)
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
