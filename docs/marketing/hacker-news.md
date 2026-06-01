# Hacker News Launch

**Title:** Show HN: Spectra — scans codebases for quantum-vulnerable cryptography (RSA, ECC, SHA-1)

**Body:**
Hey HN,

I built Spectra after reading through NIST's PQC migration guides and realizing most codebases have no idea what cryptography they're using. There's no systematic way to answer "where is RSA in my infrastructure?"

Spectra scans source code (Go, Python, Java, JS, Rust, C/C++), X.509 certificates, config files, and dependency manifests. It outputs:
- A Quantum Risk Score (0-100) per finding and aggregate
- A CycloneDX 1.7 CBOM
- A migration action plan sorted by risk × effort
- Compliance gaps against CNSA 2.0 and NIST SP 800-131A

brew install harshalpatel1972/tap/spectra
spectra scan .

Source: https://github.com/HarshalPatel1972/spectra

The most interesting technical challenge was the QRS formula — how to weight algorithm vulnerability, key size, and usage frequency into a single actionable number. Happy to discuss the methodology.
