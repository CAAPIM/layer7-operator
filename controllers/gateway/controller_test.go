/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gateway

import (
	"context"
	"os"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Gateway controller", func() {

	const (
		gatewayName         = "ssg"
		version             = "10.1.00_CR3"
		image               = "docker.io/caapim/gateway:10.1.00_CR3"
		gwLicenseSecretName = "gw-license-secret"
		namespace           = "default"
	)

	Context("When setting up the test environment", func() {
		It("Should create Gateway custom resources", func() {
			By("Creating a first Gateway custom resource")
			ctx := context.Background()
			license, err := os.ReadFile("example/base/resources/secrets/license/license.xml")
			Expect(err).NotTo(HaveOccurred())

			data := make(map[string][]byte)
			data["license.xml"] = license

			gatewayLicense := corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      gwLicenseSecretName,
					Namespace: namespace,
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Type: corev1.SecretTypeOpaque,
				Data: data,
			}
			Expect(k8sClient.Create(ctx, &gatewayLicense)).Should(Succeed())

			gw := securityv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      gatewayName,
					Namespace: namespace,
				},
				Spec: securityv1.GatewaySpec{
					License: securityv1.License{
						Accept:     true,
						SecretName: gwLicenseSecretName,
					},
					Version: version,
					App: securityv1.App{
						Image: image,
					},
				},
			}
			Expect(k8sClient.Create(ctx, &gw)).Should(Succeed())
		})
	})
})
