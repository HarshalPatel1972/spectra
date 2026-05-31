package metrics

import (
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

type AgilityReport struct {
	CAI                int // Cryptographic Agility Index (0-100)
	AbstractionGrade   string
	CentralizationGrade string
	Configurability    bool
	TotalFindings      int
}

// CalculateCAI determines how agile the codebase is cryptographically.
// Higher is better.
func CalculateCAI(findings []scanner.Finding) *AgilityReport {
	if len(findings) == 0 {
		return &AgilityReport{
			CAI:                100,
			AbstractionGrade:   "A",
			CentralizationGrade: "A",
			Configurability:    true,
		}
	}

	report := &AgilityReport{
		TotalFindings: len(findings),
	}

	// 1. Centralization: Are findings scattered or localized?
	// We count unique files containing findings.
	fileSet := make(map[string]struct{})
	for _, f := range findings {
		fileSet[f.FilePath] = struct{}{}
	}

	spreadRatio := float64(len(fileSet)) / float64(len(findings))
	centralizationScore := 100 - int(spreadRatio*100)
	if centralizationScore < 0 {
		centralizationScore = 0
	}

	// 2. Abstraction: Are they using standard interfaces or hardcoded algorithms?
	// We check how many findings are hardcoded in CODE vs coming from CONFIG or DEPS.
	hardcodedCount := 0
	for _, f := range findings {
		if f.Source == "CODE" {
			hardcodedCount++
		}
	}
	abstractionScore := 100 - int((float64(hardcodedCount)/float64(len(findings)))*100)

	// 3. Configurability: Do we have config files that define crypto?
	hasConfig := false
	for _, f := range findings {
		if f.Source == "CONFIG" {
			hasConfig = true
			break
		}
	}
	report.Configurability = hasConfig

	// Calculate CAI
	cai := (centralizationScore + abstractionScore) / 2
	if hasConfig {
		cai += 10 // Bonus for configurability
	}
	if cai > 100 {
		cai = 100
	}
	report.CAI = cai

	// Grades
	report.AbstractionGrade = scoreToGrade(abstractionScore)
	report.CentralizationGrade = scoreToGrade(centralizationScore)

	return report
}

func scoreToGrade(score int) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else if score >= 60 {
		return "D"
	}
	return "F"
}
