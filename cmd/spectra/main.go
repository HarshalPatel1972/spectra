// Package main is the entry point for the Spectra CLI.
package main

import (
	"fmt"
	"os"

	"github.com/HarshalPatel1972/spectra/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
