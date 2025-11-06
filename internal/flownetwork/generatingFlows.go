package flownetwork

import "math/rand"

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
