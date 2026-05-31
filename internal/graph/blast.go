package graph

// BlastRadius calculates the set of all nodes reachable from algorithmID
// via any path (following REVERSE edges to find dependents).
// It returns a map of node ID to minimum hop count.
func BlastRadius(g *Graph, algorithmID string) map[string]int {
	visited := map[string]int{algorithmID: 0}
	queue := []string{algorithmID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range g.InEdges[current] {
			if _, seen := visited[edge.From]; !seen {
				visited[edge.From] = visited[current] + 1
				queue = append(queue, edge.From)
			}
		}
	}
	return visited
}
