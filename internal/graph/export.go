package graph

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExportDOT returns the graph in GraphViz DOT format.
func ExportDOT(g *Graph) string {
	var b strings.Builder
	b.WriteString("digraph spectra {\n")
	b.WriteString("    rankdir=LR;\n")

	for _, n := range g.Nodes {
		shape := "box"
		if n.Type == NodeAlgorithm {
			shape = "diamond"
		}
		b.WriteString(fmt.Sprintf("    \"%s\" [shape=%s, label=\"%s\"];\n", n.ID, shape, n.Label))
	}

	for _, e := range g.Edges {
		b.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"%s\"];\n", e.From, e.To, e.Type))
	}

	b.WriteString("}\n")
	return b.String()
}

// ExportMermaid returns the graph in Mermaid format.
func ExportMermaid(g *Graph) string {
	var b strings.Builder
	b.WriteString("graph LR\n")

	for _, n := range g.Nodes {
		if n.Type == NodeAlgorithm {
			b.WriteString(fmt.Sprintf("    %s((%s))\n", n.ID, n.Label))
		} else {
			b.WriteString(fmt.Sprintf("    %s[%s]\n", n.ID, n.Label))
		}
	}

	for _, e := range g.Edges {
		b.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", e.From, e.Type, e.To))
	}

	return b.String()
}

// ExportJSONLD returns the graph in JSON-LD format.
func ExportJSONLD(g *Graph) (string, error) {
	type JSONLDNode struct {
		ID   string `json:"@id"`
		Type string `json:"@type"`
		Name string `json:"name"`
	}
	type JSONLDEdge struct {
		Subject   string `json:"subject"`
		Predicate string `json:"predicate"`
		Object    string `json:"object"`
	}
	type JSONLDGraph struct {
		Context string       `json:"@context"`
		Graph   []any        `json:"@graph"`
	}

	doc := JSONLDGraph{
		Context: "https://schema.org/",
		Graph:   []any{},
	}

	for _, n := range g.Nodes {
		doc.Graph = append(doc.Graph, JSONLDNode{
			ID:   n.ID,
			Type: string(n.Type),
			Name: n.Label,
		})
	}

	for _, e := range g.Edges {
		doc.Graph = append(doc.Graph, JSONLDEdge{
			Subject:   e.From,
			Predicate: string(e.Type),
			Object:    e.To,
		})
	}

	bytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
