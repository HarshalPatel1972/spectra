package persistence

import (
	"database/sql"
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// SaveScan saves a ScanResult into the scans table.
func (s *Store) SaveScan(scanID string, r *scanner.ScanResult) error {
	query := `
INSERT INTO scans (
	id, scan_root, started_at, completed_at, 
	total_files, scanned_files, aggregate_qrs, cps, cai
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	// Note: go_version and spectra_version can be added if available in ScanResult.
	_, err := s.db.Exec(query,
		scanID,
		r.ScanRoot,
		r.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		r.CompletedAt.Format("2006-01-02T15:04:05Z07:00"),
		r.TotalFiles,
		r.ScannedFiles,
		r.AggregateQRS,
		0, // cps placeholder
		0, // cai placeholder
	)
	if err != nil {
		return fmt.Errorf("inserting scan: %w", err)
	}
	return nil
}

// GetScan retrieves a ScanResult summary from the scans table.
func (s *Store) GetScan(scanID string) (*scanner.ScanResult, error) {
	query := `
SELECT scan_root, total_files, scanned_files, aggregate_qrs
FROM scans WHERE id = ?
`
	var root string
	var total, scanned, aggQRS int

	err := s.db.QueryRow(query, scanID).Scan(&root, &total, &scanned, &aggQRS)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scan not found")
	} else if err != nil {
		return nil, fmt.Errorf("querying scan: %w", err)
	}

	return &scanner.ScanResult{
		ScanRoot:     root,
		TotalFiles:   total,
		ScannedFiles: scanned,
		AggregateQRS: aggQRS,
	}, nil
}

// GetLatestScan retrieves the most recent ScanResult from the scans table.
func (s *Store) GetLatestScan() (*scanner.ScanResult, error) {
	query := `
SELECT scan_root, total_files, scanned_files, aggregate_qrs
FROM scans ORDER BY started_at DESC LIMIT 1
`
	var root string
	var total, scanned, aggQRS int

	err := s.db.QueryRow(query).Scan(&root, &total, &scanned, &aggQRS)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no scans found")
	} else if err != nil {
		return nil, fmt.Errorf("querying latest scan: %w", err)
	}

	return &scanner.ScanResult{
		ScanRoot:     root,
		TotalFiles:   total,
		ScannedFiles: scanned,
		AggregateQRS: aggQRS,
	}, nil
}
