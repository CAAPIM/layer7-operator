/*
* Copyright (c) 2026 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package gateway

import (
	"reflect"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func gatewayForDeploymentProbeTest() securityv1.Gateway {
	gw := getGatewayWitApp()
	gw.Name = "test"
	gw.Namespace = "testNamespace"
	gw.Spec.App.Image = "gateway:test"
	gw.Spec.App.Replicas = 1
	gw.Spec.App.ServiceAccount.Name = "testServiceAccount"
	gw.Spec.App.ServiceAccount.Create = true
	gw.Spec.App.PodSecurityContext = corev1.PodSecurityContext{}
	gw.Spec.License.SecretName = "license-secret"
	ports := securityv1.Ports{}
	ports.Name = "https"
	ports.Port = 8443
	ports.TargetPort = 8443
	ports.Protocol = "TCP"
	gw.Spec.App.Service = securityv1.Service{}
	gw.Spec.App.Service.Ports = []securityv1.Ports{ports}
	return gw
}

func TestDeploymentStartupProbe(t *testing.T) {
	t.Run("omitted leaves StartupProbe nil", func(t *testing.T) {
		gw := gatewayForDeploymentProbeTest()
		dep := NewDeployment(&gw, "kubernetes")
		c := dep.Spec.Template.Spec.Containers[0]
		if c.StartupProbe != nil {
			t.Fatalf("StartupProbe: want nil, got %+v", c.StartupProbe)
		}
	})
	t.Run("set when spec defines StartupProbe", func(t *testing.T) {
		gw := gatewayForDeploymentProbeTest()
		want := corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/health",
					Port: intstr.FromInt32(8443),
				},
			},
			FailureThreshold: 18,
		}
		gw.Spec.App.StartupProbe = want
		dep := NewDeployment(&gw, "kubernetes")
		c := dep.Spec.Template.Spec.Containers[0]
		if c.StartupProbe == nil {
			t.Fatal("StartupProbe: want non-nil")
		}
		if !reflect.DeepEqual(*c.StartupProbe, want) {
			t.Errorf("StartupProbe mismatch\nwant: %+v\ngot:  %+v", want, *c.StartupProbe)
		}
	})
}

func TestDeploymentWithPorts(t *testing.T) {
	gateway := getGatewayWitApp()
	gateway.Name = "test"
	gateway.Namespace = "testNamespace"
	gateway.Spec.App.Image = "testImage"
	gateway.Spec.App.Replicas = 5
	gateway.Spec.App.ServiceAccount.Name = "testServiceAccount"
	gateway.Spec.App.ServiceAccount.Create = true
	gateway.Spec.App.PodSecurityContext = corev1.PodSecurityContext{}
	gateway.Spec.App.TopologySpreadConstraints = []corev1.TopologySpreadConstraint{{TopologyKey: "testTopology"}}
	gateway.Spec.App.Tolerations = []corev1.Toleration{{Key: "testKey", Value: "testValue"}}
	gateway.Spec.App.NodeSelector = map[string]string{"testLabel": "testValue"}

	ports := securityv1.Ports{}
	ports.Name = "httpPort"
	ports.Port = 8443
	ports.TargetPort = 8443
	ports.Protocol = "http"

	gateway.Spec.App.Service = securityv1.Service{}
	gateway.Spec.App.Service.Ports = []securityv1.Ports{ports}

	gateway.Spec.App.Management.Service.Enabled = true
	managementPorts := securityv1.Ports{}
	managementPorts.Name = "httpPort"
	managementPorts.Port = 9443
	managementPorts.TargetPort = 9443
	managementPorts.Protocol = "http"
	gateway.Spec.App.Management.Service.Ports = []securityv1.Ports{managementPorts}
	platform := "kubernetes"

	dep := NewDeployment(&gateway, platform)

	if dep.ObjectMeta.Namespace != "testNamespace" {
		t.Errorf("expected %s, actual %s", "testNamespace", dep.ObjectMeta.Namespace)
	}

	if dep.ObjectMeta.Name != "test" {
		t.Errorf("expected %s, actual %s", "test", dep.ObjectMeta.Name)
	}

	if *dep.Spec.Replicas != 5 {
		t.Errorf("expected %d, actual %d", 5, *dep.Spec.Replicas)
	}

	if dep.Spec.Template.Spec.ServiceAccountName != "testServiceAccount" {
		t.Errorf("expected %s, actual %s", "testServiceAccount", dep.Spec.Template.Spec.ServiceAccountName)
	}

	if dep.Spec.Template.Spec.TopologySpreadConstraints[0].TopologyKey != "testTopology" {
		t.Errorf("expected %s, actual %s", "testTopology", dep.Spec.Template.Spec.TopologySpreadConstraints[0].TopologyKey)
	}
	if dep.Spec.Template.Spec.Tolerations[0].Key != "testKey" {
		t.Errorf("expected %s, actual %s", "testKey", dep.Spec.Template.Spec.Tolerations[0].Key)
	}
	if dep.Spec.Template.Spec.Tolerations[0].Value != "testValue" {
		t.Errorf("expected %s, actual %s", "testValue", dep.Spec.Template.Spec.Tolerations[0].Value)
	}
	if dep.Spec.Template.Spec.NodeSelector["testLabel"] != "testValue" {
		t.Errorf("expected %s, actual %s", "testValue", dep.Spec.Template.Spec.NodeSelector["testLabel"])
	}
	if dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort != 8443 {
		t.Errorf("expected %d, actual %d", 8443, dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
	}
	if dep.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort != 9443 {
		t.Errorf("expected %d, actual %d", 9443, dep.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort)
	}
}

func hasVolume(volumes []corev1.Volume, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

func hasVolumeMount(mounts []corev1.VolumeMount, name string) bool {
	for _, m := range mounts {
		if m.Name == name {
			return true
		}
	}
	return false
}

func countVolumeOccurrences(volumes []corev1.Volume, name string) int {
	count := 0
	for _, v := range volumes {
		if v.Name == name {
			count++
		}
	}
	return count
}

func gemfireDeploymentGateway() securityv1.Gateway {
	gw := gatewayForDeploymentProbeTest()
	gw.Spec.App.Gemfire = securityv1.GemfireConfigurations{
		Enabled: true,
		Locators: []securityv1.GemfireLocator{
			{Host: "locator-0", Port: 10334},
		},
	}
	return gw
}

func TestDeploymentGemfireVolumes(t *testing.T) {
	t.Run("operator-managed secret, ssl disabled: no keystore/truststore volumes, one shared-state volume", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes
		mounts := dep.Spec.Template.Spec.Containers[0].VolumeMounts

		if hasVolume(volumes, "gemfire-keystore") || hasVolume(volumes, "gemfire-truststore") {
			t.Errorf("did not expect keystore/truststore volumes when ssl is disabled, got: %+v", volumes)
		}
		if countVolumeOccurrences(volumes, "sharedstate-client-config") != 1 {
			t.Errorf("expected exactly one sharedstate-client-config volume, got %d: %+v", countVolumeOccurrences(volumes, "sharedstate-client-config"), volumes)
		}
		if !hasVolumeMount(mounts, "sharedstate-client-config") {
			t.Errorf("expected a sharedstate-client-config volume mount, got: %+v", mounts)
		}

		var sscVolume *corev1.Volume
		for i := range volumes {
			if volumes[i].Name == "sharedstate-client-config" {
				sscVolume = &volumes[i]
			}
		}
		if sscVolume == nil || sscVolume.Secret == nil || sscVolume.Secret.SecretName != gw.Name+"-shared-state-config" {
			t.Errorf("expected sharedstate-client-config volume to reference %s-shared-state-config, got: %+v", gw.Name, sscVolume)
		}
	})

	t.Run("operator-managed secret, ssl enabled with both stores: mounts both with default key names", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled: true,
			Keystore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-keystore-secret",
			},
			Truststore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-truststore-secret",
			},
		}
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes
		mounts := dep.Spec.Template.Spec.Containers[0].VolumeMounts

		if !hasVolume(volumes, "gemfire-keystore") || !hasVolume(volumes, "gemfire-truststore") {
			t.Fatalf("expected both gemfire-keystore and gemfire-truststore volumes, got: %+v", volumes)
		}
		if !hasVolumeMount(mounts, "gemfire-keystore") || !hasVolumeMount(mounts, "gemfire-truststore") {
			t.Fatalf("expected both gemfire-keystore and gemfire-truststore mounts, got: %+v", mounts)
		}

		for _, v := range volumes {
			switch v.Name {
			case "gemfire-keystore":
				if v.Secret.Items[0].Key != "keystore.jks" {
					t.Errorf("expected default keystore key 'keystore.jks', got %q", v.Secret.Items[0].Key)
				}
			case "gemfire-truststore":
				if v.Secret.Items[0].Key != "truststore.jks" {
					t.Errorf("expected default truststore key 'truststore.jks', got %q", v.Secret.Items[0].Key)
				}
			}
		}
	})

	t.Run("operator-managed secret, ssl enabled with keystore only: no truststore volume", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Gemfire.Ssl = securityv1.GemfireSsl{
			Enabled: true,
			Keystore: securityv1.GemfireStore{
				ExistingSecretName: "gemfire-keystore-secret",
				ExistingSecretKey:  "custom-keystore-key",
			},
		}
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes

		if !hasVolume(volumes, "gemfire-keystore") {
			t.Fatalf("expected gemfire-keystore volume, got: %+v", volumes)
		}
		if hasVolume(volumes, "gemfire-truststore") {
			t.Errorf("did not expect gemfire-truststore volume, got: %+v", volumes)
		}
		for _, v := range volumes {
			if v.Name == "gemfire-keystore" && v.Secret.Items[0].Key != "custom-keystore-key" {
				t.Errorf("expected custom keystore key to be respected, got %q", v.Secret.Items[0].Key)
			}
		}
	})

	t.Run("existing full secret: only enabled cert secrets get volumes", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Gemfire.ExistingSecret = "customer-managed-secret"
		gw.Spec.App.Gemfire.CertSecrets = []securityv1.GemfireCerts{
			{Enabled: false, SecretName: "skip-me", Key: "skip.jks"},
			{Enabled: true, SecretName: "gemfire-extra-secret", Key: "extra.jks"},
		}
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes
		mounts := dep.Spec.Template.Spec.Containers[0].VolumeMounts

		if hasVolume(volumes, "gemfire-ssl-0") {
			t.Errorf("did not expect a volume for a disabled cert secret, got: %+v", volumes)
		}
		if !hasVolume(volumes, "gemfire-ssl-1") {
			t.Errorf("expected a volume for the enabled cert secret at index 1, got: %+v", volumes)
		}
		if !hasVolumeMount(mounts, "gemfire-ssl-1") {
			t.Errorf("expected a volume mount for the enabled cert secret at index 1, got: %+v", mounts)
		}
		// The existing secret must itself contain a sharedstate_client.yaml key (per the
		// ExistingSecret field docs), so the sharedstate-client-config volume is still
		// mounted, just pointed at the customer-supplied secret instead of the
		// operator-managed one. CertSecrets only adds the extra keystore/truststore mounts.
		var sscVolume *corev1.Volume
		for i := range volumes {
			if volumes[i].Name == "sharedstate-client-config" {
				sscVolume = &volumes[i]
			}
		}
		if sscVolume == nil {
			t.Fatalf("expected a sharedstate-client-config volume pointing at the existing secret, got: %+v", volumes)
		}
		if sscVolume.Secret.SecretName != "customer-managed-secret" {
			t.Errorf("expected sharedstate-client-config volume to reference the existing secret, got %q", sscVolume.Secret.SecretName)
		}
	})

	t.Run("gemfire and redis both enabled: shared-state volume is not duplicated", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Redis = securityv1.RedisConfigurations{
			Enabled: true,
			Default: securityv1.RedisConfiguration{
				Type: securityv1.RedisTypeStandalone,
				Standalone: securityv1.RedisNode{
					Host: "redis-standalone",
					Port: 6379,
				},
			},
		}
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes

		if count := countVolumeOccurrences(volumes, "sharedstate-client-config"); count != 1 {
			t.Errorf("expected exactly one sharedstate-client-config volume when both redis and gemfire are enabled, got %d: %+v", count, volumes)
		}
	})

	t.Run("existingSecret precedence: redis existingSecret wins when both set", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Redis = securityv1.RedisConfigurations{
			Enabled:        true,
			ExistingSecret: "shared-secret",
			Default: securityv1.RedisConfiguration{
				Type:       securityv1.RedisTypeStandalone,
				Standalone: securityv1.RedisNode{Host: "redis-standalone", Port: 6379},
			},
		}
		gw.Spec.App.Gemfire.ExistingSecret = "shared-secret"
		dep := NewDeployment(&gw, "kubernetes")
		volumes := dep.Spec.Template.Spec.Volumes

		var sscVolume *corev1.Volume
		for i := range volumes {
			if volumes[i].Name == "sharedstate-client-config" {
				sscVolume = &volumes[i]
			}
		}
		if sscVolume == nil {
			t.Fatalf("expected a sharedstate-client-config volume, got: %+v", volumes)
		}
		if sscVolume.Secret.SecretName != "shared-secret" {
			t.Errorf("expected secret name 'shared-secret', got %q", sscVolume.Secret.SecretName)
		}
	})

	t.Run("only gemfire existingSecret set: it is used for the shared volume name", func(t *testing.T) {
		gw := gemfireDeploymentGateway()
		gw.Spec.App.Gemfire.ExistingSecret = "gemfire-only-secret"
		gw.Spec.App.Gemfire.CertSecrets = []securityv1.GemfireCerts{
			{Enabled: true, SecretName: "extra", Key: "extra.jks"},
		}
		dep := NewDeployment(&gw, "kubernetes")

		for _, v := range dep.Spec.Template.Spec.Volumes {
			if v.Name == "gemfire-ssl-0" && v.Secret.SecretName != "extra" {
				t.Errorf("expected cert secret volume to reference 'extra', got %q", v.Secret.SecretName)
			}
		}
	})
}
