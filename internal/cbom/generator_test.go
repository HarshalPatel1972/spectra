package cbom_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HarshalPatel1972/spectra/internal/cbom"
	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateCBOM(t *testing.T) {
	result := &scanner.ScanResult{
		ScanRoot:     "/test",
		StartedAt:    time.Time{},
		CompletedAt:  time.Time{},
		TotalFiles:   10,
		ScannedFiles: 10,
		Findings: []scanner.Finding{
			{
				Algorithm:       "RSA",
				Source:          detector.SourceCode,
				KeySize:         2048,
				QRS:             50,
				MigrationEffort: detector.EffortMedium,
				OccurrenceCount: 5,
			},
		},
	}

	bom, err := cbom.GenerateCBOM(result, "1.0.0")
	if err != nil {
		t.Fatalf("GenerateCBOM failed: %v", err)
	}

	// Make sure the UUID is stable or ignore it for comparison
	bom.SerialNumber = "urn:uuid:00000000-0000-0000-0000-000000000000"

	b, _ := json.MarshalIndent(bom, "", "  ")
	
	goldenFile := filepath.Join("../../testdata", "golden", "cbom.json")
	os.MkdirAll(filepath.Dir(goldenFile), 0755)

	if _, err := os.Stat(goldenFile); os.IsNotExist(err) {
		os.WriteFile(goldenFile, b, 0644)
		t.Logf("Created golden file %s", goldenFile)
		return
	}

	expectedBytes, _ := os.ReadFile(goldenFile)
	var expected interface{}
	var actual interface{}
	
	json.Unmarshal(expectedBytes, &expected)
	json.Unmarshal(b, &actual)

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("CBOM mismatch (-want +got):\n%s", diff)
	}
}
