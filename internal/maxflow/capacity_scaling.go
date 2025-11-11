package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type CapacityScaling struct{}

func (cs *CapacityScaling) Run(fn *flownetwork.FlowNetwork) (maxFlow int, iterations int) {
	maxCapacity := cs.findMaxCapacity(fn)
	delta := cs.initializeDelta(maxCapacity)
	iterations = 0
	for delta >= 1 {
		for {
			path, minCap := cs.findAugmentingPath(fn, delta)
			if path == nil {
				break
			}
			AugmentFlow(fn, path, minCap)
			iterations++
		}
		delta = delta / 2
	}
	maxFlow = fn.GetMaxFlowValue()
	return maxFlow, iterations
}

func (cs *CapacityScaling) findMaxCapacity(fn *flownetwork.FlowNetwork) int {
	maxCapacity := 0
	for i := 0; i < fn.N; i++ {
		for j := 0; j < fn.N; j++ {
			if fn.Capacity[i][j] > maxCapacity {
				maxCapacity = fn.Capacity[i][j]
			}
		}
	}
	return maxCapacity
}

func (cs *CapacityScaling) initializeDelta(maxCapacity int) int {
	if maxCapacity == 0 {
		return 0
	}
	return int(math.Pow(2, math.Floor(math.Log2(float64(maxCapacity)))))
}

func (cs *CapacityScaling) findAugmentingPath(fn *flownetwork.FlowNetwork, delta int) (path []int, minCapacity int) {
	parent, found := bfs(fn, delta)
	if !found {
		return nil, 0
	}
	path = []int{}
	for node := fn.Sink; node != -1; node = parent[node] {
		path = append(path, node)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	minCap := math.MaxInt
	for i := 0; i < len(path)-1; i++ {
		residual := fn.ResidualCapacity(path[i], path[i+1])
		if residual < minCap {
			minCap = residual
		}
	}
	return path, minCap
}

func bfs(fn *flownetwork.FlowNetwork, delta int) (parent map[int]int, found bool) {
	visited := make(map[int]bool)
	parent = make(map[int]int)
	for i := 0; i < fn.N; i++ {
		parent[i] = -1
	}
	queue := []int{fn.Source}
	visited[fn.Source] = true
	found = false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == fn.Sink {
			found = true
			break
		}
		for next := 0; next < fn.N; next++ {
			residual := fn.ResidualCapacity(current, next)
			if !visited[next] && residual >= delta {
				visited[next] = true
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}
	return parent, found
}
