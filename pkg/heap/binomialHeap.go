/*
	Questo file è stato generato dalla IA (Claude Sonnet 4.5) per accelerare il processo del progetto. Lo scopo è verificare la correttezza dello Heap Binomiale.
*/

package heap

import "fmt"

type binomialNode struct {
	nodeID   int           // Nome del nodo
	distance float64       // Valore del nodo (per ordinamento nello Heap)
	degree   int           // Altezza dell'albero k
	parent   *binomialNode // Padre
	child    *binomialNode // Primo figlio
	sibling  *binomialNode // Fratello successivo
}

type BinomialHeap struct {
	head *binomialNode         // testa dell'albero (radice di grado più piccolo)
	min  *binomialNode         // nodo con valore più piccolo
	size int                   // Numero nodi dell'albero
	pos  map[int]*binomialNode // Mappa da nome nodo a puntatore
}

// ---------------------- //
//     METODI PUBBLICI    //
// ---------------------- //

/*
Crea un nuovo Heap binomiale da una lista di valori.
*/
func CreateBinomialHeap(values ...float64) *BinomialHeap {
	heap := &BinomialHeap{
		head: nil,
		min:  nil,
		size: 0,
		pos:  make(map[int]*binomialNode),
	}
	for i, value := range values {
		heap.Insert(i, value)
	}
	return heap
}

/*
Inserisce un nuovo nodo all'interno dello Heap.
*/
func (heap *BinomialHeap) Insert(nodeID int, distance float64) {
	newNode := &binomialNode{
		nodeID:   nodeID,
		distance: distance,
		parent:   nil,
		child:    nil,
		sibling:  nil,
	}
	tempHeap := &BinomialHeap{
		head: newNode,
		min:  newNode,
		size: 1,
		pos:  make(map[int]*binomialNode),
	}
	tempHeap.pos[nodeID] = newNode
	heap.union(tempHeap)
}

/*
ExtractMin rimuove e restituisce il nodo con valore minimo.
*/

func (heap *BinomialHeap) ExtractMin() (nodeID int, distance float64) {
	if heap.min == nil {
		return
	}
	minNode := heap.min
	heap.removeFromRootList(minNode)
	childrenHeap := &BinomialHeap{
		head: nil,
		min:  nil,
		size: 0,
		pos:  make(map[int]*binomialNode),
	}
	if minNode.child != nil {
		child := minNode.child
		var prev *binomialNode = nil
		for child != nil {
			next := child.sibling
			child.sibling = prev
			child.parent = nil
			prev = child
			child = next
			childrenHeap.size++
		}
		childrenHeap.head = prev
		current := childrenHeap.head
		for current != nil {
			childrenHeap.pos[current.nodeID] = current
			current = current.sibling
		}
	}
	delete(heap.pos, minNode.nodeID)
	heap.size -= childrenHeap.size + 1
	heap.union(childrenHeap)
	return minNode.nodeID, minNode.distance
}

/*
Diminuisce il valore di un nodo e ripristina la proprietà Heap.
*/
func (heap *BinomialHeap) DecreaseKey(nodeID int, newDistance float64) {
	node, exists := heap.pos[nodeID]
	if !exists || newDistance > node.distance {
		return
	}
	node.distance = newDistance
	heap.moveUp(node)
	node = heap.pos[nodeID]
	if node.distance < heap.min.distance {
		heap.min = node
	}

}

func (heap *BinomialHeap) IsEmpty() bool {
	return heap.size == 0
}

func (heap *BinomialHeap) Size() int {
	return heap.size
}

func (heap *BinomialHeap) PrintHeap() {
	fmt.Println("=== Binomial Heap ===")
	fmt.Println("Numero di elementi:", heap.size)
	if heap.min != nil {
		fmt.Printf("Min: nodeID=%d, valore=%f\n", heap.min.nodeID, heap.min.distance)
	}
	current := heap.head
	treeIndex := 0
	for current != nil {
		fmt.Printf("\nAlbero %d (grado %d):\n", treeIndex, current.degree)
		heap.printTree(current, 0)
		current = current.sibling
		treeIndex++
	}
}

// ---------------------- //
//     METODI PRIVATI     //
// ---------------------- //

func (heap *BinomialHeap) union(heap2 *BinomialHeap) {
	if heap2 == nil || heap2.head == nil {
		heap.updateMin()
		return
	}
	heap.head = heap.merge(heap2.head)
	var prev *binomialNode = nil
	current := heap.head
	next := current.sibling
	for next != nil {
		if current.degree != next.degree || (next.sibling != nil && next.sibling.degree == current.degree) {
			prev = current
			current = next
		} else {
			if current.distance <= next.distance {
				current.sibling = next.sibling
				heap.link(next, current)
			} else {
				if prev == nil {
					heap.head = next
				} else {
					prev.sibling = next
				}
				heap.link(current, next)
				current = next
			}
		}
		next = current.sibling
	}
	heap.updateMin()
	heap.size += heap2.size
	for name, node := range heap2.pos {
		heap.pos[name] = node
	}
}

func (heap *BinomialHeap) merge(head2 *binomialNode) *binomialNode {
	head1 := heap.head
	if head1 == nil {
		return head2
	}
	if head2 == nil {
		return head1
	}
	var head *binomialNode
	var tail *binomialNode
	if head1.degree <= head2.degree {
		head = head1
		head1 = head1.sibling
	} else {
		head = head2
		head2 = head2.sibling
	}
	tail = head
	for head1 != nil && head2 != nil {
		if head1.degree < head2.degree {
			tail.sibling = head1
			head1 = head1.sibling
		} else {
			tail.sibling = head2
			head2 = head2.sibling
		}
		tail = tail.sibling
	}
	if head1 != nil {
		tail.sibling = head1
	} else {
		tail.sibling = head2
	}
	return head
}

func (heap *BinomialHeap) link(tree1 *binomialNode, tree2 *binomialNode) {
	tree1.parent = tree2
	tree1.sibling = tree2.child
	tree2.child = tree1
	tree2.degree++
}

func (heap *BinomialHeap) updateMin() {
	if heap.head == nil {
		heap.min = nil
		return
	}
	heap.min = heap.head
	current := heap.head.sibling
	for current != nil {
		if current.distance < heap.min.distance {
			heap.min = current
		}
		current = current.sibling
	}
}

func (heap *BinomialHeap) printTree(node *binomialNode, indent int) {
	if node == nil {
		return
	}

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("nodeID=%d, valore=%f, grado=%d\n", node.nodeID, node.distance, node.degree)

	child := node.child
	for child != nil {
		heap.printTree(child, indent+1)
		child = child.sibling
	}
}

/*
Rimuove un nodo dalla lista di radici.
*/
func (heap *BinomialHeap) removeFromRootList(node *binomialNode) {
	if heap.head == node {
		heap.head = node.sibling
	} else {
		current := heap.head
		for current.sibling != node {
			current = current.sibling
		}
		current.sibling = node.sibling
	}
}

/*
Sposta un nodo verso l'alto per ripristinare le proprietà dello Heap
*/
func (heap *BinomialHeap) moveUp(node *binomialNode) {
	current := node
	parent := current.parent
	for parent != nil && current.distance < parent.distance {
		current.nodeID, parent.nodeID = parent.nodeID, current.nodeID
		current.distance, parent.distance = parent.distance, current.distance
		heap.pos[current.nodeID] = current
		heap.pos[parent.nodeID] = parent
		current = parent
		parent = current.parent
	}
}
