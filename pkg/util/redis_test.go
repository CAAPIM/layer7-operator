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
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/caapim/layer7-operator/api/v1alpha1"
)

// newTestCert generates a self-signed certificate for use in pinning tests.
func newTestCert(t *testing.T, cn string) *x509.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("serial: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{cn},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	crt, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}
	return crt
}

func certPEM(crts ...*x509.Certificate) string {
	var b strings.Builder
	for _, c := range crts {
		b.Write(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}))
	}
	return b.String()
}

func TestVerifyPinnedLeaf(t *testing.T) {
	pinned := newTestCert(t, "redis.example")
	attacker := newTestCert(t, "attacker.example")
	verify := verifyPinnedLeaf([]*x509.Certificate{pinned})

	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{pinned}}); err != nil {
		t.Errorf("pinned leaf rejected: %v", err)
	}

	// The pinned cert is public; an attacker can append it after their own
	// leaf. Only the leaf may satisfy the pin.
	appended := tls.ConnectionState{PeerCertificates: []*x509.Certificate{attacker, pinned}}
	if err := verify(appended); err == nil {
		t.Error("attacker leaf with appended pinned cert was accepted")
	}

	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{attacker}}); err == nil {
		t.Error("unpinned leaf was accepted")
	}

	if err := verify(tls.ConnectionState{}); err == nil {
		t.Error("connection with no peer certificates was accepted")
	}
}

func TestVerifyPinnedLeafMultiplePins(t *testing.T) {
	a := newTestCert(t, "a.example")
	b := newTestCert(t, "b.example")
	verify := verifyPinnedLeaf([]*x509.Certificate{a, b})
	// Any one pinned leaf is sufficient (e.g. distinct sentinel/master certs).
	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{b}}); err != nil {
		t.Errorf("second pinned leaf rejected: %v", err)
	}
}

// newTestCA generates a CA certificate and its key.
func newTestCA(t *testing.T, cn string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: cn},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create CA: %v", err)
	}
	ca, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse CA: %v", err)
	}
	return ca, key
}

// newTestLeaf generates a server certificate issued by ca.
func newTestLeaf(t *testing.T, cn string, ca *x509.Certificate, caKey *ecdsa.PrivateKey) *x509.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{cn},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, &key.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create leaf: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse leaf: %v", err)
	}
	return leaf
}

func TestVerifyPinnedCA(t *testing.T) {
	ca, caKey := newTestCA(t, "Redis CA")
	otherCA, otherKey := newTestCA(t, "Other CA")
	leaf := newTestLeaf(t, "redis.example", ca, caKey)
	verify := verifyPinnedLeaf([]*x509.Certificate{ca})

	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{leaf}}); err != nil {
		t.Errorf("leaf issued by the pinned CA rejected: %v", err)
	}
	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{leaf, ca}}); err != nil {
		t.Errorf("leaf plus pinned CA chain rejected: %v", err)
	}

	// Appending the public CA certificate must not let another leaf through.
	attacker := newTestCert(t, "attacker.example")
	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{attacker, ca}}); err == nil {
		t.Error("self-signed leaf with appended pinned CA was accepted")
	}
	foreign := newTestLeaf(t, "redis.example", otherCA, otherKey)
	if err := verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{foreign, ca}}); err == nil {
		t.Error("leaf from another CA with appended pinned CA was accepted")
	}
}

func redisTestConfig(rType v1alpha1.RedisType, tlsCfg v1alpha1.RedisTls) *v1alpha1.Redis {
	return &v1alpha1.Redis{
		Type:       rType,
		Tls:        tlsCfg,
		Standalone: v1alpha1.RedisStandalone{Host: "127.0.0.1", Port: 6379},
		Sentinel: v1alpha1.RedisSentinel{
			Master: "mymaster",
			Nodes:  []v1alpha1.RedisSentinelNode{{Host: "127.0.0.1", Port: 26379}},
		},
	}
}

func TestRedisClientMalformedPinnedCertFailsClosed(t *testing.T) {
	cases := map[string]string{
		"no PEM marker":       "MIIBkTCB+wIJAKHBfpegPjMCMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnVu",
		"file path":           "/etc/redis/redis.crt",
		"corrupt PEM body":    "-----BEGIN CERTIFICATE-----\n!!!not base64!!!\n-----END CERTIFICATE-----\n",
		"valid PEM, bad DER":  "-----BEGIN CERTIFICATE-----\nAAECAwQF\n-----END CERTIFICATE-----\n",
		"non-certificate PEM": "-----BEGIN PUBLIC KEY-----\nAAECAwQF\n-----END PUBLIC KEY-----\n",
	}
	for name, crt := range cases {
		for _, rType := range []v1alpha1.RedisType{v1alpha1.RedisTypeStandalone, v1alpha1.RedisTypeSentinel} {
			rdb, err := RedisClient(redisTestConfig(rType, v1alpha1.RedisTls{Enabled: true, VerifyPeer: true, RedisCrt: crt}))
			if err == nil {
				if rdb != nil {
					rdb.Close()
				}
				t.Errorf("%s/%s: expected error for malformed redisCrt, got nil", name, rType)
			}
		}
	}
}

func TestRedisClientTLSConfigSetAtConstruction(t *testing.T) {
	pinned := newTestCert(t, "redis.example")

	for _, rType := range []v1alpha1.RedisType{v1alpha1.RedisTypeStandalone, v1alpha1.RedisTypeSentinel} {
		t.Run(string(rType)+"/pinned", func(t *testing.T) {
			rdb, err := RedisClient(redisTestConfig(rType, v1alpha1.RedisTls{Enabled: true, VerifyPeer: true, RedisCrt: certPEM(pinned)}))
			if err != nil {
				t.Fatalf("RedisClient: %v", err)
			}
			defer rdb.Close()
			cfg := rdb.Options().TLSConfig
			if cfg == nil {
				t.Fatal("TLSConfig not carried into client options")
			}
			if cfg.MinVersion != tls.VersionTLS12 {
				t.Errorf("MinVersion = %x, want TLS1.2", cfg.MinVersion)
			}
			if !cfg.InsecureSkipVerify || cfg.VerifyConnection == nil {
				t.Error("pinned config must replace standard verification with VerifyConnection")
			}
			if err := cfg.VerifyConnection(tls.ConnectionState{PeerCertificates: []*x509.Certificate{pinned}}); err != nil {
				t.Errorf("pinned leaf rejected: %v", err)
			}
		})

		t.Run(string(rType)+"/verify-no-pin", func(t *testing.T) {
			rdb, err := RedisClient(redisTestConfig(rType, v1alpha1.RedisTls{Enabled: true, VerifyPeer: true}))
			if err != nil {
				t.Fatalf("RedisClient: %v", err)
			}
			defer rdb.Close()
			cfg := rdb.Options().TLSConfig
			if cfg == nil {
				t.Fatal("TLSConfig not carried into client options")
			}
			if cfg.InsecureSkipVerify || cfg.VerifyConnection != nil {
				t.Error("verifyPeer without a pin must use standard verification")
			}
		})

		t.Run(string(rType)+"/no-verify", func(t *testing.T) {
			rdb, err := RedisClient(redisTestConfig(rType, v1alpha1.RedisTls{Enabled: true}))
			if err != nil {
				t.Fatalf("RedisClient: %v", err)
			}
			defer rdb.Close()
			cfg := rdb.Options().TLSConfig
			if cfg == nil || !cfg.InsecureSkipVerify {
				t.Error("TLS without verifyPeer should be encrypted with verification disabled")
			}
		})

		t.Run(string(rType)+"/tls-disabled", func(t *testing.T) {
			rdb, err := RedisClient(redisTestConfig(rType, v1alpha1.RedisTls{}))
			if err != nil {
				t.Fatalf("RedisClient: %v", err)
			}
			defer rdb.Close()
			if rdb.Options().TLSConfig != nil {
				t.Error("TLSConfig set although TLS is disabled")
			}
		})
	}
}

// firstByteListener accepts one connection and reports the first byte the
// client sends: 0x16 is a TLS handshake record, '*' would be plaintext RESP.
func firstByteListener(t *testing.T) (host string, port int, first <-chan byte) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	ch := make(chan byte, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		b := make([]byte, 1)
		if _, err := conn.Read(b); err == nil {
			ch <- b[0]
		}
	}()
	addr := ln.Addr().(*net.TCPAddr)
	return addr.IP.String(), addr.Port, ch
}

// TestRedisClientDialsTLS asserts that the first bytes on the wire are a TLS
// handshake for both standalone and sentinel clients, i.e. the TLS config is
// honoured by the dialers go-redis builds at construction time.
func TestRedisClientDialsTLS(t *testing.T) {
	for _, rType := range []v1alpha1.RedisType{v1alpha1.RedisTypeStandalone, v1alpha1.RedisTypeSentinel} {
		t.Run(string(rType), func(t *testing.T) {
			host, port, first := firstByteListener(t)
			c := redisTestConfig(rType, v1alpha1.RedisTls{Enabled: true})
			c.Standalone = v1alpha1.RedisStandalone{Host: host, Port: port}
			c.Sentinel.Nodes = []v1alpha1.RedisSentinelNode{{Host: host, Port: port}}

			rdb, err := RedisClient(c)
			if err != nil {
				t.Fatalf("RedisClient: %v", err)
			}
			defer rdb.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = rdb.Ping(ctx).Err() // expected to fail; only the wire bytes matter

			select {
			case b := <-first:
				if b != 0x16 {
					t.Fatalf("first byte on the wire = %q, want TLS handshake (0x16); connection is plaintext", b)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("client never connected")
			}
		})
	}
}
