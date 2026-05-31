package persistence

import (
	"fmt"
)

// GraphNode represents a node in the cryptographic relationship graph.
type GraphNode struct {
	ID         string
	ScanID     string
	NodeType   string
	Label      string
	Properties string // JSON blob
}

// GraphEdge represents a directed edge in the graph.
type GraphEdge struct {
	ID       string
	ScanID   string
	FromNode string
	ToNode   string
	EdgeType string
	Weight   float64
}

// SaveGraph saves nodes and edges to the database.
func (s *Store) SaveGraph(scanID string, nodes []GraphNode, edges []GraphEdge) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	nodeQuery := `INSERT INTO graph_nodes (id, scan_id, node_type, label, properties) VALUES (?, ?, ?, ?, ?)`
	nodeStmt, err := tx.Prepare(nodeQuery)
	if err != nil {
		return fmt.Errorf("preparing node statement: %w", err)
	}
	defer nodeStmt.Close()

	for _, n := range nodes {
		_, err := nodeStmt.Exec(n.ID, scanID, n.NodeType, n.Label, n.Properties)
		if err != nil {
			return fmt.Errorf("inserting node %s: %w", n.ID, err)
		}
	}

	edgeQuery := `INSERT INTO graph_edges (id, scan_id, from_node, to_node, edge_type, weight) VALUES (?, ?, ?, ?, ?, ?)`
	edgeStmt, err := tx.Prepare(edgeQuery)
	if err != nil {
		return fmt.Errorf("preparing edge statement: %w", err)
	}
	defer edgeStmt.Close()

	for _, e := range edges {
		_, err := edgeStmt.Exec(e.ID, scanID, e.FromNode, e.ToNode, e.EdgeType, e.Weight)
		if err != nil {
			return fmt.Errorf("inserting edge %s: %w", e.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing graph: %w", err)
	}
	return nil
}

// GetGraph retrieves all nodes and edges for a given scan ID.
func (s *Store) GetGraph(scanID string) ([]GraphNode, []GraphEdge, error) {
	nodeQuery := `SELECT id, node_type, label, properties FROM graph_nodes WHERE scan_id = ?`
	nodeRows, err := s.db.Query(nodeQuery, scanID)
	if err != nil {
		return nil, nil, fmt.Errorf("querying nodes: %w", err)
	}
	defer nodeRows.Close()

	var nodes []GraphNode
	for nodeRows.Next() {
		var n GraphNode
		if err := nodeRows.Scan(&n.ID, &n.NodeType, &n.Label, &n.Properties); err != nil {
			return nil, nil, err
		}
		n.ScanID = scanID
		nodes = append(nodes, n)
	}

	edgeQuery := `SELECT id, from_node, to_node, edge_type, weight FROM graph_edges WHERE scan_id = ?`
	edgeRows, err := s.db.Query(edgeQuery, scanID)
	if err != nil {
		return nil, nil, fmt.Errorf("querying edges: %w", err)
	}
	defer edgeRows.Close()

	var edges []GraphEdge
	for edgeRows.Next() {
		var e GraphEdge
		if err := edgeRows.Scan(&e.ID, &e.FromNode, &e.ToNode, &e.EdgeType, &e.Weight); err != nil {
			return nil, nil, err
		}
		e.ScanID = scanID
		edges = append(edges, e)
	}

	return nodes, edges, nil
}
