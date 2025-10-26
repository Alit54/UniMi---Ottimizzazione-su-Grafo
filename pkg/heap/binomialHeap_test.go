// pkg/heap/binomialHeap_test.go
package heap

import (
	"testing"
)

// ---------------------- //
//     TEST BASE          //
// ---------------------- //

func TestCreateBinomialHeap(t *testing.T) {
	heap := CreateBinomialHeap(5, 3, 8, 1, 9, 2)

	if heap.size != 6 {
		t.Errorf("Size attesa 6, ottenuta %d", heap.size)
	}

	if heap.min == nil || heap.min.value != 1 {
		t.Error("Minimo dovrebbe essere 1")
	}
}

func TestCreateBinomialHeapEmpty(t *testing.T) {
	heap := CreateBinomialHeap()

	if heap.size != 0 {
		t.Errorf("Heap vuoto dovrebbe avere size 0, ottenuta %d", heap.size)
	}

	if heap.min != nil {
		t.Error("Heap vuoto dovrebbe avere min = nil")
	}
}

func TestBinomialExtractMin(t *testing.T) {
	heap := CreateBinomialHeap(10, 5, 3, 15, 2, 8)

	expected := []int{2, 3, 5, 8, 10, 15}

	for i, exp := range expected {
		if heap.size != len(expected)-i {
			t.Errorf("Size attesa %d, ottenuta %d", len(expected)-i, heap.size)
		}

		node := heap.ExtractMin()
		if node == nil {
			t.Fatalf("ExtractMin restituito nil all'iterazione %d", i)
		}
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}

	if heap.size != 0 {
		t.Errorf("Heap dovrebbe essere vuoto, size: %d", heap.size)
	}
}

func TestBinomialExtractMinSingle(t *testing.T) {
	heap := CreateBinomialHeap(42)

	node := heap.ExtractMin()
	if node == nil || node.value != 42 {
		t.Error("ExtractMin fallito su heap con singolo elemento")
	}

	if heap.size != 0 {
		t.Error("Heap dovrebbe essere vuoto dopo ExtractMin")
	}
}

func TestBinomialInsert(t *testing.T) {
	heap := CreateBinomialHeap(5, 10, 3)

	heap.Insert(3, 1)

	if heap.size != 4 {
		t.Errorf("Size attesa 4, ottenuta %d", heap.size)
	}

	if heap.min.value != 1 {
		t.Errorf("Minimo atteso 1, ottenuto %d", heap.min.value)
	}

	min := heap.ExtractMin()
	if min.value != 1 {
		t.Errorf("Minimo estratto atteso 1, ottenuto %d", min.value)
	}
}

func TestBinomialInsertMultiple(t *testing.T) {
	heap := CreateBinomialHeap()

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

	if heap.size != 5 {
		t.Errorf("Size attesa 5, ottenuta %d", heap.size)
	}

	expectedOrder := []int{1, 3, 8, 15, 20}
	for i, exp := range expectedOrder {
		node := heap.ExtractMin()
		if node == nil {
			t.Fatalf("ExtractMin restituito nil all'iterazione %d", i)
		}
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialDecreaseKey(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30, 40)

	heap.DecreaseKey(3, 5)

	if heap.min.value != 5 || heap.min.name != 3 {
		t.Errorf("Min dovrebbe essere nodo 3 con valore 5, ottenuto nodo %d con valore %d",
			heap.min.name, heap.min.value)
	}

	min := heap.ExtractMin()
	if min.value != 5 || min.name != 3 {
		t.Errorf("Atteso nodo 3 con valore 5, ottenuto nodo %d con valore %d",
			min.name, min.value)
	}

	next := heap.ExtractMin()
	if next.value != 10 {
		t.Errorf("Atteso 10, ottenuto %d", next.value)
	}
}

func TestBinomialDecreaseKeyToMiddle(t *testing.T) {
	heap := CreateBinomialHeap(5, 15, 25, 35, 45)

	heap.DecreaseKey(4, 20)

	expected := []int{5, 15, 20, 25, 35}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node == nil {
			t.Fatalf("ExtractMin restituito nil all'iterazione %d", i)
		}
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialMinTracking(t *testing.T) {
	heap := CreateBinomialHeap(10, 5, 20, 3, 15)

	if heap.min == nil || heap.min.value != 3 {
		t.Error("Min dovrebbe essere 3 dopo creazione")
	}

	// Min non cambia senza estrazione
	if heap.min.value != 3 {
		t.Error("Min non dovrebbe cambiare")
	}

	heap.ExtractMin()
	if heap.min.value != 5 {
		t.Errorf("Min dovrebbe essere 5 dopo estrazione, ottenuto %d", heap.min.value)
	}
}

// ---------------------- //
//   EDGE CASES & ERRORI  //
// ---------------------- //

func TestBinomialDecreaseKeyNonExistent(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30)

	initialSize := heap.size
	heap.DecreaseKey(999, 5)

	if heap.size != initialSize {
		t.Error("DecreaseKey su nodo inesistente ha modificato l'heap")
	}
}

func TestBinomialDecreaseKeyToLarger(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30)

	heap.DecreaseKey(0, 50) // Tentativo di aumentare

	// Dovrebbe essere ignorato
	min := heap.ExtractMin()
	if min.value != 10 {
		t.Error("DecreaseKey con valore maggiore dovrebbe essere ignorato")
	}
}

func TestBinomialDecreaseKeyToSame(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30)

	heap.DecreaseKey(1, 20) // Stesso valore

	// Non dovrebbe crashare
	if heap.size != 3 {
		t.Error("DecreaseKey allo stesso valore ha modificato size")
	}
}

func TestBinomialExtractMinEmpty(t *testing.T) {
	heap := CreateBinomialHeap()

	node := heap.ExtractMin()
	if node != nil {
		t.Error("ExtractMin su heap vuoto dovrebbe restituire nil")
	}
}

func TestBinomialNegativeValues(t *testing.T) {
	heap := CreateBinomialHeap(-5, -10, 3, -1, 0)

	expected := []int{-10, -5, -1, 0, 3}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialDuplicateValues(t *testing.T) {
	heap := CreateBinomialHeap(5, 3, 5, 3, 5, 1, 1)

	expected := []int{1, 1, 3, 3, 5, 5, 5}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialAllSameValues(t *testing.T) {
	heap := CreateBinomialHeap(7, 7, 7, 7)

	for i := 0; i < 4; i++ {
		node := heap.ExtractMin()
		if node.value != 7 {
			t.Errorf("Atteso 7, ottenuto %d", node.value)
		}
	}

	if heap.size != 0 {
		t.Error("Heap dovrebbe essere vuoto")
	}
}

func TestBinomialLargeHeap(t *testing.T) {
	n := 1000
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = n - i
	}

	heap := CreateBinomialHeap(values...)

	if heap.size != n {
		t.Errorf("Size attesa %d, ottenuta %d", n, heap.size)
	}

	lastValue := 0
	for i := 0; i < n; i++ {
		node := heap.ExtractMin()
		if node.value < lastValue {
			t.Errorf("Ordine non corretto: %d < %d", node.value, lastValue)
		}
		lastValue = node.value
	}
}

func TestBinomialAlternatingOps(t *testing.T) {
	heap := CreateBinomialHeap(50)

	for i := 0; i < 100; i++ {
		heap.Insert(i+1, i*10)

		if i%3 == 0 && heap.size > 1 {
			heap.ExtractMin()
		}
	}

	lastValue := -1
	for heap.size > 0 {
		node := heap.ExtractMin()
		if node.value < lastValue {
			t.Errorf("Proprietà heap violata: %d < %d", node.value, lastValue)
		}
		lastValue = node.value
	}
}

func TestBinomialDecreaseKeyMultiple(t *testing.T) {
	heap := CreateBinomialHeap(100, 200, 300)

	heap.DecreaseKey(0, 90)
	heap.DecreaseKey(0, 80)
	heap.DecreaseKey(0, 70)
	heap.DecreaseKey(0, 60)

	min := heap.ExtractMin()
	if min.value != 60 || min.name != 0 {
		t.Errorf("DecreaseKey multipli fallito: atteso name=0 value=60, ottenuto name=%d value=%d",
			min.name, min.value)
	}
}

func TestBinomialInsertAfterExtractAll(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30)

	for heap.size > 0 {
		heap.ExtractMin()
	}

	heap.Insert(0, 5)
	heap.Insert(1, 3)
	heap.Insert(2, 7)

	if heap.size != 3 {
		t.Error("Reinserimento dopo svuotamento fallito")
	}

	min := heap.ExtractMin()
	if min.value != 3 {
		t.Errorf("Heap non valido dopo reinserimento: atteso 3, ottenuto %d", min.value)
	}
}

func TestBinomialZeroValues(t *testing.T) {
	heap := CreateBinomialHeap(0, 0, 0, 1, -1)

	expected := []int{-1, 0, 0, 0, 1}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialSingleElement(t *testing.T) {
	heap := CreateBinomialHeap(42)

	heap.DecreaseKey(0, 10)

	if heap.size != 1 {
		t.Error("Size modificata dopo DecreaseKey su singolo elemento")
	}

	node := heap.ExtractMin()
	if node.value != 10 {
		t.Errorf("Atteso 10, ottenuto %d", node.value)
	}

	if heap.size != 0 {
		t.Error("Heap non vuoto dopo estrazione ultimo elemento")
	}
}

func TestBinomialPowerOfTwo(t *testing.T) {
	// Test con potenze di 2 (struttura ottimale per binomial heap)
	for n := 1; n <= 64; n *= 2 {
		heap := CreateBinomialHeap()
		for i := 0; i < n; i++ {
			heap.Insert(i, n-i)
		}

		if heap.size != n {
			t.Errorf("Per n=%d, size attesa %d, ottenuta %d", n, n, heap.size)
		}

		for i := 1; i <= n; i++ {
			node := heap.ExtractMin()
			if node.value != i {
				t.Errorf("Per n=%d, iterazione %d: atteso %d, ottenuto %d", n, i, i, node.value)
			}
		}
	}
}

func TestBinomialInsertDuplicateNames(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30)

	// Inserisci con name già esistente
	heap.Insert(1, 5)

	// La pos map verrà sovrascritta
	if heap.size != 4 {
		t.Error("Insert con name duplicato fallito")
	}
}

func TestBinomialVeryLargeValues(t *testing.T) {
	heap := CreateBinomialHeap(1000000, 999999, 1000001, 500000)

	expected := []int{500000, 999999, 1000000, 1000001}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestBinomialExtractUntilEmpty(t *testing.T) {
	heap := CreateBinomialHeap(5, 3, 7, 1, 9)

	for heap.size > 0 {
		heap.ExtractMin()
	}

	if heap.size != 0 {
		t.Errorf("Heap dovrebbe essere vuoto, size: %d", heap.size)
	}

	if heap.min != nil {
		t.Error("Min dovrebbe essere nil dopo svuotamento")
	}
}

func TestBinomialPosMapConsistency(t *testing.T) {
	heap := CreateBinomialHeap(10, 20, 30, 40, 50)

	// Verifica che pos map sia consistente
	if len(heap.pos) != heap.size {
		t.Errorf("pos map size (%d) != heap size (%d)", len(heap.pos), heap.size)
	}

	// Estrai e verifica
	heap.ExtractMin()
	if len(heap.pos) != heap.size {
		t.Errorf("Dopo ExtractMin, pos map size (%d) != heap size (%d)", len(heap.pos), heap.size)
	}

	// DecreaseKey e verifica
	heap.DecreaseKey(1, 5)
	if len(heap.pos) != heap.size {
		t.Errorf("Dopo DecreaseKey, pos map size (%d) != heap size (%d)", len(heap.pos), heap.size)
	}
}

// ---------------------- //
//      BENCHMARK         //
// ---------------------- //

func BenchmarkBinomialCreate(b *testing.B) {
	values := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		values[i] = 1000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateBinomialHeap(values...)
	}
}

func BenchmarkBinomialInsert(b *testing.B) {
	heap := CreateBinomialHeap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}
}

func BenchmarkBinomialExtractMin(b *testing.B) {
	values := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		values[i] = b.N - i
	}
	heap := CreateBinomialHeap(values...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.ExtractMin()
	}
}

func BenchmarkBinomialDecreaseKey(b *testing.B) {
	n := 1000
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = i * 2
	}
	heap := CreateBinomialHeap(values...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.DecreaseKey(i%n, i%n)
	}
}

func BenchmarkBinomialMixed(b *testing.B) {
	heap := CreateBinomialHeap(100, 200, 300, 400, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, i%1000)
		if i%3 == 0 && heap.size > 0 {
			heap.ExtractMin()
		}
		if i%5 == 0 && heap.size > 0 {
			heap.DecreaseKey(i%heap.size, i/2)
		}
	}
}

func BenchmarkBinomialGetMin(b *testing.B) {
	heap := CreateBinomialHeap(100, 50, 200, 25, 300)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.min
	}
}
