module github.com/container-tools/spectrum

go 1.13

require (
	github.com/docker/cli v20.10.11+incompatible
	github.com/google/go-containerregistry v0.0.0-20200220215334-221517453cf9
	github.com/onsi/gomega v1.10.3
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	gotest.tools v2.2.0+incompatible
)

// Using a fork that removes the HTTPS ping before using HTTP for insecure registries
replace github.com/google/go-containerregistry => github.com/container-tools/go-containerregistry v0.7.1-0.20211124090132-40ccc94a466b
