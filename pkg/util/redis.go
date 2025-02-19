package util

import (
	"strconv"
	"strings"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/redis/go-redis/v9"
)

func RedisClient(c *v1alpha1.Redis) *redis.Client {
	username := ""
	if c.Username != "" {
		username = c.Username
	}

	password := ""
	if c.MasterPassword != "" {
		password = c.MasterPassword
	}

	database := 0
	if c.Database != 0 {
		database = c.Database
	}

	rType := strings.ToLower(string(c.Type))

	switch rType {
	case string(v1alpha1.RedisTypeStandalone):
		rdb := redis.NewClient(&redis.Options{
			Username: username,
			Addr:     c.Standalone.Host + ":" + strconv.Itoa(c.Standalone.Port),
			Password: password,
			DB:       database,
		})

		return rdb

	case string(v1alpha1.RedisTypeSentinel):
		sentinelAddrs := []string{}
		for _, sentinelAddr := range c.Sentinel.Nodes {
			sentinelAddrs = append(sentinelAddrs, sentinelAddr.Host+":"+strconv.Itoa(sentinelAddr.Port))
		}

		rdb := redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: sentinelAddrs,
			Username:      username,
			Password:      password,
			DB:            database,
		})
		return rdb
	}
	return nil
}
