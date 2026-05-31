# SPECTRA — Full Build Prompt #3
### *Distribution Architecture, Trust Engineering, and Market Adoption Playbook*

---

> **The product is complete. The category does not yet exist in the developer's mind.**
> Build Prompt #3 is not a feature document. It is a **mind-change document.**
>
> The goal is not to get downloads. The goal is to make a developer feel something —
> specifically: *"My infrastructure is exposed right now, and this tool fixes it."*
> That feeling precedes every install.

---

## 0. The Strategic Frame Before Every Decision

Before specifying any page, feature, or post, answer this question:

**"Does this make a developer feel the problem before we offer the solution?"**

If yes → build it.
If it only describes the solution → cut it or reorder it.

The quantum threat is abstract. RSA is abstract. FIPS 203 is abstract.
The following are NOT abstract:
- Stripe's payment API using RSA-2048 to encrypt card data
- Your company's VPN breaking when a quantum computer arrives
- A hostile nation-state downloading encrypted traffic today to decrypt in 5 years

**Every distribution channel must lead with the concrete before the technical.**

---

## 1. The Four-Layer Distribution Funnel

Before any build spec, define the funnel. Every deliverable in this document maps to one layer:

```
LAYER 1 — AWARENESS
  They have never heard of Spectra. They don't know PQC is urgent.
  Goal: Make the problem visceral. Make it shareable.
  Channels: The "What Happens..." page, social posts, Scan of the Month,
            conference talks, "We scanned X" articles.

LAYER 2 — CONSIDERATION
  They've heard of PQC. They're wondering if they need to act.
  Goal: Remove doubt. Show evidence. Build credibility.
  Channels: Architecture whitepaper, comparison pages, demo repo, case studies,
            benchmark reports, security audit, threat model.

LAYER 3 — ACTIVATION
  They want to try it. Something is stopping them.
  Goal: Remove every friction point between awareness and first result.
  Channels: Playground (zero install), one-line brew install, VS Code extension,
            GitHub App, pre-commit hook, Docker image.

LAYER 4 — RETENTION + ADVOCACY
  They installed it. Now make them keep using it and tell others.
  Goal: Deliver enough value that they tweet, post, or recommend it.
  Channels: Badge economy, CI/CD integration, weekly newsletter,
            Scan of the Month, Discord community, CBOM viewer.
```

---

## 2. Tier 1 — Launch Foundation

### 2.1 Domain Strategy

Check in order (first available wins):

| Priority | Domain | Why |
|---|---|---|
| 1 | `spectra.tools` | Clean, professional, no ambiguity |
| 2 | `spectractl.dev` | Developer-native (ctl = CLI convention) |
| 3 | `usespectra.dev` | Action-oriented, follows `usethis` pattern |
| 4 | `spectra-pqc.dev` | Highly searchable, keyword-rich |

**Register the domain before any public announcement.**

> Do not use `.io` domains for security tools. `.io` had DNSSEC and ownership issues in
> 2023 that damaged the reputation of security tools using them. Use `.dev` or `.tools`.

### 2.2 Website Architecture

Two separate deployable sites. Both on Vercel. Both free tier.

**Site A — Landing + Playground + Blog: `spectra.tools`**
Stack: Next.js 15 + Tailwind CSS + Framer Motion + MDX

**Site B — Documentation: `docs.spectra.tools`**
Stack: VitePress 1.6.4

Both sites share one GitHub repo: `HarshalPatel1972/spectra-site`

```
spectra-site/
├── landing/              # Next.js 15 app (Site A)
│   ├── app/
│   │   ├── page.tsx         # Main landing page
│   │   ├── what-happens/
│   │   │   └── page.tsx     # The Quantum Threat visualization
│   │   ├── playground/
│   │   │   └── page.tsx     # Live scanner playground
│   │   ├── playground-wasm/
│   │   │   └── page.tsx     # Privacy-first WASM playground (Phase 2)
│   │   ├── risk-calculator/
│   │   │   └── page.tsx     # Interactive QRS calculator
│   │   ├── pqc-readiness/
│   │   │   └── page.tsx     # 10-question PQC assessment
│   │   ├── cbom-viewer/
│   │   │   └── page.tsx     # Drag-and-drop CBOM viewer
│   │   ├── compare/
│   │   │   ├── page.tsx     # Comparison hub
│   │   │   ├── vs-grep/
│   │   │   ├── vs-diy/
│   │   │   └── vs-sca/
│   │   ├── roadmap/
│   │   │   └── page.tsx     # Public roadmap
│   │   ├── changelog/
│   │   │   └── page.tsx     # Public changelog
│   │   ├── blog/
│   │   │   ├── page.tsx     # Blog index
│   │   │   └── [...slug]/
│   │   │       └── page.tsx
│   │   ├── security/
│   │   │   └── page.tsx     # SECURITY.md rendered as page
│   │   ├── privacy/
│   │   │   └── page.tsx     # Privacy policy
│   │   └── api/
│   │       ├── scan/
│   │       │   └── route.ts   # Serverless: run Spectra on pasted code
│   │       ├── badge/
│   │       │   └── route.ts   # Serverless: return SVG badge for given QRS score
│   │       └── qrs/
│   │           └── route.ts   # Serverless: compute QRS from algorithm + keysize
│   └── public/
│       ├── spectra.wasm       # Phase 2: compiled Spectra scanner core
│       └── wasm_exec.js       # Phase 2: Go WASM bootstrap shim
│
└── docs/                 # VitePress 1.6.4 site (Site B)
    ├── .vitepress/
    │   └── config.ts
    ├── guide/
    │   ├── index.md         # What is Spectra
    │   ├── getting-started.md
    │   ├── installation.md
    │   ├── quick-start.md
    │   └── configuration.md
    ├── reference/
    │   ├── cli.md           # All commands and flags
    │   ├── algorithms.md    # Algorithm database reference
    │   ├── qrs.md           # QRS scoring methodology
    │   ├── cbom.md          # CBOM format specification
    │   ├── cai.md           # Cryptographic Agility Index
    │   └── cps.md           # Cryptographic Posture Score
    ├── integrations/
    │   ├── github-actions.md
    │   ├── vscode.md
    │   ├── docker.md
    │   ├── pre-commit.md
    │   └── gitlab-ci.md
    ├── compliance/
    │   ├── cnsa-20.md
    │   ├── nist-800-131a.md
    │   ├── pci-dss-4.md
    │   └── fips-140-3.md
    └── contributing/
        └── index.md
```

### 2.3 Landing Page (`landing/app/page.tsx`) — Full Copy Spec

The landing page has one job: get a developer from "interesting" to "installing" in under 60 seconds.

**Section 1 — Hero (Above the Fold)**

```
Headline:   "See every cipher. Own your migration."

Subheadline: "Spectra scans your codebases, certificates, and dependencies
              to find quantum-vulnerable cryptography — before it finds you."

Terminal animation (typed.js or Framer Motion):
  Show a real scan output, character by character:

  $ spectra scan .
  ▓ Scanning 1,247 files...

  CRITICAL  RSA-2048   auth/jwt.go:47          QRS: 90  → Replace with ML-KEM
  CRITICAL  RSA-2048   pkg/crypto/key.go:12    QRS: 90
  HIGH      SHA-1      legacy/hash.go:91       QRS: 70  → Replace with SHA-256
  HIGH      ECDSA/P-256 certs/api-server.pem   QRS: 85
  MEDIUM    AES-128    config/tls.yaml:14      QRS: 25  → Upgrade to AES-256

  ─────────────────────────────────────────────
  Aggregate QRS: 83/100 — CRITICAL
  Compliance: 47 gaps with CNSA 2.0 · 3 gaps with PCI DSS v4.0
  Action plan: 5 waves · Estimated 8 weeks · 2 engineers
  ─────────────────────────────────────────────
  CBOM generated → ./spectra-out/spectra-cbom.json
  HTML report  → ./spectra-out/spectra-report.html

CTA buttons:
  [ ↓ brew install spectra ]    [ Try in Browser → ]

Social proof line (update after launch):
  "Used by engineers at [logos appear as they're acquired]"
```

**Section 2 — The Problem (The Gut-Punch)**

Do NOT write a paragraph. Write a two-column layout:

```
LEFT COLUMN: "The World Your Cryptography Protects"
  Icon + label for each:
  — Payment processing
  — API authentication
  — Code signing
  — VPN tunnels
  — Health records
  — Internal secrets

RIGHT COLUMN: "What a Quantum Computer Does to It"
  Same icons, but greyed out with ❌
  Red banner across them: "Broken in minutes"

Below the two columns:
  "This is not a hypothetical. NSA's CNSA 2.0 requires post-quantum
   algorithms in new NSS systems by 2027 and exclusively by 2033.
   Most codebases have not started."

  [ → See What Happens If Quantum Arrives Tomorrow ]
```

**Section 3 — How It Works**

Three steps, each with a terminal code block:

```
Step 1: Find
  spectra scan ./myapp
  → discovers every algorithm in code, certificates, configs, dependencies

Step 2: Understand
  spectra compliance --frameworks cnsa20,nist-800-131a
  → shows which findings violate which regulatory requirements, by when

Step 3: Act
  spectra simulate --from RSA --to ML-KEM
  → shows the migration plan in ordered waves with effort estimates
```

**Section 4 — Output Formats**

Side-by-side tabs: Terminal / JSON / CBOM / HTML Report / Executive PDF
Each tab shows a real code sample or screenshot.

**Section 5 — Integrations Row**

Icons + "Works with":
GitHub Actions · VS Code · Docker · Pre-commit · GitLab CI · JetBrains · Jenkins

**Section 6 — Why Spectra (vs doing it yourself)**

Three card grid:
```
Card 1: "Not just detection"
  Other tools tell you what algorithm is present.
  Spectra tells you the risk score, migration effort,
  regulatory deadline, and action order.

Card 2: "No data leaves your machine"
  Spectra is a local CLI with no telemetry, no accounts,
  no cloud backend. Your code never leaves your machine.
  [The playground offers a server option and a WASM option.]

Card 3: "Standards-embedded"
  Every finding is cross-referenced against NIST SP 800-131A,
  NSA CNSA 2.0, and PCI DSS v4.0 with exact clause citations —
  not paraphrases.
```

**Section 7 — Community + Adoption**

GitHub stars counter (live) · Discord member count · Latest blog post

**Footer:**
MIT License · Security Policy · Privacy Policy · GitHub · Discord

---

### 2.4 Documentation Site (VitePress)

Configure `docs/.vitepress/config.ts`:

```typescript
export default defineConfig({
  title: 'Spectra',
  description: 'Cryptographic asset discovery and quantum risk intelligence',
  themeConfig: {
    logo: '/spectra-logo.svg',
    nav: [
      { text: 'Guide', link: '/guide/' },
      { text: 'Reference', link: '/reference/cli' },
      { text: 'Compliance', link: '/compliance/cnsa-20' },
      { text: 'Integrations', link: '/integrations/github-actions' },
      { text: 'Blog', link: 'https://spectra.tools/blog' },
    ],
    editLink: {
      pattern: 'https://github.com/HarshalPatel1972/spectra/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },
    search: { provider: 'local' },
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 Harshal Patel',
    }
  },
  // Enable Algolia DocSearch once search volume warrants it
})
```

**Getting Started page** must be operable in under 3 minutes:

```
## Installation

brew install harshalpatel1972/tap/spectra    # macOS/Linux
# or
go install github.com/HarshalPatel1972/spectra/cmd/spectra@latest
# or
docker pull ghcr.io/harshalpatel1972/spectra

## First Scan

cd your-project
spectra scan .

## With Output Files

spectra scan . --output cbom,html --out-dir ./spectra-out
open ./spectra-out/spectra-report.html
```

---

## 3. Tier 1b — GitHub Repository Hardening

The GitHub repo IS the product page for most developers. Treat it as such.

### 3.1 Repository README Structure

```markdown
# Spectra
### Cryptographic Asset Discovery & Quantum Risk Intelligence

![CI](badge) ![Go 1.25](badge) ![License MIT](badge) [![QRS: 8](badge)](https://spectra.tools)
![CBOM: CycloneDX 1.7](badge) ![CNSA 2.0 Compliant](badge)

> "See every cipher. Own your migration."

[Documentation](https://docs.spectra.tools) · [Live Playground](https://spectra.tools/playground) · [Discord](link) · [Blog](https://spectra.tools/blog)

---

## What is Spectra?

[One paragraph — the gut punch. Use the "world your cryptography protects" framing.]

## Quick Start

[Three code blocks: install, scan, view results. NOTHING ELSE until features.]

## Features

[Feature table with checkmarks — not bullet points]

## Output Formats

[Terminal screenshot + CBOM snippet + HTML report screenshot]

## Integration

[GitHub Actions YAML, VS Code screenshot, Docker one-liner]

## Who This Is For

[Three personas: Developer, Security Engineer, CISO]

## Architecture

[Mermaid diagram from Phase 1 build prompt]

## Contributing

[Link to CONTRIBUTING.md]

## License

MIT
```

### 3.2 Required Repository Files

Create all of these as distinct commits:

```
docs/CONTRIBUTING.md    — How to add patterns, file bugs, submit PRs
docs/SECURITY.md        — Responsible disclosure policy
docs/PRIVACY.md         — No telemetry statement, data handling
docs/THREAT_MODEL.md    — What threats Spectra is designed to address
.github/ISSUE_TEMPLATE/
  bug_report.md
  feature_request.md
  false_positive.md     — Report when Spectra misidentifies a safe algorithm
  false_negative.md     — Report when Spectra misses a vulnerable algorithm
.github/PULL_REQUEST_TEMPLATE.md
.github/CODEOWNERS
LICENSE                 — MIT
```

### 3.3 GitHub Community Infrastructure

Enable in repository settings:
- **Discussions** with categories: Announcements, Q&A, Ideas, Show & Tell
- **Sponsorship** (GitHub Sponsors) — even at $0 initial, signals the project is sustainable
- **Topics/Tags**: `pqc`, `cryptography`, `security`, `golang`, `cbom`, `quantum`, `post-quantum`

---

## 4. Tier 3 — The Interactive Web Properties

These are the highest-leverage distribution surfaces. Build these before any conference talk or blog post.

### 4.1 The Quantum Threat Visualization (`/what-happens`)

This is the most important page in the entire distribution strategy.
It does one thing: make the abstract threat VISCERAL.

**Technical build: Next.js page with Framer Motion animations.**

```
Route: /what-happens
Title: "What Happens If Quantum Computers Arrive Tomorrow?"
```

**Stage 1: The Connected World** (0:00 — 0:05)

```jsx
// Animated SVG network graph
// Nodes: Stripe, Bank, VPN, GitHub, AWS, API, Certificate Authority
// Edges: glowing green lines labeled "RSA-2048" and "ECDSA"
// Everything looks healthy and connected
```

Each node is an SVG icon with a pulsing green glow.
Lines connecting them are animated CSS strokes.
Label on each line: `RSA-2048` in tiny monospace text.

**Stage 2: Q-Day** (0:05 — 0:15, triggered by scroll or "Simulate" button)

A quantum computer icon fades in on the left.
Then, animated sequence:

```
A counter starts: "Breaking RSA-2048 keys..."
  1 key... (0.1s pause)
  100 keys... (0.1s pause)
  10,000 keys... (0.1s pause)
  ALL keys: BROKEN

Simultaneously:
  Each green line FLICKERS then turns RED
  Each service node dims and shows ❌
  A red banner appears over each: "COMPROMISED"

Final state: all nodes dark red, all lines broken
Text overlay: "Every RSA and ECC key in use today becomes worthless."
```

**Stage 3: "This Is Your Code Right Now"** (scroll-triggered)

A realistic terminal panel fades in showing Spectra output:

```
$ spectra scan ./payment-service
Scanning 847 files across 12 packages...

CRITICAL  RSA-2048     src/payments/encrypt.go:47     QRS: 90
CRITICAL  ECDSA/P-256  tls/certs/api.pem              QRS: 85
HIGH      SHA-1        src/legacy/hash_util.go:12     QRS: 70

Aggregate QRS: 88/100 — CRITICAL
47 compliance gaps with CNSA 2.0
```

Below the terminal:
```
"Spectra found this in a realistic payment service codebase.
 Does yours look different?"

[ Scan Your Code Now → ] (links to /playground)
```

**Stage 4: The Fix** (scroll-triggered)

The same network graph, but now showing:

```
spectra simulate --from RSA --to ML-KEM
→ Migration plan generated
→ 3 waves, 8 weeks, 2 engineers
```

Lines turn from red back to BLUE (PQC-safe color).
Service nodes glow blue instead of green.
New labels: `ML-KEM-1024` and `ML-DSA-87`.

```
"RSA → ML-KEM. ECDSA → ML-DSA.
 The path exists. Spectra shows you the order."

[ Start Your Migration → ] (links to /playground or install)
```

**Implementation note:** The animation must work without JavaScript disabled (use
CSS animations as fallback). The page must be shareable — OpenGraph image must be
the Stage 2 "COMPROMISED" screenshot with the Spectra logo.

---

### 4.2 Live Playground (`/playground`)

The playground removes the single biggest barrier: installation.
A developer who gets a result in the browser will install the CLI within 24 hours.

**Phase A: Server-side playground (build first)**

Architecture:
- Frontend: Next.js page with CodeMirror 6 (code editor) + output panel
- API: `/api/scan` Vercel serverless function that runs Spectra in a subprocess

```typescript
// app/api/scan/route.ts
// Accepts: { code: string, language: string, filename: string }
// Returns: { findings: Finding[], aggregateQRS: number, html: string }

export async function POST(req: Request) {
  const { code, language, filename } = await req.json()

  // Write to temp file, run spectra scan, clean up
  const tmpDir = await fs.mkdtemp('/tmp/spectra-playground-')
  const tmpFile = path.join(tmpDir, filename || `code.${langToExt(language)}`)
  await fs.writeFile(tmpFile, code)

  const result = await execSpectra(['scan', tmpDir, '--output', 'json', '--quiet'])
  await fs.rm(tmpDir, { recursive: true })

  return Response.json(JSON.parse(result))
}
```

**Rate limiting**: 20 scans per IP per hour. Return `429` with a friendly message:
"You've hit the playground rate limit. Install Spectra locally for unlimited scans."

**Privacy notice** (prominent, not buried):
```
"Your code is sent to our server, scanned, and immediately deleted.
 We do not log, store, or analyze playground submissions.
 For maximum privacy, use our WASM playground (no server required)."
```

**Phase B: WASM playground (build after Phase A ships)**

Compile Spectra's scanner core (excluding filesystem ops) to WASM:

```bash
# New build target in Makefile
wasm:
  GOOS=js GOARCH=wasm go build \
    -o landing/public/spectra.wasm \
    ./cmd/spectra-wasm

# cmd/spectra-wasm/main.go exposes:
# - scanCode(code string, lang string) string  → JSON findings
# - computeQRS(algo string, keySize int) int
```

The WASM build exposes two JavaScript-callable functions via `syscall/js`.
The browser runs the actual Go pattern matching logic locally.
No code leaves the browser.

Privacy label: `"🔒 WASM Mode: Code never leaves your browser."`

**Playground UI layout:**

```
┌─────────────────────────────────────────────────────┐
│  Spectra Playground              [Server] [🔒 WASM]  │
├──────────────────────────┬──────────────────────────┤
│ Language: [Go ▼]         │ Results                  │
│                          │                          │
│ [CodeMirror editor]      │ [Findings table]         │
│                          │ Aggregate QRS: 83/100    │
│ ...paste code here...    │                          │
│                          │ [Copy as CBOM]           │
│                          │ [View Migration Plan]    │
└──────────────────────────┴──────────────────────────┘
│ [ Scan ] [ Load Example: Go JWT ▼ ]                 │
└─────────────────────────────────────────────────────┘
```

**Pre-loaded examples** (show immediately, no interaction required):
- "Go JWT with RSA" → 3 findings (CRITICAL RSA, MEDIUM SHA)
- "Python cryptography" → 2 findings (HIGH ECDSA, SAFE AES-256)
- "Java PKCS" → 4 findings (CRITICAL RSA, HIGH DES)
- "Clean code" → 0 findings (shows QRS: 0/100 — SAFE)

---

### 4.3 Interactive QRS Calculator (`/risk-calculator`)

A visual calculator. No code required. Pure form-based.

```
"What is your Quantum Risk Score?"

Step 1: Which algorithms does your codebase use?
  [x] RSA-2048     [x] RSA-4096     [ ] ML-KEM
  [x] ECDSA/P-256  [ ] ML-DSA       [ ] SLH-DSA
  [x] SHA-1        [ ] SHA-256      [ ] SHA-3
  [ ] AES-128      [x] AES-256      [ ] DES

Step 2: How many usages (roughly)?
  RSA-2048: [slider 1-100]
  ECDSA: [slider 1-100]

Step 3: Where is it used?
  [x] Public-facing API    [ ] Internal only
  [x] Customer data        [ ] Test infrastructure

→ Real-time QRS display updates as user interacts

[ Get Your Full Migration Plan ]
   → links to /playground or install guide
```

**The calculator is shareable**: `?algorithms=RSA,ECDSA,SHA1&usages=14,7,3` URL params.
When a developer shares this link, others see the same assessment.

---

### 4.4 PQC Readiness Assessment (`/pqc-readiness`)

Ten questions. Five minutes. A score and an action plan.

```
Q1: Have you performed a cryptographic inventory in the last 12 months?
    (No / Partial / Yes)

Q2: Which algorithms do you actively use?
    (multi-select: RSA / ECC / SHA-1 / AES-128 / MD5 / PQC algorithms)

Q3: Can you change your cryptographic algorithms without code changes?
    (Never / With code changes / With config changes / Yes, runtime-configurable)

Q4: Do you have a Cryptographic Bill of Materials (CBOM)?
    (No / In progress / Yes)

Q5: Have you assessed your CNSA 2.0 compliance status?
    (Never heard of it / Aware but not assessed / Assessed / Compliant)

Q6: Do your third-party dependencies use classical cryptography?
    (Unknown / Yes definitely / Mostly / No)

Q7: Are your TLS certificates using ECDSA or RSA?
    (Unknown / RSA / ECDSA / ML-DSA / Mix)

Q8: Does your team have PQC expertise?
    (None / Basic awareness / Some expertise / Deep expertise)

Q9: Do you have CI/CD gates for cryptographic hygiene?
    (No / Planned / Partial / Yes)

Q10: When do you plan to start PQC migration?
     (No plans / 2027+ / 2025-2026 / Already started)

→ Output:
   Your PQC Readiness Score: 34/100 — DEVELOPING

   Top 3 gaps:
   1. No CBOM → "Run: spectra scan . --output cbom"
   2. No CI gate → "Add: spectra scan . --fail-on critical"
   3. No PQC algorithms detected → "See: migration roadmap"

   [ Download Full Assessment Report ] [ Start with Spectra → ]
```

---

### 4.5 CBOM Viewer (`/cbom-viewer`)

Drag-and-drop interface. Accepts CycloneDX 1.7 JSON or XML.

```
"Drop a CBOM file to visualize it"

[ Drag CBOM here or click to upload ]
→ or paste JSON:  [ textarea ]

Output:
  - Summary: X algorithms found, aggregate QRS Y
  - Table: all cryptographic components, sortable by QRS, algorithm, family
  - Risk pie chart (Chart.js, embedded)
  - Export as: CSV | PDF | Executive Summary

"No server required. Your CBOM is processed entirely in your browser."
```

---

## 5. Tier 4 — Developer Integrations (Distribution-as-Features)

Each integration is a **permanent advertising channel**. A VS Code extension that opens
silently when a developer opens a file is worth 1,000 blog posts.

### 5.1 VS Code Extension

**Package name**: `spectra-security`
**Display name**: `Spectra — Quantum Crypto Scanner`
**Publisher**: `HarshalPatel1972`

The extension activates on any supported file type and shows:
- A status bar item showing file QRS
- Inline diagnostic squiggles under crypto function calls (like a linter)
- Hover tooltip with algorithm info, QRS, and migration suggestion
- Sidebar panel with findings for the open workspace

**Critical activation trigger**: the extension shows up in VS Code's extension search
when users search: `"crypto"`, `"security"`, `"pqc"`, `"quantum"`.

**Extension manifest (`package.json`) — key fields:**

```json
{
  "name": "spectra-security",
  "displayName": "Spectra — Quantum Crypto Scanner",
  "description": "Find RSA, ECC, SHA-1 and other quantum-vulnerable cryptography in your codebase. Get QRS scores and migration plans inline.",
  "categories": ["Linters", "Security", "Other"],
  "keywords": ["cryptography", "security", "pqc", "quantum", "rsa", "ecdsa", "cbom"],
  "activationEvents": ["onLanguage:go", "onLanguage:python", "onLanguage:java",
                        "onLanguage:javascript", "onLanguage:typescript",
                        "onLanguage:rust", "onLanguage:c", "onLanguage:cpp"],
  "contributes": {
    "configuration": {
      "title": "Spectra",
      "properties": {
        "spectra.enabled": { "type": "boolean", "default": true },
        "spectra.minQRSToFlag": { "type": "number", "default": 40 },
        "spectra.binaryPath": { "type": "string", "default": "" }
      }
    },
    "commands": [
      { "command": "spectra.scanWorkspace", "title": "Spectra: Scan Workspace" },
      { "command": "spectra.openReport", "title": "Spectra: Open HTML Report" }
    ]
  }
}
```

**Publish to VS Code Marketplace** via `vsce publish`. This gives passive discovery from
the ~50 million VS Code installs.

---

### 5.2 GitHub App — "Spectra CI"

**Name**: Spectra CI
**Description**: Automatically scan pull requests for quantum-vulnerable cryptography

The app runs `spectra scan --scanners code --output json` on every PR's changed files
and posts a PR comment with results.

**Architecture:**
```
spectra-app/
├── index.ts              # Express server, handles GitHub webhook events
├── handlers/
│   ├── pull_request.ts   # Run scan on PR diff files
│   └── installation.ts   # Handle new installations
├── scanner/
│   └── scan.ts           # Spawn spectra CLI, parse JSON output
└── templates/
    └── pr-comment.md.ts  # PR comment template
```

**PR Comment Template:**

```markdown
## Spectra Quantum Risk Scan

| Algorithm | File | Line | QRS | Action |
|---|---|---|---|---|
| RSA-2048 | `auth/jwt.go` | 47 | 🔴 90 | Replace with ML-KEM |
| SHA-1 | `legacy/hash.go` | 12 | 🟡 70 | Replace with SHA-256 |

**Aggregate QRS: 83/100 — CRITICAL**
2 CRITICAL · 1 HIGH · 0 MEDIUM

[📊 View Full Report](https://spectra.tools/cbom-viewer) · [📖 Migration Guide](https://docs.spectra.tools)

_Powered by [Spectra](https://spectra.tools) v{version}_
```

**The footer "Powered by Spectra"** on every PR comment in every repository that installs
the app is persistent brand exposure. This is how Netlify and Vercel built brand recognition.

**Deploy on**: Railway free tier (250 hours/month free — sufficient for MVP).

---

### 5.3 Badge Economy

The most underrated viral mechanic for developer tools.

Add `spectra badge` command to the CLI:

```bash
spectra scan . --persist --badge
# Outputs:
# QRS: 23 — LOW
# Add this to your README:
#
# [![Spectra QRS: 23](https://badge.spectra.tools/qrs?score=23&label=Spectra)](https://spectra.tools)
```

Badge server: a Vercel Edge Function at `badge.spectra.tools/qrs` that returns an SVG badge.
The badge uses shields.io-style color coding:
- 0–20: green (`#4ADE80`)
- 21–40: yellowgreen (`#A3E635`)
- 41–60: yellow (`#FACC15`)
- 61–80: orange (`#F97316`)
- 81–100: red (`#EF4444`)

```
// Vercel Edge Function
// /api/badge/qrs/route.ts
// Query: ?score=23&label=Spectra
// Returns: SVG badge
```

**Why this works:** Every developer who adds the badge to their README is advertising
Spectra to every person who views that README. GitHub shows README badges in search
results. A popular repo with a Spectra badge gets seen by thousands of developers daily.

---

### 5.4 Docker Image

```dockerfile
# Dockerfile
FROM scratch
COPY --from=build /go/bin/spectra /spectra
ENTRYPOINT ["/spectra"]
```

Published to `ghcr.io/harshalpatel1972/spectra:latest` via GitHub Actions on every tag.

```bash
# Usage:
docker run --rm -v $(pwd):/workspace ghcr.io/harshalpatel1972/spectra scan /workspace

# Or with output files:
docker run --rm \
  -v $(pwd):/workspace \
  -v $(pwd)/spectra-out:/out \
  ghcr.io/harshalpatel1972/spectra scan /workspace --out-dir /out --output html,cbom
```

---

### 5.5 Pre-commit Hook

```yaml
# .pre-commit-hooks.yaml (in the spectra repo)
- id: spectra-scan
  name: Spectra Quantum Crypto Scanner
  description: Fail commit if new CRITICAL quantum-vulnerable cryptography is introduced
  entry: spectra scan
  args: [--fail-on=critical, --quiet, --no-progress]
  language: system
  pass_filenames: false
  always_run: false
```

Installation:
```yaml
# .pre-commit-config.yaml (in the user's project)
repos:
  - repo: https://github.com/HarshalPatel1972/spectra
    rev: v0.1.0
    hooks:
      - id: spectra-scan
```

---

### 5.6 GitHub Actions — Official Action

Create `HarshalPatel1972/spectra-action` repository:

```yaml
# action.yml
name: 'Spectra Quantum Crypto Scan'
description: 'Scan for quantum-vulnerable cryptography in your repository'
inputs:
  fail-on:
    description: 'Minimum severity to fail: any | high | critical'
    default: 'critical'
  output-formats:
    description: 'Output formats: terminal,json,cbom,html'
    default: 'terminal,cbom'
  path:
    description: 'Path to scan'
    default: '.'
outputs:
  aggregate-qrs:
    description: 'Aggregate Quantum Risk Score (0-100)'
  cbom-path:
    description: 'Path to generated CBOM file'
runs:
  using: 'docker'
  image: 'docker://ghcr.io/harshalpatel1972/spectra:latest'
  args:
    - scan
    - ${{ inputs.path }}
    - --fail-on=${{ inputs.fail-on }}
    - --output=${{ inputs.output-formats }}
```

Usage example (in the README and marketplace listing):
```yaml
- name: Scan for quantum-vulnerable crypto
  uses: HarshalPatel1972/spectra-action@v1
  with:
    fail-on: critical
```

---

## 6. Tier 2 — Trust-Building Content

### 6.1 Architecture Whitepaper

File: `docs/whitepaper/spectra-architecture.md` — published as PDF and web page.

```
Title: "Spectra: A Practical Architecture for Cryptographic Asset Discovery"

Sections:
1. The Problem Statement
   — Scope of quantum threat to deployed cryptography
   — Why discovery must precede migration
   — The CBOM gap in current tooling

2. Architecture Overview
   — Four scanner types and their detection methodology
   — Quantum Risk Score derivation and weighting rationale
   — Cryptographic Agility Index methodology
   — Compliance framework embedding

3. Comparison with Alternative Approaches
   — grep/regex (false negatives, no context, no scoring)
   — SAST tools (not crypto-specific, miss certs/deps/configs)
   — SCA tools (dependency-only, miss code-level usage)
   — Building internally (2-6 months, ongoing maintenance)

4. Algorithm Detection Accuracy
   — True positive rate by language and algorithm
   — Known limitations and false negative categories
   — How to contribute new patterns

5. Standards Alignment
   — CycloneDX 1.7 CBOM specification compliance
   — CNSA 2.0 requirement coverage
   — NIST SP 800-131A coverage

6. Security Properties of Spectra Itself
   — No network calls during scanning
   — No telemetry
   — Deterministic output
   — Input sanitization for the scan API

7. Roadmap and Contribution Opportunities

References: [20+ citations to NIST, NSA, IETF, academic papers]
```

---

### 6.2 The "Why Spectra Exists" Article

First blog post. Published on Day 1. Cross-posted to Dev.to, Medium, LinkedIn.

**Title**: "The Cryptographic Debt Nobody Is Talking About"

**Structure (not prose — actual section headings and content guidance):**

```
Opening: A specific day — the day a researcher demonstrated SHA1 collision.
         Not "in 2017, researchers showed..." but a scene: a room, a demo,
         SHA1 broken in front of a live audience.

Section 1: The quiet assumption
  Most developers assume their cryptographic choices are safe by default.
  They chose RSA-2048 because the documentation said so in 2019.
  Nobody told them to revisit it.

Section 2: The problem with "it's not broken yet"
  RSA-2048 has not been broken. It will be broken.
  The gap between "not broken yet" and "broken" for RSA
  is measured in years, not decades.
  You cannot migrate an enterprise codebase in months.
  The migration must begin before the threat materializes.

Section 3: The discovery gap
  We have CVE scanners for vulnerabilities.
  We have dependency scanners for known-bad libraries.
  We have SAST tools for code bugs.
  We have no tool specifically designed to answer:
  "What cryptography am I using, where, why, and what needs to change?"

Section 4: What Spectra does
  [Three-command demo. Show output. Let the tool speak.]

Section 5: The next step
  [Link to getting started. No sales pitch. Just the next action.]
```

---

### 6.3 Comparison Pages

Three pages, each with a specific, honest comparison.

**`/compare/vs-grep`** — "Spectra vs `grep -r RSA`"

| Capability | grep | Spectra |
|---|---|---|
| Find algorithm names | ✓ | ✓ |
| Context-aware detection | ✗ (matches comments) | ✓ |
| Key size extraction | ✗ | ✓ |
| Certificate scanning | ✗ | ✓ |
| Dependency manifest scanning | ✗ | ✓ |
| Quantum Risk Score | ✗ | ✓ |
| CBOM generation | ✗ | ✓ (CycloneDX 1.7) |
| Compliance gap mapping | ✗ | ✓ |
| CI/CD exit codes | ✗ | ✓ |
| False positive filtering | ✗ | ✓ |
| Migration action plan | ✗ | ✓ |

**`/compare/vs-diy`** — "Spectra vs Building Your Own Scanner"

Be honest about this. A senior team can build a basic scanner.
The comparison is not "can you do this" but "should you":
- 2–6 months of engineering time
- Ongoing maintenance as new crypto patterns appear
- No CBOM format compliance
- No compliance framework embedding
- No community-maintained pattern library

**`/compare/vs-sca`** — "Spectra vs Generic SCA Tools (Snyk, Trivy)"

Explicitly: Spectra is NOT a replacement for Snyk or Trivy.
They solve CVEs. Spectra solves algorithmic quantum exposure.
A dependency with no CVEs can still use RSA-2048.
A file with no dependency issues can still call `sha1.New()`.
These tools are complementary.

---

### 6.4 Security Policy (`SECURITY.md`)

```markdown
# Security Policy

## Reporting Vulnerabilities

Report security issues to: security@spectra.tools
Do not file public GitHub issues for security vulnerabilities.
We will respond within 48 hours and coordinate disclosure.

## Spectra's Security Properties

**No telemetry.** Spectra never phones home. Zero network calls during scanning.
**No account required.** No login, no API key, no data collection.
**No data retention.** The playground API deletes submitted code immediately after scanning.
**Deterministic output.** The same input always produces the same output.

## What Spectra Does Not Do

Spectra does not prevent cryptographic vulnerabilities. It discovers them.
Spectra does not store your findings. That is your responsibility.
```

---

## 7. Tier 8 — The "What Happens..." Page (Expanded Spec)

This tier was listed last but is built FIRST. It is the highest-leverage single page.

The page at `/what-happens` has one conversion goal: make a developer think
*"I need to check my code right now."*

**SEO title**: "What Happens to RSA and ECDSA When Quantum Computers Arrive?"
**Meta description**: "An interactive visualization of what a cryptographically-relevant quantum computer means for today's infrastructure — and how to prepare."

**The Flow to the Bottom of the Page:**

```
Section A: The Setup (5 seconds to read)
  "Right now, your infrastructure is protected by mathematics.
   RSA-2048 is secure because factoring large integers is computationally hard.
   ECDSA is secure because the elliptic curve discrete logarithm problem is hard.
   A quantum computer running Shor's algorithm solves both of these in polynomial time."

Section B: The Visualization (the animation described in §4.1 above)

Section C: When Does This Happen?
  Two-column layout:

  LEFT: "The Optimistic View"
    "Large-scale quantum computers capable of breaking RSA-2048 are likely
     still 10–15 years away. NIST and major cryptographers estimate the earliest
     credible timeline as early 2030s."

  RIGHT: "Why This Matters NOW"
    Countdown-style display:
    "CNSA 2.0 requires new NSS systems to use PQC by:    2027 (2 years)"
    "CNSA 2.0 requires exclusive PQC use by:            2033 (8 years)"
    "Average enterprise crypto migration timeline:       3–5 years"
    Red highlight: "If you start in 2028, you miss the 2030 window."

    "The 'Harvest Now, Decrypt Later' attack is happening TODAY.
     Adversaries collecting your RSA-encrypted traffic now
     can decrypt it the moment a quantum computer exists."

Section D: What This Means For Specific Systems
  Card grid (one per system type):

  Card: "TLS/HTTPS"
    Currently: ECDHE key exchange (ECC-based) → QUANTUM VULNERABLE
    At risk: All data in transit
    Fix: ML-KEM for key exchange (already in Chrome, Firefox)
    Timeline: Already supported in TLS 1.3 via RFC 9180

  Card: "Code Signing Certificates"
    Currently: RSA-2048 or ECDSA → QUANTUM VULNERABLE
    At risk: Software supply chain integrity
    Fix: ML-DSA-87 signatures
    Timeline: CNSA 2.0 requires this for firmware by 2030

  Card: "JWT / API Authentication"
    Currently: RS256 (RSA) or ES256 (ECDSA) → QUANTUM VULNERABLE
    At risk: All API authentication tokens
    Fix: No standard yet for ML-DSA JWTs (IETF draft in progress)
    Timeline: Hybrid approach recommended now

  Card: "Database Encryption Keys"
    Currently: RSA key wrapping → QUANTUM VULNERABLE
    At risk: All data at rest
    Fix: ML-KEM for key encapsulation
    Timeline: Migrate before 2030

  Card: "VPN Key Exchange"
    Currently: DH or ECDH → QUANTUM VULNERABLE
    At risk: All network traffic
    Fix: ML-KEM-1024 in IKEv2
    Timeline: CNSA 2.0 requires exclusive use by 2030

Section E: "Where Does Your Code Stand?"
  Final CTA:
  "Spectra scans your codebase, certificates, and dependencies to find
   exactly which of the above you need to migrate. Free. Open source. No account."

  [ brew install spectra ] or [ Try in Browser ]
```

---

## 8. The Launch Playbook

One-time opportunity. Execute it right.

### 8.1 Pre-Launch Sequence (6 Weeks Before)

**Week -6:**
- Publish "The Cryptographic Debt Nobody Is Talking About" on LinkedIn + Dev.to
- No mention of Spectra yet — just the problem
- Collect comments and reactions

**Week -4:**
- Publish "We scanned 50 popular Go open-source projects"
  (Methodology: run Spectra on the top 50 Go repos by stars on GitHub)
  (Key finding: "47 of 50 contained RSA or ECDSA usage. 31 contained SHA-1.")
  (This article names the repos — they may share the article)

**Week -3:**
- Publish `/what-happens` page
- Share it without the tool: "Built a page to visualize the quantum threat. No tool to sell yet."
- Collect traffic and discussion

**Week -2:**
- GitHub repo public with README, CONTRIBUTING.md, first release tag
- Product Hunt "Coming Soon" page activated
- Discord server created, link in GitHub README

**Week -1:**
- Email Product Hunt followers the launch date
- Write the HN "Show HN" draft (review it)
- Share the `/playground` page to select developer communities for beta feedback

### 8.2 Launch Day Timeline

**06:00 UTC:** Post to Product Hunt
**08:00 UTC:** Post to Hacker News: "Show HN: Spectra – CLI to find quantum-vulnerable cryptography in your codebase"
**10:00 UTC:** Post to Reddit (r/netsec, r/golang, r/programming) — different framing each
**12:00 UTC:** LinkedIn announcement post (professional framing, CNSA 2.0 angle)
**14:00 UTC:** Dev.to article with full tutorial
**16:00 UTC:** X/Twitter thread (technical demo, short clips of terminal output)

**Hacker News title**: `Show HN: Spectra — scans codebases for quantum-vulnerable cryptography (RSA, ECC, SHA-1)`

**HN text body** (this matters enormously):
```
Hey HN,

I built Spectra after reading through NIST's PQC migration guides and realizing
most codebases have no idea what cryptography they're using. There's no systematic
way to answer "where is RSA in my infrastructure?"

Spectra scans source code (Go, Python, Java, JS, Rust, C/C++), X.509 certificates,
config files, and dependency manifests. It outputs:
- A Quantum Risk Score (0-100) per finding and aggregate
- A CycloneDX 1.7 CBOM
- A migration action plan sorted by risk × effort
- Compliance gaps against CNSA 2.0 and NIST SP 800-131A

brew install harshalpatel1972/tap/spectra
spectra scan .

Source: github.com/HarshalPatel1972/spectra

The most interesting technical challenge was the QRS formula — how to weight
algorithm vulnerability, key size, and usage frequency into a single actionable
number. Happy to discuss the methodology.
```

**Reddit posts (different for each subreddit):**

r/netsec: Lead with CNSA 2.0 deadline. Focus on compliance risk.
r/golang: Lead with the technical implementation. Show the Go WASM compilation.
r/programming: Lead with the problem. "How do you know what cryptography is in your codebase?"
r/devops: Lead with the CI/CD integration. `--fail-on=critical`

---

### 8.3 Post-Launch Content Calendar (Months 1-6)

**Month 1:**
- "Spectra Launch Week: What We Learned" (transparency post)
- First "Scan of the Month": scan a famous open-source project, publish findings
- Fix the top 3 issues from user feedback (rapid release)

**Month 2:**
- Architecture whitepaper published
- VS Code extension published to Marketplace
- First case study (from a real user who agreed to be featured)
- Blog: "Understanding the Cryptographic Agility Index"

**Month 3:**
- GitHub App (Spectra CI) released
- Blog: "Why CBOM Matters: A Practical Guide to CycloneDX 1.7"
- Second "Scan of the Month"
- Conference talk submitted to: DEF CON 33, Black Hat USA 2026

**Month 4:**
- WASM playground launched ("Your code, never leaves your browser")
- Blog: "CNSA 2.0 Explained: What Every Developer Needs to Know"
- Collaboration with a PQC researcher or security team (co-authored article)

**Month 5:**
- Docker Desktop extension submitted
- Blog: "Scanning the Linux Kernel's TLS Implementation with Spectra"
- Monthly newsletter launched (Buttondown or Substack)

**Month 6:**
- 6-month retrospective post ("X repos scanned, Y findings discovered")
- Explore partnership with CycloneDX working group or OWASP

---

## 9. The "Scan of the Month" Format

This content format is the highest-leverage recurring content type.
It requires no new features. It demonstrates value. It gets shared by the projects scanned.

**Template for each Scan of the Month:**

```
Title: "Quantum Exposure Report: [Project Name] v[Version]"

One-paragraph intro: Why this project matters and who uses it.

Scan command:
  git clone https://github.com/[org]/[project]
  spectra scan . --output json,html

Findings summary:
  - X total findings
  - Y CRITICAL, Z HIGH
  - Top algorithm: [name] (N occurrences)
  - Aggregate QRS: [score]

Top 5 findings (with file path, line, context):
  [table]

The interesting parts:
  [1-2 paragraphs of technical analysis — not just the numbers]
  [What's particularly surprising or well-handled]

Migration difficulty assessment:
  [Using CAI — how hard would it be for this project to migrate?]

Note to project maintainers:
  "This is a non-judgmental analysis. These algorithms were correct choices
   when they were written. This report is about what changes SHOULD happen,
   not about what was done WRONG."

Full CBOM: [GitHub Gist or file download]
```

**First three candidates:**
1. `hashicorp/vault` — high-profile, security-focused, will generate attention
2. `golang/go` standard library — meta: what crypto does the Go stdlib itself use?
3. `cert-manager` — TLS certificate management tool, maximum relevance

---

## 10. The Conference Talk Framework

**Title**: "Invisible Vulnerabilities: Finding Your Cryptographic Debt Before Quantum Arrives"

**Pitch to**: DEF CON 33, Black Hat USA 2026, CCC 2026, USENIX Security Symposium, PQCrypto 2026

**Abstract (200 words):**

```
Most organizations know post-quantum cryptography exists. Few know what
classical cryptography they currently use. This talk presents Spectra,
an open-source tool for cryptographic asset discovery, and shares findings
from scanning [N] popular open-source projects totaling [X] million lines of code.

We'll cover: (1) the detection methodology — how to reliably identify RSA, ECC,
SHA-1, and other quantum-vulnerable algorithms across polyglot codebases;
(2) the Quantum Risk Score model — how to weight algorithm vulnerability,
key size, and usage frequency into an actionable priority queue;
(3) real findings from real codebases — what the open-source ecosystem
actually uses, and where the highest-density quantum debt lives;
(4) the migration path — how the Cryptographic Agility Index predicts
migration cost before you begin.

Live demo: scan a repository on stage. Show the findings. Show the CBOM.
Show the action plan. Total time from `brew install` to first report: 90 seconds.

Takeaway: every attendee leaves knowing exactly what to run on Monday morning.
```

**Demo structure for live talk:**
1. Live scan of a real repo (pre-selected for interesting findings)
2. Show CBOM in the browser viewer
3. Run compliance check against CNSA 2.0 — show deadline countdown
4. Run `spectra simulate` — show migration plan
5. Show QRS go from 83 → projected 12 after migration

---

## 11. Tier 6 — Content Marketing

### 11.1 The Blog Strategy

**SEO-targeted post series (long-tail keywords with real search volume):**

| Title | Target keyword | Search intent |
|---|---|---|
| "How to find RSA usage in a Go codebase" | find RSA go | Informational |
| "What is a Cryptographic Bill of Materials (CBOM)?" | CBOM explained | Informational |
| "CNSA 2.0 compliance checklist for developers" | CNSA 2.0 checklist | Informational |
| "SHA-1 is deprecated — here's how to find it in your code" | SHA-1 deprecated | Informational |
| "What is post-quantum cryptography and why should I care?" | post-quantum cryptography | Informational |
| "How to migrate from RSA to ML-KEM" | RSA to ML-KEM | How-to |
| "CycloneDX CBOM format explained" | CycloneDX CBOM | Informational |
| "Quantum Risk Score: how to quantify PQC exposure" | quantum risk score | Informational |

Every article ends with a section titled "Check your own code" and a code block:
`brew install spectra && spectra scan .`

### 11.2 Monthly PQC Newsletter

Platform: **Buttondown** (privacy-first, developer-friendly, free for first 100 subscribers)

Newsletter name: **Quantum Drift** — Monthly PQC and Spectra updates

Template:
```
This month in PQC:
  → [1 CNSA/NIST news item]
  → [1 industry migration story]
  → [1 new library/algorithm development]

Spectra updates:
  → What shipped in vX.Y.Z
  → What's coming

Scan of the Month:
  → [project name] scanned. [N] findings. [1 interesting insight.]

One thing to do this month:
  → [single actionable recommendation]
```

---

## 12. Distribution Atomic Commits

Each item below is one commit / deployable unit. Complete in order.

```
feat(site): initialize Next.js 15 landing site with Tailwind
feat(site): build landing page hero with Framer Motion terminal animation
feat(site): build problem section with connected-world → compromised visualization
feat(site): build "how it works" three-step section
feat(site): build integrations row and output format tabs
feat(site/what-happens): build quantum threat visualization Stage 1 (healthy network)
feat(site/what-happens): build quantum threat visualization Stage 2 (Q-Day animation)
feat(site/what-happens): build Stage 3 terminal output and Stage 4 recovery animation
feat(site/what-happens): add specific system cards (TLS, VPN, JWT, certs)
feat(site/playground): build CodeMirror editor with language selector
feat(site/playground): implement /api/scan serverless function with rate limiting
feat(site/playground): add pre-loaded examples for all five supported languages
feat(site/playground): add server vs WASM mode toggle (WASM deferred to Phase B)
feat(site/calculator): build interactive QRS calculator with URL-shareable params
feat(site/assessment): build 10-question PQC readiness assessment
feat(site/cbom-viewer): build drag-and-drop CBOM visualization widget
feat(site/compare): build vs-grep comparison page
feat(site/compare): build vs-diy and vs-sca comparison pages
feat(site/blog): build MDX blog with first three articles
feat(site): build roadmap page with public phase tracking
feat(site): build changelog page (auto-generated from CHANGELOG.md)
feat(docs): initialize VitePress documentation site
feat(docs): write getting started, installation, quick start guides
feat(docs): write CLI reference (all commands and flags)
feat(docs): write algorithm database reference
feat(docs): write QRS scoring methodology documentation
feat(docs): write CNSA 2.0, NIST 800-131A compliance documentation
feat(docs): write GitHub Actions, VS Code, Docker integration guides
feat(badge): implement /api/badge/qrs SVG badge endpoint
feat(badge): add spectra badge CLI command
feat(vscode): initialize VS Code extension with manifest and activation
feat(vscode): implement diagnostic provider (inline crypto squiggles)
feat(vscode): implement hover provider (algorithm info + QRS on hover)
feat(vscode): implement status bar item with workspace QRS
feat(vscode): publish extension to VS Code Marketplace
feat(github-app): initialize Spectra CI GitHub App with webhook server
feat(github-app): implement PR scan trigger and result parser
feat(github-app): implement PR comment template with findings table
feat(github-app): deploy GitHub App server to Railway
feat(docker): create minimal scratch-based Dockerfile
feat(docker): publish to ghcr.io via GitHub Actions on tag
feat(actions): create harshalpatel1972/spectra-action repository
feat(actions): implement action.yml with Docker runner
feat(actions): publish action to GitHub Marketplace
feat(precommit): add .pre-commit-hooks.yaml to spectra repo
feat(brew): create homebrew tap formula (harshalpatel1972/homebrew-tap)
feat(site/wasm): compile spectra scanner core to WASM
feat(site/wasm): implement cmd/spectra-wasm entry point with syscall/js bindings
feat(site/playground): add WASM mode with privacy-first badge
docs(whitepaper): publish architecture whitepaper
docs(blog): publish "The Cryptographic Debt Nobody Is Talking About"
docs(blog): publish first Scan of the Month
content(launch): submit to Product Hunt
content(launch): prepare and post Show HN
content(newsletter): launch Quantum Drift newsletter on Buttondown
```

---

## 13. Metrics: How to Know It's Working

Define success metrics before launch. Track weekly.

| Metric | Week 1 Target | Month 1 | Month 3 | Month 6 |
|---|---|---|---|---|
| GitHub stars | 50 | 250 | 500 | 1,000 |
| Playground scans | 200 | 1,000 | 5,000 | 15,000 |
| CLI installs (Homebrew) | 100 | 500 | 2,000 | 5,000 |
| VS Code installs | 50 | 300 | 1,500 | 4,000 |
| GitHub App installs | 5 | 30 | 150 | 500 |
| /what-happens pageviews | 500 | 3,000 | 10,000 | 30,000 |
| Newsletter subscribers | 50 | 200 | 800 | 2,000 |
| Documentation pageviews | 200 | 1,000 | 5,000 | 15,000 |
| Inbound links | 5 | 25 | 100 | 300 |
| Scan of Month shares | — | 50 | 200 | 600 |

**Leading indicator**: The `/what-happens` page time-on-page. If users read it for >90 seconds,
the problem is landing. If they bounce in <15 seconds, the opening animation is not working.

**Single most important metric**: Playground → CLI install conversion rate.
A developer who tries the playground and then installs the CLI is the ideal user journey.
Track this with a UTM parameter on the install CTA.

---

## 14. The Enterprise Funnel (Tier 5 Preview)

Phase 3 of Spectra is not a public roadmap item yet. But the free tool must be designed
so it naturally leads here. Every design decision should serve the funnel:

```
FREE CLI (public)
  ↓ pain point: running manually, no team sharing, no historical trends
GITHUB APP (freemium)
  ↓ pain point: single repo, no dashboard, no team management
SELF-HOSTED DASHBOARD (team tier)
  ↓ pain point: need SSO, audit logs, compliance exports, branded reports
ENTERPRISE (contract)
  — Multi-tenant, RBAC, SIEM integration, custom compliance frameworks
  — Private support SLA
  — Co-branded executive reports for CISO presentations
```

The free tier must be genuinely excellent. The value of the paid tier must be
about **organizational scale** (multi-user, persistent, auditable, reportable),
not about locking features that individuals need.

---

## 15. Summary: Execution Priority Order

Build these in this exact order. Do not skip tiers.

```
Phase A (Week 1-2) — The Foundation:
  1. GitHub repo hardened (README, SECURITY, PRIVACY, CONTRIBUTING)
  2. /what-happens page live
  3. Basic landing page live (hero + problem section + install)
  4. Basic documentation live (install + quick start)
  5. Discord server created, linked

Phase B (Week 3-4) — Activation:
  6. Playground (server-side API)
  7. QRS calculator
  8. Blog with first article ("The Cryptographic Debt")
  9. Homebrew tap + Docker image

Phase C (Week 5-6) — Launch:
  10. Product Hunt prep
  11. VS Code extension published
  12. HN Show HN + Product Hunt launch
  13. "Scan of the Month" #1 published

Phase D (Month 2-3) — Growth:
  14. GitHub App (Spectra CI)
  15. GitHub Actions official action
  16. Architecture whitepaper
  17. Badge economy

Phase E (Month 4-6) — Enterprise Foundation:
  18. WASM playground ("privacy mode")
  19. CBOM viewer
  20. PQC readiness assessment
  21. Conference talk submissions
  22. Newsletter launched
```

---

*Build Prompt #3 authored: May 2026 — the distribution architecture for Spectra,
a Cryptographic Intelligence Platform.*
*Priority order: make them feel the problem → remove friction to first result → give them a reason to share.*
