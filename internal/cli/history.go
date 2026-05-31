package cli

import (
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/HarshalPatel1972/spectra/internal/temporal"
	"github.com/spf13/cobra"
)

var (
	historyBase    string
	historyCurrent string
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Compare scans to calculate drift, velocity, and blame",
	RunE: func(cmd *cobra.Command, args []string) error {
		if historyBase == "" || historyCurrent == "" {
			return fmt.Errorf("--base and --current are required")
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		report, err := temporal.CalculateDrift(store, historyBase, historyCurrent)
		if err != nil {
			return err
		}

		fmt.Printf("History & Drift Report\n")
		fmt.Printf("Base Scan:    %s\n", report.BaseScanID)
		fmt.Printf("Current Scan: %s\n", report.CurrentScanID)
		fmt.Printf("QRS Change:   %+d\n", report.QRSChange)

		if report.QRSChange < 0 {
			fmt.Printf("Velocity:     %.2f QRS reduction/day\n", report.Velocity)
			if report.PQCReadyDate != nil {
				fmt.Printf("Projected PQC Ready Date: %s\n", report.PQCReadyDate.Format("2006-01-02"))
			}
		} else if report.QRSChange > 0 {
			fmt.Printf("Velocity:     Risk is increasing!\n")
		} else {
			fmt.Printf("Velocity:     Stagnant (No change)\n")
		}

		return nil
	},
}

func init() {
	historyCmd.Flags().StringVar(&historyBase, "base", "", "base scan ID")
	historyCmd.Flags().StringVar(&historyCurrent, "current", "", "current scan ID")
	rootCmd.AddCommand(historyCmd)
}
