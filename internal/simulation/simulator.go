package simulation

import (
	"fmt"
	_ "embed"

	"github.com/HarshalPatel1972/spectra/internal/graph"
	"gopkg.in/yaml.v3"
)

//go:embed compatibility_matrix.yaml
var compatibilityYAML []byte

type MigrationIntent struct {
	From   string
	To     string
	ScanID string
}

type SimulationResult struct {
	Intent               MigrationIntent
	Waves                []MigrationWave
	BreakingChanges      []BreakingChange
	IncompatibleContexts []ContextIncompatibility
	TotalNodes           int
	EstimatedWeeks       int
	QRSImpact            int
}

type MigrationWave struct {
	Number        int
	Description   string
	Nodes         []WaveNode
	EstimatedDays int
}

type WaveNode struct {
	ID    string
	Label string
	Type  string
}

type BreakingChange struct {
	NodeLabel  string
	NodeType   string
	Reason     string
	Workaround string
}

type ContextIncompatibility struct {
	Context string
}

type matrix struct {
	Substitutions map[string]map[string]any `yaml:"substitutions"`
}

var compatibilityMatrix matrix

func init() {
	_ = yaml.Unmarshal(compatibilityYAML, &compatibilityMatrix)
}

func checkCompatibility(from, to, context string) bool {
	key := fmt.Sprintf("%s -> %s", from, to)
	contexts, ok := compatibilityMatrix.Substitutions[key]
	if !ok {
		return false
	}
	if all, ok := contexts["all_contexts"]; ok && all == true {
		return true
	}
	val, ok := contexts[context]
	if !ok {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}

// SimulateMigration simulates a migration from one algorithm to another.
func SimulateMigration(g *graph.Graph, intent MigrationIntent) *SimulationResult {
	// A full topological sort is complex, we will create a simplified grouping here.
	result := &SimulationResult{
		Intent: intent,
	}

	// Find the source algorithm node(s)
	var fromNodes []*graph.Node
	for _, n := range g.Nodes {
		if n.Type == graph.NodeAlgorithm && n.Properties["algorithm"] == intent.From {
			fromNodes = append(fromNodes, n)
		}
	}

	if len(fromNodes) == 0 {
		return result
	}

	var allAffected []WaveNode
	
	// Collect all nodes connected to the source algorithms
	for _, fn := range fromNodes {
		radius := graph.BlastRadius(g, fn.ID)
		for nodeID, _ := range radius {
			if nodeID == fn.ID {
				continue
			}
			n := g.Nodes[nodeID]
			allAffected = append(allAffected, WaveNode{
				ID:    n.ID,
				Label: n.Label,
				Type:  string(n.Type),
			})

			// Check compatibility if it's a code finding
			if n.Type == graph.NodeFile {
				// Assuming Context is in properties, but for simplicity, we mock a check
				context := "digital_signature" 
				if propCtx, ok := n.Properties["context"]; ok && propCtx != "" {
					context = propCtx.(string)
				}
				
				if !checkCompatibility(intent.From, intent.To, context) {
					result.BreakingChanges = append(result.BreakingChanges, BreakingChange{
						NodeLabel:  n.Label,
						NodeType:   string(n.Type),
						Reason:     fmt.Sprintf("Incompatible in context: %s", context),
						Workaround: "Review algorithmic requirements or implement a hybrid scheme.",
					})
				}
			}
		}
	}

	// Group into a single wave for MVP simulation
	if len(allAffected) > 0 {
		result.Waves = append(result.Waves, MigrationWave{
			Number:        1,
			Description:   "Primary remediation wave",
			Nodes:         allAffected,
			EstimatedDays: len(allAffected) * 2,
		})
	}

	result.TotalNodes = len(allAffected)
	result.EstimatedWeeks = (len(allAffected) * 2) / 5
	if result.EstimatedWeeks == 0 && len(allAffected) > 0 {
		result.EstimatedWeeks = 1
	}

	return result
}
