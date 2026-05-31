package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store represents a connection to the SQLite database.
type Store struct {
	db *sql.DB
}

// NewStore initializes a connection to the SQLite database.
func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting user home dir: %w", err)
		}
		dbPath = filepath.Join(homeDir, ".spectra", "state.db")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("creating state directory: %w", err)
	}

	// Use PRAGMA foreign_keys = ON to enforce data integrity.
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return store, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) migrate() error {
	schema := `
-- SCANS
CREATE TABLE IF NOT EXISTS scans (
    id           TEXT PRIMARY KEY,
    scan_root    TEXT NOT NULL,
    started_at   TEXT NOT NULL,
    completed_at TEXT NOT NULL,
    go_version   TEXT,
    spectra_version TEXT,
    total_files  INTEGER,
    scanned_files INTEGER,
    aggregate_qrs INTEGER,
    cps          INTEGER,
    cai          INTEGER
);

-- FINDINGS
CREATE TABLE IF NOT EXISTS findings (
    id               TEXT PRIMARY KEY,
    scan_id          TEXT NOT NULL REFERENCES scans(id),
    algorithm        TEXT NOT NULL,
    source           TEXT NOT NULL,
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
    container_layer  TEXT,
    endpoint_host    TEXT,
    introduced_commit TEXT,
    introduced_author TEXT,
    introduced_date  TEXT,
    introduced_message TEXT,
    team_name        TEXT,
    finding_hash     TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_findings_scan_id   ON findings(scan_id);
CREATE INDEX IF NOT EXISTS idx_findings_algorithm ON findings(algorithm);
CREATE INDEX IF NOT EXISTS idx_findings_risk_band ON findings(risk_band);
CREATE INDEX IF NOT EXISTS idx_findings_hash      ON findings(finding_hash);

-- GRAPH NODES
CREATE TABLE IF NOT EXISTS graph_nodes (
    id          TEXT PRIMARY KEY,
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    node_type   TEXT NOT NULL,
    label       TEXT NOT NULL,
    properties  TEXT NOT NULL
);

-- GRAPH EDGES
CREATE TABLE IF NOT EXISTS graph_edges (
    id          TEXT PRIMARY KEY,
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    from_node   TEXT NOT NULL REFERENCES graph_nodes(id),
    to_node     TEXT NOT NULL REFERENCES graph_nodes(id),
    edge_type   TEXT NOT NULL,
    weight      REAL DEFAULT 1.0
);
CREATE INDEX IF NOT EXISTS idx_edges_from ON graph_edges(from_node);
CREATE INDEX IF NOT EXISTS idx_edges_to   ON graph_edges(to_node);

-- COMPLIANCE GAPS
CREATE TABLE IF NOT EXISTS compliance_gaps (
    id              TEXT PRIMARY KEY,
    finding_id      TEXT NOT NULL REFERENCES findings(id),
    framework       TEXT NOT NULL,
    requirement_id  TEXT NOT NULL,
    requirement_desc TEXT NOT NULL,
    severity        TEXT NOT NULL,
    deadline        TEXT,
    deadline_label  TEXT
);

-- BASELINES
CREATE TABLE IF NOT EXISTS baselines (
    name        TEXT PRIMARY KEY,
    scan_id     TEXT NOT NULL REFERENCES scans(id),
    created_at  TEXT NOT NULL,
    description TEXT
);
`
	_, err := s.db.Exec(schema)
	return err
}
