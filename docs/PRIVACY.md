# Privacy Policy

Spectra is built with a security-first, privacy-by-default architecture. We recognize that cryptographic material, even metadata, is highly sensitive.

## Core Privacy Tenets

1. **Local Execution**: The Spectra CLI binary runs entirely locally on your machine or CI/CD runner. It does not send source code, configuration files, certificates, or scan metadata to any external server.
2. **Zero Telemetry**: We do not collect crash reports, usage analytics, or IP addresses from CLI usage.
3. **Database Locality**: The internal SQLite database (`~/.spectra/state.db`) is stored locally and never synced to the cloud.

## Web Playground Privacy

If you use the online playground at `spectra.tools/playground`:
- **Server Mode**: Code submitted is written to a temporary ephemeral directory on our Vercel serverless function, scanned, and instantly deleted. We do not persist or log your code.
- **WASM Mode**: The entire scanning engine is compiled to WebAssembly and runs natively inside your browser. Your code never leaves your device.

If you have concerns about privacy, please review our open-source codebase.
