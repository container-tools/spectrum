package e2e

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/container-tools/spectrum/pkg/cmd"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/opencontainers/go-digest"
	"gotest.tools/assert"
)

func getRegistryInsecure() string {
	ins := os.Getenv("SPECTRUM_REGISTRY_INSECURE")
	if ins != "" {
		return ins
	}
	return "true"
}

func getRegistry() string {
	reg := os.Getenv("SPECTRUM_REGISTRY")
	if reg != "" {
		return reg
	}
	return "localhost:5000"
}

func spectrum(args ...string) error {
	spectrum := cmd.Spectrum()
	spectrum.SetArgs(args)
	return spectrum.Execute()
}

func assertDataMatch(t *testing.T, image, dir, expected string) {
	options := []crane.Option(nil)
	if getRegistryInsecure() == "true" {
		options = append(options, crane.Insecure)
	}

	img, err := crane.Pull(getRegistry()+"/"+image, options...)
	if err != nil {
		panic(err)
	}

	tmp, err := ioutil.TempFile("", "spectrum-*.tar")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmp.Name())

	err = crane.Export(img, tmp)
	if err != nil {
		panic(err)
	}
	if err := tmp.Close(); err != nil {
		panic(err)
	}

	tmp, err = os.Open(tmp.Name())
	if err != nil {
		panic(err)
	}
	content := tar.NewReader(tmp)
	contentMap := make(map[string]string)
	for {
		header, err := content.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else if header == nil {
			continue
		}

		if strings.HasPrefix(header.Name, dir) {
			name := header.Name[len(dir):]
			hash, err := digest.FromReader(content)
			if err != nil {
				panic(err)
			}
			contentMap[name] = hash.String()
		}
	}

	expectedMap := make(map[string]string)
	absExpected, err := filepath.Abs(expected)
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(absExpected, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		name := path[len(absExpected):]
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		hash, err := digest.FromReader(file)
		if err != nil {
			panic(err)
		}
		expectedMap[name] = hash.String()
		return nil
	})
	if err != nil {
		panic(err)
	}

	assert.DeepEqual(t, expectedMap, contentMap)
}
