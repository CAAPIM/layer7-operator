package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/caapim/layer7-operator/internal/graphman"
)

var ErrInvalidFileFormatError = errors.New("InvalidFileFormat")
var ErrInvalidArchive = errors.New("InvalidArchive")

// Download Artifact retrieves a compressed Graphman Bundle from an HTTP URL
// This is currently limited to URLs that contain the file extension as would be the
// case when targeting releases from Git releases.
// The following extensions are accepted .tar, .tar.gz, .zip
func DownloadArtifact(URL string, username string, token string, name string, forceUpdate bool) (string, error) {
	fileURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := "/tmp/" + name + "-" + segments[len(segments)-1]

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

	folderName := strings.ReplaceAll(fileName, "."+ext, "")
	if strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
		folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
	}

	existingFolder := existingFolder(folderName)

	// If the downloaded file and corresponding uncompressed folder does not represent a valid graphman bundle
	// it will be removed allowing for additional attempts to retrieve the file.
	if existingFolder && sha1sum != "" && !forceUpdate {
		err = validateGraphmanBundle(fileName, folderName, name)
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
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		f, _ := os.Open(fileName)
		defer f.Close()
		err = Unzip(fileName, folderName)
		if err != nil {
			return "", ErrInvalidArchive
		}
	case "gz":
		gz := false
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
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
			return "", ErrInvalidArchive
		}
	case "tar":
		gz := false
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		f, _ := os.Open(fileName)
		defer f.Close()

		err := Untar(folderName, name, f, gz)
		if err != nil {
			return "", ErrInvalidArchive
		}
	case "json":
		fileName = name + "-" + segments[len(segments)-1]

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

	err = validateGraphmanBundle(fileName, folderName, name)
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

func validateGraphmanBundle(fileName string, folderName string, repoName string) error {
	bundle := graphman.Bundle{}
	if _, err := os.Stat(folderName); err != nil {
		return nil
	}

	bundleBytes, err := graphman.Implode(folderName)
	if err != nil {
		return err
	}

	if len(bundleBytes) <= 2 {
		files, err := os.ReadDir(folderName)
		if err != nil {
			return err
		}

		for _, f := range files {
			segments := strings.Split(f.Name(), ".")
			ext := segments[len(segments)-1]
			if ext == "json" {
				sbb := bundleBytes
				srcBundleBytes, err := os.ReadFile(folderName + "/" + f.Name())
				if err != nil {
					return err
				}
				bundleBytes, err = graphman.ConcatBundle(srcBundleBytes, bundleBytes)
				if err != nil {
					continue
				}
				bundleBytes = sbb
			}
		}
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
