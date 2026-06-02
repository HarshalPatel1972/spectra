<img src="./docs/assets/spectra-logo.svg" alt="Spectra Logo" width="240" />

### See every cipher. Own your migration.

[![CI](https://github.com/HarshalPatel1972/spectra/actions/workflows/ci.yml/badge.svg)](https://github.com/HarshalPatel1972/spectra/actions/workflows/ci.yml)
[![Go 1.25](https://img.shields.io/badge/go-1.25-blue?style=flat&logo=go)](https://go.dev)
[![License MIT](https://img.shields.io/badge/license-MIT-green?style=flat)](LICENSE)
[![Spectra QRS: 8](https://img.shields.io/badge/Spectra_QRS-8%2F100_SAFE-16A34A?style=flat)](https://spectra-site-psi.vercel.app)

**The forensic instrument that makes an organization's cryptographic landscape visible, scorable, and navigable.**

```bash
$ brew install harshalpatel1972/tap/spectra
$ cd your-project
$ spectra scan .
```

[Documentation](https://github.com/HarshalPatel1972/spectra/tree/main/docs) • [Playground](https://spectra-site-psi.vercel.app/playground) • [Discord](https://discord.gg/spectra)

---

## The Output

```bash
$ spectra scan ./myapp
▓ Scanning 847 files in 12 packages...

CRITICAL  RSA-2048     auth/jwt.go:47              QRS: 90
CRITICAL  RSA-2048     pkg/crypto/key.go:12        QRS: 90
HIGH      SHA-1        legacy/hash_util.go:91      QRS: 70
HIGH      ECDSA/P-256  certs/api.pem               QRS: 85

──────────────────────────────────────────────────────────
Aggregate QRS: 83/100 — CRITICAL
Compliance: 47 gaps with CNSA 2.0

Run spectra simulate --from RSA --to ML-KEM to generate your migration plan.
```

## Why Spectra?

No tool combines high cryptographic specificity with high analytical precision. 

Spectra scans codebases, certificates, and dependencies for quantum-vulnerable cryptography like RSA, ECC, and SHA-1. It tells you what cryptography you use, what it means, and what to do about it.

*   **Code-level detection**: Native parsing for Go, Python, Java, JS, C++, Rust.
*   **Certificate scanning**: Full X.509 parsing for PEM/DER files to detect vulnerable key pairs.
*   **CycloneDX 1.7 CBOMs**: The industry standard Cryptographic Bill of Materials out-of-the-box.
*   **NSA CNSA 2.0 Gap Analysis**: Real-time evaluation against the 2030/2033 migration deadlines.
*   **Local-first Privacy**: No telemetry, no accounts, no analytics. Your code never leaves your machine.

## SPECTRA ATLAS (Relationship Graph)
Understand the blast radius of your cryptography. `spectra graph` maps the exact files relying on specific algorithms so you know what will break when you upgrade.

## SPECTRA FORGE (Migration Simulator)
The quantum transition is too large to do at once. `spectra simulate` produces a multi-wave migration plan prioritizing high-risk, low-effort assets first.

## SPECTRA MERIDIAN (Compliance Engine)
Run `spectra compliance` to map your current cryptographic posture against NIST SP 800-131A Rev 2, NSA CNSA 2.0, and PCI-DSS v4.0.

---

### Trust
We scan Spectra itself. Our own QRS is **8/100** (deliberate low-risk test fixtures). Every finding Spectra produces links directly to the NIST, NSA, or IETF standard that defines it as a vulnerability.

### License
MIT License. See [LICENSE](LICENSE) for details.
