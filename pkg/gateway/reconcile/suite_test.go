package reconcile

import (
	"context"
	"fmt"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var (
	k8sClient   client.Client
	testEnv     *envtest.Environment
	testScheme  *runtime.Scheme = scheme.Scheme
	ctx         context.Context
	cancel      context.CancelFunc
	instanceUID = uuid.NewUUID()
	logger      = logf.Log.WithName("unit-tests")
)

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}
	cfg, err := testEnv.Start()
	if err != nil {
		fmt.Printf("failed to start testEnv: %v", err)
		os.Exit(1)
	}
	if err = securityv1.AddToScheme(testScheme); err != nil {
		fmt.Printf("failed to register scheme: %v", err)
		os.Exit(1)
	}
	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	if err != nil {
		fmt.Printf("failed to setup a Kubernetes client: %v", err)
		os.Exit(1)
	}
	code := m.Run()

	err = testEnv.Stop()
	if err != nil {
		fmt.Printf("failed to stop testEnv: %v", err)
		os.Exit(1)
	}

	os.Exit(code)
}
func newParams() Params {
	params := Params{
		Client: k8sClient,
		Instance: &securityv1.Gateway{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Gateway",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
				UID:       instanceUID,
			},
			Spec: securityv1.GatewaySpec{
				App: securityv1.App{
					ListenPorts: securityv1.ListenPorts{
						Custom: securityv1.CustomListenPort{
							Enabled: false,
						},
					},
					ClusterProperties: securityv1.ClusterProperties{
						Properties: []securityv1.Property{
							{Name: "cwp1", Value: "test1"},
						},
					},
					Management: securityv1.Management{
						Username: "testUsername",
						Password: "testPassword",
						Cluster: securityv1.Cluster{
							Password: "testCluster",
							Hostname: "testHost",
						},
						Database: securityv1.Database{
							Username: "databaseUser",
							Password: "databasePassword",
							JDBCUrl:  "jdbc:mysql:localhost:3606",
						},
					},
					Java: securityv1.Java{
						JVMHeap: securityv1.JVMHeap{
							Calculate:  true,
							Default:    "3g",
							Percentage: 50,
						},
					},
					Resources: securityv1.PodResources{
						Limits: corev1.ResourceList{
							"memory": resource.MustParse("4Gi"),
						},
					},
					Hazelcast: securityv1.Hazelcast{
						External: true,
						Endpoint: "hazelcast.com",
					},
					System: securityv1.System{
						Properties: "testProperty",
					},
					Service: securityv1.Service{
						Ports: []securityv1.Ports{
							{Name: "http",
								Port:       443,
								TargetPort: 8443,
								Protocol:   "TCP"},
						},
					},
					PodDisruptionBudget: securityv1.PodDisruptionBudgetSpec{
						MaxUnavailable: intstr.IntOrString{IntVal: 5},
					},
				},
				License: securityv1.License{
					Accept: true,
				},
			},
			Status: securityv1.GatewayStatus{
				RepositoryStatus: []securityv1.GatewayRepositoryStatus{
					securityv1.GatewayRepositoryStatus{
						Enabled:    true,
						Name:       "testrepo",
						Commit:     "1234",
						Type:       "static",
						SecretName: "testSecret",
						Branch:     "testBranch",
						Endpoint:   "github.com",
					},
				},
			},
		},
		Scheme: testScheme,
		Log:    logger,
	}
	params.Instance.Name = "test"
	params.Instance.Namespace = "default"
	params.Client = k8sClient
	return params
}
