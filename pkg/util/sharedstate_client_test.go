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

	configBytes, err := GenerateSharedStateClientConfig("redis", rcc, nil)
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

	configBytes, err := GenerateSharedStateClientConfig("redis", rcc, nil)
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

	configBytes, err := GenerateSharedStateClientConfig("redis", rcc, nil)
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

func TestUnsupportedTypeSharedStateClientConfiguration(t *testing.T) {

	_, err := GenerateSharedStateClientConfig("hazelcast", nil, nil)
	if err == nil {
		t.Errorf("hazelcast support is not implemented, should be rejected err: %v", err)
	}

	_, err = GenerateSharedStateClientConfig("gemfire", nil, nil)
	if err == nil {
		t.Errorf("gemfire support is not implemented, should be rejected err: %v", err)
	}
}
