// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package util

import (
	"context"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GatewayDeploymentPods returns the pods that carry the Gateway's default labels
// and were created by the Gateway's own Deployment. The operator sends admin
// credentials and bundles to these pods, so labels alone are not trusted.
//
// A pod is trusted only when all of the following hold:
//   - the Gateway's Deployment exists and is controlled by the Gateway CR
//   - the pod's controller is a ReplicaSet named <deployment>-<pod-template-hash>
//   - that ReplicaSet exists, its UID matches the pod's ownerReference UID, and it
//     is controlled by the Deployment (by UID)
//
// ownerReferences are written by clients, so this raises the bar rather than
// forming a complete trust boundary: a principal that can create pods or
// ReplicaSets in the namespace and read the real UIDs can still forge a chain that
// passes. Restrict who can create pods and ReplicaSets in Gateway namespaces.
//
// On error no pods are returned, so callers never act on unverified pods.
func GatewayDeploymentPods(ctx context.Context, c client.Client, log logr.Logger, gw *securityv1.Gateway) (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(gw.Namespace),
		client.MatchingLabels(DefaultLabels(gw.Name, map[string]string{})),
	}
	if err := c.List(ctx, podList, listOpts...); err != nil {
		podList.Items = nil
		return podList, err
	}
	if len(podList.Items) == 0 {
		return podList, nil
	}

	deployment := &appsv1.Deployment{}
	if err := c.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, deployment); err != nil {
		podList.Items = nil
		if k8serrors.IsNotFound(err) {
			log.V(2).Info("gateway deployment not found, ignoring labelled pods", "name", gw.Name, "namespace", gw.Namespace)
			return podList, nil
		}
		return podList, err
	}
	if !metav1.IsControlledBy(deployment, gw) {
		log.Info("gateway deployment is not controlled by this gateway, ignoring labelled pods", "name", gw.Name, "namespace", gw.Namespace)
		podList.Items = nil
		return podList, nil
	}

	// Look each ReplicaSet up once per call. A nil entry records NotFound.
	replicaSets := map[string]*metav1.PartialObjectMetadata{}
	trustedPods := []corev1.Pod{}
	for i := range podList.Items {
		pod := &podList.Items[i]
		// A Deployment's ReplicaSets are named <deployment>-<pod-template-hash>.
		owner := metav1.GetControllerOf(pod)
		hash := pod.Labels[appsv1.DefaultDeploymentUniqueLabelKey]
		if owner == nil || owner.Kind != "ReplicaSet" || owner.APIVersion != appsv1.SchemeGroupVersion.String() ||
			hash == "" || owner.Name != deployment.Name+"-"+hash {
			log.V(2).Info("ignoring pod with gateway labels that is not owned by the gateway deployment", "pod", pod.Name, "name", gw.Name, "namespace", gw.Namespace)
			continue
		}
		rs, seen := replicaSets[owner.Name]
		if !seen {
			// Metadata only: the UID and ownerReferences are all that is needed, and a
			// metadata informer keeps the operator's cache small.
			rs = &metav1.PartialObjectMetadata{}
			rs.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("ReplicaSet"))
			if err := c.Get(ctx, types.NamespacedName{Name: owner.Name, Namespace: gw.Namespace}, rs); err != nil {
				if !k8serrors.IsNotFound(err) {
					podList.Items = nil
					return podList, err
				}
				rs = nil
			}
			replicaSets[owner.Name] = rs
		}
		if rs == nil || rs.UID != owner.UID || !metav1.IsControlledBy(rs, deployment) {
			log.V(2).Info("ignoring pod whose ReplicaSet owner is not a live ReplicaSet of the gateway deployment", "pod", pod.Name, "replicaSet", owner.Name, "name", gw.Name, "namespace", gw.Namespace)
			continue
		}
		trustedPods = append(trustedPods, *pod)
	}
	podList.Items = trustedPods
	return podList, nil
}
