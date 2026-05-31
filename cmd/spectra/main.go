// Package main is the entry point for the Spectra CLI.
package main

import (
	"os"

	"github.com/HarshalPatel1972/spectra/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
