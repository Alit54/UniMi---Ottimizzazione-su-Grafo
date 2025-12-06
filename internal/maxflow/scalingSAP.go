package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type ScalingSAP struct{}

func (csSap ScalingSAP) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, stats Stats) {
	maxCapacity := csSap.findMaxCapacity(fn)
	delta := csSap.initializeDelta(maxCapacity)
	sap := ShortestAugmentingPath{}
	for delta >= 1 {
		stats.Phases++
		tempFlow, tempStats := sap.RunWithThreshold(fn, delta, saveSteps)
		maxFlow += tempFlow
		stats.Advance += tempStats.Advance
		stats.Retreats += tempStats.Retreats
		stats.Augments += tempStats.Augments
		delta /= 2
	}
	return
}

func (csSap ScalingSAP) initializeDelta(maxCapacity int) int {
	delta := 0
	if maxCapacity > 0 {
		delta = int(math.Pow(2, math.Floor(math.Log2(float64(maxCapacity)))))
	}
	return delta
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
