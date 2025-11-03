package main

import (
	"OttimizzazioneSuGrafo/internal/graph"
	"fmt"
	"strconv"
)

func main() {
	g := graph.NewGraph(true)
	n := 50
	for i := 0; i < n; i++ {
		g.AddNode(strconv.Itoa(i))
	}
	selected := graph.CreateGraphSelectSuitablyElements(0.1, n, true)
	for i := 0; i < len(selected); i++ {
		fmt.Println(selected[i])
		g.AddEdge(selected[i], 1)
	}

	// ------------
	err := g.ExportToJSON("export/graph.json")
	if err != nil {
		fmt.Println("Errore export:", err)
		return
	}

	fmt.Println("Grafo esportato in graph.json")
}
