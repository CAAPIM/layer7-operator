package util

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	defaultName    = "default"
	verifyPeer     = false
	redisGroupName = "l7GW"
	commandTimeout = 5000
	connectTimeout = 10000
)

func TestSingleStandaloneSharedStateClientConfiguration(t *testing.T) {
	rcc := []RedisClientConfig{
		{
			Name:               defaultName,
			Type:               "standalone",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Standalone: RedisNode{
				Host: "localhost",
				Port: 6379,
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
	}

	configBytes, err := GenerateSharedStateClientConfig(rcc, nil, nil)
	if err != nil {
		t.Errorf("failed to generate sharedstate_client config for single client configuration test err: %v", err)
	}
	actual := RedisClientConfig{}

	err = yaml.Unmarshal(configBytes, &actual)
	if err != nil {
		t.Errorf("failed to unmarshal sharedstate_client config for single client configuration test err: %v", err)
	}

	expected := `redis:
  default:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    standalone:
      host: localhost
      port: 6379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false`

	if !reflect.DeepEqual(strings.TrimSpace(string(configBytes)), strings.TrimSpace(expected)) {
		t.Errorf("actual \n%v, expected \n%v", string(configBytes), expected)
	}
}

func TestSingleSentinelSharedStateClientConfiguration(t *testing.T) {
	rcc := []RedisClientConfig{
		{
			Name:               defaultName,
			Type:               "sentinel",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Sentinel: RedisSentinel{
				Master: "mymaster",
				Nodes: []RedisNode{
					{
						Host: "sentinel-redis-node-0.sentinel-redis-headless.develop.svc.cluster.local",
						Port: 26379,
					},
					{
						Host: "sentinel-redis-node-1.sentinel-redis-headless.develop.svc.cluster.local",
						Port: 26379,
					},
					{
						Host: "sentinel-redis-node-2.sentinel-redis-headless.develop.svc.cluster.local",
						Port: 26379,
					},
				},
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
	}

	configBytes, err := GenerateSharedStateClientConfig(rcc, nil, nil)
	if err != nil {
		t.Errorf("failed to generate sharedstate_client config for single client configuration test err: %v", err)
	}
	actual := RedisClientConfig{}

	err = yaml.Unmarshal(configBytes, &actual)
	if err != nil {
		t.Errorf("failed to unmarshal sharedstate_client config for single client configuration test err: %v", err)
	}

	expected := `redis:
  default:
    type: sentinel
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    sentinel:
      master: mymaster
      nodes:
        - host: sentinel-redis-node-0.sentinel-redis-headless.develop.svc.cluster.local
          port: 26379
        - host: sentinel-redis-node-1.sentinel-redis-headless.develop.svc.cluster.local
          port: 26379
        - host: sentinel-redis-node-2.sentinel-redis-headless.develop.svc.cluster.local
          port: 26379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false`

	if !reflect.DeepEqual(strings.TrimSpace(string(configBytes)), strings.TrimSpace(expected)) {
		t.Errorf("actual \n%v, expected \n%v", string(configBytes), expected)
	}
}

func TestAdditionalProviderSharedStateClientConfiguration(t *testing.T) {
	rcc := []RedisClientConfig{
		{
			Name:               defaultName,
			Type:               "standalone",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Standalone: RedisNode{
				Host: "localhost",
				Port: 6379,
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
		{
			Name:               "local",
			Type:               "standalone",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Standalone: RedisNode{
				Host: "dc1.example.com",
				Port: 6379,
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
		{
			Name:               "regional",
			Type:               "standalone",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Standalone: RedisNode{
				Host: "regional.example.com",
				Port: 6379,
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
	}

	configBytes, err := GenerateSharedStateClientConfig(rcc, nil, nil)
	if err != nil {
		t.Errorf("failed to generate sharedstate_client config for single client configuration test err: %v", err)
	}
	actual := RedisClientConfig{}

	err = yaml.Unmarshal(configBytes, &actual)
	if err != nil {
		t.Errorf("failed to unmarshal sharedstate_client config for single client configuration test err: %v", err)
	}

	expected := `redis:
  default:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    standalone:
      host: localhost
      port: 6379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false
  local:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    standalone:
      host: dc1.example.com
      port: 6379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false
  regional:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    standalone:
      host: regional.example.com
      port: 6379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false`

	if !reflect.DeepEqual(strings.TrimSpace(string(configBytes)), strings.TrimSpace(expected)) {
		t.Errorf("actual \n%v, expected \n%v", string(configBytes), expected)
	}
}

func TestEmptySharedStateClientConfiguration(t *testing.T) {

	_, err := GenerateSharedStateClientConfig(nil, nil, nil)
	if err == nil {
		t.Errorf("empty configuration should be rejected err: %v", err)
	}
}

func TestGemfireSharedStateClientConfiguration(t *testing.T) {
	gcc := &GemfireClientConfig{
		TestOnStart: true,
		Locators: []GemfireLocator{
			{Host: "locator-0.gemfire.develop.svc.cluster.local", Port: 10334},
			{Host: "locator-1.gemfire.develop.svc.cluster.local", Port: 10334},
		},
		Username:                "gateway",
		Password:                "7layer",
		GwKeyValueRegionName:    "layer7gw_keyvalue",
		GwCounterRegionName:     "layer7gw_counter",
		GwRateLimiterRegionName: "layer7gw_ratelimiter",
		GwSortedSetRegionName:   "layer7gw_sortedset",
		Ssl: GemfireSsl{
			Enabled:            true,
			EnabledComponents:  "all",
			Keystore:           "keystore.jks",
			KeystorePassword:   "keystorepass",
			KeystoreType:       "JKS",
			Truststore:         "truststore.jks",
			TruststorePassword: "truststorepass",
			TruststoreType:     "JKS",
		},
	}

	configBytes, err := GenerateSharedStateClientConfig(nil, nil, gcc)
	if err != nil {
		t.Errorf("failed to generate sharedstate_client config for gemfire configuration test err: %v", err)
	}

	expected := `gemfire:
  testOnStart: true
  locators:
    - host: locator-0.gemfire.develop.svc.cluster.local
      port: 10334
    - host: locator-1.gemfire.develop.svc.cluster.local
      port: 10334
  username: gateway
  password: 7layer
  gwKeyValueRegionName: layer7gw_keyvalue
  gwCounterRegionName: layer7gw_counter
  gwRateLimiterRegionName: layer7gw_ratelimiter
  gwSortedSetRegionName: layer7gw_sortedset
  ssl:
    enabled: true
    enabledComponents: all
    keystore: keystore.jks
    keystorePassword: keystorepass
    keystoreType: JKS
    truststore: truststore.jks
    truststorePassword: truststorepass
    truststoreType: JKS`

	if !reflect.DeepEqual(strings.TrimSpace(string(configBytes)), strings.TrimSpace(expected)) {
		t.Errorf("actual \n%v, expected \n%v", string(configBytes), expected)
	}
}

func TestCombinedRedisGemfireSharedStateClientConfiguration(t *testing.T) {
	rcc := []RedisClientConfig{
		{
			Name:               defaultName,
			Type:               "standalone",
			KeyPrefixGroupName: redisGroupName,
			CommandTimeout:     commandTimeout,
			ConnectTimeout:     connectTimeout,
			Password:           "7layer",
			Standalone: RedisNode{
				Host: "localhost",
				Port: 6379,
			},
			Ssl: RedisSsl{
				Enabled: true, Crt: "------",
				VerifyPeer: &verifyPeer,
			},
		},
	}

	gcc := &GemfireClientConfig{
		Locators: []GemfireLocator{
			{Host: "locator-0.gemfire.develop.svc.cluster.local", Port: 10334},
		},
		Username: "gateway",
		Password: "7layer",
	}

	configBytes, err := GenerateSharedStateClientConfig(rcc, nil, gcc)
	if err != nil {
		t.Errorf("failed to generate sharedstate_client config for combined redis/gemfire configuration test err: %v", err)
	}

	expected := `gemfire:
  testOnStart: false
  locators:
    - host: locator-0.gemfire.develop.svc.cluster.local
      port: 10334
  username: gateway
  password: 7layer
redis:
  default:
    type: standalone
    commandTimeout: 5000
    connectTimeout: 10000
    keyPrefixGroupName: l7GW
    testOnStart: false
    password: 7layer
    standalone:
      host: localhost
      port: 6379
    ssl:
      enabled: true
      cert: '------'
      verifyPeer: false`

	if !reflect.DeepEqual(strings.TrimSpace(string(configBytes)), strings.TrimSpace(expected)) {
		t.Errorf("actual \n%v, expected \n%v", string(configBytes), expected)
	}
}

// redisVerifyPeerFromYAML returns ssl.verifyPeer for each redis provider in the
// rendered config (nil if the key is absent).
func redisVerifyPeerFromYAML(t *testing.T, configBytes []byte) map[string]*bool {
	t.Helper()
	var doc struct {
		Redis map[string]struct {
			Ssl struct {
				VerifyPeer *bool `yaml:"verifyPeer"`
			} `yaml:"ssl"`
		} `yaml:"redis"`
	}
	if err := yaml.Unmarshal(configBytes, &doc); err != nil {
		t.Fatalf("unmarshal rendered config: %v", err)
	}
	out := map[string]*bool{}
	for name, rc := range doc.Redis {
		out[name] = rc.Ssl.VerifyPeer
	}
	return out
}

func TestSharedStateClientVerifyPeerDefaultsTrue(t *testing.T) {
	explicitFalse := false
	explicitTrue := true
	node := RedisNode{Host: "localhost", Port: 6379}

	t.Run("single", func(t *testing.T) {
		cases := []struct {
			name string
			in   *bool
			want bool
		}{
			{"unset defaults to true", nil, true},
			{"explicit false preserved", &explicitFalse, false},
			{"explicit true preserved", &explicitTrue, true},
		}
		for _, tc := range cases {
			rcc := []RedisClientConfig{{
				Name: defaultName, Type: "standalone", Standalone: node,
				Ssl: RedisSsl{Enabled: true, VerifyPeer: tc.in},
			}}
			configBytes, err := GenerateSharedStateClientConfig(rcc, nil, nil)
			if err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			got := redisVerifyPeerFromYAML(t, configBytes)["default"]
			if got == nil || *got != tc.want {
				t.Errorf("%s: verifyPeer = %v, want %v", tc.name, got, tc.want)
			}
		}
	})

	t.Run("multiple", func(t *testing.T) {
		rcc := []RedisClientConfig{
			{Name: defaultName, Type: "standalone", Standalone: node, Ssl: RedisSsl{Enabled: true}},
			{Name: "unset", Type: "standalone", Standalone: node, Ssl: RedisSsl{Enabled: true}},
			{Name: "optout", Type: "standalone", Standalone: node, Ssl: RedisSsl{Enabled: true, VerifyPeer: &explicitFalse}},
			{Name: "plain", Type: "standalone", Standalone: node},
		}
		configBytes, err := GenerateSharedStateClientConfig(rcc, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		got := redisVerifyPeerFromYAML(t, configBytes)
		for name, want := range map[string]bool{"default": true, "unset": true, "optout": false} {
			if got[name] == nil || *got[name] != want {
				t.Errorf("%s: verifyPeer = %v, want %v", name, got[name], want)
			}
		}
		// SSL disabled: verifyPeer is not rendered at all.
		if got["plain"] != nil {
			t.Errorf("plain: verifyPeer rendered (%v) although ssl is disabled", *got["plain"])
		}
	})
}
