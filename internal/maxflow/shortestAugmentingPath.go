package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type ShortestAugmentingPath struct{}

func (sap *ShortestAugmentingPath) Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats) {
	return sap.RunWithThreshold(fn, 1)
}

func (sap *ShortestAugmentingPath) RunWithThreshold(fn *flownetwork.FlowNetwork, threshold int) (maxFlow int, stats Stats) {
	distance, number := sap.exactDistance(fn, threshold)
	if distance[fn.Source] >= fn.N {
		return
	}

	predecessor := make([]int, fn.N)
	edgeIndices := make([]int, fn.N)
	currentArc := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		predecessor[i] = -1
		edgeIndices[i] = -1
	}

	current := fn.Source
	for distance[fn.Source] < fn.N {
		admissibleEdge := sap.findAdmissibleEdge(fn, current, distance, currentArc, threshold)
		if admissibleEdge != -1 {
			stats.Advances++
			edge := fn.OutStars[current][admissibleEdge]
			next := edge.To
			predecessor[next] = current
			edgeIndices[next] = admissibleEdge
			current = next

			if current == fn.Sink {
				delta := sap.augment(fn, predecessor, edgeIndices)
				maxFlow += delta
				stats.Augments++
				current = fn.Source
			}
		} else {
			endTest := sap.retreat(fn, current, distance, number, currentArc, threshold)
			if endTest {
				break
			}
			stats.Retreats++
			if current != fn.Source {
				current = predecessor[current]
			}
		}
	}
	return
}

func (sap *ShortestAugmentingPath) exactDistance(fn *flownetwork.FlowNetwork, threshold int) ([]int, []int) {
	distance := make([]int, fn.N)
	number := make([]int, fn.N+1)
	for i := 0; i < fn.N; i++ {
		distance[i] = fn.N
	}
	distance[fn.Sink] = 0
	number[0] = 1

	queue := make([]int, 0, fn.N)
	queue = append(queue, fn.Sink)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range fn.OutStars[current] {
			neighbor := edge.To
			reverseEdge := fn.OutStars[neighbor][edge.Reverse]
			residual := reverseEdge.Capacity - reverseEdge.Flow

			if distance[neighbor] == fn.N && residual >= threshold {
				distance[neighbor] = distance[current] + 1
				number[distance[neighbor]]++
				queue = append(queue, neighbor)
			}
		}
	}
	return distance, number
}

func (sap *ShortestAugmentingPath) findAdmissibleEdge(fn *flownetwork.FlowNetwork, node int, distance []int, currentArc []int, threshold int) int {
	for i := currentArc[node]; i < len(fn.OutStars[node]); i++ {
		edge := fn.OutStars[node][i]
		residual := edge.Capacity - edge.Flow
		if residual >= threshold && distance[node] == distance[edge.To]+1 {
			currentArc[node] = i
			return i
		}
	}
	currentArc[node] = len(fn.OutStars[node])
	return -1
}

func (sap *ShortestAugmentingPath) retreat(fn *flownetwork.FlowNetwork, node int, distance []int, number []int, currentArc []int, threshold int) bool {
	minDistance := fn.N
	for _, edge := range fn.OutStars[node] {
		residual := edge.Capacity - edge.Flow
		if residual >= threshold {
			if distance[edge.To] < minDistance {
				minDistance = distance[edge.To]
			}
		}
	}

	number[distance[node]]--
	if number[distance[node]] == 0 {
		return true
	}

	distance[node] = minDistance + 1
	number[distance[node]]++
	currentArc[node] = 0
	return false
}

func (sap *ShortestAugmentingPath) augment(fn *flownetwork.FlowNetwork, predecessor []int, edgeIndices []int) int {
	delta := math.MaxInt

	curr := fn.Sink
	for curr != fn.Source {
		prev := predecessor[curr]
		idx := edgeIndices[curr]
		edge := fn.OutStars[prev][idx]
		residual := edge.Capacity - edge.Flow
		if residual < delta {
			delta = residual
		}
		curr = prev
	}

	curr = fn.Sink
	for curr != fn.Source {
		prev := predecessor[curr]
		idx := edgeIndices[curr]

		edge := fn.OutStars[prev][idx]
		edge.Flow += delta
		fn.OutStars[edge.To][edge.Reverse].Flow -= delta

		curr = prev
	}
	return delta
}
