// pkg/heap/binaryHeap_test.go
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
	min := heap.ExtractMin()
	if min.value != 1 {
		t.Errorf("Minimo atteso 1, ottenuto %d", min.value)
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

// ---------------------- //
//   EDGE CASES & ERRORI  //
// ---------------------- //

func TestDecreaseKeyNonExistentNode(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// Prova a fare DecreaseKey su un nodo che non esiste
	initialLen := heap.Len()
	heap.DecreaseKey(999, 5) // Nodo 999 non esiste

	// L'heap non dovrebbe cambiare
	if heap.Len() != initialLen {
		t.Errorf("DecreaseKey su nodo inesistente ha cambiato la lunghezza")
	}

	// Verifica che l'heap sia ancora valido
	min := heap.ExtractMin()
	if min.value != 10 {
		t.Errorf("Heap corrotto dopo DecreaseKey su nodo inesistente")
	}
}

func TestDecreaseKeyAfterExtract(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// Estrai il nodo 0
	extracted := heap.ExtractMin()
	extractedName := extracted.name

	// Prova a fare DecreaseKey sul nodo estratto
	heap.DecreaseKey(extractedName, 5)

	// Non dovrebbe crashare e l'heap dovrebbe essere ancora valido
	if heap.Len() != 2 {
		t.Errorf("Lunghezza heap errata dopo DecreaseKey su nodo estratto")
	}
}

func TestExtractMinUntilEmpty(t *testing.T) {
	heap := CreateBinaryHeap(5, 3, 7, 1)

	// Estrai tutti gli elementi
	for heap.Len() > 0 {
		heap.ExtractMin()
	}

	if heap.Len() != 0 {
		t.Errorf("Heap dovrebbe essere vuoto, lunghezza: %d", heap.Len())
	}

	if len(heap.pos) != 0 {
		t.Errorf("pos dovrebbe essere vuoto, lunghezza: %d", len(heap.pos))
	}
}

func TestNegativeValues(t *testing.T) {
	heap := CreateBinaryHeap(-5, -10, 3, -1, 0)

	expected := []int{-10, -5, -1, 0, 3}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestAllSameValues(t *testing.T) {
	heap := CreateBinaryHeap(5, 5, 5, 5, 5)

	if heap.Len() != 5 {
		t.Errorf("Lunghezza attesa 5, ottenuta %d", heap.Len())
	}

	// Tutti gli elementi dovrebbero essere 5
	for i := 0; i < 5; i++ {
		node := heap.ExtractMin()
		if node.value != 5 {
			t.Errorf("Atteso valore 5, ottenuto %d", node.value)
		}
	}
}

func TestDecreaseKeyToSameValue(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// DecreaseKey al valore attuale (non dovrebbe fare nulla)
	heap.DecreaseKey(1, 20)

	node := heap.ExtractMin()
	if node.value != 10 {
		t.Errorf("Heap modificato erroneamente da DecreaseKey allo stesso valore")
	}
}

func TestDecreaseKeyToLargerValue(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// DecreaseKey a un valore MAGGIORE (tecnicamente sbagliato, ma non dovrebbe crashare)
	heap.DecreaseKey(0, 50)

	// L'heap potrebbe non essere più valido, ma non dovrebbe crashare
	if heap.Len() != 3 {
		t.Errorf("Lunghezza heap cambiata inaspettatamente")
	}
}

func TestInsertWithNegativeName(t *testing.T) {
	heap := CreateBinaryHeap(10, 20)

	// Inserisci con name negativo
	heap.Insert(-1, 5)

	if heap.Len() != 3 {
		t.Errorf("Insert con name negativo fallito")
	}

	// Verifica che sia estratto correttamente
	min := heap.ExtractMin()
	if min.value != 5 || min.name != -1 {
		t.Errorf("Nodo con name negativo non gestito correttamente")
	}
}

func TestInsertDuplicateNames(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// Inserisci con un name già esistente (name=1 esiste già)
	heap.Insert(1, 5)

	// Questo sovrascriverà la posizione in pos[1]
	// L'heap dovrebbe comunque funzionare, anche se semanticamente è scorretto
	if heap.Len() != 4 {
		t.Errorf("Insert con name duplicato ha fallito")
	}
}

func TestVeryLargeValues(t *testing.T) {
	heap := CreateBinaryHeap(1000000, 999999, 1000001, 500000)

	expected := []int{500000, 999999, 1000000, 1000001}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestAlternatingInsertExtract(t *testing.T) {
	heap := CreateBinaryHeap(50)

	for i := 0; i < 100; i++ {
		heap.Insert(i+1, i*10)

		if i%3 == 0 && heap.Len() > 1 {
			heap.ExtractMin()
		}
	}

	// Verifica che l'heap sia ancora valido
	lastValue := -1
	for heap.Len() > 0 {
		node := heap.ExtractMin()
		if node.value < lastValue {
			t.Errorf("Proprietà heap violata: %d < %d", node.value, lastValue)
		}
		lastValue = node.value
	}
}

func TestDecreaseKeyMultipleTimes(t *testing.T) {
	heap := CreateBinaryHeap(100, 200, 300)

	// DecreaseKey più volte sullo stesso nodo
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

func TestInsertAfterExtractAll(t *testing.T) {
	heap := CreateBinaryHeap(10, 20, 30)

	// Estrai tutto
	for heap.Len() > 0 {
		heap.ExtractMin()
	}

	// Inserisci di nuovo
	heap.Insert(0, 5)
	heap.Insert(1, 3)
	heap.Insert(2, 7)

	if heap.Len() != 3 {
		t.Errorf("Reinserimento dopo svuotamento fallito")
	}

	min := heap.ExtractMin()
	if min.value != 3 {
		t.Errorf("Heap non valido dopo reinserimento: atteso 3, ottenuto %d", min.value)
	}
}

func TestZeroValues(t *testing.T) {
	heap := CreateBinaryHeap(0, 0, 0, 1, -1)

	expected := []int{-1, 0, 0, 0, 1}
	for i, exp := range expected {
		node := heap.ExtractMin()
		if node.value != exp {
			t.Errorf("Iterazione %d: atteso %d, ottenuto %d", i, exp, node.value)
		}
	}
}

func TestSingleElementOperations(t *testing.T) {
	heap := CreateBinaryHeap(42)

	// DecreaseKey su singolo elemento
	heap.DecreaseKey(0, 10)

	if heap.Len() != 1 {
		t.Error("Lunghezza cambiata dopo DecreaseKey su singolo elemento")
	}

	node := heap.ExtractMin()
	if node.value != 10 {
		t.Errorf("Atteso 10, ottenuto %d", node.value)
	}

	if heap.Len() != 0 {
		t.Error("Heap non vuoto dopo estrazione ultimo elemento")
	}
}

func TestStressDecreaseKey(t *testing.T) {
	n := 100
	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = n - i
	}

	heap := CreateBinaryHeap(values...)

	// DecreaseKey casuale su molti nodi
	for i := 0; i < n/2; i++ {
		heap.DecreaseKey(i, i/2)
	}

	// Verifica che l'heap sia ancora valido
	lastValue := -1000000
	for heap.Len() > 0 {
		node := heap.ExtractMin()
		if node.value < lastValue {
			t.Errorf("Heap property violata dopo stress test DecreaseKey")
			break
		}
		lastValue = node.value
	}
}
