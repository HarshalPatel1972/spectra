// Package cbom generates CycloneDX Cryptographic Bill of Materials (CBOM)
// documents from Spectra scan results.
package cbom

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// cbomOutputFilename is the default file name for the CBOM export.
const cbomOutputFilename = "spectra-cbom.json"

// componentKey is a deduplication key for CBOM components.
type componentKey struct {
	Algorithm string
	KeySize   int
}

// componentAccum accumulates data for a single CBOM component across multiple
// findings.
type componentAccum struct {
	Info        detector.AlgorithmInfo
	KeySize     int
	MaxQRS      int
	RiskBand    detector.RiskBand
	Effort      detector.EffortLevel
	Occurrences int
}

// GenerateCBOM builds a CycloneDX BOM from the scan result. Each unique
// (Algorithm, KeySize) pair is mapped to a single cdx.Component with
// Spectra-specific properties attached.
func GenerateCBOM(result *scanner.ScanResult, version string) (*cdx.BOM, error) {
	// Deduplicate findings by (Algorithm, KeySize).
	accum := make(map[componentKey]*componentAccum)
	for _, f := range result.Findings {
		key := componentKey{Algorithm: f.Algorithm, KeySize: f.KeySize}
		if a, ok := accum[key]; ok {
			a.Occurrences++
			if f.QRS > a.MaxQRS {
				a.MaxQRS = f.QRS
				a.RiskBand = f.RiskBand
				a.Effort = f.MigrationEffort
			}
		} else {
			accum[key] = &componentAccum{
				Info:        f.AlgorithmInfo,
				KeySize:     f.KeySize,
				MaxQRS:      f.QRS,
				RiskBand:    f.RiskBand,
				Effort:      f.MigrationEffort,
				Occurrences: 1,
			}
		}
	}

	// Sort keys for deterministic output.
	keys := make([]componentKey, 0, len(accum))
	for k := range accum {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Algorithm != keys[j].Algorithm {
			return keys[i].Algorithm < keys[j].Algorithm
		}
		return keys[i].KeySize < keys[j].KeySize
	})

	// Build components.
	components := make([]cdx.Component, 0, len(keys))
	for _, key := range keys {
		a := accum[key]

		compVersion := ""
		if a.KeySize > 0 {
			compVersion = strconv.Itoa(a.KeySize)
		}

		description := a.Info.DisplayName
		if description == "" {
			description = key.Algorithm
		}

		props := &[]cdx.Property{
			{Name: "spectra:qrs", Value: strconv.Itoa(a.MaxQRS)},
			{Name: "spectra:risk-band", Value: string(a.RiskBand)},
			{Name: "spectra:effort", Value: string(a.Effort)},
			{Name: "spectra:occurrences", Value: strconv.Itoa(a.Occurrences)},
			{Name: "spectra:family", Value: a.Info.Family},
			{Name: "spectra:pqc-safe", Value: strconv.FormatBool(a.Info.PQCSafe)},
		}
		if a.Info.CWE != "" {
			*props = append(*props, cdx.Property{Name: "spectra:cwe", Value: a.Info.CWE})
		}

		comp := cdx.Component{
			Type:        cdx.ComponentTypeLibrary,
			Name:        key.Algorithm,
			Version:     compVersion,
			Description: description,
			Properties:  props,
		}
		components = append(components, comp)
	}

	serialNumber := "urn:uuid:" + uuid.New().String()

	bom := &cdx.BOM{
		BOMFormat:    "CycloneDX",
		SpecVersion:  cdx.SpecVersion1_6,
		Version:      1,
		SerialNumber: serialNumber,
		Metadata: &cdx.Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Tools: &cdx.ToolsChoice{
				Components: &[]cdx.Component{
					{
						Type:    cdx.ComponentTypeApplication,
						Name:    "spectra",
						Version: version,
					},
				},
			},
		},
		Components: &components,
	}

	return bom, nil
}

// WriteCBOM generates a CycloneDX CBOM from the scan result and writes it as
// JSON to outDir. The output directory is created if it does not already exist.
// Returns the path of the written file.
func WriteCBOM(result *scanner.ScanResult, outDir string, version string) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory %s: %w", outDir, err)
	}

	bom, err := GenerateCBOM(result, version)
	if err != nil {
		return "", fmt.Errorf("generating CBOM: %w", err)
	}

	outPath := filepath.Join(outDir, cbomOutputFilename)
	f, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("creating CBOM file %s: %w", outPath, err)
	}
	defer f.Close()

	encoder := cdx.NewBOMEncoder(f, cdx.BOMFileFormatJSON)
	encoder.SetPretty(true)
	if err := encoder.Encode(bom); err != nil {
		return "", fmt.Errorf("encoding CBOM: %w", err)
	}

	return outPath, nil
}
