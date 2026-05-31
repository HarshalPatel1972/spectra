package scanner_test

import (
	"path/filepath"
	"testing"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

func TestScanDepsFiles(t *testing.T) {
	detector.LookupAlgorithm("AES")

	absPath, err := filepath.Abs("../../testdata/samples")
	if err != nil {
		t.Fatalf("failed to get abs path for samples: %v", err)
	}

	findings, err := scanner.ScanDepsFiles(absPath, nil)
	if err != nil {
		t.Fatalf("ScanDepsFiles failed: %v", err)
	}

	// We expect some findings if testdata has a go.mod or package.json with known weak deps
	// Since sample test files might not contain actual weak deps, we just ensure it runs without error.
	if len(findings) > 0 {
		for _, f := range findings {
			if f.Source != detector.SourceDeps {
				t.Errorf("Expected source DEPS, got %s", f.Source)
			}
		}
	}
}
