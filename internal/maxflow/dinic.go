package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type Dinic struct{}

func (d *Dinic) Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats) {
	for {
		level := d.bfs(fn)
		stats.Phases++

		if level[fn.Sink] == -1 {
			break
		}

		work := make([]int, fn.N)
		for {
			flow := d.dfs(fn, level, work)
			if flow == 0 {
				break
			}
			maxFlow += flow
			stats.Augments++
		}
	}
	return
}

func (d *Dinic) bfs(fn *flownetwork.FlowNetwork) []int {
	label := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		label[i] = -1
	}
	queue := make([]int, 0, fn.N)
	queue = append(queue, fn.Source)
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

func (d *Dinic) dfs(fn *flownetwork.FlowNetwork, level []int, work []int) int {
	predecessor := make([]int, fn.N)
	edgeIndex := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		predecessor[i] = -1
		edgeIndex[i] = -1
	}

	stack := make([]int, 0, fn.N)
	stack = append(stack, fn.Source)

	for len(stack) > 0 {
		current := stack[len(stack)-1]

		if current == fn.Sink {
			break
		}
		found := false

		for work[current] < len(fn.OutStars[current]) {
			idx := work[current]
			edge := fn.OutStars[current][idx]
			residual := edge.Capacity - edge.Flow

			if residual > 0 && level[edge.To] == level[current]+1 {
				predecessor[edge.To] = current
				edgeIndex[edge.To] = idx
				stack = append(stack, edge.To)
				found = true
				break
			}
			work[current]++
		}

		if !found {
			stack = stack[:len(stack)-1]

			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				work[parent]++
			}
		}
	}

	if predecessor[fn.Sink] == -1 {
		return 0
	}

	flow := math.MaxInt
	curr := fn.Sink
	for curr != fn.Source {
		prev := predecessor[curr]
		idx := edgeIndex[curr]
		edge := fn.OutStars[prev][idx]
		residual := edge.Capacity - edge.Flow
		if residual < flow {
			flow = residual
		}
		curr = prev
	}

	curr = fn.Sink
	for curr != fn.Source {
		prev := predecessor[curr]
		idx := edgeIndex[curr]

		edge := fn.OutStars[prev][idx]
		edge.Flow += flow
		fn.OutStars[edge.To][edge.Reverse].Flow -= flow

		curr = prev
	}

	return flow
}
