# Maintenance Commitment

Spectra is an active, commercially-backed open-source project. This document outlines our commitments regarding the long-term maintenance and viability of the project.

## Core Commitments

1.  **Rule Viability**: We commit to updating Spectra's cryptographic detection rules to match the latest iterations of NIST FIPS and NSA CNSA standards within 30 days of their official publication.
2.  **Backwards Compatibility**: We commit to semantic versioning. The CLI interface, JSON output format, and CBOM generation will not introduce breaking changes without a major version bump.
3.  **Local-First Guarantee**: We commit that the open-source CLI will always remain local-first. We will never introduce mandatory telemetry or remote code-upload requirements into the core execution engine.

## Governance

Spectra is currently maintained by the core Spectra Tools team. We actively welcome community contributions for:
*   New language AST parsers.
*   Framework-specific cryptographic signatures.
*   Output formatting plugins.

## Sunset Policy

In the unlikely event that the core team can no longer maintain Spectra:
*   We will announce the sunsetting phase at least 6 months in advance.
*   We will ensure the repository remains publicly accessible in a read-only state.
*   We will actively seek to transfer ownership to a reputable foundation (e.g., the Cloud Native Computing Foundation or OpenSSF).
