package cli

import (
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/metrics"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/spf13/cobra"
)

var agilityScanID string

var agilityCmd = &cobra.Command{
	Use:   "agility",
	Short: "Calculate the Cryptographic Agility Index (CAI) of the codebase",
	RunE: func(cmd *cobra.Command, args []string) error {
		if agilityScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		findings, err := store.GetFindings(agilityScanID)
		if err != nil {
			return fmt.Errorf("getting findings: %w", err)
		}

		report := metrics.CalculateCAI(findings)

		fmt.Printf("Cryptographic Agility Index (CAI)\n")
		fmt.Printf("===================================\n")
		fmt.Printf("Score:             %d / 100\n", report.CAI)
		fmt.Printf("Abstraction:       %s\n", report.AbstractionGrade)
		fmt.Printf("Centralization:    %s\n", report.CentralizationGrade)
		fmt.Printf("Configurability:   %v\n", report.Configurability)
		
		if report.CAI >= 80 {
			fmt.Println("\nStatus: Highly Agile. You are well positioned for a PQC migration.")
		} else {
			fmt.Println("\nStatus: Monolithic Crypto. Significant refactoring is required before transitioning algorithms.")
		}

		return nil
	},
}

func init() {
	agilityCmd.Flags().StringVar(&agilityScanID, "scan-id", "", "which scan to measure (required)")
	rootCmd.AddCommand(agilityCmd)
}
