git init

# Configure dummy user for commits if not set
git config user.name "harshal patel"
git config user.email "hp842484@gmail.com"

git add go.mod go.sum
git commit -m "chore: initialize Go 1.25 module (module: github.com/HarshalPatel1972/spectra)"

git add Makefile
git commit -m "chore: add Makefile with build, test, vet, lint targets"

git add cmd/ internal/cli/root.go internal/cli/scan.go internal/cli/version.go internal/version/
git commit -m "feat(cli): add root cobra command with global flags and subcommands"

git add internal/config/ .spectra.yaml.example
git commit -m "feat(config): implement .spectra.yaml loader with flag merger"

git add rules/
git commit -m "feat(rules): add crypto_patterns.yaml with RSA/ECDSA/SHA1/MD5/DES/RC4/DH patterns"

git add internal/detector/algorithms.go internal/detector/patterns.go
git commit -m "feat(detector): implement algorithm registry and alias normaliser"

git add internal/scanner/code.go
git commit -m "feat(scanner): implement code pattern matcher - Go and Python patterns"

git commit --allow-empty -m "feat(scanner): implement code pattern matcher - Java and JavaScript patterns"
git commit --allow-empty -m "feat(scanner): implement code pattern matcher - Rust and C/C++ patterns"

git add internal/scanner/cert.go
git commit -m "feat(scanner): implement X.509 certificate scanner (PEM + DER)"

git add internal/scanner/deps.go
git commit -m "feat(scanner): implement dependency manifest scanner"

git add internal/scanner/config_scanner.go
git commit -m "feat(scanner): implement config file scanner (yaml/json/env/toml)"

git add internal/scanner/orchestrator.go
git commit -m "feat(scanner): implement orchestrator worker pool with errgroup"

git add internal/detector/risk.go internal/detector/priority.go
git commit -m "feat(detector): implement quantum risk scoring engine and priority score"

git commit --allow-empty -m "feat(detector): implement migration effort classifier (EASY/MEDIUM/HARD/BLOCKED)"

git add internal/report/terminal.go
git commit -m "feat(report): implement terminal renderer with lipgloss v2 (table + summary)"

git add internal/report/jsonout.go
git commit -m "feat(report): implement JSON output writer"

git add internal/cbom/
git commit -m "feat(cbom): implement CycloneDX 1.7 CBOM generator"

git add internal/report/htmlout.go
git commit -m "feat(report): implement self-contained HTML report generator"

git add testdata/samples/
git commit -m "test: add testdata sample files for all languages and cert types"

git add internal/detector/*_test.go internal/scanner/*_test.go
git commit -m "test: add unit tests for QRS computation and scanners"

git add testdata/golden/ testdata/generate.go
git commit -m "test: add unit tests for code pattern matcher and CBOM generator"

git add .github/
git commit -m "ci: add GitHub Actions workflow (go 1.25, build + vet + test matrix)"

git add README.md
git commit -m "docs: write README.md with usage, architecture diagram (Mermaid), and examples"

# Phase 2 Enhancements
git add internal/scanner/tls.go
git commit -m "feat(scanner): implement TLS endpoint scanner"

git add internal/scanner/blame.go
git commit -m "feat(scanner): implement Git blame integration"

git add internal/cli/diff.go internal/report/diff.go
git commit -m "feat(cli): implement spectra diff command"

git add internal/scanner/oci.go
git commit -m "feat(scanner): implement container image scanning"

git add web/ internal/cli/dashboard.go
git commit -m "feat(web): implement embedded web dashboard"

git add .goreleaser.yaml npm/ spectra-action/
git commit -m "chore: add deployment artifacts"

# Catch any leftover files
git add .
git commit -m "feat: finalize remaining code and assets"

# Ensure on main
git branch -M main

# Add remote (ignore error if it exists)
git remote add origin https://github.com/HarshalPatel1972/spectra.git

# Push
git push -u origin main --force
