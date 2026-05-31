package cli

import (
	"fmt"
	"os"

	"github.com/HarshalPatel1972/spectra/internal/report"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <old-cbom.json> <new-cbom.json>",
	Short: "Compare two Spectra CBOM files to track migration progress",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldPath := args[0]
		newPath := args[1]

		diff, err := report.CompareCBOMs(oldPath, newPath)
		if err != nil {
			return fmt.Errorf("diff failed: %w", err)
		}

		report.RenderDiffTerminal(os.Stdout, diff)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
