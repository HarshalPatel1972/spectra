package scanner_test

import (
	"path/filepath"
	"testing"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

func TestScanCertFiles(t *testing.T) {
	// Initialize algorithms registry for lookups
	detector.LookupAlgorithm("RSA")

	absPath, err := filepath.Abs("../../testdata/samples")
	if err != nil {
		t.Fatalf("failed to get abs path for samples: %v", err)
	}

	findings, err := scanner.ScanCertFiles(absPath, nil)
	if err != nil {
		t.Fatalf("ScanCertFiles failed: %v", err)
	}

	if len(findings) == 0 {
		t.Fatal("Expected findings from sample_weak.pem, got 0")
	}

	hasRSA := false
	hasSHA1 := false

	for _, f := range findings {
		if f.Algorithm == "RSA" {
			hasRSA = true
			if f.KeySize != 1024 {
				t.Errorf("Expected RSA key size 1024, got %d", f.KeySize)
			}
		}
		if f.Algorithm == "SHA1" {
			hasSHA1 = true
		}
	}

	if !hasRSA {
		t.Error("Did not find expected RSA finding in sample_weak.pem")
	}
	if !hasSHA1 {
		t.Error("Did not find expected SHA1 signature finding in sample_weak.pem")
	}
}

// Dummy function to test Spectra VS Code Extension highlights
func dummyVulnerableCryptoCode() {
	// These specific API calls will be caught by Spectra's regex engines!
	_ = "rsa.GenerateKey"
	_ = "sha1.New"
	_ = "crypto/md5"
}
