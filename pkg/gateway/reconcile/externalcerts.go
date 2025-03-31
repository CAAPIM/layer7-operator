package reconcile

import (
	"context"
)

func ExternalCerts(ctx context.Context, params Params) error {
	gateway := params.Instance
	if len(gateway.Spec.App.ExternalCerts) == 0 {
		for _, v := range gateway.Status.LastAppliedExternalCerts {
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
		WithBundleType(BundleTypeExternalCert),
	)

	if err != nil {
		return err
	}

	for _, extCert := range gwUpdReq.externalEntities {
		extCertUpdReq := gwUpdReq
		extCertUpdReq.bundle = extCert.Bundle
		extCertUpdReq.bundleName = extCert.Name
		extCertUpdReq.checksum = extCert.Checksum
		extCertUpdReq.cacheEntry = extCert.CacheEntry
		extCertUpdReq.patchAnnotation = extCert.Annotation
		err = SyncGateway(ctx, params, *extCertUpdReq)
		if err != nil {
			return err
		}
	}
	return nil
}
