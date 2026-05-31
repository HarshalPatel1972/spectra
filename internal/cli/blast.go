package cli

import (
	"fmt"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/graph"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/spf13/cobra"
)

var (
	blastAlgorithm string
	blastScanID    string
)

var blastCmd = &cobra.Command{
	Use:   "blast",
	Short: "Calculate the blast radius of replacing an algorithm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if blastAlgorithm == "" {
			return fmt.Errorf("--algorithm is required")
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		if blastScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}

		nodes, edges, err := store.GetGraph(blastScanID)
		if err != nil {
			return err
		}

		g := graph.NewGraph()
		var targetNodeID string

		for _, n := range nodes {
			g.AddNode(&graph.Node{
				ID:    n.ID,
				Type:  graph.NodeType(n.NodeType),
				Label: n.Label,
			})
			// Prefix match since algorithm could have -2048 attached
			if n.NodeType == string(graph.NodeAlgorithm) && strings.HasPrefix(n.Label, blastAlgorithm) {
				// We'll take the first match for simplicity, or we should handle multiple
				targetNodeID = n.ID
			}
		}

		if targetNodeID == "" {
			return fmt.Errorf("algorithm %q not found in graph", blastAlgorithm)
		}

		for _, e := range edges {
			g.AddEdge(&graph.Edge{
				ID:     e.ID,
				From:   e.FromNode,
				To:     e.ToNode,
				Type:   graph.EdgeType(e.EdgeType),
				Weight: e.Weight,
			})
		}

		radius := graph.BlastRadius(g, targetNodeID)

		fmt.Printf("Blast Radius for %s:\n", blastAlgorithm)
		fmt.Printf("Total affected nodes: %d\n", len(radius)-1) // Exclude the algorithm node itself

		for nodeID, hops := range radius {
			if nodeID == targetNodeID {
				continue
			}
			node := g.Nodes[nodeID]
			fmt.Printf("  [%d hops] %s: %s\n", hops, node.Type, node.Label)
		}

		return nil
	},
}

func init() {
	blastCmd.Flags().StringVar(&blastAlgorithm, "algorithm", "", "algorithm to analyze (required)")
	blastCmd.Flags().StringVar(&blastScanID, "scan-id", "", "which scan to use (required)")
	rootCmd.AddCommand(blastCmd)
}
