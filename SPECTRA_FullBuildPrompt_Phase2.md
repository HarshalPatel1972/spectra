# SPECTRA — Full Build Prompt #2
### *Cryptographic Intelligence Platform: Phase 2 Strategic Directive*

---

> **Stop thinking of Spectra as a scanner.**
> A scanner tells you what exists. An intelligence platform tells you what it means,
> what to do about it, in what order, at what cost, by what deadline, for whom, and why.
>
> The question Phase 2 answers is not "what cryptography is in use?"
> The question Phase 2 answers is: **"What is our quantum exposure, and how do we survive it?"**

---

## 0. The Reframe: From Scanner to Operating System

Phase 1 delivered: detection, scoring, and reporting. These are now **table stakes**,
not differentiators. Every serious security team can build a grepping scanner with an LLM.

**What they cannot quickly build:**

1. A persistent, queryable **Cryptographic Knowledge Graph** of their entire infrastructure
2. A **Migration Simulation Engine** that predicts the blast radius of any algorithm change
3. A **Temporal Intelligence Layer** that maps who introduced each vulnerability, when, and why
4. A **Compliance Gap Engine** with embedded, authoritative regulatory timelines
5. A **Cryptographic Agility Index** that measures how painful your next migration will be
6. An **Evidence-Based Finding System** where every score has a verifiable chain of reasoning
7. A **Posture Score** that translates technical exposure into a board-level risk metric

None of these can be rebuilt quickly. Each requires weeks of domain expertise to encode
correctly. Together they create a **data and reasoning moat** that compounds over time.

**The competitive advantage is not the features. It is the knowledge graph.**

Once Spectra has mapped an organization's cryptographic graph, that graph takes months
to recreate. The longer they use Spectra, the richer the history, the higher the switching cost.

---

## 1. Phase 2 Capabilities — The Full Expansion

Harshal listed five items. Phase 2 expands and reframes them into thirteen capabilities,
each chosen because it answers a question a company would pay to answer:

| # | What Harshal Listed | What Phase 2 Delivers | Business Question Answered |
|---|---|---|---|
| 1 | Cross-source correlation | **Cryptographic Relationship Graph** | What breaks if I change X? |
| 2 | TLS endpoint scanner | **Live Infrastructure Interrogation** | What does our perimeter expose? |
| 3 | Git blame integration | **Temporal Intelligence Engine** | Who owns this debt, and since when? |
| 4 | `spectra diff` | **Migration Progress Intelligence** | Are we getting better or worse? |
| 5 | Container image scanning | **OCI Intelligence Layer** | Is our runtime environment safe? |
| 6 | *(new)* | **SQLite Persistence Layer** | How has our posture changed over time? |
| 7 | *(new)* | **Cryptographic Agility Index (CAI)** | How hard will the next migration be? |
| 8 | *(new)* | **Compliance Mapping Engine** | Are we meeting our regulatory obligations? |
| 9 | *(new)* | **Migration Simulation Engine** | What is the plan and in what order? |
| 10 | *(new)* | **Evidence-Based Finding Enrichment** | Why is this a problem? Prove it. |
| 11 | *(new)* | **Cryptographic Posture Score (CPS)** | How do I explain this to the board? |
| 12 | *(new)* | **Executive Audit Report Generator** | What does the CISO present to leadership? |
| 13 | *(new)* | **Team Attribution Engine** | Who owns the most quantum debt? |

---

## 2. Architectural Foundation — The Persistence Layer

**Every Phase 2 capability requires state.** Phase 1 was stateless: scan → output → done.
Phase 2 needs to remember, compare, trend, and reason across time.

### 2.1 Storage Engine

Use `modernc.org/sqlite` (SQLite 3.53.1 embedded as pure Go — no CGO, zero external
dependencies, full cross-compilation support). This preserves the single-binary guarantee.

State file location: `~/.spectra/state.db` (configurable via `--state-db` flag or
`SPECTRA_STATE_DB` environment variable).

> **Why not PostgreSQL?** Spectra must work offline, on developer laptops, in air-gapped
> environments, and in CI with no infrastructure provisioned. SQLite is the right choice.
> If an enterprise wants a shared server database, that is a future `spectra server` mode.

### 2.2 Schema — Complete Definition

```sql
-- ─────────────────────────────────────────────────────────
-- SCANS
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS scans (
    id           TEXT PRIMARY KEY,         -- UUID v4
    scan_root    TEXT NOT NULL,
    started_at   TEXT NOT NULL,            -- RFC3339
    completed_at TEXT NOT NULL,
    go_version   TEXT,
    spectra_version TEXT,
    total_files  INTEGER,
    scanned_files INTEGER,
    aggregate_qrs INTEGER,
    cps          INTEGER,                  -- Cryptographic Posture Score
    cai          INTEGER                   -- Cryptographic Agility Index
);

-- ─────────────────────────────────────────────────────────
-- FINDINGS
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS findings (
    id               TEXT PRIMARY KEY,
    scan_id          TEXT NOT NULL REFERENCES scans(id),
    algorithm        TEXT NOT NULL,
    source           TEXT NOT NULL,        -- CODE|CERT|DEPS|CONFIG|ENDPOINT|CONTAINER
    file_path        TEXT,
    line_number      INTEGER,
    line_content     TEXT,
    language         TEXT,
    key_size         INTEGER,
    context          TEXT,
    qrs              INTEGER NOT NULL,
    risk_band        TEXT NOT NULL,
    effort           TEXT NOT NULL,
    effort_rationale TEXT,
    container_layer  TEXT,                 -- for OCI findings
    endpoint_host    TEXT,                 -- for endpoint findings
    introduced_commit TEXT,               -- from git blame
    introduced_author TEXT,
    introduced_date  TEXT,
    introduced_message TEXT,
    team_name        TEXT,                 -- from attribution
    finding_hash     TEXT NOT NULL        -- SHA256(algo+path+line) for deduplication
);
CREATE INDEX IF NOT EXISTS idx_findings_scan_id   ON findings(scan_id);
CREATE INDEX IF NOT EXISTS idx_findings_algorithm ON findings(algorithm);
CREATE INDEX IF NOT EXISTS idx_findings_risk_band ON findings(risk_band);
CREATE INDEX IF NOT EXISTS idx_findings_hash      ON findings(finding_hash);

-- ─────────────────────────────────────────────────────────
-- GRAPH NODES
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS graph_nodes (
    id          TEXT PRIMARY KEY,         -- UUID v4
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    node_type   TEXT NOT NULL,            -- ALGORITHM|FILE|DEPENDENCY|CERT|ENDPOINT|TEAM
    label       TEXT NOT NULL,
    properties  TEXT NOT NULL             -- JSON blob of node properties
);

-- ─────────────────────────────────────────────────────────
-- GRAPH EDGES
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS graph_edges (
    id          TEXT PRIMARY KEY,
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    from_node   TEXT NOT NULL REFERENCES graph_nodes(id),
    to_node     TEXT NOT NULL REFERENCES graph_nodes(id),
    edge_type   TEXT NOT NULL,            -- USES|DEPENDS_ON|SIGNS|CONFIGURES|OWNS|PROVIDES
    weight      REAL DEFAULT 1.0
);
CREATE INDEX IF NOT EXISTS idx_edges_from ON graph_edges(from_node);
CREATE INDEX IF NOT EXISTS idx_edges_to   ON graph_edges(to_node);

-- ─────────────────────────────────────────────────────────
-- COMPLIANCE GAPS
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS compliance_gaps (
    id              TEXT PRIMARY KEY,
    finding_id      TEXT NOT NULL REFERENCES findings(id),
    framework       TEXT NOT NULL,        -- CNSA_2_0|NIST_800_131A|PCI_DSS_4|FIPS_140_3
    requirement_id  TEXT NOT NULL,
    requirement_desc TEXT NOT NULL,
    severity        TEXT NOT NULL,        -- MANDATORY|RECOMMENDED
    deadline        TEXT,                 -- RFC3339 or NULL if no specific date
    deadline_label  TEXT                  -- e.g. "CNSA 2.0 Exclusive Use Deadline"
);

-- ─────────────────────────────────────────────────────────
-- BASELINES (for diff / migration progress)
-- ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS baselines (
    name        TEXT PRIMARY KEY,         -- user-defined name, e.g. "Q1-2026"
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    created_at  TEXT NOT NULL,
    description TEXT
);
```

### 2.3 New `--persist` Flag (added to `spectra scan`)

When `--persist` is passed, the orchestrator writes the completed `ScanResult` and
associated graph to the state database after scanning. All Phase 2 commands read from
state; they do not re-scan unless `--re-scan` is specified.

```bash
spectra scan . --persist                     # scan and save to state
spectra scan . --persist --baseline Q1-2026  # scan, save, and tag as baseline
```

---

## 3. Updated Directory Structure — Phase 2 Additions

```
spectra/
├── ... (all Phase 1 directories unchanged)
│
├── internal/
│   ├── ... (all Phase 1 packages unchanged)
│   │
│   ├── persistence/
│   │   ├── store.go          # SQLiteStore: Open, Close, migrations, schema init
│   │   ├── scans.go          # CRUD for scans table
│   │   ├── findings.go       # CRUD + query for findings table
│   │   ├── graph.go          # CRUD + traversal for graph_nodes/graph_edges
│   │   ├── compliance.go     # CRUD for compliance_gaps table
│   │   └── baseline.go       # CRUD for baselines table
│   │
│   ├── graph/
│   │   ├── model.go          # Node, Edge, Graph types
│   │   ├── builder.go        # Build graph from ScanResult
│   │   ├── traversal.go      # BFS/DFS, reachability, blast radius
│   │   ├── export.go         # DOT (GraphViz), Mermaid, JSON-LD export
│   │   └── blast.go          # Blast radius calculator for algorithm substitution
│   │
│   ├── temporal/
│   │   ├── git.go            # Git log and blame executor + parser
│   │   ├── attribution.go    # Enrich findings with commit/author/team metadata
│   │   ├── drift.go          # Drift detector: compare current vs previous scan
│   │   └── velocity.go       # Migration velocity: trend analysis + ETA calculator
│   │
│   ├── compliance/
│   │   ├── framework.go      # ComplianceFramework type and registry
│   │   ├── cnsa20.go         # NSA CNSA 2.0 v2.1 embedded requirements + timelines
│   │   ├── nist800131a.go    # NIST SP 800-131A Rev 2 embedded requirements
│   │   ├── pcidss4.go        # PCI DSS v4.0 cryptographic requirements
│   │   ├── fips1403.go       # FIPS 140-3 approved algorithm list
│   │   ├── mapper.go         # Finding → []ComplianceGap mapper
│   │   └── timeline.go       # Deadline urgency scorer and Gantt data generator
│   │
│   ├── simulation/
│   │   ├── simulator.go      # MigrationSimulator: entry point
│   │   ├── intent.go         # MigrationIntent: {From, To} algorithm pair
│   │   ├── waves.go          # Wave planner: topological sort of dependency DAG
│   │   ├── compatibility.go  # Algorithm substitution compatibility matrix
│   │   └── impact.go         # Breaking change detector + effort estimator per node
│   │
│   ├── agility/
│   │   ├── index.go          # CAI calculator: dimensions + scoring + improvement tips
│   │   ├── abstraction.go    # Detector: are algorithms behind interfaces?
│   │   ├── centralization.go # Detector: is crypto logic centralized?
│   │   └── configurability.go# Detector: are algo choices runtime-configurable?
│   │
│   ├── endpoint/
│   │   ├── scanner.go        # EndpointScanner: orchestrator
│   │   ├── tls.go            # TLS dial + cipher suite + certificate chain analysis
│   │   ├── ssh.go            # SSH banner + host key + key exchange negotiation
│   │   └── batch.go          # Batch scan from --file=hosts.txt
│   │
│   ├── container/
│   │   ├── scanner.go        # OCI image scanner: orchestrator
│   │   ├── pull.go           # Pull from registry or load from .tar (go-containerregistry)
│   │   ├── layers.go         # Layer unpacker: stream tar → temp files → scanner pipeline
│   │   └── attribution.go    # Map finding → layer digest + image reference
│   │
│   ├── posture/
│   │   └── score.go          # CPS: Cryptographic Posture Score calculator
│   │
│   ├── evidence/
│   │   ├── enricher.go       # Enrich each Finding with standards references + guidance
│   │   └── database.go       # Load evidence.yaml into in-memory Evidence registry
│   │
│   ├── diff/
│   │   ├── diff.go           # CBOM/findings diff engine
│   │   ├── metrics.go        # Migration velocity, net QRS change, ETA
│   │   └── report.go         # Diff report terminal + HTML output
│   │
│   └── executive/
│       ├── report.go         # Executive report generator (board-ready HTML)
│       ├── narrative.go      # Business-language risk narrative generator
│       └── gantt.go          # Compliance deadline Gantt chart (HTML table)
│
└── data/
    ├── evidence.yaml         # Pre-authored evidence: explanation + refs + migration paths
    ├── cnsa20_requirements.yaml  # CNSA 2.0 requirements (machine-readable)
    ├── nist_800_131a.yaml        # NIST SP 800-131A Rev 2 requirements
    ├── pcidss4_crypto.yaml       # PCI DSS v4.0 crypto requirements
    └── compatibility_matrix.yaml # Algorithm substitution compatibility
```

---

## 4. Updated Dependency Table — Phase 2 Additions

| Dependency | Module Path | Version | Why |
|---|---|---|---|
| Persistence | `modernc.org/sqlite` | **v1.37.0** | Pure Go SQLite, no CGO, full cross-compile |
| Container registry | `github.com/google/go-containerregistry` | **latest** | OCI pull + layer unpack |
| Errgroup | `golang.org/x/sync` | latest | Already in Phase 1 |
| Netip | standard library `net` | Go 1.25 | TLS dial, SSH client |

> **Kept from Phase 1:** cobra v1.10.2, cyclonedx-go v0.9.2, lipgloss v2, gopkg.in/yaml.v3.
> **Constraint**: still zero CGO; still one binary; still no external network calls during
> non-endpoint scans. Endpoint scanning is opt-in via `spectra endpoint` subcommand only.

---

## 5. Subsystem Specifications

---

### 5.1 Cryptographic Relationship Graph (`internal/graph/`)

#### 5.1.1 Data Model

```go
type NodeType string
const (
    NodeAlgorithm  NodeType = "ALGORITHM"
    NodeFile       NodeType = "FILE"
    NodeDependency NodeType = "DEPENDENCY"
    NodeCert       NodeType = "CERT"
    NodeEndpoint   NodeType = "ENDPOINT"
    NodeTeam       NodeType = "TEAM"
    NodeLayer      NodeType = "LAYER"     // OCI image layer
)

type EdgeType string
const (
    EdgeUses       EdgeType = "USES"        // FILE uses ALGORITHM
    EdgeDependsOn  EdgeType = "DEPENDS_ON"  // FILE depends on DEPENDENCY
    EdgeProvides   EdgeType = "PROVIDES"    // DEPENDENCY provides ALGORITHM
    EdgeSigns      EdgeType = "SIGNS"       // CERT signs something
    EdgeConfigures EdgeType = "CONFIGURES"  // CONFIG configures ALGORITHM
    EdgeOwns       EdgeType = "OWNS"        // TEAM owns FILE
    EdgeServes     EdgeType = "SERVES"      // ENDPOINT serves CERT
    EdgeContains   EdgeType = "CONTAINS"    // LAYER contains FILE
)

type Node struct {
    ID         string
    Type       NodeType
    Label      string
    Properties map[string]any
}

type Edge struct {
    ID     string
    From   string    // Node ID
    To     string    // Node ID
    Type   EdgeType
    Weight float64
}

type Graph struct {
    Nodes map[string]*Node
    Edges []*Edge
    // Adjacency indices for fast traversal
    OutEdges map[string][]*Edge    // node ID → outgoing edges
    InEdges  map[string][]*Edge    // node ID → incoming edges
}
```

#### 5.1.2 Graph Builder (`graph/builder.go`)

After a scan completes, build the graph:

1. For each finding, create/upsert:
   - A **FILE node** (or CERT/ENDPOINT/DEPENDENCY node based on source)
   - An **ALGORITHM node** (one per unique `(algorithm, key_size)` pair)
   - A **USES edge** from the file/cert/endpoint to the algorithm
2. For each dependency finding, also create:
   - A **DEPENDENCY node** for the library
   - A **PROVIDES edge** from the dependency to each algorithm it provides
   - A **DEPENDS_ON edge** from files that import the dependency to the dependency node
3. For git attribution data (if temporal scan was run):
   - A **TEAM node** (grouped by git author email domain or configurable team map)
   - An **OWNS edge** from team to each file they primarily authored

#### 5.1.3 Blast Radius Calculator (`graph/blast.go`)

The blast radius of replacing algorithm A with algorithm B is the set of all nodes
that must be updated as a consequence.

```go
// BlastRadius returns the set of all nodes reachable from `algorithmID`
// via any path, along with the minimum hop count from the algorithm node.
func BlastRadius(g *Graph, algorithmID string) map[string]int {
    // BFS from the algorithm node following REVERSE edges
    // (what depends on this algorithm?)
    visited := map[string]int{algorithmID: 0}
    queue := []string{algorithmID}
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        for _, edge := range g.InEdges[current] {
            if _, seen := visited[edge.From]; !seen {
                visited[edge.From] = visited[current] + 1
                queue = append(queue, edge.From)
            }
        }
    }
    return visited
}
```

**Output**: a list of affected nodes, sorted by hop distance, with their type and label.
This is the "impact report" for any migration decision.

#### 5.1.4 Export Formats (`graph/export.go`)

**DOT (GraphViz)**:
```
digraph spectra {
    rankdir=LR;
    "RSA:2048" [shape=diamond, color=red, label="RSA-2048\nQRS: 90"];
    "auth/jwt.go" [shape=box, label="auth/jwt.go"];
    "auth/jwt.go" -> "RSA:2048" [label="USES"];
}
```

**Mermaid** (for README and HTML reports):
```
graph LR
    A[auth/jwt.go] -->|USES| B((RSA-2048\nQRS:90))
    C[pkg/crypto/keypair.go] -->|USES| B
    B -->|VIOLATES| D{CNSA 2.0}
```

**JSON-LD** (for structured data interchange):
Standard graph schema with `@type`, `@id`, and relationship predicates.

---

### 5.2 Temporal Intelligence Engine (`internal/temporal/`)

#### 5.2.1 Git Blame Integration (`temporal/git.go`)

> **Prerequisite**: the scan must be running inside a git repository. If no `.git` directory
> is found in the scan root or any parent, temporal features are silently skipped and a
> warning is emitted.

**Process for each finding with a file path:**

1. Execute: `git log --follow --diff-filter=A --format='%H|%ae|%ai|%s' -- <filepath>`
   to find the commit that first introduced the file.
2. Execute: `git blame --line-porcelain -L <line>,<line> -- <filepath>` to find the
   specific commit that introduced the matched line.
3. Parse the output into a `GitAttribution` struct:
   ```go
   type GitAttribution struct {
       CommitHash    string
       AuthorEmail   string
       AuthorName    string
       CommitDate    time.Time
       CommitMessage string
       TeamName      string    // derived from email domain or config
   }
   ```
4. Enrich the finding in the database with this data.

**Performance**: run git commands concurrently (up to `--concurrency` workers).
Cache results per `(filePath, lineNumber)` to avoid redundant git calls.

**Team derivation**: `teamMap.yaml` in the repo root (optional). If absent, derive team
from email domain (`@security.company.com` → "Security Team"). If no map and no domain
pattern, use email prefix as team name.

```yaml
# .spectra-teams.yaml (optional, place in repo root)
teams:
  - name: "Auth Team"
    email_patterns: ["@authteam.company.com", "alice@", "bob@"]
    paths: ["internal/auth/**", "pkg/jwt/**"]
  - name: "Infrastructure"
    email_patterns: ["@infra.company.com"]
    paths: ["cmd/**", "deploy/**"]
```

#### 5.2.2 Drift Detection (`temporal/drift.go`)

Compares the current scan to the most recent previous scan in the state database.

```go
type DriftReport struct {
    BaselineScanID string
    CurrentScanID  string
    PeriodDays     float64

    NewFindings      []Finding   // in current, not in baseline
    ResolvedFindings []Finding   // in baseline, not in current
    UnchangedCount   int

    NetQRSChange    int          // positive = improving
    NetBandChange   map[RiskBand]int  // per-band delta

    NewCritical     int
    ResolvedCritical int
}
```

Finding identity for comparison: use `finding_hash` (SHA256 of `algorithm + filepath + line_number`).
Two findings are "the same" if their hash matches, regardless of scan ID.

Emit a **drift alert** if:
- `NewCritical > 0` — someone introduced a new critical finding since last scan
- `NewFindings > ResolvedFindings` — net regression, findings are growing

#### 5.2.3 Migration Velocity (`temporal/velocity.go`)

Queries the last N scans from the state database and computes:

```go
type MigrationVelocity struct {
    ScansAnalyzed     int
    PeriodDays        float64
    CriticalResolved  float64    // per month
    HighResolved      float64    // per month
    NetQRSImprovment  float64    // per month
    ProjectedQRS      []QRSPoint // time-series projection
    PQCReadyDate      *time.Time // when AggQRS < 10, nil if not trending toward it
    CNSADeadline      time.Time  // 2030 for most systems
    OnTrack           bool       // PQCReadyDate < CNSADeadline
}
```

If velocity data suggests the organization is NOT on track for CNSA 2.0 by 2030,
emit a `OFF_TRACK` warning in the terminal output and executive report.

---

### 5.3 Compliance Mapping Engine (`internal/compliance/`)

#### 5.3.1 Framework Data Model

```go
type ComplianceFramework struct {
    ID          string
    Name        string
    Version     string
    Authority   string
    URL         string
    LastUpdated time.Time
    Deadlines   []Deadline
    Requirements []Requirement
}

type Requirement struct {
    ID                  string
    Description         string
    DisallowedAlgos     []string        // any finding matching these → gap
    DisallowedBelow     map[string]int  // algo → min key size (finding below = gap)
    Severity            string          // MANDATORY | RECOMMENDED
    EffectiveDate       time.Time
    ExpiresDate         *time.Time      // nil = permanent
    MandatoryReplacement []string       // recommended replacements
}

type Deadline struct {
    Label     string     // e.g. "CNSA 2.0 Network Equipment Exclusive Use"
    Date      time.Time
    Milestone string     // "prefer" | "require" | "exclusive"
    Systems   string     // description of what systems this applies to
    Urgency   string     // RED|AMBER|GREEN based on time remaining
}
```

#### 5.3.2 Embedded CNSA 2.0 Requirements (`compliance/cnsa20.go`)

Hard-code the following requirements and deadlines. Source: NSA CNSA 2.0 v2.1 (Dec 2024)
and the CNSA 2.0 Complete Guide published May 2026.

**Approved algorithms only:**
- Key Establishment: **ML-KEM-1024** (FIPS 203)
- Digital Signatures: **ML-DSA-87** (FIPS 204), **XMSS/LMS** (SP 800-208, firmware only)
- Symmetric Encryption: **AES-256** (FIPS 197)
- Hashing: **SHA-384** or **SHA-512** (FIPS 180-4)

**Disallowed under CNSA 2.0 (all RSA, all ECC, all DH, DSA, SHA-1, MD5, DES, 3DES, RC4)**

**Deadlines (embed these as structured data, not prose):**

```go
var CNSA20Deadlines = []Deadline{
    // Firmware / Software Code Signing
    {Label: "CNSA 2.0: Code Signing — Prefer",    Date: date(2025, 1, 1),  Milestone: "prefer",    Systems: "firmware and software code signing"},
    {Label: "CNSA 2.0: Code Signing — Exclusive",  Date: date(2030, 1, 1),  Milestone: "exclusive", Systems: "firmware and software code signing"},
    // Web / Cloud
    {Label: "CNSA 2.0: Web/Cloud — Prefer",        Date: date(2025, 1, 1),  Milestone: "prefer",    Systems: "web browsers and cloud services"},
    {Label: "CNSA 2.0: Web/Cloud — Exclusive",     Date: date(2033, 1, 1),  Milestone: "exclusive", Systems: "web browsers and cloud services"},
    // Network Equipment
    {Label: "CNSA 2.0: Network Equipment — Prefer", Date: date(2026, 1, 1), Milestone: "prefer",    Systems: "VPNs, routers, network infrastructure"},
    {Label: "CNSA 2.0: Network Equipment — Exclusive", Date: date(2030, 1, 1), Milestone: "exclusive", Systems: "VPNs, routers, network infrastructure"},
    // General Software / Applications
    {Label: "CNSA 2.0: Applications — Prefer",     Date: date(2027, 1, 1),  Milestone: "prefer",    Systems: "custom applications"},
    {Label: "CNSA 2.0: Applications — Require NSS", Date: date(2027, 1, 1), Milestone: "require",   Systems: "new NSS acquisitions"},
    {Label: "CNSA 2.0: Applications — Exclusive",   Date: date(2033, 1, 1), Milestone: "exclusive", Systems: "all custom applications"},
    // Ultimate goal
    {Label: "NSM-10: All Federal Systems",          Date: date(2035, 1, 1),  Milestone: "exclusive", Systems: "all US national security systems"},
}
```

**Urgency scoring for compliance gaps:**

```go
func DeadlineUrgency(deadline time.Time, now time.Time) string {
    months := deadline.Sub(now).Hours() / (24 * 30)
    switch {
    case months < 0:   return "OVERDUE"   // already passed
    case months < 12:  return "RED"       // less than 1 year
    case months < 30:  return "AMBER"     // 1–2.5 years
    default:           return "GREEN"     // more than 2.5 years
    }
}
```

#### 5.3.3 NIST SP 800-131A Rev 2 (`compliance/nist800131a.go`)

Embed the disallowed/deprecated algorithm table from NIST SP 800-131A Rev 2 (2019):

| Algorithm | Status | Effective |
|---|---|---|
| SHA-1 (for digital signatures) | DISALLOWED | Jan 1, 2014 |
| SHA-1 (for all uses) | DEPRECATED | 2030 |
| 2-key 3DES | DISALLOWED | Jan 1, 2024 |
| 3-key 3DES | DISALLOWED | Jan 1, 2024 |
| DES | DISALLOWED | Immediately |
| RC4 | DISALLOWED | Immediately |
| RSA < 2048 bits | DISALLOWED | Immediately |
| DSA (key generation) | DISALLOWED | Jan 1, 2024 |
| ECDSA < 224-bit | DISALLOWED | Immediately |
| MD5 | DISALLOWED | Immediately |
| ECDH (classical) | DEPRECATED (quantum) | Transitional |

#### 5.3.4 PCI DSS v4.0 (`compliance/pcidss4.go`)

Key requirements affecting cryptography:

| Requirement | Constraint | Spectra Finding Trigger |
|---|---|---|
| Req 4.2.1 | Strong cryptography for data in transit | TLS 1.0/1.1 detected |
| Req 4.2.1 | No RC4 | RC4 algorithm detected |
| Req 4.2.1 | No SSL | SSL protocol in use |
| Req 6.3.3 | Patching within 1 month (critical) | Outdated crypto library detected |
| Req 12.3.3 | Annual cryptography review | (Spectra satisfies this requirement) |
| General | SHA-1 certificates disallowed | SHA-1 signing algorithm in cert |

#### 5.3.5 Compliance Gap Mapper (`compliance/mapper.go`)

For each finding, iterate through all enabled frameworks and their requirements.
If the finding's algorithm is in a framework's `DisallowedAlgos` list (and the
`EffectiveDate` has passed), create a `ComplianceGap`.

Return gaps sorted by urgency (OVERDUE first, then RED, AMBER, GREEN).

CLI flag: `--frameworks` (comma-separated list of framework IDs to check).
Default: all frameworks.

---

### 5.4 Migration Simulation Engine (`internal/simulation/`)

#### 5.4.1 Algorithm Compatibility Matrix (`data/compatibility_matrix.yaml`)

Before simulating a migration, the engine must know whether the target algorithm is
compatible with each usage context. Embed this matrix as YAML data:

```yaml
# Can algorithm A be substituted with algorithm B in context C?
# true = compatible | false = incompatible | hybrid = hybrid approach required

substitutions:
  RSA → ML-KEM:
    key_establishment: true      # ML-KEM is a key encapsulation mechanism
    digital_signature: false     # RSA signs; ML-KEM does not — use ML-DSA instead
    certificate_pubkey: true     # ML-KEM public keys can appear in certificates
    tls_key_exchange: true       # ML-KEM is supported in TLS 1.3 (RFC 9180)
    jwt_signing: false           # ML-KEM does not sign JWTs; no standard yet
    ssh_host_key: false          # ML-KEM cannot be an SSH host key; use ML-DSA

  RSA → ML-DSA:
    key_establishment: false     # ML-DSA signs, does not encapsulate
    digital_signature: true
    certificate_pubkey: true
    jwt_signing: false           # no JWT standard for ML-DSA yet → flag as PENDING
    code_signing: true

  ECDSA → ML-DSA:
    digital_signature: true
    certificate_pubkey: true
    jwt_signing: false           # pending JOSE ML-DSA RFC
    code_signing: true
    tls_certificate: true

  ECDH → ML-KEM:
    key_establishment: true
    tls_key_exchange: true

  SHA1 → SHA256:
    hmac: true
    certificate_signature: true  # requires cert re-issuance
    checksum: true
    pbkdf: true

  MD5 → SHA256:
    hmac: true
    checksum: true
    password_hash: false         # MD5 for passwords → bcrypt/Argon2, not SHA-256

  AES128 → AES256:
    all_contexts: true           # drop-in replacement, same API
```

#### 5.4.2 Simulation Entry Point (`simulation/simulator.go`)

```go
type MigrationIntent struct {
    From        string    // source algorithm canonical name
    To          string    // target algorithm canonical name
    DryRun      bool      // if true, calculate only, do not write to DB
    ScanID      string    // which scan to base the simulation on
}

type SimulationResult struct {
    Intent         MigrationIntent
    Waves          []MigrationWave
    BreakingChanges []BreakingChange
    IncompatibleContexts []ContextIncompatibility
    TotalNodes     int
    EstimatedWeeks int         // rough estimate from effort × node count
    QRSImpact      int         // new aggregate QRS after migration
}

type MigrationWave struct {
    Number      int
    Description string
    Nodes       []WaveNode
    EstimatedDays int
}

type BreakingChange struct {
    NodeLabel   string
    NodeType    NodeType
    Reason      string
    Workaround  string
}
```

#### 5.4.3 Wave Planner (`simulation/waves.go`)

1. Load all findings for the `From` algorithm from the graph.
2. Build a dependency DAG: for each finding, add edges from findings that "use" it
   to findings that "depend on" the same dependency or cert.
3. Topological sort the DAG → nodes with no dependents go into Wave 1.
4. Within each wave, sort by:
   - Effort level (EASY first)
   - QRS descending (highest risk first within same effort)
5. Assign estimated days per wave: EASY=1-2d, MEDIUM=3-7d, HARD=14-30d, BLOCKED=defer.

#### 5.4.4 Breaking Change Detector (`simulation/impact.go`)

For each node about to be migrated, check the compatibility matrix for the usage context.
Usage context is inferred from:
- The matched pattern in the source code (e.g., `createSign` → `digital_signature`)
- The certificate Extended Key Usage (code signing, TLS server, etc.)
- The dependency's known context from the dependency registry

If `compatibility_matrix[from → to][context] == false`, create a `BreakingChange`
with a workaround suggestion.

---

### 5.5 Cryptographic Agility Index (`internal/agility/`)

#### 5.5.1 Definition

The **Cryptographic Agility Index (CAI)** measures how easily an organization can
replace cryptographic algorithms in the future, independent of what algorithms they
currently use. A codebase with RSA but high CAI can migrate faster than a codebase
with RSA and low CAI.

Scale: **0–100**. Higher = more agile.

#### 5.5.2 Measurement Dimensions (25 points each)

**Dimension 1: Abstraction** (`agility/abstraction.go`)

Are algorithms accessed through abstractions (interfaces, factories, helper functions)
or referenced directly as magic strings/constants?

Scoring signals (positive):
- Algorithm names passed through configuration files rather than hardcoded (+5 per config-driven algo)
- Crypto utility package exists (`pkg/crypto/`, `internal/security/`, `utils/crypto.go`) (+5)
- Algorithm selection via interface/strategy pattern (e.g., `type Signer interface`) (+5)
- No algorithm name appears in more than 3 files (+5)
- All key generation goes through one function (+5)

Scoring penalties (negative):
- Same algorithm string literal (`"RSA"`, `"SHA-1"`) in 5+ files (-5 each, max -20)
- Algorithm hardcoded in test + production + config all separately (-5)

**Dimension 2: Centralization** (`agility/centralization.go`)

Is crypto logic consolidated or scattered?

Scoring signals (positive):
- Single crypto package imports from a central utility (+10)
- Key material generated/loaded in one place (+5)
- Cipher suite configuration in one config file (+5)
- Evidence of key rotation logic (+5)

Scoring penalties:
- Crypto imports scattered across >5 different packages (-5)
- Key generation logic duplicated across files (-5 each, max -15)

**Dimension 3: Configurability** (`agility/configurability.go`)

Can algorithms be changed without code modifications?

Scoring signals (positive):
- Algorithm names appear in config files that are also scanned in findings (+5 per case)
- TLS cipher suites configurable via YAML/ENV (+5)
- Key size configurable via environment variable or config (+5)
- Feature flags present for crypto changes (+5)
- Environment-specific crypto configs (dev vs prod) (+5)

**Dimension 4: Migration Readiness** (`index.go`)

Fraction of findings that are actionable:

```
readiness = (EASY + MEDIUM) / total_findings
score = readiness × 25
```

If any PQC-safe algorithms (ML-KEM, ML-DSA, SLH-DSA) are detected in the codebase,
add a +5 readiness bonus (already started migrating).

#### 5.5.3 CAI Output

```
Cryptographic Agility Index: 34/100 — LOW AGILITY

Dimension Breakdown:
  Abstraction:      8/25  — Algorithms hardcoded in 14 files (major refactoring needed)
  Centralization:  12/25  — Crypto logic scattered across 8 packages
  Configurability:  6/25  — No algorithm configuration detected
  Migration Ready:  8/25  — 64% of findings are HARD or BLOCKED effort

Top improvement actions:
  1. Centralize all crypto into pkg/crypto/ (+12 pts)
  2. Move algorithm names to config.yaml (+8 pts)
  3. Replace direct library calls with interface wrappers (+6 pts)

Projected CAI after improvements: 68/100 (MODERATE)
```

---

### 5.6 Evidence-Based Finding Enrichment (`internal/evidence/`)

#### 5.6.1 Evidence Database (`data/evidence.yaml`)

Pre-author this file. It is embedded at compile time using Go's `//go:embed` directive.
One entry per algorithm, with full standards citations.

```yaml
# data/evidence.yaml
RSA:
  plain_english: >
    RSA is based on the difficulty of factoring large integers. A sufficiently powerful
    quantum computer running Shor's algorithm can factor any RSA key — regardless of
    key size — in polynomial time. RSA-2048 offers no post-quantum security.
    The "Harvest Now, Decrypt Later" (HNDL) threat means adversaries collecting RSA-encrypted
    traffic today can decrypt it once a quantum computer is available.
  standards:
    - framework: "NIST SP 800-131A Rev 2"
      section: "Section 3"
      clause: "RSA key pairs having a security strength of less than 112 bits are disallowed."
      url: "https://csrc.nist.gov/publications/detail/sp/800-131a/rev-2/final"
    - framework: "NSA CNSA 2.0"
      section: "Algorithm List"
      clause: "RSA is not approved for use in National Security Systems under CNSA 2.0."
      url: "https://media.defense.gov/2022/Sep/07/2003071834/-1/-1/0/CSA_CNSA_2.0_ALGORITHMS_AND_KEY_SIZES.PDF"
  migration_path:
    - "Inventory all RSA key usages (completed by Spectra scan)."
    - "Determine usage context: key establishment or digital signature."
    - "For key establishment: replace with ML-KEM-1024 (FIPS 203)."
    - "For digital signatures: replace with ML-DSA-87 (FIPS 204)."
    - "Update certificate infrastructure: request new CA-signed ML-KEM/ML-DSA certificates."
    - "Update TLS configuration to prefer ML-KEM for key exchange."
    - "Update clients/consumers to accept new key types before revoking RSA keys."
    - "Rotate all existing RSA keys and destroy old key material."
  code_before_go: |
    key, _ := rsa.GenerateKey(rand.Reader, 2048)
    sig, _ := rsa.SignPSS(rand.Reader, key, crypto.SHA256, digest, nil)
  code_after_go: |
    // github.com/open-quantum-safe/liboqs-go or cloudflare/circl
    sk, pk, _ := mldsago.KeyGen()
    sig := mldsago.Sign(sk, message)

SHA1:
  plain_english: >
    SHA-1 was broken by researchers in 2017 (SHAttered attack) using chosen-prefix
    collision. Generating SHA-1 collisions is now a practical attack, not a theoretical
    one. NIST deprecated SHA-1 for all digital signature applications in 2014.
    While SHA-1 for non-signature uses (checksums) remained allowed, NIST SP 800-131A
    Rev 2 deprecated it further. All SHA-1 usage should be migrated to SHA-256 or SHA-3.
  standards:
    - framework: "NIST SP 800-131A Rev 2"
      section: "Section 9"
      clause: "The use of SHA-1 is deprecated for all applications."
      url: "https://csrc.nist.gov/publications/detail/sp/800-131a/rev-2/final"
  migration_path:
    - "Replace hashlib.sha1 / sha1.New() with sha256.New() or sha3.New256()."
    - "Update all certificate signing requests to use SHA-256 or SHA-384."
    - "Re-issue any certificates signed with SHA1withRSA."
  code_before_go: |
    h := sha1.New()
    h.Write(data)
    digest := h.Sum(nil)
  code_after_go: |
    h := sha256.New()
    h.Write(data)
    digest := h.Sum(nil)

# ... (MD5, ECDSA, DES, 3DES, RC4, DH, AES128 all follow same structure)
```

#### 5.6.2 Enricher (`evidence/enricher.go`)

After the scan completes, for each finding, look up its algorithm in the evidence
database and attach the evidence to the `Finding` struct before writing to state.

The enriched finding displays in:
- Terminal: collapsed by default; expand with `--verbose` or `--evidence`
- JSON: always included
- HTML report: shown as expandable accordion per finding
- Executive report: simplified version with business-language explanation only

---

### 5.7 Live Infrastructure Interrogation (`internal/endpoint/`)

#### 5.7.1 TLS Scanner (`endpoint/tls.go`)

**Input**: hostname, port (default: 443). Optionally read from `--file=hosts.txt`
(one `host:port` per line).

**Process:**
1. Dial TLS with a permissive `tls.Config` (all cipher suites, TLS 1.0 through 1.3,
   `InsecureSkipVerify: false` — use system trust store).
2. Record:
   - `ConnectionState.Version` (TLS version negotiated)
   - `ConnectionState.CipherSuite` (negotiated cipher suite)
   - `ConnectionState.PeerCertificates` (full chain)
   - Supported cipher suites: probe by negotiating with each suite explicitly
     (use separate dials with forced `CipherSuites` and `MaxVersion` settings)
3. For each certificate in the chain, run the certificate scanner (Phase 1 logic).
4. Check for:
   - TLS 1.0 or 1.1 support (HIGH vulnerability)
   - NULL cipher suites (CRITICAL)
   - Export-grade cipher suites (CRITICAL)
   - RC4 cipher suites (CRITICAL)
   - 3DES cipher suites (HIGH)
   - Non-forward-secret key exchange (no ECDHE/DHE) (HIGH)
   - Certificate using SHA-1 signature (HIGH)
   - Certificate expiring within 30 days (MEDIUM)
   - No Certificate Transparency SCTs (LOW)

**Timeout**: 10 seconds per host. Skip gracefully if timeout.
**Concurrency**: use worker pool from Phase 1 orchestrator.

#### 5.7.2 SSH Scanner (`endpoint/ssh.go`)

**Input**: hostname, port (default: 22).

Use `golang.org/x/crypto/ssh` to connect and read:
- Server host key type (RSA, ECDSA, Ed25519)
- Supported key exchange algorithms
- Supported host key algorithms
- Supported MAC algorithms

Flag RSA host keys (CRITICAL), DSA host keys (CRITICAL),
ECDSA host keys (HIGH), Ed25519 (SAFE, though not PQC).

**Note**: do not authenticate — just complete the handshake banner exchange and
read the server's advertised capabilities.

---

### 5.8 OCI Intelligence Layer (`internal/container/`)

#### 5.8.1 Image Sources

Accept:
- `spectra container scan --image=myapp:latest` (pull from local Docker daemon)
- `spectra container scan --image=gcr.io/project/api:v1` (pull from registry)
- `spectra container scan --tar=myimage.tar` (load from saved tar file)

Use `github.com/google/go-containerregistry`:
- `crane.Load()` for tar files
- `crane.Pull()` for registry images (requires auth via `~/.docker/config.json`)

#### 5.8.2 Layer Processing (`container/layers.go`)

1. Load image → iterate over layers in order.
2. For each layer, decompress (gzip) and untar to a temp directory.
3. Pass each extracted file to the Phase 1 scanner pipeline (code, cert, config, deps).
4. Augment findings with `container_layer = <layer digest>` and `container_image = <ref>`.
5. Clean up temp directory after each layer.

**Important**: handle whiteout files (`.wh.` prefix) — these represent deleted files
in overlay layers. Mark findings from whiteout layers as RESOLVED.

**Output enhancement**: show image layer in the terminal output:
```
Layer sha256:abc123 [/app/server] RSA   auth/handler.go:47   CRITICAL QRS:90
```

---

### 5.9 Migration Progress Intelligence (`internal/diff/`)

#### 5.9.1 CBOM Diff Engine (`diff/diff.go`)

**Input**: two CBOM JSON files or two scan IDs from the state database.

```go
type DiffResult struct {
    Before ScanSummary
    After  ScanSummary
    Period time.Duration

    Resolved   []FindingDiff    // in before, not in after
    Introduced []FindingDiff    // in after, not in before
    Unchanged  int

    QRSDelta       int          // negative = improving
    BandDelta      map[RiskBand]int
    CPSDelta        int

    Velocity       MigrationVelocity
    ComplianceGapsDelta int    // net change in compliance gaps
}
```

#### 5.9.2 Migration Velocity (`diff/metrics.go`)

```
MonthlyResolutionRate = len(Resolved) / (Period.Days / 30)

ProjectedMonthsToSafe = (CurrentCriticalCount + CurrentHighCount) / MonthlyResolutionRate

PQCReadyDate = now + ProjectedMonthsToSafe

OnTrack = PQCReadyDate < nearestCNSA20Deadline
```

**Terminal output example:**
```
SPECTRA DIFF — 17-month migration progress
────────────────────────────────────────────
Before (Jan 2025): QRS=73, CPS=31 (WEAK)
After  (May 2026): QRS=48, CPS=58 (DEVELOPING)
Net improvement:   QRS -25 pts  CPS +27 pts  ✓ IMPROVING

Resolved:  23 findings (12 CRITICAL, 8 HIGH, 3 MEDIUM)
New:        7 findings (2 CRITICAL, 5 HIGH)  ← requires attention
Unchanged: 47 findings

Migration velocity: ~1.4 critical findings resolved/month
Projected PQC-ready: July 2029 (38 months)
CNSA 2.0 (applications): Jan 2033 — ON TRACK ✓

New critical findings (action required):
  1. RSA-2048 in new-service/api/handler.go  (introduced in v2.3.0, Jan 2026)
  2. SHA-1 in vendor/third-party-sdk@v1.2.0  (dependency update)
```

---

### 5.10 Cryptographic Posture Score (`internal/posture/`)

#### 5.10.1 CPS Formula

```go
func ComputeCPS(r ScanResult, cai int, complianceGaps int, totalFindings int, qrsDelta int) int {
    // Invert QRS: high QRS = low posture score
    qrsContrib := (100 - r.AggregateQRS) * 0.40

    // Agility: higher agility = better posture
    agilityContrib := float64(cai) * 0.25

    // Compliance: fewer gaps = better posture
    complianceScore := 100 - min(100, complianceGaps*5)  // each gap costs 5 pts, max 100 penalty
    complianceContrib := float64(complianceScore) * 0.20

    // Drift: positive = improving (bonus), negative = regressing (penalty)
    driftContrib := 0.0
    if qrsDelta < 0 {  // getting better
        driftContrib = float64(min(10, -qrsDelta)) * 0.15
    } else {  // getting worse
        driftContrib = float64(max(-10, -qrsDelta)) * 0.15
    }

    cps := qrsContrib + agilityContrib + complianceContrib + driftContrib
    return max(0, min(100, int(math.Round(cps))))
}
```

#### 5.10.2 CPS Interpretation

| Score | Band | Description |
|---|---|---|
| 80–100 | STRONG | PQC-ready or nearly so; minor gaps remain |
| 60–79 | ADEQUATE | Manageable exposure; migration underway or feasible |
| 40–59 | DEVELOPING | Significant gaps; migration not yet underway |
| 20–39 | WEAK | Major cryptographic debt; urgent action required |
| 0–19 | CRITICAL | Systemic quantum exposure; remediation overdue |

---

### 5.11 Executive Audit Report Generator (`internal/executive/`)

The executive report must speak the language of business risk, not cryptographic theory.
Every technical fact must be translated into organizational consequence.

#### 5.11.1 Report Sections

**Cover Page:**
- Organization name (from `--org-name` flag)
- Scan date range
- "Prepared by Spectra v{version} — Cryptographic Intelligence Platform"
- Classification banner (CONFIDENTIAL / INTERNAL / PUBLIC) via `--classification` flag

**Section 1 — Executive Summary (1 page)**
- CPS traffic light (RED/AMBER/GREEN)
- 3 bullet risks: "Your organization has X quantum-vulnerable assets protecting Y services"
- 3 bullet actions: "The highest-priority action is..."
- Compliance status: "You meet/do not meet CNSA 2.0 Year 1 requirements"

**Section 2 — Cryptographic Risk Landscape**
- Doughnut chart: findings by algorithm family
- Table: top 10 most critical findings with business context
- Comparison to industry benchmarks (use synthetic anonymised data in v1)

**Section 3 — Compliance Alignment Matrix**

```
Framework           │ Required By │ Status      │ Gaps
────────────────────┼─────────────┼─────────────┼──────
NIST SP 800-131A    │ Now         │ NON-COMPLIANT│ 23
NSA CNSA 2.0 Pref.  │ 2027        │ NOT MET      │ 47
NSA CNSA 2.0 Excl.  │ 2033        │ ON TRACK     │ 47
PCI DSS v4.0 4.2.1  │ Mar 2025    │ NON-COMPLIANT│  3
FIPS 140-3          │ Ongoing     │ PARTIAL      │ 12
```

**Section 4 — Migration Roadmap**
Gantt-style table (HTML) showing:
- Wave 1 (Immediate / EASY): items that can be resolved in Sprint 1–2
- Wave 2 (Short-term / MEDIUM): 1–3 months
- Wave 3 (Medium-term / HARD): 3–12 months
- Wave 4 (Long-term / BLOCKED): dependency on third parties

**Section 5 — Cryptographic Agility Index**
Narrative: "Your organization would need approximately X months and Y engineers to
complete a full PQC migration. The primary obstacle is Z."

**Section 6 — Team Ownership**
Bar chart: which teams own the most quantum debt.
(Only shows if git attribution data is available.)

**Section 7 — Risk Narrative**
Auto-generated business-language text from `executive/narrative.go`:

```go
// RiskNarrative generates business-language risk description
func RiskNarrative(r ScanResult) string {
    critCount := r.FindingsByBand[CRITICAL]
    return fmt.Sprintf(
        "Your organization has %d cryptographic assets that quantum computers "+
        "will be able to compromise. These assets currently protect data transmitted "+
        "across your %d scanned services. Under the NSA's CNSA 2.0 mandate, organizations "+
        "serving national security systems must eliminate these vulnerabilities by "+
        "January 2033. At the current rate of remediation, your organization is "+
        "%s to meet this deadline.",
        critCount, inferServiceCount(r), inferTrackStatus(r),
    )
}
```

**Output**: `spectra-executive-report.html` — self-contained, printable (print CSS included).

---

## 6. New CLI Commands

```
spectra graph [flags]
  --format      dot|mermaid|json-ld (default: terminal)
  --scan-id     which scan to graph (default: latest)
  --algorithm   filter to a specific algorithm
  --depth       max traversal depth for blast radius (default: unlimited)

spectra blast [flags]
  --algorithm   algorithm to analyze (required)
  --replace-with target algorithm (optional, shows compatibility)
  --scan-id     which scan to use (default: latest)

spectra simulate [flags]
  --from        source algorithm (required)
  --to          target algorithm (required)
  --scan-id     which scan to base the simulation on (default: latest)
  --format      terminal|json|html

spectra diff [flags]
  --before      CBOM file or scan ID (required)
  --after       CBOM file or scan ID (default: latest)
  --format      terminal|json|html

spectra compliance [flags]
  --frameworks  comma-separated framework IDs (default: all)
  --scan-id     which scan to analyze (default: latest)
  --format      terminal|json|html
  --deadline-filter  overdue|red|amber|green (filter by urgency)

spectra endpoint scan [flags]
  --url         URL to scan (https://example.com)
  --host        host:port to scan directly
  --file        file with one host:port per line
  --protocol    tls|ssh|both (default: tls)
  --format      terminal|json|cbom
  --persist     save endpoint findings to state

spectra agility [flags]
  --scan-id     which scan to analyze (default: latest)
  --format      terminal|json
  --tips        show improvement recommendations

spectra history [flags]
  --path        file or directory path to analyze
  --blame       enrich with git blame data (requires git repo)
  --team-map    path to .spectra-teams.yaml
  --format      terminal|json

spectra container scan [flags]
  --image       image reference (e.g. myapp:latest, gcr.io/project/api:v1)
  --tar         path to saved .tar image file
  --format      terminal|json|cbom|html
  --persist     save container findings to state

spectra report executive [flags]
  --scan-id     which scan to use (default: latest)
  --org-name    organization name for cover page
  --classification  CONFIDENTIAL|INTERNAL|PUBLIC
  --out-dir     output directory (default: ./spectra-out)

spectra posture [flags]
  --scan-id     which scan to analyze (default: latest)
  --history     show posture trend over last N scans (default: 5)
```

Also, add `--persist` and `--baseline <name>` flags to the existing `spectra scan` command.

---

## 7. Phase 2 Atomic Commits

Each commit below is one atomic PR-ready commit following Conventional Commits.

```
feat(persistence): add modernc.org/sqlite state database with schema migrations
feat(persistence): implement FindingsStore CRUD and deduplication by hash
feat(persistence): implement GraphStore with node/edge CRUD
feat(persistence): implement ComplianceStore and BaselineStore
feat(cli/scan): add --persist and --baseline flags to scan command
feat(graph): define Node/Edge/Graph types and build tags
feat(graph): implement graph builder from ScanResult (nodes + edges)
feat(graph): implement BFS blast radius calculator
feat(graph): add GraphViz DOT and Mermaid diagram export
feat(graph): add JSON-LD export for interoperability
feat(cli): add spectra graph subcommand with --format flag
feat(cli): add spectra blast subcommand with compatibility check
feat(temporal): implement git log parser and GitAttribution extractor
feat(temporal): enrich findings with commit, author, date, team metadata
feat(temporal): implement drift detector comparing current vs prior scan
feat(temporal): implement migration velocity and PQC-ready ETA calculator
feat(cli): add spectra history subcommand with --blame flag
feat(compliance): implement ComplianceFramework and Requirement data model
feat(compliance): embed NIST SP 800-131A Rev 2 requirements and disallowed algorithms
feat(compliance): embed NSA CNSA 2.0 v2.1 requirements and milestone deadlines
feat(compliance): embed PCI DSS v4.0 cryptographic requirements
feat(compliance): embed FIPS 140-3 approved algorithm list
feat(compliance): implement compliance gap mapper (finding → violations)
feat(compliance): implement deadline urgency scorer (OVERDUE/RED/AMBER/GREEN)
feat(cli): add spectra compliance subcommand with --frameworks flag
feat(simulation): implement MigrationIntent and compatibility matrix loader
feat(simulation): implement dependency DAG builder from graph store
feat(simulation): implement topological sort wave planner
feat(simulation): implement breaking change detector using compatibility matrix
feat(simulation): add effort estimator and timeline calculator per wave
feat(cli): add spectra simulate subcommand with --from and --to flags
feat(agility): implement Cryptographic Agility Index framework (4 dimensions)
feat(agility): implement abstraction detector (interface and factory patterns)
feat(agility): implement centralization detector (crypto package analysis)
feat(agility): implement configurability detector (config-driven algorithm detection)
feat(agility): add improvement tips generator with projected CAI score
feat(cli): add spectra agility subcommand
feat(evidence): embed evidence.yaml with standards references for all algorithms
feat(evidence): implement finding enricher that attaches evidence post-scan
feat(posture): implement CPS formula combining QRS, CAI, compliance, and drift
feat(posture): add posture history trend from state database
feat(cli): add spectra posture subcommand with --history flag
feat(endpoint): implement TLS scanner with cipher suite analysis
feat(endpoint): implement certificate chain analyser
feat(endpoint): add SSH key exchange scanner
feat(endpoint): implement batch host scanning from file
feat(cli): add spectra endpoint subcommand
feat(container): implement OCI image pull and tar loading via go-containerregistry
feat(container): implement layer-by-layer tar extraction and scanner pipeline wiring
feat(container): add whiteout file handling and layer attribution
feat(cli): add spectra container subcommand
feat(diff): implement CBOM and state-based diff engine
feat(diff): implement migration velocity metrics and PQC-ready date projection
feat(diff): add diff terminal output with net change summary
feat(cli): add spectra diff subcommand
feat(executive): implement board-ready HTML executive report generator
feat(executive): implement business-language risk narrative generator
feat(executive): add compliance alignment matrix section
feat(executive): add migration roadmap Gantt table section
feat(executive): add team attribution section (conditional on git data)
feat(cli): add spectra report executive subcommand
test: add unit tests for graph builder and blast radius calculator
test: add unit tests for compliance gap mapper with golden fixtures
test: add unit tests for simulation wave planner and breaking change detector
test: add unit tests for CAI calculator per dimension
test: add unit tests for CPS formula
test: add unit tests for CBOM diff engine
test: add integration test for full Phase 2 pipeline (scan → persist → compliance → simulate)
docs: update README with Phase 2 architecture and capability overview
docs: add COMPLIANCE.md documenting all embedded framework coverage with citations
docs: add ARCHITECTURE.md with relationship graph entity model diagram
ci: extend GitHub Actions to run integration tests with temp SQLite state
```

---

## 8. Non-Negotiable Quality Rules for Phase 2

All Phase 1 rules apply, plus:

1. **Graph integrity**: every edge must reference valid node IDs. Enforce with foreign
   key constraints in SQLite (`PRAGMA foreign_keys = ON`). Any code that inserts graph
   data must use transactions to maintain referential integrity.

2. **Compliance data accuracy**: every requirement embedded in `compliance/*.go` must
   include a `SourceURL` field pointing to the authoritative public document. No
   compliance requirements may be invented or paraphrased without a source citation.

3. **No false compliance claims**: never output "COMPLIANT" for a framework unless
   zero findings violate any of its requirements. Use "NO VIOLATIONS FOUND" instead
   — there may be violations Spectra could not detect.

4. **Deterministic graph IDs**: node IDs must be derived deterministically from the
   node's unique properties (e.g., `SHA256(nodeType + ":" + label)`), not random UUIDs.
   This ensures that the same scan produces the same graph, and that diff operations
   correctly identify matching nodes across scans.

5. **Temporal data is optional, never blocking**: if a git repository is not available,
   git attribution features silently degrade to showing "unknown" — they must never
   cause the scan to fail.

6. **Executive report accuracy**: the narrative generator must not invent organizational
   facts. Every quantitative claim in the executive report must be derivable from the
   scan data. No hallucinated statistics.

7. **Container scanning is ephemeral**: layer extraction must use `os.MkdirTemp` and
   call `defer os.RemoveAll(tempDir)` immediately after creation. Never leave extracted
   image contents on disk after scanning.

8. **The state database is the user's data**: on `spectra purge`, ask for confirmation
   before deleting state. Never delete state as a side effect of any other command.
   Expose `spectra state export` for backup.

---

## 9. Summary of Phase 2 Strategic Value

| Dimension | Phase 1 | Phase 2 |
|---|---|---|
| Scope | Single scan | Persistent intelligence platform |
| Output | Report | Decision support system |
| Time horizon | Point in time | Historical + predictive |
| Addressee | Security engineer | CISO, CTO, Board |
| Unique asset | None | Cryptographic Knowledge Graph |
| Defensibility | Replicable | Graph data compounds over time |
| Enterprise fit | Utility | Governance infrastructure |
| Key differentiator | QRS scoring | CAI + CPS + Simulation + Evidence chain |

**The core insight:** every enterprise security team can build a scanner in weeks.
No team can rebuild an organization's accumulated cryptographic relationship graph,
temporal history, compliance gap ledger, and migration simulation model in anything
less than months. That is the moat. That is what Spectra Phase 2 builds.

---

*Phase 2 build prompt authored: May 2026 — based on CNSA 2.0 v2.1 (Dec 2024),
NIST FIPS 203/204/205 (Aug 2024), NIST SP 800-131A Rev 2, PCI DSS v4.0,
and the strategic reframing of Spectra as a Cryptographic Intelligence Platform.*
