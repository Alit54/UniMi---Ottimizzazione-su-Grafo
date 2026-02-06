package main

import (
	"OttimizzazioneSuGrafo/internal/benchmark"
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"fmt"
)

func main() {
	//generateCustomBenchmarkProblems()
	//generateCapacityProblem(1000, 0.5, 100000, 0, 1, 10)
	benchmark.Run("data/flownetwork/custom", "export/benchmark_results_2.csv")
	//benchmark.Run("data/flownetwork/BVZ-tsukuba", "export/benchmark_results_2.csv")
	//benchmark.Run("data/flownetwork/KZ2-venus", "export/benchmark_results_2.csv")
	//benchmark.Run("data/flownetwork/babyface", "export/benchmark_results_2.csv")
}

func generateCustomBenchmarkProblems() {
	/*
		Generazione di diversi grafi, spaziando su pochi/tanti nodi, grafi sparsi/densi e capacità grandi/piccole
	*/
	source := 0
	sink := 1

	// Grafo Piccolo e Denso, Capacità piccole
	generateProblem(1000, 0.25, 10, source, sink)
	// Grafo Piccolo e Denso, Capacità grandi
	generateProblem(1000, 0.25, 100000, source, sink)
	// Grafo Piccolo e Sparso, Capacità piccole
	generateProblem(1000, 0.005, 10, source, sink)
	// Grafo Piccolo e Sparso, Capacità grandi
	generateProblem(1000, 0.005, 100000, source, sink)
	// Grafo Piccolo e Quasi Completo, Capacità piccole
	generateProblem(1000, 0.9, 10, source, sink)
	// Grafo Piccolo e Quasi Completo, Capacità grandi
	generateProblem(1000, 0.9, 100000, source, sink)

	// Grafo Medio e Denso, Capacità piccole
	generateProblem(5000, 0.5, 10, source, sink)
	// Grafo Medio e Denso, Capacità grandi
	generateProblem(5000, 0.5, 100000, source, sink)
	// Grafo Medio e Sparso, Capacità piccole
	generateProblem(5000, 0.005, 10, source, sink)
	// Grafo Medio e Sparso, Capacità grandi
	generateProblem(5000, 0.005, 100000, source, sink)

	// Grafo Grande e Sparso, Capacità piccole
	generateProblem(50000, 0.0001, 10, source, sink)
	// Grafo Grande e Sparso, Capacità grandi
	generateProblem(50000, 0.0001, 100000, source, sink)
	// Grafo Grande e Medio, Capacità piccole
	generateProblem(50000, 0.01, 10, source, sink)
	// Grafo Grande e Medio, Capacità grandi
	generateProblem(50000, 0.01, 100000, source, sink)

}

func generateProblem(nNodes int, density float64, maxCap int, source int, sink int) {
	for i := 0; i < 4; i++ {
		flownetwork.GenerateDIMACS(fmt.Sprintf("data/flownetwork/custom/custom.n%dd%.2fc%dmode%d.max", nNodes, density*100, maxCap, i), nNodes, density, maxCap, source, sink, i)
	}
}

func generateCapacityProblem(numNodes int, density float64, maxCap int, source int, sink int, lowFactor int) {
	fnHigh := flownetwork.NewFlowNetwork(numNodes, source, sink)
	fnHigh.GenerateRandomArcs(density, 1, maxCap)
	fnHigh.ToDIMACS("graph_high_cap", "data/flownetwork/capacity", "Grafo High Capacity")

	fnLow := flownetwork.NewFlowNetwork(numNodes, source, sink)

	for u := 0; u < fnHigh.N; u++ {
		for _, edge := range fnHigh.OutStars[u] {
			if edge.Capacity > 0 {
				newCap := edge.Capacity % lowFactor
				if newCap == 0 {
					newCap++
				}
				fnLow.AddEdge(edge.From, edge.To, newCap)
			}
		}
	}

	fnLow.ToDIMACS("graph_low_cap", "data/flownetwork/capacity", "Grafo Low Capacity")
}
