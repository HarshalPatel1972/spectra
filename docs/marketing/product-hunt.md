# Product Hunt Launch

## Tagline
Spectra – The enterprise scanner for quantum-vulnerable cryptography.

## Description
Most organizations know post-quantum cryptography is coming, but few know what classical cryptography they are actually using. Spectra solves this. It's an open-source, local-first CLI that scans codebases, certificates, and dependencies for quantum-vulnerable cryptography like RSA, ECC, and SHA-1.

Spectra calculates a Quantum Risk Score, generates a CycloneDX Cryptographic Bill of Materials (CBOM), and maps findings against NSA CNSA 2.0 deadlines.

No data leaves your machine. Try it now in the WASM playground.

## First Comment (Maker Comment)
Hi hunters! 👋 

I built Spectra after reading through NIST's PQC migration guides and realizing that most developers have no idea what cryptography is actually running in their infrastructure. We have SAST tools for bugs and SCA tools for CVEs, but nothing specifically designed for cryptographic asset discovery.

Spectra is built in Go. It's extremely fast and supports scanning 5+ languages, X.509 certs, and Docker images. 

I'd love your feedback! Let me know what algorithms or frameworks we should add next.
