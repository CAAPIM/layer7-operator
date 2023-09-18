package portal

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewConfigMap
func NewConfigMap(portal *securityv1alpha1.L7Portal, apiSummary []byte) *corev1.ConfigMap {
	data := make(map[string]string)
	dataCheckSum := ""

	data["apiSummary"] = base64.StdEncoding.EncodeToString(apiSummary)

	dataBytes, _ := json.Marshal(data)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	dataCheckSum = sha1Sum

	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        portal.Name + "-api-summary",
			Namespace:   portal.Namespace,
			Labels:      util.DefaultLabels(portal.Name, map[string]string{}),
			Annotations: map[string]string{"checksum/data": dataCheckSum},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Data: data,
	}
	return cmap
}
