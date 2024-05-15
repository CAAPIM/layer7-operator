package reconcile

import (
	"context"
	"encoding/json"
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
		_ = removeJob(params.Instance.Name + "-sync-portal")
		return
	}

	if !l7Portal.Spec.Enabled {
		_ = removeJob(params.Instance.Name + "-sync-portal")
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
		params.Log.Info("failed to read api from temp storage directory", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
		return
	}
	err = ConfigMap(ctx, params, portalApiBytes)
	if err != nil {
		params.Log.Info("failed to reconcile configmap", "name", l7Portal.Name, "namespace", l7Portal.Namespace)
	}

}

func getConfigmap(ctx context.Context, params Params, name string) (*corev1.ConfigMap, error) {
	shortSummary := &corev1.ConfigMap{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: params.Instance.Namespace}, shortSummary)
	if err != nil {
		return shortSummary, err
	}
	return shortSummary, nil
}
