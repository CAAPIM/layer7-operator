package gateway

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"regexp"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret
func NewSecret(gw *securityv1.Gateway, name string) (*corev1.Secret, error) {

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
	case gw.Name + "-node-properties":
		nodeProperties := fmt.Sprintf("node.cluster.pass=%s\nadmin.user=%s\nadmin.pass=%s\n", gw.Spec.App.Management.Cluster.Password, gw.Spec.App.Management.Username, gw.Spec.App.Management.Password)
		if gw.Spec.App.Management.Database.Enabled {
			nodeProperties = fmt.Sprintf("%sl7.mysql.connection.url=%s\nnode.db.config.main.user=%s\nnode.db.config.main.pass=%s", nodeProperties, gw.Spec.App.Management.Database.JDBCUrl, gw.Spec.App.Management.Database.Username, gw.Spec.App.Management.Database.Password)
		} else {
			nodeProperties = fmt.Sprintf("%snode.db.type=%s\nnode.db.config.main.user=%s", nodeProperties, "derby", "gateway")

		}
		data["node.properties"] = []byte(nodeProperties)

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
		if gw.Spec.App.Otk.Database.Auth.ClientReadOnlyUser != (securityv1.OtkDatabaseAuthCredentials{}) {
			data["OTK_CLIENT_READ_DATABASE_USERNAME"] = []byte(gw.Spec.App.Otk.Database.Auth.ClientReadOnlyUser.Username)
			data["OTK_CLIENT_READ_DATABASE_PASSWORD"] = []byte(gw.Spec.App.Otk.Database.Auth.ClientReadOnlyUser.Password)
		}

	case gw.Name + "-shared-state-config":
		redisGroupName := "l7GW"
		sentinelMasterSet := "mymaster"
		commandTimeout := 5000
		connectTimeout := 10000

		if gw.Spec.App.Redis.Default.Ssl.Enabled {
			if gw.Spec.App.Redis.Default.Ssl.Crt != "" && gw.Spec.App.Redis.Default.Ssl.ExistingSecretName == "" {
				data["redis.crt"] = []byte(gw.Spec.App.Redis.Default.Ssl.Crt)
				gw.Spec.App.Redis.Default.Ssl.Crt = "redis.crt"
			}
			if gw.Spec.App.Redis.Default.Ssl.ExistingSecretName != "" {
				gw.Spec.App.Redis.Default.Ssl.Crt = "redis.crt"
			}
		} else {
			gw.Spec.App.Redis.Default.Ssl = securityv1.RedisSsl{}
		}

		redisConfigs := []util.RedisClientConfig{}
		defaultRedisConfig := util.RedisClientConfig{}

		gw.Spec.App.Redis.Default.Name = "default"
		defaultRedisBytes, err := json.Marshal(gw.Spec.App.Redis.Default)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(defaultRedisBytes, &defaultRedisConfig)
		if err != nil {
			return nil, err
		}

		defaultRedisConfig.KeyPrefixGroupName = gw.Spec.App.Redis.Default.GroupName

		if defaultRedisConfig.KeyPrefixGroupName == "" {
			defaultRedisConfig.KeyPrefixGroupName = redisGroupName
		} else {
			defaultRedisConfig.KeyPrefixGroupName = gw.Spec.App.Redis.Default.GroupName
		}

		if gw.Spec.App.Redis.Default.CommandTimeout == 0 {
			defaultRedisConfig.CommandTimeout = commandTimeout
		}

		if gw.Spec.App.Redis.Default.ConnectTimeout == 0 {
			defaultRedisConfig.ConnectTimeout = connectTimeout
		}

		if gw.Spec.App.Redis.Default.Auth.Enabled {
			if gw.Spec.App.Redis.Default.Auth.Username != "" {
				defaultRedisConfig.Username = gw.Spec.App.Redis.Default.Auth.Username
			}

			if gw.Spec.App.Redis.Default.Auth.PasswordEncoded != "" && gw.Spec.App.Redis.Default.Auth.PasswordPlainText != "" {
				return nil, fmt.Errorf("invalid redis configuration for %s provide one password type", gw.Spec.App.Redis.Default.Name)
			}
			if gw.Spec.App.Redis.Default.Auth.PasswordEncoded != "" {
				defaultRedisConfig.EncodedPassword = gw.Spec.App.Redis.Default.Auth.PasswordEncoded
			}
			if gw.Spec.App.Redis.Default.Auth.PasswordPlainText != "" {
				defaultRedisConfig.Password = gw.Spec.App.Redis.Default.Auth.PasswordPlainText
			}
		}

		if len(gw.Spec.App.Redis.Default.Sentinel.Nodes) > 0 {
			defaultRedisConfig.Sentinel.Master = gw.Spec.App.Redis.Default.Sentinel.MasterSet
			if gw.Spec.App.Redis.Default.Sentinel.MasterSet == "" {
				defaultRedisConfig.Sentinel.Master = sentinelMasterSet
			}
		}

		redisConfigs = append(redisConfigs, defaultRedisConfig)

		if len(gw.Spec.App.Redis.AdditionalConfigs) > 0 {
			for _, rc := range gw.Spec.App.Redis.AdditionalConfigs {
				if rc.Enabled {
					if rc.Ssl.Enabled {
						if rc.Ssl.Crt != "" && rc.Ssl.ExistingSecretName == "" {
							data[rc.Name+"-redis.crt"] = []byte(rc.Ssl.Crt)
							rc.Ssl.Crt = rc.Name + "-redis.crt"
						}
						if rc.Ssl.ExistingSecretName != "" {
							rc.Ssl.Crt = rc.Name + "-redis.crt"
						}
					} else {
						rc.Ssl = securityv1.RedisSsl{}
					}
					redisConfig := util.RedisClientConfig{}
					redisBytes, err := json.Marshal(rc)
					if err != nil {
						return nil, err
					}
					err = json.Unmarshal(redisBytes, &redisConfig)
					if err != nil {
						return nil, err
					}

					redisConfig.KeyPrefixGroupName = rc.GroupName
					if redisConfig.KeyPrefixGroupName == "" {
						redisConfig.KeyPrefixGroupName = redisGroupName
					}

					if rc.CommandTimeout == 0 {
						redisConfig.CommandTimeout = commandTimeout
					}

					if rc.ConnectTimeout == 0 {
						redisConfig.ConnectTimeout = connectTimeout
					}

					if rc.Auth.Enabled {
						if rc.Auth.Username != "" {
							redisConfig.Username = rc.Auth.Username
						}

						if rc.Auth.PasswordEncoded != "" && rc.Auth.PasswordPlainText != "" {
							return nil, fmt.Errorf("invalid redis configuration for %s provide one password type", rc.Name)
						}
						if rc.Auth.PasswordEncoded != "" {
							redisConfig.EncodedPassword = rc.Auth.PasswordEncoded
						}
						if rc.Auth.PasswordPlainText != "" {
							redisConfig.Password = rc.Auth.PasswordPlainText
						}
					}

					switch rc.Type {
					case securityv1.RedisTypeSentinel:
						if len(rc.Sentinel.Nodes) == 0 {
							return nil, fmt.Errorf("redis %s sentinel requires an array of nodes that contain host and port", rc.Name)
						}
						for i, node := range rc.Sentinel.Nodes {
							if node.Host == "" || node.Port == 0 {
								return nil, fmt.Errorf("redis %s sentinel node %d requires host and port to be set", rc.Name, i)
							}
						}
					case securityv1.RedisTypeStandalone:
						if rc.Standalone.Host == "" || rc.Standalone.Port == 0 {
							return nil, fmt.Errorf("redis %s standalone requires host and port to be set", rc.Name)
						}
					default:
						return nil, fmt.Errorf("redis %s requires a type, valid options are sentinel or standalone", rc.Name)
					}

					if len(rc.Sentinel.Nodes) > 0 {
						redisConfig.Sentinel.Master = rc.Sentinel.MasterSet
						if rc.Sentinel.MasterSet == "" {
							redisConfig.Sentinel.Master = sentinelMasterSet
						}
					}

					redisConfigs = append(redisConfigs, redisConfig)
				}
			}
		}
		sharedStateClientBytes, err := util.GenerateSharedStateClientConfig(string(util.SharedStateClientConfigTypeRedis), redisConfigs, nil)
		if err != nil {
			return nil, err
		}

		ssc := string(sharedStateClientBytes)
		re1 := regexp.MustCompile(`(?m)(password:)\s(.*)`)
		re2 := regexp.MustCompile(`(?m)(encodedPassword:)\s(.*)`)
		sub1 := "password: \"$2\""
		sub2 := "encodedPassword: \"$2\""
		ssc = re1.ReplaceAllString(ssc, sub1)
		ssc = re2.ReplaceAllString(ssc, sub2)
		data["sharedstate_client.yaml"] = []byte(ssc)
		//data["sharedstate_client.yaml"] = sharedStateClientBytes
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

	return secret, nil
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
