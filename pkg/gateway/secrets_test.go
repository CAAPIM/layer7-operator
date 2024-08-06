package gateway

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaultGatewayDbSecret(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}
	gateway.Spec.App.Management = securityv1.Management{
		Username: "testUser",
		Password: "testPassword",
		Cluster: securityv1.Cluster{
			Password: "testClusterPassword",
		},
		Database: securityv1.Database{
			Enabled:  true,
			JDBCUrl:  "jdbc:mysql:localhost:3306",
			Username: "testDBUser",
			Password: "testDBPassword"},
	}

	secret, _ := NewSecret(&gateway, gateway.Name)

	expectedSecretData := map[string][]byte{
		"SSG_ADMIN_USERNAME":    []byte(gateway.Spec.App.Management.Username),
		"SSG_ADMIN_PASSWORD":    []byte(gateway.Spec.App.Management.Password),
		"SSG_CLUSTER_PASSWORD":  []byte(gateway.Spec.App.Management.Cluster.Password),
		"SSG_DATABASE_PASSWORD": []byte(gateway.Spec.App.Management.Database.Password),
		"SSG_DATABASE_USER":     []byte(gateway.Spec.App.Management.Database.Username),
	}

	for i := range secret.Data {
		if string(secret.Data[i]) != string(expectedSecretData[i]) {
			t.Errorf("expected %s, actual %s", string(expectedSecretData[i]), string(secret.Data[i]))

		}
	}
}

func TestDefaultGatewayEphemeralSecret(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}
	gateway.Spec.App.Management = securityv1.Management{
		Username: "testUser",
		Password: "testPassword",
		Cluster: securityv1.Cluster{
			Password: "testClusterPassword",
		},
	}

	secret, _ := NewSecret(&gateway, gateway.Name)

	expectedSecretData := map[string][]byte{
		"SSG_ADMIN_USERNAME":   []byte(gateway.Spec.App.Management.Username),
		"SSG_ADMIN_PASSWORD":   []byte(gateway.Spec.App.Management.Password),
		"SSG_CLUSTER_PASSWORD": []byte(gateway.Spec.App.Management.Cluster.Password),
	}

	for i := range secret.Data {
		if string(secret.Data[i]) != string(expectedSecretData[i]) {
			t.Errorf("expected %s, actual %s", string(expectedSecretData[i]), string(secret.Data[i]))

		}
	}
}

func TestNodePropertiesSecretDb(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}
	gateway.Spec.App.Management = securityv1.Management{
		Username: "testUser",
		Password: "testPassword",
		DisklessConfig: securityv1.DisklessConfig{
			Disabled: true,
		},
		Cluster: securityv1.Cluster{
			Password: "testClusterPassword",
		},
		Database: securityv1.Database{
			Enabled:  true,
			JDBCUrl:  "jdbc:mysql:localhost:3306",
			Username: "testDBUser",
			Password: "testDBPassword"},
	}

	secret, _ := NewSecret(&gateway, gateway.Name+"-node-properties")

	nodeProperties :=
		`
node.cluster.pass=%s
admin.user=%s
admin.pass=%s
l7.mysql.connection.url=%s
node.db.config.main.user=%s
node.db.config.main.pass=%s
`

	expectedSecretData := map[string][]byte{
		"node.properties": []byte(fmt.Sprintf(nodeProperties, gateway.Spec.App.Management.Cluster.Password, gateway.Spec.App.Management.Username, gateway.Spec.App.Management.Password, gateway.Spec.App.Management.Database.JDBCUrl, gateway.Spec.App.Management.Database.Username, gateway.Spec.App.Management.Database.Password)),
	}
	if !reflect.DeepEqual(strings.TrimSpace(string(secret.Data["node.properties"])), strings.TrimSpace(string(expectedSecretData["node.properties"]))) {
		t.Errorf("expected %s, actual %s", string(expectedSecretData["node.properties"]), string(secret.Data["node.properties"]))
	}
}

func TestNodePropertiesSecretEphemeral(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}
	gateway.Spec.App.Management = securityv1.Management{
		Username: "testUser",
		Password: "testPassword",
		DisklessConfig: securityv1.DisklessConfig{
			Disabled: true,
		},
		Cluster: securityv1.Cluster{
			Password: "testClusterPassword",
		},
	}

	secret, _ := NewSecret(&gateway, gateway.Name+"-node-properties")

	nodeProperties :=
		`
node.cluster.pass=%s
admin.user=%s
admin.pass=%s
node.db.type=derby
node.db.config.main.user=gateway
`

	expectedSecretData := fmt.Sprintf(nodeProperties, gateway.Spec.App.Management.Cluster.Password, gateway.Spec.App.Management.Username, gateway.Spec.App.Management.Password)
	if strings.TrimSpace(string(secret.Data["node.properties"])) != strings.TrimSpace(expectedSecretData) {
		t.Errorf("expected %s, actual %s", strings.TrimSpace(expectedSecretData), string(secret.Data["node.properties"]))
	}
}

func TestSharedStateClientSecretRedis(t *testing.T) {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}

	redisStandalone := securityv1.RedisConfigurations{
		Enabled: true,
		Default: securityv1.RedisConfiguration{
			Type: securityv1.RedisTypeStandalone,
			Auth: securityv1.RedisAuth{
				Enabled:           true,
				Username:          "redis-user",
				PasswordPlainText: "redis-password",
			},
			Standalone: securityv1.RedisNode{
				Host: "redis-standalone",
				Port: 6379,
			},
		},
	}

	redisSentinel := securityv1.RedisConfigurations{
		Enabled: true,
		Default: securityv1.RedisConfiguration{
			Type: securityv1.RedisTypeSentinel,
			Auth: securityv1.RedisAuth{
				Enabled:         true,
				Username:        "redis-user",
				PasswordEncoded: "wyN0kCr15hI.O37BlXCmrYS5V24l2MH1yg",
			},
			Sentinel: securityv1.RedisSentinel{
				MasterSet: "mymaster",
				Nodes: []securityv1.RedisNode{
					{
						Host: "redis-sentinel-1",
						Port: 26379,
					},
					{
						Host: "redis-sentinel-2",
						Port: 26379,
					},
					{
						Host: "redis-sentinel-3",
						Port: 26379,
					},
				},
			},
		},
	}

	gateway.Spec.App.Redis = redisStandalone

	expectedRedisConfig :=
		`
redis:
  default:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    username: redis-user
    password: "redis-password"
    standalone:
      host: redis-standalone
      port: 6379
    ssl:
      enabled: false
`

	secret, _ := NewSecret(&gateway, gateway.Name+"-shared-state-client-configuration")

	if strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])) != strings.TrimSpace(expectedRedisConfig) {
		t.Errorf("expected %s, actual %s", strings.TrimSpace(expectedRedisConfig), strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])))
	}

	gateway.Spec.App.Redis = redisSentinel
	expectedRedisConfig =
		`
redis:
  default:
    type: sentinel
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    username: redis-user
    encodedPassword: "wyN0kCr15hI.O37BlXCmrYS5V24l2MH1yg"
    sentinel:
      master: mymaster
      nodes:
        - host: redis-sentinel-1
          port: 26379
        - host: redis-sentinel-2
          port: 26379
        - host: redis-sentinel-3
          port: 26379
    ssl:
      enabled: false
`

	secret, _ = NewSecret(&gateway, gateway.Name+"-shared-state-client-configuration")

	if strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])) != strings.TrimSpace(expectedRedisConfig) {
		t.Errorf("expected %s, actual %s", strings.TrimSpace(expectedRedisConfig), strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])))
	}

}
