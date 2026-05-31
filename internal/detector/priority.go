package detector

import (
	"math"
	"sort"
)

// ActionPlanItem represents a single entry in the prioritised remediation action plan.
type ActionPlanItem struct {
	Algorithm      string
	DisplayName    string
	RiskBand       RiskBand
	Occurrences    int
	Effort         EffortLevel
	PriorityScore  float64
	Recommendation string
	Replacement    []string
}

// FindingSummary contains the minimal fields needed to build an action plan
// from a collection of scan findings.
type FindingSummary struct {
	Algorithm       string
	QRS             int
	MigrationEffort EffortLevel
}

// PriorityScore computes the priority score for action plan ordering.
// Higher scores indicate items that should be fixed first. The score combines
// quantum risk, ease of remediation (effort inverse), and a logarithmic
// frequency boost.
func PriorityScore(qrs int, effort EffortLevel, occurrences int) float64 {
	effortMap := map[EffortLevel]float64{
		EffortEasy:    1.0,
		EffortMedium:  0.6,
		EffortHard:    0.3,
		EffortBlocked: 0.1,
	}
	effortWeight, ok := effortMap[effort]
	if !ok {
		effortWeight = 0.6 // default to medium if unknown
	}
	freqBoost := math.Log1p(float64(occurrences))
	return float64(qrs) * effortWeight * (1 + freqBoost*0.05)
}

// BuildActionPlan groups findings by algorithm, computes priority scores,
// and returns a sorted list of action items ordered by descending priority
// (highest priority first). For each algorithm group it takes the maximum
// QRS, counts total occurrences, and looks up display names and PQC
// replacements from the algorithm registry.
func BuildActionPlan(findings []FindingSummary) []ActionPlanItem {
	if len(findings) == 0 {
		return nil
	}

	// Group by canonical algorithm name
	type group struct {
		maxQRS      int
		effort      EffortLevel
		occurrences int
	}
	groups := make(map[string]*group)

	for _, f := range findings {
		canonical, _ := NormaliseAlgorithm(f.Algorithm)
		g, ok := groups[canonical]
		if !ok {
			g = &group{
				maxQRS: f.QRS,
				effort: f.MigrationEffort,
			}
			groups[canonical] = g
		}
		g.occurrences++
		if f.QRS > g.maxQRS {
			g.maxQRS = f.QRS
		}
		// Use the hardest effort level in the group
		if effortOrdinal(f.MigrationEffort) > effortOrdinal(g.effort) {
			g.effort = f.MigrationEffort
		}
	}

	// Build action items
	items := make([]ActionPlanItem, 0, len(groups))
	for algo, g := range groups {
		info, found := LookupAlgorithm(algo)

		displayName := algo
		var replacement []string
		var recommendation string

		if found {
			displayName = info.DisplayName
			replacement = info.PQCReplacement
			if info.PQCSafe {
				recommendation = "Algorithm is PQC-safe; no action required."
			} else if len(info.PQCReplacement) > 0 {
				recommendation = "Migrate to: " + info.PQCReplacement[0]
			} else {
				recommendation = "Review and replace with a quantum-safe alternative."
			}
		} else {
			recommendation = "Unknown algorithm; manual review recommended."
		}

		band := QRSToBand(g.maxQRS)
		pScore := PriorityScore(g.maxQRS, g.effort, g.occurrences)

		items = append(items, ActionPlanItem{
			Algorithm:      algo,
			DisplayName:    displayName,
			RiskBand:       band,
			Occurrences:    g.occurrences,
			Effort:         g.effort,
			PriorityScore:  pScore,
			Recommendation: recommendation,
			Replacement:    replacement,
		})
	}

	// Sort by priority score descending
	sort.Slice(items, func(i, j int) bool {
		if items[i].PriorityScore != items[j].PriorityScore {
			return items[i].PriorityScore > items[j].PriorityScore
		}
		// Tie-break by algorithm name for deterministic output
		return items[i].Algorithm < items[j].Algorithm
	})

	return items
}

// effortOrdinal returns a numeric ordinal for an effort level for comparison.
func effortOrdinal(e EffortLevel) int {
	switch e {
	case EffortEasy:
		return 0
	case EffortMedium:
		return 1
	case EffortHard:
		return 2
	case EffortBlocked:
		return 3
	default:
		return 1
	}
}
