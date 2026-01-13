package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math/bits"
)

type ScalingSAP struct{}

func (csSap ScalingSAP) Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats) {
	maxCapacity := csSap.findMaxCapacity(fn)
	delta := csSap.initializeDelta(maxCapacity)
	sap := ShortestAugmentingPath{}
	for delta >= 1 {
		stats.Phases++
		_, tempStats := sap.RunWithThreshold(fn, delta)
		stats.Advances += tempStats.Advances
		stats.Retreats += tempStats.Retreats
		stats.Augments += tempStats.Augments
		delta /= 2
	}
	return fn.GetMaxFlowValue(), stats
}

func (csSap ScalingSAP) initializeDelta(maxCapacity int) int {
	if maxCapacity <= 0 {
		return 0
	}
	return 1 << (bits.Len(uint(maxCapacity)) - 1)
}

func (csSap ScalingSAP) findMaxCapacity(fn *flownetwork.FlowNetwork) int {
	maxCapacity := 0
	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.OutStars[i] {
			if edge.Capacity > maxCapacity {
				maxCapacity = edge.Capacity
			}
		}
	}
	return maxCapacity
}
