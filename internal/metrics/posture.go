package metrics

import (
	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

type PostureReport struct {
	CPS        int // Cryptographic Posture Score (0-100)
	Grade      string
	RiskStatus string
}

// CalculateCPS determines the overarching Cryptographic Posture Score.
func CalculateCPS(findings []scanner.Finding, aggregateQRS int) *PostureReport {
	if len(findings) == 0 {
		return &PostureReport{
			CPS:        100,
			Grade:      "A+",
			RiskStatus: "Quantum Safe",
		}
	}

	// Base score is 100. We subtract based on QRS and Critical findings.
	cps := 100

	// 1. Subtract for Aggregate QRS (capped penalty of 40)
	qrsPenalty := aggregateQRS / 10
	if qrsPenalty > 40 {
		qrsPenalty = 40
	}
	cps -= qrsPenalty

	// 2. Subtract for Critical and High findings
	criticalCount := 0
	highCount := 0
	for _, f := range findings {
		if f.RiskBand == detector.BandCritical {
			criticalCount++
		} else if f.RiskBand == detector.BandHigh {
			highCount++
		}
	}

	cps -= criticalCount * 5
	cps -= highCount * 2

	if cps < 0 {
		cps = 0
	}

	report := &PostureReport{
		CPS: cps,
	}

	if cps >= 90 {
		report.Grade = "A"
		report.RiskStatus = "Excellent"
	} else if cps >= 80 {
		report.Grade = "B"
		report.RiskStatus = "Good"
	} else if cps >= 70 {
		report.Grade = "C"
		report.RiskStatus = "Fair"
	} else if cps >= 60 {
		report.Grade = "D"
		report.RiskStatus = "At Risk"
	} else {
		report.Grade = "F"
		report.RiskStatus = "Critical Vulnerability"
	}

	return report
}
