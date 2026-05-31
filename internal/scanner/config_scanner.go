// Package scanner implements the config file scanner.
package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/google/uuid"
)

// ConfigPattern defines a regex to extract algorithm names from config files.
type ConfigPattern struct {
	Regex *regexp.Regexp
	Group int
}

// Compile config patterns from the spec §7.5.
var configPatterns = []ConfigPattern{
	{Regex: regexp.MustCompile(`(?i)algorithm["\s:=]+["']?(\w[\w\-]*)`), Group: 1},
	{Regex: regexp.MustCompile(`(?i)cipher["\s:=]+["']?(\w[\w\-]*)`), Group: 1},
	{Regex: regexp.MustCompile(`(?i)hash[_\-]?algorithm["\s:=]+["']?(\w[\w\-]*)`), Group: 1},
	{Regex: regexp.MustCompile(`(?i)signature[_\-]?algorithm["\s:=]+["']?(\w[\w\-]*)`), Group: 1},
	// key size and tls version skipped since we only care about algorithms for now,
	// or we can detect known algo names:
	{Regex: regexp.MustCompile(`(?i)(SHA[-_]?1|MD[-_]?5|DES|3DES|RC4|RSA[-_]?\d*|ECDSA)\b`), Group: 1},
}

// targetConfigExtensions defines which file extensions we scan for configurations.
var targetConfigExtensions = map[string]bool{
	".yaml":       true,
	".yml":        true,
	".json":       true,
	".toml":       true,
	".env":        true,
	".ini":        true,
	".conf":       true,
	".properties": true,
	".xml":        true,
	".cfg":        true,
}

// ScanConfigFiles walks the directory and parses configuration files to detect
// cryptographic algorithms referenced by name.
func ScanConfigFiles(root string, excludes []string, patternRegistry *detector.PatternRegistry) ([]Finding, error) {
	var findings []Finding

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !targetConfigExtensions[ext] {
			// Also check for files like .env which might not have an extension if we consider .env the name
			name := d.Name()
			if !strings.HasPrefix(name, ".env") && !targetConfigExtensions[name] {
				return nil
			}
		}

		f, err := scanConfigFile(path)
		if err == nil {
			findings = append(findings, f...)
		}

		return nil
	})

	return findings, err
}

func scanConfigFile(path string) ([]Finding, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []Finding
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		for _, cp := range configPatterns {
			matches := cp.Regex.FindStringSubmatch(line)
			if len(matches) > cp.Group {
				candidate := matches[cp.Group]

				// Normalize and look up candidate
				algoInfo, found := detector.LookupAlgorithm(candidate)
				if found {
					f := Finding{
						ID:              uuid.NewString(),
						Algorithm:       algoInfo.Name,
						AlgorithmInfo:   algoInfo,
						Source:          detector.SourceConfig,
						FilePath:        path,
						LineNumber:      lineNum,
						LineContent:     strings.TrimSpace(line),
						Language:        "config",
						Context:         "Config Reference",
						OccurrenceCount: 1,
					}
					findings = append(findings, f)
					break // Break to avoid multiple findings on the same line if one matches
				}
			}
		}
		lineNum++
	}

	return findings, scanner.Err()
}
