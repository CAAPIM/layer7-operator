package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServiceAccount(gw *securityv1.Gateway) *corev1.ServiceAccount {

	serviceAccountName := gw.Spec.App.ServiceAccount.Name
	if gw.Spec.App.ServiceAccount.Name == "" {
		serviceAccountName = gw.Name
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: gw.Namespace,
			Labels:    ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
	}
	return sa
}
