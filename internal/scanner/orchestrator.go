// Package scanner orchestrates parallel file scanning, collects cryptographic
// findings, and produces aggregate results for downstream reporting.
package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/HarshalPatel1972/spectra/internal/detector"
)

// Finding represents a single cryptographic usage discovered during a scan.
type Finding struct {
	ID              string                 `json:"id"`
	Algorithm       string                 `json:"algorithm"`
	AlgorithmInfo   detector.AlgorithmInfo `json:"algorithm_info"`
	Source          detector.SourceType    `json:"source"`
	FilePath        string                 `json:"file_path"`
	LineNumber      int                    `json:"line_number"`
	LineContent     string                 `json:"line_content"`
	Language        string                 `json:"language"`
	KeySize         int                    `json:"key_size,omitempty"`
	Context         string                 `json:"context,omitempty"`
	QRS             int                    `json:"qrs"`
	RiskBand        detector.RiskBand      `json:"risk_band"`
	MigrationEffort detector.EffortLevel   `json:"migration_effort"`
	EffortRationale string                 `json:"effort_rationale"`
	OccurrenceCount int                    `json:"occurrence_count"`
	Timestamp       time.Time              `json:"timestamp"`
	Author          string                 `json:"author,omitempty"`
	CommitHash      string                 `json:"commit_hash,omitempty"`
	IntroducedAt    time.Time              `json:"introduced_at,omitempty"`
	RelatedIDs      []string               `json:"related_ids,omitempty"`
	ContainerLayer  string                 `json:"container_layer,omitempty"`
	EndpointHost    string                 `json:"endpoint_host,omitempty"`
}

// ActionItem represents a prioritised remediation action derived from findings.
type ActionItem struct {
	Rank           int                  `json:"rank"`
	Algorithm      string               `json:"algorithm"`
	RiskBand       detector.RiskBand    `json:"risk_band"`
	Occurrences    int                  `json:"occurrences"`
	Effort         detector.EffortLevel `json:"effort"`
	PriorityScore  float64              `json:"priority_score"`
	Recommendation string               `json:"recommendation"`
	Replacement    []string             `json:"replacement"`
}

// ScanResult holds the complete output of a Spectra scan run.
type ScanResult struct {
	ScanRoot       string                       `json:"scan_root"`
	StartedAt      time.Time                    `json:"started_at"`
	CompletedAt    time.Time                    `json:"completed_at"`
	TotalFiles     int                          `json:"total_files"`
	ScannedFiles   int                          `json:"scanned_files"`
	SkippedFiles   int                          `json:"skipped_files"`
	Findings       []Finding                    `json:"findings"`
	AggregateQRS   int                          `json:"aggregate_qrs"`
	FindingsByBand map[detector.RiskBand]int     `json:"findings_by_band"`
	ActionPlan     []ActionItem                 `json:"action_plan"`
}

// FileStats tracks file counts across scanner passes.
type FileStats struct {
	TotalFiles   int
	ScannedFiles int
	SkippedFiles int
}

// scannerEnabled returns true if the named scanner is in the enabled list,
// or if the enabled list is empty (meaning all scanners are active).
func scannerEnabled(name string, enabled []string) bool {
	if len(enabled) == 0 {
		return true
	}
	for _, s := range enabled {
		if s == name {
			return true
		}
	}
	return false
}

// ScanString processes a single string of code and returns a complete ScanResult.
// This is primarily used for WASM and API integrations where writing to disk is not possible.
func ScanString(code string, filename string, language string, patternRegistry *detector.PatternRegistry) (*ScanResult, error) {
	startedAt := time.Now()
	
	findings, err := ScanCodeString(code, filename, language, patternRegistry)
	if err != nil {
		return nil, err
	}
	
	// Post-process findings
	now := time.Now()
	for i := range findings {
		f := &findings[i]
		f.Timestamp = now
		if f.OccurrenceCount < 1 {
			f.OccurrenceCount = 1
		}
		f.QRS = detector.ComputeQRS(f.AlgorithmInfo, f.KeySize, f.OccurrenceCount)
		f.RiskBand = detector.QRSToBand(f.QRS)
		effort, rationale := detector.ClassifyEffort(f.Source, f.Algorithm)
		f.MigrationEffort = effort
		f.EffortRationale = rationale
	}
	
	// Sort findings
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].QRS != findings[j].QRS {
			return findings[i].QRS > findings[j].QRS
		}
		if findings[i].FilePath != findings[j].FilePath {
			return findings[i].FilePath < findings[j].FilePath
		}
		return findings[i].LineNumber < findings[j].LineNumber
	})
	
	// Aggregate QRS
	qrsValues := make([]int, len(findings))
	for i, f := range findings {
		qrsValues[i] = f.QRS
	}
	aggregateQRS := detector.AggregateQRS(qrsValues)
	
	// Findings by band
	findingsByBand := make(map[detector.RiskBand]int)
	for _, f := range findings {
		findingsByBand[f.RiskBand]++
	}
	
	// Action plan
	summaries := make([]detector.FindingSummary, len(findings))
	for i, f := range findings {
		summaries[i] = detector.FindingSummary{
			Algorithm:       f.Algorithm,
			QRS:             f.QRS,
			MigrationEffort: f.MigrationEffort,
		}
	}
	planItems := detector.BuildActionPlan(summaries)
	actionPlan := make([]ActionItem, len(planItems))
	for i, item := range planItems {
		actionPlan[i] = ActionItem{
			Rank:           i + 1,
			Algorithm:      item.Algorithm,
			RiskBand:       item.RiskBand,
			Occurrences:    item.Occurrences,
			Effort:         item.Effort,
			PriorityScore:  item.PriorityScore,
			Recommendation: item.Recommendation,
			Replacement:    item.Replacement,
		}
	}
	
	return &ScanResult{
		ScanRoot:       filename,
		StartedAt:      startedAt,
		CompletedAt:    time.Now(),
		TotalFiles:     1,
		ScannedFiles:   1,
		SkippedFiles:   0,
		Findings:       findings,
		AggregateQRS:   aggregateQRS,
		FindingsByBand: findingsByBand,
		ActionPlan:     actionPlan,
	}, nil
}

// ScanDirectory orchestrates a full scan of the given root directory. It fans
// out to each enabled scanner (code, cert, deps, config), collects findings,
// computes Quantum Risk Scores, classifies migration effort, and builds a
// prioritised action plan.
//
// Parameters:
//   - root: directory path to scan
//   - targetURL: optional URL for TLS scanning
//   - excludes: glob patterns for directories/files to skip
//   - scanners: list of enabled scanner names; empty means all
//   - concurrency: max parallel workers (passed to code scanner)
//   - patternRegistry: compiled pattern registry from crypto_patterns.yaml
//   - verbose: if true, individual file processing may be logged
func ScanDirectory(root string, targetURL string, excludes []string, scanners []string, concurrency int, patternRegistry *detector.PatternRegistry, verbose bool) (*ScanResult, error) {
	startedAt := time.Now()

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var (
		mu       sync.Mutex
		findings []Finding
	)

	// Counters for total / scanned / skipped files across all scanners.
	var totalFiles, scannedFiles, skippedFiles int

	// Run code scanner.
	if scannerEnabled("code", scanners) {
		codeFindings, stats, codeErr := ScanCodeFiles(absRoot, excludes, patternRegistry, concurrency, verbose)
		if codeErr != nil {
			return nil, codeErr
		}
		mu.Lock()
		findings = append(findings, codeFindings...)
		totalFiles += stats.TotalFiles
		scannedFiles += stats.ScannedFiles
		skippedFiles += stats.SkippedFiles
		mu.Unlock()
	}

	// Run certificate scanner.
	if scannerEnabled("cert", scanners) {
		certFindings, certErr := ScanCertFiles(absRoot, excludes)
		if certErr != nil {
			return nil, certErr
		}
		mu.Lock()
		findings = append(findings, certFindings...)
		mu.Unlock()
	}

	// Run dependency scanner.
	if scannerEnabled("deps", scanners) {
		depsFindings, depsErr := ScanDepsFiles(absRoot, excludes)
		if depsErr != nil {
			return nil, depsErr
		}
		mu.Lock()
		findings = append(findings, depsFindings...)
		mu.Unlock()
	}

	// Run config scanner.
	if scannerEnabled("config", scanners) {
		configFindings, configErr := ScanConfigFiles(absRoot, excludes, patternRegistry)
		if configErr != nil {
			return nil, configErr
		}
		mu.Lock()
		findings = append(findings, configFindings...)
		mu.Unlock()
	}

	// Run TLS scanner.
	if targetURL != "" && scannerEnabled("tls", scanners) {
		tlsFindings, tlsErr := ScanTLS(targetURL)
		if tlsErr != nil {
			return nil, tlsErr
		}
		mu.Lock()
		findings = append(findings, tlsFindings...)
		mu.Unlock()
	}

	// Check if this is a Git repository to enable Blame
	isGitRepo := false
	if stat, err := os.Stat(filepath.Join(absRoot, ".git")); err == nil && stat.IsDir() {
		isGitRepo = true
	}

	// Post-process: compute QRS, risk band, effort, and blame for each finding.
	now := time.Now()
	for i := range findings {
		f := &findings[i]

		// Set timestamp.
		if f.Timestamp.IsZero() {
			f.Timestamp = now
		}

		// Ensure occurrence count is at least 1.
		if f.OccurrenceCount < 1 {
			f.OccurrenceCount = 1
		}

		// Compute QRS.
		f.QRS = detector.ComputeQRS(f.AlgorithmInfo, f.KeySize, f.OccurrenceCount)

		// Compute risk band.
		f.RiskBand = detector.QRSToBand(f.QRS)
		effort, rationale := detector.ClassifyEffort(f.Source, f.Algorithm)
		f.MigrationEffort = effort
		f.EffortRationale = rationale
		f.Timestamp = now
		if isGitRepo {
			BlameFinding(absRoot, f)
		}
	}

	// Correlate findings across sources.
	algoGroups := make(map[string][]int)
	for i, f := range findings {
		algoGroups[f.Algorithm] = append(algoGroups[f.Algorithm], i)
	}
	for _, indices := range algoGroups {
		if len(indices) <= 1 {
			continue
		}
		// Check if there are different sources
		hasMultiSource := false
		firstSource := findings[indices[0]].Source
		for _, idx := range indices[1:] {
			if findings[idx].Source != firstSource {
				hasMultiSource = true
				break
			}
		}
		if hasMultiSource {
			// Extract all IDs
			var ids []string
			for _, idx := range indices {
				ids = append(ids, findings[idx].ID)
			}
			// Add related IDs (excluding self)
			for _, idx := range indices {
				var related []string
				for _, id := range ids {
					if id != findings[idx].ID {
						related = append(related, id)
					}
				}
				findings[idx].RelatedIDs = related
			}
		}
	}

	// Sort findings by (QRS desc, FilePath asc, LineNumber asc).
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].QRS != findings[j].QRS {
			return findings[i].QRS > findings[j].QRS
		}
		if findings[i].FilePath != findings[j].FilePath {
			return findings[i].FilePath < findings[j].FilePath
		}
		return findings[i].LineNumber < findings[j].LineNumber
	})

	// Compute aggregate QRS.
	qrsValues := make([]int, len(findings))
	for i, f := range findings {
		qrsValues[i] = f.QRS
	}
	aggregateQRS := detector.AggregateQRS(qrsValues)

	// Count findings by band.
	findingsByBand := make(map[detector.RiskBand]int)
	for _, f := range findings {
		findingsByBand[f.RiskBand]++
	}

	// Build action plan using the detector's priority engine.
	summaries := make([]detector.FindingSummary, len(findings))
	for i, f := range findings {
		summaries[i] = detector.FindingSummary{
			Algorithm:       f.Algorithm,
			QRS:             f.QRS,
			MigrationEffort: f.MigrationEffort,
		}
	}
	planItems := detector.BuildActionPlan(summaries)

	actionPlan := make([]ActionItem, len(planItems))
	for i, item := range planItems {
		actionPlan[i] = ActionItem{
			Rank:           i + 1,
			Algorithm:      item.Algorithm,
			RiskBand:       item.RiskBand,
			Occurrences:    item.Occurrences,
			Effort:         item.Effort,
			PriorityScore:  item.PriorityScore,
			Recommendation: item.Recommendation,
			Replacement:    item.Replacement,
		}
	}

	return &ScanResult{
		ScanRoot:       absRoot,
		StartedAt:      startedAt,
		CompletedAt:    time.Now(),
		TotalFiles:     totalFiles,
		ScannedFiles:   scannedFiles,
		SkippedFiles:   skippedFiles,
		Findings:       findings,
		AggregateQRS:   aggregateQRS,
		FindingsByBand: findingsByBand,
		ActionPlan:     actionPlan,
	}, nil
}
