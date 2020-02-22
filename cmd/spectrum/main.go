package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nicolaferraro/spectrum/pkg/spectrum"
	"github.com/spf13/cobra"
)

func main() {
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
					return errors.New("wrong format for dir " + dir + ". Expeced: \"local:remote\"")
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return spectrum.Build(options, args...)
		},
	}

	cmd.Flags().StringVarP(&options.Base, "base", "b", "", "Base container image to use")
	cmd.Flags().BoolVarP(&options.BaseInsecure, "base-insecure", "", false, "If the base image is hosted in an insecure registry")
	cmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target container image to use")
	cmd.Flags().BoolVarP(&options.TargetInsecure, "target-insecure", "", false, "If the target image will be pushed to an insecure registry")

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
