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
//
// The bundle's trusted certificates are pointers, so the entries are changed in place and a
// caller that still holds the same bundle sees the filled-in thumbprint too. That is harmless (it
// only adds the value Graphman would report), but it is not a copy.
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
