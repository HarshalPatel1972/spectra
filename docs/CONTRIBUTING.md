# Contributing to Spectra

First off, thank you for considering contributing to Spectra! The goal of this project is to build the definitive open-source Cryptographic Intelligence Platform.

## How to Contribute

### 1. Adding Cryptographic Patterns
The core of Spectra's detection engine relies on regular expressions defined in `rules/crypto_patterns.yaml`. If you discover a cryptographic API or pattern that Spectra misses, please submit a Pull Request adding it to this file. 
- Ensure you test the pattern using a regex debugger.
- Add a test case to the `testdata/` directory to prove the pattern works.

### 2. Filing Bugs and False Positives
If Spectra misidentifies a safe algorithm as vulnerable (False Positive) or completely misses a vulnerable algorithm (False Negative), please use the specific Issue Templates provided in the `.github/ISSUE_TEMPLATE` directory.

### 3. Submitting Pull Requests
- Fork the repository and create a branch from `main`.
- Write tests for any new features or bug fixes.
- Run `go test ./...` and ensure all tests pass.
- Run `golangci-lint run` to ensure your code matches our style guidelines.
- Your commit messages should follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

## Development Setup

```bash
git clone https://github.com/HarshalPatel1972/spectra.git
cd spectra
go mod download
go test ./...
```
