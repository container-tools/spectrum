package e2e

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSimplePublish(t *testing.T) {
	RegisterTestingT(t)

	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", getRegistry()+"/publish/simple",
		"--target-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	assertDataMatch(t, "publish/simple", "/app", "./files/01-simple")
}

func TestDirectOverride(t *testing.T) {
	RegisterTestingT(t)

	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", getRegistry()+"/publish/override",
		"--target-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app", "./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, "publish/override", "/app", "./files/03-merge")
}

func TestLayerComposition(t *testing.T) {
	RegisterTestingT(t)

	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", getRegistry()+"/publish/layer1",
		"--target-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	Expect(spectrum("-b", getRegistry()+"/publish/layer1",
		"--base-insecure="+getRegistryInsecure(),
		"-t", getRegistry()+"/publish/layer2",
		"--target-insecure="+getRegistryInsecure(),
		"./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, "publish/layer2", "/app", "./files/03-merge")
}
