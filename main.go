package main

import (
	"OttimizzazioneSuGrafo/pkg/heap"
)

func main() {
	heap := heap.CreateFibonacciHeap(3, 6, 7, 4, 1, 9, 12, 52, 324, 143, 343)
	heap.PrintHeap()
	heap.ExtractMin()
	heap.PrintHeap()
}
