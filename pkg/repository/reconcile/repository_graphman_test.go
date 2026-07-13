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
package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-logr/logr"
)

func TestAnalyzeGraphmanCloneRoot_EmptyClone(t *testing.T) {
	d := t.TempDir()
	_, _, err := analyzeGraphmanCloneRoot(d, logr.Discard())
	if err == nil {
		t.Fatal("expected bundle build to fail for empty clone")
	}
}

func TestAnalyzeGraphmanCloneRoot_FlatJSONBundle(t *testing.T) {
	d := t.TempDir()
	// No graphman *project* dirs (DetectGraphmanFolders empty), but root-level JSON graphman file
	// that BuildAndValidateBundle still merges — same pattern as a "pure" combined layout.
	j := `{"services":[{"name":"s","resolutionPath":"/x","enabled":true}]}`
	if err := os.WriteFile(filepath.Join(d, "svc.json"), []byte(j), 0644); err != nil {
		t.Fatal(err)
	}
	projects, err := util.DetectGraphmanFolders(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected no project dirs, got %v", projects)
	}
	warn, useRoot, err := analyzeGraphmanCloneRoot(d, logr.Discard())
	if err != nil {
		t.Fatal(err)
	}
	if warn || !useRoot {
		t.Fatalf("expected root bundle without warn, got warn=%v useRoot=%v", warn, useRoot)
	}
}

func TestAnalyzeGraphmanCloneRoot_ClusterPropertiesJSON(t *testing.T) {
	d := t.TempDir()
	// Flat repo: only test.json at root (no per-project dirs); clusterProperties is valid graphman content.
	j := `{
    "clusterProperties": [
        {
            "goid": "982cc1ee7369c6ca5a7ae1e4ad866070",
            "name": "myDemoConfigVal123",
            "hiddenProperty": false,
            "value": "suspiciousLlama - changed"
        }
    ]
}`
	if err := os.WriteFile(filepath.Join(d, "test.json"), []byte(j), 0644); err != nil {
		t.Fatal(err)
	}
	warnNoProjectDirs, useCloneRootAsProject, err := analyzeGraphmanCloneRoot(d, logr.Discard())
	if err != nil {
		t.Fatal(err)
	}
	if warnNoProjectDirs || !useCloneRootAsProject {
		t.Fatalf("expected root bundle: warnNoProjectDirs=%v useCloneRootAsProject=%v", warnNoProjectDirs, useCloneRootAsProject)
	}
	b, err := util.BuildAndValidateBundle(d, false, logr.Discard())
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 || util.GraphmanBundleBytesHaveNoEntities(b) {
		t.Fatal("clusterProperties bundle should not be empty")
	}
}
