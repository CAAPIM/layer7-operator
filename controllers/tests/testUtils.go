package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
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
}

func createGatewayLicenseSecret(secret Secret) {
	// Gateway licence
	license, found := os.LookupEnv("GW_LICENSE")
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
			Name:     repo.Name,
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
	repositorySecret := &corev1.Secret{}
	err := repo.Client.Get(context.Background(), types.NamespacedName{Name: repo.SecretName, Namespace: repo.Namespace}, repositorySecret)
	Expect(err).NotTo(HaveOccurred())

	token := string(repositorySecret.Data["TOKEN"])
	if token == "" {
		token = string(repositorySecret.Data["PASSWORD"])
	}

	username := string(repositorySecret.Data["USERNAME"])
	sshKey := repositorySecret.Data["SSH_KEY"]
	sshKeyPass := string(repositorySecret.Data["SSH_KEY_PASS"])
	knownHosts := repositorySecret.Data["KNOWN_HOSTS"]
	var commit string
	commit, err = util.CloneRepository(repo.Url, username, token, sshKey, sshKeyPass, repo.Branch, "", "", repo.Name, "Github", string(securityv1.RepositoryAuthTypeBasic), knownHosts)

	if err == git.NoErrAlreadyUpToDate || err == git.ErrRemoteExists {
		fmt.Print(err.Error())
	}

	GinkgoWriter.Printf("commit version %s", commit)

	r, err := git.PlainOpen(repo.CheckoutPath)
	Expect(err).NotTo(HaveOccurred())

	w, err := r.Worktree()
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
	obj, _ := r.CommitObject(commitHash)
	fmt.Println(obj)

	auth := &gitHttp.BasicAuth{
		Username: username,
		Password: token,
	}
	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	Expect(err).NotTo(HaveOccurred())
	return commitHash.String()
}

func cleanupRepo(repo Repo) {
	repositorySecret := &corev1.Secret{}

	err := repo.Client.Get(context.Background(), types.NamespacedName{Name: repo.SecretName, Namespace: repo.Namespace}, repositorySecret)
	Expect(err).NotTo(HaveOccurred())

	token := string(repositorySecret.Data["TOKEN"])
	if token == "" {
		token = string(repositorySecret.Data["PASSWORD"])
	}

	username := string(repositorySecret.Data["USERNAME"])

	r, err := git.PlainOpen(repo.CheckoutPath)

	w, err := r.Worktree()
	filename := filepath.Join(repo.CheckoutPath, "clusterProperties", "c.json")
	err = os.Remove(filename)
	_, err = w.Remove("clusterProperties/c.json")

	// We can verify the current status of the worktree using the method Status.
	status, err := w.Status()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file

	commitHash, err := w.Commit("clean up the file created", &git.CommitOptions{
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
	Expect(err).NotTo(HaveOccurred())
}
