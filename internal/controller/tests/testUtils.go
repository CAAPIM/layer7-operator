package tests

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Secret struct {
	Client    client.Client
	Ctx       context.Context
	Name      string
	Namespace string
}

type Repo struct {
	Client       client.Client
	Ctx          context.Context
	Name         string
	Url          string
	Branch       string
	SecretName   string
	CheckoutPath string
	Namespace    string
	Type         securityv1.RepositoryType
}

type TestAPIResp struct {
	Client    string `json:"client"`
	Plan      string `json:"plan"`
	Service   string `json:"service"`
	ConfigVal string `json:"myDemoConfigVal"`
}

func createGatewayLicenseSecret(secret Secret) error {
	license, found := os.LookupEnv("LICENSE")
	if found {
		data := make(map[string][]byte)
		data["license.xml"] = []byte(license)
		gatewayLicense := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			Type: corev1.SecretTypeOpaque,
			Data: data,
		}
		err := secret.Client.Create(secret.Ctx, &gatewayLicense)
		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

func getGatewayPods(ctx context.Context, name string, namespace string, k8sClient client.Client) (*corev1.PodList, error) {
	podList := &corev1.PodList{}

	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(util.DefaultLabels(name, map[string]string{})),
	}
	if err := k8sClient.List(ctx, podList, listOpts...); err != nil {
		return podList, err
	}
	return podList, nil
}

func getGatewayDeployment(ctx context.Context, name string, namespace string, k8sClient client.Client) (*appsv1.Deployment, error) {
	gatewayDeployment := &appsv1.Deployment{}

	err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, gatewayDeployment)
	if err != nil {
		return nil, err
	}
	return gatewayDeployment, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
	return nil
}
