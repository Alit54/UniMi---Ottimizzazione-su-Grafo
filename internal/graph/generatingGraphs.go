package graph

import (
	"math/rand"
)

/*
Algoritmo 1 presentato a lezione: ha complessità O(p) ma non garantisce univocità degli archi.
*/
func CreateGraphSelectRandom(density float64, nNode int, directed bool) []int {
	numberArcs, mMax := calculateNumberArcs(density, nNode, directed)
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
	numberArcs, mMax := calculateNumberArcs(density, nNode, directed)
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

/*
Algoritmo 3 presentato a lezione: nel caso medio ha complessità O(m) ma potrebbe entrare in un loop infinito
*/
func CreateGraphDiscardRepetition(density float64, nNode int, directed bool) []int {
	numberArcs, mMax := calculateNumberArcs(density, nNode, directed)
	selected := make([]int, numberArcs)
	flag := make(map[int]bool, mMax)
	k := 0
	for k < numberArcs {
		arc := rand.Intn(mMax)
		if !flag[arc] {
			selected[k] = arc
			k++
			flag[arc] = true
		}
	}
	return selected
}

/*
Algoritmo 4 presentato a lezione: usa una hash function per rilevare archi già selezionati
*/
func CreateGraphHashFunction(density float64, nNode int, directed bool, hashDimension int) []int {
	numberArcs, mMax := calculateNumberArcs(density, nNode, directed)
	selected := make([]int, 0, numberArcs)
	flag := make(map[int]bool, hashDimension)
	hash := make([]map[int]bool, hashDimension)
	for i := 0; i < hashDimension; i++ {
		hash[i] = make(map[int]bool, numberArcs)
	}
	k := 0
	for k < numberArcs {
		arc := rand.Intn(mMax)
		hashArc := arc % hashDimension // h[arc] = arc mod hashDimension
		if !flag[hashArc] || !hash[hashArc][arc] {
			k++
			flag[hashArc] = true
			hash[hashArc][arc] = true
		}
	}
	// NOTA: complessità di questo for annidato è O(p) perché la mappa contiene solo i valori per cui hash[i] = True. Quindi l'append viene eseguita p volte.
	for i := 0; i < hashDimension; i++ {
		for arc := range hash[i] {
			selected = append(selected, arc)
		}
	}
	return selected
}

func calculateNumberArcs(density float64, nNode int, directed bool) (numberArcs int, mMax int) {
	if density < 0 || density > 0.5 {
		return 0, 0
	}
	if directed {
		mMax = nNode * (nNode - 1)
	} else {
		mMax = nNode * (nNode - 1) / 2
	}
	numberArcs = int(density * float64(mMax))
	return numberArcs, mMax
}
