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
	"os"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway controller", func() {
	Context("When repo of type static is updated", func() {
		var (
			gwLicenseSecretName = "gateway-license"
			repoSecretName      = "test-repository-secret"
			encSecretName       = "graphman-encryption-secret"
			namespace           = "l7operator"
			gatewayName         = "ssg"
			version             = "11.1.1"
			image               = "docker.io/caapim/gateway:11.1.1"
			repoName            = "l7-gw-myframework"
			repoCheckoutPath    = "/tmp/l7GWMyFramework"
			repoGitUrl          = "https://github.com/uppoju/l7GWMyFramework"
			repoType            = "git"
			repo                Repo
		)

		BeforeEach(func() {
			var found bool
			branchName, found := os.LookupEnv("TEST_BRANCH")
			Expect(found).NotTo(BeFalse())
			repo = Repo{k8sClient, ctx, repoName, repoGitUrl, branchName, repoSecretName, repoCheckoutPath, namespace, repoType}
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
			})
		})

		It("Should pick up changes in repo with gw restart", func() {

			By("Creating repository CRD")
			//Repository resource
			createRepository(repo)

			var repository securityv1.Repository
			repoReq := types.NamespacedName{
				Name:      repoName,
				Namespace: namespace,
			}
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, repoReq, &repository); err != nil {
					return false
				}
				return repository.Status.Ready
			}).WithTimeout(time.Second * 180).Should(BeTrue())

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
								Encryption: securityv1.BundleEncryption{
									ExistingSecret: encSecretName,
									Key:            "FRAMEWORK_ENCRYPTION_PASSPHRASE",
								},
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
								InitContainerImage: "docker.io/caapim/graphman-static-init:1.0.2",
							},
							Username: "admin",
							Password: "7layer",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &gw)).Should(Succeed())

			var gateway securityv1.Gateway
			gwRequest := types.NamespacedName{
				Name:      gatewayName,
				Namespace: namespace,
			}

			Eventually(func() bool {
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				deployment := appsv1.Deployment{}

				if err := k8sClient.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gateway.Namespace}, &deployment); err != nil {
					return false
				}

				if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
					return false
				}

				return true

			}).WithTimeout(time.Second * 120).Should(BeTrue())

			By("Updating the repo")
			var commitHash = commitAndPushNewFile(repo)

			By("Gateway CRD should have new commit")
			GinkgoWriter.Println("Repo name %s and %s", repoName, commitHash)
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				for _, repoStatus := range gateway.Status.RepositoryStatus {
					if repoStatus.Name == repoName && repoStatus.Commit == commitHash {
						return true
					}
				}
				return false

			}).WithTimeout(time.Second * 180).Should(BeTrue())

			By("Gateway pod should restart")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				deployment := appsv1.Deployment{}

				if err := k8sClient.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gateway.Namespace}, &deployment); err != nil {
					return false
				}

				if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
					return false
				}

				return true
			}).WithTimeout(time.Second * 380).Should(BeTrue())

		})
	})
})
