package detector

import "math"

// ComputeQRS calculates the Quantum Risk Score for a finding based on the
// algorithm's base QRS, key size adjustments (for asymmetric/signature/key-agreement
// families), and occurrence frequency.
func ComputeQRS(algo AlgorithmInfo, keySize int, occurrenceCount int) int {
	score := algo.BaseQRS

	// Key size adjustment (only for asymmetric/signature/key-agreement)
	if keySize > 0 && (algo.Family == "asymmetric-encryption" || algo.Family == "signature" || algo.Family == "key-agreement") {
		switch algo.Name {
		case "RSA":
			if keySize <= 512 {
				score = min(100, score+10)
			}
			if keySize >= 3072 {
				score -= 10
			}
			if keySize >= 4096 {
				score -= 20
			}
		case "ECDSA", "ECDH", "ECC":
			if keySize <= 192 {
				score = min(100, score+5)
			}
			if keySize >= 384 {
				score -= 5
			}
			if keySize >= 521 {
				score -= 10
			}
		}
	}

	// Frequency adjustment
	switch {
	case occurrenceCount >= 21:
		score = min(100, score+10)
	case occurrenceCount >= 6:
		score = min(100, score+5)
	case occurrenceCount >= 2:
		score = min(100, score+2)
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

// QRSToBand maps a Quantum Risk Score value to its corresponding risk band.
func QRSToBand(qrs int) RiskBand {
	switch {
	case qrs >= 80:
		return BandCritical
	case qrs >= 60:
		return BandHigh
	case qrs >= 40:
		return BandMedium
	case qrs >= 20:
		return BandLow
	default:
		return BandSafe
	}
}

// AggregateQRS computes the weighted mean QRS across all findings. Each score
// is weighted by w = (qrs/100)^2 + 0.01, so higher-risk findings dominate the
// aggregate. Returns 0 for an empty input slice.
func AggregateQRS(qrsValues []int) int {
	if len(qrsValues) == 0 {
		return 0
	}

	var weightedSum, totalWeight float64
	for _, qrs := range qrsValues {
		q := float64(qrs)
		w := (q/100.0)*(q/100.0) + 0.01
		weightedSum += q * w
		totalWeight += w
	}
	if totalWeight == 0 {
		return 0
	}
	return int(math.Round(weightedSum / totalWeight))
}

// ClassifyEffort determines the migration effort and rationale based on where
// a cryptographic finding was discovered (source type) and which algorithm is
// involved. The classification follows the spec §8.4 effort matrix.
func ClassifyEffort(source SourceType, algorithm string) (EffortLevel, string) {
	canonical, found := NormaliseAlgorithm(algorithm)
	if !found {
		return EffortMedium, "unknown algorithm; manual review recommended"
	}

	info, ok := LookupAlgorithm(canonical)
	if !ok {
		return EffortMedium, "algorithm not in registry; manual review recommended"
	}

	// PQC-safe algorithms need no migration
	if info.PQCSafe {
		return EffortEasy, "algorithm is already PQC-safe; no migration needed"
	}

	switch source {
	case SourceDeps:
		return EffortBlocked, "migration blocked until upstream dependency updates"
	case SourceConfig:
		return EffortEasy, "configuration change; update cipher suite or algorithm setting"
	case SourceCert:
		switch info.Family {
		case "asymmetric-encryption", "signature", "key-agreement":
			return EffortMedium, "certificate re-issuance required with PQC algorithm"
		default:
			return EffortEasy, "certificate hash algorithm can be updated at renewal"
		}
	case SourceCode:
		switch info.Family {
		case "asymmetric-encryption", "key-agreement":
			return EffortHard, "asymmetric algorithm embedded in code; requires architectural changes"
		case "signature":
			return EffortHard, "signature algorithm embedded in code; requires protocol updates"
		case "hash":
			if canonical == "MD5" || canonical == "SHA1" {
				return EffortEasy, "hash function swap; drop-in replacement available"
			}
			return EffortMedium, "hash algorithm change may affect downstream consumers"
		case "symmetric-encryption":
			if canonical == "DES" || canonical == "3DES" || canonical == "RC2" {
				return EffortMedium, "symmetric cipher upgrade; review mode-of-operation compatibility"
			}
			return EffortEasy, "symmetric key size increase; minimal code changes"
		case "stream-cipher":
			return EffortMedium, "stream cipher replacement; review protocol compatibility"
		default:
			return EffortMedium, "code-level change required; review impact"
		}
	default:
		return EffortMedium, "unknown source type; manual review recommended"
	}
}
