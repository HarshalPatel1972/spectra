//go:build ignore

// Command generate_cert creates a self-signed RSA-1024 / SHA1WithRSA
// certificate and writes it to testdata/samples/sample_weak.pem.
//
// Usage:
//
//	go run testdata/cmd/generate_cert/main.go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "generate_cert: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Generate a deliberately weak 1024-bit RSA key for testing.
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return fmt.Errorf("generating RSA-1024 key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("generating serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Spectra Test CA"},
			Country:      []string{"US"},
			Province:     []string{"Test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SignatureAlgorithm:    x509.SHA1WithRSA,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("creating certificate: %w", err)
	}

	outPath := "testdata/samples/sample_weak.pem"
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", outPath, err)
	}
	defer f.Close()

	if err := pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("encoding PEM: %w", err)
	}

	fmt.Fprintf(os.Stdout, "wrote %s (RSA-1024 / SHA1WithRSA self-signed)\n", outPath)
	return nil
}
