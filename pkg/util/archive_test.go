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
	"archive/tar"
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeTar builds an uncompressed tar with the given regular-file entries.
func makeTar(t *testing.T, entries map[string]int) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, size := range entries {
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Mode:     0644,
			Size:     int64(size),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatalf("write header: %v", err)
		}
		if _, err := tw.Write(bytes.Repeat([]byte("a"), size)); err != nil {
			t.Fatalf("write body: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	return buf.Bytes()
}

// withLimits temporarily lowers the archive caps for a test and restores them.
func withLimits(t *testing.T, maxSize int64, maxEntries int) {
	t.Helper()
	origSize, origEntries := maxDecompressedSize, maxArchiveEntries
	maxDecompressedSize, maxArchiveEntries = maxSize, maxEntries
	t.Cleanup(func() { maxDecompressedSize, maxArchiveEntries = origSize, origEntries })
}

func TestUntarCumulativeSizeExceeded(t *testing.T) {
	withLimits(t, 100, 10_000)

	// Three 40-byte files = 120 bytes total, over the 100-byte cumulative cap,
	// even though no single file exceeds it.
	data := makeTar(t, map[string]int{"a.txt": 40, "b.txt": 40, "c.txt": 40})

	dest := filepath.Join(t.TempDir(), "out")
	err := Untar(dest, "repo", bytes.NewReader(data), false)
	if err == nil {
		t.Fatal("expected cumulative size error, got nil (silent truncation?)")
	}
	if !strings.Contains(err.Error(), "maximum decompressed size") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUntarWithinBudget(t *testing.T) {
	withLimits(t, 1000, 10_000)

	data := makeTar(t, map[string]int{"a.txt": 40, "b.txt": 40})
	dest := filepath.Join(t.TempDir(), "out")
	if err := Untar(dest, "repo", bytes.NewReader(data), false); err != nil {
		t.Fatalf("Untar within budget failed: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(dest, "a.txt"))
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if len(b) != 40 {
		t.Fatalf("extracted file truncated: got %d bytes, want 40", len(b))
	}
}

func TestUntarEntryCountExceeded(t *testing.T) {
	withLimits(t, 1<<20, 2)

	data := makeTar(t, map[string]int{"a": 1, "b": 1, "c": 1})
	dest := filepath.Join(t.TempDir(), "out")
	err := Untar(dest, "repo", bytes.NewReader(data), false)
	if err == nil || !strings.Contains(err.Error(), "maximum entry count") {
		t.Fatalf("expected entry-count error, got %v", err)
	}
}

// writeZip writes a zip archive with the given regular-file entries (name ->
// size, contents all zeros so they compress very well) and returns its path.
func writeZip(t *testing.T, entries map[string]int) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "archive.zip")
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	for name, size := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create entry: %v", err)
		}
		if _, err := w.Write(make([]byte, size)); err != nil {
			t.Fatalf("write entry: %v", err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}
	return p
}

func TestUnzipBombRejected(t *testing.T) {
	withLimits(t, 64<<10, 10_000)

	// A single highly-compressible 1 MiB entry, far over the 64 KiB budget.
	src := writeZip(t, map[string]int{"bomb.bin": 1 << 20})
	dest := filepath.Join(t.TempDir(), "out")
	err := Unzip(src, dest)
	if err == nil || !strings.Contains(err.Error(), "maximum decompressed size") {
		t.Fatalf("expected decompressed-size error, got %v", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("partial extraction left behind at %s (a retry would treat it as success)", dest)
	}
	// A retry must still fail rather than short-circuit on an existing dest.
	if err := Unzip(src, dest); err == nil {
		t.Fatal("retry after rejected archive returned nil")
	}
}

func TestUnzipCumulativeSizeExceeded(t *testing.T) {
	withLimits(t, 100, 10_000)

	// Three 40-byte files = 120 bytes total; no single file exceeds the cap.
	src := writeZip(t, map[string]int{"a.txt": 40, "b.txt": 40, "c.txt": 40})
	err := Unzip(src, filepath.Join(t.TempDir(), "out"))
	if err == nil || !strings.Contains(err.Error(), "maximum decompressed size") {
		t.Fatalf("expected cumulative size error, got %v", err)
	}
}

func TestUnzipEntryCountExceeded(t *testing.T) {
	withLimits(t, 1<<20, 2)

	src := writeZip(t, map[string]int{"a": 1, "b": 1, "c": 1})
	err := Unzip(src, filepath.Join(t.TempDir(), "out"))
	if err == nil || !strings.Contains(err.Error(), "maximum entry count") {
		t.Fatalf("expected entry-count error, got %v", err)
	}
}

func TestUnzipWithinBudget(t *testing.T) {
	withLimits(t, 1000, 10)

	src := writeZip(t, map[string]int{"a.txt": 40, "dir/b.txt": 40})
	dest := filepath.Join(t.TempDir(), "out")
	if err := Unzip(src, dest); err != nil {
		t.Fatalf("Unzip within budget failed: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(dest, "dir", "b.txt"))
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if len(b) != 40 {
		t.Fatalf("extracted file truncated: got %d bytes, want 40", len(b))
	}
}
