package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
)

// DiffResult holds the results of comparing two CBOM files.
type DiffResult struct {
	OldComponents int
	NewComponents int
	Resolved      []cyclonedx.Component
	Added         []cyclonedx.Component
	NetChange     int
}

// CompareCBOMs reads two CycloneDX CBOM files and computes the differences.
func CompareCBOMs(oldPath, newPath string) (*DiffResult, error) {
	oldBOM, err := loadCBOM(oldPath)
	if err != nil {
		return nil, fmt.Errorf("loading old CBOM: %w", err)
	}

	newBOM, err := loadCBOM(newPath)
	if err != nil {
		return nil, fmt.Errorf("loading new CBOM: %w", err)
	}

	oldComps := make(map[string]cyclonedx.Component)
	if oldBOM.Components != nil {
		for _, c := range *oldBOM.Components {
			// Create a unique key for the component based on Name and Version (KeySize)
			key := fmt.Sprintf("%s-%s", c.Name, c.Version)
			oldComps[key] = c
		}
	}

	newComps := make(map[string]cyclonedx.Component)
	if newBOM.Components != nil {
		for _, c := range *newBOM.Components {
			key := fmt.Sprintf("%s-%s", c.Name, c.Version)
			newComps[key] = c
		}
	}

	res := &DiffResult{
		OldComponents: len(oldComps),
		NewComponents: len(newComps),
	}

	// Find added
	for key, c := range newComps {
		if _, exists := oldComps[key]; !exists {
			res.Added = append(res.Added, c)
		}
	}

	// Find resolved
	for key, c := range oldComps {
		if _, exists := newComps[key]; !exists {
			res.Resolved = append(res.Resolved, c)
		}
	}

	res.NetChange = res.NewComponents - res.OldComponents

	return res, nil
}

func loadCBOM(path string) (*cyclonedx.BOM, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bom cyclonedx.BOM
	if err := json.NewDecoder(f).Decode(&bom); err != nil {
		return nil, err
	}
	return &bom, nil
}

// RenderDiffTerminal prints the diff results to the terminal.
func RenderDiffTerminal(w io.Writer, diff *DiffResult) {
	fmt.Fprintf(w, "Spectra CBOM Diff\n")
	fmt.Fprintf(w, "=================\n\n")

	fmt.Fprintf(w, "Old Findings: %d\n", diff.OldComponents)
	fmt.Fprintf(w, "New Findings: %d\n", diff.NewComponents)
	fmt.Fprintf(w, "Net Change:   %+d\n\n", diff.NetChange)

	if len(diff.Resolved) > 0 {
		fmt.Fprintf(w, "Resolved Findings (Good job!):\n")
		for _, c := range diff.Resolved {
			fmt.Fprintf(w, "  - %s (Size: %s)\n", c.Name, c.Version)
		}
		fmt.Fprintln(w)
	}

	if len(diff.Added) > 0 {
		fmt.Fprintf(w, "Newly Introduced Findings (Warning!):\n")
		for _, c := range diff.Added {
			fmt.Fprintf(w, "  + %s (Size: %s)\n", c.Name, c.Version)
		}
		fmt.Fprintln(w)
	}
	
	if len(diff.Resolved) == 0 && len(diff.Added) == 0 {
		fmt.Fprintf(w, "No changes in cryptographic assets.\n")
	}
}
