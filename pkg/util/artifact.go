package util

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/caapim/layer7-operator/internal/graphman"
)

// Download Artifact retrieves a compressed Graphman Bundle from an HTTP URL
// This is currently limited to URLs that contain the file extension as would be the
// case when targeting releases from Git releases.
// The following extensions are accepted .tar, .tar.gz, .zip
func DownloadArtifact(URL string, username string, token string, name string) (string, error) {
	fileURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := "/tmp/" + name + "-" + segments[len(segments)-1]

	// Downloaded artifacts are treated as immutable, once a given release has been downloaded it will not be
	// retrieved again.
	sha1sum := existingFileSha(fileName)
	ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]

	if ext != "zip" && ext != "tar" && ext != "gz" {
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
	if existingFolder && sha1sum != "" {
		err = validateGraphmanBundle(fileName, folderName, name)
		if err != nil {
			return "", err
		}
		return sha1sum, nil
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
			return "", err
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
			return "", err
		}
	case "tar":
		gz := false
		folderName := strings.ReplaceAll(fileName, "."+ext, "")
		f, _ := os.Open(fileName)
		defer f.Close()

		err := Untar(folderName, name, f, gz)
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
	_, err := graphman.Implode(folderName)
	if err != nil {
		os.Remove(fileName)
		err = os.RemoveAll(folderName)
		if err != nil {
			return err
		}
		return fmt.Errorf("repository %s does not contain a valid graphman bundle", repoName)
	}
	return nil
}
