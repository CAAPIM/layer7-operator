package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
}

func createGatewayLicenseSecret(secret Secret) {
	// Gateway licence
	license, found := os.LookupEnv("LICENSE")
	Expect(found).NotTo(BeFalse())

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
	Expect(secret.Client.Create(secret.Ctx, &gatewayLicense)).Should(Succeed())
}

func createRepositorySecret(secret Secret) {
	//Repository secret
	ghUname, found := os.LookupEnv("TESTREPO_USER")
	Expect(found).NotTo(BeFalse())
	ghToken, found := os.LookupEnv("TESTREPO_TOKEN")
	Expect(found).NotTo(BeFalse())
	repSecretData := make(map[string][]byte)

	repSecretData["USERNAME"] = []byte(ghUname)
	repSecretData["TOKEN"] = []byte(ghToken)

	repSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: repSecretData,
	}
	Expect(secret.Client.Create(secret.Ctx, &repSecret)).Should(Succeed())
}

func createGraphmanEncSecret(secret Secret) {
	//Graphman enc secret
	encSecretData := map[string][]byte{"FRAMEWORK_ENCRYPTION_PASSPHRASE": []byte("7layer")}
	encSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		Type: corev1.SecretTypeOpaque,
		Data: encSecretData,
	}
	Expect(secret.Client.Create(secret.Ctx, &encSecret)).Should(Succeed())
}

func createRepository(repo Repo) {
	//Repository resource
	repository := securityv1.Repository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      repo.Name,
			Namespace: repo.Namespace,
		},
		Spec: securityv1.RepositorySpec{
			Enabled:  true,
			Endpoint: repo.Url,
			Branch:   repo.Branch,
			Auth: securityv1.RepositoryAuth{
				Vendor:             "Github",
				ExistingSecretName: repo.SecretName,
			},
		},
	}
	Expect(repo.Client.Create(repo.Ctx, &repository)).Should(Succeed())
}

func commitAndPushNewFile(repo Repo) string {
	ghUname, found := os.LookupEnv("TESTREPO_USER")
	Expect(found).NotTo(BeFalse())
	ghToken, found := os.LookupEnv("TESTREPO_TOKEN")
	Expect(found).NotTo(BeFalse())
	token := string(ghToken)
	username := string(ghUname)

	gRepo, err := git.PlainOpen(repo.CheckoutPath)
	Expect(err).NotTo(HaveOccurred())

	w, err := gRepo.Worktree()
	Expect(err).NotTo(HaveOccurred())

	filename := filepath.Join(repo.CheckoutPath, "clusterProperties", "c.json")
	err = os.WriteFile(filename, []byte("{\n  \"goid\": \"84449671abe2a5b143051dbdfdf5e5f4\",\n  \"name\": \"c\",\n  \"checksum\": \"b77d1a0eca5224e5a33453b8fa5ace8fcbb1ce4e\",\n  \"description\": \"c cwp\",\n  \"hiddenProperty\": false,\n  \"value\": \"c\"\n}"), 0644)
	Expect(err).NotTo(HaveOccurred())

	_, err = w.Add("clusterProperties/c.json")
	Expect(err).NotTo(HaveOccurred())

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file

	commitHash, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.org",
			When:  time.Now(),
		},
	})
	Expect(err).NotTo(HaveOccurred())

	// Prints the current HEAD to verify that all worked well.
	obj, _ := gRepo.CommitObject(commitHash)
	GinkgoWriter.Printf("Commit hash %s", obj.Hash)

	auth := &gitHttp.BasicAuth{
		Username: username,
		Password: token,
	}
	err = gRepo.Push(&git.PushOptions{
		Auth: auth,
	})
	Expect(err).NotTo(HaveOccurred())
	return commitHash.String()
}

func commitAndPushUpdatedFile(repo Repo) string {
	ghUname, found := os.LookupEnv("TESTREPO_USER")
	Expect(found).NotTo(BeFalse())
	ghToken, found := os.LookupEnv("TESTREPO_TOKEN")
	Expect(found).NotTo(BeFalse())
	token := string(ghToken)
	username := string(ghUname)

	r, err := git.PlainOpen(repo.CheckoutPath)
	Expect(err).NotTo(HaveOccurred())

	w, err := r.Worktree()
	Expect(err).NotTo(HaveOccurred())

	filename := filepath.Join("/tmp/l7GWMyAPIs/tree/myApis/", "Rest Api 3-+api3.webapi.json")
	err = os.Remove(filename)
	err = os.WriteFile(filename, []byte("{\n  \"goid\": \"84449671abe2a5b143051dbdfdf7e684\",\n  \"name\": \"Rest Api 3\",\n  \"resolutionPath\": \"/api3\",\n  \"checksum\": \"ad069ae7b081636f7334ff76b99d09b75dd78b81\",\n  \"enabled\": true,\n  \"folderPath\": \"/myApis\",\n  \"methodsAllowed\": [\n    \"GET\",\n    \"POST\",\n    \"PUT\",\n    \"DELETE\"\n  ],\n  \"tracingEnabled\": false,\n  \"wssProcessingEnabled\": false,\n  \"policy\": {\n    \"xml\": \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n<wsp:Policy xmlns:L7p=\\\"http://www.layer7tech.com/ws/policy\\\" xmlns:wsp=\\\"http://schemas.xmlsoap.org/ws/2002/12/policy\\\">\\n    <wsp:All wsp:Usage=\\\"Required\\\">\\n        <L7p:HardcodedResponse><L7p:Base64ResponseBody stringValue=\\\"aGVsbG8gdGVzdA==\\\"/>    </L7p:HardcodedResponse>    </wsp:All>\\n</wsp:Policy>\\n\"\n  }\n}"), 0644)
	_, err = w.Add("tree/myApis/Rest Api 3-+api3.webapi.json")

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file

	commitHash, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.org",
			When:  time.Now(),
		},
	})
	Expect(err).NotTo(HaveOccurred())

	// Prints the current HEAD to verify that all worked well.
	obj, _ := r.CommitObject(commitHash)

	fmt.Println(obj)
	auth := &gitHttp.BasicAuth{
		Username: username,
		Password: token,
	}
	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	return commitHash.String()
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
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
