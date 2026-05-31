package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// htmlOutputFilename is the default file name for the HTML report.
const htmlOutputFilename = "spectra-report.html"

// hexColor maps risk bands to hex colour strings for the HTML template.
func hexColor(band detector.RiskBand) string {
	switch band {
	case detector.BandCritical:
		return "#FF4444"
	case detector.BandHigh:
		return "#FF8C00"
	case detector.BandMedium:
		return "#FFD700"
	case detector.BandLow:
		return "#4ADE80"
	case detector.BandSafe:
		return "#38BDF8"
	default:
		return "#F8FAFC"
	}
}

// htmlData holds the preprocessed data passed to the HTML template.
type htmlData struct {
	Version       string
	ScanRoot      string
	StartedAt     string
	CompletedAt   string
	Duration      string
	TotalFiles    int
	ScannedFiles  int
	SkippedFiles  int
	AggregateQRS  int
	AggBand       string
	AggColor      string
	BandCounts    []bandCount
	Findings      []htmlFinding
	ActionPlan    []scanner.ActionItem
	AlgoInventory []algoEntry
	TotalFindings int
	GeneratedAt   string
}

// bandCount holds a band label, count, and colour for template rendering.
type bandCount struct {
	Band       string
	Count      int
	Color      string
	Percentage float64
}

// htmlFinding is a template-friendly representation of a Finding.
type htmlFinding struct {
	Algorithm  string
	Source     string
	FilePath   string
	LineNumber int
	KeySize    string
	QRS        int
	RiskBand   string
	BandColor  string
	Effort     string
}

// algoEntry is a unique algorithm summary row for the inventory table.
type algoEntry struct {
	Algorithm   string
	DisplayName string
	Family      string
	Occurrences int
	MaxQRS      int
	RiskBand    string
	BandColor   string
	PQCSafe     bool
	Replacement string
}

// WriteHTML generates a self-contained HTML report of scan results and writes
// it to outDir. The report uses inline CSS with a dark theme suitable for
// executive presentation. Returns the path of the written file.
func WriteHTML(result *scanner.ScanResult, outDir string, version string) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory %s: %w", outDir, err)
	}

	data := buildHTMLData(result, version)

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"bandColorFn": hexColor,
	}).Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing HTML template: %w", err)
	}

	outPath := filepath.Join(outDir, htmlOutputFilename)
	f, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("creating HTML file %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("executing HTML template: %w", err)
	}

	return outPath, nil
}

// buildHTMLData preprocesses the ScanResult into the structure consumed by the
// HTML template.
func buildHTMLData(result *scanner.ScanResult, version string) htmlData {
	bands := []detector.RiskBand{
		detector.BandCritical, detector.BandHigh, detector.BandMedium,
		detector.BandLow, detector.BandSafe,
	}

	totalFindings := len(result.Findings)
	var counts []bandCount
	for _, b := range bands {
		c := result.FindingsByBand[b]
		pct := 0.0
		if totalFindings > 0 {
			pct = float64(c) / float64(totalFindings) * 100
		}
		counts = append(counts, bandCount{
			Band:       string(b),
			Count:      c,
			Color:      hexColor(b),
			Percentage: pct,
		})
	}

	// Sort findings by QRS desc for display.
	sorted := make([]scanner.Finding, len(result.Findings))
	copy(sorted, result.Findings)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].QRS > sorted[j].QRS
	})

	var findings []htmlFinding
	for _, f := range sorted {
		ks := "-"
		if f.KeySize > 0 {
			ks = fmt.Sprintf("%d", f.KeySize)
		}
		findings = append(findings, htmlFinding{
			Algorithm:  f.Algorithm,
			Source:     string(f.Source),
			FilePath:   f.FilePath,
			LineNumber: f.LineNumber,
			KeySize:    ks,
			QRS:        f.QRS,
			RiskBand:   string(f.RiskBand),
			BandColor:  hexColor(f.RiskBand),
			Effort:     string(f.MigrationEffort),
		})
	}

	// Algorithm inventory: unique algorithms with aggregate stats.
	algoMap := make(map[string]*algoEntry)
	for _, f := range result.Findings {
		key := f.Algorithm
		if entry, ok := algoMap[key]; ok {
			entry.Occurrences++
			if f.QRS > entry.MaxQRS {
				entry.MaxQRS = f.QRS
				entry.RiskBand = string(f.RiskBand)
				entry.BandColor = hexColor(f.RiskBand)
			}
		} else {
			replacement := "-"
			if len(f.AlgorithmInfo.PQCReplacement) > 0 {
				replacement = strings.Join(f.AlgorithmInfo.PQCReplacement, ", ")
			}
			algoMap[key] = &algoEntry{
				Algorithm:   f.Algorithm,
				DisplayName: f.AlgorithmInfo.DisplayName,
				Family:      f.AlgorithmInfo.Family,
				Occurrences: 1,
				MaxQRS:      f.QRS,
				RiskBand:    string(f.RiskBand),
				BandColor:   hexColor(f.RiskBand),
				PQCSafe:     f.AlgorithmInfo.PQCSafe,
				Replacement: replacement,
			}
		}
	}
	var inventory []algoEntry
	for _, e := range algoMap {
		inventory = append(inventory, *e)
	}
	sort.Slice(inventory, func(i, j int) bool {
		return inventory[i].MaxQRS > inventory[j].MaxQRS
	})

	aggBand := bandFromQRS(result.AggregateQRS)
	duration := result.CompletedAt.Sub(result.StartedAt).Round(time.Millisecond)

	return htmlData{
		Version:       version,
		ScanRoot:      result.ScanRoot,
		StartedAt:     result.StartedAt.Format(time.RFC3339),
		CompletedAt:   result.CompletedAt.Format(time.RFC3339),
		Duration:      duration.String(),
		TotalFiles:    result.TotalFiles,
		ScannedFiles:  result.ScannedFiles,
		SkippedFiles:  result.SkippedFiles,
		AggregateQRS:  result.AggregateQRS,
		AggBand:       string(aggBand),
		AggColor:      hexColor(aggBand),
		BandCounts:    counts,
		Findings:      findings,
		ActionPlan:    result.ActionPlan,
		AlgoInventory: inventory,
		TotalFindings: totalFindings,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// htmlTemplate is the self-contained HTML template for the Spectra report.
// All CSS is inline — no external dependencies or CDN calls.
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Spectra — Cryptographic Risk Report</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif;
    background: #0f172a; color: #f8fafc; line-height: 1.6;
  }
  .container { max-width: 1280px; margin: 0 auto; padding: 2rem; }
  h1, h2, h3 { font-weight: 700; }

  /* Header */
  .header {
    text-align: center; padding: 3rem 0 2rem;
    border-bottom: 1px solid #334155;
  }
  .header h1 { font-size: 3rem; letter-spacing: 0.4em; margin-bottom: 0.5rem; }
  .header .subtitle { color: #94a3b8; font-size: 0.875rem; }
  .meta-grid {
    display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem; margin-top: 1.5rem; text-align: left;
  }
  .meta-item { background: #1e293b; border-radius: 0.5rem; padding: 0.75rem 1rem; }
  .meta-item .label { color: #94a3b8; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em; }
  .meta-item .value { font-size: 1rem; font-weight: 600; margin-top: 0.25rem; }

  /* Cards */
  .card {
    background: #1e293b; border-radius: 0.75rem;
    padding: 1.5rem; margin-top: 2rem;
  }
  .card h2 { font-size: 1.25rem; margin-bottom: 1rem; color: #e2e8f0; }

  /* Executive summary */
  .exec-summary {
    display: flex; flex-wrap: wrap; gap: 2rem; align-items: center;
    justify-content: center;
  }
  .qrs-circle {
    width: 160px; height: 160px; border-radius: 50%;
    display: flex; flex-direction: column; align-items: center; justify-content: center;
    border: 4px solid; font-weight: 700;
  }
  .qrs-circle .number { font-size: 3rem; line-height: 1; }
  .qrs-circle .label { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; margin-top: 0.25rem; }

  .badge-row { display: flex; flex-wrap: wrap; gap: 0.75rem; }
  .badge {
    display: inline-flex; align-items: center; gap: 0.4rem;
    padding: 0.5rem 1rem; border-radius: 9999px; font-weight: 600; font-size: 0.875rem;
    border: 1px solid; background: rgba(0,0,0,0.2);
  }
  .badge .dot {
    width: 10px; height: 10px; border-radius: 50%; display: inline-block;
  }

  /* Risk distribution bar */
  .risk-bar-container { margin-top: 1rem; }
  .risk-bar {
    display: flex; height: 32px; border-radius: 6px; overflow: hidden; width: 100%;
  }
  .risk-bar-segment {
    display: flex; align-items: center; justify-content: center;
    font-size: 0.75rem; font-weight: 600; color: #0f172a;
    min-width: 0; transition: width 0.3s;
  }
  .risk-bar-legend {
    display: flex; flex-wrap: wrap; gap: 1rem; margin-top: 0.75rem; font-size: 0.8rem;
  }
  .risk-bar-legend .item { display: flex; align-items: center; gap: 0.3rem; }
  .risk-bar-legend .swatch { width: 12px; height: 12px; border-radius: 3px; display: inline-block; }

  /* Tables */
  table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
  th {
    text-align: left; padding: 0.6rem 0.75rem; color: #94a3b8;
    text-transform: uppercase; font-size: 0.7rem; letter-spacing: 0.05em;
    border-bottom: 1px solid #334155;
  }
  td { padding: 0.55rem 0.75rem; border-bottom: 1px solid #1e293b; }
  tr:hover td { background: rgba(255,255,255,0.03); }
  .mono { font-family: "SF Mono", "Fira Code", "Cascadia Code", Consolas, monospace; font-size: 0.8rem; }
  .band-badge {
    display: inline-block; padding: 0.15rem 0.6rem; border-radius: 4px;
    font-weight: 700; font-size: 0.7rem; text-transform: uppercase;
  }

  /* Action plan */
  .action-item {
    padding: 1rem; border-left: 4px solid; margin-bottom: 0.75rem;
    background: rgba(0,0,0,0.15); border-radius: 0 0.5rem 0.5rem 0;
  }
  .action-item .title { font-weight: 700; font-size: 1rem; }
  .action-item .detail { color: #94a3b8; font-size: 0.85rem; margin-top: 0.25rem; }
  .action-item .replacement { color: #38bdf8; font-size: 0.85rem; margin-top: 0.25rem; }

  /* Footer */
  .footer {
    text-align: center; color: #475569; font-size: 0.75rem;
    margin-top: 3rem; padding-top: 1.5rem; border-top: 1px solid #1e293b;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .container { padding: 1rem; }
    .header h1 { font-size: 2rem; }
    .meta-grid { grid-template-columns: 1fr; }
    .exec-summary { flex-direction: column; }
  }
</style>
</head>
<body>
<div class="container">

  <!-- Header -->
  <div class="header">
    <h1>SPECTRA</h1>
    <div class="subtitle">Post-Quantum Cryptographic Risk Assessment Report</div>
    <div class="meta-grid">
      <div class="meta-item">
        <div class="label">Scan Root</div>
        <div class="value mono">{{.ScanRoot}}</div>
      </div>
      <div class="meta-item">
        <div class="label">Duration</div>
        <div class="value">{{.Duration}}</div>
      </div>
      <div class="meta-item">
        <div class="label">Files Scanned</div>
        <div class="value">{{.ScannedFiles}} / {{.TotalFiles}}</div>
      </div>
      <div class="meta-item">
        <div class="label">Spectra Version</div>
        <div class="value">{{.Version}}</div>
      </div>
    </div>
  </div>

  <!-- Executive Summary -->
  <div class="card">
    <h2>Executive Summary</h2>
    <div class="exec-summary">
      <div class="qrs-circle" style="border-color: {{.AggColor}}; color: {{.AggColor}};">
        <div class="number">{{.AggregateQRS}}</div>
        <div class="label">QRS / 100</div>
      </div>
      <div>
        <div style="margin-bottom:0.75rem; font-size:1.1rem;">
          Risk Band: <strong style="color:{{.AggColor}}">{{.AggBand}}</strong>
        </div>
        <div class="badge-row">
          {{range .BandCounts}}
          <span class="badge" style="border-color:{{.Color}}; color:{{.Color}};">
            <span class="dot" style="background:{{.Color}};"></span>
            {{.Band}}: {{.Count}}
          </span>
          {{end}}
        </div>
      </div>
    </div>
  </div>

  <!-- Risk Distribution -->
  <div class="card">
    <h2>Risk Distribution</h2>
    <div class="risk-bar-container">
      <div class="risk-bar">
        {{range .BandCounts}}{{if gt .Count 0}}
        <div class="risk-bar-segment" style="width:{{printf "%.1f" .Percentage}}%; background:{{.Color}};">
          {{if gt .Percentage 5.0}}{{.Count}}{{end}}
        </div>
        {{end}}{{end}}
      </div>
      <div class="risk-bar-legend">
        {{range .BandCounts}}
        <div class="item">
          <span class="swatch" style="background:{{.Color}};"></span>
          {{.Band}} ({{.Count}})
        </div>
        {{end}}
      </div>
    </div>
  </div>

  <!-- Findings Table -->
  <div class="card">
    <h2>Findings ({{.TotalFindings}})</h2>
    {{if .Findings}}
    <div style="overflow-x:auto;">
    <table>
      <thead>
        <tr>
          <th>Algorithm</th><th>Source</th><th>File</th><th>Line</th>
          <th>Key Size</th><th>QRS</th><th>Band</th><th>Effort</th>
        </tr>
      </thead>
      <tbody>
        {{range .Findings}}
        <tr>
          <td class="mono">{{.Algorithm}}</td>
          <td>{{.Source}}</td>
          <td class="mono" title="{{.FilePath}}">{{.FilePath}}</td>
          <td>{{.LineNumber}}</td>
          <td>{{.KeySize}}</td>
          <td><strong>{{.QRS}}</strong></td>
          <td><span class="band-badge" style="background:{{.BandColor}}; color:#0f172a;">{{.RiskBand}}</span></td>
          <td>{{.Effort}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
    </div>
    {{else}}
    <p style="color:#94a3b8;">No cryptographic findings detected.</p>
    {{end}}
  </div>

  <!-- Action Plan -->
  {{if .ActionPlan}}
  <div class="card">
    <h2>Remediation Action Plan</h2>
    {{range .ActionPlan}}
    <div class="action-item" style="border-left-color:{{bandColorFn .RiskBand}};">
      <div class="title" style="color:{{bandColorFn .RiskBand}};">
        #{{.Rank}} — {{.Algorithm}}
      </div>
      <div class="detail">
        {{.Recommendation}} &middot; {{.Occurrences}} occurrence{{if ne .Occurrences 1}}s{{end}} &middot; Effort: {{.Effort}} &middot; Priority Score: {{printf "%.1f" .PriorityScore}}
      </div>
      {{if .Replacement}}
      <div class="replacement">→ Replace with: {{.Replacement}}</div>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}

  <!-- Algorithm Inventory -->
  {{if .AlgoInventory}}
  <div class="card">
    <h2>Algorithm Inventory</h2>
    <div style="overflow-x:auto;">
    <table>
      <thead>
        <tr>
          <th>Algorithm</th><th>Family</th><th>Occurrences</th>
          <th>Max QRS</th><th>Band</th><th>PQC-Safe</th><th>Replacement</th>
        </tr>
      </thead>
      <tbody>
        {{range .AlgoInventory}}
        <tr>
          <td class="mono" title="{{.DisplayName}}">{{.Algorithm}}</td>
          <td>{{.Family}}</td>
          <td>{{.Occurrences}}</td>
          <td><strong>{{.MaxQRS}}</strong></td>
          <td><span class="band-badge" style="background:{{.BandColor}}; color:#0f172a;">{{.RiskBand}}</span></td>
          <td>{{if .PQCSafe}}✓{{else}}✗{{end}}</td>
          <td class="mono">{{.Replacement}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>
    </div>
  </div>
  {{end}}

  <!-- Footer -->
  <div class="footer">
    Generated by Spectra {{.Version}} on {{.GeneratedAt}}<br>
    Post-Quantum Cryptographic Bill of Materials &amp; Risk Assessment
  </div>

</div>
</body>
</html>`
