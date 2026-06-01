# Quantum Drift - Monthly Newsletter

**Subject:** What you missed in PQC this month 🔐

Hey everyone,

Welcome to this month's edition of **Quantum Drift**. Here is your monthly update on post-quantum cryptography and what we're building at Spectra.

### This month in PQC:
- **News:** [Insert NSA CNSA 2.0 or NIST latest announcement]
- **Industry Migration:** [Insert story about a major enterprise migrating to PQC]
- **Algorithm Development:** [Insert update on ML-KEM or ML-DSA]

### Spectra Updates:
- **What shipped in vX.Y.Z:** We just launched the WebAssembly Playground. You can now scan your codebase for quantum vulnerabilities entirely in your browser. Code never leaves your machine!
- **What's coming:** We're working on expanding our detection rules for Rust and C++ macro-based cryptography.

### Scan of the Month:
This month we ran Spectra on **[Project Name]**.
- **Findings:** [N] total findings detected.
- **Insight:** We found a surprising amount of legacy SHA-1 usage buried inside test utilities that were inadvertently shipped to production.
- **Report:** Check out the full CBOM here: [Link]

### One thing to do this month:
Run a baseline scan on your primary API gateway repository.
\`spectra scan . --output html\`

Stay secure,
Harshal
