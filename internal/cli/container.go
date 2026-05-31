package cli

import (
	"fmt"
	"os"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/report"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/HarshalPatel1972/spectra/rules"
	"github.com/spf13/cobra"
)

var containerImage string

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Container image commands",
}

var containerScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan an OCI container image layer-by-layer",
	RunE: func(cmd *cobra.Command, args []string) error {
		if containerImage == "" {
			return fmt.Errorf("--image is required")
		}

		fmt.Printf("Pulling and scanning %s...\n", containerImage)

		registry, err := detector.LoadPatternsFromBytes(rules.CryptoPatternsYAML)
		if err != nil {
			return fmt.Errorf("failed to load patterns: %w", err)
		}

		findings, err := scanner.ScanImage(containerImage, registry, concurrency)
		if err != nil {
			return fmt.Errorf("scanning container: %w", err)
		}

		// Group by layer for output or just treat as a single scan
		result := &scanner.ScanResult{
			ScanRoot: containerImage,
			Findings: findings,
		}

		// We could enrich the result but for now we just dump to terminal
		report.RenderTerminal(os.Stdout, result, "v0.1.0")

		return nil
	},
}

func init() {
	containerScanCmd.Flags().StringVar(&containerImage, "image", "", "container image to scan (e.g., ubuntu:latest)")
	containerScanCmd.Flags().IntVar(&concurrency, "concurrency", 0, "number of parallel scanner goroutines (0 = auto)")
	containerCmd.AddCommand(containerScanCmd)
	rootCmd.AddCommand(containerCmd)
}
