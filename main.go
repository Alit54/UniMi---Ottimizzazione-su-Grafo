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

	// Grafo piccolo e denso, capacità piccole
	source := 0
	sink := 1
	nNodes := 1000
	density := 0.5
	minCap := 1
	maxCap := 10
	fn := &flownetwork.FlowNetwork{}
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo piccolo e denso, capacità grandi
	nNodes = 1000
	density = 0.5
	minCap = 1
	maxCap = 100000
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo piccolo e sparso, capacità piccole
	nNodes = 1000
	density = 0.01
	minCap = 1
	maxCap = 10
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo piccolo e sparso, capacità grandi
	nNodes = 1000
	density = 0.01
	minCap = 1
	maxCap = 100000
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo medio e denso, capacità piccole
	nNodes = 20000
	density = 0.5
	minCap = 1
	maxCap = 10
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo medio e denso, capacità grandi
	nNodes = 20000
	density = 0.5
	minCap = 1
	maxCap = 100000
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo grande e sparso, capacità piccole
	nNodes = 1000000
	density = 0.0001
	minCap = 1
	maxCap = 10
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo grande e sparso, capacità grandi
	nNodes = 1000000
	density = 0.0001
	minCap = 1
	maxCap = 100000
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo grande e "meno denso", capacità piccole
	nNodes = 500000
	density = 0.002
	minCap = 1
	maxCap = 10
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))

	// Grafo grande e "meno denso", capacità grandi
	nNodes = 500000
	density = 0.002
	minCap = 1
	maxCap = 100000
	fn = flownetwork.NewFlowNetwork(nNodes, source, sink)
	fn.GenerateRandomArcs(density, minCap, maxCap)
	fn.ToDIMACS(fmt.Sprintf("custom.n%dd%.2fc%d", nNodes, density*100, maxCap), "data/flownetwork/custom", fmt.Sprintf("Grafo generato casualmente dalla funzione GenerateRandomArcs() con %d nodi e una densità di %.2f. La massima capacità degli archi è %d", nNodes, density, maxCap))
}
