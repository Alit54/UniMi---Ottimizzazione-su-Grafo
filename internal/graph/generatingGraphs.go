package graph

import "math/rand"

/*
Algoritmo 1 presentato a lezione: ha complessità O(p) ma non garantisce univocità degli archi.
*/
func CreateGraphSelectRandom(density float64, nNode int, directed bool) []int {
	if density < 0 || density > 1 {
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
