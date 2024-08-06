package util

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type SharedStateClientConfigType string

var (
	SharedStateClientConfigTypeRedis     SharedStateClientConfigType = "redis"
	SharedStateClientConfigTypeHazelcast SharedStateClientConfigType = "hazelcast"
)

type RedisNode struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type RedisSsl struct {
	Enabled    bool   `yaml:"enabled"`
	Crt        string `yaml:"crt,omitempty"`
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

func appendRedisConfig(redisConfigBlock map[string]interface{}, new map[string]RedisClientConfig) interface{} {
	for k, v := range new {
		redisConfigBlock[k] = v
	}
	return redisConfigBlock
}

// Redis is currently the only supported shared state client that supports the new configuration.
func GenerateSharedStateClientConfig(configType string, redisConfigs []RedisClientConfig, hazelcastConfigs interface{}) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)

	defer yamlEncoder.Close()

	dynamicRedisConfig := map[string]interface{}{}
	combinedRedisConfig := map[string]interface{}{}
	var newRedisConfigBlock interface{}
	t := new(bool)
	f := new(bool)
	*t = true
	*f = false

	switch configType {
	case "redis":
		if len(redisConfigs) > 1 {
			for _, rc := range redisConfigs {
				switch rc.Name {
				case "default":
					rc.Name = ""

					if rc.Ssl.Enabled && rc.Ssl.VerifyPeer != t {
						rc.Ssl.VerifyPeer = f
					}
					dynamicRedisConfig["default"] = rc
				default:
					name := rc.Name
					rc.Name = ""
					if rc.Ssl.Enabled && rc.Ssl.VerifyPeer != t {
						rc.Ssl.VerifyPeer = f
					}
					newRedisConfig := map[string]RedisClientConfig{name: rc}
					newRedisConfigBlock = appendRedisConfig(dynamicRedisConfig, newRedisConfig)
				}
			}

			combinedRedisConfig["redis"] = newRedisConfigBlock

			err := yamlEncoder.Encode(combinedRedisConfig)
			if err != nil {
				return nil, err
			}

			return b.Bytes(), nil

		} else {
			redisConfigs[0].Name = ""
			dynamicRedisConfig["default"] = redisConfigs[0]
			combinedRedisConfig["redis"] = dynamicRedisConfig
			err := yamlEncoder.Encode(combinedRedisConfig)
			if err != nil {
				return nil, err
			}
			return b.Bytes(), nil
		}

	case "hazelcast":
		return nil, fmt.Errorf("%s is not a supported shared state client type", configType)
	default:
		return nil, fmt.Errorf("%s is not a supported shared state client type", configType)
	}

}
