package secrets

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"testing"
)

func TestNewSecret(t *testing.T) {
	repository := getRepositoryWithAuth()
	repository.Spec.Auth.Username = "testUser"
	repository.Spec.Auth.Password = "testPassword"
	repository.Spec.Auth.Token = "testToken"

	secret := NewSecret(&repository)

	if string(secret.Data["USERNAME"]) != repository.Spec.Auth.Username {
		t.Errorf("expected %s, actual %s", repository.Spec.Auth.Username, string(secret.Data["SSG_ADMIN_USERNAME"]))
	}
	if string(secret.Data["PASSWORD"]) != repository.Spec.Auth.Password {
		t.Errorf("expected %s, actual %s", repository.Spec.Auth.Password, string(secret.Data["SSG_ADMIN_PASSWORD"]))
	}

	if string(secret.Data["TOKEN"]) != repository.Spec.Auth.Token {
		t.Errorf("expected %s, actual %s", repository.Spec.Auth.Token, string(secret.Data["SSG_CLUSTER_PASSWORD"]))
	}

}

func getRepositoryWithAuth() securityv1.Repository {
	repository := securityv1.Repository{}
	repository.Spec = securityv1.RepositorySpec{}
	repository.Spec.Auth = securityv1.RepositoryAuth{}
	return repository
}