/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
package reconcile

import (
	"context"
	"sort"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
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

func controllerRef(apiVersion, kind, name string, uid types.UID) []metav1.OwnerReference {
	controller := true
	return []metav1.OwnerReference{{APIVersion: apiVersion, Kind: kind, Name: name, UID: uid, Controller: &controller}}
}

// newGatewayPodsTestParams builds a fake cluster containing the operator managed Deployment for a Gateway, two of its
// pods, and pods that carry the Gateway's labels but are not owned by that Deployment.
func newGatewayPodsTestParams(t *testing.T) Params {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := securityv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	ns := "default"
	gw := &securityv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: ns, UID: "gw-uid"}}
	labels := util.DefaultLabels(gw.Name, map[string]string{})
	podLabels := func(extra map[string]string) map[string]string {
		l := map[string]string{}
		for k, v := range labels {
			l[k] = v
		}
		for k, v := range extra {
			l[k] = v
		}
		return l
	}
	running := corev1.PodStatus{Phase: corev1.PodRunning}

	hash := func(h string, extra map[string]string) map[string]string {
		l := podLabels(extra)
		l[appsv1.DefaultDeploymentUniqueLabelKey] = h
		return l
	}
	objs := []client.Object{
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: gw.Name, Namespace: ns, UID: "dep-uid",
			OwnerReferences: controllerRef("security.brcmlabs.com/v1", "Gateway", gw.Name, gw.UID)}},
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "gw-abc-1", Namespace: ns, Labels: hash("abc", nil),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-abc", "rs-uid")}, Status: running},
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "gw-abc-2", Namespace: ns, Labels: hash("abc", nil),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-abc", "rs-uid")}, Status: running},
		// labels only, already claiming leadership
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-labels", Namespace: ns, Labels: hash("abc", map[string]string{"management-access": "leader"})}, Status: running},
		// controlled by a ReplicaSet of another Deployment
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-other-rs", Namespace: ns, Labels: hash("abc", nil),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "other-abc", "other-rs-uid")}, Status: running},
		// owner name doesn't match its pod-template-hash label
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-hash", Namespace: ns, Labels: hash("xyz", nil),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-abc", "rs-uid")}, Status: running},
		// no pod-template-hash label
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-no-hash", Namespace: ns, Labels: podLabels(nil),
			OwnerReferences: controllerRef("apps/v1", "ReplicaSet", "gw-", "rs-uid")}, Status: running},
		// claims the Deployment directly
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "spoof-dep", Namespace: ns, Labels: hash("abc", nil),
			OwnerReferences: controllerRef("apps/v1", "Deployment", gw.Name, "dep-uid")}, Status: running},
	}

	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	return Params{Client: cl, Scheme: scheme, Log: logr.Discard(), Instance: gw}
}

func podNames(pods []corev1.Pod) []string {
	names := []string{}
	for _, p := range pods {
		names = append(names, p.Name)
	}
	sort.Strings(names)
	return names
}

func TestGetGatewayPodsOnlyReturnsDeploymentPods(t *testing.T) {
	params := newGatewayPodsTestParams(t)
	podList, err := getGatewayPods(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	got := podNames(podList.Items)
	if len(got) != 2 || got[0] != "gw-abc-1" || got[1] != "gw-abc-2" {
		t.Fatalf("expected only deployment pods, got %v", got)
	}
}

func TestGetGatewayPodsDeploymentNotControlledByGateway(t *testing.T) {
	params := newGatewayPodsTestParams(t)
	params.Instance.UID = "different-gw-uid"
	podList, err := getGatewayPods(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) != 0 {
		t.Fatalf("expected no pods, got %v", podNames(podList.Items))
	}
}

func TestGetGatewayPodsNoDeployment(t *testing.T) {
	params := newGatewayPodsTestParams(t)
	if err := params.Client.Delete(context.Background(), &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "gw", Namespace: "default"}}); err != nil {
		t.Fatal(err)
	}
	podList, err := getGatewayPods(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) != 0 {
		t.Fatalf("expected no pods, got %v", podNames(podList.Items))
	}
}

func TestManagementPodIgnoresSpoofedLeader(t *testing.T) {
	params := newGatewayPodsTestParams(t)
	ctx := context.Background()
	if err := ManagementPod(ctx, params); err != nil {
		t.Fatal(err)
	}

	pods := &corev1.PodList{}
	if err := params.Client.List(ctx, pods, client.InNamespace("default")); err != nil {
		t.Fatal(err)
	}
	leaders := []string{}
	for _, p := range pods.Items {
		if p.Labels["management-access"] == "leader" && p.Name != "spoof-labels" {
			leaders = append(leaders, p.Name)
		}
	}
	if len(leaders) != 1 || (leaders[0] != "gw-abc-1" && leaders[0] != "gw-abc-2") {
		t.Fatalf("expected exactly one deployment pod elected leader, got %v", leaders)
	}
}
