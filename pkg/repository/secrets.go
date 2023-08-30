package repository

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret
func NewSecret(repository *securityv1.Repository, name string, data map[string][]byte) *corev1.Secret {

	dataBytes, _ := json.Marshal(data)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   repository.Namespace,
			Labels:      util.DefaultLabels(name, repository.Spec.Labels),
			Annotations: map[string]string{"checksum/data": sha1Sum},
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
