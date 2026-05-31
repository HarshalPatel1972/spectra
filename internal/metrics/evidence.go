package metrics

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed evidence.yaml
var evidenceYAML []byte

type EvidenceItem struct {
	URL         string `yaml:"url"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type evidenceConfig struct {
	Evidence map[string][]EvidenceItem `yaml:"evidence"`
}

var eConfig evidenceConfig

func init() {
	_ = yaml.Unmarshal(evidenceYAML, &eConfig)
}

// GetEvidence returns a list of evidence items for a given algorithm.
func GetEvidence(algorithm string) []EvidenceItem {
	// Simple matching, could be extended to prefix matching like in compliance
	items, ok := eConfig.Evidence[algorithm]
	if !ok {
		return nil
	}
	return items
}
