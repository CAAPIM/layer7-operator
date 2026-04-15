/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* AI assistance has been used to generate some or all contents of this file.
 */
package reconcile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommitMarkerPresent(t *testing.T) {
	tmp := t.TempDir()
	commit := "commit"

	if commitMarkerPresent(tmp, commit) {
		t.Fatal("missing marker should be false")
	}

	path := filepath.Join(tmp, commit+".txt")
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	if !commitMarkerPresent(tmp, commit) {
		t.Fatal("existing marker should be true")
	}
}
