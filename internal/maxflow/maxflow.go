package maxflow

import "OttimizzazioneSuGrafo/internal/flownetwork"

// MaxFlowAlgorithm è l'interfaccia che ogni classe rappresentante un algoritmo di maxFlow deve implementare
type MaxFlowAlgorithm interface {
	// Run è il metodo da invocare per eseguire l'algoritmo.
	Run(fn *flownetwork.FlowNetwork) (maxFlow int, iterations int)
}

// AugmentFlow aumenta il flusso di δ lungo il cammino 'path'
func AugmentFlow(fn *flownetwork.FlowNetwork, path []int, delta int) {
	for i := 0; i < len(path)-1; i++ {
		edgeIndex := findEdgeIndex(fn, path[i], path[i+1])
		if edgeIndex == -1 {
			panic("Arco non trovato nel path.")
		}
		fn.PushFlow(path[i], edgeIndex, delta)
	}
}

func findEdgeIndex(fn *flownetwork.FlowNetwork, from int, to int) int {
	for i, edge := range fn.Arcs[from] {
		if edge.To == to {
			return i
		}
	}
	return -1
}
