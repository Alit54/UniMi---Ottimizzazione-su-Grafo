package heap

import "fmt"

type Node struct {
	name  int
	value int
}

/*
TODO: Documentazione
*/
type BinaryHeap struct {
	nodes  []Node
	length int
	pos    []int
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) Len() int {
	return heap.length
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) Left(pos int) int {
	/*if 2*pos+1 > heap.length {
		return -1
	}*/
	return 2*pos + 1
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) Right(pos int) int {
	/*if 2*pos+2 > heap.length {
		return -1
	}*/
	return 2*pos + 2
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) Parent(pos int) int {
	if pos == 0 {
		return -1
	}
	return (pos - 1) / 2
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) Swap(i, j int) {
	heap.nodes[i], heap.nodes[j] = heap.nodes[j], heap.nodes[i]
	heap.pos[i], heap.pos[j] = heap.pos[j], heap.pos[i]
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) Insert(name int, value int) {
	heap.length++
	heap.nodes = append(heap.nodes, Node{name, value})
	heap.pos[name] = heap.length
	heap.MoveUp(heap.length - 1)
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) ExtractMin() {
	// TODO: implementazione
}

/*
TODO: Documentazione
*/
func (heap BinaryHeap) DecreaseKey(name int, value int) {
	nodePosition := heap.pos[name]
	heap.nodes[nodePosition].value = value
	heap.MoveUp(nodePosition)
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) MoveUp(pos int) {
	stop := false
	for pos != 0 && !stop {
		parent := heap.Parent(pos)
		if heap.nodes[parent].value > heap.nodes[pos].value {
			heap.Swap(pos, parent)
		} else {
			stop = true
		}
		pos = parent
	}
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) MoveDown(pos int) {
	stop := false
	minSon := 0
	for heap.Left(pos) < heap.length && !stop {
		if heap.Right(pos) >= heap.length || heap.nodes[heap.Left(pos)].value < heap.nodes[heap.Right(pos)].value {
			minSon = heap.Left(pos)
		} else {
			minSon = heap.Right(pos)
		}
		if heap.nodes[pos].value > heap.nodes[minSon].value {
			heap.Swap(pos, minSon)
		} else {
			stop = true
		}
		pos = minSon
	}
}

/*
TODO: documentazione
*/
func CreateBinaryHeap(values ...int) BinaryHeap {
	n := len(values)
	heap := BinaryHeap{
		nodes:  make([]Node, 0, n),
		length: n,
		pos:    make([]int, 0, n),
	}
	for i := 0; i < n; i++ {
		heap.nodes = append(heap.nodes, Node{i, values[i]})
		heap.pos = append(heap.pos, i)
	}
	for k := n/2 - 1; k >= 0; k-- {
		heap.MoveDown(k)
	}
	return heap
}

func (heap BinaryHeap) PrintNodes() {
	for _, node := range heap.nodes {
		fmt.Println(node.value)
	}
}
