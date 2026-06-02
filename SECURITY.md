# Security Policy

## Supported Versions

We provide security updates for the current major release of Spectra.

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within Spectra, please DO NOT open a public issue.

Instead, please send an email to **security@spectra.tools**. We will respond within 48 hours to acknowledge receipt of the vulnerability report.

### Our Commitment
*   We will rapidly investigate all legitimate reports.
*   We will provide regular updates on our progress toward a fix.
*   We will publicly acknowledge your contribution (if desired) once the vulnerability has been patched and disclosed.

## Scope
Please report issues related to:
*   AST parser crashes resulting in Denial of Service.
*   Bypasses of the local-first execution model (e.g., unintended network calls).
*   Path traversal or arbitrary file read/write vulnerabilities.
*   Dependency vulnerabilities in our published binaries.

Please DO NOT report:
*   False negatives (Spectra failing to detect a specific cipher). These should be reported as standard bug reports or feature requests via GitHub Issues.
*   Theoretical vulnerabilities without a proof of concept.
