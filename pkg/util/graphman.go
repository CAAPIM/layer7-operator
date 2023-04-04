package util

import (
	"bytes"
	"compress/gzip"
	"errors"

	graphman "gitlab.sutraone.com/gazza/go-graphman"
)

func ApplyToGraphmanTarget(path string, username string, password string, target string, encpass string) error {
	_, err := graphman.Apply(path, username, password, "https://"+target, encpass)
	if err != nil {
		return err
	}
	return nil
}

func CompressGraphmanBundle(path string) ([]byte, error) {
	bundle, err := graphman.Implode(path)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(bundle)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	if buf.Len() > 900000 {
		return nil, errors.New("this bundle would exceed the maximum Kubernetes secret size.")

	}

	return buf.Bytes(), nil
}
