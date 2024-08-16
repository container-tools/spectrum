> NOTE: this project is deprecated and will be soon archived. Please, use any other tool like [go-containerregistry](https://github.com/google/go-containerregistry).

# Spectrum

Spectrum is a lightweight super-fast image builder tailored for application development use cases, where
you need to add application artifacts on top of a base image and do it as fast as possible.

Spectrum is able to build and push images to your registry in few seconds.

For example:

```
$ spectrum build -b adoptopenjdk/openjdk8:slim \
  -t local.dev/myorg/myapp \
  ./dist:/deployments
```

You need to specify the base image (`-b`), the target image (`-t`) and the directory of your file system that you want to copy
together with the location on the image file system (`/path/to/source-dir:/path/to/dest-dir`).

Additional options can be specified:

```
$ spectrum build --help
Build an image and publish it

Usage:
  spectrum build [flags]

Flags:
  -b, --base string              Base container image to use
  -h, --help                     help for build
      --pull-config-dir string   A directory containing the docker config.json file that will be used for pulling the base image, in case authentication is required
      --pull-insecure            If the base image is hosted in an insecure registry
      --push-config-dir string   A directory containing the docker config.json file that will be used for pushing the target image, in case authentication is required
      --push-insecure            If the target image will be pushed to an insecure registry
  -t, --target string            Target container image to use
```

## Credits

The Spectrum tool is based on [go-containerregistry](https://github.com/google/go-containerregistry).
