package persistence

import (
	"fmt"
	"strings"
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
		dbID := scanID + "-" + n.ID
		_, err := nodeStmt.Exec(dbID, scanID, n.NodeType, n.Label, n.Properties)
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
		dbID := scanID + "-" + e.ID
		fromNode := scanID + "-" + e.FromNode
		toNode := scanID + "-" + e.ToNode
		_, err := edgeStmt.Exec(dbID, scanID, fromNode, toNode, e.EdgeType, e.Weight)
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
		var dbID string
		if err := nodeRows.Scan(&dbID, &n.NodeType, &n.Label, &n.Properties); err != nil {
			return nil, nil, err
		}
		n.ID = strings.TrimPrefix(dbID, scanID+"-")
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
		var dbID, dbFrom, dbTo string
		if err := edgeRows.Scan(&dbID, &dbFrom, &dbTo, &e.EdgeType, &e.Weight); err != nil {
			return nil, nil, err
		}
		e.ID = strings.TrimPrefix(dbID, scanID+"-")
		e.FromNode = strings.TrimPrefix(dbFrom, scanID+"-")
		e.ToNode = strings.TrimPrefix(dbTo, scanID+"-")
		e.ScanID = scanID
		edges = append(edges, e)
	}

	return nodes, edges, nil
}
