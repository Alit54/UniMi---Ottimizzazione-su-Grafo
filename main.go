package main

import (
	"OttimizzazioneSuGrafo/pkg/heap"
)

func main() {
	binaryHeap := heap.CreateBinomialHeap(4, 1, 5, 2, 3, 0, 7)
	binaryHeap.PrintHeap()
}
