package cli

import (
	"fmt"
	"os"

	"github.com/HarshalPatel1972/spectra/internal/report"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/spf13/cobra"
)

var endpointURL string

var endpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Scan a live TLS endpoint for cryptographic assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if endpointURL == "" {
			return fmt.Errorf("--url is required")
		}

		fmt.Printf("Dialing %s...\n", endpointURL)

		findings, err := scanner.ScanTLS(endpointURL)
		if err != nil {
			return fmt.Errorf("scanning endpoint: %w", err)
		}

		result := &scanner.ScanResult{
			ScanRoot: endpointURL,
			Findings: findings,
		}

		report.RenderTerminal(os.Stdout, result, "v0.1.0")

		return nil
	},
}

func init() {
	endpointCmd.Flags().StringVar(&endpointURL, "url", "", "TLS endpoint to scan (e.g., example.com:443)")
	rootCmd.AddCommand(endpointCmd)
}
