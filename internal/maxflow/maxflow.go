package maxflow

import "OttimizzazioneSuGrafo/internal/flownetwork"

type Stats struct {
	Advances int64
	Augments int64
	Retreats int64
	Phases   int64
}

type MaxFlowAlgorithm interface {
	Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats)
}
