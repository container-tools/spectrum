package builder

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

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
		localPath, targetPath, err := getPaths(spec, runtime.GOOS)
		if err != nil {
			return "", err
		}

		tarFile, err := tarPackage(localPath, targetPath, options.Recursive)
		if err != nil {
			return "", errors.Wrapf(err, "cannot package dir %s as tar file", localPath)
		}
		defer os.Remove(tarFile)
		tarFiles = append(tarFiles, tarFile)
	}
	newImage, err := appendPaths(base, options.Annotations, tarFiles...)
	if err != nil {
		return "", errors.Wrap(err, "could not append tar layers to base image")
	}
	confFile, err := newImage.ConfigFile()
	if err != nil {
		panic(err)
	}

	if options.ClearEntrypoint == true {
		StepLogger.Println("Clearing entrypoint...")
		// Change the entry point
		confFile.Config.Entrypoint = nil
		newImage, err = mutate.Config(newImage, confFile.Config)
		if err != nil {
			panic(err)
		}
	}

	if options.RunAs != "" {
		StepLogger.Printf("Setting user as %s", options.RunAs)
		// Change the configured USER
		confFile.Config.User = options.RunAs
		newImage, err = mutate.Config(newImage, confFile.Config)
		if err != nil {
			panic(err)
		}
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

func getPaths(paths string, os string) (localPath string, targetPath string, err error) {
	parts := strings.Split(paths, ":")
	if len(parts) != 2 && (len(parts) == 3 && os != "windows") {
		return "", "", errors.New("wrong dir format for " + paths + " (expected \"local:remote\")")
	}
	localPath = parts[0]
	targetPath = parts[1]
	if os == "windows" && len(parts) == 3 {
		localPath = fmt.Sprintf("%s:%s", parts[0], parts[1])
		targetPath = parts[2]
	}
	return localPath, targetPath, nil
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

	header := prepareHeader(
		targetPath,
		path.Join(targetPath, filepath.Base(file.Name())),
		fileInfo,
	)

	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareHeader(tp, name string, fi fs.FileInfo) *tar.Header {
	// prepare the tar header
	header := new(tar.Header)
	header.Name = name
	header.Size = fi.Size()
	header.Mode = int64(fi.Mode().Perm())
	// Non portable way of retrieving uid/gid, but Golang does not offer any other way programmatically
	fileSys := fi.Sys()
	if fileSys != nil {
		header.Uid = int(fileSys.(*syscall.Stat_t).Uid)
		header.Gid = int(fileSys.(*syscall.Stat_t).Gid)
	} else {
		StepLogger.Printf("Warning: could not read UID/GID. Assuming default (root) permissions.")
	}
	header.ModTime = fi.ModTime()

	return header
}

func tarPackageRecursive(dirName, targetPath string, writer *tar.Writer) error {

	filepath.Walk(dirName, func(filePath string, fileInfo os.FileInfo, err error) error {
		fileRelPath := strings.Replace(filePath, path.Clean(dirName), "", 1)
		header := prepareHeader(
			targetPath,
			path.Join(targetPath, fileRelPath),
			fileInfo,
		)
		if fileInfo.IsDir() {
			header.Name = header.Name + "/"
		}

		err = writer.WriteHeader(header)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
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
