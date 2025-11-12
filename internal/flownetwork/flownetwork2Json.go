/*
	Questo file è stato generato dalla IA (Claude Sonnet 4.5) per accelerare il processo del progetto. Lo scopo è generare un Json che rappresenti il flowNetwork (per poterlo visualizzare su Python).
*/

package flownetwork

import "encoding/json"

// ToJSON esporta il FlowNetwork in formato JSON per NetworkX
func (fn *FlowNetwork) ToJSON() string {
	type Edge struct {
		Source   int `json:"source"`
		Target   int `json:"target"`
		Capacity int `json:"capacity"`
		Flow     int `json:"flow"`
	}

	type Graph struct {
		Directed bool   `json:"directed"`
		Nodes    []int  `json:"nodes"`
		Edges    []Edge `json:"edges"`
		Source   int    `json:"source"`
		Sink     int    `json:"sink"`
	}

	// Costruisci lista nodi
	nodes := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		nodes[i] = i
	}

	// Costruisci lista archi (solo archi originali con capacità > 0)
	edges := []Edge{}
	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.Arcs[i] {
			// Includi solo archi con capacità > 0 (archi originali)
			// Gli archi backward hanno Capacity = 0 inizialmente
			if edge.Capacity > 0 {
				edges = append(edges, Edge{
					Source:   i,
					Target:   edge.To,
					Capacity: edge.Capacity,
					Flow:     edge.Flow,
				})
			}
		}
	}

	// Crea la struttura del grafo
	graph := Graph{
		Directed: true,
		Nodes:    nodes,
		Edges:    edges,
		Source:   fn.Source,
		Sink:     fn.Sink,
	}

	// Serializza in JSON
	jsonData, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
