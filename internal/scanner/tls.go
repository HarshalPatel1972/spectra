package scanner

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/HarshalPatel1972/spectra/internal/detector"
)

// ScanTLS connects to the target URL, extracts certificates and cipher suite information,
// and returns a list of cryptographic findings.
func ScanTLS(targetURL string) ([]Finding, error) {
	// Strip protocols if any
	host := targetURL
	if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	} else if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	}
	
	// Add port if missing
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conf := &tls.Config{
		InsecureSkipVerify: true, // We want to scan it regardless of validity
		MinVersion:         tls.VersionSSL30,
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("connecting to TLS endpoint %s: %w", host, err)
	}
	defer conn.Close()

	var findings []Finding
	state := conn.ConnectionState()

	// 1. Check TLS Version
	var tlsVersionStr string
	switch state.Version {
	case tls.VersionTLS10:
		tlsVersionStr = "TLS 1.0"
	case tls.VersionTLS11:
		tlsVersionStr = "TLS 1.1"
	case tls.VersionTLS12:
		tlsVersionStr = "TLS 1.2"
	case tls.VersionTLS13:
		tlsVersionStr = "TLS 1.3"
	default:
		tlsVersionStr = "Unknown TLS"
	}

	// Flag old TLS
	if state.Version < tls.VersionTLS12 {
		f := Finding{
			ID:              uuid.New().String(),
			Algorithm:       tlsVersionStr,
			Source:          "CONFIG", // Treating protocol config as config
			FilePath:        targetURL,
			Context:         "Endpoint negotiated legacy TLS protocol",
			Timestamp:       time.Now(),
		}
		// Try to lookup
		if info, ok := detector.LookupAlgorithm(tlsVersionStr); ok {
			f.AlgorithmInfo = info
			f.QRS = detector.ComputeQRS(info, 0, 1)
		} else {
			f.AlgorithmInfo = detector.AlgorithmInfo{
				Name: tlsVersionStr, Family: "protocol", QuantumThreat: detector.ThreatHigh, BaseQRS: 80,
			}
			f.QRS = 80
		}
		findings = append(findings, f)
	}

	// 2. Process Certificates using existing Cert Logic
	for i, cert := range state.PeerCertificates {
		// Just call the cert analysis logic
		certFindings := processCertificate(cert, targetURL)
		
		// Add context
		for j := range certFindings {
			if i == 0 {
				certFindings[j].Context = "Endpoint Leaf Certificate"
			} else {
				certFindings[j].Context = "Endpoint Intermediate Certificate"
			}
		}
		findings = append(findings, certFindings...)
	}

	return findings, nil
}
