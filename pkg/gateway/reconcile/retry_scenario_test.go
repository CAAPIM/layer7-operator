/*
*
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
*
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/go-logr/logr"
)

func TestCheckRetryScenario_commitMarkerRequired(t *testing.T) {
	tmp := t.TempDir()
	commit := "newsha1"
	repoRef := "myrepo"
	params := Params{Log: logr.Discard()}

	gw := &securityv1.Gateway{
		Status: securityv1.GatewayStatus{
			RepositoryStatus: []securityv1.GatewayRepositoryStatus{
				{
					Name:   repoRef,
					Commit: commit,
					Conditions: []securityv1.RepositoryCondition{
						{Reason: "GraphQL apply failed", Status: "False"},
					},
				},
			},
		},
	}

	should, bundle := checkRetryScenario(gw, repoRef, commit, tmp, params)
	if should || len(bundle) != 0 {
		t.Fatalf("without marker: want shouldRetry=false and empty bundle, got should=%v len=%d", should, len(bundle))
	}
}

func TestCheckRetryScenario_markerAndLastApplied_retries(t *testing.T) {
	tmp := t.TempDir()
	commit := "abc123"
	repoRef := "myrepo"
	want := []byte(`{"bundle":true}`)
	params := Params{Log: logr.Discard()}

	if err := os.WriteFile(filepath.Join(tmp, commit+".txt"), []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "last_applied_"+repoRef+".json"), want, 0o600); err != nil {
		t.Fatal(err)
	}

	gw := &securityv1.Gateway{
		Status: securityv1.GatewayStatus{
			RepositoryStatus: []securityv1.GatewayRepositoryStatus{
				{
					Name:   repoRef,
					Commit: commit,
					Conditions: []securityv1.RepositoryCondition{
						{Reason: "error: transient", Status: "False"},
					},
				},
			},
		},
	}

	should, bundle := checkRetryScenario(gw, repoRef, commit, tmp, params)
	if !should {
		t.Fatal("with marker and last_applied: expected shouldRetry true")
	}
	if string(bundle) != string(want) {
		t.Fatalf("bundle: got %q want %q", bundle, want)
	}
}

func TestCheckRetryScenario_markerNoLastApplied_shouldRetryTrueNilBundle(t *testing.T) {
	tmp := t.TempDir()
	commit := "abc123"
	repoRef := "myrepo"
	params := Params{Log: logr.Discard()}

	if err := os.WriteFile(filepath.Join(tmp, commit+".txt"), []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}

	gw := &securityv1.Gateway{
		Status: securityv1.GatewayStatus{
			RepositoryStatus: []securityv1.GatewayRepositoryStatus{
				{
					Name:   repoRef,
					Commit: commit,
					Conditions: []securityv1.RepositoryCondition{
						{Reason: "failed apply", Status: "False"},
					},
				},
			},
		},
	}

	should, bundle := checkRetryScenario(gw, repoRef, commit, tmp, params)
	if !should || bundle != nil {
		t.Fatalf("marker without last_applied: want should=true, nil bundle; got should=%v bundle=%v", should, bundle)
	}
}
