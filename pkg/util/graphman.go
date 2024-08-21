package util

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	graphman "github.com/caapim/layer7-operator/internal/graphman"
)

const fipsProviderGuid = "41e5cacd15f86758f03ff2952616d4f3"

var internalPolicies = []string{"#OTK Client Context Variables", "OTK FIP Client Authentication Extension"}
var externalPolicies = []string{"#OTK OVP Configuration", "#OTK Storage Configuration", "OTK Client DB GET", ""}

type GraphmanSecret struct {
	Name        string `json:"name,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Description string `json:"description,omitempty"`

	VariableReferencable bool `json:"variableReferencable,omitempty"`
}

type GraphmanKey struct {
	Name      string `json:"name,omitempty"`
	Crt       string `json:"crt,omitempty"`
	Key       string `json:"key,omitempty"`
	Port      string `json:"port,omitempty"`
	Alias     string `json:"alias,omitempty"`
	UsageType string `json:"usageType,omitempty"`
}

type GraphmanOtkConfig struct {
	Type                     string `json:"type,omitempty"`
	InternalGatewayReference string `json:"internalGatewayReference,omitempty"`
}

func ApplyToGraphmanTarget(path string, secretBundle []byte, singleton bool, username string, password string, target string, encpass string) error {
	bundle := graphman.Bundle{}

	bundleBytes, err := BuildAndValidateBundle(path)
	if err != nil {
		return err
	}

	if bundleBytes == nil && len(secretBundle) > 0 {
		bundleBytes = secretBundle
	}

	if !singleton {
		scheduledTasks := []*graphman.ScheduledTaskInput{}
		jmsListeners := []*graphman.JmsDestinationInput{}

		err := json.Unmarshal(bundleBytes, &bundle)
		if err != nil {
			return err
		}

		for _, st := range bundle.ScheduledTasks {
			if !st.ExecuteOnSingleNode {
				scheduledTasks = append(scheduledTasks, st)
			}
		}

		bundle.ScheduledTasks = scheduledTasks

		for _, jmsl := range bundle.JmsDestinations {
			if jmsl.Direction != "OUTBOUND" {
				jmsListeners = append(jmsListeners, jmsl)
			}
		}
		bundle.JmsDestinations = jmsListeners

		bundleBytes, err = json.Marshal(bundle)
		if err != nil {
			return err
		}
	}

	_, err = graphman.ApplyDynamicBundle(username, password, "https://"+target, encpass, bundleBytes)
	if err != nil {
		return err
	}
	return nil
}

func ConvertX509ToGraphmanBundle(keys []GraphmanKey) ([]byte, error) {
	bundle := graphman.Bundle{}

	for _, key := range keys {
		//TODO: revisit this
		crtStrings := strings.SplitAfter(string(key.Crt), "-----END CERTIFICATE-----")
		crtStrings = crtStrings[:len(crtStrings)-1]
		crtsX509 := []x509.Certificate{}
		certsChain := []string{}
		for crt := range crtStrings {
			b, _ := pem.Decode([]byte(crtStrings[crt]))
			crtX509, _ := x509.ParseCertificate(b.Bytes)
			crtsX509 = append(crtsX509, *crtX509)
			certsChain = append(certsChain, crtStrings[crt])
		}

		certDN := ""
		for i := range crtsX509 {
			if i == 0 {
				certDN = crtsX509[i].Subject.CommonName
			}
		}

		gmanKey := graphman.KeyInput{
			KeystoreId: "00000000000000000000000000000002",
			Pem:        key.Key,
			Alias:      key.Name,
			KeyType:    "RSA",
			SubjectDn:  "CN=" + certDN,
			CertChain:  certsChain,
		}

		if key.Alias != "" {
			gmanKey.Alias = key.Alias
		}

		switch strings.ToUpper(key.UsageType) {
		case "SSL", "CA", "AUDIT_SIGNING", "AUDIT_VIEWER":
			gmanKey.UsageTypes = []graphman.KeyUsageType{graphman.KeyUsageType(key.UsageType)}
		}

		bundle.Keys = append(bundle.Keys, &gmanKey)

	}

	bundleBytes, _ := json.Marshal(bundle)
	return bundleBytes, nil
}

func ConvertOpaqueMapToGraphmanBundle(secrets []GraphmanSecret) ([]byte, error) {
	bundle := graphman.Bundle{}
	for _, secret := range secrets {
		description := "layer7 operator managed secret"
		if secret.Description != "" {
			description = secret.Description
		}

		variableReferencable := false
		if secret.VariableReferencable {
			variableReferencable = secret.VariableReferencable
		}

		// basic check to determine if secret is a private key
		// this doesn't cover keys that are encrypted at rest
		// additional checks will be added if there is demand.
		secretType := graphman.SecretTypePassword

		if strings.Contains(secret.Secret, "-----BEGIN") {
			secretType = graphman.SecretTypePemPrivateKey
		}

		bundle.Secrets = append(bundle.Secrets, &graphman.SecretInput{
			Name:                 secret.Name,
			SecretType:           secretType,
			Secret:               secret.Secret,
			VariableReferencable: variableReferencable,
			Description:          description,
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

func RemoveL7API(username string, password string, target string, apiName string, policyFragmentName string) error {
	_, err := graphman.RemoveL7PortalApi(username, password, "https://"+target, apiName, policyFragmentName)
	if err != nil {
		return err
	}
	return nil
}

func CompressGraphmanBundle(path string) ([]byte, error) {

	bundleBytes, err := BuildAndValidateBundle(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(bundleBytes)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	if buf.Len() > 900000 {
		return nil, errors.New("this bundle would exceed the maximum Kubernetes secret size")

	}

	return buf.Bytes(), nil
}

func ConcatBundles(bundleMap map[string][]byte) ([]byte, error) {
	var combinedBundle []byte

	for _, bundle := range bundleMap {
		newBundle, err := graphman.ConcatBundle(combinedBundle, bundle)
		if err != nil {
			return nil, err
		}
		combinedBundle = newBundle
	}

	return combinedBundle, nil

}

func BuildAndValidateBundle(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	bundle := graphman.Bundle{}
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	bundleBytes, err := graphman.Implode(path)
	if err != nil {
		return nil, err
	}

	_ = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			segments := strings.Split(d.Name(), ".")
			ext := segments[len(segments)-1]
			if ext == "json" && !strings.Contains(strings.ToLower(d.Name()), "sourcesummary.json") && !strings.Contains(strings.ToLower(d.Name()), "bundle-properties.json") {
				//sbb := bundleBytes
				srcBundleBytes, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				sbb, err := graphman.ConcatBundle(srcBundleBytes, bundleBytes)
				if err != nil {
					return nil
				}
				bundleBytes = sbb
			}
		}
		return nil
	})

	// if the bundle is still empty after parsing all of the directory files
	// return an error
	if len(bundleBytes) <= 2 {
		return nil, errors.New("no valid graphman bundles were found")
	}

	r := bytes.NewReader(bundleBytes)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	_ = json.Unmarshal(bundleBytes, &bundle)
	// check the graphman bundle for errors
	err = d.Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return bundleBytes, nil
}

func BuildOtkOverrideBundle(mode string, gatewayHost string, otkPort int) ([]byte, string, error) {
	var bundle graphman.Bundle
	var policyXml []byte
	switch mode {
	case "INTERNAL":
		for _, internalPolicy := range internalPolicies {
			switch internalPolicy {
			case "#OTK Client Context Variables":
				policyXml, _ = BuildLayer7PolicyXml(internalPolicy, gatewayHost, "")
				bundle.PolicyFragments = append(bundle.PolicyFragments, &graphman.PolicyFragmentInput{
					FolderPath: "/OTK/Customizations",
					Guid:       "105d3617-d61c-4c83-a952-2ed5a93425e9",
					Goid:       "bc9a31b7578652a08a514d7d4fef1fb7",
					Name:       internalPolicy,
					Policy:     &graphman.PolicyInput{Xml: string(policyXml)},
					Soap:       false,
				})
			case "OTK FIP Client Authentication Extension":
				policyXml, _ = BuildLayer7PolicyXml(internalPolicy, "", fipsProviderGuid)
				bundle.PolicyFragments = append(bundle.PolicyFragments, &graphman.PolicyFragmentInput{
					FolderPath: "/OTK/Customizations/authentication",
					Guid:       "7847c7a6-ac68-456b-841a-122726323efd",
					Goid:       "bc9a31b7578652a08a514d7d4fef30e1",
					Name:       internalPolicy,
					Policy:     &graphman.PolicyInput{Xml: string(policyXml)},
					Soap:       false,
				})
			}
		}

		bundle.Fips = append(bundle.Fips, &graphman.FipInput{
			Name:                     "otk-fips-provider",
			Goid:                     fipsProviderGuid,
			EnableCredentialTypeSaml: false,
			EnableCredentialTypeX509: true,
			CertificateValidation:    graphman.CertificateValidationTypeCertificateOnly,
			CertificateReferences:    []*graphman.FipCertInput{},
		})
	case "DMZ":
		for _, externalPolicy := range externalPolicies {
			switch externalPolicy {
			case "#OTK OVP Configuration":
				policyXml, _ = BuildLayer7PolicyXml(externalPolicy, gatewayHost, "")
				bundle.PolicyFragments = append(bundle.PolicyFragments, &graphman.PolicyFragmentInput{
					FolderPath: "/OTK/Customizations",
					Name:       externalPolicy,
					Guid:       "a4448be1-9b0e-417f-b498-8a268cadf8a5",
					Goid:       "24e6fd7c5b6fb3a96690246c8ac492ec",
					Policy:     &graphman.PolicyInput{Xml: string(policyXml)},
					Soap:       false,
				})
			case "#OTK Storage Configuration":
				policyXml, _ = BuildLayer7PolicyXml(externalPolicy, gatewayHost, "")
				bundle.PolicyFragments = append(bundle.PolicyFragments, &graphman.PolicyFragmentInput{
					FolderPath: "/OTK/Customizations",
					Name:       externalPolicy,
					Guid:       "cfa7239a-60e4-483a-9d45-c364f2fb673d",
					Goid:       "24e6fd7c5b6fb3a96690246c8ac49304",
					Policy:     &graphman.PolicyInput{Xml: string(policyXml)},
					Soap:       false,
				})
			}
		}
	case "SINGLE":
		bundle.ClusterProperties = append(bundle.ClusterProperties, &graphman.ClusterPropertyInput{
			Name:        "otk.port",
			Value:       strconv.Itoa(otkPort),
			Description: "OTK Port",
		})
	default:
		return nil, "", fmt.Errorf("invalid otk installation type %s. Valid types are single, dmz and internal", mode)
	}

	bundleBytes := new(bytes.Buffer)
	enc := json.NewEncoder(bundleBytes)
	enc.SetEscapeHTML(false)
	enc.Encode(&bundle)

	//bundleBytes, _ := json.Marshal(bundle)
	h := sha1.New()
	h.Write(bundleBytes.Bytes())
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	bundleCheckSum := sha1Sum

	return bundleBytes.Bytes(), bundleCheckSum, nil
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

// Legacy PEM to P12 conversion
// func ConvertX509ToGraphmanBundle(keys []GraphmanKey) ([]byte, error) {
// 	bundle := graphman.Bundle{}

// 	for _, key := range keys {

// 		crtStrings := strings.SplitAfter(string(key.Crt), "-----END CERTIFICATE-----")
// 		crtStrings = crtStrings[:len(crtStrings)-1]

// 		// flip the chain order for pfx.
// 		for i, j := 0, len(crtStrings)-1; i < j; i, j = i+1, j-1 {
// 			crtStrings[i], crtStrings[j] = crtStrings[j], crtStrings[i]
// 		}
// 		crtsX509 := []x509.Certificate{}
// 		certsChain := []string{}
// 		for crt := range crtStrings {
// 			b, _ := pem.Decode([]byte(crtStrings[crt]))
// 			crtX509, _ := x509.ParseCertificate(b.Bytes)
// 			crtsX509 = append(crtsX509, *crtX509)
// 			certsChain = append(certsChain, crtStrings[crt])
// 		}

// 		block, _ := pem.Decode([]byte(key.Key))
// 		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 		if err != nil {
// 			panic(err)
// 		}
// 		privKey := parseResult.(*rsa.PrivateKey)

// 		if err != nil {
// 			return nil, err
// 		}

// 		// Create a P12 to marshal the new p12 into
// 		p12 := pkcs12.NewWithPassword("7layer")
// 		certs := []pkcs12.CertEntry{}
// 		certDN := ""
// 		for i := range crtsX509 {
// 			certs = append(certs, pkcs12.CertEntry{Cert: &crtsX509[i], FriendlyName: key.Name})
// 			if i == 0 {
// 				certDN = crtsX509[i].Subject.CommonName
// 			}
// 		}

// 		p12.CertEntries = append(p12.CertEntries, certs...)
// 		p12.KeyEntries = append(p12.KeyEntries, pkcs12.KeyEntry{Key: privKey, FriendlyName: key.Name})
// 		p12Bytes, _ := pkcs12.Marshal(&p12)

// 		if err != nil {
// 			return nil, err
// 		}

// 		gmanKey := graphman.KeyInput{
// 			KeystoreId: "00000000000000000000000000000002",
// 			Alias:      key.Name,
// 			KeyType:    "RSA",
// 			SubjectDn:  "CN=" + certDN,
// 			P12:        base64.RawURLEncoding.EncodeToString(p12Bytes),
// 			CertChain:  certsChain,
// 		}

// 		bundle.Keys = append(bundle.Keys, &gmanKey)

// 	}
