package util

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strings"

	graphman "github.com/caapim/layer7-operator/internal/graphman"
	"github.com/gazza7205/go-pkcs12"
)

// implement back off loop

type GraphmanSecret struct {
	Name   string `json:"name,omitempty"`
	Secret string `json:"secret,omitempty"`
}

type GraphmanKey struct {
	Name string `json:"name,omitempty"`
	Crt  string `json:"crt,omitempty"`
	Key  string `json:"key,omitempty"`
	Port string `json:"port,omitempty"`
}

func ApplyToGraphmanTarget(path string, username string, password string, target string, encpass string) error {
	_, err := graphman.Apply(path, username, password, "https://"+target, encpass)
	if err != nil {
		return err
	}
	return nil
}

func ConvertX509ToGraphmanBundle(keys []GraphmanKey) ([]byte, error) {
	bundle := graphman.Bundle{}

	for _, key := range keys {

		crtStrings := strings.SplitAfter(string(key.Crt), "-----END CERTIFICATE-----")
		crtStrings = crtStrings[:len(crtStrings)-1]

		// flip the chain order for pfx.
		for i, j := 0, len(crtStrings)-1; i < j; i, j = i+1, j-1 {
			crtStrings[i], crtStrings[j] = crtStrings[j], crtStrings[i]
		}
		crtsX509 := []x509.Certificate{}
		certsChain := []string{}
		for crt := range crtStrings {
			b, _ := pem.Decode([]byte(crtStrings[crt]))
			crtX509, _ := x509.ParseCertificate(b.Bytes)
			crtsX509 = append(crtsX509, *crtX509)
			certsChain = append(certsChain, crtStrings[crt])
		}

		block, _ := pem.Decode([]byte(key.Key))
		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		privKey := parseResult.(*rsa.PrivateKey)

		if err != nil {
			return nil, err
		}

		// Create a P12 to marshal the new p12 into
		p12 := pkcs12.NewWithPassword("7layer")
		certs := []pkcs12.CertEntry{}
		certDN := ""
		for i := range crtsX509 {
			certs = append(certs, pkcs12.CertEntry{Cert: &crtsX509[i], FriendlyName: key.Name})
			if i == 0 {
				certDN = crtsX509[i].Subject.CommonName
			}
		}

		p12.CertEntries = append(p12.CertEntries, certs...)
		p12.KeyEntries = append(p12.KeyEntries, pkcs12.KeyEntry{Key: privKey, FriendlyName: key.Name})
		p12Bytes, _ := pkcs12.Marshal(&p12)

		if err != nil {
			return nil, err
		}

		gmanKey := graphman.KeyInput{
			KeystoreId: "00000000000000000000000000000002",
			Alias:      key.Name,
			KeyType:    "RSA",
			SubjectDn:  "CN=" + certDN,
			P12:        base64.RawURLEncoding.EncodeToString(p12Bytes),
			CertChain:  certsChain,
		}

		bundle.Keys = append(bundle.Keys, &gmanKey)

	}

	// bundle.ListenPorts = append(bundle.ListenPorts, &graphman.ListenPortInput{
	// 	Name:            "Default HTTPS (8443)",
	// 	Port:            8443,
	// 	Enabled:         true,
	// 	Protocol:        "HTTPS",
	// 	EnabledFeatures: []graphman.ListenPortFeature{"PUBLISHED_SERVICE_MESSAGE_INPUT"},

	// 	TlsSettings: &graphman.ListenPortTlsSettingsInput{
	// 		KeystoreId:           "00000000000000000000000000000002",
	// 		KeyAlias:             "brcmlabs",
	// 		TlsVersions:          []string{"TLSv1.2", "TLSv1.3"},
	// 		UseCipherSuitesOrder: true,
	// 		ClientAuthentication: "OPTIONAL",
	// 		CipherSuites: []string{
	// 			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	// 			"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384",
	// 			"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	// 			"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
	// 			"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256",
	// 			"TLS_DHE_RSA_WITH_AES_256_CBC_SHA",
	// 			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	// 			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
	// 			"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	// 			"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	// 			"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
	// 			"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256",
	// 			"TLS_DHE_RSA_WITH_AES_128_CBC_SHA",
	// 			"TLS_AES_256_GCM_SHA384",
	// 			"TLS_AES_128_GCM_SHA256",
	// 		},
	// 	},
	// })

	bundleBytes, _ := json.Marshal(bundle)

	return bundleBytes, nil
}

func ConvertOpaqueMapToGraphmanBundle(secrets []GraphmanSecret) ([]byte, error) {
	bundle := graphman.Bundle{}
	for _, secret := range secrets {
		bundle.Secrets = append(bundle.Secrets, &graphman.SecretInput{
			Name:                 secret.Name,
			SecretType:           graphman.SecretTypePassword,
			Secret:               secret.Secret,
			VariableReferencable: true,
			Description:          "layer7 operator managed secret",
		})
	}

	bundleBytes, err := json.Marshal(bundle)

	if err != nil {
		return nil, err
	}

	return bundleBytes, nil
}

func ApplyGraphmanBundle(username string, password string, target string, encpass string, bundle []byte) error {
	_, err := graphman.ApplyDynamicBundle(username, password, "https://"+target, encpass, bundle)
	if err != nil {
		return err
	}
	return nil
}

func CompressGraphmanBundle(path string) ([]byte, error) {
	bundle, err := graphman.Implode(path)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(bundle)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	if buf.Len() > 900000 {
		return nil, errors.New("this bundle would exceed the maximum Kubernetes secret size.")

	}

	return buf.Bytes(), nil
}

// Reserved for future use.
// // Brotli can compress an 11mb restman bundle down to 550-600kb
// func CompressGraphmanBundle(path string) ([]byte, error) {
// 	bundle, err := graphman.Implode(path)

// 	if err != nil {
// 		return nil, err
// 	}

// 	bytes, err := cbrotli.Encode(bundle, cbrotli.WriterOptions{Quality: 6})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(bytes) > 900000 {
// 		return nil, errors.New("this bundle would exceed the maximum Kubernetes secret size.")
// 	}

// 	return bytes, nil
// }
