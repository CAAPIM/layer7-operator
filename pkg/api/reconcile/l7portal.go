package reconcile

import (
	"context"
	"fmt"

	v1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
)

func L7Portal(ctx context.Context, params Params) error {
	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
		currentPortal := v1alpha1.L7Portal{}
		err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Spec.L7Portal, Namespace: params.Instance.Namespace}, &currentPortal)
		if err != nil {
			if k8serrors.IsNotFound(err) {
				newPortal := l7Portal(params)
				if err = params.Client.Create(ctx, newPortal); err != nil {
					return fmt.Errorf("failed creating l7Portal: %w", err)
				}
			}
		}
	}
	return nil
}

func l7Portal(params Params) *v1alpha1.L7Portal {
	portal := v1alpha1.L7Portal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Instance.Spec.L7Portal,
			Namespace: params.Instance.Namespace,
			Labels:    util.DefaultLabels(params.Instance.Spec.L7Portal, map[string]string{}),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1alpha1",
			Kind:       "L7Portal",
		},
		Spec: v1alpha1.L7PortalSpec{
			Enabled:        true,
			DeploymentTags: params.Instance.Spec.DeploymentTags,
		},
	}
	return &portal
}
