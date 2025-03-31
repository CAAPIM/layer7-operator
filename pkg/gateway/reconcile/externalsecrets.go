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
