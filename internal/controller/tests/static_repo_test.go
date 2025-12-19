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

package tests

import (
	"encoding/json"
	"os"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway controller", func() {
	Context("When repo of type static is updated", func() {
		var (
			gwLicenseSecretName = "gateway-license"
			repoSecretName      = "test-static-repository-secret"
			namespace           = "l7operator"
			gatewayName         = "static-repo-ssg"
			version             = "11.1.3"
			image               = "docker.io/caapim/gateway:11.1.3"
			repoName            = "local-repo"
			repoType            = securityv1.RepositoryTypeLocal
			staticChecksum      = ""
		)

		// just check static repo commit is different
		// static repo will fail on gateway because of graphman-static-init handling.
		BeforeEach(func() {
			DeferCleanup(func() {
				k8sClient.Delete(ctx, &securityv1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name:      gatewayName,
						Namespace: namespace,
					},
				})
				k8sClient.Delete(ctx, &securityv1.Repository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      repoName,
						Namespace: namespace,
					},
				})
				k8sClient.Delete(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      repoSecretName,
						Namespace: namespace,
					},
				})
			})
		})

		It("Should pick up changes in repo with gw restart", func() {

			By("Creating repo secret")

			bundle := graphman.Bundle{
				ClusterProperties: []*graphman.ClusterPropertyInput{
					{Name: "test-cwp", Value: "test-cwp-value"},
				},
			}

			bundleBytes, err := json.Marshal(bundle)

			if err != nil {
				os.Exit(1)
			}

			repoSecret := corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoSecretName,
					Namespace: namespace,
				},
				Data: map[string][]byte{"test.json": bundleBytes},
			}

			Expect(k8sClient.Create(ctx, &repoSecret)).Should(Succeed())

			By("Creating repository CRD")
			repo := securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoName,
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Enabled:        true,
					Type:           repoType,
					LocalReference: securityv1.LocalReference{SecretName: repoSecretName},
				},
			}
			Expect(k8sClient.Create(ctx, &repo)).Should(Succeed())

			By("Creating Gateway custom resource with a repository")
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
						Image:    image,
						Replicas: 1,
						Service: securityv1.Service{
							Type: corev1.ServiceTypeLoadBalancer,
							Ports: []securityv1.Ports{
								{
									Name:       "https",
									Port:       8443,
									TargetPort: 8443,
									Protocol:   "TCP",
								},
								{
									Name:       "management",
									Port:       9443,
									TargetPort: 9443,
									Protocol:   "TCP",
								},
							},
						},

						RepositoryReferences: []securityv1.RepositoryReference{
							{
								Name:    repoName,
								Enabled: true,
								Type:    "static",
							},
						},
						Management: securityv1.Management{
							Restman: securityv1.Restman{
								Enabled: true,
							},
							Cluster: securityv1.Cluster{
								Hostname: "gateway.brcmlabs.com",
								Password: "7layer",
							},
							Graphman: securityv1.Graphman{
								Enabled:            true,
								InitContainerImage: "docker.io/caapim/graphman-static-init:1.0.4",
							},
							Username: "admin",
							Password: "7layer",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &gw)).Should(Succeed())

			Eventually(func() bool {
				var repository securityv1.Repository
				repoReq := types.NamespacedName{
					Name:      repoName,
					Namespace: namespace,
				}
				if err := k8sClient.Get(ctx, repoReq, &repository); err != nil {
					return false
				}

				if !repository.Status.Ready {
					return false
				}

				dep, err := getGatewayDeployment(ctx, gatewayName, namespace, k8sClient)
				if err != nil {
					return false
				}
				staticChecksum = dep.ObjectMeta.Labels["security.brcmlabs.com/static-repositories-checksum"]

				return true
			}).Within(time.Second * 180).WithPolling(3 * time.Second).Should(BeTrue())

			By("Updating the repo secret")

			bundle = graphman.Bundle{
				ClusterProperties: []*graphman.ClusterPropertyInput{
					{Name: "test-cwp", Value: "test-cwp-value-1"},
				},
			}

			bundleBytes, err = json.Marshal(bundle)

			if err != nil {
				os.Exit(1)
			}

			repoSecret = corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoSecretName,
					Namespace: namespace,
				},
				Data: map[string][]byte{"test.json": bundleBytes},
			}

			Expect(k8sClient.Update(ctx, &repoSecret)).Should(Succeed())

			Eventually(func() bool {
				var repository securityv1.Repository
				repoReq := types.NamespacedName{
					Name:      repoName,
					Namespace: namespace,
				}
				if err := k8sClient.Get(ctx, repoReq, &repository); err != nil {
					return false
				}

				if !repository.Status.Ready {
					return false
				}

				dep, err := getGatewayDeployment(ctx, gatewayName, namespace, k8sClient)
				if err != nil {
					return false
				}

				if staticChecksum != dep.ObjectMeta.Labels["security.brcmlabs.com/static-repositories-checksum"] {
					return true
				}

				return false
			}).Within(time.Second * 180).WithPolling(3 * time.Second).Should(BeTrue())
		})
	})
})
