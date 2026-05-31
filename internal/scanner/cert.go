package scanner

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/HarshalPatel1972/spectra/internal/detector"
)

// certExtensions lists file extensions that may contain X.509 certificates.
var certExtensions = map[string]bool{
	".pem": true,
	".crt": true,
	".cer": true,
	".der": true,
}

// x509AlgorithmMap maps x509.SignatureAlgorithm values to compound algorithm
// strings. These are split into individual algorithms when creating findings.
var x509AlgorithmMap = map[x509.SignatureAlgorithm]string{
	x509.SHA1WithRSA:     "RSA+SHA1",
	x509.MD5WithRSA:      "RSA+MD5",
	x509.SHA256WithRSA:   "RSA+SHA256",
	x509.SHA384WithRSA:   "RSA+SHA384",
	x509.SHA512WithRSA:   "RSA+SHA512",
	x509.ECDSAWithSHA1:   "ECDSA+SHA1",
	x509.ECDSAWithSHA256: "ECDSA+SHA256",
	x509.ECDSAWithSHA384: "ECDSA+SHA384",
	x509.ECDSAWithSHA512: "ECDSA+SHA512",
	x509.DSAWithSHA1:     "DSA+SHA1",
	x509.DSAWithSHA256:   "DSA+SHA256",
}

// pubKeyAlgoName returns the canonical algorithm name for an x509 public key algorithm.
func pubKeyAlgoName(algo x509.PublicKeyAlgorithm) string {
	switch algo {
	case x509.RSA:
		return "RSA"
	case x509.ECDSA:
		return "ECDSA"
	case x509.DSA:
		return "DSA"
	case x509.Ed25519:
		return "Ed25519"
	default:
		return ""
	}
}

// ScanCertFiles walks the root directory looking for certificate files,
// parses PEM/DER encoded certificates, and creates findings for the
// public key algorithm and any weak signature algorithms discovered.
func ScanCertFiles(root string, excludes []string) ([]Finding, error) {
	var findings []Finding

	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}

		if d.IsDir() {
			name := d.Name()
			if defaultSkipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !certExtensions[ext] {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		if shouldExcludeFile(relPath, excludes) {
			return nil
		}

		certFindings, parseErr := parseCertFile(path, relPath, ext)
		if parseErr != nil {
			// Skip unparseable cert files instead of failing the entire scan.
			return nil
		}
		findings = append(findings, certFindings...)
		return nil
	})

	return findings, walkErr
}

// parseCertFile reads a certificate file and creates findings for each
// certificate found within it.
func parseCertFile(absPath, relPath, ext string) ([]Finding, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var certs []*x509.Certificate

	if ext == ".der" {
		// DER is raw binary ASN.1.
		cert, derErr := x509.ParseCertificate(data)
		if derErr != nil {
			return nil, derErr
		}
		certs = append(certs, cert)
	} else {
		// PEM: decode all CERTIFICATE blocks.
		rest := data
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				continue
			}
			cert, parseErr := x509.ParseCertificate(block.Bytes)
			if parseErr != nil {
				continue // skip malformed certs
			}
			certs = append(certs, cert)
		}
	}

	var findings []Finding
	for _, cert := range certs {
		certFindings := processCertificate(cert, relPath)
		findings = append(findings, certFindings...)
	}

	return findings, nil
}

// processCertificate extracts cryptographic findings from a parsed X.509
// certificate, including the public key algorithm and signature hash.
func processCertificate(cert *x509.Certificate, relPath string) []Finding {
	var findings []Finding

	// Build context string from certificate metadata.
	context := buildCertContext(cert)

	// Extract key size from the public key.
	keySize := extractCertKeySize(cert)

	// --- Finding 1: Public key algorithm ---
	pkAlgo := pubKeyAlgoName(cert.PublicKeyAlgorithm)
	if pkAlgo != "" {
		canonical, found := detector.NormaliseAlgorithm(pkAlgo)
		if found {
			info, ok := detector.LookupAlgorithm(canonical)
			if ok {
				findings = append(findings, Finding{
					ID:              uuid.New().String(),
					Algorithm:       canonical,
					AlgorithmInfo:   info,
					Source:          detector.SourceCert,
					FilePath:        relPath,
					LineNumber:      0,
					LineContent:     fmt.Sprintf("PublicKeyAlgorithm: %s", cert.PublicKeyAlgorithm),
					Language:        "x509",
					KeySize:         keySize,
					Context:         context,
					OccurrenceCount: 1,
				})
			}
		}
	}

	// --- Finding 2+: Signature algorithm components ---
	compound, sigKnown := x509AlgorithmMap[cert.SignatureAlgorithm]
	if sigKnown {
		parts := strings.Split(compound, "+")
		for _, part := range parts {
			canonical, found := detector.NormaliseAlgorithm(part)
			if !found {
				continue
			}
			info, ok := detector.LookupAlgorithm(canonical)
			if !ok {
				continue
			}

			// Skip if we already created a finding for this exact algorithm
			// from the public key step above (avoid duplicate RSA findings).
			if canonical == pkAlgo && len(findings) > 0 {
				// Check if last finding already covers this algorithm.
				alreadyCovered := false
				for _, f := range findings {
					if f.Algorithm == canonical && f.Source == detector.SourceCert && f.FilePath == relPath {
						alreadyCovered = true
						break
					}
				}
				if alreadyCovered {
					continue
				}
			}

			sigKeySize := 0
			if canonical == pkAlgo {
				sigKeySize = keySize
			}

			findings = append(findings, Finding{
				ID:              uuid.New().String(),
				Algorithm:       canonical,
				AlgorithmInfo:   info,
				Source:          detector.SourceCert,
				FilePath:        relPath,
				LineNumber:      0,
				LineContent:     fmt.Sprintf("SignatureAlgorithm: %s", cert.SignatureAlgorithm),
				Language:        "x509",
				KeySize:         sigKeySize,
				Context:         context,
				OccurrenceCount: 1,
			})
		}
	}

	// --- Finding 3: Expired or near-expiry certificate ---
	expiryFinding := checkCertExpiry(cert, relPath, context, keySize, pkAlgo)
	if expiryFinding != nil {
		findings = append(findings, *expiryFinding)
	}

	return findings
}

// buildCertContext constructs a human-readable context string from certificate
// metadata including subject, issuer, and validity dates.
func buildCertContext(cert *x509.Certificate) string {
	var parts []string

	if cert.Subject.CommonName != "" {
		parts = append(parts, "CN="+cert.Subject.CommonName)
	}
	if cert.Issuer.CommonName != "" {
		parts = append(parts, "Issuer="+cert.Issuer.CommonName)
	}
	if len(cert.DNSNames) > 0 {
		parts = append(parts, "SANs="+strings.Join(cert.DNSNames, ","))
	}
	parts = append(parts, fmt.Sprintf("NotAfter=%s", cert.NotAfter.Format("2006-01-02")))

	return strings.Join(parts, "; ")
}

// extractCertKeySize returns the key size in bits for the certificate's public key.
func extractCertKeySize(cert *x509.Certificate) int {
	switch key := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		return key.N.BitLen()
	case *ecdsa.PublicKey:
		return key.Curve.Params().BitSize
	default:
		return 0
	}
}

// checkCertExpiry creates a finding if the certificate is expired or will
// expire within 90 days.
func checkCertExpiry(cert *x509.Certificate, relPath, context string, keySize int, pkAlgo string) *Finding {
	now := time.Now()
	ninetyDays := now.Add(90 * 24 * time.Hour)

	var expiryContext string
	if cert.NotAfter.Before(now) {
		expiryContext = fmt.Sprintf("EXPIRED on %s", cert.NotAfter.Format("2006-01-02"))
	} else if cert.NotAfter.Before(ninetyDays) {
		daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
		expiryContext = fmt.Sprintf("expires in %d days (%s)", daysLeft, cert.NotAfter.Format("2006-01-02"))
	} else {
		return nil
	}

	canonical := pkAlgo
	if canonical == "" {
		canonical = "RSA"
	}
	info, ok := detector.LookupAlgorithm(canonical)
	if !ok {
		return nil
	}

	return &Finding{
		ID:              uuid.New().String(),
		Algorithm:       canonical,
		AlgorithmInfo:   info,
		Source:          detector.SourceCert,
		FilePath:        relPath,
		LineNumber:      0,
		LineContent:     expiryContext,
		Language:        "x509",
		KeySize:         keySize,
		Context:         context + "; " + expiryContext,
		OccurrenceCount: 1,
	}
}
