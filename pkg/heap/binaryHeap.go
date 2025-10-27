package heap

import (
	"fmt"
	"strings"
)

type Node struct {
	name  int // Nome del nodo
	value int // Valore del nodo (per ordinamento nello heap)
}

/*
	Implementazione di uno Heap Binario.

Uno heap binario è un albero dove il valore di un nodo padre è sempre minore o al più uguale di entrambi i valori dei nodi figli. Inoltre, il massimo numero di nodi figlio è 2.
*/
type BinaryHeap struct {
	nodes  []Node      // nodi dello heap
	length int         // lunghezza dello heap
	pos    map[int]int // posizioni dei nodi nello heap. La mappa associa il nome del nodo alla sua posizione nello heap
}

// ---------------------- //
//     METODI PUBBLICI    //
// ---------------------- //

/*
	Creazione di uno heap binario sistemando una volta sola lo heap.

Complessità: O(n)
*/
func CreateBinaryHeap(values ...int) *BinaryHeap {
	n := len(values)
	heap := &BinaryHeap{
		nodes:  make([]Node, 0, n),
		length: n,
		pos:    make(map[int]int, n),
	}
	for i, value := range values {
		heap.nodes = append(heap.nodes, Node{i, value})
		heap.pos[i] = i
	}
	for k := n/2 - 1; k >= 0; k-- {
		heap.moveDown(k)
	}
	return heap
}

/*
	Restituisce il numero di nodi nello heap.

Complessità: O(1)
*/
func (heap *BinaryHeap) Len() int {
	return heap.length
}

/*
	Inserisce un nuovo nodo nello heap.

Complessità: O(log n)
*/
func (heap *BinaryHeap) Insert(name int, value int) {
	heap.length++
	heap.nodes = append(heap.nodes, Node{name, value})
	heap.pos[name] = heap.length - 1
	heap.moveUp(heap.length - 1)
}

/*
	Estrae il nodo con valore minore.

Complessità: O(log n)
*/
func (heap *BinaryHeap) ExtractMin() Node {
	node := heap.nodes[0]
	heap.swap(0, heap.length-1)
	delete(heap.pos, node.name)
	heap.nodes = heap.nodes[:len(heap.nodes)-1]
	heap.length--
	heap.moveDown(0)
	return node
}

/*
	Decrementa il valore di un nodo presente nello heap.

Complessità: O(log n)
*/
func (heap *BinaryHeap) DecreaseKey(name int, newValue int) {
	nodePosition, exists := heap.pos[name]
	if !exists || newValue > heap.nodes[nodePosition].value {
		return
	}
	heap.nodes[nodePosition].value = newValue
	heap.moveUp(nodePosition)
}

/*
	Stampa lo heap.

Complessità: O(n)
*/
func (heap *BinaryHeap) PrintHeap() {
	heap.printNode(heap.nodes[0], 0)
}

// ---------------------- //
//     METODI PRIVATI     //
// ---------------------- //

/*
	Restituisce l'indice del figlio di sinistra

Complessità: O(1)
*/
func (heap *BinaryHeap) left(pos int) int {
	/*if 2*pos+1 > heap.length {
		return -1
	}*/
	return 2*pos + 1
}

/*
	Restituisce l'indice del figlio di destra

Complessità: O(1)
*/
func (heap *BinaryHeap) right(pos int) int {
	/*if 2*pos+2 > heap.length {
		return -1
	}*/
	return 2*pos + 2
}

/*
	Restituisce l'indice del padre

Complessità: O(1)
*/
func (heap *BinaryHeap) parent(pos int) int {
	if pos == 0 {
		return -1
	}
	return (pos - 1) / 2
}

/*
	Scambia due nodi, prendendo come input le posizioni nello heap dei due nodi.

Complessità: O(1)
*/
func (heap *BinaryHeap) swap(i, j int) {
	heap.nodes[i], heap.nodes[j] = heap.nodes[j], heap.nodes[i]
	heap.pos[heap.nodes[i].name] = i
	heap.pos[heap.nodes[j].name] = j
}

/*
	Sposta verso l'alto un nodo con valore non corretto per le proprietà dello heap.

Complessità: O(log n)
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
	Sposta verso il basso un nodo con valore non corretto per le proprietà dello heap.

Complessità: O(log n)
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

/*
	stampa un nodo e tutti i suoi figli ricorsivamente

Complessità: O(n)
*/
func (heap *BinaryHeap) printNode(node Node, indent int) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Printf("Nome: %d Valore: %d\n", node.name, node.value)
	if heap.left(heap.pos[node.name]) < heap.length {
		leftSon := heap.nodes[heap.left(heap.pos[node.name])]
		heap.printNode(leftSon, indent+1)
	}
	if heap.right(heap.pos[node.name]) < heap.length {
		rightSon := heap.nodes[heap.right(heap.pos[node.name])]
		heap.printNode(rightSon, indent+1)
	}
}
