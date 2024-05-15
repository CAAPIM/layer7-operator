package tests

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gateway controller support for http repo", func() {
	Context("When repo of type static is updated", func() {
		var (
			gwLicenseSecretName = "gateway-license"
			encSecretName       = "graphman-encryption-secret"
			namespace           = "l7operator"
			gatewayName         = "ssg-repo"
			version             = "11.1.00"
			image               = "docker.io/caapim/gateway:11.1.00"
			repoName            = "http-repo"
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
						Name:      repoName,
						Namespace: namespace,
					},
				})
			})
		})

		It("Should be able to pick up changes from http type repo", func() {
			By("Creating repository CRD")
			//Repository resource
			repo := securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoName,
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Enabled:  true,
					Type:     "http",
					Endpoint: "https://raw.githubusercontent.com/uppoju/l7GWMyAPIs/main/implodedbundle.zip",
				},
			}
			Expect(k8sClient.Create(ctx, &repo)).Should(Succeed())

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
								Type:    "dynamic",
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
							Graphman: securityv1.Graphman{
								Enabled:            true,
								InitContainerImage: "docker.io/layer7api/graphman-static-init:1.0.1",
							},
							Cluster: securityv1.Cluster{
								Hostname: "gateway.brcmlabs.com",
							},
							Username: "admin",
							Password: "7layer",
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, &gw)).Should(Succeed())

			gwRequest := types.NamespacedName{
				Name:      gatewayName,
				Namespace: namespace,
			}

			By("Verify Gateway status")

			Eventually(func() bool {
				var gateway securityv1.Gateway
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				if gateway.Status.State == corev1.PodReady && gateway.Status.RepositoryStatus[0].Enabled == true {
					return true
				}
				return false
			}).WithTimeout(time.Second * 180).Should(BeTrue())

			By("Verify service deployed to Gateway")
			currentService := &corev1.Service{}
			Eventually(func() int {
				if err := k8sClient.Get(ctx, gwRequest, currentService); err != nil {
					return 0
				}
				return len(currentService.Status.LoadBalancer.Ingress)
			}).WithTimeout(time.Second * 120).Should(BeNumerically("==", 1))

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			httpclient := &http.Client{
				CheckRedirect: redirectPolicyFunc,
				Transport:     tr,
			}

			Eventually(func() bool {
				requestURL := fmt.Sprintf("https://%s:8443/restman/1.0/services/84449671abe2a5b143051dbdfdf7e684", currentService.Status.LoadBalancer.Ingress[0].IP)
				req, err := http.NewRequest("GET", requestURL, nil)
				if err != nil {
					return false
				}
				req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
				_, err = httpclient.Do(req)

				return err == nil

			}).WithTimeout(time.Second * 180).Should(BeTrue())

			requestURL := fmt.Sprintf("https://%s:8443/api3", currentService.Status.LoadBalancer.Ingress[0].IP)
			req, err := http.NewRequest("GET", requestURL, nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
			resp, err := httpclient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			resBody, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			GinkgoWriter.Printf("Response %s", resBody)

		})
	})
})
