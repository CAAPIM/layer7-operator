package tests

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
			encSecretName       = "graphman-encryption-secret"
			namespace           = "l7operator"
			gatewayName         = "ssg"
			version             = "10.1.00_CR4"
			image               = "docker.io/caapim/gateway:10.1.00_CR4"
			repoName            = "file-repository"
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

		It("Should pick up changes in repo with gw restart", func() {
			By("Creating repository CRD")
			//Repository resource
			repo := securityv1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoName,
					Namespace: namespace,
				},
				Spec: securityv1.RepositorySpec{
					Name:     repoName,
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
			}).WithTimeout(time.Second * 120).Should(BeTrue())

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

				for _, pod := range gateway.Status.Gateway {
					if pod.Ready == false {
						return false
					}
				}
				return true

			}).WithTimeout(time.Second * 180).Should(BeTrue())

			By("Verify Gateway status")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				if gateway.Status.State == "Ready" {
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

			requestURL := fmt.Sprintf("https://%s:8443/api3", currentService.Status.LoadBalancer.Ingress[0].IP)
			req, err := http.NewRequest("GET", requestURL, nil)
			req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
			resp, err := httpclient.Do(req)
			if err != nil {
				GinkgoWriter.Printf("client: request failed: %s\n", err)
				os.Exit(1)
			}
			resBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				GinkgoWriter.Printf("client: could not read response body: %s\n", err)
				os.Exit(1)
			}

			GinkgoWriter.Printf("Response %s", resBody)
			Expect(strings.Contains(string(resBody), "hello world")).Should(BeTrue())
		})
	})
})
