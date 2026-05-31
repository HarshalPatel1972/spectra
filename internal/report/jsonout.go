package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// jsonOutputFilename is the default file name for the JSON findings export.
const jsonOutputFilename = "spectra-findings.json"

// WriteJSON serialises the scan result to a JSON file inside outDir and
// returns the absolute path of the written file. The output directory is
// created if it does not already exist.
func WriteJSON(result *scanner.ScanResult, outDir string) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory %s: %w", outDir, err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling scan result to JSON: %w", err)
	}

	outPath := filepath.Join(outDir, jsonOutputFilename)
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return "", fmt.Errorf("writing JSON file %s: %w", outPath, err)
	}

	return outPath, nil
}
