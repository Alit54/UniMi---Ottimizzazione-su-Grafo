package main

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"OttimizzazioneSuGrafo/internal/maxflow"
	"fmt"
	"os"
	"time"
)

func main() {
	fn := generateLectureExample(false)

	toRun := map[string]maxflow.MaxFlowAlgorithm{}
	toRun["Capacity Scaling"] = &maxflow.CapacityScaling{}
	toRun["Shortest Augmenting Path"] = &maxflow.ShortestAugmentingPath{}
	toRun["OutStars SAP"] = &maxflow.OutStarsShortestAugmentingPath{}
	toRun["Dinic"] = &maxflow.Dinic{}

	multipleBenchmark(fn, 1e6, toRun)

}

func multipleBenchmark(fn *flownetwork.FlowNetwork, iterations int, toRun map[string]maxflow.MaxFlowAlgorithm) {
	for i, alg := range toRun {
		timer := benchmark(fn, alg, iterations)
		fmt.Println("Average Time:", timer, "of algorithm:", i)
		fn.Reset()
	}
}

func benchmark(fn *flownetwork.FlowNetwork, algorithm interface{ maxflow.MaxFlowAlgorithm }, iterations int) time.Duration {
	totalTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, _ = algorithm.Run(fn, false)
		end := time.Now()
		totalTime += end.Sub(start)
	}
	averageTime := totalTime / time.Duration(iterations)
	return averageTime
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
