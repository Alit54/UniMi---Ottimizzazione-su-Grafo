package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type Dinic struct{}

func (d *Dinic) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, iterations int) {
	iterations = 0
	for {
		level := d.bfs(fn)
		if level[fn.Sink] == -1 {
			break
		}
		for {
			flow := d.dfs(fn, level)
			if flow == 0 {
				break
			}
			maxFlow += flow
			iterations++
		}
	}
	return maxFlow, iterations
}

func (d *Dinic) bfs(fn *flownetwork.FlowNetwork) []int {
	label := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		label[i] = -1
	}
	queue := []int{fn.Source}
	label[fn.Source] = 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, edge := range fn.OutStars[current] {
			residual := edge.Capacity - edge.Flow
			next := edge.To
			if label[next] == -1 && residual > 0 {
				label[next] = label[current] + 1
				queue = append(queue, next)
			}
		}
	}
	return label
}

func (d *Dinic) dfs(fn *flownetwork.FlowNetwork, level []int) int {
	predecessor := make([]int, fn.N)
	edgeIndex := make([]int, fn.N)
	stack := []int{fn.Source}
	for i := 0; i < fn.N; i++ {
		predecessor[i] = -1
		edgeIndex[i] = -1
	}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current == fn.Sink {
			break
		}
		for id, edge := range fn.OutStars[current] {
			residual := edge.Capacity - edge.Flow
			next := edge.To
			if predecessor[next] == -1 && residual > 0 && level[next] == level[current]+1 {
				predecessor[next] = current
				edgeIndex[next] = id
				stack = append(stack, next)
			}
		}
	}
	if predecessor[fn.Sink] == -1 {
		return 0
	}
	flow := math.MaxInt
	current := fn.Sink
	for current != fn.Source {
		previous := predecessor[current]
		id := edgeIndex[current]
		edge := fn.OutStars[previous][id]
		residual := edge.Capacity - edge.Flow
		if residual < flow {
			flow = residual
		}
		current = previous
	}
	current = fn.Sink
	for current != fn.Source {
		previous := predecessor[current]
		id := edgeIndex[current]
		fn.PushFlowWithIndex(previous, id, flow)
		current = previous
	}
	return flow
}
