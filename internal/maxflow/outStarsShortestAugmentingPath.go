package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type OutStarsShortestAugmentingPath struct{}

func (sap *OutStarsShortestAugmentingPath) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, iterations int) {
	advances, retreats, augments := 0, 0, 0
	distance := sap.exactDistance(fn)
	if distance[fn.Source] >= fn.N {
		return 0, 0
	}
	iterations = 0
	current := fn.Source
	predecessor := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		predecessor[i] = -1
	}
	for distance[fn.Source] < fn.N {
		admissibleEdge := sap.findAdmissibleEdge(fn, current, distance)
		if admissibleEdge != -1 {
			// ADVANCE
			advances++
			edge := fn.OutStars[current][admissibleEdge]
			next := edge.To
			predecessor[next] = current
			current = next
			if current == fn.Sink {
				sap.augment(fn, predecessor)
				augments++
				iterations++
				current = fn.Source
			}
		} else {
			sap.retreat(fn, current, distance)
			retreats++
			if current != fn.Source {
				current = predecessor[current]
			}
		}
	}
	maxFlow = fn.GetMaxFlowValue()
	/*fmt.Println("Advances: ", advances)
	fmt.Println("Retreats: ", retreats)
	fmt.Println("Augments: ", augments)
	*/return
}

func (sap *OutStarsShortestAugmentingPath) exactDistance(fn *flownetwork.FlowNetwork) []int {
	distance := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		distance[i] = fn.N
	}
	distance[fn.Sink] = 0
	queue := []int{fn.Sink}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for previous := 0; previous < fn.N; previous++ {
			if distance[previous] < fn.N {
				continue
			}
			for _, edge := range fn.OutStars[previous] {
				if edge.To == current && edge.Capacity-edge.Flow > 0 {
					distance[previous] = distance[current] + 1
					queue = append(queue, previous)
					break
				}
			}
		}
	}
	return distance
}

func (sap *OutStarsShortestAugmentingPath) findAdmissibleEdge(fn *flownetwork.FlowNetwork, node int, distance []int) int {
	for id, edge := range fn.OutStars[node] {
		residual := edge.Capacity - edge.Flow
		if residual > 0 && distance[node] == distance[edge.To]+1 {
			return id
		}
	}
	return -1
}

func (sap *OutStarsShortestAugmentingPath) retreat(fn *flownetwork.FlowNetwork, node int, distance []int) {
	minDistance := fn.N
	for _, edge := range fn.OutStars[node] {
		residual := edge.Capacity - edge.Flow
		if residual > 0 {
			if distance[edge.To] < minDistance {
				minDistance = distance[edge.To]
			}
		}
	}
	distance[node] = minDistance + 1
}

func (sap *OutStarsShortestAugmentingPath) augment(fn *flownetwork.FlowNetwork, predecessor []int) int {
	path := []int{}
	node := fn.Sink
	for node != fn.Source {
		path = append([]int{node}, path...)
		node = predecessor[node]
	}
	path = append([]int{fn.Source}, path...)
	delta := math.MaxInt
	for i := 0; i < len(path)-1; i++ {
		from := path[i]
		to := path[i+1]
		for _, edge := range fn.OutStars[from] {
			if edge.To == to {
				residual := edge.Capacity - edge.Flow
				if residual < delta {
					delta = residual
				}
				break
			}
		}
	}
	AugmentFlow(fn, path, delta)
	return delta
}
