package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type RecursiveDinic struct{}

func (d *RecursiveDinic) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, stats Stats) {
	step := 0
	if saveSteps {
		initLevels := make([]int, fn.N)
		for i := range initLevels {
			initLevels[i] = -1
		}
		d.saveStep(fn, step, "Start", "Stato Iniziale", initLevels, nil)
	}
	for {
		level := d.bfs(fn)
		stats.Phases++
		if level[fn.Sink] == -1 {
			break
		}
		if saveSteps {
			step++
			d.saveStep(fn, step, "BFS", "Level Graph costruito", level, nil)
		}
		currentArc := make([]int, fn.N)
		for {
			flow, path := d.dfs(fn, fn.Source, math.MaxInt, level, currentArc)
			if flow == 0 {
				break
			}
			maxFlow += flow
			stats.Augments++
			if saveSteps {
				step++
				d.saveStep(fn, step, "Flow Pushed", fmt.Sprintf("Flusso aumentato di %d", flow), level, path)
			}
		}
	}
	if saveSteps {
		step++
		finalLevels := make([]int, fn.N)
		for i := range finalLevels {
			finalLevels[i] = -1
		}
		d.saveStep(fn, step, "End", fmt.Sprintf("Terminato. MaxFlow: %d", maxFlow), finalLevels, nil)
	}
	return
}

func (d *RecursiveDinic) bfs(fn *flownetwork.FlowNetwork) []int {
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

func (d *RecursiveDinic) dfs(fn *flownetwork.FlowNetwork, current int, flow int, level []int, currentArc []int) (int, []int) {
	if current == fn.Sink {
		return flow, []int{current}
	}
	for i := currentArc[current]; i < len(fn.OutStars[current]); i++ {
		edge := fn.OutStars[current][i]
		if edge.Capacity-edge.Flow > 0 && level[edge.To] == level[current]+1 {
			pushed := min(flow, edge.Capacity-edge.Flow)
			returnedFlow, path := d.dfs(fn, edge.To, pushed, level, currentArc)
			if returnedFlow > 0 {
				edge.Flow += returnedFlow
				currentArc[current] = i
				if path != nil {
					path = append([]int{current}, path...)
				}
				return returnedFlow, path
			}
		}
		currentArc[current]++
	}
	return 0, nil
}

func (d *RecursiveDinic) saveStep(
	fn *flownetwork.FlowNetwork,
	step int,
	stepType string,
	description string,
	levels []int,
	path []int,
) {
	type NodeInfo struct {
		ID       int  `json:"id"`
		Level    int  `json:"level"`
		IsSource bool `json:"is_source"`
		IsSink   bool `json:"is_sink"`
	}

	type EdgeInfo struct {
		From     int  `json:"from"`
		To       int  `json:"to"`
		Capacity int  `json:"capacity"`
		Flow     int  `json:"flow"`
		Residual int  `json:"residual"`
		InPath   bool `json:"in_path"`
	}

	type Snapshot struct {
		Iteration   int        `json:"step"`
		StepType    string     `json:"step_type"`
		Description string     `json:"description"`
		Nodes       []NodeInfo `json:"nodes"`
		Edges       []EdgeInfo `json:"edges"`
	}

	nodes := make([]NodeInfo, fn.N)
	for i := 0; i < fn.N; i++ {
		nodes[i] = NodeInfo{
			ID:       i,
			Level:    levels[i],
			IsSource: i == fn.Source,
			IsSink:   i == fn.Sink,
		}
	}

	pathSet := make(map[string]bool)
	if path != nil {
		for i := 0; i < len(path)-1; i++ {
			key := fmt.Sprintf("%d-%d", path[i], path[i+1])
			pathSet[key] = true
		}
	}

	edges := []EdgeInfo{}
	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.OutStars[i] {
			if edge.Capacity > 0 {
				residual := edge.Capacity - edge.Flow
				key := fmt.Sprintf("%d-%d", i, edge.To)
				edges = append(edges, EdgeInfo{
					From:     i,
					To:       edge.To,
					Capacity: edge.Capacity,
					Flow:     edge.Flow,
					Residual: residual,
					InPath:   pathSet[key],
				})
			}
		}
	}

	snap := Snapshot{
		Iteration:   step,
		StepType:    stepType,
		Description: description,
		Nodes:       nodes,
		Edges:       edges,
	}

	// Cartella RecursiveDinic
	folder := "export/graphical_steps/dinic"
	_ = os.MkdirAll(folder, os.ModePerm)

	filename := fmt.Sprintf("%s/step_%04d.json", folder, step)
	file, _ := os.Create(filename)
	defer file.Close()
	jsonData, _ := json.MarshalIndent(snap, "", "  ")
	file.Write(jsonData)
}
