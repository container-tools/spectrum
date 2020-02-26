package cmd

import (
	"errors"
	"strings"

	"github.com/container-tools/spectrum/pkg/spectrum"
	"github.com/spf13/cobra"
)

func Spectrum() *cobra.Command {
	options := spectrum.Options{}
	cmd := cobra.Command{
		Use:   "spectrum",
		Short: "Spectrum can publish simple container images in a few of seconds",
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
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return spectrum.Build(options, args...)
		},
	}

	cmd.Flags().StringVarP(&options.Base, "base", "b", "", "Base container image to use")
	cmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target container image to use")
	cmd.Flags().BoolVarP(&options.PullInsecure, "pull-insecure", "", false, "If the base image is hosted in an insecure registry")
	cmd.Flags().BoolVarP(&options.PushInsecure, "push-insecure", "", false, "If the target image will be pushed to an insecure registry")
	cmd.Flags().StringVarP(&options.PullConfigDir, "pull-config-dir", "", "", "A directory containing the docker config.json file that will be used for pulling the base image, in case authentication is required")
	cmd.Flags().StringVarP(&options.PushConfigDir, "push-config-dir", "", "", "A directory containing the docker config.json file that will be used for pushing the target image, in case authentication is required")

	return &cmd
}
