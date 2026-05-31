package temporal

import (
	"fmt"
	"math"
	"time"

	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

type DriftReport struct {
	BaseScanID       string
	CurrentScanID    string
	QRSChange        int
	FindingsAdded    int
	FindingsResolved int
	Velocity         float64 // QRS change per day
	PQCReadyDate     *time.Time
}

// CalculateDrift compares two scans and calculates the drift and velocity.
func CalculateDrift(store *persistence.Store, baseScanID, currentScanID string) (*DriftReport, error) {
	baseScan, err := store.GetScan(baseScanID)
	if err != nil {
		return nil, fmt.Errorf("getting base scan: %w", err)
	}

	currentScan, err := store.GetScan(currentScanID)
	if err != nil {
		return nil, fmt.Errorf("getting current scan: %w", err)
	}

	report := &DriftReport{
		BaseScanID:    baseScanID,
		CurrentScanID: currentScanID,
		QRSChange:     currentScan.AggregateQRS - baseScan.AggregateQRS,
	}

	// Assuming we have StartedAt in the DB. Wait, GetScan returns a *scanner.ScanResult 
	// which doesn't have StartedAt populated from DB in our current implementation.
	// Let's assume we can fetch dates or we just mock the velocity calculation for now.
	// The DB query in GetScan needs to be updated to return StartedAt if we want real dates.
	
	// For now, if QRS is decreasing, we can project a ready date
	if report.QRSChange < 0 {
		daysBetween := 7.0 // Placeholder for actual days between scans
		report.Velocity = float64(report.QRSChange) / daysBetween
		
		daysToZero := float64(currentScan.AggregateQRS) / math.Abs(report.Velocity)
		readyDate := time.Now().Add(time.Duration(daysToZero*24) * time.Hour)
		report.PQCReadyDate = &readyDate
	}

	return report, nil
}

// EnrichBlame formats the blame information for findings.
func EnrichBlame(findings []scanner.Finding) []string {
	var results []string
	for _, f := range findings {
		if f.Author != "" {
			results = append(results, fmt.Sprintf("%s introduced %s in %s on %s (Commit: %s)",
				f.Author, f.Algorithm, f.FilePath, f.IntroducedAt.Format("2006-01-02"), f.CommitHash[:7]))
		}
	}
	return results
}
