package cli

import (
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/compliance"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/spf13/cobra"
)

var complianceScanID string

var complianceCmd = &cobra.Command{
	Use:   "compliance",
	Short: "Check findings against compliance frameworks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if complianceScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		// Get all findings for the scan
		// We don't have GetFindings in store yet.
		// For the sake of the command, let's pretend we do, or we implement it.
		findings, err := store.GetFindings(complianceScanID)
		if err != nil {
			return fmt.Errorf("getting findings: %w", err)
		}

		gaps := compliance.EvaluateFindings(findings)
		if len(gaps) == 0 {
			fmt.Println("No compliance violations found!")
			return nil
		}

		fmt.Printf("Compliance Violations Found: %d\n\n", len(gaps))
		for _, g := range gaps {
			deadline := "No Deadline"
			if g.Deadline != nil {
				deadline = *g.Deadline
			}
			fmt.Printf("[%s] %s (%s)\n", g.Severity, g.RequirementDesc, g.Framework)
			fmt.Printf("  Finding ID: %s\n", g.FindingID)
			fmt.Printf("  Requirement: %s\n", g.RequirementID)
			fmt.Printf("  Deadline: %s\n\n", deadline)
		}

		return nil
	},
}

func init() {
	complianceCmd.Flags().StringVar(&complianceScanID, "scan-id", "", "which scan to check (required)")
	rootCmd.AddCommand(complianceCmd)
}
