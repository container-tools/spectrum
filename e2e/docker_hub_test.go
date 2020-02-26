package e2e

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDockerHubPush(t *testing.T) {
	RegisterTestingT(t)

	user := os.Getenv("TEST_DOCKER_HUB_USERNAME")
	pass := os.Getenv("TEST_DOCKER_HUB_PASSWORD")
	if user == "" || pass == "" {
		t.Skip("No docker credentials found")
	}

	configDir := createRegistryConfigDir(t, "https://index.docker.io/v1/", user, pass)
	target := fmt.Sprintf("%s/spectrum-push", user)

	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-config-dir", configDir,
		"./files/01-simple:/app")).To(BeNil())

	assertDataMatch(t, target, false, "/app", "./files/01-simple")
}
