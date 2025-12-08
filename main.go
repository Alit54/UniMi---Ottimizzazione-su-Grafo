package main

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"fmt"
)

func main() {
	//benchmark.Run()
	generateCustomBenchmarkProblems()
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
	fn := flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, 1, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))
}
