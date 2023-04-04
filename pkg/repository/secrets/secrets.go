package secrets

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret
func NewSecret(repository *securityv1.Repository) *corev1.Secret {

	data := make(map[string][]byte)

	data["USERNAME"] = []byte(repository.Spec.Auth.Username)
	data["PASSWORD"] = []byte(repository.Spec.Auth.Password)
	data["TOKEN"] = []byte(repository.Spec.Auth.Token)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      repository.Name,
			Namespace: repository.Namespace,
			Labels:    util.DefaultLabels(repository.Name, repository.Spec.Labels),
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

// NewSecret
func NewStorageSecret(repository *securityv1.Repository, bundle []byte) *corev1.Secret {

	data := make(map[string][]byte)

	data[repository.Name+".gz"] = bundle

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      repository.Name + "-repository",
			Namespace: repository.Namespace,
			Labels:    util.DefaultLabels(repository.Name, repository.Spec.Labels),
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
