package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type CapacityScaling struct{}

func (cs *CapacityScaling) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, stats Stats) {
	maxCapacity := cs.findMaxCapacity(fn)
	delta := cs.initializeDelta(maxCapacity)
	step := 0
	if saveSteps {
		cs.saveStep(fn, step, nil, 0, delta, "Fase iniziale")
	}
	for delta >= 1 {
		stats.Phases++
		for {
			path, minCap := cs.findAugmentingPath(fn, delta)
			if path == nil {
				break
			}
			if saveSteps {
				step++
				cs.saveStep(fn, step, path, minCap, delta, fmt.Sprintf("Trovato percorso con capacita' %d", minCap))
			}
			AugmentFlow(fn, path, minCap)
			stats.Augments++
			if saveSteps {
				step++
				cs.saveStep(fn, step, path, minCap, delta, fmt.Sprintf("Aumento di %d unita'", minCap))
			}
		}
		delta = delta / 2
		if saveSteps {
			step++
			cs.saveStep(fn, step, nil, 0, delta, fmt.Sprintf("Nessun percorso con capacita' superiore a %d, passiamo a %d", delta*2, delta))
		}
	}
	maxFlow = fn.GetMaxFlowValue()
	if saveSteps {
		step++
		cs.saveStep(fn, step, nil, 0, 0, fmt.Sprintf("Fase finale - MaxFlow = %d", maxFlow))
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
		from := path[i]
		to := path[i+1]
		for _, edge := range fn.OutStars[from] {
			if edge.To == to {
				residual := edge.Capacity - edge.Flow
				if residual < minCap {
					minCap = residual
				}
				break
			}
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
		for _, edge := range fn.OutStars[current] {
			residual := edge.Capacity - edge.Flow
			next := edge.To
			if !visited[next] && residual >= delta {
				visited[next] = true
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}
	return parent, found
}

func (cs *CapacityScaling) saveStep(
	fn *flownetwork.FlowNetwork,
	step int,
	path []int,
	delta int,
	scalingDelta int,
	description string) {
	snapshot := cs.createSnapshot(fn, path, delta, scalingDelta, description)
	filename := fmt.Sprintf("export/graphical_steps/capacity_scaling/step_%04d.json", step)
	file, _ := os.Create(filename)
	defer file.Close()
	file.WriteString(snapshot)
}

func (cs *CapacityScaling) createSnapshot(
	fn *flownetwork.FlowNetwork,
	path []int,
	pathCapacity int,
	scalingDelta int,
	description string,
) string {
	type EdgeInfo struct {
		Source      int  `json:"source"`
		Target      int  `json:"target"`
		Capacity    int  `json:"capacity"`
		Flow        int  `json:"flow"`
		Residual    int  `json:"residual"`
		InPath      bool `json:"in_path"`
		IsSaturated bool `json:"is_saturated"`
	}

	type Snapshot struct {
		Description  string     `json:"description"`
		ScalingDelta int        `json:"scaling_delta"`
		PathCapacity int        `json:"path_capacity"`
		CurrentFlow  int        `json:"current_flow"`
		Path         []int      `json:"path"`
		Nodes        []int      `json:"nodes"`
		Edges        []EdgeInfo `json:"edges"`
		Source       int        `json:"source"`
		Sink         int        `json:"sink"`
	}

	// Crea set dei nodi nel path per lookup veloce
	pathSet := make(map[int]bool)
	for _, node := range path {
		pathSet[node] = true
	}

	// Crea set degli archi nel path
	pathEdges := make(map[string]bool)
	for i := 0; i < len(path)-1; i++ {
		key := fmt.Sprintf("%d-%d", path[i], path[i+1])
		pathEdges[key] = true
	}

	// Costruisci lista nodi
	nodes := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		nodes[i] = i
	}

	// Costruisci lista archi
	edges := []EdgeInfo{}
	for i := 0; i < fn.N; i++ {
		for _, edge := range fn.OutStars[i] {
			if edge.Capacity > 0 { // Solo archi originali
				residual := edge.Capacity - edge.Flow
				key := fmt.Sprintf("%d-%d", i, edge.To)

				edges = append(edges, EdgeInfo{
					Source:      i,
					Target:      edge.To,
					Capacity:    edge.Capacity,
					Flow:        edge.Flow,
					Residual:    residual,
					InPath:      pathEdges[key],
					IsSaturated: residual == 0,
				})
			}
		}
	}

	snapshot := Snapshot{
		Description:  description,
		ScalingDelta: scalingDelta,
		PathCapacity: pathCapacity,
		CurrentFlow:  fn.GetMaxFlowValue(),
		Path:         path,
		Nodes:        nodes,
		Edges:        edges,
		Source:       fn.Source,
		Sink:         fn.Sink,
	}

	// Serializza in JSON
	jsonData, _ := json.MarshalIndent(snapshot, "", "  ")
	return string(jsonData)
}
