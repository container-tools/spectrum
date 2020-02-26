package builder

import (
	"encoding/json"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/name"
)

// dirKeychain implements Keychain with standard docker config files in a given directory.
type dirKeychain struct {
	dir string
}

func NewDirKeyChain(dir string) authn.Keychain {
	return &dirKeychain{
		dir: dir,
	}
}

// Resolve implements Keychain.
func (fk *dirKeychain) Resolve(target authn.Resource) (authn.Authenticator, error) {
	cf, err := config.Load(fk.dir)
	if err != nil {
		return nil, err
	}

	// See:
	// https://github.com/google/ko/issues/90
	// https://github.com/moby/moby/blob/fc01c2b481097a6057bec3cd1ab2d7b4488c50c4/registry/config.go#L397-L404
	key := target.RegistryStr()
	if key == name.DefaultRegistry {
		key = authn.DefaultAuthKey
	}

	cfg, err := cf.GetAuthConfig(key)
	if err != nil {
		return nil, err
	}
	if logs.Enabled(logs.Debug) {
		b, err := json.Marshal(cfg)
		if err == nil {
			logs.Debug.Printf("dirKeychain.Resolve(%q) = %s", key, string(b))
		}
	}

	empty := types.AuthConfig{}
	if cfg == empty {
		return authn.Anonymous, nil
	}
	return authn.FromConfig(authn.AuthConfig{
		Username:      cfg.Username,
		Password:      cfg.Password,
		Auth:          cfg.Auth,
		IdentityToken: cfg.IdentityToken,
		RegistryToken: cfg.RegistryToken,
	}), nil
}
