package e2e

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSimplePublish(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/simple"
	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	assertDataMatch(t, target, isRegistryInsecure(), "/app", "./files/01-simple")
}

func TestDirectOverride(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/override"
	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app", "./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, target, isRegistryInsecure(), "/app", "./files/03-merge")
}

func TestLayerComposition(t *testing.T) {
	RegisterTestingT(t)

	target1 := getRegistry() + "/publish/layer1"
	Expect(spectrum("-b", "adoptopenjdk/openjdk8:slim",
		"-t", target1,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	target2 := getRegistry() + "/publish/layer2"
	Expect(spectrum("-b", getRegistry()+"/publish/layer1",
		"--pull-insecure="+getRegistryInsecure(),
		"-t", target2,
		"--push-insecure="+getRegistryInsecure(),
		"./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, target2, isRegistryInsecure(), "/app", "./files/03-merge")
}
