# SPECTRA — Full Build Prompt
### *Cryptographic Asset Discovery & Quantum Risk Intelligence CLI*

> **"See every cipher. Own your migration."**

---

## 0. Preamble — What You Are Building

**Spectra** is a production-grade command-line tool that scans codebases, X.509 certificates,
configuration files, and dependency manifests to discover every cryptographic algorithm in use,
score each finding against quantum-computing threat models, and generate a Cryptographic Bill of
Materials (CBOM) in the CycloneDX v1.7 format.

The output is not a raw dump of grep results. Spectra delivers a **Quantum Risk Score (QRS)**
per finding and in aggregate, a **migration effort classification** (EASY / MEDIUM / HARD /
BLOCKED), a **prioritised action plan**, and machine-readable output suitable for CI/CD
pipeline gates and compliance audits.

**Why this matters:**
Post-quantum cryptography standards now exist (NIST FIPS 203/204/205). Implementations exist.
What organisations are missing is **visibility into what they actually use today**. You cannot
migrate what you cannot see. Spectra is the first tool in every PQC migration — without it,
every subsequent step is blind.

---

## 1. The Innovation (What Makes Spectra Different)

| Feature | Existing Grep-Based Scanners | Spectra |
|---|---|---|
| Algorithm detection | ✓ | ✓ |
| Key size awareness | ✗ | ✓ (adjusts score by key size) |
| Quantum Risk Score (QRS) | ✗ | ✓ (0–100 numeric, per-finding + aggregate) |
| Migration effort rating | ✗ | ✓ (EASY/MEDIUM/HARD/BLOCKED) |
| Prioritised action plan | ✗ | ✓ (sorted by Risk × Frequency ÷ Effort) |
| CycloneDX 1.7 CBOM output | ✗ | ✓ (algorithm families + elliptic curves) |
| CI/CD gate mode | ✗ | ✓ (`--fail-on=critical`, exit codes 0/1/2/3) |
| Dependency-aware detection | ✗ | ✓ (cross-reference known crypto libraries) |
| Multi-scanner correlation | ✗ | Phase 2 |

---

## 2. Verified Latest Dependency Versions

> **RULE**: Never use a version not on this list. If in doubt, `go get module@latest` and
> check the resolved version before pinning.

| Dependency | Module Path | Version |
|---|---|---|
| Go toolchain | — | **1.25.8** |
| CLI framework | `github.com/spf13/cobra` | **v1.10.2** |
| CycloneDX Go library | `github.com/CycloneDX/cyclonedx-go` | **v0.9.2** |
| Terminal styling | `charm.land/lipgloss/v2` | **v2.x (latest)** |
| YAML parsing | `gopkg.in/yaml.v3` | **v3.0.1** |
| Go crypto stdlib | `golang.org/x/crypto` | **latest** |
| CycloneDX spec target | — | **1.7** (ECMA-424 2nd Ed.) |

> **Do not use** `github.com/charmbracelet/lipgloss` (v1); use `charm.land/lipgloss/v2`.
> **Do not use** `github.com/spf13/viper` (heavyweight); use standard `encoding/json` +
> `gopkg.in/yaml.v3` for config parsing.

---

## 3. Project Structure — Every File, Every Directory

```
spectra/
├── .github/
│   └── workflows/
│       └── ci.yml                  # Build, vet, test on push/PR
├── cmd/
│   └── spectra/
│       └── main.go                 # Binary entry point — calls Execute()
├── internal/
│   ├── cli/
│   │   ├── root.go                 # Root cobra command, global flags
│   │   ├── scan.go                 # `spectra scan` subcommand
│   │   └── version.go              # `spectra version` subcommand
│   ├── config/
│   │   └── config.go               # .spectra.yaml loader, CLI flag merger
│   ├── scanner/
│   │   ├── orchestrator.go         # Fan-out worker pool; owns ScanResult
│   │   ├── code.go                 # Source code regex scanner (all languages)
│   │   ├── cert.go                 # X.509 / PEM certificate scanner
│   │   ├── deps.go                 # Dependency manifest scanner
│   │   └── config_scanner.go       # Config/env file scanner
│   ├── detector/
│   │   ├── algorithms.go           # Algorithm registry: name → base QRS + metadata
│   │   ├── patterns.go             # Language-specific regex pattern loader
│   │   ├── risk.go                 # QRS calculation + migration effort classifier
│   │   └── priority.go             # Prioritised action plan builder
│   ├── cbom/
│   │   └── generator.go            # CycloneDX 1.7 CBOM builder using cyclonedx-go
│   ├── report/
│   │   ├── terminal.go             # Lipgloss-styled terminal output
│   │   ├── jsonout.go              # JSON findings output
│   │   └── htmlout.go              # Self-contained HTML report (Chart.js embedded)
│   └── version/
│       └── version.go              # Build-time version injection
├── rules/
│   └── crypto_patterns.yaml        # Declarative pattern rules (see §6)
├── testdata/
│   ├── samples/
│   │   ├── sample_go.go            # Go file with deliberate RSA/SHA1/ECC usage
│   │   ├── sample_python.py        # Python file with deliberate crypto usage
│   │   ├── sample_java.java        # Java file with deliberate crypto usage
│   │   ├── sample_js.js            # JS/TS file with deliberate crypto usage
│   │   ├── sample_weak.pem         # SHA1withRSA certificate for testing
│   │   ├── sample_go.mod           # Go module with crypto deps
│   │   └── sample_requirements.txt # Python requirements with crypto libs
│   └── golden/
│       ├── findings.json           # Expected JSON output for testdata/samples
│       └── cbom.json               # Expected CBOM output for testdata/samples
├── .spectra.yaml.example           # Annotated config file template
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

> **Rule**: No file in `internal/` may import from `cmd/`. No circular imports.
> Every package must be independently testable.

---

## 4. CLI Design — Commands, Flags, Exit Codes

### 4.1 Primary Command: `spectra scan`

```
spectra scan [path] [flags]

Scan a directory or file for quantum-vulnerable cryptographic assets.

Arguments:
  path    Target path to scan (default: current working directory ".")

Flags:
  -o, --output strings      Output formats: terminal,json,cbom,html
                            (default "terminal") Multiple allowed, comma-separated.
      --out-dir string      Directory for file-based outputs (default "./spectra-out")
      --exclude strings     Glob patterns to exclude (e.g. "vendor/**,*.min.js")
      --scanners strings    Comma-separated scanner list: code,cert,deps,config
                            (default: all four)
      --fail-on string      Minimum QRS band to fail exit code: any|high|critical
                            (for CI/CD mode; no output written unless --quiet=false)
  -q, --quiet               Suppress terminal output (useful with --fail-on in CI)
      --no-progress         Disable progress spinner
  -c, --config string       Path to .spectra.yaml (default: ./.spectra.yaml)
      --concurrency int     Parallel scan workers (default: runtime.NumCPU())
  -v, --verbose             Print each file as it is scanned
      --version             Print version and exit
  -h, --help                Print help
```

### 4.2 Secondary Command: `spectra version`

```
spectra version
```
Prints:
```
Spectra v0.1.0 (commit: abc1234, built: 2026-05-31, go: 1.25.8)
```

### 4.3 Exit Codes

| Code | Meaning |
|---|---|
| **0** | No findings, or all findings below `--fail-on` threshold |
| **1** | Findings at or above `--fail-on` threshold (use in CI `if` conditions) |
| **2** | Scanner runtime error (path not found, permission denied, etc.) |
| **3** | Invalid flags or configuration |

### 4.4 Example Invocations

```bash
# Basic terminal scan of current directory
spectra scan .

# Output JSON + CBOM to ./spectra-out/
spectra scan ./myapp --output json,cbom --out-dir ./spectra-out

# CI/CD gate: fail if any CRITICAL findings exist
spectra scan . --fail-on=critical --quiet

# Scan only certificates and dependency manifests
spectra scan . --scanners cert,deps --output html

# Exclude vendored and built assets
spectra scan . --exclude "vendor/**,dist/**,node_modules/**"
```

---

## 5. Algorithm Registry — The Ground Truth

Define this in `internal/detector/algorithms.go` as a Go map. Every scanner references this
registry by algorithm name (a canonical string key).

### 5.1 Algorithm Entry Structure

```go
type AlgorithmInfo struct {
    Name           string           // canonical ID, e.g. "RSA"
    DisplayName    string           // e.g. "RSA (Rivest–Shamir–Adleman)"
    Family         string           // "asymmetric-encryption", "hash", "symmetric-encryption"
    QuantumThreat  ThreatLevel      // CRITICAL / HIGH / MEDIUM / LOW / SAFE
    BaseQRS        int              // 0–100 base score before adjustments
    PQCSafe        bool             // is this algorithm already PQC-safe?
    PQCReplacement []string         // suggested replacements
    CWE            string           // relevant CWE if applicable
}
```

### 5.2 Complete Algorithm Table

| Canonical Name | Family | BaseQRS | Quantum Threat | PQC-Safe | Suggested Replacement |
|---|---|---|---|---|---|
| `RSA` | asymmetric-encryption | 90 | CRITICAL | No | ML-KEM (FIPS 203) |
| `ECDSA` | signature | 90 | CRITICAL | No | ML-DSA (FIPS 204) |
| `ECDH` | key-agreement | 90 | CRITICAL | No | ML-KEM (FIPS 203) |
| `ECC` | generic | 90 | CRITICAL | No | ML-KEM / ML-DSA |
| `DSA` | signature | 90 | CRITICAL | No | ML-DSA (FIPS 204) |
| `DH` | key-agreement | 85 | CRITICAL | No | ML-KEM (FIPS 203) |
| `ElGamal` | asymmetric-encryption | 85 | CRITICAL | No | ML-KEM (FIPS 203) |
| `SHA1` | hash | 70 | HIGH | No | SHA-256 or SHA-3 |
| `MD5` | hash | 80 | HIGH | No | SHA-256 |
| `DES` | symmetric-encryption | 85 | HIGH | No | AES-256-GCM |
| `3DES` | symmetric-encryption | 70 | HIGH | No | AES-256-GCM |
| `RC4` | stream-cipher | 85 | HIGH | No | ChaCha20-Poly1305 |
| `RC2` | symmetric-encryption | 80 | HIGH | No | AES-256-GCM |
| `AES128` | symmetric-encryption | 25 | MEDIUM | Partial | AES-256 (Grover halves key space) |
| `AES192` | symmetric-encryption | 10 | LOW | Partial | AES-256 |
| `AES256` | symmetric-encryption | 0 | SAFE | Yes | — |
| `SHA256` | hash | 0 | SAFE | Yes | — |
| `SHA384` | hash | 0 | SAFE | Yes | — |
| `SHA512` | hash | 0 | SAFE | Yes | — |
| `SHA3` | hash | 0 | SAFE | Yes | — |
| `ChaCha20` | stream-cipher | 0 | SAFE | Yes | — |
| `ML-KEM` | asymmetric-encryption | 0 | SAFE | Yes | — |
| `ML-DSA` | signature | 0 | SAFE | Yes | — |
| `SLH-DSA` | signature | 0 | SAFE | Yes | — |
| `BLAKE2` | hash | 0 | SAFE | Yes | — |

> **Implementation note**: Normalise all algorithm names to canonical form before lookup.
> Accept aliases: `SHA-1`, `sha1`, `SHA1`, `"SHA1"` all → canonical `SHA1`.
> Maintain an alias map in `algorithms.go`.

---

## 6. Pattern Rules — `rules/crypto_patterns.yaml`

The YAML file is the declarative source of truth for code-level detection. `patterns.go`
loads it at startup into a compiled `[]*regexp.Regexp` slice per algorithm per language.

### 6.1 Schema

```yaml
# rules/crypto_patterns.yaml
# Format: algorithm: language: [regex patterns]
# Each pattern is a raw Go regex string matched against file lines.

RSA:
  go:
    - 'rsa\.GenerateKey\b'
    - 'rsa\.GenerateMultiPrimeKey\b'
    - 'rsa\.EncryptPKCS1v15\b'
    - 'rsa\.EncryptOAEP\b'
    - 'rsa\.SignPKCS1v15\b'
    - 'rsa\.SignPSS\b'
    - '"RSA"'
  python:
    - 'from\s+Crypto\.PublicKey\s+import\s+RSA'
    - 'from\s+cryptography\.hazmat\.primitives\.asymmetric\s+import\s+rsa'
    - 'RSA\.generate\b'
    - 'rsa\.generate_private_key\b'
    - 'load_pem_private_key.*RSA'
    - 'algorithms\.RS256'
    - 'algorithms\.RS384'
    - 'algorithms\.RS512'
    - 'algorithms\.PS256'
  java:
    - 'KeyPairGenerator\.getInstance\("RSA"\)'
    - 'Cipher\.getInstance\("RSA'
    - 'RSAKeyPairGenerator\b'
    - 'RSAKeyGenParameterSpec\b'
    - 'RSACryptoServiceProvider\b'
    - 'Signature\.getInstance\(".*RSA"\)'
    - '"RSA"'
  javascript:
    - 'crypto\.createSign\("RSA'
    - 'crypto\.createVerify\("RSA'
    - 'generateKeyPair.*"rsa"'
    - 'jose.*RS[0-9]'
    - '"RSA-OAEP"'
    - '"RSA-PKCS1-v1_5"'
  rust:
    - 'rsa::RsaPrivateKey\b'
    - 'rsa::RsaPublicKey\b'
    - 'rsa::pkcs1\b'
  c_cpp:
    - 'RSA_generate_key\b'
    - 'RSA_new\b'
    - 'RSA_public_encrypt\b'
    - 'EVP_PKEY_RSA\b'
    - 'NID_rsaEncryption\b'
  generic:
    - '(?i)\bRSA[-_]?\d{3,4}\b'

ECDSA:
  go:
    - 'ecdsa\.GenerateKey\b'
    - 'ecdsa\.Sign\b'
    - 'ecdsa\.Verify\b'
    - 'elliptic\.P256\(\)'
    - 'elliptic\.P384\(\)'
    - 'elliptic\.P521\(\)'
  python:
    - 'from\s+cryptography\.hazmat\.primitives\.asymmetric\s+import\s+ec'
    - 'ec\.SECP256R1\b'
    - 'ec\.SECP384R1\b'
    - 'ec\.SECP521R1\b'
    - 'ec\.generate_private_key\b'
    - 'algorithms\.ES256'
    - 'algorithms\.ES384'
    - 'algorithms\.ES512'
  java:
    - 'KeyPairGenerator\.getInstance\("EC"\)'
    - 'ECGenParameterSpec\b'
    - 'ECDSA\b'
    - '"EC"\s*[,\)]'
  javascript:
    - 'generateKeyPair.*"ec"'
    - '"ECDSA"'
    - '"ECDH"'
    - '"P-256"'
    - '"P-384"'
    - '"P-521"'
    - 'jose.*ES[0-9]'
  rust:
    - 'p256::\b'
    - 'p384::\b'
    - 'k256::\b'
  c_cpp:
    - 'EC_KEY_new_by_curve_name\b'
    - 'ECDSA_sign\b'
    - 'NID_X9_62_prime256v1\b'
    - 'NID_secp384r1\b'

SHA1:
  go:
    - '"crypto/sha1"'
    - 'sha1\.New\(\)'
    - 'sha1\.Sum\b'
  python:
    - 'hashlib\.sha1\b'
    - 'SHA1\(\)'
    - 'hashes\.SHA1\b'
    - 'algorithms\.HS1'
  java:
    - 'MessageDigest\.getInstance\("SHA-1"\)'
    - 'MessageDigest\.getInstance\("SHA1"\)'
    - '"SHA1withRSA"'
    - '"SHA1withECDSA"'
    - '"SHA1withDSA"'
  javascript:
    - 'crypto\.createHash\("sha1"\)'
    - '"SHA-1"'
  rust:
    - 'sha1::Sha1\b'
  c_cpp:
    - 'SHA1_Init\b'
    - 'SHA1_Update\b'
    - 'EVP_sha1\(\)'
  generic:
    - '(?i)\bSHA[-_]?1\b'

MD5:
  go:
    - '"crypto/md5"'
    - 'md5\.New\(\)'
    - 'md5\.Sum\b'
  python:
    - 'hashlib\.md5\b'
    - 'MD5\(\)'
    - 'hashes\.MD5\b'
  java:
    - 'MessageDigest\.getInstance\("MD5"\)'
    - '"MD5withRSA"'
  javascript:
    - 'crypto\.createHash\("md5"\)'
    - '"MD5"'
  rust:
    - 'md5::Md5\b'
  c_cpp:
    - 'MD5_Init\b'
    - 'EVP_md5\(\)'
  generic:
    - '(?i)\bMD[-_]?5\b'

DES:
  go:
    - '"crypto/des"'
    - 'des\.NewCipher\b'
    - 'des\.NewTripleDESCipher\b'
  python:
    - 'from\s+Crypto\.Cipher\s+import\s+DES\b'
    - 'DES\.new\b'
  java:
    - 'Cipher\.getInstance\("DES'
    - 'Cipher\.getInstance\("DESede'
    - '"DESede"'
    - 'SecretKeyFactory\.getInstance\("DES'
  c_cpp:
    - 'EVP_des_\b'
    - 'EVP_des_ede3\b'
  generic:
    - '(?i)\b3DES\b'
    - '(?i)\bTripleDES\b'
    - '(?i)\bDESede\b'

RC4:
  go:
    - '"golang.org/x/crypto/rc4"'
    - 'rc4\.NewCipher\b'
  python:
    - 'from\s+Crypto\.Cipher\s+import\s+ARC4'
    - 'ARC4\.new\b'
  java:
    - 'Cipher\.getInstance\("RC4"\)'
    - 'Cipher\.getInstance\("ARCFOUR"\)'
  generic:
    - '(?i)\bRC[-_]?4\b'
    - '(?i)\barcfour\b'

DH:
  go:
    - '"golang.org/x/crypto/dh"'
    - 'dh\.GenerateKey\b'
  java:
    - 'KeyPairGenerator\.getInstance\("DH"\)'
    - 'DHParameterSpec\b'
    - 'KeyAgreement\.getInstance\("DH"\)'
  generic:
    - '(?i)\bDiffie.Hellman\b'
    - '(?i)\bDHKeyPair\b'

AES128:
  go:
    - 'aes\.NewCipher\b'   # note: combined with key-size check
  python:
    - 'AES\.new.*16\)'     # 16 bytes = 128 bits
  generic:
    - '(?i)\bAES[-_]?128\b'
    - '(?i)\bAES[-_]128[-_]'

ML-KEM:
  generic:
    - '(?i)\bML[-_]?KEM\b'
    - '(?i)\bKyber\b'
    - '(?i)\bFIPS[-_]?203\b'

ML-DSA:
  generic:
    - '(?i)\bML[-_]?DSA\b'
    - '(?i)\bDilithium\b'
    - '(?i)\bFIPS[-_]?204\b'

SLH-DSA:
  generic:
    - '(?i)\bSLH[-_]?DSA\b'
    - '(?i)\bSPHINCS\b'
    - '(?i)\bFIPS[-_]?205\b'
```

> **Implementation note**: The `generic` language key is applied to ALL files regardless of
> extension — useful for config values and comments. Language-specific patterns are only applied
> to files with matching extensions. Define the extension → language map in `scanner/code.go`.

### 6.2 Extension → Language Mapping (in `scanner/code.go`)

```go
var extensionToLanguage = map[string]string{
    ".go":   "go",
    ".py":   "python",
    ".java": "java",
    ".kt":   "java",      // Kotlin: same Java API patterns
    ".scala":"java",
    ".js":   "javascript",
    ".ts":   "javascript",
    ".tsx":  "javascript",
    ".jsx":  "javascript",
    ".rs":   "rust",
    ".c":    "c_cpp",
    ".cpp":  "c_cpp",
    ".cc":   "c_cpp",
    ".h":    "c_cpp",
    ".cs":   "csharp",
    ".rb":   "ruby",
    ".php":  "php",
    ".swift":"swift",
}
```

---

## 7. Scanner Specifications

### 7.1 Orchestrator (`scanner/orchestrator.go`)

Owns the top-level `ScanResult` type and fans out to individual scanners.

```go
type Finding struct {
    ID              string         // UUID v4
    Algorithm       string         // canonical name from algorithm registry
    AlgorithmInfo   AlgorithmInfo  // full metadata from registry
    Source          SourceType     // CODE | CERT | DEPS | CONFIG
    FilePath        string         // relative to scan root
    LineNumber      int            // 0 if not applicable
    LineContent     string         // trimmed matched line (redacted if secret-like)
    Language        string         // detected language
    KeySize         int            // bits, if detectable; 0 if unknown
    Context         string         // extra context (e.g. cert subject, dep version)
    QRS             int            // Quantum Risk Score 0–100
    RiskBand        RiskBand       // CRITICAL|HIGH|MEDIUM|LOW|SAFE
    MigrationEffort EffortLevel    // EASY|MEDIUM|HARD|BLOCKED
    EffortRationale string
    Timestamp       time.Time
}

type ScanResult struct {
    ScanRoot        string
    StartedAt       time.Time
    CompletedAt     time.Time
    TotalFiles      int
    ScannedFiles    int
    SkippedFiles    int
    Findings        []Finding
    AggregateQRS    int              // weighted average across all findings
    FindingsByBand  map[RiskBand]int // count per severity band
    ActionPlan      []ActionItem     // sorted prioritised migration steps
}
```

**Worker pool**: Use `errgroup.Group` (from `golang.org/x/sync/errgroup`) with a semaphore
channel capped at `--concurrency` value. Each goroutine processes one file. Collect results
via a mutex-protected slice.

### 7.2 Code Scanner (`scanner/code.go`)

1. Walk the target directory recursively using `filepath.WalkDir`.
2. Skip directories: `.git`, `node_modules`, `vendor`, `dist`, `build`, `__pycache__`,
   `.idea`, `.vscode`, directories matching any `--exclude` glob.
3. Skip binary files: check first 512 bytes for NUL characters.
4. For each text file with a known extension, load language patterns from the registry.
5. Always apply `generic` patterns regardless of extension.
6. Scan line-by-line. For each matched line:
   - Look up the algorithm in the registry.
   - Attempt key size extraction via a secondary pattern (e.g. `2048`, `4096`, `256`
     near the match within ±3 lines).
   - Create a `Finding`.
7. **Deduplication**: if the same algorithm appears 3+ consecutive lines in the same block
   (e.g. a test file repeating the same call), collapse into one finding with
   `occurrenceCount` incremented.

### 7.3 Certificate Scanner (`scanner/cert.go`)

Files to scan: `*.pem`, `*.crt`, `*.cer`, `*.der`, `*.p12`, `*.pfx`, `*.jks` (JKS is
binary; skip unless `--scanners cert` explicitly requested and the file begins with `FEED`
magic bytes — JKS support is Phase 2).

For PEM files:
1. Read the file, decode all PEM blocks using `encoding/pem`.
2. For each `CERTIFICATE` block, parse with `crypto/x509.ParseCertificate`.
3. Extract:
   - `cert.PublicKeyAlgorithm` → map to canonical algorithm name.
   - `cert.SignatureAlgorithm` → detect SHA1 or MD5 signatures.
   - Key size: for `*rsa.PublicKey`, use `.N.BitLen()`; for `*ecdsa.PublicKey`, use
     `key.Curve.Params().BitSize`.
   - `cert.NotAfter` → flag if expired or expiring within 90 days (INFO finding).
   - Subject CN and SANs for context.

Algorithm mapping:
```go
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
}
```

### 7.4 Dependency Scanner (`scanner/deps.go`)

Parse the following manifest files and flag known cryptographic dependencies.

**Manifest files to look for:**

| File | Package Manager | Parser |
|---|---|---|
| `go.mod` | Go modules | Line-scan for `require` blocks |
| `go.sum` | Go modules | Line-scan |
| `package.json` | npm/yarn | `encoding/json` |
| `yarn.lock` | yarn | Line-scan |
| `requirements.txt` | pip | Line-scan |
| `Pipfile` | pipenv | `gopkg.in/yaml.v3` |
| `pyproject.toml` | poetry | Line-scan |
| `pom.xml` | Maven | `encoding/xml` (partial) |
| `build.gradle` / `build.gradle.kts` | Gradle | Line-scan regex |
| `Cargo.toml` | Cargo | Line-scan |
| `composer.json` | Composer | `encoding/json` |

**Known cryptographic dependency map** (embed in `deps.go`):

```go
var knownCryptoDeps = map[string]DependencyInfo{
    // Go
    "golang.org/x/crypto":            {Algorithms: []string{"various"}, Risk: "check-usage"},
    "github.com/square/go-jose":       {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "github.com/dvsekhvalnov/jose2go":  {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},

    // npm
    "node-forge":     {Algorithms: []string{"RSA","DES","3DES","RC4","MD5","SHA1"}, Risk: "CRITICAL"},
    "jsencrypt":      {Algorithms: []string{"RSA"}, Risk: "CRITICAL"},
    "crypto-js":      {Algorithms: []string{"MD5","SHA1","DES","3DES","RC4"}, Risk: "HIGH"},
    "elliptic":       {Algorithms: []string{"ECDSA","ECDH"}, Risk: "HIGH"},
    "bn.js":          {Algorithms: []string{"RSA","DH"}, Risk: "HIGH"},
    "forge":          {Algorithms: []string{"RSA","DES","MD5"}, Risk: "CRITICAL"},
    "jose":           {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},

    // Python
    "pycryptodome":      {Algorithms: []string{"RSA","DES","3DES","RC4","MD5","SHA1"}, Risk: "HIGH"},
    "pycrypto":          {Algorithms: []string{"RSA","DES","RC4","MD5","SHA1"}, Risk: "CRITICAL"},
    "pyOpenSSL":         {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "paramiko":          {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "cryptography":      {Algorithms: []string{"RSA","ECC","various"}, Risk: "check-usage"},
    "PyJWT":             {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},

    // Java (Maven artifact IDs)
    "bcprov-jdk15on":     {Algorithms: []string{"RSA","ECDSA","DES","MD5","SHA1"}, Risk: "HIGH"},
    "bcpkix-jdk15on":     {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "nimbus-jose-jwt":    {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "java-jwt":           {Algorithms: []string{"RSA","ECDSA"}, Risk: "HIGH"},
    "commons-codec":      {Algorithms: []string{"MD5","SHA1"}, Risk: "HIGH"},

    // Rust (crate names)
    "rsa":      {Algorithms: []string{"RSA"}, Risk: "CRITICAL"},
    "p256":     {Algorithms: []string{"ECDSA","ECDH"}, Risk: "HIGH"},
    "p384":     {Algorithms: []string{"ECDSA","ECDH"}, Risk: "HIGH"},
    "k256":     {Algorithms: []string{"ECDSA","ECDH"}, Risk: "HIGH"},
    "sha1":     {Algorithms: []string{"SHA1"}, Risk: "HIGH"},
    "md5":      {Algorithms: []string{"MD5"}, Risk: "HIGH"},
    "des":      {Algorithms: []string{"DES"}, Risk: "CRITICAL"},
}
```

### 7.5 Config Scanner (`scanner/config_scanner.go`)

Scan: `*.yaml`, `*.yml`, `*.json`, `*.toml`, `*.env`, `*.ini`, `*.conf`, `*.properties`,
`*.xml`, `.env*`, `*.cfg`.

Use line-by-line regex matching (do not fully parse the document; scanning for values
is sufficient). Patterns to detect:

```go
var configPatterns = []ConfigPattern{
    {Regex: `(?i)algorithm["\s:=]+["']?(\w[\w\-]*)`, Group: 1},
    {Regex: `(?i)cipher["\s:=]+["']?(\w[\w\-]*)`, Group: 1},
    {Regex: `(?i)hash[_\-]?algorithm["\s:=]+["']?(\w[\w\-]*)`, Group: 1},
    {Regex: `(?i)signature[_\-]?algorithm["\s:=]+["']?(\w[\w\-]*)`, Group: 1},
    {Regex: `(?i)key[_\-]?size["\s:=]+["']?(\d+)`, Group: 1},
    {Regex: `(?i)tls[_\-]?version["\s:=]+["']?([\d\.]+|TLSv?\d[\.\d]*)`, Group: 1},
    {Regex: `(?i)min[_\-]?version["\s:=]+["']?([\d\.]+|TLSv?\d[\.\d]*)`, Group: 1},
    {Regex: `(?i)(SHA[-_]?1|MD[-_]?5|DES|3DES|RC4|RSA[-_]?\d*|ECDSA)\b`, Group: 0},
}
```

After extracting a candidate value, normalise it and look it up in the algorithm registry.
Only create a finding if the lookup matches a known algorithm.

---

## 8. Quantum Risk Scoring Engine (`detector/risk.go`)

### 8.1 QRS Formula

```go
func ComputeQRS(algo AlgorithmInfo, keySize int, occurrenceCount int) int {
    score := algo.BaseQRS

    // Key size adjustment (only for asymmetric)
    if keySize > 0 && algo.Family == "asymmetric-encryption" || algo.Family == "signature" {
        switch algo.Name {
        case "RSA":
            if keySize <= 512  { score = min(100, score+10) }
            if keySize >= 3072 { score -= 10 }
            if keySize >= 4096 { score -= 20 }
        case "ECDSA", "ECDH", "ECC":
            if keySize <= 192  { score = min(100, score+5) }
            if keySize >= 384  { score -= 5 }
            if keySize >= 521  { score -= 10 }
        }
    }

    // Frequency adjustment (max +10 for widely used algorithms)
    switch {
    case occurrenceCount >= 21: score = min(100, score+10)
    case occurrenceCount >= 6:  score = min(100, score+5)
    case occurrenceCount >= 2:  score = min(100, score+2)
    }

    if score < 0 { score = 0 }
    if score > 100 { score = 100 }
    return score
}
```

### 8.2 Risk Band Mapping

```go
func QRSToBand(qrs int) RiskBand {
    switch {
    case qrs >= 80: return CRITICAL
    case qrs >= 60: return HIGH
    case qrs >= 40: return MEDIUM
    case qrs >= 20: return LOW
    default:        return SAFE
    }
}
```

### 8.3 Aggregate QRS

```go
// AggregateQRS is the weighted mean, giving higher weight to critical findings.
func AggregateQRS(findings []Finding) int {
    if len(findings) == 0 { return 0 }
    totalWeight := 0.0
    weightedSum := 0.0
    for _, f := range findings {
        w := math.Pow(float64(f.QRS)/100.0, 2) + 0.01 // weight heavier findings more
        weightedSum += w * float64(f.QRS)
        totalWeight += w
    }
    return int(math.Round(weightedSum / totalWeight))
}
```

### 8.4 Migration Effort Classifier

```go
func ClassifyEffort(f Finding) (EffortLevel, string) {
    switch f.Source {
    case DEPS:
        return BLOCKED, "Algorithm used inside a third-party dependency; " +
            "remediation requires upgrading or replacing the package."
    case CERT:
        if f.Algorithm == "RSA" || f.Algorithm == "ECDSA" {
            return HARD, "Certificate replacement requires PKI coordination: " +
                "CA re-issuance, trust store updates, and client compatibility verification."
        }
        return MEDIUM, "Certificate re-issuance needed."
    case CONFIG:
        return EASY, "Algorithm reference in configuration file; " +
            "update the configuration value and validate the consuming service."
    case CODE:
        switch f.Algorithm {
        case "SHA1", "MD5":
            return EASY, "Drop-in hash function replacement with no interface change."
        case "AES128":
            return EASY, "Increase AES key size from 128 to 256 bits."
        case "DES", "3DES", "RC4", "RC2":
            return MEDIUM, "Symmetric cipher replacement; update key management and IV/nonce handling."
        default: // RSA, ECDSA, ECDH, DH
            return HARD, "Asymmetric algorithm replacement requires new key types, " +
                "updated serialisation formats, and potential protocol changes."
        }
    }
    return MEDIUM, ""
}
```

### 8.5 Priority Score (for Action Plan ordering)

```go
// PriorityScore: higher = fix first.
// Balances risk against ease of fixing.
func PriorityScore(qrs int, effort EffortLevel, occurrences int) float64 {
    effortMap := map[EffortLevel]float64{
        EASY: 1.0, MEDIUM: 0.6, HARD: 0.3, BLOCKED: 0.1,
    }
    freqBoost := math.Log1p(float64(occurrences))
    return float64(qrs) * effortMap[effort] * (1 + freqBoost*0.05)
}
```

---

## 9. Output Formats

### 9.1 Terminal Output (`report/terminal.go`)

Use `charm.land/lipgloss/v2`.

**Palette**:
```go
const (
    ColorCritical = "#FF4444"
    ColorHigh     = "#FF8C00"
    ColorMedium   = "#FFD700"
    ColorLow      = "#4ADE80"
    ColorSafe     = "#38BDF8"
    ColorDim      = "#6B7280"
    ColorBold     = "#F8FAFC"
)
```

**Terminal output structure:**
1. **Header block**: "SPECTRA" ASCII wordmark + version + scan root + timestamp.
2. **Progress spinner** (while scanning): "Scanning… [N files]"
3. **Summary table** (post-scan):
   - Columns: ALGORITHM | SOURCE | FILE | LINE | KEY SIZE | QRS | BAND | EFFORT
   - Rows sorted by QRS descending.
   - Band column coloured by severity.
   - Maximum 50 rows in terminal view; `… and N more findings. Run with --output json for all.`
4. **Aggregate block**:
   - "Aggregate QRS: 73/100 (HIGH)"
   - Count per band: "CRITICAL: 3 | HIGH: 12 | MEDIUM: 7 | LOW: 2 | SAFE: 1"
5. **Top 5 Action Items** from prioritised plan.
6. **Output paths** if --output includes file formats.

Do not use `fmt.Println` for any user-facing output. All terminal output goes through
`report/terminal.go`.

### 9.2 JSON Output (`report/jsonout.go`)

Write to `<out-dir>/spectra-findings.json`. Use `encoding/json` with `json.MarshalIndent`.

Schema: serialise the full `ScanResult` struct. Ensure all time fields use RFC3339.

### 9.3 CBOM Output (`report/cbom.go` + `cbom/generator.go`)

Target: CycloneDX **v1.7** using `github.com/CycloneDX/cyclonedx-go` v0.9.2.

Map each unique `(Algorithm, KeySize)` pair to one `cdx.Component` of
`Type: cdx.ComponentTypeCryptographicAsset`. Set:
- `CryptoProperties.AssetType` = `algorithm`
- `CryptoProperties.AlgorithmProperties.Primitive` = appropriate primitive
- `CryptoProperties.AlgorithmProperties.ParameterSetIdentifier` = key size if known
- `CryptoProperties.AlgorithmProperties.ExecutionEnvironment` = `"software"`
- `Properties` = Spectra-specific metadata (QRS, effort, occurrence count)

Write to `<out-dir>/spectra-cbom.json`. Validate with CycloneDX CLI if available.

### 9.4 HTML Report (`report/htmlout.go`)

Generate a **self-contained** HTML file (all CSS and JS inline; no external CDN calls).
Use Chart.js embedded as a minified `<script>` block (download at build time; do not
fetch at runtime).

**Report sections:**
1. **Header**: Spectra logo + scan metadata (root, duration, timestamp, version).
2. **Executive Summary**: Aggregate QRS gauge chart (0–100), counts per band.
3. **Risk Distribution**: Doughnut chart of findings by algorithm family.
4. **Top Findings Table**: Sortable HTML table (all findings).
5. **Action Plan**: Numbered list sorted by priority score.
6. **Algorithm Inventory**: Summary table of unique algorithms found.

**Design**: Dark theme. Monospace accents for code snippets. Colour-coded severity badges.
Professional enough to share with a CISO or compliance officer.

---

## 10. Configuration File (`.spectra.yaml`)

```yaml
# .spectra.yaml — Spectra configuration
# All settings can be overridden by CLI flags.

scan:
  # Glob patterns to exclude from scanning (appended to built-in defaults)
  exclude:
    - "vendor/**"
    - "node_modules/**"
    - "dist/**"
    - "*.min.js"
    - "testdata/**"

  # File extensions to scan for code patterns (defaults to built-in list)
  # extensions:
  #   - ".go"
  #   - ".py"

  # Scanners to enable: code, cert, deps, config
  scanners: [code, cert, deps, config]

  # Worker concurrency (0 = use runtime.NumCPU())
  concurrency: 0

output:
  # Default output formats
  formats: [terminal, cbom]

  # Output directory for file-based reports
  dir: "./spectra-out"

ci:
  # Fail exit code if findings at or above this band exist
  # Values: any | high | critical
  fail_on: "critical"
```

---

## 11. Phase Breakdown and Atomic Commits

### Phase 1 — MVP (scope for initial agent session)

Every item below is **one atomic commit**. Follow Conventional Commits strictly.

```
chore: initialize Go 1.25 module (module: github.com/HarshalPatel1972/spectra)
chore: add Makefile with build, test, vet, lint targets
feat(cli): add root cobra command with global flags
feat(cli): add version subcommand with build-time injection
feat(cli): add scan subcommand skeleton with all flags
feat(config): implement .spectra.yaml loader with flag merger
feat(rules): add crypto_patterns.yaml with RSA/ECDSA/SHA1/MD5/DES/RC4/DH patterns
feat(detector): implement algorithm registry and alias normaliser
feat(scanner): implement code file walker with extension filter and binary skip
feat(scanner): implement code pattern matcher — Go and Python patterns
feat(scanner): implement code pattern matcher — Java and JavaScript patterns
feat(scanner): implement code pattern matcher — Rust and C/C++ patterns
feat(scanner): implement X.509 certificate scanner (PEM + DER)
feat(scanner): implement dependency manifest scanner (go.mod, package.json, requirements.txt)
feat(scanner): implement config file scanner (yaml/json/env/toml)
feat(scanner): implement orchestrator worker pool with errgroup
feat(detector): implement quantum risk scoring engine (QRS formula)
feat(detector): implement migration effort classifier (EASY/MEDIUM/HARD/BLOCKED)
feat(detector): implement priority score and action plan builder
feat(report): implement terminal renderer with lipgloss v2 (table + summary)
feat(report): implement JSON output writer
feat(cbom): implement CycloneDX 1.7 CBOM generator
feat(report): implement self-contained HTML report generator
feat(cli): add --fail-on CI/CD gate with exit codes 0/1/2/3
test: add testdata sample files for all languages and cert types
test: add unit tests for QRS computation and band mapping
test: add unit tests for code pattern matcher (golden file comparison)
test: add unit tests for certificate scanner
test: add unit tests for dependency manifest scanner
test: add unit tests for CBOM generator output against golden fixture
ci: add GitHub Actions workflow (go 1.25, build + vet + test matrix)
docs: write README.md with usage, architecture diagram (Mermaid), and examples
```

### Phase 2 — Post-MVP Enhancements

- Cross-source correlation engine (link code + cert + dep findings by algorithm)
- TLS endpoint scanner (`spectra scan --url=https://example.com`)
- Git blame integration (when was the vulnerable crypto introduced?)
- `spectra diff` — compare two CBOM files to track migration progress
- Container image scanning (parse OCI layer tarballs)
- Optional web dashboard (Next.js, served by embedded Go server)

---

## 12. Testing Requirements

### 12.1 Unit Tests

Every function in `internal/detector/` must have table-driven unit tests with at least:
- 1 CRITICAL finding case
- 1 SAFE/PQC finding case
- 1 edge case (empty input, zero key size, unknown algorithm)

Pattern matchers: test against `testdata/samples/` files using golden outputs in
`testdata/golden/`. Use `github.com/google/go-cmp/cmp` for structured diffs.

Certificate scanner: test with `testdata/samples/sample_weak.pem` (SHA1withRSA, RSA-1024).
Generate this certificate in `testdata/generate.go` using `crypto/x509` so it never
expires and the test suite always passes.

### 12.2 Integration Tests

One end-to-end test in `internal/scanner/orchestrator_test.go` that:
1. Calls `ScanDirectory(testdata/samples)`.
2. Verifies aggregate QRS > 70.
3. Verifies at least one CRITICAL finding for RSA.
4. Verifies at least one finding with Source = CERT.
5. Verifies CBOM output is valid JSON parseable by `cyclonedx-go`.

### 12.3 CI Validation

The GitHub Actions workflow must:
1. Run `go vet ./...` — zero warnings enforced.
2. Run `go test ./... -race -count=1` — all tests pass.
3. Build for `linux/amd64`, `darwin/arm64`, `windows/amd64` using `GOOS/GOARCH`.
4. Produce release binaries on tag push via `goreleaser`.

---

## 13. GitHub Actions CI (`/.github/workflows/ci.yml`)

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go: ["1.25.x"]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
          cache: true

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test ./... -race -count=1 -coverprofile=coverage.out

      - name: Coverage report
        run: go tool cover -func=coverage.out

  build:
    runs-on: ubuntu-latest
    needs: test
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: darwin
            goarch: arm64
          - goos: windows
            goarch: amd64

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25.x"
          cache: true
      - name: Build
        run: |
          GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} \
          go build -ldflags="-X github.com/HarshalPatel1972/spectra/internal/version.Version=$(git describe --tags --always --dirty)" \
          -o dist/spectra-${{ matrix.goos }}-${{ matrix.goarch }} \
          ./cmd/spectra
```

---

## 14. README Requirements

The README must include:

1. **Badge row**: CI status, Go version, License (MIT), CycloneDX v1.7.
2. **Tagline**: "See every cipher. Own your migration."
3. **One-paragraph problem statement** (from the research context above).
4. **Quick start** (3 commands: install, scan, view report).
5. **CLI reference** (table of all flags).
6. **Architecture diagram** (Mermaid flowchart: CLI → Orchestrator → 4 Scanners → Detector
   → Risk Scorer → 4 Output Formats).
7. **Output format examples** (terminal screenshot in code block, CBOM JSON snippet).
8. **QRS explanation** table (0–100 ranges → band → meaning).
9. **Roadmap** (Phase 2 items).
10. **Contributing** section.

---

## 15. go.mod Skeleton

```go
module github.com/HarshalPatel1972/spectra

go 1.25

require (
    charm.land/lipgloss/v2                    v2.x.x
    github.com/CycloneDX/cyclonedx-go         v0.9.2
    github.com/spf13/cobra                    v1.10.2
    golang.org/x/crypto                       latest
    golang.org/x/sync                         latest
    gopkg.in/yaml.v3                          v3.0.1
)
```

> After `go mod tidy`, all versions will resolve to exact semver. Pin them in `go.sum`.

---

## 16. Deployment

- **Binary distribution**: GitHub Releases via `goreleaser`. No server required.
- **Homebrew tap** (Phase 2): `brew install harshalpatel1972/tap/spectra`.
- **npm package** (Phase 2): Wrap the binary with a thin npm shim for JS ecosystem users.
- **GitHub Action** (Phase 2): `harshalpatel1972/spectra-action@v1` — wraps
  `spectra scan . --fail-on=critical --output cbom` in one Action step.

---

## 17. Non-Negotiable Quality Rules

1. **No `fmt.Println` in library code.** All user output goes through `report/terminal.go`.
2. **No `panic` outside of `main`.** Return errors; let `cobra` handle exit codes.
3. **No init() functions** except for the algorithm registry pre-compilation of regexes.
4. **Every exported function has a Go doc comment.**
5. **No global mutable state.** Pass config and registries as function arguments or struct fields.
6. **Sensitive data protection**: If a matched line looks like a private key or secret
   (matches `PRIVATE KEY`, `SECRET`, `PASSWORD`, `TOKEN`), redact it to `[REDACTED]`
   in all outputs including the CBOM.
7. **Zero external network calls** during a scan. Spectra is fully offline-capable.
8. **Deterministic output**: given the same input, the same JSON/CBOM must be produced.
   Sort findings by `(QRS desc, FilePath asc, LineNumber asc)` before serialisation.

---

## 18. Summary

| Property | Value |
|---|---|
| Project Name | **Spectra** |
| Tagline | "See every cipher. Own your migration." |
| Language | Go 1.25.8 |
| CLI Framework | Cobra v1.10.2 |
| CBOM Format | CycloneDX 1.7 (ECMA-424 2nd Ed.) |
| Terminal Styling | Lipgloss v2 |
| Output Formats | terminal, JSON, CBOM, HTML |
| Primary innovation | Quantum Risk Score + Migration Effort + CI gate |
| Repository | `github.com/HarshalPatel1972/spectra` |
| License | MIT |
| Phase 1 deliverable | Single binary, 4 scanners, 4 output formats, CI pipeline |

---

*Build prompt authored: May 2026 — based on NIST PQC standards (FIPS 203/204/205),
CycloneDX 1.7 specification, and the Spectra project research findings.*
