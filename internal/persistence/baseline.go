package persistence

import (
	"fmt"
	"time"
)

// Baseline represents a saved snapshot of a scan.
type Baseline struct {
	Name        string
	ScanID      string
	CreatedAt   time.Time
	Description string
}

// SaveBaseline saves a baseline tag for a scan.
func (s *Store) SaveBaseline(b Baseline) error {
	query := `INSERT INTO baselines (name, scan_id, created_at, description) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, b.Name, b.ScanID, b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), b.Description)
	if err != nil {
		return fmt.Errorf("inserting baseline: %w", err)
	}
	return nil
}

// GetBaseline retrieves a baseline by name.
func (s *Store) GetBaseline(name string) (*Baseline, error) {
	query := `SELECT scan_id, created_at, description FROM baselines WHERE name = ?`
	var scanID, createdAtStr, desc string

	err := s.db.QueryRow(query, name).Scan(&scanID, &createdAtStr, &desc)
	if err != nil {
		return nil, fmt.Errorf("querying baseline: %w", err)
	}

	createdAt, _ := time.Parse("2006-01-02T15:04:05Z07:00", createdAtStr)
	return &Baseline{
		Name:        name,
		ScanID:      scanID,
		CreatedAt:   createdAt,
		Description: desc,
	}, nil
}
