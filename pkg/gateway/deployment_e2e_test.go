package gateway

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"testing"
	"time"
)

func TestDeploymentE2EWithPorts(t *testing.T) {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "../../testdata/kubeconfig", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}
	//create repository resource
	err = createRepositoryResource(client)

	//create gateway resource
	err = createGatewayResource(client)
	checkIfError(err)
	time.Sleep(2 * time.Minute)

	currentDeployment := &appsv1.Deployment{}
	err = client.Get(context.Background(), types.NamespacedName{Name: "demo", Namespace: "l7operator"}, currentDeployment)
	if err != nil {
		panic(err)
	}
	if currentDeployment.Status.ReadyReplicas != 1 {
		t.Errorf("expected %d, actual %d", 1, currentDeployment.Status.ReadyReplicas)
	}

	currentService := &corev1.Service{}
	err = client.Get(context.Background(), types.NamespacedName{Name: "demo", Namespace: "l7operator"}, currentService)
	checkIfError(err)
	t.Log("IP: " + currentService.Status.LoadBalancer.Ingress[0].IP)

	err = commitAndPushNewFile(client)
	checkIfError(err)

	time.Sleep(1 * time.Minute)
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
	checkIfError(err)
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
	if !strings.Contains(string(resBody), "hello test") {
		t.Errorf("expected hello test response from api3")
	}
	t.Cleanup(func() {
		cleanupRepo(client)
		//deleteGateway(client)
		//deleteRepository(client)
	})

}

/*func deleteGateway(client client.Client) {
	gateway := &securityv1.Gateway{}
	err := client.Get(context.Background(), types.NamespacedName{Name: "demo", Namespace: "l7operator"}, gateway)
	checkIfError(err)
	client.Delete(context.Background(), gateway)
}
func deleteRepository(client client.Client) {
	repository := &securityv1.Repository{}
	err := client.Get(context.Background(), types.NamespacedName{Name: "l7-gw-myframework", Namespace: "l7operator"}, repository)
	checkIfError(err)
	client.Delete(context.Background(), repository)
}*/

func cleanupRepo(client client.Client) error {
	repositorySecret := &corev1.Secret{}

	err := client.Get(context.Background(), types.NamespacedName{Name: "test-repository-secret", Namespace: "l7operator"}, repositorySecret)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			panic(err)
		}
	}
	token := string(repositorySecret.Data["TOKEN"])
	if token == "" {
		token = string(repositorySecret.Data["PASSWORD"])
	}

	username := string(repositorySecret.Data["USERNAME"])

	r, err := git.PlainOpen("/tmp/l7-gw-myapis-main")

	w, err := r.Worktree()
	filename := filepath.Join("/tmp/l7-gw-myapis-main/tree/myApis", "Rest Api 3-+api3.webapi.json")
	err = os.Remove(filename)
	err = os.WriteFile(filename, []byte("{\n  \"goid\": \"84449671abe2a5b143051dbdfdf7e684\",\n  \"name\": \"Rest Api 3\",\n  \"resolutionPath\": \"/api3\",\n  \"checksum\": \"ad069ae7b081636f7334ff76b99d09b75dd78b81\",\n  \"enabled\": true,\n  \"folderPath\": \"/myApis\",\n  \"methodsAllowed\": [\n    \"GET\",\n    \"POST\",\n    \"PUT\",\n    \"DELETE\"\n  ],\n  \"tracingEnabled\": false,\n  \"wssProcessingEnabled\": false,\n  \"policy\": {\n    \"xml\": \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n<wsp:Policy xmlns:L7p=\\\"http://www.layer7tech.com/ws/policy\\\" xmlns:wsp=\\\"http://schemas.xmlsoap.org/ws/2002/12/policy\\\">\\n    <wsp:All wsp:Usage=\\\"Required\\\">\\n    <L7p:HardcodedResponse><L7p:Base64ResponseBody stringValue=\\\"aGVsbG8gd29ybGQ=\\\"/>    </L7p:HardcodedResponse>    </wsp:All>\\n</wsp:Policy>\\n\"\n  }\n}"), 0644)
	_, err = w.Add("tree/myApis/Rest Api 3-+api3.webapi.json")

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	checkIfError(err)
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file

	commitHash, err := w.Commit("clean up the file created", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.org",
			When:  time.Now(),
		},
	})
	checkIfError(err)

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
	return err
}

func checkIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("admin", "7layer"))
	return nil
}

func createRepositoryResource(client client.Client) error {
	repo := &unstructured.Unstructured{}
	repo.Object = map[string]interface{}{
		"apiVersion": "security.brcmlabs.com/v1",
		"kind":       "Repository",
		"metadata": map[string]interface{}{
			"name": "l7-gw-myapis",
		},
		"spec": map[string]interface{}{
			"name":     "l7-gw-myapis",
			"enabled":  true,
			"endpoint": "https://github.com/uppoju/l7GWMyAPIs",
			"branch":   "main",
			"auth": map[string]interface{}{
				"vendor":             "Github",
				"existingSecretName": "test-repository-secret",
			},
		},
	}
	repo.SetNamespace("l7operator")
	err := client.Create(context.Background(), repo)
	return err
}

func createGatewayResource(client client.Client) error {
	u := &unstructured.Unstructured{}
	u.Object = map[string]interface{}{
		"apiVersion": "security.brcmlabs.com/v1",
		"kind":       "Gateway",
		"metadata": map[string]interface{}{
			"name": "demo",
		},
		"spec": map[string]interface{}{
			"license": map[string]interface{}{
				"accept":     true,
				"secretName": "gateway-license",
			},
			"app": map[string]interface{}{
				"image":    "docker.io/caapim/gateway:10.1.00_CR4",
				"replicas": 1,
				"service": map[string]interface{}{
					"type": "LoadBalancer",
					"ports": []interface{}{
						map[string]interface{}{
							"name":       "https",
							"port":       8443,
							"targetPort": 8443,
							"protocol":   "TCP",
						},
						map[string]interface{}{
							"name":       "management",
							"port":       9443,
							"targetPort": 9443,
							"protocol":   "TCP",
						},
					},
				},
				"repositoryReferences": []interface{}{
					map[string]interface{}{
						"name":    "l7-gw-myapis",
						"enabled": true,
						"type":    "dynamic",
						"encryption": map[string]interface{}{
							"existingSecret": "graphman-encryption-secret",
							"key":            "FRAMEWORK_ENCRYPTION_PASSPHRASE",
						},
					},
				},
				"management": map[string]interface{}{
					"restman": map[string]interface{}{
						"enabled": true,
					},
					"graphman": map[string]interface{}{
						"enabled":            true,
						"initContainerImage": "docker.io/layer7api/graphman-static-init:1.0.1",
					},
					"username": "admin",
					"password": "7layer",
				},
			},
		},
	}
	u.SetNamespace("l7operator")
	err := client.Create(context.Background(), u)
	return err
}

func commitAndPushNewFile(client client.Client) error {
	repositorySecret := &corev1.Secret{}

	err := client.Get(context.Background(), types.NamespacedName{Name: "test-repository-secret", Namespace: "l7operator"}, repositorySecret)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			panic(err)
		}
	}
	token := string(repositorySecret.Data["TOKEN"])
	if token == "" {
		token = string(repositorySecret.Data["PASSWORD"])
	}

	username := string(repositorySecret.Data["USERNAME"])
	sshKey := repositorySecret.Data["SSH_KEY"]
	sshKeyPass := string(repositorySecret.Data["SSH_KEY_PASS"])
	knownHosts := repositorySecret.Data["KNOWN_HOSTS"]
	var commit string
	commit, err = util.CloneRepository("https://github.com/uppoju/l7GWMyAPIs", username, token, sshKey, sshKeyPass, "main", "", "", "l7-gw-myapis", "Github", string(securityv1.RepositoryAuthTypeBasic), knownHosts)
	if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
		fmt.Print(err.Error())
	}
	fmt.Print("commit version" + commit)

	r, err := git.PlainOpen("/tmp/l7-gw-myapis-main")

	w, err := r.Worktree()
	filename := filepath.Join("/tmp/l7-gw-myapis-main/tree/myApis/", "Rest Api 3-+api3.webapi.json")
	err = os.Remove(filename)
	err = os.WriteFile(filename, []byte("{\n  \"goid\": \"84449671abe2a5b143051dbdfdf7e684\",\n  \"name\": \"Rest Api 3\",\n  \"resolutionPath\": \"/api3\",\n  \"checksum\": \"ad069ae7b081636f7334ff76b99d09b75dd78b81\",\n  \"enabled\": true,\n  \"folderPath\": \"/myApis\",\n  \"methodsAllowed\": [\n    \"GET\",\n    \"POST\",\n    \"PUT\",\n    \"DELETE\"\n  ],\n  \"tracingEnabled\": false,\n  \"wssProcessingEnabled\": false,\n  \"policy\": {\n    \"xml\": \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n<wsp:Policy xmlns:L7p=\\\"http://www.layer7tech.com/ws/policy\\\" xmlns:wsp=\\\"http://schemas.xmlsoap.org/ws/2002/12/policy\\\">\\n    <wsp:All wsp:Usage=\\\"Required\\\">\\n        <L7p:HardcodedResponse><L7p:Base64ResponseBody stringValue=\\\"aGVsbG8gdGVzdA==\\\"/>    </L7p:HardcodedResponse>    </wsp:All>\\n</wsp:Policy>\\n\"\n  }\n}"), 0644)
	_, err = w.Add("tree/myApis/Rest Api 3-+api3.webapi.json")

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	checkIfError(err)
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file

	commitHash, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.org",
			When:  time.Now(),
		},
	})
	checkIfError(err)

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
	return err
}
