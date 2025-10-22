package heap

import (
	"testing"
)

func TestCreateBinaryHeap(t *testing.T) {
	// Test creazione heap con valori
	heap := CreateBinaryHeap(5, 3, 8, 1, 9, 2)

	if heap.Len() != 6 {
		t.Errorf("Lunghezza attesa 6, ottenuta %d", heap.Len())
	}

	// Verifica che il minimo sia in cima (proprietà heap)
	if heap.nodes[0].value != 1 {
		t.Errorf("Minimo atteso 1, ottenuto %d", heap.nodes[0].value)
	}
}

func TestCreateBinaryHeapEmpty(t *testing.T) {
	// Test creazione heap vuoto
	heap := CreateBinaryHeap()

	if heap.Len() != 0 {
		t.Errorf("Heap vuoto dovrebbe avere lunghezza 0, ottenuta %d", heap.Len())
	}
}

func TestExtractMin(t *testing.T) {
	heap := CreateBinaryHeap(10, 5, 3, 15, 2, 8)

	// Estrai in ordine crescente
	expected := []int{2, 3, 5, 8, 10, 15}

	for i, exp := range expected {
		if heap.Len() != len(expected)-i {
			t.Errorf("Lunghezza attesa %d, ottenuta %d", len(expected)-i, heap.Len())
		}

		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}

	// Heap dovrebbe essere vuoto
	if heap.Len() != 0 {
		t.Errorf("Heap dovrebbe essere vuoto, lunghezza: %d", heap.Len())
	}
}

func TestExtractMinSingleElement(t *testing.T) {
	heap := CreateBinaryHeap(42)

	node := heap.ExtractMin()
	if node.value != 42 {
		t.Errorf("Atteso 42, ottenuto %d", node.value)
	}

	if heap.Len() != 0 {
		t.Errorf("Heap dovrebbe essere vuoto dopo ExtractMin")
	}
}

func TestInsert(t *testing.T) {
	heap := CreateBinaryHeap(5, 10, 3)

	// Inserisci un valore minore del minimo attuale
	heap.Insert(3, 1)

	if heap.Len() != 4 {
		t.Errorf("Lunghezza attesa 4, ottenuta %d", heap.Len())
	}

	// Il nuovo minimo dovrebbe essere 1
	minimum := heap.ExtractMin()
	if minimum.value != 1 {
		t.Errorf("Minimo atteso 1, ottenuto %d", minimum.value)
	}
}

func TestInsertMultiple(t *testing.T) {
	heap := CreateBinaryHeap()

	// Inserisci elementi uno alla volta
	values := []struct {
		name  int
		value int
	}{
		{0, 15},
		{1, 3},
		{2, 8},
		{3, 1},
		{4, 20},
	}

	for _, v := range values {
		heap.Insert(v.name, v.value)
	}

	if heap.Len() != 5 {
		t.Errorf("Lunghezza attesa 5, ottenuta %d", heap.Len())
	}

	// Estrai e verifica ordine
	expectedOrder := []int{1, 3, 8, 15, 20}
	for i, exp := range expectedOrder {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestDecreaseKey(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30, 40)

	// Diminuisci il valore del nodo 3 (40) a 5
	heap.DecreaseKey(3, 5)

	// Il nuovo minimo dovrebbe essere 5
	min := heap.ExtractMin()
	if min.value != 5 || min.name != 3 {
		t.Errorf("Atteso nodo 3 con valore 5, ottenuto nodo %d con valore %d",
			min.name, min.value)
	}

	// Il prossimo dovrebbe essere 10
	next := heap.ExtractMin()
	if next.value != 10 {
		t.Errorf("Atteso 10, ottenuto %d", next.value)
	}
}

func TestDecreaseKeyToMiddle(t *testing.T) {
	heap := CreateBinaryHeap(5, 15, 25, 35, 45)

	// Diminuisci 45 (nodo 4) a 20
	heap.DecreaseKey(4, 20)

	// Estrai e verifica ordine
	expected := []int{5, 15, 20, 25, 35}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestHeapPropertyAfterOperations(t *testing.T) {
	heap := CreateBinaryHeap(50, 30, 40, 10, 20)

	// Mix di operazioni
	heap.Insert(5, 5)
	heap.DecreaseKey(0, 25) // 50 -> 25
	heap.ExtractMin()       // Rimuovi 5
	heap.Insert(6, 15)

	// Verifica che la proprietà heap sia mantenuta
	// (ogni padre è minore dei figli)
	lastValue := -1
	for heap.Len() > 0 {
		node := heap.ExtractMin()
		if node.value < lastValue {
			t.Errorf("Proprietà heap violata: %d < %d", node.value, lastValue)
		}
		lastValue = node.value
	}
}

func TestDuplicateValues(t *testing.T) {
	// Test con valori duplicati
	heap := CreateBinaryHeap(5, 3, 5, 3, 5, 1, 1)

	expected := []int{1, 1, 3, 3, 5, 5, 5}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestLargeHeap(t *testing.T) {
	// Test con heap grande
	n := 1000
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = n - i // Valori in ordine decrescente
	}

	heap := CreateBinaryHeap(values...)

	if heap.Len() != n {
		t.Errorf("Lunghezza attesa %d, ottenuta %d", n, heap.Len())
	}

	// Verifica che vengano estratti in ordine crescente
	lastValue := 0
	for i := 0; i < n; i++ {
		node := heap.ExtractMin()
		if node.value <= lastValue && i > 0 {
			t.Errorf("Ordine non corretto: %d <= %d", node.value, lastValue)
		}
		lastValue = node.value
	}
}

// ---------------------- //
//      BENCHMARK         //
// ---------------------- //

func BenchmarkCreateBinaryHeap(b *testing.B) {
	values := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		values[i] = 1000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateBinaryHeap(values...)
	}
}

func BenchmarkExtractMin(b *testing.B) {
	values := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		values[i] = b.N - i
	}
	heap := CreateBinaryHeap(values...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.ExtractMin()
	}
}

func BenchmarkInsert(b *testing.B) {
	heap := CreateBinaryHeap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}
}

func BenchmarkDecreaseKey(b *testing.B) {
	n := 1000
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = i * 2
	}
	heap := CreateBinaryHeap(values...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.DecreaseKey(i%n, i%n)
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	heap := CreateBinaryHeap(100, 200, 300, 400, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, i%1000)
		if i%3 == 0 && heap.Len() > 0 {
			heap.ExtractMin()
		}
		if i%5 == 0 && heap.Len() > 0 {
			heap.DecreaseKey(i%heap.Len(), i/2)
		}
	}
}

// Test per verificare che pos sia aggiornato correttamente
func TestPositionTracking(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// Dopo la creazione, verifica che pos tenga traccia corretta
	if len(heap.pos) != heap.Len() {
		t.Errorf("pos dovrebbe avere lunghezza %d, ha %d", heap.Len(), len(heap.pos))
	}

	// Inserisci e verifica
	heap.Insert(3, 5)
	heap.DecreaseKey(0, 2)

	// ExtractMin e verifica che pos sia aggiornato
	initialLen := len(heap.pos)
	heap.ExtractMin()

	if len(heap.pos) != initialLen-1 {
		t.Errorf("Dopo ExtractMin, pos dovrebbe avere lunghezza %d, ha %d",
			initialLen-1, len(heap.pos))
	}
}
