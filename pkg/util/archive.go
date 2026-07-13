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
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) error {
	folderExists, _ := os.Stat(dest)
	if folderExists != nil {
		return nil
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		if strings.Contains(f.Name, "__MACOSX") || strings.Contains(f.Name, "._") || strings.Contains(f.Name, "._") || strings.Contains(f.Name, "./._") || strings.Contains(f.Name, "/._") {
			continue
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// These are vars (not consts) only so tests can lower them; treat them as
// constants in production code.
var (
	// maxDecompressedSize is the maximum total number of decompressed bytes
	// accepted from an archive — cumulative across all tar entries, or the size
	// of a single gzip payload. 500 MiB prevents decompression bombs from
	// exhausting memory or disk. Exceeding it is a hard error, not a truncation.
	maxDecompressedSize int64 = 500 << 20 // 500 MiB

	// maxArchiveEntries is the maximum number of entries accepted from a tar
	// archive. Many tiny files can be used as a zip-bomb variant; this cap
	// prevents that class of attack.
	maxArchiveEntries = 10_000
)

// safeJoin joins base and name and verifies the result is contained within base.
// It returns an error when name contains path traversal sequences that would
// escape the intended target directory (TarSlip / ZipSlip protection).
func safeJoin(base, name string) (string, error) {
	p := filepath.Clean(filepath.Join(base, name))
	if !strings.HasPrefix(p, filepath.Clean(base)+string(os.PathSeparator)) {
		return "", fmt.Errorf("archive entry %q escapes target directory", name)
	}
	return p, nil
}

func Untar(folderName string, repoName string, tarStream io.Reader, gz bool) error {
	folderExists, _ := os.Stat(folderName)
	if folderExists != nil {
		return nil
	}

	tarReader := tar.NewReader(tarStream)

	if gz {
		uncompressedStream, err := gzip.NewReader(tarStream)
		if err != nil {
			return err
		}
		tarReader = tar.NewReader(uncompressedStream)
	}

	// Create the folder up front and return an error if it doesn't exist while iterating over tar entries
	_ = os.Mkdir(folderName, 0755)

	entryCount := 0
	var totalWritten int64
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		entryCount++
		if entryCount > maxArchiveEntries {
			return fmt.Errorf("archive exceeds maximum entry count of %d", maxArchiveEntries)
		}

		switch header.Typeflag {
		case tar.TypeXGlobalHeader:
			continue
		case tar.TypeDir:
			if header.Name == "./" {
				continue
			}

			entryName := strings.TrimPrefix(header.Name, "./")
			path, err := safeJoin(folderName, entryName)
			if err != nil {
				return err
			}

			if err := os.Mkdir(path, 0755); err != nil {
				if _, err = os.Stat(path); err != nil {
					return fmt.Errorf("failed to create folder %s", header.Name)
				}
			}
		case tar.TypeReg:
			// Ignore BSD Tar macOS metadata entries
			if strings.HasPrefix(header.Name, "./._") || strings.HasPrefix(header.Name, "._") || strings.Contains(header.Name, "/._") {
				continue
			}

			entryName := strings.TrimPrefix(header.Name, "./")
			path, err := safeJoin(folderName, entryName)
			if err != nil {
				return err
			}

			outFile, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create file %s", header.Name)
			}
			// Enforce a cumulative decompressed-size budget across all entries so
			// that many individually-small files can't sum past the cap. Read one
			// byte past the remaining budget so an overflow is detected and
			// surfaced as an error rather than silently truncated. Close the file
			// each iteration instead of deferring to avoid holding up to
			// maxArchiveEntries file descriptors open until Untar returns.
			remaining := maxDecompressedSize - totalWritten
			written, err := io.Copy(outFile, io.LimitReader(tarReader, remaining+1))
			outFile.Close()
			if err != nil {
				return fmt.Errorf("copy failed: %s", err)
			}
			if written > remaining {
				return fmt.Errorf("archive exceeds maximum decompressed size of %d bytes", maxDecompressedSize)
			}
			totalWritten += written
		default:
			return fmt.Errorf("uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}

func GzipDecompress(gzipBundle []byte) (bundleBytes []byte, err error) {
	r := bytes.NewReader(gzipBundle)
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	// Read one byte past the cap so an oversized payload is surfaced as an error
	// rather than silently truncated.
	bundleBytes, err = io.ReadAll(io.LimitReader(gzr, maxDecompressedSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(bundleBytes)) > maxDecompressedSize {
		return nil, fmt.Errorf("gzip payload exceeds maximum decompressed size of %d bytes", maxDecompressedSize)
	}

	return bundleBytes, nil
}

func GzipCompress(gzipBundle []byte) (gzipBytes []byte, err error) {

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(gzipBundle)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	gzipBytes = buf.Bytes()
	buf.Reset()

	return gzipBytes, nil
}
