package gateway

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

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
	case gw.Name + "-redis-properties":
		redisConfig := ""

		redisGroupName := "l7GW"

		if gw.Spec.App.Redis.GroupName != "" {
			redisGroupName = gw.Spec.App.Redis.GroupName
		}

		switch strings.ToLower(string(gw.Spec.App.Redis.Type)) {

		case string(securityv1.RedisTypeStandalone):
			redisConfig = fmt.Sprintf("redis.type=%s\nredis.standalone.hostname=%s\nredis.standalone.port=%d\nredis.key.prefix.grpname=%s", gw.Spec.App.Redis.Type, gw.Spec.App.Redis.Standalone.Hostname, gw.Spec.App.Redis.Standalone.Port, redisGroupName)
		case string(securityv1.RedisTypeSentinel):
			nodes := strings.Join(gw.Spec.App.Redis.Sentinel.Nodes, ",")
			auth := ""
			tlsConfig := ""
			publicCrt := ""
			if gw.Spec.App.Redis.Tls.Enabled {
				tlsConfig = fmt.Sprintf("redis.ssl.cert=redis.crt\nredis.ssl.verifypeer=%v", gw.Spec.App.Redis.Tls.VerifyPeer)
				if gw.Spec.App.Redis.Tls.Crt != "" {
					publicCrt = gw.Spec.App.Redis.Tls.Crt
					data["redis.crt"] = []byte(publicCrt)
				}
			}

			if gw.Spec.App.Redis.Auth.Enabled {
				username := ""
				password := ""

				if gw.Spec.App.Redis.Auth.Username != "" {
					username = "redis.sentinel.username=" + gw.Spec.App.Redis.Auth.Username + "\n"
				}

				if gw.Spec.App.Redis.Auth.PasswordPlainText != "" {
					password = "redis.sentinel.password=" + gw.Spec.App.Redis.Auth.PasswordPlainText
				}

				if gw.Spec.App.Redis.Auth.PasswordEncoded != "" {
					password = "redis.sentinel.encodedPassword=" + gw.Spec.App.Redis.Auth.PasswordEncoded
				}

				auth = fmt.Sprintf("%s%s", username, password)
			}

			redisConfig = fmt.Sprintf("redis.type=%s\nredis.sentinel.master=%s\nredis.sentinel.nodes=%s\nredis.ssl=%v\n%s\n%s", gw.Spec.App.Redis.Type, gw.Spec.App.Redis.Sentinel.MasterSet, nodes, gw.Spec.App.Redis.Tls.Enabled, tlsConfig, auth)
		}
		data["redis.properties"] = []byte(redisConfig)
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
