package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func Status(ctx context.Context, params Params) error {
	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
		traceId := params.Instance.Annotations["app.l7.traceId"]

		if params.Instance.Status.Checksum != traceId {
			params.Instance.Status.Ready = true
			params.Instance.Status.Checksum = traceId
			err := params.Client.Status().Update(ctx, params.Instance)
			if err != nil {
				params.Log.V(2).Info("failed to update api status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
				return err
			}
		}
		return nil
	}

	graphmanBundleBytes, err := base64.StdEncoding.DecodeString(params.Instance.Spec.GraphmanBundle)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(graphmanBundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	if params.Instance.Status.Checksum != sha1Sum {
		params.Instance.Status.Ready = true
		params.Instance.Status.Checksum = sha1Sum
		err := params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			return err
		}
	}
	return nil
}
