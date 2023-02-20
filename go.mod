module github.com/container-tools/spectrum

go 1.16

require (
	github.com/docker/cli v23.0.1+incompatible
	github.com/google/go-containerregistry v0.13.0
	github.com/onsi/gomega v1.27.1
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.1
	gotest.tools v2.2.0+incompatible
)

// Using a fork that removes the HTTPS ping before using HTTP for insecure registries
replace github.com/google/go-containerregistry => github.com/container-tools/go-containerregistry v0.7.1-0.20211124090132-40ccc94a466b
