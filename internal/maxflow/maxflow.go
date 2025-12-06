package maxflow

import "OttimizzazioneSuGrafo/internal/flownetwork"

type Stats struct {
	Advances int64
	Augments int64
	Retreats int64
	Phases   int64
}

// MaxFlowAlgorithm è l'interfaccia che ogni classe rappresentante un algoritmo di maxFlow deve implementare
type MaxFlowAlgorithm interface {
	// Run è il metodo da invocare per eseguire l'algoritmo.
	Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, stats Stats)
}

// AugmentFlow aumenta il flusso di δ lungo il cammino 'path'
func AugmentFlow(fn *flownetwork.FlowNetwork, path []int, delta int) {
	for i := 0; i < len(path)-1; i++ {
		fn.PushFlow(path[i], path[i+1], delta)
	}
}
