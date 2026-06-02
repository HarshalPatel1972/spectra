# Spectra Architecture Whitepaper

## Executive Summary
Spectra is a forensic instrument designed to discover, quantify, and map cryptographic assets across enterprise environments. As the quantum threat (specifically Shor's Algorithm running on a CRQC) approaches, organizations face an unprecedented migration challenge: moving from classical asymmetric cryptography (RSA, ECC, DSA) to Post-Quantum Cryptography (ML-KEM, ML-DSA).

This whitepaper outlines the architectural framework Spectra uses to guarantee deterministic asset discovery and objective risk scoring.

## 1. The Discovery Engine
Spectra does not rely on regular expressions or heuristic text matching. It uses an Abstract Syntax Tree (AST) parsing engine for compiled and interpreted languages (Go, Python, Java, JS/TS, C++, Rust).

When Spectra encounters `crypto/rsa` in Go or `cryptography.hazmat` in Python, it walks the AST to determine exactly how the algorithm is instantiated, the key sizes used, and the operational modes (e.g., ECB vs. GCM).

## 2. Certificate & Manifest Analysis
Beyond source code, Spectra parses X.509 certificates (PEM/DER) to extract signing algorithms and key lengths. It also parses `go.mod`, `package.json`, `pom.xml`, and `requirements.txt` to identify known vulnerable cryptographic dependencies using the National Vulnerability Database (NVD) and GitHub Advisory Database.

## 3. The Quantum Risk Score (QRS)
Spectra introduces the Quantum Risk Score (QRS), an objective 0-100 metric calculated via a composite function:
*   **Base Algorithm Vulnerability (0-100)**: Evaluates Shor's algorithm impact (e.g., RSA = 90) and Grover's algorithm impact (e.g., AES-128 = 20).
*   **Key Size Penalty**: Deducts points for insufficient key lengths based on NIST SP 800-131A Rev 2.
*   **Usage Frequency & Context**: Amplifies the score based on the blast radius within the codebase.

## 4. The CBOM Standard
Spectra natively outputs CycloneDX 1.7 Cryptographic Bill of Materials (CBOM). This guarantees that Spectra's findings can be ingested by any enterprise compliance tool, SIEM, or vulnerability management platform.

## 5. Migration Simulation (Spectra Forge)
Spectra constructs a dependency graph of cryptographic assets (Spectra Atlas). Using this graph, Spectra Forge simulates the blast radius of migrating an algorithm. It generates a multi-wave migration plan that prioritizes high-risk, low-effort assets (e.g., internal APIs) before tackling complex, high-effort assets (e.g., external legacy clients).

## Conclusion
Spectra provides the foundational intelligence required to execute a successful, mathematically sound transition to a post-quantum cryptographic architecture.
