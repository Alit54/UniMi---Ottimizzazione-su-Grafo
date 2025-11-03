package graph

type Node struct {
	name   string
	nodeID int
}

type Edge struct {
	startNode *Node
	endNode   *Node
	cost      int
}

type Graph struct {
	nodes    []*Node
	edges    []*Edge
	directed bool
}

func NewGraph(directed bool) *Graph {
	graph := new(Graph)
	graph.directed = directed
	return graph
}

func (graph *Graph) AddNode(name string) {
	graph.nodes = append(graph.nodes, &Node{name, len(graph.nodes)})
}

func (graph *Graph) AddEdge(edgeID int, cost int) {
	var inNode, outNode *Node
	n := len(graph.nodes)
	inName := edgeID / (n - 1)
	outName := edgeID % (n - 1)
	if edgeID >= inName*n {
		outName++
	}
	for _, node := range graph.nodes {
		if node.nodeID == inName {
			inNode = node
		}
		if node.nodeID == outName {
			outNode = node
		}
	}
	edge := &Edge{inNode, outNode, cost}
	graph.edges = append(graph.edges, edge)
}
