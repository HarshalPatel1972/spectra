package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
)

var badgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Generate a markdown badge for the latest scan",
	Long: `Reads the latest scan result from the local database and generates a Markdown badge snippet
that you can embed in your README.md. The badge uses the spectra.tools API to render the Quantum Risk Score.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %v", err)
		}
		
		dbPath := filepath.Join(homeDir, ".spectra", "state.db")
		store, err := persistence.NewStore(dbPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %v. Have you run a scan yet?", err)
		}
		defer store.Close()

		latestScan, err := store.GetLatestScan()
		if err != nil {
			return fmt.Errorf("failed to get latest scan: %v. Run 'spectra scan' first.", err)
		}
		
		markdown := fmt.Sprintf("[![Spectra QRS](https://spectra.tools/api/badge?qrs=%d)](https://spectra.tools)", latestScan.AggregateQRS)
		
		if !quiet {
			fmt.Println("Add the following markdown to your README.md to display your Quantum Risk Score:")
			fmt.Println()
		}
		
		fmt.Println(markdown)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(badgeCmd)
}
