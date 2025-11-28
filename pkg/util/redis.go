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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/redis/go-redis/v9"
)

func RedisClient(c *v1alpha1.Redis) (rdb *redis.Client, err error) {
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

	tlsConfig := tls.Config{}

	if c.Tls.Enabled {
		tlsConfig.MinVersion = tls.VersionTLS12
		localCrts := []x509.Certificate{}

		tlsConfig.InsecureSkipVerify = false

		if !c.Tls.VerifyPeer {
			tlsConfig.InsecureSkipVerify = true
		}

		if c.Tls.VerifyPeer && c.Tls.RedisCrt != "" {
			certValidations := 0
			tlsConfig.InsecureSkipVerify = true
			if strings.Contains(string(c.Tls.RedisCrt), "-----BEGIN CERTIFICATE-----") {
				crtStrings := strings.SplitAfter(string(c.Tls.RedisCrt), "-----END CERTIFICATE-----")
				crtStrings = crtStrings[:len(crtStrings)-1]
				for crt := range crtStrings {
					b, _ := pem.Decode([]byte(crtStrings[crt]))
					crtX509, err := x509.ParseCertificate(b.Bytes)
					if err != nil {
						return nil, err
					}
					localCrts = append(localCrts, *crtX509)
				}

				tlsConfig.VerifyConnection = func(cs tls.ConnectionState) error {
					// validate local certs against server certs
					for _, localCert := range localCrts {
						for _, peerCert := range cs.PeerCertificates {
							if bytes.Equal(localCert.Raw, peerCert.Raw) {
								certValidations = certValidations + 1
							}
						}
					}

					if certValidations != len(localCrts) {
						return fmt.Errorf("certificate validation failed")
					}
					return nil
				}
			}
		}
	}

	switch rType {
	case string(v1alpha1.RedisTypeStandalone):
		rdb := redis.NewClient(&redis.Options{
			Username: username,
			Addr:     c.Standalone.Host + ":" + strconv.Itoa(c.Standalone.Port),
			Password: password,
			DB:       database,
		})

		if c.Tls.Enabled {
			rdb.Options().TLSConfig = &tlsConfig
		}

		return rdb, nil

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

		if c.Tls.Enabled {
			rdb.Options().TLSConfig = &tlsConfig
		}

		return rdb, nil
	}
	return nil, err
}
