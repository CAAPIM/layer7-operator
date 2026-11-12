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
	"testing"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/types"
)

func TestNewDeployment(t *testing.T) {
	t.Run("should create Deployment", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Deployment is created
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &appsv1.Deployment{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if *got.Spec.Replicas != int32(5) {
			t.Errorf("Expected %d, Actual %d", int32(5), *got.Spec.Replicas)
		}
	})

	t.Run("should update Deployment", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()
		err := Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		params.Instance.Spec.App.ServiceAccount.Name = "modified"
		err = Deployment(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		//verify that Deployment is updated
		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &appsv1.Deployment{}
		err = k8sClient.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}
		if got.Spec.Template.Spec.ServiceAccountName != "modified" {
			t.Errorf("Expected %s, Actual %s", "modified", got.Spec.Template.Spec.ServiceAccountName)
		}
	})

	// Guards the full-Spec-replace behaviour in updateDeployment: the pod template
	// Spec is replaced wholesale with the desired one and patched via MergeFrom, so a
	// PodLabel removed from the CR must be dropped from the live template (emitted as
	// an explicit null in the merge patch), not left behind by an additive merge.
	t.Run("should remove a PodLabel from the pod template on update", func(t *testing.T) {
		params := newParams()
		ctx := context.Background()

		// Create with two PodLabels present.
		params.Instance.Spec.App.PodLabels = map[string]string{"team": "payments", "tier": "backend"}
		if err := Deployment(ctx, params); err != nil {
			t.Fatal(err)
		}

		nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &appsv1.Deployment{}
		if err := k8sClient.Get(ctx, nns, got); err != nil {
			t.Fatal(err)
		}
		if got.Spec.Template.Labels["team"] != "payments" || got.Spec.Template.Labels["tier"] != "backend" {
			t.Fatalf("expected both PodLabels on template, got %+v", got.Spec.Template.Labels)
		}

		// Remove "tier" from the CR and reconcile again.
		params.Instance.Spec.App.PodLabels = map[string]string{"team": "payments"}
		if err := Deployment(ctx, params); err != nil {
			t.Fatal(err)
		}

		got = &appsv1.Deployment{}
		if err := k8sClient.Get(ctx, nns, got); err != nil {
			t.Fatal(err)
		}
		if _, ok := got.Spec.Template.Labels["tier"]; ok {
			t.Errorf("removed PodLabel %q still present on template: %+v", "tier", got.Spec.Template.Labels)
		}
		if got.Spec.Template.Labels["team"] != "payments" {
			t.Errorf("retained PodLabel %q missing: %+v", "team", got.Spec.Template.Labels)
		}
		// Default selector label must survive the removal round-trip.
		if got.Spec.Template.Labels["app.kubernetes.io/name"] != "test" {
			t.Errorf("default label lost after PodLabel removal: %+v", got.Spec.Template.Labels)
		}
	})
}
