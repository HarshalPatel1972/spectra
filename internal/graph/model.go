package graph

type NodeType string

const (
	NodeAlgorithm  NodeType = "ALGORITHM"
	NodeFile       NodeType = "FILE"
	NodeDependency NodeType = "DEPENDENCY"
	NodeCert       NodeType = "CERT"
	NodeEndpoint   NodeType = "ENDPOINT"
	NodeTeam       NodeType = "TEAM"
	NodeLayer      NodeType = "LAYER" // OCI image layer
)

type EdgeType string

const (
	EdgeUses       EdgeType = "USES"       // FILE uses ALGORITHM
	EdgeDependsOn  EdgeType = "DEPENDS_ON" // FILE depends on DEPENDENCY
	EdgeProvides   EdgeType = "PROVIDES"   // DEPENDENCY provides ALGORITHM
	EdgeSigns      EdgeType = "SIGNS"      // CERT signs something
	EdgeConfigures EdgeType = "CONFIGURES" // CONFIG configures ALGORITHM
	EdgeOwns       EdgeType = "OWNS"       // TEAM owns FILE
	EdgeServes     EdgeType = "SERVES"     // ENDPOINT serves CERT
	EdgeContains   EdgeType = "CONTAINS"   // LAYER contains FILE
)

type Node struct {
	ID         string         `json:"id"`
	Type       NodeType       `json:"type"`
	Label      string         `json:"label"`
	Properties map[string]any `json:"properties"`
}

type Edge struct {
	ID     string   `json:"id"`
	From   string   `json:"from"`
	To     string   `json:"to"`
	Type   EdgeType `json:"type"`
	Weight float64  `json:"weight"`
}

type Graph struct {
	Nodes    map[string]*Node
	Edges    []*Edge
	OutEdges map[string][]*Edge
	InEdges  map[string][]*Edge
}

// NewGraph initializes an empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Edges:    []*Edge{},
		OutEdges: make(map[string][]*Edge),
		InEdges:  make(map[string][]*Edge),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(n *Node) {
	g.Nodes[n.ID] = n
}

// AddEdge adds an edge to the graph and updates indices.
func (g *Graph) AddEdge(e *Edge) {
	g.Edges = append(g.Edges, e)
	g.OutEdges[e.From] = append(g.OutEdges[e.From], e)
	g.InEdges[e.To] = append(g.InEdges[e.To], e)
}
