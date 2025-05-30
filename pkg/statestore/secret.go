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

package statestore

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret
func NewSecret(statestore *securityv1alpha1.L7StateStore, data map[string][]byte, name string) *corev1.Secret {
	dataBytes, _ := json.Marshal(data)
	h := sha1.New()
	h.Write(dataBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	//name := statestore.Name + "-secret"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   statestore.Namespace,
			Labels:      util.DefaultLabels(name, nil),
			Annotations: map[string]string{"checksum/data": sha1Sum},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}

	return secret
}
