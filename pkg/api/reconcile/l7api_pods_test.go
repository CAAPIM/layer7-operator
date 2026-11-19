// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
package reconcile

import (
	"context"
	"testing"

	v1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// The L7Api controller sends bundles and gateway credentials to the pods it
// finds, so it must apply the same ownership checks as the Gateway controller.
func TestL7ApiGetGatewayPodsIgnoresUnownedPods(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := v1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	ns := "apis"
	gw := &v1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: ns, UID: "gw-uid"}}
	controllerRef := func(apiVersion, kind, name string, uid types.UID) []metav1.OwnerReference {
		controller := true
		return []metav1.OwnerReference{{APIVersion: apiVersion, Kind: kind, Name: name, UID: uid, Controller: &controller}}
	}
	labels := func() map[string]string {
		l := util.DefaultLabels(gw.Name, map[string]string{})
		l[appsv1.DefaultDeploymentUniqueLabelKey] = "abc"
		return l
	}
	objs := []client.Object{
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: gw.Name, Namespace: ns, UID: "dep-uid",
			OwnerReferences: controllerRef("security.brcmlabs.com/v1", "Gateway", gw.Name, gw.UID)}},
		&appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: "gw-abc", Namespace: ns, UID: "rs-uid",
			OwnerReferences: controllerRef("apps/v1", "Deployment", gw.Name, "dep-uid")}},
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "gw-abc-1", Namespace: ns, Labels: labels(),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-abc", "rs-uid")}},
		// labels only
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-labels", Namespace: ns, Labels: labels()}},
		// right ReplicaSet name, forged UID
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-wrong-uid", Namespace: ns, Labels: labels(),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-abc", "forged-uid")}},
	}
	params := Params{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
		Scheme: scheme,
		Log:    logr.Discard(),
	}
	podList, err := getGatewayPods(context.Background(), params, gw)
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) != 1 || podList.Items[0].Name != "gw-abc-1" {
		names := []string{}
		for _, p := range podList.Items {
			names = append(names, p.Name)
		}
		t.Fatalf("expected only gw-abc-1, got %v", names)
	}
}
