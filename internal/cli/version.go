package cli

import (
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd prints the build version, commit, and Go runtime information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and build information",
	Long:  "Display the Spectra version, git commit, build date, and Go runtime version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), version.Info())
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
