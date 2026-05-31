# Spectra GitHub Action

[![Build Status](https://github.com/HarshalPatel1972/spectra/actions/workflows/ci.yml/badge.svg)](https://github.com/HarshalPatel1972/spectra/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official GitHub Action for [Spectra](https://github.com/HarshalPatel1972/spectra).

Run blazing-fast Cryptographic Asset Discovery & Post-Quantum Risk Intelligence scans directly in your CI/CD pipelines. It automatically generates a CycloneDX CBOM (Cryptographic Bill of Materials) and outputs risk reports.

## Usage

Add this step to your GitHub Actions workflow:

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Run Spectra Crypto Scanner
    uses: HarshalPatel1972/spectra-action@main
    with:
      target: '.'            # Optional: Directory to scan (default: '.')
      format: 'json'         # Optional: Output format: terminal, json, both (default: 'terminal')
      qrs_threshold: '40'    # Optional: Fail build if QRS exceeds this (default: '0' = don't fail)
```

## Inputs

| Input | Description | Default |
| --- | --- | --- |
| `target` | Directory to scan for cryptographic assets | `.` |
| `format` | Output format (`terminal`, `json`, `both`) | `terminal` |
| `qrs_threshold` | Fail build if the aggregate QRS score is strictly greater than this value (0-100). | `0` |

## Outputs

The action will automatically drop a `spectra-report.json` and `cbom.json` in the root of your workspace if you enable JSON format, which can be uploaded as workflow artifacts.
