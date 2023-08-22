package reconcile

import (
	"context"
	"fmt"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"testing"
)

func TestNewSecret(t *testing.T) {
	t.Run("should create secret", func(t *testing.T) {
		//gatewayCRD := buildCRD("gateways.security.brcmlabs.com", "Gateway", "gateways", "security.brcmlabs.com")
		//repositoryCRD := buildCRD("repositories.security.brcmlabs.com", "Repository", "repositories", "security.brcmlabs.com")
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
			ErrorIfCRDPathMissing: true,
			/*CRDInstallOptions: envtest.CRDInstallOptions{
				CRDs: []*apiextensionsv1.CustomResourceDefinition{
					gatewayCRD,
					repositoryCRD,
				},
			},*/
		}
		cfg, err := testEnv.Start()
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()

		params, err := newParams()
		params.Instance.Name = "test"
		params.Instance.Namespace = "default"
		if err != nil {
			t.Fatal(err)
		}
		params.Scheme.AddKnownTypes(securityv1.GroupVersion, params.Instance)
		k8sClient, err = client.New(cfg, client.Options{Scheme: params.Scheme})
		params.Client = k8sClient
		if err != nil {
			fmt.Printf("failed to setup a Kubernetes client: %v", err)
			os.Exit(1)
		}
		err = Secret(ctx, params)
		if err != nil {
			t.Fatal(err)
		}
		/*nns := types.NamespacedName{Namespace: "default", Name: "test"}
		got := &corev1.Secret{}
		err = params.Client.Get(ctx, nns, got)
		if err != nil {
			t.Fatal(err)
		}*/
	})
}

func buildCRD(name, kind, plural, group string) *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind: kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: group,
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
			Scope: apiextensionsv1.NamespaceScoped,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural: plural,
				Kind:   kind,
			},
		},
	}
}
