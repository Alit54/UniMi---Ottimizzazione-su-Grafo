package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type ScalingSAP struct{}

func (csSap ScalingSAP) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, iterations int) {
	maxCapacity := csSap.findMaxCapacity(fn)
	delta := csSap.initializeDelta(maxCapacity)
	sap := ShortestAugmentingPath{}
	for delta >= 1 {
		_, iter := sap.RunWithThreshold(fn, delta, saveSteps)
		iterations += iter
		delta /= 2
	}
	return fn.GetMaxFlowValue(), iterations
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
