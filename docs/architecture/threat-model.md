# Spectra Threat Model

## 1. System Overview
Spectra is a local-first command-line application that analyzes source code, certificates, and configuration files to discover cryptographic assets. 

## 2. Trust Boundaries
*   **Execution Environment**: Spectra executes within the user's local terminal or CI/CD runner. It assumes the execution environment is trusted.
*   **Source Code**: Spectra assumes the source code it is scanning is trusted (i.e., not actively malicious code designed to exploit the AST parser).
*   **External APIs**: Spectra makes NO external API calls during a standard scan. It does not send telemetry, analytics, or code snippets to any external server.

## 3. Potential Threats & Mitigations

### 3.1 Denial of Service (DoS) via Malicious Source Files
*   **Threat**: An attacker includes a deeply nested or infinitely recursive source file designed to exhaust memory or CPU during AST parsing.
*   **Mitigation**: Spectra enforces strict depth limits during AST traversal and implements execution timeouts for individual files.

### 3.2 Information Disclosure
*   **Threat**: Spectra inadvertently leaks sensitive proprietary code or hardcoded secrets in its output reports.
*   **Mitigation**: Spectra strictly limits its output to the *names* of cryptographic algorithms, their key sizes, and file line numbers. It does not output the surrounding source code logic or extract string literals (keys/secrets).

### 3.3 False Negatives (Missed Vulnerabilities)
*   **Threat**: Spectra fails to identify a vulnerable cryptographic implementation (e.g., a custom or highly obfuscated implementation of RSA).
*   **Mitigation**: Spectra is designed to detect standard library and common third-party cryptographic framework usage. It explicitly states that it does not detect custom, obfuscated, or non-standard cryptographic implementations.

### 3.4 Supply Chain Attacks on Spectra
*   **Threat**: The Spectra binary or its dependencies are compromised.
*   **Mitigation**: Spectra's own codebase is continuously scanned. All releases are signed, and we provide reproducible build instructions. We publish a CycloneDX CBOM and standard SBOM for every release.

## 4. Out of Scope
*   Dynamic analysis of running applications.
*   Detection of hardcoded secrets or passwords (use a dedicated secret scanner).
*   Verification of the mathematical correctness of cryptographic implementations.
