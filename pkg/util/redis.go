/*
* Copyright (c) 2026 Broadcom. All rights reserved.
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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
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

// RedisClient builds a short-lived go-redis client. Callers must Close() when done to release the connection pool.
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

	var tlsConfig *tls.Config

	if c.Tls.Enabled {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}

		if !c.Tls.VerifyPeer {
			tlsConfig.InsecureSkipVerify = true
		}

		if c.Tls.VerifyPeer && c.Tls.RedisCrt != "" {
			// Parse the pin set before touching verification settings so that a
			// configured-but-unusable redisCrt fails closed instead of leaving
			// the client with verification disabled and no substitute check.
			pinnedCrts, err := parsePinnedCertificates(c.Tls.RedisCrt)
			if err != nil {
				return nil, err
			}
			// Standard chain/hostname verification is replaced by leaf pinning;
			// InsecureSkipVerify is only ever set together with VerifyConnection.
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.VerifyConnection = verifyPinnedLeaf(pinnedCrts)
		}
	}

	// TLSConfig must be set in the options literal: go-redis captures it when the
	// client (and, for sentinel, its master/replica and sentinel dialers) is
	// constructed, so assigning rdb.Options().TLSConfig afterwards is not
	// reliably honoured and can leave connections in plaintext.
	switch rType {
	case string(v1alpha1.RedisTypeStandalone):
		rdb := redis.NewClient(&redis.Options{
			Username:  username,
			Addr:      c.Standalone.Host + ":" + strconv.Itoa(c.Standalone.Port),
			Password:  password,
			DB:        database,
			TLSConfig: tlsConfig,
		})

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
			TLSConfig:     tlsConfig,
		})

		return rdb, nil
	}
	return nil, err
}

// parsePinnedCertificates decodes every PEM CERTIFICATE block in redisCrt. It
// returns an error if a certificate block fails to parse or if no certificate
// is found at all, so pinning can never silently degrade to "accept any peer".
func parsePinnedCertificates(redisCrt string) ([]*x509.Certificate, error) {
	pinned := []*x509.Certificate{}
	rest := []byte(redisCrt)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		crt, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("redis tls.redisCrt contains an invalid certificate: %w", err)
		}
		pinned = append(pinned, crt)
	}
	if len(pinned) == 0 {
		return nil, fmt.Errorf("redis tls.verifyPeer is enabled but tls.redisCrt contains no parseable PEM certificate")
	}
	return pinned, nil
}

// verifyPinnedLeaf returns a tls.Config.VerifyConnection callback that accepts
// the connection only if the server's leaf certificate (PeerCertificates[0],
// the one whose private key the handshake proved possession of) is one of the
// pinned certificates or chains to one of them. A pinned certificate merely
// appearing elsewhere in the peer-supplied chain is not enough.
func verifyPinnedLeaf(pinned []*x509.Certificate) func(tls.ConnectionState) error {
	roots := x509.NewCertPool()
	for _, p := range pinned {
		roots.AddCert(p)
	}
	return func(cs tls.ConnectionState) error {
		if len(cs.PeerCertificates) == 0 {
			return fmt.Errorf("certificate validation failed: redis server presented no certificate")
		}
		leaf := cs.PeerCertificates[0]
		for _, p := range pinned {
			if bytes.Equal(p.Raw, leaf.Raw) {
				return nil
			}
		}
		// redisCrt may also hold the CA that issued the server certificate. The leaf
		// must chain to a pinned certificate; extra certificates the server sends are
		// only used as intermediates. Hostnames aren't checked, as before.
		intermediates := x509.NewCertPool()
		for _, c := range cs.PeerCertificates[1:] {
			intermediates.AddCert(c)
		}
		if _, err := leaf.Verify(x509.VerifyOptions{Roots: roots, Intermediates: intermediates, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny}}); err == nil {
			return nil
		}
		return fmt.Errorf("certificate validation failed: redis server certificate does not match or chain to a pinned certificate")
	}
}
