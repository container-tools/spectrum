package builder

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/pkg/errors"
)

func Build(options Options, dirs ...string) error {
	base, err := Pull(options)
	if err != nil {
		return errors.Wrapf(err, "could not pull base image image %s", options.Base)
	}

	tarFiles := make([]string, 0)
	for _, spec := range dirs {
		parts := strings.Split(spec, ":")
		if len(parts) != 2 {
			return errors.New("wrong dir format for " + spec + " (expected \"local:remote\")")
		}
		tarFile, err := tarPackage(parts[0], parts[1])
		if err != nil {
			return errors.Wrapf(err, "cannot package dir %s as tar file", parts[0])
		}
		defer os.Remove(tarFile)
		tarFiles = append(tarFiles, tarFile)
	}
	newImage, err := crane.Append(base, tarFiles...)
	if err != nil {
		return errors.Wrap(err, "could not append tar layers to base image")
	}

	return Push(newImage, options)
}

func tarPackage(dirName, targetPath string) (file string, err error) {
	layerFile, err := ioutil.TempFile("", "spectrum-layer-*.tar")
	if err != nil {
		return "", err
	}
	defer layerFile.Close()

	writer := tar.NewWriter(layerFile)

	dir, err := os.Open(dirName)
	if err != nil {
		return "", err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			return "", err
		}
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = path.Join(targetPath, filepath.Base(file.Name()))
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = writer.WriteHeader(header)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return "", err
		}
	}

	return layerFile.Name(), nil
}
