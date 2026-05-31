package persistence

import (
	"fmt"
)

// ComplianceGap represents a compliance violation.
type ComplianceGap struct {
	ID              string
	FindingID       string
	Framework       string
	RequirementID   string
	RequirementDesc string
	Severity        string
	Deadline        *string // RFC3339 or nil
	DeadlineLabel   *string
}

// SaveComplianceGaps saves a slice of gaps to the database.
func (s *Store) SaveComplianceGaps(gaps []ComplianceGap) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
INSERT INTO compliance_gaps (
	id, finding_id, framework, requirement_id,
	requirement_desc, severity, deadline, deadline_label
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing compliance statement: %w", err)
	}
	defer stmt.Close()

	for _, g := range gaps {
		_, err := stmt.Exec(
			g.ID, g.FindingID, g.Framework, g.RequirementID,
			g.RequirementDesc, g.Severity, g.Deadline, g.DeadlineLabel,
		)
		if err != nil {
			return fmt.Errorf("inserting compliance gap %s: %w", g.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing compliance gaps: %w", err)
	}
	return nil
}
