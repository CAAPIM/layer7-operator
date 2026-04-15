package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caapim/layer7-operator/pkg/util"
)

func TestShouldWarnNoGraphmanProjectDirs_EmptyClone(t *testing.T) {
	d := t.TempDir()
	if !shouldWarnNoGraphmanProjectDirs(d) {
		t.Fatal("expected warning for empty clone with no bundles")
	}
}

func TestShouldWarnNoGraphmanProjectDirs_FlatJSONBundle(t *testing.T) {
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
	if shouldWarnNoGraphmanProjectDirs(d) {
		t.Fatal("should not warn when root produces a non-empty graphman bundle")
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
	warnNoProjectDirs, useCloneRootAsProject, b := analyzeGraphmanCloneRoot(d)
	if warnNoProjectDirs || !useCloneRootAsProject || len(b) == 0 {
		t.Fatalf("expected root bundle: warnNoProjectDirs=%v useCloneRootAsProject=%v len=%d", warnNoProjectDirs, useCloneRootAsProject, len(b))
	}
	if util.GraphmanBundleBytesHaveNoEntities(b) {
		t.Fatal("clusterProperties bundle should not be empty")
	}
}
