package heap

import "fmt"

type BinomialNode struct {
	name    int           // Nome del nodo
	value   int           // Valore del nodo (per ordinamento nello heap)
	degree  int           // Altezza dell'albero k
	parent  *BinomialNode // Padre
	child   *BinomialNode // Primo figlio
	sibling *BinomialNode // Fratello successivo
}

type BinomialHeap struct {
	head *BinomialNode // testa dell'albero (radice di grado più piccolo)
	min  *BinomialNode // nodo con valore più piccolo
}

// ---------------------- //
//     METODI PUBBLICI    //
// ---------------------- //

/*
	Crea un nuovo heap binomiale da una lista di valori.

Complessità: O(n * log n)
*/
func CreateBinomialHeap(values ...int) *BinomialHeap {
	heap := &BinomialHeap{
		head: nil,
		min:  nil,
	}
	for i, value := range values {
		heap.Insert(i, value)
	}
	return heap
}

/*
	Inserisce un nuovo nodo all'interno dello heap.

TODO: Complessità
*/
func (heap *BinomialHeap) Insert(name int, value int) {
	newNode := &BinomialNode{
		name:    name,
		value:   value,
		parent:  nil,
		child:   nil,
		sibling: nil,
	}
	tempHeap := &BinomialHeap{
		head: newNode,
		min:  nil,
	}
	heap.union(tempHeap)
}

func (heap *BinomialHeap) PrintHeap() {
	fmt.Println("=== Binomial Heap ===")
	if heap.min != nil {
		fmt.Printf("Min: name=%d, value=%d\n", heap.min.name, heap.min.value)
	}

	current := heap.head
	treeIndex := 0
	for current != nil {
		fmt.Printf("\nTree %d (degree %d):\n", treeIndex, current.degree)
		heap.printTree(current, 0)
		current = current.sibling
		treeIndex++
	}
}

// ---------------------- //
//     METODI PRIVATI     //
// ---------------------- //

func (heap *BinomialHeap) union(heap2 *BinomialHeap) {
	if heap == nil || heap2 == nil || heap2.head == nil {
		return
	}
	heap.head = heap.merge(heap2.head)
	var prev *BinomialNode = nil
	current := heap.head
	next := current.sibling
	for next != nil {
		if current.degree != next.degree || (next.sibling != nil && next.sibling.degree == current.degree) {
			prev = current
			current = next
		} else {
			if current.value <= next.value {
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
	// heap.size += heap2.size
	/*
		for name, node := range heap2.pos {
			heap.pos[name] = node
		}
	*/
}

func (heap *BinomialHeap) merge(head2 *BinomialNode) *BinomialNode {
	head1 := heap.head
	if head1 == nil {
		return head2
	}
	if head2 == nil {
		return head1
	}
	var head *BinomialNode
	var tail *BinomialNode
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

func (heap *BinomialHeap) link(tree1 *BinomialNode, tree2 *BinomialNode) {
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
		if current.value < heap.min.value {
			heap.min = current
		}
		current = current.sibling
	}
}

func (heap *BinomialHeap) printTree(node *BinomialNode, indent int) {
	if node == nil {
		return
	}

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("name=%d, value=%d, degree=%d\n", node.name, node.value, node.degree)

	child := node.child
	for child != nil {
		heap.printTree(child, indent+1)
		child = child.sibling
	}
}
