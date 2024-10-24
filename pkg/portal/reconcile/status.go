package reconcile

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/caapim/layer7-operator/internal/templategen"
)

func Status(ctx context.Context, params Params) error {

	summaryCm, err := getConfigmap(ctx, params, params.Instance.Name+"-api-summary")
	if err != nil {
		params.Log.V(2).Info("failed to retrieve summary configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return err
	}

	portalStatus := params.Instance.Status
	portalStatus.ApiSummaryConfigMap = params.Instance.Name + "-api-summary"

	portalAPIs := []templategen.PortalAPI{}

	portalSummaryGz, err := base64.StdEncoding.DecodeString(summaryCm.Data["apis"])
	if err != nil {
		params.Log.V(2).Info("failed to decode summary configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return err
	}
	tarStream := bytes.NewReader(portalSummaryGz)
	gReader, err := gzip.NewReader(tarStream)
	if err != nil {
		params.Log.V(2).Info("failed to decompress summary configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return err
	}

	portalSummaryBytes, err := io.ReadAll(gReader)
	if err != nil {
		params.Log.V(2).Info("failed to read summary configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return err
	}

	err = json.Unmarshal(portalSummaryBytes, &portalAPIs)
	if err != nil {
		params.Log.V(2).Info("failed to unmarshal summary configmap", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return err
	}

	portalStatus.ApiCount = len(portalAPIs)
	portalStatus.Checksum = summaryCm.ObjectMeta.Annotations["checksum/data"]

	if portalStatus.Checksum != params.Instance.Status.Checksum {
		portalStatus.LastUpdated = time.Now().UnixMilli()
		params.Instance.Status = portalStatus
		err := params.Client.Status().Update(ctx, params.Instance)
		if err != nil {
			params.Log.Info("failed to update portal status", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "message", err.Error())
			return err
		}
	}
	return nil
}
