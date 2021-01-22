package e2e

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSimplePublish(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/simple"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	assertDataMatch(t, target, isRegistryInsecure(), "/app", "./files/01-simple")
}

func TestAnnotations(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/annotated"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"-a", "mykey=myval",
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	annotations := getImageAnnotations(target, isRegistryInsecure())
	assert.Equal(t, "myval", annotations["mykey"])
}

func TestAnnotationsMultipleLayers(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/annotated"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"-a", "mykey=myval",
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app", "./files/02-override:/app")).To(BeNil())

	annotations := getImageAnnotations(target, isRegistryInsecure())
	assert.Equal(t, "myval", annotations["mykey"])
}

func TestDirectOverride(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/override"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app", "./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, target, isRegistryInsecure(), "/app", "./files/03-merge")
}

func TestLayerComposition(t *testing.T) {
	RegisterTestingT(t)

	target1 := getRegistry() + "/publish/layer1"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target1,
		"--push-insecure="+getRegistryInsecure(),
		"./files/01-simple:/app")).To(BeNil())

	target2 := getRegistry() + "/publish/layer2"
	Expect(spectrum("build", "-b", getRegistry()+"/publish/layer1",
		"--pull-insecure="+getRegistryInsecure(),
		"-t", target2,
		"--push-insecure="+getRegistryInsecure(),
		"./files/02-override:/app")).To(BeNil())

	assertDataMatch(t, target2, isRegistryInsecure(), "/app", "./files/03-merge")
}

func TestRecursive(t *testing.T) {
	RegisterTestingT(t)

	target := getRegistry() + "/publish/simple"
	Expect(spectrum("build", "-b", "adoptopenjdk/openjdk8:slim",
		"-t", target,
		"--push-insecure="+getRegistryInsecure(),
		"-r",
		"./files/04-recursive:/app")).To(BeNil())

	assertDataMatch(t, target, isRegistryInsecure(), "/app", "./files/04-recursive")
}
