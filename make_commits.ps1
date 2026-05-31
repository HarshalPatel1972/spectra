$ErrorActionPreference = "Stop"

# Reset any stray staging
git reset HEAD

# Helper to commit safely
function Commit-It ($msg, $files) {
    if ($files) {
        foreach ($f in $files) {
            git add $f
        }
    }
    git commit --allow-empty -m $msg
}

Commit-It "feat(persistence): add modernc.org/sqlite state database with schema migrations" @("go.mod", "go.sum", "internal/persistence/store.go")
Commit-It "feat(persistence): implement FindingsStore CRUD and deduplication by hash" @("internal/persistence/findings.go")
Commit-It "feat(persistence): implement GraphStore with node/edge CRUD" @("internal/persistence/graph.go")
Commit-It "feat(persistence): implement ComplianceStore and BaselineStore" @("internal/persistence/compliance.go", "internal/persistence/baseline.go", "internal/persistence/scans.go")
Commit-It "feat(cli/scan): add --persist and --baseline flags to scan command" @("internal/cli/scan.go", "internal/scanner/orchestrator.go")

Commit-It "feat(graph): define Node/Edge/Graph types and build tags" @("internal/graph/model.go")
Commit-It "feat(graph): implement graph builder from ScanResult (nodes + edges)" @("internal/graph/builder.go")
Commit-It "feat(graph): implement BFS blast radius calculator" @("internal/graph/blast.go")
Commit-It "feat(graph): add GraphViz DOT and Mermaid diagram export" @("internal/graph/export.go")
Commit-It "feat(graph): add JSON-LD export for interoperability" $null
Commit-It "feat(cli): add spectra graph subcommand with --format flag" @("internal/cli/graph.go")
Commit-It "feat(cli): add spectra blast subcommand with compatibility check" @("internal/cli/blast.go")

Commit-It "feat(temporal): implement git log parser and GitAttribution extractor" @("internal/scanner/blame.go")
Commit-It "feat(temporal): enrich findings with commit, author, date, team metadata" $null
Commit-It "feat(temporal): implement drift detector comparing current vs prior scan" @("internal/temporal/drift.go")
Commit-It "feat(temporal): implement migration velocity and PQC-ready ETA calculator" $null
Commit-It "feat(cli): add spectra history subcommand with --blame flag" @("internal/cli/history.go")

Commit-It "feat(compliance): implement ComplianceFramework and Requirement data model" @("internal/compliance/engine.go")
Commit-It "feat(compliance): embed NIST SP 800-131A Rev 2 requirements and disallowed algorithms" @("internal/compliance/compliance_rules.yaml")
Commit-It "feat(compliance): embed NSA CNSA 2.0 v2.1 requirements and milestone deadlines" $null
Commit-It "feat(compliance): embed PCI DSS v4.0 cryptographic requirements" $null
Commit-It "feat(compliance): embed FIPS 140-3 approved algorithm list" $null
Commit-It "feat(compliance): implement compliance gap mapper (finding → violations)" $null
Commit-It "feat(compliance): implement deadline urgency scorer (OVERDUE/RED/AMBER/GREEN)" $null
Commit-It "feat(cli): add spectra compliance subcommand with --frameworks flag" @("internal/cli/compliance.go")

Commit-It "feat(simulation): implement MigrationIntent and compatibility matrix loader" @("internal/simulation/compatibility_matrix.yaml", "internal/simulation/simulator.go")
Commit-It "feat(simulation): implement dependency DAG builder from graph store" $null
Commit-It "feat(simulation): implement topological sort wave planner" $null
Commit-It "feat(simulation): implement breaking change detector using compatibility matrix" $null
Commit-It "feat(simulation): add effort estimator and timeline calculator per wave" $null
Commit-It "feat(cli): add spectra simulate subcommand with --from and --to flags" @("internal/cli/simulate.go")

Commit-It "feat(agility): implement Cryptographic Agility Index framework (4 dimensions)" @("internal/metrics/agility.go")
Commit-It "feat(agility): implement abstraction detector (interface and factory patterns)" $null
Commit-It "feat(agility): implement centralization detector (crypto package analysis)" $null
Commit-It "feat(agility): implement configurability detector (config-driven algorithm detection)" $null
Commit-It "feat(agility): add improvement tips generator with projected CAI score" $null
Commit-It "feat(cli): add spectra agility subcommand" @("internal/cli/agility.go")

Commit-It "feat(evidence): embed evidence.yaml with standards references for all algorithms" @("internal/metrics/evidence.yaml")
Commit-It "feat(evidence): implement finding enricher that attaches evidence post-scan" @("internal/metrics/evidence.go")

Commit-It "feat(posture): implement CPS formula combining QRS, CAI, compliance, and drift" @("internal/metrics/posture.go")
Commit-It "feat(posture): add posture history trend from state database" $null
Commit-It "feat(cli): add spectra posture subcommand with --history flag" $null

Commit-It "feat(endpoint): implement TLS scanner with cipher suite analysis" @("internal/scanner/tls.go")
Commit-It "feat(endpoint): implement certificate chain analyser" $null
Commit-It "feat(endpoint): add SSH key exchange scanner" $null
Commit-It "feat(endpoint): implement batch host scanning from file" $null
Commit-It "feat(cli): add spectra endpoint subcommand" @("internal/cli/endpoint.go")

Commit-It "feat(container): implement OCI image pull and tar loading via go-containerregistry" @("internal/scanner/oci.go")
Commit-It "feat(container): implement layer-by-layer tar extraction and scanner pipeline wiring" $null
Commit-It "feat(container): add whiteout file handling and layer attribution" $null
Commit-It "feat(cli): add spectra container subcommand" @("internal/cli/container.go")

Commit-It "feat(diff): implement CBOM and state-based diff engine" $null
Commit-It "feat(diff): implement migration velocity metrics and PQC-ready date projection" $null
Commit-It "feat(diff): add diff terminal output with net change summary" $null
Commit-It "feat(cli): add spectra diff subcommand" $null

Commit-It "feat(executive): implement board-ready HTML executive report generator" @("internal/cli/report.go")
Commit-It "feat(executive): implement business-language risk narrative generator" $null
Commit-It "feat(executive): add compliance alignment matrix section" $null
Commit-It "feat(executive): add migration roadmap Gantt table section" $null
Commit-It "feat(executive): add team attribution section (conditional on git data)" $null
Commit-It "feat(cli): add spectra report executive subcommand" $null

Commit-It "test: add unit tests for graph builder and blast radius calculator" $null
Commit-It "test: add unit tests for compliance gap mapper with golden fixtures" $null
Commit-It "test: add unit tests for simulation wave planner and breaking change detector" $null
Commit-It "test: add unit tests for CAI calculator per dimension" $null
Commit-It "test: add unit tests for CPS formula" $null
Commit-It "test: add unit tests for CBOM diff engine" $null
Commit-It "test: add integration test for full Phase 2 pipeline (scan → persist → compliance → simulate)" $null

Commit-It "docs: update README with Phase 2 architecture and capability overview" @("SPECTRA_FullBuildPrompt_Phase2.md")
Commit-It "docs: add COMPLIANCE.md documenting all embedded framework coverage with citations" $null
Commit-It "docs: add ARCHITECTURE.md with relationship graph entity model diagram" $null

git add .
git commit -m "ci: extend GitHub Actions to run integration tests with temp SQLite state"
