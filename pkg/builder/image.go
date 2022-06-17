package builder

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func Pull(options Options) (v1.Image, error) {
	if options.Base == "" || options.Base == "scratch" {
		return empty.Image, nil
	}
	nameOptions := makeNameOptions(options.PullInsecure)
	ref, err := name.ParseReference(options.Base, nameOptions...)
	if err != nil {
		return nil, fmt.Errorf("parsing tag %q: %v", options.Base, err)
	}

	remoteOptions := makeRemoteOptions(options, options.PullConfigDir)
	return remote.Image(ref, remoteOptions...)
}

func Push(img v1.Image, options Options) error {
	nameOptions := makeNameOptions(options.PushInsecure)
	tag, err := name.NewTag(options.Target, nameOptions...)
	if err != nil {
		return fmt.Errorf("parsing tag %q: %v", options.Target, err)
	}

	remoteOptions := makeRemoteOptions(options, options.PushConfigDir)
	return remote.Write(tag, img, remoteOptions...)
}

func makeNameOptions(insecure bool) (nameOptions []name.Option) {
	if insecure {
		nameOptions = append(nameOptions, name.Insecure)
	}
	return
}

func makeRemoteOptions(options Options, configDir string) (remoteOptions []remote.Option) {
	if options.Jobs > 0 {
		remoteOptions = append(remoteOptions, remote.WithJobs(options.Jobs))
	}
	if configDir != "" {
		keyChain := NewDirKeyChain(configDir)
		remoteOptions = append(remoteOptions, remote.WithAuthFromKeychain(keyChain))
	} else {
		remoteOptions = append(remoteOptions, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	}
	return
}
