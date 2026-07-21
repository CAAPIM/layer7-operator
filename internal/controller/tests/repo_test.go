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
package tests

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
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
			namespace           = "l7operator"
			gatewayName         = "git-repo-ssg"
			version             = "11.2.1"
			image               = "docker.io/caapim/gateway:11.2.1"
			repoType            = securityv1.RepositoryTypeGit
		)

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
						Name:      "l7-gw-myframework",
						Namespace: namespace,
					},
				})
				k8sClient.Delete(ctx, &securityv1.Repository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "l7-gw-mysubscriptions",
						Namespace: namespace,
					},
				})
				k8sClient.Delete(ctx, &securityv1.Repository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "l7-gw-myapis",
						Namespace: namespace,
					},
				})
			})
		})

		It("Should create and apply repositories", func() {
			By("Creating framework repository")
			Expect(k8sClient.Create(ctx, &securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "l7-gw-myframework",
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Enabled:  true,
					Type:     repoType,
					Endpoint: "https://github.com/Gazza7205/l7GWMyFramework",
					Branch:   "main",
				},
			})).Should(Succeed())

			By("Creating subscription repository")
			Expect(k8sClient.Create(ctx, &securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "l7-gw-mysubscriptions",
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Enabled:  true,
					Type:     repoType,
					Endpoint: "https://github.com/Gazza7205/l7GWMySubscriptions",
					Branch:   "main",
				},
			})).Should(Succeed())

			By("Creating api repository")
			Expect(k8sClient.Create(ctx, &securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "l7-gw-myapis",
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Enabled:  true,
					Type:     repoType,
					Endpoint: "https://github.com/Gazza7205/l7GWMyAPIs",
					Branch:   "main",
				},
			})).Should(Succeed())

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
								Name:    "l7-gw-myframework",
								Enabled: true,
								Type:    "static",
							},
							{
								Name:    "l7-gw-mysubscriptions",
								Enabled: true,
								Type:    "dynamic",
							},
							{
								Name:    "l7-gw-myapis",
								Enabled: true,
								Type:    "dynamic",
							},
						},
						Management: securityv1.Management{
							Graphman: securityv1.Graphman{
								Enabled:            true,
								InitContainerImage: "docker.io/caapim/graphman-static-init:1.0.5",
							},
							Cluster: securityv1.Cluster{
								Hostname: "gateway.brcmlabs.com",
								Password: "7layer",
							},
							Username: "admin",
							Password: "7layer",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &gw)).Should(Succeed())

			By("Repositories are deployed to Gateway")
			Eventually(func() bool {
				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				httpclient := &http.Client{
					CheckRedirect: redirectPolicyFunc,
					Transport:     tr,
				}

				svc := &corev1.Service{}

				if err := k8sClient.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: namespace}, svc); err != nil {
					return false
				}
				testApiResp := TestAPIResp{}

				requestURL := fmt.Sprintf("https://%s:8443/api1", svc.Status.LoadBalancer.Ingress[0].IP)
				req, err := http.NewRequest("GET", requestURL, nil)
				if err != nil {
					return false
				}
				//req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
				req.Header.Add("client-id", "D63FA04C8447")
				resp, err := httpclient.Do(req)
				if err != nil {
					return false
				}
				err = json.NewDecoder(resp.Body).Decode(&testApiResp)
				if err != nil {
					return false
				}
				if testApiResp.Client == "D63FA04C8447" {
					return true
				}
				return false
			}).Within(time.Second * 180).WithPolling(3 * time.Second).Should(BeTrue())
		})
	})
})
