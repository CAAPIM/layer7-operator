package util

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type SharedStateClientConfigType string

var (
	SharedStateClientConfigTypeRedis     SharedStateClientConfigType = "redis"
	SharedStateClientConfigTypeGemfire   SharedStateClientConfigType = "gemfire"
	SharedStateClientConfigTypeHazelcast SharedStateClientConfigType = "hazelcast"
)

type RedisNode struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type RedisSsl struct {
	Enabled    bool   `yaml:"enabled"`
	Crt        string `yaml:"cert,omitempty"`
	VerifyPeer *bool  `yaml:"verifyPeer,omitempty"`
}

type RedisSentinel struct {
	Master string      `yaml:"master,omitempty"`
	Nodes  []RedisNode `yaml:"nodes,omitempty"`
}

type RedisClientConfig struct {
	Name               string        `yaml:"name,omitempty"`
	Type               string        `yaml:"type,omitempty"`
	CommandTimeout     int           `yaml:"commandTimeout,omitempty"`
	ConnectTimeout     int           `yaml:"connectTimeout,omitempty"`
	KeyPrefixGroupName string        `yaml:"keyPrefixGroupName,omitempty"`
	TestOnStart        bool          `yaml:"testOnStart"`
	Username           string        `yaml:"username,omitempty"`
	EncodedPassword    string        `yaml:"encodedPassword,omitempty"`
	Password           string        `yaml:"password,omitempty"`
	Standalone         RedisNode     `yaml:"standalone,omitempty"`
	Sentinel           RedisSentinel `yaml:"sentinel,omitempty"`
	Ssl                RedisSsl      `yaml:"ssl"`
}

type RedisConfigBlock struct {
	Default RedisClientConfig `yaml:"default,omitempty"`
}

type GemfireLocator struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type GemfireSsl struct {
	Enabled            bool   `yaml:"enabled,omitempty"`
	EnabledComponents  string `yaml:"enabledComponents,omitempty"`
	Keystore           string `yaml:"keystore,omitempty"`
	KeystorePassword   string `yaml:"keystorePassword,omitempty"`
	KeystoreType       string `yaml:"keystoreType,omitempty"`
	Truststore         string `yaml:"truststore,omitempty"`
	TruststorePassword string `yaml:"truststorePassword,omitempty"`
	TruststoreType     string `yaml:"truststoreType,omitempty"`
}

type GemfireClientConfig struct {
	TestOnStart             bool              `yaml:"testOnStart"`
	Locators                []GemfireLocator  `yaml:"locators,omitempty"`
	Username                string            `yaml:"username,omitempty"`
	EncodedPassword         string            `yaml:"encodedPassword,omitempty"`
	Password                string            `yaml:"password,omitempty"`
	GwKeyValueRegionName    string            `yaml:"gwKeyValueRegionName,omitempty"`
	GwCounterRegionName     string            `yaml:"gwCounterRegionName,omitempty"`
	GwRateLimiterRegionName string            `yaml:"gwRateLimiterRegionName,omitempty"`
	GwSortedSetRegionName   string            `yaml:"gwSortedSetRegionName,omitempty"`
	DynamicProperties       map[string]string `yaml:"dynamicProperties,omitempty"`
	Ssl                     GemfireSsl        `yaml:"ssl,omitempty"`
}

func appendRedisConfig(redisConfigBlock map[string]interface{}, new map[string]RedisClientConfig) interface{} {
	for k, v := range new {
		redisConfigBlock[k] = v
	}
	return redisConfigBlock
}

// GenerateSharedStateClientConfig renders the sharedstate_client.yaml contents for the Gateway.
// Redis and Gemfire configuration can be combined in the same document, each under its own
// top-level key, since the Gateway loads both from the same file.
func GenerateSharedStateClientConfig(configType string, redisConfigs []RedisClientConfig, hazelcastConfigs interface{}, gemfireConfig *GemfireClientConfig) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)

	defer yamlEncoder.Close()

	dynamicRedisConfig := map[string]interface{}{}
	combinedConfig := map[string]interface{}{}
	var newRedisConfigBlock interface{}
	t := new(bool)
	f := new(bool)
	*t = true
	*f = false

	if len(redisConfigs) > 0 {
		if len(redisConfigs) > 1 {
			for _, rc := range redisConfigs {
				switch rc.Name {
				case "default":
					rc.Name = ""
					if rc.Ssl.Enabled {
						if rc.Ssl.VerifyPeer == nil || !*rc.Ssl.VerifyPeer {
							rc.Ssl.VerifyPeer = f
						}
					}
					dynamicRedisConfig["default"] = rc
				default:
					name := rc.Name
					rc.Name = ""
					if rc.Ssl.Enabled {
						if rc.Ssl.VerifyPeer == nil || !*rc.Ssl.VerifyPeer {
							rc.Ssl.VerifyPeer = f
						}
					}

					newRedisConfig := map[string]RedisClientConfig{name: rc}
					newRedisConfigBlock = appendRedisConfig(dynamicRedisConfig, newRedisConfig)
				}
			}

			combinedConfig["redis"] = newRedisConfigBlock
		} else {
			redisConfigs[0].Name = ""
			if redisConfigs[0].Ssl.Enabled {
				if redisConfigs[0].Ssl.VerifyPeer == nil || !*redisConfigs[0].Ssl.VerifyPeer {
					redisConfigs[0].Ssl.VerifyPeer = f
				}
			}
			dynamicRedisConfig["default"] = redisConfigs[0]
			combinedConfig["redis"] = dynamicRedisConfig
		}
	}

	if gemfireConfig != nil {
		combinedConfig["gemfire"] = gemfireConfig
	}

	if len(combinedConfig) == 0 {
		return nil, fmt.Errorf("%s is not a supported shared state client type", configType)
	}

	if err := yamlEncoder.Encode(combinedConfig); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
