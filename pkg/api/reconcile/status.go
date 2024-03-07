package reconcile

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/caapim/layer7-operator/internal/templategen"
	"github.com/caapim/layer7-operator/pkg/api"
)

// TODO: Move status updates here.
func Status(ctx context.Context, params Params) error {
	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
		portalMeta := templategen.PortalAPI{}
		portalMetaBytes, err := json.Marshal(params.Instance.Spec.PortalMeta)
		if err != nil {
			return err
		}
		err = json.Unmarshal(portalMetaBytes, &portalMeta)
		if err != nil {
			return err
		}

		policyXml := templategen.BuildTemplate(portalMeta)
		_, sha1sum, err := api.ConvertPortalPolicyXmlToGraphman(policyXml)
		if err != nil {
			return err
		}

		if params.Instance.Status.Checksum != sha1sum {
			params.Instance.Status.Ready = true
			params.Instance.Status.Checksum = sha1sum
			err = params.Client.Status().Update(ctx, params.Instance)
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

// func Self(ctx context.Context, params Params) error {
// 	if params.Instance.Spec.PortalPublished && params.Instance.Spec.L7Portal != "" {
// 		portalMeta := templategen.PortalAPI{}
// 		portalMetaBytes, err := json.Marshal(params.Instance.Spec.PortalMeta)
// 		if err != nil {
// 			return err
// 		}
// 		err = json.Unmarshal(portalMetaBytes, &portalMeta)
// 		if err != nil {
// 			return err
// 		}

// 		policyXml := templategen.BuildTemplate(portalMeta)
// 		graphmanBundleBytes, sha1sum, err := api.ConvertPortalPolicyXmlToGraphman(policyXml)
// 		if err != nil {
// 			return err
// 		}

// 		if params.Instance.Status.Checksum != sha1sum {
// 			params.Instance.Spec.GraphmanBundle = base64.StdEncoding.EncodeToString(graphmanBundleBytes)
// 			params.Instance.Status.Ready = true
// 			params.Instance.Status.Checksum = sha1sum
// 			err = params.Client.Update(ctx, params.Instance)
// 			if err != nil {
// 				params.Log.V(2).Info("failed to update api", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
// 				return err
// 			}
// 		}
// 		return nil
// 	}

// 	if !params.Instance.Status.Ready {
// 		params.Instance.Status.Ready = true
// 		err := params.Client.Status().Update(ctx, params.Instance)
// 		if err != nil {
// 			params.Log.V(2).Info("failed to update api status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
// 			return err
// 		}
// 	}

// 	return nil
// }
