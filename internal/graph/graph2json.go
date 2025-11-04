package graph

import (
	"encoding/json"
	"os"
)

// GraphData rappresenta il grafo in formato esportabile
type GraphData struct {
	Nodes    []NodeData `json:"nodes"`
	Edges    []EdgeData `json:"edges"`
	Directed bool       `json:"directed"`
}

// NodeData rappresenta un nodo per l'export
type NodeData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// EdgeData rappresenta un arco per l'export
type EdgeData struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Weight int `json:"weight"`
}

// ExportToJSON esporta il grafo in formato JSON per NetworkX
func (graph *Graph) ExportToJSON(filename string) error {
	// Costruisci la struttura dati
	data := GraphData{
		Nodes:    make([]NodeData, 0, len(graph.nodes)),
		Edges:    make([]EdgeData, 0, len(graph.edges)),
		Directed: graph.directed,
	}

	// Aggiungi nodi
	for _, node := range graph.nodes {
		data.Nodes = append(data.Nodes, NodeData{
			ID:   node.nodeID,
			Name: node.name,
		})
	}

	// Aggiungi archi
	for _, edge := range graph.edges {
		data.Edges = append(data.Edges, EdgeData{
			Source: edge.startNode.nodeID,
			Target: edge.endNode.nodeID,
			Weight: edge.cost,
		})
	}

	// Scrivi su file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Formattazione leggibile
	return encoder.Encode(data)
}

// ExportToGraphML esporta in formato GraphML (alternativa a JSON)
func (graph *Graph) ExportToGraphML(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Header GraphML
	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<graphml xmlns="http://graphml.graphdrawing.org/xmlns">
  <key id="weight" for="edge" attr.name="weight" attr.type="double"/>
  <key id="name" for="node" attr.name="name" attr.type="string"/>
`)

	// Tipo di grafo
	if graph.directed {
		file.WriteString(`  <graph id="G" edgedefault="directed">` + "\n")
	} else {
		file.WriteString(`  <graph id="G" edgedefault="undirected">` + "\n")
	}

	// Nodi
	for _, node := range graph.nodes {
		file.WriteString("    <node id=\"n" + string(rune(node.nodeID+'0')) + "\">\n")
		file.WriteString("      <data key=\"name\">" + node.name + "</data>\n")
		file.WriteString("    </node>\n")
	}

	// Archi
	edgeID := 0
	for _, edge := range graph.edges {
		file.WriteString("    <edge id=\"e" + string(rune(edgeID+'0')) +
			"\" source=\"n" + string(rune(edge.startNode.nodeID+'0')) +
			"\" target=\"n" + string(rune(edge.endNode.nodeID+'0')) + "\">\n")
		file.WriteString("      <data key=\"weight\">" + string(rune(edge.cost+'0')) + "</data>\n")
		file.WriteString("    </edge>\n")
		edgeID++
	}

	// Footer
	file.WriteString("  </graph>\n</graphml>\n")

	return nil
}

// ExportToAdjacencyList esporta come lista di adiacenza semplice
func (graph *Graph) ExportToAdjacencyList(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Costruisci mappa di adiacenza
	adjList := make(map[int][]struct {
		target int
		weight int
	})

	for _, edge := range graph.edges {
		source := edge.startNode.nodeID
		target := edge.endNode.nodeID
		weight := edge.cost

		adjList[source] = append(adjList[source], struct {
			target int
			weight int
		}{target, weight})

		// Se non diretto, aggiungi anche l'inverso
		if !graph.directed {
			adjList[target] = append(adjList[target], struct {
				target int
				weight int
			}{source, weight})
		}
	}

	// Scrivi formato: "nodeID: neighbor1,weight1 neighbor2,weight2 ..."
	for nodeID := 0; nodeID < len(graph.nodes); nodeID++ {
		file.WriteString(string(rune(nodeID+'0')) + ":")
		for _, neighbor := range adjList[nodeID] {
			file.WriteString(" " + string(rune(neighbor.target+'0')) +
				"," + string(rune(neighbor.weight+'0')))
		}
		file.WriteString("\n")
	}

	return nil
}
