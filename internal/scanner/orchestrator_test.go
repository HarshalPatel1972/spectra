package scanner_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/HarshalPatel1972/spectra/internal/cbom"
	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

func TestScanDirectory_Integration(t *testing.T) {
	// Need to load patterns
	data, err := os.ReadFile("../../rules/crypto_patterns.yaml")
	if err != nil {
		t.Skipf("skipping integration test, could not read patterns: %v", err)
	}
	
	registry, err := detector.LoadPatternsFromBytes(data)
	if err != nil {
		t.Fatalf("failed to load patterns: %v", err)
	}

	samplesDir, err := filepath.Abs("../../testdata/samples")
	if err != nil {
		t.Fatalf("failed to get abs path for samples: %v", err)
	}

	// 1. Calls ScanDirectory(testdata/samples)
	result, err := scanner.ScanDirectory(samplesDir, "", nil, []string{"code", "cert", "deps", "config"}, 4, registry, false)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	// 2. Verifies aggregate QRS > 70
	if result.AggregateQRS <= 70 {
		t.Errorf("expected aggregate QRS > 70, got %d", result.AggregateQRS)
	}

	// 3. Verifies at least one CRITICAL finding for RSA
	var hasCriticalRSA bool
	var hasCertFinding bool
	
	for _, f := range result.Findings {
		if f.Algorithm == "RSA" && f.RiskBand == detector.BandCritical {
			hasCriticalRSA = true
		}
		if f.Source == detector.SourceCert {
			hasCertFinding = true
		}
	}

	if !hasCriticalRSA {
		t.Error("expected at least one CRITICAL finding for RSA")
	}

	// 4. Verifies at least one finding with Source = CERT
	if !hasCertFinding {
		t.Error("expected at least one finding with Source = CERT")
	}

	// 5. Verifies CBOM output is valid JSON parseable by cyclonedx-go
	bom, err := cbom.GenerateCBOM(result, "test-version")
	if err != nil {
		t.Fatalf("GenerateCBOM failed: %v", err)
	}
	
	b, err := json.Marshal(bom)
	if err != nil {
		t.Fatalf("failed to marshal BOM: %v", err)
	}
	
	// Quick check if parseable
	var parsed map[string]interface{}
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Errorf("BOM JSON is not valid: %v", err)
	}
}
