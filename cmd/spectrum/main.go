package main

import (
	"os"

	"github.com/container-tools/spectrum/pkg/cmd"
)

func main() {
	err := cmd.Spectrum().Execute()
	if err != nil {
		os.Exit(1)
	}
}
