package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/container-tools/spectrum/pkg/builder"
	"github.com/container-tools/spectrum/pkg/util"
	"github.com/spf13/cobra"
)

type CommandOptions struct {
	builder.Options

	annotationList []string
	quiet          bool
}

func Spectrum() *cobra.Command {
	cmd := cobra.Command{
		Use:   "spectrum",
		Short: "Spectrum can publish simple container images in a few seconds",
	}

	options := CommandOptions{}
	build := cobra.Command{
		Use:   "build",
		Short: "Build an image and publish it",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("at least one argument is required")
			}
			for _, dir := range args {
				parts := strings.Split(dir, ":")
				if len(parts) != 2 {
					return errors.New("wrong format for dir " + dir + ". Expected: \"local:remote\"")
				}
			}

			// Configure output
			if !options.quiet {
				options.Stdout = cmd.OutOrStdout()
				options.Stderr = cmd.ErrOrStderr()
			}

			for _, akv := range options.annotationList {
				if options.Annotations == nil {
					options.Annotations = make(map[string]string)
				}
				parts := strings.SplitN(akv, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf(`wrong format for the annotation: expected "key=value", got %q`, akv)
				}
				options.Annotations[parts[0]] = parts[1]
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			digest, err := builder.Build(options.Options, args...)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), digest)
			return nil
		},
	}

	build.Flags().StringVarP(&options.Base, "base", "b", "", "Base container image to use")
	build.Flags().StringVarP(&options.Target, "target", "t", "", "Target container image to use")
	build.Flags().BoolVarP(&options.PullInsecure, "pull-insecure", "", false, "If the base image is hosted in an insecure registry")
	build.Flags().BoolVarP(&options.PushInsecure, "push-insecure", "", false, "If the target image will be pushed to an insecure registry")
	build.Flags().StringVarP(&options.PullConfigDir, "pull-config-dir", "", "", "A directory containing the docker config.json file that will be used for pulling the base image, in case authentication is required")
	build.Flags().StringVarP(&options.PushConfigDir, "push-config-dir", "", "", "A directory containing the docker config.json file that will be used for pushing the target image, in case authentication is required")
	build.Flags().StringSliceVarP(&options.annotationList, "annotations", "a", nil, "A list of annotations in the key=value format to add to the final image")
	build.Flags().BoolVarP(&options.quiet, "quiet", "q", false, "Do not print logs to stdout and stderr")
	build.Flags().BoolVarP(&options.Recursive, "recursive", "r", false, "Copy content from the source filesystem directory recursively")
	build.Flags().BoolVar(&options.ClearEntrypoint, "clear-entrypoint", false, "Clear any entrypoint defined")
	build.Flags().StringVar(&options.RunAs, "run-as", "", "User id/name used to run the container image")
	cmd.AddCommand(&build)

	version := cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), util.Version)
		},
	}
	cmd.AddCommand(&version)

	return &cmd
}
