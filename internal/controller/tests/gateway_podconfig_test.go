package tests

import (
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
	Context("When pod config is mentioned in gateway custom resource", func() {
		var (
			gwLicenseSecretName = "gateway-license"
			namespace           = "l7operator"
			gatewayName         = "podconfig-ssg"
			version             = "11.1.2"
			image               = "docker.io/caapim/gateway:11.1.2"
		)

		BeforeEach(func() {
			DeferCleanup(func() {
				k8sClient.Delete(ctx, &securityv1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name:      gatewayName,
						Namespace: namespace,
					},
				})
			})
		})

		It("Should deploy gateway with given pod options", func() {
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
						PodAnnotations: map[string]string{
							"testAnnotation": "test1",
						},
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

				podList, err := getGatewayPods(ctx, gatewayName, namespace, k8sClient)
				if err != nil {
					return false
				}

				for _, pod := range podList.Items {
					for _, podStatus := range pod.Status.ContainerStatuses {
						if podStatus.Ready {
							if pod.ObjectMeta.Annotations["testAnnotation"] == "test1" {
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
