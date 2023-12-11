package tests

import (
	"fmt"
	"os"
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
			repoSecretName      = "test-repository-secret"
			encSecretName       = "graphman-encryption-secret"
			namespace           = "l7operator"
			gatewayName         = "dynamic"
			version             = "10.1.00_CR4"
			image               = "docker.io/caapim/gateway:10.1.00_CR4"
			repoName            = "l7-gw-myapis"
			repoCheckoutPath    = "/tmp/l7GWMyAPIs"
			repoGitUrl          = "https://github.com/uppoju/l7GWMyAPIs"
			repo                Repo
		)

		BeforeEach(func() {
			var found bool
			branchName, found := os.LookupEnv("BRANCH_NAME")
			Expect(found).NotTo(BeFalse())
			repo = Repo{k8sClient, ctx, repoName, repoGitUrl, branchName, repoSecretName, repoCheckoutPath, namespace}
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

			By("Updating the repo")
			var commitHash = commitAndPushUpdatedFile(repo)

			By("Gateway CRD should have new commit")
			GinkgoWriter.Println("Repo name %s and %s", repoName, commitHash)
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, gwRequest, &gateway); err != nil {
					return false
				}

				for _, repoStatus := range gateway.Status.RepositoryStatus {
					if repoStatus.Name == repoName && repoStatus.Commit == commitHash {
						fmt.Printf("Gateway Repo status %s and %s", repoStatus.Name, repoStatus.Commit)
						return true
					}
				}
				return false

			}).WithTimeout(time.Second * 180).Should(BeTrue())

		})
	})
})
