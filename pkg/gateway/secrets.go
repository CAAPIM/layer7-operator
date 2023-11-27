package gateway

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
func NewSecret(gw *securityv1.Gateway, name string) *corev1.Secret {

	data := make(map[string][]byte)
	dataCheckSum := ""
	switch name {
	case gw.Name:
		data["SSG_ADMIN_USERNAME"] = []byte(gw.Spec.App.Management.Username)
		data["SSG_ADMIN_PASSWORD"] = []byte(gw.Spec.App.Management.Password)
		data["SSG_CLUSTER_PASSWORD"] = []byte(gw.Spec.App.Management.Cluster.Password)

		if gw.Spec.App.Management.Database.Enabled {
			data["SSG_DATABASE_PASSWORD"] = []byte(gw.Spec.App.Management.Database.Password)
			data["SSG_DATABASE_USER"] = []byte(gw.Spec.App.Management.Database.Username)
		}
	case gw.Name + "-otk-db-credentials":
		if gw.Spec.App.Otk.Database.Auth.GatewayUser != (securityv1.OtkDatabaseAuthCredentials{}) {
			data["OTK_DATABASE_USERNAME"] = []byte(gw.Spec.App.Otk.Database.Auth.GatewayUser.Username)
			data["OTK_DATABASE_PASSWORD"] = []byte(gw.Spec.App.Otk.Database.Auth.GatewayUser.Password)
		}
		if gw.Spec.App.Otk.Database.Auth.AdminUser != (securityv1.OtkDatabaseAuthCredentials{}) {
			data["OTK_DATABASE_DDL_USERNAME"] = []byte(gw.Spec.App.Otk.Database.Auth.AdminUser.Username)
			data["OTK_DATABASE_DDL_PASSWORD"] = []byte(gw.Spec.App.Otk.Database.Auth.AdminUser.Password)
		}
		if gw.Spec.App.Otk.Database.Auth.ReadOnlyUser != (securityv1.OtkDatabaseAuthCredentials{}) {
			data["OTK_RO_DATABASE_USERNAME"] = []byte(gw.Spec.App.Otk.Database.Auth.ReadOnlyUser.Username)
			data["OTK_RO_DATABASE_PASSWORD"] = []byte(gw.Spec.App.Otk.Database.Auth.ReadOnlyUser.Password)
		}

	case gw.Name + "-otk-dmz-certificates":

	case gw.Name + "-otk-internal-certificates":
	}

	if dataCheckSum == "" {
		dataBytes, _ := json.Marshal(data)
		h := sha1.New()
		h.Write(dataBytes)
		sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
		dataCheckSum = sha1Sum
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   gw.Namespace,
			Labels:      util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
			Annotations: map[string]string{"checksum/data": dataCheckSum}, // TODO: add default annotations
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

func NewOtkCertificateSecret(gw *securityv1.Gateway, name string, data map[string][]byte) *corev1.Secret {
	dataBytes, _ := json.Marshal(data)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	dataCheckSum := sha1Sum
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   gw.Namespace,
			Labels:      util.DefaultLabels(gw.Name, gw.Spec.App.Labels),
			Annotations: map[string]string{"checksum/data": dataCheckSum}, // TODO: add default annotations
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
