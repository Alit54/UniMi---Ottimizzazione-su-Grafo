package maxflow

/*import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type BasicShortestAugmentingPath struct{}

func (sap *BasicShortestAugmentingPath) Run(fn *flownetwork.FlowNetwork, saveSteps bool) (maxFlow int, stats Stats) {
	current := fn.Source
	step := 0
	if saveSteps {
		initialDist := make([]int, fn.N)
		sap.saveStep(fn, step, "Start", "Stato iniziale", current, initialDist, nil)
	}
	distance, number := sap.exactDistance(fn)
	if distance[fn.Source] >= fn.N {
		return
	}
	predecessor := make([]int, fn.N)
	for i := 0; i < fn.N; i++ {
		predecessor[i] = -1
	}
	if saveSteps {
		step++
		sap.saveStep(fn, step, "Start", "Distanze calcolate", current, distance, nil)
	}
	for distance[fn.Source] < fn.N {
		admissibleEdge := sap.findAdmissibleEdge(fn, current, distance)
		if admissibleEdge != -1 {
			// ADVANCE
			stats.Advances++
			edge := fn.OutStars[current][admissibleEdge]
			next := edge.To
			predecessor[next] = current
			current = next
			if saveSteps {
				step++
				sap.saveStep(fn, step, "Advances", fmt.Sprintf("Advances al nodo %d", current), current, distance, nil)
			}
			if current == fn.Sink {
				delta, path := sap.augment(fn, predecessor)
				stats.Augments++
				if saveSteps {
					step++
					sap.saveStep(fn, step, "Augment", fmt.Sprintf("Augmenting Path trovato (Cap: %d)", delta), current, distance, path)
				}
				current = fn.Source
			}
		} else {
			oldLabel := distance[current]
			endTest := sap.retreat(fn, current, distance, number)
			if endTest {
				break
			}
			stats.Retreats++
			if saveSteps {
				step++
				sap.saveStep(fn, step, "Retreat", fmt.Sprintf("Retreat del nodo %d: (%d -> %d)", current, oldLabel, distance[current]), current, distance, nil)
			}
			if current != fn.Source {
				current = predecessor[current]
			}
		}
	}
	maxFlow = fn.GetMaxFlowValue()
	return
}

func (sap *BasicShortestAugmentingPath) exactDistance(fn *flownetwork.FlowNetwork) ([]int, []int) {
	distance := make([]int, fn.N)
	number := make([]int, fn.N+2)
	for i := 0; i < fn.N; i++ {
		distance[i] = fn.N
		number[fn.N]++
	}
	distance[fn.Sink] = 0
	number[0]++
	queue := []int{fn.Sink}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, edge := range fn.InStars[current] {
			previous := edge.From
			residual := edge.Capacity - edge.Flow
			if distance[previous] > distance[current]+1 && residual > 0 {
				number[distance[previous]]--
				distance[previous] = distance[current] + 1
				number[distance[previous]]++
				queue = append(queue, previous)
			}
		}
	}
	return distance, number
}

func (sap *BasicShortestAugmentingPath) findAdmissibleEdge(fn *flownetwork.FlowNetwork, node int, distance []int) int {
	for id, edge := range fn.OutStars[node] {
		residual := edge.Capacity - edge.Flow
		if residual > 0 && distance[node] == distance[edge.To]+1 {
			return id
		}
	}
	return -1
}

func (sap *BasicShortestAugmentingPath) retreat(fn *flownetwork.FlowNetwork, node int, distance []int, number []int) bool {
	minDistance := fn.N
	for _, edge := range fn.OutStars[node] {
		residual := edge.Capacity - edge.Flow
		if residual > 0 {
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
	return false
}

func (sap *BasicShortestAugmentingPath) augment(fn *flownetwork.FlowNetwork, predecessor []int) (int, []int) {
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

func (sap *BasicShortestAugmentingPath) saveStep(
	fn *flownetwork.FlowNetwork,
	step int,
	stepType string,
	description string,
	currentNode int,
	distances []int,
	path []int,
) {
	// Strutture Dati per JSON
	type NodeInfo struct {
		ID        int  `json:"id"`
		Distance  int  `json:"distance"` // La label richiesta
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
		Iteration   int        `json:"step"`
		StepType    string     `json:"step_type"` // Advances, Retreat, Augment
		Description string     `json:"description"`
		Nodes       []NodeInfo `json:"nodes"`
		Edges       []EdgeInfo `json:"edges"`
	}

	// Costruzione Snapshot
	nodes := make([]NodeInfo, fn.N)
	for i := 0; i < fn.N; i++ {
		nodes[i] = NodeInfo{
			ID:        i,
			Distance:  distances[i], // Qui salviamo la label distanza!
			IsSource:  i == fn.Source,
			IsSink:    i == fn.Sink,
			IsCurrent: i == currentNode,
		}
	}

	// Set per lookup veloce del path
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
			if edge.Capacity > 0 { // Solo archi 'reali'
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

	// Scrittura su file
	filename := fmt.Sprintf("export/graphical_steps/shortest_augmenting_path/step_%04d.json", step)
	file, _ := os.Create(filename)
	defer file.Close()
	jsonData, _ := json.MarshalIndent(snap, "", "  ")
	file.Write(jsonData)
}
*/
