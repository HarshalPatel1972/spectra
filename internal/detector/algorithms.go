// Package detector provides the cryptographic algorithm registry, quantum risk
// scoring engine, migration effort classifier, and priority action plan builder.
package detector

// ThreatLevel represents the quantum computing threat level of an algorithm.
type ThreatLevel string

const (
	// ThreatCritical indicates the algorithm is fully broken by quantum computers.
	ThreatCritical ThreatLevel = "CRITICAL"
	// ThreatHigh indicates the algorithm is significantly weakened by quantum computers.
	ThreatHigh ThreatLevel = "HIGH"
	// ThreatMedium indicates the algorithm has reduced security margins under quantum attack.
	ThreatMedium ThreatLevel = "MEDIUM"
	// ThreatLow indicates the algorithm has minimal quantum vulnerability.
	ThreatLow ThreatLevel = "LOW"
	// ThreatSafe indicates the algorithm is considered quantum-safe.
	ThreatSafe ThreatLevel = "SAFE"
)

// RiskBand represents the severity band derived from a Quantum Risk Score.
type RiskBand string

const (
	// BandCritical represents QRS >= 80.
	BandCritical RiskBand = "CRITICAL"
	// BandHigh represents QRS 60-79.
	BandHigh RiskBand = "HIGH"
	// BandMedium represents QRS 40-59.
	BandMedium RiskBand = "MEDIUM"
	// BandLow represents QRS 20-39.
	BandLow RiskBand = "LOW"
	// BandSafe represents QRS 0-19.
	BandSafe RiskBand = "SAFE"
)

// EffortLevel represents the migration effort classification.
type EffortLevel string

const (
	// EffortEasy indicates a drop-in replacement with minimal effort.
	EffortEasy EffortLevel = "EASY"
	// EffortMedium indicates moderate effort requiring some code changes.
	EffortMedium EffortLevel = "MEDIUM"
	// EffortHard indicates significant effort requiring architectural changes.
	EffortHard EffortLevel = "HARD"
	// EffortBlocked indicates remediation is blocked by external dependencies.
	EffortBlocked EffortLevel = "BLOCKED"
)

// SourceType represents where a cryptographic finding was discovered.
type SourceType string

const (
	// SourceCode indicates the finding was discovered in source code.
	SourceCode SourceType = "CODE"
	// SourceCert indicates the finding was discovered in a certificate.
	SourceCert SourceType = "CERT"
	// SourceDeps indicates the finding was discovered in a dependency manifest.
	SourceDeps SourceType = "DEPS"
	// SourceConfig indicates the finding was discovered in a configuration file.
	SourceConfig SourceType = "CONFIG"
)

// AlgorithmInfo holds metadata about a cryptographic algorithm including its
// quantum risk profile, PQC safety status, and recommended replacements.
type AlgorithmInfo struct {
	Name           string      // canonical ID, e.g. "RSA"
	DisplayName    string      // e.g. "RSA (Rivest–Shamir–Adleman)"
	Family         string      // "asymmetric-encryption", "hash", "symmetric-encryption", etc.
	QuantumThreat  ThreatLevel // CRITICAL / HIGH / MEDIUM / LOW / SAFE
	BaseQRS        int         // 0–100 base score before adjustments
	PQCSafe        bool        // is this algorithm already PQC-safe?
	PQCReplacement []string    // suggested replacements
	CWE            string      // relevant CWE if applicable
}

// algorithmRegistry is the ground truth for all known cryptographic algorithms.
var algorithmRegistry = map[string]AlgorithmInfo{
	"RSA": {
		Name: "RSA", DisplayName: "RSA (Rivest–Shamir–Adleman)",
		Family: "asymmetric-encryption", QuantumThreat: ThreatCritical,
		BaseQRS: 90, PQCSafe: false,
		PQCReplacement: []string{"ML-KEM (FIPS 203)"},
		CWE: "CWE-327",
	},
	"ECDSA": {
		Name: "ECDSA", DisplayName: "ECDSA (Elliptic Curve Digital Signature Algorithm)",
		Family: "signature", QuantumThreat: ThreatCritical,
		BaseQRS: 90, PQCSafe: false,
		PQCReplacement: []string{"ML-DSA (FIPS 204)"},
		CWE: "CWE-327",
	},
	"ECDH": {
		Name: "ECDH", DisplayName: "ECDH (Elliptic Curve Diffie–Hellman)",
		Family: "key-agreement", QuantumThreat: ThreatCritical,
		BaseQRS: 90, PQCSafe: false,
		PQCReplacement: []string{"ML-KEM (FIPS 203)"},
		CWE: "CWE-327",
	},
	"ECC": {
		Name: "ECC", DisplayName: "ECC (Elliptic Curve Cryptography)",
		Family: "generic", QuantumThreat: ThreatCritical,
		BaseQRS: 90, PQCSafe: false,
		PQCReplacement: []string{"ML-KEM (FIPS 203)", "ML-DSA (FIPS 204)"},
		CWE: "CWE-327",
	},
	"DSA": {
		Name: "DSA", DisplayName: "DSA (Digital Signature Algorithm)",
		Family: "signature", QuantumThreat: ThreatCritical,
		BaseQRS: 90, PQCSafe: false,
		PQCReplacement: []string{"ML-DSA (FIPS 204)"},
		CWE: "CWE-327",
	},
	"DH": {
		Name: "DH", DisplayName: "DH (Diffie–Hellman)",
		Family: "key-agreement", QuantumThreat: ThreatCritical,
		BaseQRS: 85, PQCSafe: false,
		PQCReplacement: []string{"ML-KEM (FIPS 203)"},
		CWE: "CWE-327",
	},
	"ElGamal": {
		Name: "ElGamal", DisplayName: "ElGamal Encryption",
		Family: "asymmetric-encryption", QuantumThreat: ThreatCritical,
		BaseQRS: 85, PQCSafe: false,
		PQCReplacement: []string{"ML-KEM (FIPS 203)"},
		CWE: "CWE-327",
	},
	"SHA1": {
		Name: "SHA1", DisplayName: "SHA-1",
		Family: "hash", QuantumThreat: ThreatHigh,
		BaseQRS: 70, PQCSafe: false,
		PQCReplacement: []string{"SHA-256", "SHA-3"},
		CWE: "CWE-328",
	},
	"MD5": {
		Name: "MD5", DisplayName: "MD5",
		Family: "hash", QuantumThreat: ThreatHigh,
		BaseQRS: 80, PQCSafe: false,
		PQCReplacement: []string{"SHA-256"},
		CWE: "CWE-328",
	},
	"DES": {
		Name: "DES", DisplayName: "DES (Data Encryption Standard)",
		Family: "symmetric-encryption", QuantumThreat: ThreatHigh,
		BaseQRS: 85, PQCSafe: false,
		PQCReplacement: []string{"AES-256-GCM"},
		CWE: "CWE-327",
	},
	"3DES": {
		Name: "3DES", DisplayName: "Triple DES (3DES/DESede)",
		Family: "symmetric-encryption", QuantumThreat: ThreatHigh,
		BaseQRS: 70, PQCSafe: false,
		PQCReplacement: []string{"AES-256-GCM"},
		CWE: "CWE-327",
	},
	"RC4": {
		Name: "RC4", DisplayName: "RC4 (Rivest Cipher 4)",
		Family: "stream-cipher", QuantumThreat: ThreatHigh,
		BaseQRS: 85, PQCSafe: false,
		PQCReplacement: []string{"ChaCha20-Poly1305"},
		CWE: "CWE-327",
	},
	"RC2": {
		Name: "RC2", DisplayName: "RC2 (Rivest Cipher 2)",
		Family: "symmetric-encryption", QuantumThreat: ThreatHigh,
		BaseQRS: 80, PQCSafe: false,
		PQCReplacement: []string{"AES-256-GCM"},
		CWE: "CWE-327",
	},
	"AES128": {
		Name: "AES128", DisplayName: "AES-128",
		Family: "symmetric-encryption", QuantumThreat: ThreatMedium,
		BaseQRS: 25, PQCSafe: false,
		PQCReplacement: []string{"AES-256"},
		CWE: "",
	},
	"AES192": {
		Name: "AES192", DisplayName: "AES-192",
		Family: "symmetric-encryption", QuantumThreat: ThreatLow,
		BaseQRS: 10, PQCSafe: false,
		PQCReplacement: []string{"AES-256"},
		CWE: "",
	},
	"AES256": {
		Name: "AES256", DisplayName: "AES-256",
		Family: "symmetric-encryption", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"SHA256": {
		Name: "SHA256", DisplayName: "SHA-256",
		Family: "hash", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"SHA384": {
		Name: "SHA384", DisplayName: "SHA-384",
		Family: "hash", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"SHA512": {
		Name: "SHA512", DisplayName: "SHA-512",
		Family: "hash", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"SHA3": {
		Name: "SHA3", DisplayName: "SHA-3",
		Family: "hash", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"ChaCha20": {
		Name: "ChaCha20", DisplayName: "ChaCha20-Poly1305",
		Family: "stream-cipher", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"ML-KEM": {
		Name: "ML-KEM", DisplayName: "ML-KEM (FIPS 203 / Kyber)",
		Family: "asymmetric-encryption", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"ML-DSA": {
		Name: "ML-DSA", DisplayName: "ML-DSA (FIPS 204 / Dilithium)",
		Family: "signature", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"SLH-DSA": {
		Name: "SLH-DSA", DisplayName: "SLH-DSA (FIPS 205 / SPHINCS+)",
		Family: "signature", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
	"BLAKE2": {
		Name: "BLAKE2", DisplayName: "BLAKE2",
		Family: "hash", QuantumThreat: ThreatSafe,
		BaseQRS: 0, PQCSafe: true,
		PQCReplacement: nil,
	},
}

// aliasMap maps non-canonical algorithm names to their canonical form.
var aliasMap = map[string]string{
	// RSA aliases
	"rsa": "RSA", "Rsa": "RSA",
	"RSA-1024": "RSA", "RSA-2048": "RSA", "RSA-3072": "RSA", "RSA-4096": "RSA",
	"RSA1024": "RSA", "RSA2048": "RSA", "RSA3072": "RSA", "RSA4096": "RSA",
	"rsa1024": "RSA", "rsa2048": "RSA", "rsa3072": "RSA", "rsa4096": "RSA",

	// ECDSA aliases
	"ecdsa": "ECDSA", "Ecdsa": "ECDSA", "EC-DSA": "ECDSA",

	// ECDH aliases
	"ecdh": "ECDH", "Ecdh": "ECDH", "EC-DH": "ECDH",

	// ECC aliases
	"ecc": "ECC", "Ecc": "ECC",

	// DSA aliases
	"dsa": "DSA", "Dsa": "DSA",

	// DH aliases
	"dh": "DH", "Dh": "DH",
	"Diffie-Hellman": "DH", "diffie-hellman": "DH", "DiffieHellman": "DH",

	// ElGamal aliases
	"elgamal": "ElGamal", "ELGAMAL": "ElGamal", "Elgamal": "ElGamal", "el-gamal": "ElGamal",

	// SHA1 aliases
	"sha1": "SHA1", "SHA-1": "SHA1", "sha-1": "SHA1", "Sha1": "SHA1",

	// MD5 aliases
	"md5": "MD5", "Md5": "MD5", "MD-5": "MD5", "md-5": "MD5",

	// DES aliases
	"des": "DES", "Des": "DES",

	// 3DES aliases
	"3des": "3DES", "triple-des": "3DES", "TripleDES": "3DES", "tripledes": "3DES",
	"DESede": "3DES", "desede": "3DES", "DESEDE": "3DES",

	// RC4 aliases
	"rc4": "RC4", "Rc4": "RC4", "RC-4": "RC4", "rc-4": "RC4",
	"ARCFOUR": "RC4", "arcfour": "RC4", "ARC4": "RC4", "arc4": "RC4",

	// RC2 aliases
	"rc2": "RC2", "Rc2": "RC2", "RC-2": "RC2", "rc-2": "RC2",

	// AES aliases
	"AES-128": "AES128", "aes-128": "AES128", "aes128": "AES128", "AES_128": "AES128",
	"AES-192": "AES192", "aes-192": "AES192", "aes192": "AES192", "AES_192": "AES192",
	"AES-256": "AES256", "aes-256": "AES256", "aes256": "AES256", "AES_256": "AES256",

	// SHA2 family aliases
	"sha256": "SHA256", "SHA-256": "SHA256", "sha-256": "SHA256", "Sha256": "SHA256",
	"sha384": "SHA384", "SHA-384": "SHA384", "sha-384": "SHA384", "Sha384": "SHA384",
	"sha512": "SHA512", "SHA-512": "SHA512", "sha-512": "SHA512", "Sha512": "SHA512",

	// SHA3 aliases
	"sha3": "SHA3", "SHA-3": "SHA3", "sha-3": "SHA3", "Sha3": "SHA3",

	// ChaCha20 aliases
	"chacha20": "ChaCha20", "CHACHA20": "ChaCha20", "chacha20-poly1305": "ChaCha20",
	"ChaCha20Poly1305": "ChaCha20", "CHACHA20-POLY1305": "ChaCha20",

	// PQC aliases
	"ML_KEM": "ML-KEM", "ml-kem": "ML-KEM", "ml_kem": "ML-KEM", "MLKEM": "ML-KEM",
	"Kyber": "ML-KEM", "kyber": "ML-KEM", "KYBER": "ML-KEM",
	"FIPS-203": "ML-KEM", "FIPS203": "ML-KEM", "fips203": "ML-KEM",

	"ML_DSA": "ML-DSA", "ml-dsa": "ML-DSA", "ml_dsa": "ML-DSA", "MLDSA": "ML-DSA",
	"Dilithium": "ML-DSA", "dilithium": "ML-DSA", "DILITHIUM": "ML-DSA",
	"FIPS-204": "ML-DSA", "FIPS204": "ML-DSA", "fips204": "ML-DSA",

	"SLH_DSA": "SLH-DSA", "slh-dsa": "SLH-DSA", "slh_dsa": "SLH-DSA", "SLHDSA": "SLH-DSA",
	"SPHINCS": "SLH-DSA", "sphincs": "SLH-DSA", "SPHINCS+": "SLH-DSA",
	"FIPS-205": "SLH-DSA", "FIPS205": "SLH-DSA", "fips205": "SLH-DSA",

	// BLAKE2 aliases
	"blake2": "BLAKE2", "Blake2": "BLAKE2", "BLAKE2b": "BLAKE2", "BLAKE2s": "BLAKE2",
	"blake2b": "BLAKE2", "blake2s": "BLAKE2",
}

// NormaliseAlgorithm converts an algorithm name to its canonical form using
// the alias map. Returns the canonical name and true if found, or the
// original name and false if unknown.
func NormaliseAlgorithm(name string) (string, bool) {
	// Check if already canonical
	if _, ok := algorithmRegistry[name]; ok {
		return name, true
	}
	// Check alias map
	if canonical, ok := aliasMap[name]; ok {
		return canonical, true
	}
	return name, false
}

// LookupAlgorithm returns the AlgorithmInfo for a given canonical or aliased
// algorithm name. Returns the info and true if found.
func LookupAlgorithm(name string) (AlgorithmInfo, bool) {
	canonical, found := NormaliseAlgorithm(name)
	if !found {
		return AlgorithmInfo{}, false
	}
	info, ok := algorithmRegistry[canonical]
	return info, ok
}

// AllAlgorithms returns a copy of all algorithms in the registry.
func AllAlgorithms() map[string]AlgorithmInfo {
	result := make(map[string]AlgorithmInfo, len(algorithmRegistry))
	for k, v := range algorithmRegistry {
		result[k] = v
	}
	return result
}
