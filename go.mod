module github.com/container-tools/spectrum

go 1.18

require (
	github.com/docker/cli v23.0.1+incompatible
	github.com/google/go-containerregistry v0.13.0
	github.com/onsi/gomega v1.27.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.2
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/containerd/stargz-snapshotter/estargz v0.14.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v23.0.1+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/image-spec v1.0.2-0.20211117181255-693428a734f5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Using a fork that removes the HTTPS ping before using HTTP for insecure registries
replace github.com/google/go-containerregistry => github.com/container-tools/go-containerregistry v0.7.1-0.20211124090132-40ccc94a466b
