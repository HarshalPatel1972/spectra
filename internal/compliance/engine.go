package compliance

import (
	_ "embed"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

//go:embed compliance_rules.yaml
var complianceRulesYAML []byte

type RuleCondition struct {
	Source     string   `yaml:"source"`
	Context    string   `yaml:"context"`
	Algorithms []string `yaml:"algorithms"`
}

type Rule struct {
	Framework       string        `yaml:"framework"`
	RequirementID   string        `yaml:"requirement_id"`
	RequirementDesc string        `yaml:"requirement_desc"`
	Severity        string        `yaml:"severity"`
	Deadline        string        `yaml:"deadline"`
	DeadlineLabel   string        `yaml:"deadline_label"`
	Conditions      RuleCondition `yaml:"conditions"`
}

type ruleConfig struct {
	Rules []Rule `yaml:"rules"`
}

var config ruleConfig

func init() {
	_ = yaml.Unmarshal(complianceRulesYAML, &config)
}

// EvaluateFindings evaluates findings against compliance rules and generates gaps.
func EvaluateFindings(findings []scanner.Finding) []persistence.ComplianceGap {
	var gaps []persistence.ComplianceGap

	for _, f := range findings {
		for _, r := range config.Rules {
			if matchesCondition(f, r.Conditions) {
				var deadline *string
				if r.Deadline != "" {
					d := r.Deadline
					deadline = &d
				}
				var deadlineLabel *string
				if r.DeadlineLabel != "" {
					dl := r.DeadlineLabel
					deadlineLabel = &dl
				}
				
				gaps = append(gaps, persistence.ComplianceGap{
					ID:              uuid.New().String(),
					FindingID:       f.ID,
					Framework:       r.Framework,
					RequirementID:   r.RequirementID,
					RequirementDesc: r.RequirementDesc,
					Severity:        r.Severity,
					Deadline:        deadline,
					DeadlineLabel:   deadlineLabel,
				})
			}
		}
	}

	return gaps
}

func matchesCondition(f scanner.Finding, cond RuleCondition) bool {
	if cond.Source != "" && string(f.Source) != cond.Source {
		return false
	}
	if cond.Context != "" && f.Context != cond.Context {
		return false
	}
	
	if len(cond.Algorithms) > 0 {
		algoMatched := false
		for _, a := range cond.Algorithms {
			if strings.HasPrefix(f.Algorithm, a) { // Handle AES matching AES128
				algoMatched = true
				break
			}
		}
		if !algoMatched {
			return false
		}
	}

	return true
}
