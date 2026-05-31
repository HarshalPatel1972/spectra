# Spectra Threat Model

Spectra is a static analysis and discovery tool, not an active defense mechanism. This document outlines the threats Spectra mitigates and the threats it does not cover.

## Threats in Scope (Mitigated)

1. **Unknown Cryptographic Debt**: Organizations using outdated algorithms (e.g., MD5, SHA-1, DES) without realizing where they exist.
2. **Quantum Vulnerability**: Organizations relying on asymmetric algorithms (RSA, ECC, Diffie-Hellman) that will be broken by Cryptographically Relevant Quantum Computers (CRQCs).
3. **Compliance Violations**: Codebases inadvertently violating NSA CNSA 2.0, NIST SP 800-131A, or PCI DSS v4.0 cryptographic standards.
4. **Weak Key Sizes**: The use of approved algorithms with insufficiently long keys (e.g., RSA-1024).

## Threats Out of Scope

1. **Implementation Flaws**: Spectra cannot detect if an AES-256 implementation is subject to timing attacks, side-channel leaks, or uses a hardcoded initialization vector (IV).
2. **Key Management Issues**: Spectra detects the presence of keys and certificates, but cannot evaluate if your key rotation policies or HSM integrations are secure.
3. **Runtime Encryption Evasion**: Spectra performs static analysis. It cannot detect cryptography dynamically loaded or decrypted at runtime (e.g., packed malware).
4. **Malicious Insider Bypassing**: Developers can intentionally obfuscate cryptographic calls to bypass Spectra's regex patterns. Spectra assumes good faith from the scanned codebase.
