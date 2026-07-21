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
package util

import (
	"strings"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
)

func TestBuildCWPBundle(t *testing.T) {
	cwps := []securityv1.Property{{Name: "test1", Value: "value1"}}
	bundleBytes, sha1, err := BuildCWPBundle(cwps)
	if err != nil {
		t.Errorf("Error getting cwp bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "test1") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "test1")
	}
}

func TestBuildDefaultListenPortBundle(t *testing.T) {
	bundleBytes, sha1, err := BuildDefaultListenPortBundle(true)
	if err != nil {
		t.Errorf("Error getting default listen port bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "8443") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "8443")
	}
	if !strings.Contains(bundle, "9443") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "9443")
	}

	if !strings.Contains(bundle, "8080") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "8080")
	}
}

func TestBuildCustomListenPortBundle(t *testing.T) {
	gateway := securityv1.Gateway{}
	gateway.Spec = securityv1.GatewaySpec{}
	gateway.Spec.App = securityv1.App{}
	gateway.Spec.App.ListenPorts = securityv1.ListenPorts{}
	gateway.Spec.App.ListenPorts.Custom = securityv1.CustomListenPort{Enabled: true} //ecurityv1.CustomListenPort{true}

	port := securityv1.ListenPort{}
	port.Name = "custom1"
	port.Enabled = true
	port.Protocol = "http"
	port.Port = 9090
	gateway.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{port}

	bundleBytes, sha1, err := BuildCustomListenPortBundle(&gateway, false)
	if err != nil {
		t.Errorf("Error getting default listen port bundle")
	}
	bundle := string(bundleBytes)
	if !strings.Contains(bundle, "9090") {
		t.Errorf("bundle %s, sha1 %s, expected key %s", bundle, sha1, "9090")
	}
}

func TestListenPortNamesExcludedFromGraphmanSync(t *testing.T) {
	gw := securityv1.Gateway{}
	gw.Spec.App.ListenPorts.Ports = []securityv1.ListenPort{
		{Name: "a", Port: 8443},
		{Name: "graphman-lp", Port: 9443},
	}
	ex := ListenPortNamesExcludedFromGraphmanSync(&gw, 9443)
	if _, ok := ex["graphman-lp"]; !ok {
		t.Fatalf("expected graphman-lp excluded for port 9443, got %v", ex)
	}
	if _, ok := ex["a"]; ok {
		t.Fatalf("did not expect port 8443 in excluded set")
	}
}

func TestFilterListenPortBundleForGraphmanSync(t *testing.T) {
	b := graphman.Bundle{
		ListenPorts: []*graphman.ListenPortInput{
			{Name: "keep", Port: 8443},
			{Name: "drop", Port: 9443},
		},
	}
	FilterListenPortBundleForGraphmanSync(&b, 9443)
	if len(b.ListenPorts) != 1 || b.ListenPorts[0].Name != "keep" {
		t.Fatalf("expected single listen port keep, got %+v", b.ListenPorts)
	}
}
