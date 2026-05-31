package cli

import (
	"fmt"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/graph"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/spf13/cobra"
)

var (
	graphFormat string
	graphScanID string
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Export the cryptographic relationship graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		// If no scan ID is provided, we should ideally fetch the latest.
		// For now, if empty, we just return an error asking for it, or
		// implement a GetLatestScanID() in store.
		if graphScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}

		// Retrieve graph from DB
		nodes, edges, err := store.GetGraph(graphScanID)
		if err != nil {
			return err
		}

		g := graph.NewGraph()
		for _, n := range nodes {
			g.AddNode(&graph.Node{
				ID:    n.ID,
				Type:  graph.NodeType(n.NodeType),
				Label: n.Label,
			})
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

		switch strings.ToLower(graphFormat) {
		case "dot":
			fmt.Print(graph.ExportDOT(g))
		case "mermaid":
			fmt.Print(graph.ExportMermaid(g))
		case "json-ld":
			j, err := graph.ExportJSONLD(g)
			if err != nil {
				return err
			}
			fmt.Println(j)
		default:
			return fmt.Errorf("unknown format %q (use dot, mermaid, json-ld)", graphFormat)
		}

		return nil
	},
}

func init() {
	graphCmd.Flags().StringVar(&graphFormat, "format", "mermaid", "output format: dot, mermaid, json-ld")
	graphCmd.Flags().StringVar(&graphScanID, "scan-id", "", "which scan to graph (required)")
	rootCmd.AddCommand(graphCmd)
}
