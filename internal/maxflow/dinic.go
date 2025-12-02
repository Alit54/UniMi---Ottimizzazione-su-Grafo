package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Dinic struct{}

func (d *Dinic) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, iterations int) {
	iterations = 0
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
		if level[fn.Sink] == -1 {
			break
		}
		if saveSteps {
			step++
			d.saveStep(fn, step, "BFS", "Level Graph costruito", level, nil)
		}
		for {
			flow, path := d.dfs(fn, level, saveSteps, &step)
			if flow == 0 {
				break
			}
			maxFlow += flow
			iterations++
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

func (d *Dinic) dfs(fn *flownetwork.FlowNetwork, level []int, saveSteps bool, step *int) (int, []int) {
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
		return 0, nil
	}
	flow := math.MaxInt
	path := []int{}
	tempCurr := fn.Sink
	for tempCurr != fn.Source {
		path = append([]int{tempCurr}, path...)
		prev := predecessor[tempCurr]
		id := edgeIndex[tempCurr]
		edge := fn.OutStars[prev][id]
		if edge.Capacity-edge.Flow < flow {
			flow = edge.Capacity - edge.Flow
		}
		tempCurr = prev
	}
	path = append([]int{fn.Source}, path...)
	if saveSteps {
		*step++
		d.saveStep(fn, *step, "Augment Found", fmt.Sprintf("Blocking Path trovato (Cap: %d)", flow), level, path)
	}
	current := fn.Sink
	for current != fn.Source {
		previous := predecessor[current]
		id := edgeIndex[current]
		fn.PushFlowWithIndex(previous, id, flow)
		current = previous
	}
	return flow, path
}

func (d *Dinic) saveStep(
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

	// Cartella Dinic
	folder := "export/maxflow/dinic"
	_ = os.MkdirAll(folder, os.ModePerm)

	filename := fmt.Sprintf("%s/step_%04d.json", folder, step)
	file, _ := os.Create(filename)
	defer file.Close()
	jsonData, _ := json.MarshalIndent(snap, "", "  ")
	file.Write(jsonData)
}
