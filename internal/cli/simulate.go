package cli

import (
	"encoding/json"
	"fmt"

	"github.com/HarshalPatel1972/spectra/internal/graph"
	"github.com/HarshalPatel1972/spectra/internal/persistence"
	"github.com/HarshalPatel1972/spectra/internal/simulation"
	"github.com/spf13/cobra"
)

var (
	simFrom   string
	simTo     string
	simScanID string
	simFormat string
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate a migration from one algorithm to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		if simFrom == "" || simTo == "" {
			return fmt.Errorf("--from and --to are required")
		}
		if simScanID == "" {
			return fmt.Errorf("--scan-id is required for now")
		}

		store, err := persistence.NewStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		nodes, edges, err := store.GetGraph(simScanID)
		if err != nil {
			return err
		}

		g := graph.NewGraph()
		for _, n := range nodes {
			props := map[string]any{}
			// In a real app we'd unmarshal n.Properties JSON string
			// For simplicity we extract algorithm if available
			if n.NodeType == string(graph.NodeAlgorithm) {
				props["algorithm"] = n.Label
			}
			g.AddNode(&graph.Node{
				ID:         n.ID,
				Type:       graph.NodeType(n.NodeType),
				Label:      n.Label,
				Properties: props,
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

		intent := simulation.MigrationIntent{
			From:   simFrom,
			To:     simTo,
			ScanID: simScanID,
		}

		result := simulation.SimulateMigration(g, intent)

		if simFormat == "json" {
			b, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		fmt.Printf("Migration Simulation: %s -> %s\n", simFrom, simTo)
		fmt.Printf("Total nodes affected: %d\n", result.TotalNodes)
		fmt.Printf("Estimated timeframe:  %d weeks\n", result.EstimatedWeeks)
		
		if len(result.BreakingChanges) > 0 {
			fmt.Println("\nWARNING: Breaking Changes Detected:")
			for _, bc := range result.BreakingChanges {
				fmt.Printf("  - %s (%s): %s\n", bc.NodeLabel, bc.NodeType, bc.Reason)
				fmt.Printf("    Workaround: %s\n", bc.Workaround)
			}
		}

		fmt.Println("\nMigration Waves:")
		for _, w := range result.Waves {
			fmt.Printf("  Wave %d: %s (%d nodes, ~%d days)\n", w.Number, w.Description, len(w.Nodes), w.EstimatedDays)
			for _, n := range w.Nodes {
				fmt.Printf("    - %s: %s\n", n.Type, n.Label)
			}
		}

		return nil
	},
}

func init() {
	simulateCmd.Flags().StringVar(&simFrom, "from", "", "source algorithm (required)")
	simulateCmd.Flags().StringVar(&simTo, "to", "", "target algorithm (required)")
	simulateCmd.Flags().StringVar(&simScanID, "scan-id", "", "which scan to base the simulation on (required)")
	simulateCmd.Flags().StringVar(&simFormat, "format", "terminal", "output format: terminal|json")
	rootCmd.AddCommand(simulateCmd)
}
