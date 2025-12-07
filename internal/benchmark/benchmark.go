package benchmark

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"OttimizzazioneSuGrafo/internal/maxflow"
	"fmt"
	"time"
)

type Result struct {
	Graph     string
	Nodes     int
	Edges     int
	Algorithm string
	MaxFlow   int
	Time      time.Duration
	Augments  int64
	Retreats  int64
	Advances  int64
	Phases    int64
}

func Benchmark(fn *flownetwork.FlowNetwork, nameAlg string, algorithm interface{ maxflow.MaxFlowAlgorithm }, path string, iterations int) Result {
	totalTime := time.Duration(0)
	maxFlow := 0
	stats := maxflow.Stats{}
	for i := 0; i < iterations; i++ {
		fmt.Println("RUNNING", nameAlg, "on Instance", path)
		start := time.Now()
		maxFlow, stats = algorithm.Run(fn, false)
		end := time.Now()
		totalTime += end.Sub(start)
	}
	averageTime := totalTime / time.Duration(iterations)
	return Result{
		Graph:     path,
		Nodes:     fn.N,
		Edges:     len(fn.OutStars),
		Algorithm: nameAlg,
		MaxFlow:   maxFlow,
		Time:      averageTime,
		Augments:  stats.Augments,
		Retreats:  stats.Retreats,
		Advances:  stats.Advances,
		Phases:    stats.Phases,
	}
}
