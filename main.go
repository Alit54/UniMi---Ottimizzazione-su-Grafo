package main

import (
	"OttimizzazioneSuGrafo/pkg/heap"
	"fmt"
)

func main() {
	heap := heap.CreateBinomialHeap()

	heap.PrintHeap()
	fmt.Println("---")

	// Alterna Insert ed ExtractMin
	heap.Insert(1, 50.0)
	heap.Insert(2, 30.0)
	heap.Insert(3, 70.0)

	heap.PrintHeap()
	fmt.Println("---")

	heap.ExtractMin() // Estrai 30.0

	heap.PrintHeap()
	fmt.Println("---")

	heap.Insert(4, 20.0)
	heap.Insert(5, 40.0)

	heap.PrintHeap()
	fmt.Println("---")

	heap.ExtractMin() // Estrai 20.0

	heap.PrintHeap()
	fmt.Println("---")

	heap.DecreaseKey(3, 25.0)

	heap.PrintHeap()
}
