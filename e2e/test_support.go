package e2e

import (
	"archive/tar"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/container-tools/spectrum/pkg/cmd"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/opencontainers/go-digest"
	"gotest.tools/assert"
)

func isRegistryInsecure() bool {
	return getRegistryInsecure() == "true"
}

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

func getImageEntrypoint(image string, insecure bool) []string {
	options := []crane.Option(nil)
	if insecure {
		options = append(options, crane.Insecure)
	}

	img, err := crane.Pull(image, options...)
	if err != nil {
		panic(err)
	}

	configFile, err := img.ConfigFile()
	if err != nil {
		panic(err)
	}
	return configFile.Config.Entrypoint
}

func getImageAnnotations(image string, insecure bool) map[string]string {
	options := []crane.Option(nil)
	if insecure {
		options = append(options, crane.Insecure)
	}

	img, err := crane.Pull(image, options...)
	if err != nil {
		panic(err)
	}

	manifest, err := img.Manifest()
	if err != nil {
		panic(err)
	}
	annotations := make(map[string]string)
	for _, layer := range manifest.Layers {
		for k, v := range layer.Annotations {
			annotations[k] = v
		}
	}
	return annotations
}

func assertDataMatch(t *testing.T, image string, insecure bool, dir, expected string) {
	options := []crane.Option(nil)
	if insecure {
		options = append(options, crane.Insecure)
	}

	img, err := crane.Pull(image, options...)
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

func createRegistryConfigDir(t *testing.T, registry, user, pass string) string {
	var params = struct {
		Token    string
		Registry string
	}{
		Token:    base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass))),
		Registry: registry,
	}
	templateFile, err := ioutil.ReadFile("./files/.docker/config.json.tmpl")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tmpl, err := template.New(".dockerconfigjson").Parse(string(templateFile))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tmpDir, err := ioutil.TempDir("", "spectrum-docker-config-")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	configFile, err := os.Create(filepath.Join(tmpDir, "config.json"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer configFile.Close()
	if err := tmpl.Execute(configFile, params); err != nil {
		t.Error(err)
		t.FailNow()
	}
	return tmpDir
}
