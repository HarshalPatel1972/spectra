// Package cli implements the Spectra command-line interface using cobra.
package cli

import (
	"github.com/spf13/cobra"
)

// Persistent flag variables shared across commands.
var (
	cfgFile    string
	verbose    bool
	quiet      bool
	noProgress bool
)

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "spectra",
	Short: "Cryptographic asset discovery & quantum risk intelligence",
	Long: `Spectra scans your codebase, certificates, dependencies, and configuration
files to discover every cryptographic asset and assess its quantum risk.

It produces actionable intelligence for PQC (Post-Quantum Cryptography)
migration planning, including:

  • Quantum Risk Scores (QRS) for each finding
  • NIST-recommended PQC replacement suggestions
  • CBOM (Cryptography Bill of Materials) export
  • CI/CD gate integration with configurable fail thresholds

Run 'spectra scan' to analyse a project or 'spectra version' for build info.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: .spectra.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&noProgress, "no-progress", false, "disable progress indicators")
}

// Execute runs the root command and returns any error.
func Execute() error {
	return rootCmd.Execute()
}
