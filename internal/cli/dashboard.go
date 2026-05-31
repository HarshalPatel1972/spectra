package cli

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/HarshalPatel1972/spectra/web"
	"github.com/spf13/cobra"
)

var dashboardPort int
var dashboardOutDir string

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start the Spectra Web Dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Serve embedded files
		dist, err := fs.Sub(web.DistFS, "dist")
		if err != nil {
			return err
		}

		http.Handle("/", http.FileServer(http.FS(dist)))

		// API endpoint to serve findings
		http.HandleFunc("/api/findings", func(w http.ResponseWriter, r *http.Request) {
			findingsPath := filepath.Join(dashboardOutDir, "spectra-findings.json")
			
			// If missing, check if it's in the current dir
			if _, err := os.Stat(findingsPath); os.IsNotExist(err) {
				findingsPath = "spectra-findings.json"
			}

			data, err := os.ReadFile(findingsPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading %s. Ensure you have run a scan with JSON output.", findingsPath), http.StatusNotFound)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		})

		addr := fmt.Sprintf(":%d", dashboardPort)
		fmt.Printf("Starting Spectra Web Dashboard at http://localhost%s ...\n", addr)
		fmt.Printf("Press Ctrl+C to stop.\n")
		
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	dashboardCmd.Flags().IntVarP(&dashboardPort, "port", "p", 8080, "Port to serve the dashboard on")
	dashboardCmd.Flags().StringVar(&dashboardOutDir, "out-dir", "./spectra-out", "Directory where spectra-findings.json is located")
	rootCmd.AddCommand(dashboardCmd)
}
