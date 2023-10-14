package util

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func DownloadArtifact(URL string, username string, token string, name string) (string, error) {
	/////TODO: REFACTOR
	fileURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := "/tmp/" + name + "-" + segments[len(segments)-1]

	sha1sum := existingFileSha(fileName)

	ext := strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1]
	folderName := strings.ReplaceAll(fileName, "."+ext, "")
	if strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-2] == "tar" {
		folderName = strings.ReplaceAll(fileName, ".tar.gz", "")
	}

	existingFolder := existingFolder(folderName)

	if existingFolder && sha1sum != "" {
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

	case "json":
		break
	}

	sha1sum = existingFileSha(fileName)
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
