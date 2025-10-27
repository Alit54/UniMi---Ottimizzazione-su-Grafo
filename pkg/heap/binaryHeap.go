package heap

import (
	"fmt"
	"strings"
)

type item struct {
	nodeID   int     // Nome del nodo
	distance float64 // Valore
}

/*
BinaryHeap Implementa uno heap binario. Uno heap binario è un albero il cui valore di un nodo padre è minore o uguale del valore di un nodo figlio.
*/
type BinaryHeap struct {
	items    []item      // Implementazione dello heap tramite array
	position map[int]int // Mappa da nodeID a indice dell'array
	size     int         // Numero di elementi nello heap
}

func CreateBinaryHeap(values ...float64) *BinaryHeap {
	n := len(values)
	heap := &BinaryHeap{
		items:    make([]item, 0, n),
		position: make(map[int]int, n),
		size:     n,
	}
	for i, value := range values {
		heap.items = append(heap.items, item{nodeID: i, distance: value})
		heap.position[i] = i
	}
	for k := n/2 - 1; k >= 0; k-- {
		heap.moveDown(k)
	}
	return heap
}

func (heap *BinaryHeap) Insert(nodeID int, distance float64) {
	heap.items = append(heap.items, item{nodeID: nodeID, distance: distance})
	heap.position[nodeID] = heap.size
	heap.moveUp(heap.size)
	heap.size++
}

func (heap *BinaryHeap) ExtractMin() (nodeID int, distance float64) {
	item := heap.items[0]
	heap.swap(0, heap.size-1)
	delete(heap.position, item.nodeID)
	heap.items = heap.items[:heap.size-1]
	heap.size--
	heap.moveDown(0)
	return item.nodeID, item.distance
}

func (heap *BinaryHeap) DecreaseKey(nodeID int, newDistance float64) {
	nodePosition, _ := heap.position[nodeID]
	heap.items[nodePosition].distance = newDistance
	heap.moveUp(nodePosition)
}

func (heap *BinaryHeap) Size() int {
	return heap.size
}

func (heap *BinaryHeap) IsEmpty() bool {
	return heap.size == 0
}

func (heap *BinaryHeap) PrintHeap() {
	fmt.Println("=== Heap binario ===")
	fmt.Println("Numero di elementi: ", heap.size)
	fmt.Println("Array:")
	for i, item := range heap.items {
		fmt.Printf("  [%d] nodeID=%d, dist=%.2f\n", i, item.nodeID, item.distance)
	}
	fmt.Println("Albero:")
	heap.printNode(0, 0)
}

func (heap *BinaryHeap) moveDown(position int) {
	stop := false
	minSon := 0
	for heap.left(position) < heap.size && !stop {
		if heap.right(position) >= heap.size || heap.items[heap.left(position)].distance < heap.items[heap.right(position)].distance {
			minSon = heap.left(position)
		} else {
			minSon = heap.right(position)
		}
		if heap.items[position].distance > heap.items[minSon].distance {
			heap.swap(position, minSon)
		} else {
			stop = true
		}
		position = minSon
	}
}

func (heap *BinaryHeap) moveUp(position int) {
	stop := false
	for position != 0 && !stop {
		parent := heap.parent(position)
		if heap.items[parent].distance > heap.items[position].distance {
			heap.swap(position, parent)
		} else {
			stop = true
		}
		position = parent
	}
}

func (heap *BinaryHeap) swap(position1 int, position2 int) {
	heap.items[position1], heap.items[position2] = heap.items[position2], heap.items[position1]
	heap.position[heap.items[position1].nodeID] = position1
	heap.position[heap.items[position2].nodeID] = position2
}

func (heap *BinaryHeap) printNode(index int, level int) {
	if index >= heap.size {
		return
	}
	item := heap.items[index]
	fmt.Printf("[%d] nodeID=%d, dist=%.2f\n", index, item.nodeID, item.distance)
	leftChild := heap.left(index)
	rightChild := heap.right(index)
	if leftChild < heap.size || rightChild < heap.size {
		prefix := strings.Repeat(" ", level+1)
		if leftChild < heap.size {
			fmt.Printf("%s", prefix)
			heap.printNode(leftChild, level+1)
		}
		if rightChild < heap.size {
			fmt.Printf("%s", prefix)
			heap.printNode(rightChild, level+1)
		}
	}
}

func (heap *BinaryHeap) left(position int) int {
	return 2*position + 1
}

func (heap *BinaryHeap) right(position int) int {
	return 2*position + 2
}

func (heap *BinaryHeap) parent(position int) int {
	return (position - 1) / 2
}
