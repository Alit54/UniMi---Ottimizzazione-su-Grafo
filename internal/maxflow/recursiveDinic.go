package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"math"
)

type RecursiveDinic struct{}

func (d *RecursiveDinic) Run(fn *flownetwork.FlowNetwork) (maxFlow int, stats Stats) {
	for {
		level := d.bfs(fn)
		stats.Phases++
		if level[fn.Sink] == -1 {
			break
		}

		currentArc := make([]int, fn.N)
		for {
			flow := d.dfs(fn, fn.Source, math.MaxInt, level, currentArc)
			if flow == 0 {
				break
			}
			maxFlow += flow
			stats.Augments++
		}
	}
	return
}

func (d *RecursiveDinic) bfs(fn *flownetwork.FlowNetwork) []int {
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

func (d *RecursiveDinic) dfs(fn *flownetwork.FlowNetwork, current int, flow int, level []int, currentArc []int) int {
	if current == fn.Sink {
		return flow
	}

	for i := currentArc[current]; i < len(fn.OutStars[current]); i++ {
		currentArc[current] = i
		edge := fn.OutStars[current][i]
		residual := edge.Capacity - edge.Flow

		if residual > 0 && level[edge.To] == level[current]+1 {
			pushed := flow
			if residual < pushed {
				pushed = residual
			}

			returnedFlow := d.dfs(fn, edge.To, pushed, level, currentArc)

			if returnedFlow > 0 {
				fn.OutStars[current][i].Flow += returnedFlow
				fn.OutStars[edge.To][edge.Reverse].Flow -= returnedFlow
				return returnedFlow
			}
		}
	}

	currentArc[current] = len(fn.OutStars[current])
	return 0
}
