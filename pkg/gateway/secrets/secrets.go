package secrets

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret
func NewSecret(gw *securityv1.Gateway) *corev1.Secret {

	data := make(map[string][]byte)

	data["SSG_ADMIN_USERNAME"] = []byte(gw.Spec.App.Management.Username)
	data["SSG_ADMIN_PASSWORD"] = []byte(gw.Spec.App.Management.Password)
	data["SSG_CLUSTER_PASSWORD"] = []byte(gw.Spec.App.Management.Cluster.Password)

	if gw.Spec.App.Management.Database.Enabled {
		data["SSG_DATABASE_PASSWORD"] = []byte(gw.Spec.App.Management.Database.Password)
		data["SSG_DATABASE_USER"] = []byte(gw.Spec.App.Management.Database.Username)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}

	return secret
}
