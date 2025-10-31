package main

import (
	"OttimizzazioneSuGrafo/internal/graph"
	"fmt"
)

func main() {
	fmt.Println("First method")
	a := graph.CreateGraphSelectRandom(0.5, 4, true)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%2d ", a[i])
	}

	fmt.Println("\nSecond method")
	b := graph.CreateGraphSelectRandomSubset(0.5, 4, true)
	for i := 0; i < len(b); i++ {
		fmt.Printf("%2d ", b[i])
	}
	fmt.Println("\nThird method")
	c := graph.CreateGraphDiscardRepetition(0.5, 4, true)
	for i := 0; i < len(c); i++ {
		fmt.Printf("%2d ", c[i])
	}
	fmt.Println("\nFourth method")
	d := graph.CreateGraphHashFunction(0.5, 4, true, 2)
	for i := 0; i < len(d); i++ {
		fmt.Printf("%2d ", d[i])
	}
	fmt.Println("\nFifth method")
	e := graph.CreateGraphSelectSuitablyElements(0.5, 4, true)
	for i := 0; i < len(e); i++ {
		fmt.Printf("%2d ", e[i])
	}
}
