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
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

func TestNewIngress(t *testing.T) {
	gateway := securityv1.Gateway{}
	gateway.Name = "test"

	gateway.Spec.App.Ingress = securityv1.Ingress{
		Enabled:          true,
		Annotations:      map[string]string{"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS"},
		IngressClassName: "nginx",
		TLS:              []networkingv1.IngressTLS{{Hosts: []string{"testing.com"}, SecretName: "default"}},
		Rules:            []networkingv1.IngressRule{{Host: "testing.com"}},
	}

	gateway.Spec.App.Service = securityv1.Service{
		Enabled: true,
		Type:    "ClusterIP",
		Ports:   []securityv1.Ports{{Port: 8443, TargetPort: 8443, Protocol: "TCP", Name: "https"}},
	}

	ingress := NewIngress(&gateway)
	if ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name != gateway.Spec.App.Service.Ports[0].Name {
		t.Errorf("expected %s, actual %s", ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name, gateway.Spec.App.Service.Ports[0].Name)
	}

}

func TestNewIngress_emptyIngressClassNameIsNil(t *testing.T) {
	gw := securityv1.Gateway{}
	gw.Name = "test"
	gw.Namespace = "default"
	gw.Spec.App.Ingress = securityv1.Ingress{
		Enabled:          true,
		IngressClassName: "",
		Rules:            []networkingv1.IngressRule{{Host: "example.com"}},
	}
	gw.Spec.App.Service = securityv1.Service{
		Enabled: true,
		Type:    "ClusterIP",
		Ports:   []securityv1.Ports{{Port: 8443, TargetPort: 8443, Protocol: "TCP", Name: "https"}},
	}
	ing := NewIngress(&gw)
	if ing.Spec.IngressClassName != nil {
		t.Fatalf("expected nil IngressClassName when empty, got %#v", ing.Spec.IngressClassName)
	}
}
