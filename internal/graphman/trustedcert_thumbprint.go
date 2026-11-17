package graphman

import (
	"crypto/sha1"
	"encoding/base64"
)

// normalizeTrustedCertThumbprints fills in ThumbprintSha1 for trusted certificates that arrive
// without one. The value is the base64 encoding of the SHA-1 of the DER certificate, which is
// the format Graphman itself reports.
//
// Why: bundle merging identifies a trusted certificate only by ThumbprintSha1 and skips any
// entity whose identity is empty. The Portal enrolment bundle carries its trusted certificate
// (cn=tssg) without a thumbprint, so the certificate was silently dropped on merge and the
// gateway never trusted the Portal.
//
// Entries that already have a thumbprint, or whose certBase64 cannot be decoded, are left
// unchanged.
func normalizeTrustedCertThumbprints(b *Bundle) {
	if b == nil {
		return
	}
	for _, c := range b.TrustedCerts {
		if c == nil || c.ThumbprintSha1 != "" || c.CertBase64 == "" {
			continue
		}
		der, err := base64.StdEncoding.DecodeString(c.CertBase64)
		if err != nil || len(der) == 0 {
			continue
		}
		sum := sha1.Sum(der)
		c.ThumbprintSha1 = base64.StdEncoding.EncodeToString(sum[:])
	}
}
