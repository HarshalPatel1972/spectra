package graph

import (
	"crypto/sha256"
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// BuildGraph constructs a relationship graph from a ScanResult.
func BuildGraph(result *scanner.ScanResult) *Graph {
	g := NewGraph()

	for _, f := range result.Findings {
		// 1. Create/Upsert ALGORITHM Node
		algoID := fmt.Sprintf("algo:%s:%d", f.Algorithm, f.KeySize)
		algoHash := fmt.Sprintf("%x", sha256.Sum256([]byte(algoID)))
		
		if _, exists := g.Nodes[algoHash]; !exists {
			label := f.Algorithm
			if f.KeySize > 0 {
				label = fmt.Sprintf("%s-%d", f.Algorithm, f.KeySize)
			}
			g.AddNode(&Node{
				ID:    algoHash,
				Type:  NodeAlgorithm,
				Label: label,
				Properties: map[string]any{
					"algorithm": f.Algorithm,
					"key_size":  f.KeySize,
					"qrs":       f.QRS,
					"risk_band": f.RiskBand,
				},
			})
		}

		// 2. Create/Upsert SOURCE Node
		sourceType := NodeFile
		if f.Source == "CERT" {
			sourceType = NodeCert
		} else if f.Source == "DEPS" {
			sourceType = NodeDependency
		} else if f.Source == "ENDPOINT" {
			sourceType = NodeEndpoint
		}

		sourceID := fmt.Sprintf("%s:%s", sourceType, f.FilePath)
		sourceHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sourceID)))

		if _, exists := g.Nodes[sourceHash]; !exists {
			g.AddNode(&Node{
				ID:    sourceHash,
				Type:  sourceType,
				Label: f.FilePath,
				Properties: map[string]any{
					"file_path": f.FilePath,
					"source":    f.Source,
				},
			})
		}

		// 3. Create USES (or PROVIDES) Edge
		edgeType := EdgeUses
		if sourceType == NodeDependency {
			edgeType = EdgeProvides
		} else if sourceType == NodeCert {
			edgeType = EdgeUses
		}

		edgeID := fmt.Sprintf("%s->%s:%s", sourceHash, algoHash, edgeType)
		edgeHash := fmt.Sprintf("%x", sha256.Sum256([]byte(edgeID)))

		// Avoid duplicate edges
		edgeExists := false
		for _, e := range g.OutEdges[sourceHash] {
			if e.To == algoHash && e.Type == edgeType {
				edgeExists = true
				break
			}
		}

		if !edgeExists {
			g.AddEdge(&Edge{
				ID:     edgeHash,
				From:   sourceHash,
				To:     algoHash,
				Type:   edgeType,
				Weight: 1.0,
			})
		}
	}

	return g
}
