package reconcile

import (
	"context"
	"fmt"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	flag := false
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    &flag,
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
		Instance: &securityv1.Repository{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Repository",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
				UID:       instanceUID,
			},
			Spec: securityv1.RepositorySpec{
				Enabled:  true,
				Name:     "test",
				Branch:   "testbranch",
				Endpoint: "testing.com",
				Auth: securityv1.RepositoryAuth{
					Username: "testUser",
					Password: "testPassword",
					Token:    "testToken",
					Type:     securityv1.RepositoryAuthTypeBasic,
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
