package cli

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/HarshalPatel1972/spectra/internal/metrics"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/spf13/cobra"
)

var reportScanID string
var reportOut string

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report generation commands",
}

var reportExecutiveCmd = &cobra.Command{
	Use:   "executive",
	Short: "Generate an executive HTML risk report",
	RunE: func(cmd *cobra.Command, args []string) error {
		if reportScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}
		if reportOut == "" {
			reportOut = "executive_report.html"
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		scan, err := store.GetScan(reportScanID)
		if err != nil {
			return fmt.Errorf("getting scan: %w", err)
		}

		findings, err := store.GetFindings(reportScanID)
		if err != nil {
			return fmt.Errorf("getting findings: %w", err)
		}

		agility := metrics.CalculateCAI(findings)
		posture := metrics.CalculateCPS(findings, scan.AggregateQRS)

		tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>Spectra Executive Risk Report</title>
	<style>
		body { font-family: -apple-system, sans-serif; margin: 40px; background: #f9f9fb; color: #333; }
		.container { max-width: 900px; margin: 0 auto; background: white; padding: 40px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
		h1 { color: #1a1a2e; border-bottom: 2px solid #eaeaea; padding-bottom: 10px; }
		.metrics { display: flex; gap: 20px; margin-top: 30px; }
		.metric-card { flex: 1; padding: 20px; border-radius: 8px; color: white; text-align: center; }
		.cps-card { background: linear-gradient(135deg, #e53935, #b71c1c); }
		.cai-card { background: linear-gradient(135deg, #1e88e5, #0d47a1); }
		.value { font-size: 48px; font-weight: bold; margin: 10px 0; }
		.label { font-size: 14px; text-transform: uppercase; letter-spacing: 1px; opacity: 0.9; }
		table { width: 100%; border-collapse: collapse; margin-top: 30px; }
		th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
		th { background-color: #f1f1f1; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Cryptographic Executive Risk Report</h1>
		<p>Scan ID: {{ .ScanID }}</p>
		
		<div class="metrics">
			<div class="metric-card cps-card">
				<div class="label">Cryptographic Posture Score (CPS)</div>
				<div class="value">{{ .CPS }} / 100</div>
				<div>Grade: {{ .CPSGrade }} | Status: {{ .CPSStatus }}</div>
			</div>
			<div class="metric-card cai-card">
				<div class="label">Cryptographic Agility Index (CAI)</div>
				<div class="value">{{ .CAI }} / 100</div>
				<div>Abstraction: {{ .Abstraction }} | Centralization: {{ .Centralization }}</div>
			</div>
		</div>

		<h2>Summary</h2>
		<p>The codebase contains <strong>{{ .TotalFindings }}</strong> cryptographic findings across <strong>{{ .ScannedFiles }}</strong> scanned files.</p>
		
		<h2>Key Vulnerabilities</h2>
		<table>
			<tr>
				<th>Algorithm</th>
				<th>Occurrences</th>
				<th>Evidence / Citation</th>
			</tr>
			{{ range .TopFindings }}
			<tr>
				<td>{{ .Algorithm }}</td>
				<td>{{ .Count }}</td>
				<td>
					{{ if .Evidence }}
						<a href="{{ .Evidence.URL }}" target="_blank">{{ .Evidence.Title }}</a>
					{{ else }}
						N/A
					{{ end }}
				</td>
			</tr>
			{{ end }}
		</table>
	</div>
</body>
</html>
`

		// Aggregate top findings for the report
		type TopFinding struct {
			Algorithm string
			Count     int
			Evidence  *metrics.EvidenceItem
		}
		
		counts := make(map[string]int)
		for _, f := range findings {
			counts[f.Algorithm]++
		}

		var topFindings []TopFinding
		for algo, count := range counts {
			var ev *metrics.EvidenceItem
			evList := metrics.GetEvidence(algo)
			if len(evList) > 0 {
				ev = &evList[0]
			}
			topFindings = append(topFindings, TopFinding{
				Algorithm: algo,
				Count:     count,
				Evidence:  ev,
			})
		}

		t, err := template.New("report").Parse(tmpl)
		if err != nil {
			return err
		}

		f, err := os.Create(reportOut)
		if err != nil {
			return err
		}
		defer f.Close()

		data := struct {
			ScanID         string
			CPS            int
			CPSGrade       string
			CPSStatus      string
			CAI            int
			Abstraction    string
			Centralization string
			TotalFindings  int
			ScannedFiles   int
			TopFindings    []TopFinding
		}{
			ScanID:         reportScanID,
			CPS:            posture.CPS,
			CPSGrade:       posture.Grade,
			CPSStatus:      posture.RiskStatus,
			CAI:            agility.CAI,
			Abstraction:    agility.AbstractionGrade,
			Centralization: agility.CentralizationGrade,
			TotalFindings:  len(findings),
			ScannedFiles:   scan.ScannedFiles,
			TopFindings:    topFindings,
		}

		if err := t.Execute(f, data); err != nil {
			return err
		}

		absPath, _ := filepath.Abs(reportOut)
		fmt.Printf("Executive report generated at: %s\n", absPath)
		return nil
	},
}

func init() {
	reportExecutiveCmd.Flags().StringVar(&reportScanID, "scan-id", "", "scan ID to report on")
	reportExecutiveCmd.Flags().StringVar(&reportOut, "out", "executive_report.html", "output HTML file path")
	reportCmd.AddCommand(reportExecutiveCmd)
	rootCmd.AddCommand(reportCmd)
}
