package main

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"fmt"
	"os"
)

func main() {
	generateRandomFlow(10, 0, 5, 0.3, 0, 100, true)
}

func generateLectureExample(save bool) *flownetwork.FlowNetwork {
	fn := flownetwork.NewFlowNetwork(7, 0, 6)
	fn.AddEdge(0, 1, 6)
	fn.AddEdge(0, 2, 10)
	fn.AddEdge(0, 3, 5)
	fn.AddEdge(1, 3, 6)
	fn.AddEdge(1, 4, 1)
	fn.AddEdge(2, 3, 6)
	fn.AddEdge(2, 5, 6)
	fn.AddEdge(3, 4, 3)
	fn.AddEdge(3, 5, 4)
	fn.AddEdge(3, 6, 5)
	fn.AddEdge(4, 6, 10)
	fn.AddEdge(5, 6, 1)
	if save {
		jsonStr := fn.ToJSON()
		err := os.WriteFile("export/lectureExample.json", []byte(jsonStr), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("FlowNetwork generato e salvato in export/lectureExample.json")
	}
	return fn
}

func generateRandomFlow(numberNode int, source int, sink int, density float64, minCap int, maxCap int, save bool) *flownetwork.FlowNetwork {
	fn := flownetwork.NewFlowNetwork(numberNode, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	if save {
		jsonStr := fn.ToJSON()
		err := os.WriteFile("export/flownetwork.json", []byte(jsonStr), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("FlowNetwork generato e salvato in flownetwork.json")
		fmt.Printf("Nodi: %d, Archi generati: circa %.0f\n", fn.N, 0.3*float64(fn.N*(fn.N-1)))
	}
	return fn
}
