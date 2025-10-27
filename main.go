package main

import (
	"OttimizzazioneSuGrafo/pkg/heap"
	"fmt"
)

func main() {
	binaryHeap := heap.CreateBinomialHeap(5, 15, 25, 35, 45)
	binaryHeap.PrintHeap()
	fmt.Println("---")
	binaryHeap.DecreaseKey(4, 0)
	binaryHeap.ExtractMin()
	binaryHeap.PrintHeap()
}
