package flownetwork

import (
	"math/rand"
)

// GenerateRandomArcs genera un FlowNetwork casuale usando l'algoritmo SelectSuitablyElements
func (fn *FlowNetwork) GenerateRandomArcs(density float64, minCap int, maxCap int) {
	if density < 0 || density > 1 {
		panic("density deve essere tra 0 e 1")
	}
	if minCap > maxCap {
		panic("minCap deve essere minore o uguale a maxCap")
	}
	n := fn.N
	mMax := n * (n - 1)
	numberArcs := int(density * float64(mMax))
	selectedIDs := selectSuitablyElements(numberArcs, mMax)
	for _, arcID := range selectedIDs {
		from, to := arcIDToNodes(arcID, n)
		capacity := rand.Intn(maxCap-minCap+1) + minCap
		fn.AddEdge(from, to, capacity)
	}
	fn.ensureConnectivity(minCap, maxCap)
}

// selectSuitablyElements implementa l'algoritmo 5 presentato a lezione per la generazione di archi casuali in un grafo
func selectSuitablyElements(numberArcs, mMax int) []int {
	selected := make([]int, numberArcs)
	elements := make([]int, mMax)
	for i := 0; i < mMax; i++ {
		elements[i] = i
	}
	for k := 0; k < numberArcs; k++ {
		arc := rand.Intn(mMax-k) + k
		selected[k] = elements[arc]
		elements[arc] = elements[k]
	}
	return selected
}

// arcIDToNodes converte un arcID numerico in una coppia (from, to)
func arcIDToNodes(arcID, n int) (from, to int) {
	from = arcID / (n - 1)
	to = arcID % (n - 1)
	if arcID >= from*n {
		to++
	}
	return from, to
}

// isReachable controlla, tramite una BFS, se il nodo sink è raggiungibile da source, per permettere almeno una soluzione ammissibile nel problema di maxflow.
func (fn *FlowNetwork) isReachable(source, sink int) bool {
	visited := make(map[int]bool, fn.N)
	queue := []int{source}
	visited[source] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == sink {
			return true
		}
		for _, edge := range fn.OutStars[current] {
			residual := edge.Capacity - edge.Flow
			if !visited[edge.To] && residual > 0 {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
	}
	return false
}

func (fn *FlowNetwork) getMostFarFromSource(source int) int {
	visited := make([]bool, fn.N)
	lastExtracted := source
	queue := []int{source}
	visited[source] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, edge := range fn.OutStars[current] {
			residual := edge.Capacity
			if !visited[edge.To] && residual > 0 {
				visited[edge.To] = true
				queue = append(queue, edge.To)
			}
		}
		lastExtracted = current
	}
	return lastExtracted
}

func (fn *FlowNetwork) getMostFarToSink(sink int) int {
	visited := make([]bool, fn.N)
	queue := []int{sink}
	lastExtracted := sink
	visited[sink] = true
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for previous := 0; previous < fn.N; previous++ {
			if visited[previous] {
				continue
			}
			for _, edge := range fn.OutStars[previous] {
				if edge.To == current && edge.Capacity > 0 {
					visited[previous] = true
					queue = append(queue, previous)
					break
				}
			}
		}
		lastExtracted = current
	}
	return lastExtracted
}

func (fn *FlowNetwork) ensureConnectivity(minCap int, maxCap int) {
	if fn.isReachable(fn.Source, fn.Sink) {
		return
	}
	mostFarFromSource := fn.getMostFarFromSource(fn.Source)
	mostFarToSink := fn.getMostFarToSink(fn.Sink)
	randomCapacity := rand.Intn(maxCap-minCap+1) + minCap
	fn.AddEdge(mostFarFromSource, mostFarToSink, randomCapacity)
}
