package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"testing"
)

func TestNewSecret(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Spec.App.Management = securityv1.Management{}
	gateway.Spec.App.Management.Username = "testUser"
	gateway.Spec.App.Management.Password = "testPassword"
	gateway.Spec.App.Management.Cluster = securityv1.Cluster{}
	gateway.Spec.App.Management.Cluster.Password = "testClusterPassword"
	gateway.Spec.App.Management.Database = securityv1.Database{true, "jdbc:mysql:localhost:3606", "testDBUser", "testDBPassword"}

	secret := NewSecret(&gateway)

	if string(secret.Data["SSG_ADMIN_USERNAME"]) != gateway.Spec.App.Management.Username {
		t.Errorf("expected %s, actual %s", gateway.Spec.App.Management.Username, string(secret.Data["SSG_ADMIN_USERNAME"]))
	}
	if string(secret.Data["SSG_ADMIN_PASSWORD"]) != gateway.Spec.App.Management.Password {
		t.Errorf("expected %s, actual %s", gateway.Spec.App.Management.Password, string(secret.Data["SSG_ADMIN_PASSWORD"]))
	}

	if string(secret.Data["SSG_CLUSTER_PASSWORD"]) != gateway.Spec.App.Management.Cluster.Password {
		t.Errorf("expected %s, actual %s", gateway.Spec.App.Management.Cluster.Password, string(secret.Data["SSG_CLUSTER_PASSWORD"]))
	}

	if string(secret.Data["SSG_DATABASE_PASSWORD"]) != gateway.Spec.App.Management.Database.Password {
		t.Errorf("expected %s, actual %s", gateway.Spec.App.Management.Database.Password, string(secret.Data["SSG_DATABASE_PASSWORD"]))
	}

	if string(secret.Data["SSG_DATABASE_USER"]) != gateway.Spec.App.Management.Database.Username {
		t.Errorf("expected %s, actual %s", gateway.Spec.App.Management.Database.Username, string(secret.Data["SSG_DATABASE_USER"]))
	}
}
