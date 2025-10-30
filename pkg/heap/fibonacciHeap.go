package heap

import (
	"fmt"
	"math"
)

// Implementazione di uno heap di Fibonacci tramite una lista di alberi le cui radici sono disposte in una lista circolare doppiamente concatenata.
// La root list contiene i nodi radice di tutti gli alberi.
type fibonacciNode struct {
	nodeID   int            // Identificatore del nodo
	distance float64        // Priorità (distanza/costo)
	degree   int            // Numero di figli
	marked   bool           // Flag per cascading cut
	parent   *fibonacciNode // Puntatore al padre
	child    *fibonacciNode // Puntatore a un figlio (lista circolare)
	left     *fibonacciNode // Puntatore al fratello sinistro (lista circolare)
	right    *fibonacciNode // Puntatore al fratello destro (lista circolare)
}

type FibonacciHeap struct {
	min      *fibonacciNode         // Puntatore al nodo con distanza minima
	size     int                    // Numero totale di nodi
	position map[int]*fibonacciNode // nodeID -> puntatore al nodo (per DecreaseKey O(1) lookup)
}

// Crea uno Heap di Fibonacci a partire da una lista di distanze.
func CreateFibonacciHeap(distances ...float64) *FibonacciHeap {
	heap := &FibonacciHeap{
		min:      nil,
		size:     0,
		position: make(map[int]*fibonacciNode),
	}

	for i, distance := range distances {
		heap.Insert(i, distance)
	}

	return heap
}

/*
Inserisce un nuovo elemento nello heap
*/
func (heap *FibonacciHeap) Insert(nodeID int, distance float64) {
	newNode := &fibonacciNode{
		nodeID:   nodeID,
		distance: distance,
		degree:   0,
		marked:   false,
		parent:   nil,
		child:    nil,
	}
	newNode.left = newNode
	newNode.right = newNode
	if heap.min == nil {
		heap.min = newNode
	} else {
		heap.insertIntoRootList(newNode)
		if newNode.distance < heap.min.distance {
			heap.min = newNode
		}
	}
	heap.position[nodeID] = newNode
	heap.size++
}

// Rimuove e restituisce l'elemento con distanza minima.
func (heap *FibonacciHeap) ExtractMin() (nodeID int, distance float64) {
	minNode := heap.min
	if minNode == nil {
		return 0, 0.0
	}
	minNodeID := minNode.nodeID
	minDistance := minNode.distance
	if minNode.child != nil {
		child := minNode.child
		for {
			nextChild := child.right
			child.parent = nil
			child.marked = false
			heap.insertIntoRootList(child)
			child = nextChild
			if child == minNode.child {
				break
			}
		}
	}
	heap.removeFromRootList(minNode)
	delete(heap.position, minNodeID)
	heap.size--
	if heap.size == 0 {
		heap.min = nil
	} else {
		heap.merge()
	}
	return minNodeID, minDistance
}

// Diminuisce la distanza di un nodo.
func (heap *FibonacciHeap) DecreaseKey(nodeID int, newDistance float64) {
	node, exists := heap.position[nodeID]
	if !exists || newDistance > node.distance {
		return
	}
	node.distance = newDistance
	parent := node.parent
	if parent != nil && node.distance < parent.distance {
		heap.cut(node, parent)
		heap.cascadingCut(parent)
	}
	if node.distance < heap.min.distance {
		heap.min = node
	}
}

// Stampa la struttura dell'heap.
func (heap *FibonacciHeap) PrintHeap() {
	fmt.Println("=== Fibonacci Heap ===")
	fmt.Println("Numero di elementi:", heap.size)
	fmt.Printf("Min: nodeID=%d, valore=%.2f\n", heap.min.nodeID, heap.min.distance)
	if heap.min != nil {
		fmt.Println("Lista di radici:")
		current := heap.min
		treeIndex := 0
		for {
			fmt.Printf("\nAlbero %d (grado=%d):\n", treeIndex, current.degree)
			heap.printTree(current, 0)
			current = current.right
			treeIndex++
			if current == heap.min {
				break
			}
		}
	}
}

// Restituisce il numero di elementi.
func (heap *FibonacciHeap) Size() int {
	return heap.size
}

// Controlla se lo Heap è vuoto.
func (heap *FibonacciHeap) IsEmpty() bool {
	return heap.size == 0
}

// ---------------------- //
//    METODI PRIVATI      //
// ---------------------- //

// Inserisce un nodo nella root list accanto al minimo.
func (heap *FibonacciHeap) insertIntoRootList(node *fibonacciNode) {
	if heap.min == nil {
		heap.min = node
		node.left = node
		node.right = node
	} else {
		node.left = heap.min
		node.right = heap.min.right
		heap.min.right.left = node
		heap.min.right = node
	}
}

// Rimuove un nodo dalla root list.
func (heap *FibonacciHeap) removeFromRootList(node *fibonacciNode) {
	if node.right == node {
		heap.min = nil
	} else {
		node.left.right = node.right
		node.right.left = node.left
		if heap.min == node {
			heap.min = node.right
		}
	}
}

// Riduce il numero di alberi nella root list unendo alberi dello stesso grado.
func (heap *FibonacciHeap) merge() {
	if heap.min == nil {
		return
	}
	// log_φ(n) dove φ = golden ratio
	maxDegree := int(math.Floor(math.Log(float64(heap.size))/math.Log((1.0+math.Sqrt(5.0))/2.0))) + 1
	degreeTable := make([]*fibonacciNode, maxDegree+1)
	var nodes []*fibonacciNode
	current := heap.min
	for {
		nodes = append(nodes, current)
		current = current.right
		if current == heap.min {
			break
		}
	}
	for _, node := range nodes {
		degree := node.degree
		for degreeTable[degree] != nil {
			other := degreeTable[degree]
			if node.distance > other.distance {
				node, other = other, node
			}
			heap.link(other, node)
			degreeTable[degree] = nil
			degree++
		}
		degreeTable[degree] = node
	}
	heap.min = nil
	for _, node := range degreeTable {
		if node != nil {
			if heap.min == nil {
				heap.min = node
				node.left = node
				node.right = node
			} else {
				heap.insertIntoRootList(node)
				if node.distance < heap.min.distance {
					heap.min = node
				}
			}
		}
	}
}

// Collega child come figlio di parent.
func (heap *FibonacciHeap) link(child, parent *fibonacciNode) {
	child.left.right = child.right
	child.right.left = child.left
	child.parent = parent
	child.marked = false
	if parent.child == nil {
		parent.child = child
		child.left = child
		child.right = child
	} else {
		child.left = parent.child
		child.right = parent.child.right
		parent.child.right.left = child
		parent.child.right = child
	}
	parent.degree++
}

// Rimuove node dal suo parent e lo aggiunge alla root list.
func (heap *FibonacciHeap) cut(node, parent *fibonacciNode) {
	if node.right == node {
		parent.child = nil
	} else {
		node.left.right = node.right
		node.right.left = node.left
		if parent.child == node {
			parent.child = node.right
		}
	}
	parent.degree--
	node.parent = nil
	node.marked = false
	heap.insertIntoRootList(node)
}

// Esegue tagli a cascata per mantenere la struttura ottimale.
func (heap *FibonacciHeap) cascadingCut(node *fibonacciNode) {
	parent := node.parent
	if parent != nil {
		if !node.marked {
			node.marked = true
		} else {
			heap.cut(node, parent)
			heap.cascadingCut(parent)
		}
	}
}

// Stampa ricorsivamente un albero.
func (heap *FibonacciHeap) printTree(node *fibonacciNode, level int) {
	if node == nil {
		return
	}
	for i := 0; i < level; i++ {
		fmt.Print("  ")
	}
	markedStr := ""
	if node.marked {
		markedStr = " [marcato]"
	}
	fmt.Printf("nodeID=%d, valore=%.2f, grado=%d%s\n",
		node.nodeID, node.distance, node.degree, markedStr)
	if node.child != nil {
		child := node.child
		for {
			heap.printTree(child, level+1)
			child = child.right
			if child == node.child {
				break
			}
		}
	}
}
