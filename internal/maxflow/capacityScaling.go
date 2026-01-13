package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
	"math/bits"
)

type CapacityScaling struct{}

func (cs *CapacityScaling) Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats) {
	maxCapacity := cs.findMaxCapacity(fn)
	delta := cs.initializeDelta(maxCapacity)

	for delta >= 1 {
		stats.Phases++
		for {
			parent, edgeIndices := cs.bfs(fn, delta)
			if parent[fn.Sink] == -1 {
				break
			}

			minCap := math.MaxInt
			curr := fn.Sink
			for curr != fn.Source {
				prev := parent[curr]
				idx := edgeIndices[curr]
				edge := fn.OutStars[prev][idx]
				residual := edge.Capacity - edge.Flow
				if residual < minCap {
					minCap = residual
				}
				curr = prev
			}

			curr = fn.Sink
			for curr != fn.Source {
				prev := parent[curr]
				idx := edgeIndices[curr]
				edge := fn.OutStars[prev][idx]

				edge.Flow += minCap
				fn.OutStars[edge.To][edge.Reverse].Flow -= minCap

				curr = prev
			}

			maxFlow += minCap
			stats.Augments++
		}
		delta = delta / 2
	}
	return
}

func (cs *CapacityScaling) findMaxCapacity(fn *flownetwork.FlowNetwork) int {
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

func (cs *CapacityScaling) initializeDelta(maxCapacity int) int {
	if maxCapacity <= 0 {
		return 0
	}
	return 1 << (bits.Len(uint(maxCapacity)) - 1)
}

func (cs *CapacityScaling) bfs(fn *flownetwork.FlowNetwork, delta int) ([]int, []int) {
	parent := make([]int, fn.N)
	edgeIndices := make([]int, fn.N)
	visited := make([]bool, fn.N)
	for i := range parent {
		parent[i] = -1
	}

	queue := make([]int, 0, fn.N)
	queue = append(queue, fn.Source)
	visited[fn.Source] = true
	parent[fn.Source] = -1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == fn.Sink {
			break
		}

		for i, edge := range fn.OutStars[current] {
			residual := edge.Capacity - edge.Flow
			next := edge.To
			if !visited[next] && residual >= delta {
				visited[next] = true
				parent[next] = current
				edgeIndices[next] = i
				queue = append(queue, next)
			}
		}
	}
	return parent, edgeIndices
}
