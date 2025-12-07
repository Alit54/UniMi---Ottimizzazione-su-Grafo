/*
	Questo file è stato generato dalla IA (Claude Sonnet 4.5) per accelerare il processo del progetto. Lo scopo è generare un Json che rappresenti il flowNetwork (per poterlo visualizzare su Python).
*/

package flownetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

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
		for _, edge := range fn.OutStars[i] {
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

func (fn *FlowNetwork) ToDIMACS(graphName string, path string, comment string) {
	fileName := path + "/" + graphName + ".max"
	file, _ := os.Create(fileName)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	edgeCount := 0
	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.OutStars[i] {
			if edge.Capacity > 0 {
				edgeCount++
			}
		}
	}
	maxEdges := fn.N * (fn.N - 1)
	density := float64(edgeCount) / float64(maxEdges)
	fmt.Fprintf(writer, "c %s\n", comment)
	fmt.Fprintf(writer, "c Nome del grafo: %s\n", graphName)
	fmt.Fprintf(writer, "c Densità: %f\n", density)
	fmt.Fprintf(writer, "p max %d %d\n", fn.N, edgeCount)
	fmt.Fprintf(writer, "n %d s\n", fn.Source+1)
	fmt.Fprintf(writer, "n %d t\n", fn.Sink+1)

	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.OutStars[i] {
			if edge.Capacity > 0 {
				fmt.Fprintf(writer, "a %d %d %d\n", edge.From+1, edge.To+1, edge.Capacity)
			}
		}
	}
}
