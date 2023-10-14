package util

import (
	"archive/tar"
	"archive/zip"
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

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeXGlobalHeader:
			continue
		case tar.TypeDir:
			if err := os.Mkdir("/tmp/"+repoName+"-"+header.Name, 0755); err != nil {
				return fmt.Errorf("failed to create folder %s", header.Name)
			}
		case tar.TypeReg:
			outFile, err := os.Create("/tmp/" + repoName + "-" + header.Name)
			if err != nil {
				return fmt.Errorf("failed to create file %s", header.Name)
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("copy failed: %s", err)

			}
		default:
			return fmt.Errorf("uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}
