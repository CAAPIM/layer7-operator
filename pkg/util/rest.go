package util

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func RestCall(method string, URL string, insecureSkipVerify bool, headers map[string]string, contentType string, data []byte, username string, password string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, URL, strings.NewReader(string(data)))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", contentType)

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return []byte{}, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
