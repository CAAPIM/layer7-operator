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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package reconcile

import (
	"context"
)

func OtkFipsCerts(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.Otk.FipsCertificates) == 0 {
		if len(gateway.Status.LastAppliedOtkFipsCerts) == 0 {
			return nil
		}
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeOTKFips),
	)

	if err != nil {
		return err
	}

	if gwUpdReq == nil {
		return nil
	}

	for _, fipsCert := range gwUpdReq.externalEntities {
		fipsCertUpdReq := gwUpdReq
		fipsCertUpdReq.bundle = fipsCert.Bundle
		fipsCertUpdReq.bundleName = fipsCert.Name
		fipsCertUpdReq.checksum = fipsCert.Checksum
		fipsCertUpdReq.cacheEntry = fipsCert.CacheEntry
		fipsCertUpdReq.patchAnnotation = fipsCert.Annotation
		err = SyncGateway(ctx, params, *fipsCertUpdReq)
		if err != nil {
			return err
		}
	}
	return nil
}
