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

	// Costruisci lista archi (solo quelli con capacità > 0)
	edges := []Edge{}
	for i := 0; i < fn.N; i++ {
		for j := 0; j < fn.N; j++ {
			if fn.Capacity[i][j] > 0 {
				edges = append(edges, Edge{
					Source:   i,
					Target:   j,
					Capacity: fn.Capacity[i][j],
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
