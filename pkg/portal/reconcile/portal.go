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
package reconcile

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/internal/templategen"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

const tempDirectoryBase = "/tmp/portalapis/"

type RawPortalAPISummary struct {
	APIs []templategen.PortalAPI `json:"results"`
}

func syncPortal(ctx context.Context, params Params) {

	l7Portal := &v1alpha1.L7Portal{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, l7Portal)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Info("portal not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Namespace + "/" + params.Instance.Name + "-sync-portal")
		return
	}

	if !l7Portal.Spec.Enabled {
		_ = removeJob(params.Instance.Namespace + "/" + params.Instance.Name + "-sync-portal")
		return
	}

	portalAPIs := []templategen.PortalAPI{}
	portalTempDirectory := tempDirectoryBase + l7Portal.Name

	folderInfo, err := os.Stat(portalTempDirectory)
	if err != nil {
		params.Log.Info("failed to scan temp storage", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	_, err = getConfigmap(ctx, params, l7Portal.Name+"-api-summary")

	if folderInfo.ModTime().Add(30*time.Second).Before(time.Now()) && err == nil {
		return
	}

	dInfo, err := os.ReadDir(portalTempDirectory)
	if err != nil {
		params.Log.V(2).Info("failed to read temp storage directory", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	for _, f := range dInfo {
		fBytes, err := os.ReadFile(portalTempDirectory + "/" + f.Name())
		if err != nil {
			params.Log.V(2).Info("failed to read portal api from temp storage", "name", l7Portal.Name, "summary file", f.Name(), "namespace", l7Portal.Namespace)
		}

		portalAPI := templategen.PortalAPI{}
		err = json.Unmarshal(fBytes, &portalAPI)
		if err != nil {
			params.Log.Info("failed to unmarshal portal api summary", "name", l7Portal.Name, "summary file", f.Name(), "namespace", l7Portal.Namespace)
			continue
		}
		portalAPIs = append(portalAPIs, portalAPI)
	}

	portalApiBytes, err := json.Marshal(portalAPIs)
	if err != nil {
		params.Log.Info("failed to marshal api summary", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(portalApiBytes)
	if err != nil {
		params.Log.Error(err, "failed to compress api summary")
		return
	}

	if err := zw.Close(); err != nil {
		params.Log.Error(err, "failed to compress api summary")
		return
	}
	if buf.Len() > 900000 {
		params.Log.Error(errors.New("this bundle would exceed the maximum Kubernetes secret size"), "failed to compress api summary")
		return
	}

	err = ConfigMap(ctx, params, buf.Bytes())
	if err != nil {
		params.Log.Info("failed to reconcile configmap", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}
	buf.Reset()

}

func getConfigmap(ctx context.Context, params Params, name string) (*corev1.ConfigMap, error) {
	shortSummary := &corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, shortSummary)
	if err != nil {
		return shortSummary, err
	}
	return shortSummary, nil
}
