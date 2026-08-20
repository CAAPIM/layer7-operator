package gateway

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"gopkg.in/yaml.v3"
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

	secret, _ := NewSecret(&gateway, gateway.Name+"-shared-state-config")

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

	secret, _ = NewSecret(&gateway, gateway.Name+"-shared-state-config")

	if strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])) != strings.TrimSpace(expectedRedisConfig) {
		t.Errorf("expected %s, actual %s", strings.TrimSpace(expectedRedisConfig), strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])))
	}

}

func gemfireGateway() securityv1.Gateway {
	gateway := securityv1.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		}}
	gateway.Spec.App.Gemfire = securityv1.GemfireConfigurations{
		Enabled: true,
		Locators: []securityv1.GemfireLocator{
			{Host: "locator-0.gemfire.svc.cluster.local", Port: 10334},
			{Host: "locator-1.gemfire.svc.cluster.local", Port: 10334},
		},
		Auth: securityv1.GemfireAuth{
			Enabled:           true,
			Username:          "gateway",
			PasswordPlainText: "7layer",
		},
	}
	return gateway
}

func TestSharedStateClientSecretGemfire(t *testing.T) {
	t.Run("renders default region names", func(t *testing.T) {
		gateway := gemfireGateway()

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := `
gemfire:
  testOnStart: false
  locators:
    - host: locator-0.gemfire.svc.cluster.local
      port: 10334
    - host: locator-1.gemfire.svc.cluster.local
      port: 10334
  username: gateway
  password: "7layer"
  gwKeyValueRegionName: layer7gw_keyvalue
  gwCounterRegionName: layer7gw_counter
  gwRateLimiterRegionName: layer7gw_ratelimiter
  gwSortedSetRegionName: layer7gw_sortedset
`
		if strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])) != strings.TrimSpace(expected) {
			t.Errorf("expected %s, actual %s", strings.TrimSpace(expected), strings.TrimSpace(string(secret.Data["sharedstate_client.yaml"])))
		}
	})

	t.Run("custom region names override defaults", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.GwKeyValueRegionName = "custom_keyvalue"
		gateway.Spec.App.Gemfire.GwCounterRegionName = "custom_counter"
		gateway.Spec.App.Gemfire.GwRateLimiterRegionName = "custom_ratelimiter"
		gateway.Spec.App.Gemfire.GwSortedSetRegionName = "custom_sortedset"

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		for _, want := range []string{
			"gwKeyValueRegionName: custom_keyvalue",
			"gwCounterRegionName: custom_counter",
			"gwRateLimiterRegionName: custom_ratelimiter",
			"gwSortedSetRegionName: custom_sortedset",
		} {
			if !strings.Contains(rendered, want) {
				t.Errorf("expected rendered config to contain %q, got:\n%s", want, rendered)
			}
		}
	})

	t.Run("auth encoded password", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.Auth.PasswordPlainText = ""
		gateway.Spec.App.Gemfire.Auth.PasswordEncoded = "wyN0kCr15hI.O37BlXCmrYS5V24l2MH1yg"

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		if !strings.Contains(rendered, `encodedPassword: "wyN0kCr15hI.O37BlXCmrYS5V24l2MH1yg"`) {
			t.Errorf("expected rendered config to contain encodedPassword, got:\n%s", rendered)
		}
		if strings.Contains(rendered, "password: ") {
			t.Errorf("did not expect plaintext password field, got:\n%s", rendered)
		}
	})

	t.Run("dynamic properties passthrough", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.DynamicProperties = map[string]string{"gemfire.some-property": "value"}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		if !strings.Contains(rendered, "dynamicProperties:") || !strings.Contains(rendered, "gemfire.some-property: value") {
			t.Errorf("expected rendered config to contain dynamicProperties, got:\n%s", rendered)
		}
	})

	t.Run("ssl enabled with keystore only omits truststore key", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled: true,
			Keystore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-keystore-secret",
				PasswordPlainText:  "keystorepass",
			},
		}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		if !strings.Contains(rendered, "keystore: "+gemfireProvidersDir+"keystore.jks") {
			t.Errorf("expected rendered config to reference the full keystore path, got:\n%s", rendered)
		}
		if strings.Contains(rendered, "truststore:") {
			t.Errorf("did not expect a truststore reference when only keystore is configured, got:\n%s", rendered)
		}
	})

	t.Run("ssl enabled with truststore only omits keystore key", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled: true,
			Truststore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-truststore-secret",
				PasswordPlainText:  "truststorepass",
			},
		}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		if !strings.Contains(rendered, "truststore: "+gemfireProvidersDir+"truststore.jks") {
			t.Errorf("expected rendered config to reference the full truststore path, got:\n%s", rendered)
		}
		if strings.Contains(rendered, "keystore:") {
			t.Errorf("did not expect a keystore reference when only truststore is configured, got:\n%s", rendered)
		}
	})

	t.Run("ssl enabled with both stores", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled:           true,
			EnabledComponents: "all",
			KeystoreType:      "JKS",
			TruststoreType:    "JKS",
			Keystore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-keystore-secret",
				PasswordPlainText:  "keystorepass",
			},
			Truststore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-truststore-secret",
				PasswordEncoded:    "dHJ1c3RzdG9yZXBhc3M=",
			},
		}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		for _, want := range []string{
			"keystore: " + gemfireProvidersDir + "keystore.jks",
			"keystorePassword: \"keystorepass\"",
			"truststore: " + gemfireProvidersDir + "truststore.jks",
			"truststorePassword: \"dHJ1c3RzdG9yZXBhc3M=\"",
			"enabledComponents: all",
		} {
			if !strings.Contains(rendered, want) {
				t.Errorf("expected rendered config to contain %q, got:\n%s", want, rendered)
			}
		}
	})

	t.Run("keystore/truststore passwords with yaml-significant characters round-trip intact", func(t *testing.T) {
		const keystorePass = "pass: word#1"
		const truststorePass = "trust: word#2"

		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled: true,
			Keystore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-keystore-secret",
				PasswordPlainText:  keystorePass,
			},
			Truststore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-truststore-secret",
				PasswordPlainText:  truststorePass,
			},
		}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := secret.Data["sharedstate_client.yaml"]
		var parsed struct {
			Gemfire struct {
				Ssl struct {
					KeystorePassword   string `yaml:"keystorePassword"`
					TruststorePassword string `yaml:"truststorePassword"`
				} `yaml:"ssl"`
			} `yaml:"gemfire"`
		}
		if err := yaml.Unmarshal(rendered, &parsed); err != nil {
			t.Fatalf("rendered sharedstate_client.yaml is not valid yaml: %v\n%s", err, string(rendered))
		}

		if parsed.Gemfire.Ssl.KeystorePassword != keystorePass {
			t.Errorf("expected keystorePassword %q, got %q (rendered:\n%s)", keystorePass, parsed.Gemfire.Ssl.KeystorePassword, string(rendered))
		}
		if parsed.Gemfire.Ssl.TruststorePassword != truststorePass {
			t.Errorf("expected truststorePassword %q, got %q (rendered:\n%s)", truststorePass, parsed.Gemfire.Ssl.TruststorePassword, string(rendered))
		}
	})

	t.Run("existing secret skips operator-managed rendering", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Gemfire.ExistingSecret = "customer-managed-secret"

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, ok := secret.Data["sharedstate_client.yaml"]; ok {
			t.Errorf("expected no sharedstate_client.yaml to be rendered when gemfire.existingSecret is set, got:\n%s", string(secret.Data["sharedstate_client.yaml"]))
		}
	})

	t.Run("combined with redis renders both blocks", func(t *testing.T) {
		gateway := gemfireGateway()
		gateway.Spec.App.Redis = securityv1.RedisConfigurations{
			Enabled: true,
			Default: securityv1.RedisConfiguration{
				Type: securityv1.RedisTypeStandalone,
				Standalone: securityv1.RedisNode{
					Host: "redis-standalone",
					Port: 6379,
				},
			},
		}

		secret, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rendered := string(secret.Data["sharedstate_client.yaml"])
		if !strings.Contains(rendered, "gemfire:") || !strings.Contains(rendered, "redis:") {
			t.Errorf("expected rendered config to contain both gemfire and redis blocks, got:\n%s", rendered)
		}
	})
}

func TestSharedStateClientSecretGemfireValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(gw *securityv1.Gateway)
		wantErr string
	}{
		{
			name: "redis and gemfire existing secrets mismatch",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Redis = securityv1.RedisConfigurations{Enabled: true, ExistingSecret: "redis-secret"}
				gw.Spec.App.Gemfire.ExistingSecret = "gemfire-secret"
			},
			wantErr: "redis and gemfire share a single sharedstate_client.yaml secret",
		},
		{
			name: "no locators",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Locators = nil
			},
			wantErr: "gemfire requires an array of locators",
		},
		{
			name: "locator missing host",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Locators = []securityv1.GemfireLocator{{Host: "", Port: 10334}}
			},
			wantErr: "gemfire locator 0 requires host and port to be set",
		},
		{
			name: "locator missing port",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Locators = []securityv1.GemfireLocator{{Host: "locator-0", Port: 0}}
			},
			wantErr: "gemfire locator 0 requires host and port to be set",
		},
		{
			name: "auth both password types set",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Auth.PasswordEncoded = "encoded"
			},
			wantErr: "invalid gemfire configuration provide one password type",
		},
		{
			name: "ssl enabled with neither store configured",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{Enabled: true}
			},
			wantErr: "gemfire ssl is enabled but neither keystore.existingSecretName nor truststore.existingSecretName is set",
		},
		{
			name: "ssl keystore both password types set",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
					Enabled: true,
					Keystore: securityv1.GemfireStore{
						ExistingSecretName: "gemfire-keystore-secret",
						PasswordPlainText:  "a",
						PasswordEncoded:    "b",
					},
				}
			},
			wantErr: "invalid gemfire keystore configuration provide one password type",
		},
		{
			name: "ssl truststore both password types set",
			mutate: func(gw *securityv1.Gateway) {
				gw.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
					Enabled: true,
					Truststore: securityv1.GemfireStore{
						ExistingSecretName: "gemfire-truststore-secret",
						PasswordPlainText:  "a",
						PasswordEncoded:    "b",
					},
				}
			},
			wantErr: "invalid gemfire truststore configuration provide one password type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := gemfireGateway()
			tt.mutate(&gateway)

			_, err := NewSecret(&gateway, gateway.Name+"-shared-state-config")
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
