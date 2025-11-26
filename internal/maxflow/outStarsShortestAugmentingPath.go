package maxflow

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"encoding/json"
	"fmt"
	"math"
	"os"
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
	step := 0
	if saveSteps {
		sap.saveStep(fn, step, "Start", "Stato Iniziale", current, distance, nil)
	}
	for distance[fn.Source] < fn.N {
		admissibleEdge := sap.findAdmissibleEdge(fn, current, distance)
		if admissibleEdge != -1 {
			// ADVANCE
			advances++
			iterations++
			edge := fn.OutStars[current][admissibleEdge]
			next := edge.To
			predecessor[next] = current
			current = next
			if saveSteps {
				step++
				sap.saveStep(fn, step, "Advance", fmt.Sprintf("Avanzamento -> Nodo %d", current), current, distance, nil)
			}
			if current == fn.Sink {
				delta, path := sap.augment(fn, predecessor)
				augments++
				iterations++
				if saveSteps {
					step++
					sap.saveStep(fn, step, "Augment", fmt.Sprintf("Augmenting Path (Cap: %d)", delta), current, distance, path)
				}
				current = fn.Source
				iterations++
				if saveSteps {
					step++
					sap.saveStep(fn, step, "Restart", "Ritorno a Source dopo Augment", current, distance, nil)
				}
			}
		} else {
			oldDist := distance[current]
			sap.retreat(fn, current, distance)
			retreats++
			iterations++
			if saveSteps {
				step++
				sap.saveStep(fn, step, "Retreat", fmt.Sprintf("Retreat Nodo %d (h: %d -> %d)", current, oldDist, distance[current]), current, distance, nil)
			}
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

func (sap *OutStarsShortestAugmentingPath) augment(fn *flownetwork.FlowNetwork, predecessor []int) (int, []int) {
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
	return delta, path
}

func (sap *OutStarsShortestAugmentingPath) saveStep(
	fn *flownetwork.FlowNetwork,
	step int,
	stepType string,
	description string,
	currentNode int,
	distances []int,
	path []int,
) {
	type NodeInfo struct {
		ID        int  `json:"id"`
		Distance  int  `json:"distance"`
		IsSource  bool `json:"is_source"`
		IsSink    bool `json:"is_sink"`
		IsCurrent bool `json:"is_current"`
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
		Iteration   int        `json:"iteration"`
		StepType    string     `json:"step_type"`
		Description string     `json:"description"`
		Nodes       []NodeInfo `json:"nodes"`
		Edges       []EdgeInfo `json:"edges"`
	}

	// Nodi
	nodes := make([]NodeInfo, fn.N)
	for i := 0; i < fn.N; i++ {
		nodes[i] = NodeInfo{
			ID:        i,
			Distance:  distances[i],
			IsSource:  i == fn.Source,
			IsSink:    i == fn.Sink,
			IsCurrent: i == currentNode,
		}
	}

	// Path Lookup
	pathSet := make(map[string]bool)
	if path != nil {
		for i := 0; i < len(path)-1; i++ {
			key := fmt.Sprintf("%d-%d", path[i], path[i+1])
			pathSet[key] = true
		}
	}

	// Archi
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

	// Assicurati che la cartella esista
	folder := "export/maxflow/outstars_shortest_augmenting_path"
	_ = os.MkdirAll(folder, os.ModePerm)

	filename := fmt.Sprintf("%s/step_%04d.json", folder, step)
	file, _ := os.Create(filename)
	defer file.Close()

	jsonData, _ := json.MarshalIndent(snap, "", "  ")
	file.Write(jsonData)
}
