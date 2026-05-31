// Package report provides terminal, JSON, and HTML output renderers for
// Spectra scan results.
package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// ANSI escape codes for terminal colouring.
const (
	ColorReset    = "\033[0m"
	ColorBold     = "\033[1m"
	ColorDim      = "\033[2m"
	ColorCritical = "\033[38;2;255;68;68m"   // #FF4444
	ColorHigh     = "\033[38;2;255;140;0m"   // #FF8C00
	ColorMedium   = "\033[38;2;255;215;0m"   // #FFD700
	ColorLow      = "\033[38;2;74;222;128m"  // #4ADE80
	ColorSafe     = "\033[38;2;56;189;248m"  // #38BDF8
	ColorWhite    = "\033[38;2;248;250;252m" // #F8FAFC
)

// spectraASCII is the ASCII art wordmark printed at the top of every terminal
// report.
const spectraASCII = `
 ███████╗██████╗ ███████╗ ██████╗████████╗██████╗  █████╗
 ██╔════╝██╔══██╗██╔════╝██╔════╝╚══██╔══╝██╔══██╗██╔══██╗
 ███████╗██████╔╝█████╗  ██║        ██║   ██████╔╝███████║
 ╚════██║██╔═══╝ ██╔══╝  ██║        ██║   ██╔══██╗██╔══██║
 ███████║██║     ███████╗╚██████╗   ██║   ██║  ██║██║  ██║
 ╚══════╝╚═╝     ╚══════╝ ╚═════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝`

// maxFindingsDisplay is the maximum number of findings rows displayed in the
// terminal table before truncation.
const maxFindingsDisplay = 50

// bandColor returns the ANSI colour code for the given risk band.
func bandColor(band detector.RiskBand) string {
	switch band {
	case detector.BandCritical:
		return ColorCritical
	case detector.BandHigh:
		return ColorHigh
	case detector.BandMedium:
		return ColorMedium
	case detector.BandLow:
		return ColorLow
	case detector.BandSafe:
		return ColorSafe
	default:
		return ColorWhite
	}
}

// bandFromQRS derives a RiskBand from a numeric QRS value.
func bandFromQRS(qrs int) detector.RiskBand {
	switch {
	case qrs >= 80:
		return detector.BandCritical
	case qrs >= 60:
		return detector.BandHigh
	case qrs >= 40:
		return detector.BandMedium
	case qrs >= 20:
		return detector.BandLow
	default:
		return detector.BandSafe
	}
}

// truncatePath shortens a file path for display if it exceeds maxLen.
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

// RenderTerminal writes a full terminal report of scan results to the provided
// writer. The output includes an ASCII header, findings table, aggregate
// summary, and top action items.
func RenderTerminal(w io.Writer, result *scanner.ScanResult, version string) error {
	// ── Header ──────────────────────────────────────────────────────────
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ColorBold+ColorWhite, spectraASCII, ColorReset); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	if _, err := fmt.Fprintf(w, "\n  %sSpectra %s%s\n", ColorDim, version, ColorReset); err != nil {
		return fmt.Errorf("writing version: %w", err)
	}
	if _, err := fmt.Fprintf(w, "  %sScan root : %s%s\n", ColorDim, result.ScanRoot, ColorReset); err != nil {
		return fmt.Errorf("writing scan root: %w", err)
	}
	if _, err := fmt.Fprintf(w, "  %sCompleted : %s%s\n\n",
		ColorDim, result.CompletedAt.Format(time.RFC3339), ColorReset); err != nil {
		return fmt.Errorf("writing timestamp: %w", err)
	}

	// ── Summary line ────────────────────────────────────────────────────
	if _, err := fmt.Fprintf(w, "  %sFiles: %d scanned, %d skipped, %d total%s\n\n",
		ColorWhite, result.ScannedFiles, result.SkippedFiles, result.TotalFiles, ColorReset); err != nil {
		return fmt.Errorf("writing summary: %w", err)
	}

	// ── Findings table ──────────────────────────────────────────────────
	if len(result.Findings) == 0 {
		if _, err := fmt.Fprintf(w, "  %s✓ No cryptographic findings detected.%s\n\n",
			ColorSafe, ColorReset); err != nil {
			return fmt.Errorf("writing no-findings: %w", err)
		}
		return nil
	}

	// Sort by QRS descending.
	sorted := make([]scanner.Finding, len(result.Findings))
	copy(sorted, result.Findings)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].QRS > sorted[j].QRS
	})

	tw := tabwriter.NewWriter(w, 2, 4, 2, ' ', 0)
	header := fmt.Sprintf("  %s%sALGORITHM\tSOURCE\tFILE\tLINE\tKEY SIZE\tQRS\tBAND\tEFFORT%s",
		ColorBold, ColorWhite, ColorReset)
	if _, err := fmt.Fprintln(tw, header); err != nil {
		return fmt.Errorf("writing table header: %w", err)
	}

	separator := fmt.Sprintf("  %s─────────\t──────\t────\t────\t────────\t───\t────\t──────%s",
		ColorDim, ColorReset)
	if _, err := fmt.Fprintln(tw, separator); err != nil {
		return fmt.Errorf("writing table separator: %w", err)
	}

	displayCount := len(sorted)
	if displayCount > maxFindingsDisplay {
		displayCount = maxFindingsDisplay
	}

	for i := 0; i < displayCount; i++ {
		f := sorted[i]
		bc := bandColor(f.RiskBand)
		keySizeStr := "-"
		if f.KeySize > 0 {
			keySizeStr = fmt.Sprintf("%d", f.KeySize)
		}
		row := fmt.Sprintf("  %s\t%s\t%s\t%d\t%s\t%d\t%s%s%s\t%s",
			f.Algorithm,
			string(f.Source),
			truncatePath(f.FilePath, 40),
			f.LineNumber,
			keySizeStr,
			f.QRS,
			bc, string(f.RiskBand), ColorReset,
			string(f.MigrationEffort),
		)
		if _, err := fmt.Fprintln(tw, row); err != nil {
			return fmt.Errorf("writing finding row: %w", err)
		}
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flushing table: %w", err)
	}

	remaining := len(sorted) - displayCount
	if remaining > 0 {
		if _, err := fmt.Fprintf(w, "\n  %s... and %d more findings%s\n",
			ColorDim, remaining, ColorReset); err != nil {
			return fmt.Errorf("writing truncation notice: %w", err)
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return fmt.Errorf("writing newline: %w", err)
	}

	// ── Aggregate block ─────────────────────────────────────────────────
	aggBand := bandFromQRS(result.AggregateQRS)
	aggColor := bandColor(aggBand)
	if _, err := fmt.Fprintf(w, "  %s%sAggregate QRS: %d/100 (%s)%s\n",
		ColorBold, aggColor, result.AggregateQRS, string(aggBand), ColorReset); err != nil {
		return fmt.Errorf("writing aggregate QRS: %w", err)
	}

	bands := []detector.RiskBand{
		detector.BandCritical, detector.BandHigh, detector.BandMedium,
		detector.BandLow, detector.BandSafe,
	}
	var parts []string
	for _, b := range bands {
		count := result.FindingsByBand[b]
		bc := bandColor(b)
		parts = append(parts, fmt.Sprintf("%s%s: %d%s", bc, string(b), count, ColorReset))
	}
	if _, err := fmt.Fprintf(w, "  %s\n\n", strings.Join(parts, " | ")); err != nil {
		return fmt.Errorf("writing band counts: %w", err)
	}

	// ── Action items (top 5) ────────────────────────────────────────────
	if len(result.ActionPlan) > 0 {
		if _, err := fmt.Fprintf(w, "  %s%sTop Action Items%s\n", ColorBold, ColorWhite, ColorReset); err != nil {
			return fmt.Errorf("writing action header: %w", err)
		}
		if _, err := fmt.Fprintf(w, "  %s────────────────%s\n", ColorDim, ColorReset); err != nil {
			return fmt.Errorf("writing action separator: %w", err)
		}

		actionCount := len(result.ActionPlan)
		if actionCount > 5 {
			actionCount = 5
		}
		for i := 0; i < actionCount; i++ {
			a := result.ActionPlan[i]
			bc := bandColor(a.RiskBand)
			if _, err := fmt.Fprintf(w,
				"  %s%d.%s %s%s%s — %s (%d occurrences, effort: %s)\n     %s→ %s%s\n",
				ColorBold, a.Rank, ColorReset,
				bc, a.Algorithm, ColorReset,
				a.Recommendation,
				a.Occurrences,
				string(a.Effort),
				ColorDim, a.Replacement, ColorReset,
			); err != nil {
				return fmt.Errorf("writing action item %d: %w", i+1, err)
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return fmt.Errorf("writing trailing newline: %w", err)
		}
	}

	return nil
}
