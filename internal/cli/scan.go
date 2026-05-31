package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/cbom"
	"github.com/HarshalPatel1972/spectra/internal/compliance"
	"github.com/HarshalPatel1972/spectra/internal/config"
	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/graph"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/HarshalPatel1972/spectra/internal/report"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/HarshalPatel1972/spectra/rules"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Exit codes as specified in the Spectra specification.
const (
	exitOK            = 0
	exitFindings      = 1
	exitRuntimeError  = 2
	exitInvalidConfig = 3
)

// Scan flag variables.
var (
	outputFormats []string
	outDir        string
	excludePaths  []string
	scanners      []string
	failOn        string
	concurrency   int
	targetURL     string
	targetImage   string
	persist       bool
	baselineName  string
)

// scanCmd implements the 'spectra scan' subcommand that performs cryptographic
// asset discovery and quantum risk assessment on a target directory.
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a directory for cryptographic assets and assess quantum risk",
	Long: `Scan analyses source code, certificates, dependency manifests, and
configuration files to discover cryptographic assets. Each finding is scored
with a Quantum Risk Score (QRS) and mapped to NIST PQC replacements.

If no path is given, the current directory is scanned.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringSliceVarP(&outputFormats, "output", "o", nil, "output format(s): terminal, json, cbom, html")
	scanCmd.Flags().StringVar(&outDir, "out-dir", "", "directory for file-based outputs (default: ./spectra-out)")
	scanCmd.Flags().StringSliceVar(&excludePaths, "exclude", nil, "glob patterns to exclude from scanning")
	scanCmd.Flags().StringSliceVar(&scanners, "scanners", nil, "scanner modules to enable: code, cert, deps, config, tls")
	scanCmd.Flags().StringVar(&failOn, "fail-on", "", "minimum severity to cause non-zero exit: critical, high, medium, low")
	scanCmd.Flags().IntVar(&concurrency, "concurrency", 0, "number of parallel scanner goroutines (0 = auto)")
	scanCmd.Flags().StringVar(&targetURL, "url", "", "TLS endpoint to scan (e.g. example.com:443)")
	scanCmd.Flags().StringVar(&targetImage, "image", "", "Container image to scan (e.g. ubuntu:latest)")
	scanCmd.Flags().BoolVar(&persist, "persist", false, "save scan results to state database")
	scanCmd.Flags().StringVar(&baselineName, "baseline", "", "tag the scan with a baseline name (requires --persist)")

	rootCmd.AddCommand(scanCmd)
}

// runScan is the main entry point for the scan command.
func runScan(cmd *cobra.Command, args []string) error {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	// Verify the target path exists.
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("target path %q: %w", targetPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("target path %q is not a directory", targetPath)
	}

	// Load configuration.
	cfg, err := loadScanConfig()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}

	// Merge CLI flags into config (flags take precedence).
	config.MergeFlags(cfg, config.FlagValues{
		OutputFormats: outputFormats,
		OutDir:        outDir,
		Exclude:       excludePaths,
		Scanners:      scanners,
		FailOn:        failOn,
		Concurrency:   concurrency,
	})

	// Validate fail-on value if provided.
	if cfg.CI.FailOn != "" {
		if err := validateFailOn(cfg.CI.FailOn); err != nil {
			return err
		}
	}

	// Resolve effective concurrency.
	_ = config.EffectiveConcurrency(cfg)

	// Create output directory if file-based formats are requested.
	if needsOutputDir(cfg.Output.Formats) {
		if err := os.MkdirAll(cfg.Output.Dir, 0o755); err != nil {
			return fmt.Errorf("creating output directory %q: %w", cfg.Output.Dir, err)
		}
	}

	registry, err := detector.LoadPatternsFromBytes(rules.CryptoPatternsYAML)
	if err != nil {
		return fmt.Errorf("failed to load patterns: %w", err)
	}

	// Handle Image extraction if targetImage is provided
	if targetImage != "" {
		fmt.Fprintf(os.Stderr, "Pulling and extracting image %s...\n", targetImage)
		extractedDir, err := scanner.ExtractImage(targetImage)
		if err != nil {
			return fmt.Errorf("failed to extract image: %w", err)
		}
		defer os.RemoveAll(extractedDir)
		targetPath = extractedDir
	}

	// Call scanner.ScanDirectory with cfg and targetPath
	effConcurrency := config.EffectiveConcurrency(cfg)
	result, err := scanner.ScanDirectory(targetPath, targetURL, cfg.Scan.Exclude, cfg.Scan.Scanners, effConcurrency, registry, verbose)
	if err != nil {
		return fmt.Errorf("scan failed (exit code %d): %w", exitRuntimeError, err)
	}

	// For each output format, call the appropriate writer
	for _, f := range cfg.Output.Formats {
		switch strings.ToLower(f) {
		case "terminal":
			report.RenderTerminal(os.Stdout, result, "v0.1.0")
		case "json":
			report.WriteJSON(result, cfg.Output.Dir)
		case "cbom":
			cbom.WriteCBOM(result, cfg.Output.Dir, "v0.1.0")
		case "html":
			report.WriteHTML(result, cfg.Output.Dir, "v0.1.0")
		}
	}

	if persist {
		if err := saveToStateDB(result); err != nil {
			return fmt.Errorf("saving to state db: %w", err)
		}
	}

	// If fail-on is set, check findings and return appropriate exit code
	if cfg.CI.FailOn != "" {
		failThreshold := parseFailOnBand(cfg.CI.FailOn)
		hasFindings := false
		for _, f := range result.Findings {
			if f.RiskBand <= failThreshold {
				hasFindings = true
				break
			}
		}
		if hasFindings {
			os.Exit(exitFindings)
		}
	}


	if !quiet {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Scanning %s ...\n", targetPath)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Scanners: %s\n", strings.Join(cfg.Scan.Scanners, ", "))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Output:   %s\n", strings.Join(cfg.Output.Formats, ", "))
	}

	return nil
}

func saveToStateDB(result *scanner.ScanResult) error {
	store, err := persistence.NewStore("")
	if err != nil {
		return err
	}
	defer store.Close()

	scanID := uuid.New().String()
	if err := store.SaveScan(scanID, result); err != nil {
		return err
	}

	if len(result.Findings) > 0 {
		if err := store.SaveFindings(scanID, result.Findings); err != nil {
			return err
		}
	}

	g := graph.BuildGraph(result)
	var nodes []persistence.GraphNode
	for _, n := range g.Nodes {
		nodes = append(nodes, persistence.GraphNode{
			ID:         n.ID,
			NodeType:   string(n.Type),
			Label:      n.Label,
			Properties: "{}", // Could be serialized n.Properties
		})
	}
	var edges []persistence.GraphEdge
	for _, e := range g.Edges {
		edges = append(edges, persistence.GraphEdge{
			ID:       e.ID,
			FromNode: e.From,
			ToNode:   e.To,
			EdgeType: string(e.Type),
			Weight:   e.Weight,
		})
	}

	if len(nodes) > 0 || len(edges) > 0 {
		if err := store.SaveGraph(scanID, nodes, edges); err != nil {
			return err
		}
	}

	gaps := compliance.EvaluateFindings(result.Findings)
	if len(gaps) > 0 {
		if err := store.SaveComplianceGaps(gaps); err != nil {
			return err
		}
	}

	if baselineName != "" {
		err := store.SaveBaseline(persistence.Baseline{
			Name:        baselineName,
			ScanID:      scanID,
			CreatedAt:   result.CompletedAt,
			Description: "Baseline taken via CLI",
		})
		if err != nil {
			return fmt.Errorf("saving baseline: %w", err)
		}
	}

	return nil
}

// loadScanConfig loads configuration from the config file path (--config flag)
// or falls back to .spectra.yaml in the current directory. If no config file
// exists, default configuration is returned.
func loadScanConfig() (*config.Config, error) {
	if cfgFile != "" {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("exit code %d: %w", exitInvalidConfig, err)
		}
		return cfg, nil
	}

	// Try default location.
	cfg, err := config.LoadConfig(".spectra.yaml")
	if err != nil {
		if os.IsNotExist(errors.Unwrap(err)) {
			return config.DefaultConfig(), nil
		}
		// File exists but is invalid.
		return nil, fmt.Errorf("exit code %d: %w", exitInvalidConfig, err)
	}
	return cfg, nil
}

// validateFailOn checks that the fail-on value is a recognised severity level.
func validateFailOn(level string) error {
	switch strings.ToLower(level) {
	case "critical", "high", "medium", "low":
		return nil
	default:
		return fmt.Errorf("invalid --fail-on value %q: must be critical, high, medium, or low", level)
	}
}

// needsOutputDir returns true if any of the requested formats produce files.
func needsOutputDir(formats []string) bool {
	for _, f := range formats {
		switch strings.ToLower(f) {
		case "json", "cbom", "html":
			return true
		}
	}
	return false
}

// parseFailOnBand maps a string level to a detector.RiskBand
func parseFailOnBand(level string) detector.RiskBand {
	switch strings.ToLower(level) {
	case "critical":
		return detector.BandCritical
	case "high":
		return detector.BandHigh
	case "medium":
		return detector.BandMedium
	case "low":
		return detector.BandLow
	default:
		return detector.BandSafe
	}
}
