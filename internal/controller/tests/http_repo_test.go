package tests

import (
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
			gatewayName         = "http-repo-ssg"
			version             = "11.1.2"
			image               = "docker.io/caapim/gateway:11.1.2"
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

		It("Should apply http type repo", func() {
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

			By("Creating Gateway custom resource with a repository")
			Expect(createGatewayLicenseSecret(Secret{k8sClient, ctx, gwLicenseSecretName, namespace})).Should(Succeed())
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
							},
						},
						Management: securityv1.Management{
							Restman: securityv1.Restman{
								Enabled: true,
							},
							Graphman: securityv1.Graphman{
								Enabled:            true,
								InitContainerImage: "docker.io/caapim/graphman-static-init:1.0.2",
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

			By("Verify Gateway Pod has been updated")
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

				podList, err := getGatewayPods(ctx, gatewayName, namespace, k8sClient)
				if err != nil {
					return false
				}

				for _, pod := range podList.Items {
					for _, podStatus := range pod.Status.ContainerStatuses {
						if podStatus.Ready {
							if pod.ObjectMeta.Annotations["security.brcmlabs.com/http-repo-dynamic"] == repository.Status.Commit {
								return true
							}
						}
					}
				}
				return false
			}).Within(time.Second * 180).WithPolling(3 * time.Second).Should(BeTrue())
		})
	})
})
