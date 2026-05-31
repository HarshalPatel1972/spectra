// Package config implements the .spectra.yaml configuration loader, default
// values, and CLI flag merge logic for the Spectra cryptographic scanner.
package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure loaded from .spectra.yaml.
type Config struct {
	Scan   ScanConfig   `yaml:"scan"`
	Output OutputConfig `yaml:"output"`
	CI     CIConfig     `yaml:"ci"`
}

// ScanConfig controls which files and scanners are used during a scan.
type ScanConfig struct {
	Exclude     []string `yaml:"exclude"`
	Extensions  []string `yaml:"extensions"`
	Scanners    []string `yaml:"scanners"`
	Concurrency int      `yaml:"concurrency"`
}

// OutputConfig controls the output formats and directory for scan results.
type OutputConfig struct {
	Formats []string `yaml:"formats"`
	Dir     string   `yaml:"dir"`
}

// CIConfig holds CI/CD pipeline integration options.
type CIConfig struct {
	FailOn string `yaml:"fail_on"`
}

// DefaultConfig returns a Config populated with sensible defaults.
//
// Default values:
//   - Scanners: [code, cert, deps, config]
//   - Formats:  [terminal]
//   - Dir:      ./spectra-out
//   - Concurrency: 0 (resolved to runtime.NumCPU() at scan time)
func DefaultConfig() *Config {
	return &Config{
		Scan: ScanConfig{
			Scanners:    []string{"code", "cert", "deps", "config"},
			Concurrency: 0,
		},
		Output: OutputConfig{
			Formats: []string{"terminal"},
			Dir:     "./spectra-out",
		},
	}
}

// LoadConfig reads a YAML configuration file from path and returns the
// parsed Config. If the file does not exist or cannot be parsed, an error
// is returned.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}

// FlagValues holds CLI flag values that may override configuration. A zero
// value for any field means "not set by the user."
type FlagValues struct {
	OutputFormats []string
	OutDir        string
	Exclude       []string
	Scanners      []string
	FailOn        string
	Concurrency   int
}

// MergeFlags merges non-zero CLI flag values into cfg, with flags taking
// precedence over config file values.
func MergeFlags(cfg *Config, flags FlagValues) {
	if len(flags.OutputFormats) > 0 {
		cfg.Output.Formats = flags.OutputFormats
	}
	if flags.OutDir != "" {
		cfg.Output.Dir = flags.OutDir
	}
	if len(flags.Exclude) > 0 {
		cfg.Scan.Exclude = flags.Exclude
	}
	if len(flags.Scanners) > 0 {
		cfg.Scan.Scanners = flags.Scanners
	}
	if flags.FailOn != "" {
		cfg.CI.FailOn = flags.FailOn
	}
	if flags.Concurrency != 0 {
		cfg.Scan.Concurrency = flags.Concurrency
	}
}

// EffectiveConcurrency returns the concurrency value to use. If cfg.Scan.Concurrency
// is 0 (the default), it returns runtime.NumCPU().
func EffectiveConcurrency(cfg *Config) int {
	if cfg.Scan.Concurrency > 0 {
		return cfg.Scan.Concurrency
	}
	return runtime.NumCPU()
}
