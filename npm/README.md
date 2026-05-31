# Spectra

[![Build Status](https://github.com/HarshalPatel1972/spectra/actions/workflows/ci.yml/badge.svg)](https://github.com/HarshalPatel1972/spectra/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CycloneDX 1.7](https://img.shields.io/badge/CycloneDX-1.7-blue)](https://cyclonedx.org)

**Cryptographic Asset Discovery & Quantum Risk Intelligence CLI**

> "See every cipher. Own your migration."

Spectra is a production-grade command-line tool that scans codebases, X.509 certificates, configuration files, and dependency manifests to discover every cryptographic algorithm in use, score each finding against quantum-computing threat models, and generate a Cryptographic Bill of Materials (CBOM) in the CycloneDX format.

## Architecture

```text
[Input: Repo, Image, TLS URL] 
 └── Orchestrator
      ├── Code Scanner
      ├── Cert Scanner
      ├── Deps Scanner
      └── Config Scanner
           └── Detector Engine
                └── QRS & Priority Scoring
                     ├── Terminal Report
                     ├── JSON Output
                     ├── CBOM Generator
                     └── Web Dashboard
```

## Quick Start

### Install via NPM
```bash
npm install -g spectra-crypto-cli
```

### Usage
```bash
# Basic terminal scan of current directory
spectra scan .

# Scan an OCI Container Image
spectra scan --image alpine:latest

# Scan a TLS Endpoint
spectra scan --url example.com:443

# Start the interactive web dashboard
spectra dashboard
```

## Output Formats

### Terminal Output
Spectra provides a rich, ANSI-colored terminal UI showing findings, aggregate risks, and a prioritized action plan.

### CBOM Output
Generates CycloneDX formatted Cryptographic Bill of Materials.
```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.6",
  "version": 1,
  "components": [
    {
      "type": "cryptographic-asset",
      "name": "RSA",
      "version": "2048",
      "description": "RSA (2048-bit)",
      "properties": [
        { "name": "spectra:qrs", "value": "50" }
      ]
    }
  ]
}
```

## Quantum Risk Score (QRS)

Spectra evaluates algorithms using a proprietary 0-100 scoring model:

| Score Range | Risk Band | Meaning |
| :--- | :--- | :--- |
| **80-100** | CRITICAL | Broken algorithm (MD5) or highly vulnerable to Shor's algorithm |
| **60-79** | HIGH | Weak algorithm (SHA1, RSA-1024) |
| **40-59** | MEDIUM | Acceptable for now, but not quantum-safe (RSA-2048) |
| **20-39** | LOW | Strong algorithm (AES-256, SHA-3) |
| **0-19** | SAFE | Post-Quantum Cryptography (ML-KEM, ML-DSA) |

## CLI Reference

| Command | Description | Example |
| :--- | :--- | :--- |
| `spectra scan [dir]` | Scans a directory | `spectra scan ./src` |
| `spectra scan --image` | Scans a container image | `spectra scan --image ubuntu:latest` |
| `spectra scan --url` | Scans a TLS endpoint | `spectra scan --url google.com:443` |
| `spectra dashboard` | Starts the embedded Web Dashboard | `spectra dashboard --port 8080` |
| `spectra diff` | Compares two CBOMs to track progress | `spectra diff old.json new.json` |
| `spectra version` | Displays the current version | `spectra version` |

## Roadmap
- [x] Phase 1: MVP Core Scanners
- [x] Phase 2: Git Blame Integration, TLS Scanning, Container Scanning, Web Dashboard
- [ ] Phase 3: Advanced CI/CD integrations and API Server

## Contributing
Contributions are welcome! Please open an issue before submitting a large PR. Ensure you run `go test ./...` and `golangci-lint` before committing.
