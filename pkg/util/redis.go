/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
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
