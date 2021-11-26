package builder

import (
	"archive/tar"
	"fmt"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pkg/errors"
)

const LogPrefix = "spectrum - "

var StepLogger = log.New(ioutil.Discard, LogPrefix, log.LstdFlags)

// Build executes the full build cycle and returns the image digest
func Build(options Options, dirs ...string) (string, error) {
	configureLogging(options)
	StepLogger.Printf("Pulling base image %s (insecure=%v)...", options.Base, options.PullInsecure)
	base, err := Pull(options)
	if err != nil {
		return "", errors.Wrapf(err, "could not pull base image image %s", options.Base)
	}

	StepLogger.Println("Composing layers...")
	tarFiles := make([]string, 0)
	for _, spec := range dirs {
		parts := strings.Split(spec, ":")
		if len(parts) != 2 {
			return "", errors.New("wrong dir format for " + spec + " (expected \"local:remote\")")
		}
		tarFile, err := tarPackage(parts[0], parts[1], options.Recursive)
		if err != nil {
			return "", errors.Wrapf(err, "cannot package dir %s as tar file", parts[0])
		}
		defer os.Remove(tarFile)
		tarFiles = append(tarFiles, tarFile)
	}
	newImage, err := appendPaths(base, options.Annotations, tarFiles...)
	if err != nil {
		return "", errors.Wrap(err, "could not append tar layers to base image")
	}

	StepLogger.Printf("Pushing image %s (insecure=%v)...", options.Target, options.PushInsecure)
	if err := Push(newImage, options); err != nil {
		return "", err
	}
	var hash v1.Hash
	if hash, err = newImage.Digest(); err != nil {
		return "", err
	}
	return hash.String(), nil
}

func configureLogging(options Options) {
	stdout := options.Stdout
	if stdout == nil {
		stdout = ioutil.Discard
	}
	logs.Progress = log.New(stdout, LogPrefix, log.LstdFlags)
	StepLogger = log.New(stdout, LogPrefix, log.LstdFlags)

	stderr := options.Stderr
	if stderr == nil {
		stderr = ioutil.Discard
	}
	logs.Warn = log.New(stderr, LogPrefix, log.LstdFlags)
}

func tarPackage(name, targetPath string, recursive bool) (file string, err error) {
	layerFile, err := ioutil.TempFile("", "spectrum-layer-*.tar")
	if err != nil {
		return "", err
	}
	defer layerFile.Close()

	writer := tar.NewWriter(layerFile)
	defer writer.Close()
	fileInfo, err := os.Stat(name)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		err := writeFileToTar(name, targetPath, writer, fileInfo)
		if err != nil {
			return "", err
		}
	} else if recursive {
		err = tarPackageRecursive(name, targetPath, writer)
		if err != nil {
			return "", err
		}
	} else {
		err = tarPackageNonRecursive(name, targetPath, writer)
		if err != nil {
			return "", err
		}
	}

	return layerFile.Name(), nil
}

func tarPackageNonRecursive(dirName, targetPath string, writer *tar.Writer) error {
	dir, err := os.Open(dirName)
	if err != nil {
		return err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		err := writeFileToTar(dir.Name()+string(filepath.Separator)+fileInfo.Name(), targetPath, writer, fileInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFileToTar(name, targetPath string, writer *tar.Writer, fileInfo fs.FileInfo) error {
	file, err := os.Open(name)
	if err != nil {
		return err
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
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return nil
}

func tarPackageRecursive(dirName, targetPath string, writer *tar.Writer) error {
	filepath.Walk(dirName, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			fileRelPath := strings.Replace(filePath, path.Clean(dirName), "", 1)

			// prepare the tar header
			header := new(tar.Header)
			header.Name = path.Join(targetPath, fileRelPath)
			header.Size = fileInfo.Size()
			header.Mode = int64(fileInfo.Mode())
			header.ModTime = fileInfo.ModTime()

			err = writer.WriteHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func appendPaths(base v1.Image, annotations map[string]string, paths ...string) (v1.Image, error) {
	additions := make([]mutate.Addendum, 0, len(paths))
	for idx, path := range paths {
		layer, err := tarball.LayerFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading tar %q: %v", path, err)
		}
		addendum := mutate.Addendum{
			Layer: layer,
		}
		if len(annotations) > 0 && idx == len(paths)-1 {
			addendum.Annotations = annotations
		}
		additions = append(additions, addendum)
	}

	return mutate.Append(base, additions...)
}
