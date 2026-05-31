package scanner_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestScanCodeFiles(t *testing.T) {
	data, err := os.ReadFile("../../rules/crypto_patterns.yaml")
	if err != nil {
		t.Skipf("skipping code test, could not read patterns: %v", err)
	}
	
	registry, err := detector.LoadPatternsFromBytes(data)
	if err != nil {
		t.Fatalf("failed to load patterns: %v", err)
	}

	samplesDir, _ := filepath.Abs("../../testdata/samples")
	findings, _, err := scanner.ScanCodeFiles(samplesDir, nil, registry, 1, false)
	if err != nil {
		t.Fatalf("ScanCodeFiles failed: %v", err)
	}

	goldenFile := filepath.Join("../../testdata", "golden", "code_findings.json")
	
	// Create golden dir if not exists
	os.MkdirAll(filepath.Dir(goldenFile), 0755)

	if _, err := os.Stat(goldenFile); os.IsNotExist(err) {
		// Create golden file
		b, _ := json.MarshalIndent(findings, "", "  ")
		os.WriteFile(goldenFile, b, 0644)
		t.Logf("Created golden file %s", goldenFile)
		return
	}

	// Compare with golden file
	expectedBytes, _ := os.ReadFile(goldenFile)
	var expected []scanner.Finding
	json.Unmarshal(expectedBytes, &expected)

	// Ignore dynamic fields like ID and Timestamp
	opts := cmpopts.IgnoreFields(scanner.Finding{}, "ID", "Timestamp")
	if diff := cmp.Diff(expected, findings, opts); diff != "" {
		t.Errorf("ScanCodeFiles mismatch (-want +got):\n%s", diff)
	}
}
