# Spectra GitHub Action

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Run Spectra Cryptographic Intelligence Platform directly in your CI/CD pipelines.**

This GitHub Action wraps the Spectra CLI, allowing you to automatically scan your repository for cryptographic assets, evaluate Quantum Risk Scores (QRS), enforce compliance (CNSA 2.0, NIST), and generate CBOMs on every Pull Request.

## Quick Start

Add the following step to your `.github/workflows/` YAML file:

```yaml
name: Cryptographic Audit
on: [push, pull_request]

jobs:
  spectra-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Required for Git Blame temporal tracking

      - name: Run Spectra Scan
        uses: HarshalPatel1972/spectra-action@main
        with:
          fail-on-critical: 'true'
          generate-cbom: 'true'
```

## Inputs

| Input | Description | Default |
| --- | --- | --- |
| `scan-path` | Directory to scan | `.` |
| `fail-on-critical` | Fail the build if CRITICAL findings are detected | `false` |
| `generate-cbom` | Generate a CycloneDX CBOM | `true` |
| `output-format` | Output format (`terminal`, `json`, `cbom`) | `terminal` |

## Artifacts

If you enable CBOM generation or HTML reporting, you can upload them as build artifacts:

```yaml
      - name: Upload CBOM Artifact
        uses: actions/upload-artifact@v4
        with:
          name: spectra-cbom
          path: spectra-cbom.json
```

## Phase 2 Capabilities Supported

This action natively supports all Spectra Phase 2 capabilities:
- **Temporal Tracking**: Git blame is extracted automatically (ensure you use `fetch-depth: 0` in checkout).
- **Compliance Rules**: Evaluates against built-in rulesets.
- **Dependency & Config Scanning**: Deep scanning across manifests and infra-as-code files.
