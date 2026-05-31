package persistence

import (
	"crypto/sha256"
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// SaveFindings saves a slice of Findings into the findings table.
func (s *Store) SaveFindings(scanID string, findings []scanner.Finding) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
INSERT INTO findings (
	id, scan_id, algorithm, source, file_path, line_number,
	line_content, language, key_size, context, qrs, risk_band,
	effort, effort_rationale, finding_hash
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing finding statement: %w", err)
	}
	defer stmt.Close()

	for _, f := range findings {
		// Calculate a deterministic hash for deduplication and tracking
		hashStr := fmt.Sprintf("%s|%s|%d", f.Algorithm, f.FilePath, f.LineNumber)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(hashStr)))

		_, err := stmt.Exec(
			f.ID, scanID, f.Algorithm, f.Source, f.FilePath, f.LineNumber,
			f.LineContent, f.Language, f.KeySize, f.Context, f.QRS, f.RiskBand,
			f.MigrationEffort, f.EffortRationale, hash,
		)
		if err != nil {
			return fmt.Errorf("inserting finding %s: %w", f.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing findings: %w", err)
	}
	return nil
}

// GetFindings retrieves all findings for a scan.
func (s *Store) GetFindings(scanID string) ([]scanner.Finding, error) {
	query := `
SELECT
	id, algorithm, source, file_path, line_number, line_content,
	language, key_size, context, qrs, risk_band, effort, effort_rationale,
	introduced_author, introduced_commit
FROM findings WHERE scan_id = ?
`
	rows, err := s.db.Query(query, scanID)
	if err != nil {
		return nil, fmt.Errorf("querying findings: %w", err)
	}
	defer rows.Close()

	var findings []scanner.Finding
	for rows.Next() {
		var f scanner.Finding
		var author, commit *string
		if err := rows.Scan(
			&f.ID, &f.Algorithm, &f.Source, &f.FilePath, &f.LineNumber, &f.LineContent,
			&f.Language, &f.KeySize, &f.Context, &f.QRS, &f.RiskBand, &f.MigrationEffort, &f.EffortRationale,
			&author, &commit,
		); err != nil {
			return nil, err
		}
		if author != nil {
			f.Author = *author
		}
		if commit != nil {
			f.CommitHash = *commit
		}
		findings = append(findings, f)
	}

	return findings, nil
}
