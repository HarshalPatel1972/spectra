# Security Policy

## Reporting Vulnerabilities

Please report security issues directly to: **security@spectra.tools**

Do **not** file public GitHub issues for security vulnerabilities.
We take security seriously and will respond within 48 hours to coordinate disclosure and a patch.

## Spectra's Security Properties

- **No telemetry.** Spectra never phones home. Zero network calls are made during local scanning.
- **No account required.** No login, no API keys, and no data collection.
- **No data retention.** The playground API deletes submitted code immediately after scanning.
- **Deterministic output.** The same input always produces the exact same output.

## What Spectra Does Not Do

Spectra does **not** prevent cryptographic vulnerabilities from entering your code. It discovers them.
Spectra does **not** store your findings or remediate your code automatically. That is your responsibility.
