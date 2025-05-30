/*
* Copyright (c) 2025 Broadcom. All rights reserved.
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
 */
package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/caapim/layer7-operator/internal/graphman"
)

var ErrInvalidFileFormatError = errors.New("InvalidFileFormat")
var ErrInvalidTarArchive = errors.New("InvalidTarArchive")
var ErrInvalidZipArchive = errors.New("InvalidZipArchive")

// Download Artifact retrieves a compressed Graphman Bundle from an HTTP URL
// This is currently limited to URLs that contain the file extension as would be the
// case when targeting releases from Git releases.
// The following extensions are accepted .tar, .tar.gz, .zip
func DownloadArtifact(URL string, username string, token string, name string, forceUpdate bool, namespace string) (string, error) {
	fileURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	folderName := ""
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := "/tmp/" + name + "-" + namespace + "-" + segments[len(segments)-1]

	// Downloaded artifacts are treated as immutable, once a given release has been downloaded it will not be
	// retrieved again unless there is a URL change.
	// If a downloaded artifact is invalid it is initially assumned that the download failed, if a file fails validation checks
	// multiple times a backoff is initiated preventing the invalid file from being downloaded multiple times.
	sha1sum := existingFileSha(fileName)
	ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]

	if ext != "zip" && ext != "tar" && ext != "gz" && ext != "json" {
		return "", fmt.Errorf("unsupported file type %s", ext)
	}

	if ext == "gz" && strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] != "tar" {
		return "", fmt.Errorf("unsupported file type %s", ext)
	}

	folderName = strings.ReplaceAll(fileName, "."+ext, "")
	if strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
		folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
	}

	existingFolder := existingFolder(folderName)

	// If the downloaded file and corresponding uncompressed folder does not represent a valid graphman bundle
	// it will be removed allowing for additional attempts to retrieve the file.
	if existingFolder && sha1sum != "" && !forceUpdate {
		err = validateGraphmanBundle(fileName, folderName)
		if err != nil {
			return "", err
		}
		return sha1sum, nil
	}

	if forceUpdate {
		os.RemoveAll(folderName)
	}

	resp, err := RestCall("GET", URL, true, map[string]string{}, "text/plain", []byte{}, username, token)

	if err != nil {
		return "", err
	}

	err = os.WriteFile(fileName, resp, 0666)
	if err != nil {
		return "", err
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}

	if fileInfo.Size() == 0 {
		return "", errors.New("file is empty")
	}

	switch ext {
	case "zip":
		folderName = strings.ReplaceAll(fileName, "."+ext, "")
		f, _ := os.Open(fileName)
		defer f.Close()
		err = Unzip(fileName, folderName)
		if err != nil {
			return "", ErrInvalidZipArchive
		}
	case "gz":
		gz := false
		f, _ := os.Open(fileName)
		defer f.Close()
		if strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
			folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
			gz = true
		} else {
			return "", errors.New(".gz is an unsupported file format")
		}

		err := Untar(folderName, name, f, gz)
		if err != nil {
			return "", ErrInvalidTarArchive
		}
	case "tar":
		gz := false
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		f, _ := os.Open(fileName)
		defer f.Close()

		err := Untar(folderName, name, f, gz)
		if err != nil {
			return "", ErrInvalidTarArchive
		}
	case "json":
		fileName = name + "-" + namespace + "-" + segments[len(segments)-1]

		folderName := strings.ReplaceAll("/tmp/"+fileName, "."+ext, "")
		os.Mkdir(folderName, 0755)
		fBytes, err := os.ReadFile("/tmp/" + fileName)
		if err != nil {
			return "", err
		}
		err = os.WriteFile(folderName+"/"+fileName, fBytes, 0755)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("extension %s not supported", ext)
	}

	sha1sum = existingFileSha(fileName)

	folderName = strings.ReplaceAll(fileName, "."+ext, "")
	if strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
		folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
	}

	err = validateGraphmanBundle(fileName, folderName)
	if err != nil {
		return "", err
	}
	return sha1sum, nil
}

func existingFileSha(fileName string) string {
	_, err := os.Stat(fileName)
	if err != nil {
		return ""
	}

	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}

	h := sha1.New()
	h.Write(fileBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))
	fileCheckSum := sha1Sum

	return fileCheckSum
}

func existingFolder(folderName string) bool {
	folderExists, _ := os.Stat(folderName)
	return folderExists != nil
}

func validateGraphmanBundle(fileName string, folderName string) error {
	bundle := graphman.Bundle{}
	if _, err := os.Stat(folderName); err != nil {
		return nil
	}

	bundleBytes, err := graphman.Implode(folderName)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(folderName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "/.git") {
			return nil
		}
		if !d.IsDir() {
			segments := strings.Split(d.Name(), ".")
			ext := segments[len(segments)-1]
			if ext == "json" && !strings.Contains(strings.ToLower(d.Name()), "sourcesummary.json") && !strings.Contains(strings.ToLower(d.Name()), "bundle-properties.json") {
				//sbb := bundleBytes
				srcBundleBytes, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				sbb, err := graphman.ConcatBundle(srcBundleBytes, bundleBytes)
				if err != nil {
					return nil
				}
				bundleBytes = sbb
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(bundleBytes) <= 2 {
		os.Remove(fileName)
		err = os.RemoveAll(folderName)
		if err != nil {
			return err
		}
		return ErrInvalidFileFormatError
	}

	r := bytes.NewReader(bundleBytes)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	_ = json.Unmarshal(bundleBytes, &bundle)
	// check the graphman bundle for errors
	err = d.Decode(&bundle)
	if err != nil || len(bundleBytes) <= 2 {
		os.Remove(fileName)
		fErr := os.RemoveAll(folderName)
		if fErr != nil {
			return err
		}
		return ErrInvalidFileFormatError
	}
	return nil
}
