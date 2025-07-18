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
package reconcile

import (
	"context"
)

func ExternalSecrets(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalSecrets) == 0 {
		for _, v := range gateway.Status.LastAppliedExternalSecrets {
			if len(v) != 0 {
				continue
			}
			return nil
		}
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeExternalSecret),
	)

	if err != nil {
		return err
	}

	for _, extSecret := range gwUpdReq.externalEntities {
		extSecretUpdReq := gwUpdReq
		extSecretUpdReq.bundle = extSecret.Bundle
		extSecretUpdReq.bundleName = extSecret.Name
		extSecretUpdReq.checksum = extSecret.Checksum
		extSecretUpdReq.cacheEntry = extSecret.CacheEntry
		extSecretUpdReq.patchAnnotation = extSecret.Annotation
		extSecretUpdReq.graphmanEncryptionPassphrase = extSecret.EncryptionPassphrase
		err = SyncGateway(ctx, params, *extSecretUpdReq)
		if err != nil {
			return err
		}
	}

	return nil
}
