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

// ---------------------- //
//     METODI PUBBLICI    //
// ---------------------- //

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
		heap.moveDown(k)
	}
	return heap
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) Len() int {
	return heap.length
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) Insert(name int, value int) {
	heap.length++
	heap.nodes = append(heap.nodes, Node{name, value})
	heap.pos[name] = heap.length
	heap.moveUp(heap.length - 1)
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) ExtractMin() Node {
	node := heap.nodes[0]
	heap.swap(0, heap.length-1)
	heap.pos = heap.pos[:heap.length-1]
	heap.length--
	heap.moveDown(0)
	return node
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) DecreaseKey(name int, value int) {
	nodePosition := heap.pos[name]
	heap.nodes[nodePosition].value = value
	heap.moveUp(nodePosition)
}

/*
TODO: documentazione
*/
func (heap *BinaryHeap) PrintNodes() {
	for _, node := range heap.nodes {
		fmt.Println(node.value)
	}
}

// ---------------------- //
//     METODI PRIVATI     //
// ---------------------- //

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) left(pos int) int {
	/*if 2*pos+1 > heap.length {
		return -1
	}*/
	return 2*pos + 1
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) right(pos int) int {
	/*if 2*pos+2 > heap.length {
		return -1
	}*/
	return 2*pos + 2
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) parent(pos int) int {
	if pos == 0 {
		return -1
	}
	return (pos - 1) / 2
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) swap(i, j int) {
	heap.nodes[i], heap.nodes[j] = heap.nodes[j], heap.nodes[i]
	heap.pos[i], heap.pos[j] = heap.pos[j], heap.pos[i]
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) moveUp(pos int) {
	stop := false
	for pos != 0 && !stop {
		parent := heap.parent(pos)
		if heap.nodes[parent].value > heap.nodes[pos].value {
			heap.swap(pos, parent)
		} else {
			stop = true
		}
		pos = parent
	}
}

/*
TODO: Documentazione
*/
func (heap *BinaryHeap) moveDown(pos int) {
	stop := false
	minSon := 0
	for heap.left(pos) < heap.length && !stop {
		if heap.right(pos) >= heap.length || heap.nodes[heap.left(pos)].value < heap.nodes[heap.right(pos)].value {
			minSon = heap.left(pos)
		} else {
			minSon = heap.right(pos)
		}
		if heap.nodes[pos].value > heap.nodes[minSon].value {
			heap.swap(pos, minSon)
		} else {
			stop = true
		}
		pos = minSon
	}
}
