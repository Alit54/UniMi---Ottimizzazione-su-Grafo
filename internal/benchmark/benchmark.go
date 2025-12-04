package benchmark

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"OttimizzazioneSuGrafo/internal/maxflow"
	"fmt"
	"time"
)

func Benchmark(iterations int, path string) {
	networks := LoadSuite(path)
	toRun := map[string]maxflow.MaxFlowAlgorithm{}
	toRun["Capacity Scaling"] = &maxflow.CapacityScaling{}
	toRun["Shortest Augmenting Path"] = &maxflow.ShortestAugmentingPath{}
	toRun["Dinic"] = &maxflow.Dinic{}
	toRun["Capacity Scaling con SAP"] = &maxflow.ScalingSAP{}
	toRun["Recursive Dinic"] = &maxflow.RecursiveDinic{}

	for _, network := range networks {
		fmt.Println("ALGORITMO ----", network)
		fn := flownetwork.NewNetworkFromDIMACS(network)
		fmt.Println(fn.N, len(fn.OutStars), len(fn.InStars), fn.Source, fn.Sink)
		multipleBenchmark(fn, iterations, toRun)
	}
}

func multipleBenchmark(fn *flownetwork.FlowNetwork, iterations int, toRun map[string]maxflow.MaxFlowAlgorithm) {
	for i, alg := range toRun {
		fmt.Println("Running algorithm", i)
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
