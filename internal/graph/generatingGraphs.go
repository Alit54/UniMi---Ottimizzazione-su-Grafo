package graph

import "math/rand"

/*
Algoritmo 1 presentato a lezione: ha complessità O(p) ma non garantisce univocità degli archi.
*/
func CreateGraphSelectRandom(density float64, nNode int, directed bool) []int {
	if density < 0 || density > 0.5 {
		return nil
	}
	var mMax int
	if directed {
		mMax = nNode * (nNode - 1)
	} else {
		mMax = nNode * (nNode - 1) / 2
	}
	numberArcs := int(density * float64(mMax))
	selected := make([]int, numberArcs)
	for i := 0; i < numberArcs; i++ {
		selected[i] = rand.Intn(mMax)
	}
	return selected
}

/*
Algoritmo 2 presentato a lezione: ha complessità O(p) ma non genera gli archi in modo indipendente
*/
func CreateGraphSelectRandomSubset(density float64, nNode int, directed bool) []int {
	if density < 0 || density > 0.5 {
		return nil
	}
	var mMax int
	if directed {
		mMax = nNode * (nNode - 1)
	} else {
		mMax = nNode * (nNode - 1) / 2
	}
	numberArcs := int(density * float64(mMax))
	selected := make([]int, numberArcs)
	arc := rand.Intn(mMax)
	for i := 0; i < numberArcs; i++ {
		selected[i] = arc
		if arc == mMax {
			arc = 0
		} else {
			arc++
		}
	}
	return selected
}
