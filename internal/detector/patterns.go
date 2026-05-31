// Package detector provides the cryptographic algorithm registry, quantum risk
// scoring engine, migration effort classifier, and priority action plan builder.
package detector

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// PatternEntry holds compiled regex patterns for a specific algorithm and language.
type PatternEntry struct {
	Algorithm string
	Language  string
	Patterns  []*regexp.Regexp
}

// PatternRegistry holds all compiled patterns organized by language.
type PatternRegistry struct {
	// ByLanguage maps language -> slice of PatternEntry
	ByLanguage map[string][]PatternEntry
	// Generic patterns applied to all files
	Generic []PatternEntry
}

// LoadPatterns reads the crypto patterns YAML from the given file path,
// compiles every regex, and returns a PatternRegistry ready for scanning.
func LoadPatterns(yamlPath string) (*PatternRegistry, error) {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("reading patterns file %s: %w", yamlPath, err)
	}
	return LoadPatternsFromBytes(data)
}

// LoadPatternsFromBytes parses crypto patterns from raw YAML bytes, compiles
// every regex, and returns a PatternRegistry. This is useful for testing
// without touching the filesystem.
func LoadPatternsFromBytes(data []byte) (*PatternRegistry, error) {
	// Top-level: algorithm name -> (language -> []pattern-string)
	var raw map[string]map[string][]string
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing patterns YAML: %w", err)
	}

	registry := &PatternRegistry{
		ByLanguage: make(map[string][]PatternEntry),
	}

	for algo, langPatterns := range raw {
		for lang, patterns := range langPatterns {
			entry := PatternEntry{
				Algorithm: algo,
				Language:  lang,
				Patterns:  make([]*regexp.Regexp, 0, len(patterns)),
			}
			for _, p := range patterns {
				re, err := regexp.Compile(p)
				if err != nil {
					return nil, fmt.Errorf("compiling regex %q for %s/%s: %w", p, algo, lang, err)
				}
				entry.Patterns = append(entry.Patterns, re)
			}

			if lang == "generic" {
				registry.Generic = append(registry.Generic, entry)
			} else {
				registry.ByLanguage[lang] = append(registry.ByLanguage[lang], entry)
			}
		}
	}

	return registry, nil
}
