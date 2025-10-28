/*
	Questo file è stato generato dalla IA (Claude Sonnet 4.5) per accelerare il processo del progetto. Lo scopo è verificare la correttezza dello Heap Binario.
*/

package heap

import (
	"math/rand"
	"testing"
)

// ============================================ //
//          TEST COSTRUTTORE                    //
// ============================================ //

func TestCreateBinaryHeap(t *testing.T) {
	values := []float64{50.0, 30.0, 70.0, 10.0, 40.0}
	heap := CreateBinaryHeap(values...)

	if heap.Size() != 5 {
		t.Errorf("Size attesa 5, ottenuta %d", heap.Size())
	}

	// Verifica proprietà min-heap: radice deve essere il minimo
	if heap.items[0].distance != 10.0 {
		t.Errorf("Radice dovrebbe essere 10.0, ottenuto %.1f", heap.items[0].distance)
	}
}

func TestCreateBinaryHeapEmpty(t *testing.T) {
	heap := CreateBinaryHeap()

	if heap.Size() != 0 {
		t.Error("Heap vuoto dovrebbe avere size 0")
	}

	if !heap.IsEmpty() {
		t.Error("Heap vuoto dovrebbe risultare isEmpty")
	}
}

func TestCreateBinaryHeapSingle(t *testing.T) {
	heap := CreateBinaryHeap(42.0)

	if heap.Size() != 1 {
		t.Error("Heap con singolo elemento dovrebbe avere size 1")
	}

	nodeID, dist := heap.ExtractMin()
	if nodeID != 0 || dist != 42.0 {
		t.Errorf("Atteso (0, 42.0), ottenuto (%d, %.1f)", nodeID, dist)
	}
}

func TestCreateBinaryHeapOrdered(t *testing.T) {
	// Test con valori già ordinati (caso edge)
	heap := CreateBinaryHeap(10.0, 20.0, 30.0, 40.0, 50.0)

	if heap.items[0].distance != 10.0 {
		t.Error("Heap dovrebbe mantenere 10.0 come minimo")
	}
}

func TestCreateBinaryHeapReverse(t *testing.T) {
	// Test con valori in ordine inverso
	heap := CreateBinaryHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	if heap.items[0].distance != 10.0 {
		t.Error("Heap dovrebbe riorganizzare con 10.0 come minimo")
	}
}

// ============================================ //
//          TEST INSERT                         //
// ============================================ //

func TestInsertSingle(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(1, 10.0)

	if heap.Size() != 1 {
		t.Errorf("Size attesa 1, ottenuta %d", heap.Size())
	}

	if heap.items[0].nodeID != 1 || heap.items[0].distance != 10.0 {
		t.Error("Elemento inserito non corretto")
	}
}

func TestInsertMultiple(t *testing.T) {
	heap := CreateBinaryHeap()

	testData := []struct {
		nodeID   int
		distance float64
	}{
		{5, 50.0},
		{3, 30.0},
		{7, 70.0},
		{1, 10.0},
		{9, 90.0},
	}

	for _, data := range testData {
		heap.Insert(data.nodeID, data.distance)
	}

	if heap.Size() != 5 {
		t.Errorf("Size attesa 5, ottenuta %d", heap.Size())
	}

	// Il minimo dovrebbe essere nodeID=1 con distance=10.0
	if heap.items[0].nodeID != 1 || heap.items[0].distance != 10.0 {
		t.Errorf("Minimo atteso (1, 10.0), ottenuto (%d, %.1f)",
			heap.items[0].nodeID, heap.items[0].distance)
	}
}

func TestInsertMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinaryHeap()

	// Insert elementi casuali
	values := []float64{50.0, 30.0, 70.0, 20.0, 60.0, 10.0, 80.0}
	for i, val := range values {
		heap.Insert(i, val)
	}

	// Verifica proprietà heap per ogni nodo
	for i := 0; i < heap.Size(); i++ {
		leftChild := 2*i + 1
		rightChild := 2*i + 2

		if leftChild < heap.Size() {
			if heap.items[i].distance > heap.items[leftChild].distance {
				t.Errorf("Violazione heap property: parent[%d]=%.1f > leftChild[%d]=%.1f",
					i, heap.items[i].distance, leftChild, heap.items[leftChild].distance)
			}
		}

		if rightChild < heap.Size() {
			if heap.items[i].distance > heap.items[rightChild].distance {
				t.Errorf("Violazione heap property: parent[%d]=%.1f > rightChild[%d]=%.1f",
					i, heap.items[i].distance, rightChild, heap.items[rightChild].distance)
			}
		}
	}
}

func TestInsertUpdatePosition(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 200.0)
	heap.Insert(30, 300.0)

	// Verifica che position map sia aggiornata
	for nodeID := range heap.position {
		idx := heap.position[nodeID]
		if heap.items[idx].nodeID != nodeID {
			t.Errorf("Position map inconsistente per nodeID=%d", nodeID)
		}
	}
}

// ============================================ //
//          TEST EXTRACTMIN                     //
// ============================================ //

func TestExtractMinOrder(t *testing.T) {
	heap := CreateBinaryHeap(100.0, 50.0, 150.0, 30.0, 70.0, 10.0)

	// ExtractMin dovrebbe restituire valori in ordine crescente
	expectedDistances := []float64{10.0, 30.0, 50.0, 70.0, 100.0, 150.0}

	for i, expectedDist := range expectedDistances {
		_, dist := heap.ExtractMin()

		if dist != expectedDist {
			t.Errorf("Iterazione %d: distanza attesa %.1f, ottenuta %.1f",
				i, expectedDist, dist)
		}
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto dopo aver estratto tutti gli elementi")
	}
}

func TestExtractMinCorrectNodeIDs(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 50.0)
	heap.Insert(30, 150.0)

	// Primo ExtractMin: nodeID=20 (dist=50.0)
	nodeID, dist := heap.ExtractMin()
	if nodeID != 20 || dist != 50.0 {
		t.Errorf("Primo ExtractMin: atteso (20, 50.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}

	// Secondo ExtractMin: nodeID=10 (dist=100.0)
	nodeID, dist = heap.ExtractMin()
	if nodeID != 10 || dist != 100.0 {
		t.Errorf("Secondo ExtractMin: atteso (10, 100.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}

	// Terzo ExtractMin: nodeID=30 (dist=150.0)
	nodeID, dist = heap.ExtractMin()
	if nodeID != 30 || dist != 150.0 {
		t.Errorf("Terzo ExtractMin: atteso (30, 150.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestExtractMinUpdatesPosition(t *testing.T) {
	heap := CreateBinaryHeap(50.0, 30.0, 70.0, 10.0, 40.0)

	initialSize := heap.Size()
	minNodeID := heap.items[0].nodeID

	heap.ExtractMin()

	// Verifica che nodeID estratto non sia più in position map
	if _, exists := heap.position[minNodeID]; exists {
		t.Error("NodeID estratto dovrebbe essere rimosso dalla position map")
	}

	// Verifica che size sia diminuita
	if heap.Size() != initialSize-1 {
		t.Error("Size dovrebbe diminuire dopo ExtractMin")
	}

	// Verifica consistenza position map
	if len(heap.position) != heap.Size() {
		t.Errorf("Position map size (%d) != heap size (%d)",
			len(heap.position), heap.Size())
	}
}

func TestExtractMinSingle(t *testing.T) {
	heap := CreateBinaryHeap(42.0)

	nodeID, dist := heap.ExtractMin()

	if nodeID != 0 || dist != 42.0 {
		t.Errorf("ExtractMin su singolo elemento: atteso (0, 42.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto dopo ExtractMin su singolo elemento")
	}
}

func TestExtractMinMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinaryHeap()

	// Insert molti elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64((i*7)%20))
	}

	// Estrai alcuni elementi
	for i := 0; i < 10; i++ {
		heap.ExtractMin()
	}

	// Verifica che l'heap mantenga la proprietà
	for i := 0; i < heap.Size(); i++ {
		leftChild := 2*i + 1
		rightChild := 2*i + 2

		if leftChild < heap.Size() {
			if heap.items[i].distance > heap.items[leftChild].distance {
				t.Error("Proprietà heap violata dopo ExtractMin multipli")
			}
		}

		if rightChild < heap.Size() {
			if heap.items[i].distance > heap.items[rightChild].distance {
				t.Error("Proprietà heap violata dopo ExtractMin multipli")
			}
		}
	}
}

// ============================================ //
//          TEST DECREASEKEY                    //
// ============================================ //

func TestDecreaseKeyBasic(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(1, 100.0)
	heap.Insert(2, 200.0)
	heap.Insert(3, 300.0)

	// Diminuisci distanza di nodeID=3
	heap.DecreaseKey(3, 50.0)

	// Ora nodeID=3 dovrebbe essere il minimo
	nodeID, dist := heap.ExtractMin()
	if nodeID != 3 || dist != 50.0 {
		t.Errorf("Dopo DecreaseKey, atteso (3, 50.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestDecreaseKeyToMinimum(t *testing.T) {
	heap := CreateBinaryHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	// Diminuisci l'ultimo nodo (nodeID=4, dist=10.0) a 5.0
	heap.DecreaseKey(4, 5.0)

	nodeID, dist := heap.ExtractMin()
	if nodeID != 4 || dist != 5.0 {
		t.Errorf("DecreaseKey a nuovo minimo: atteso (4, 5.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestDecreaseKeyMiddle(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(1, 100.0)
	heap.Insert(2, 200.0)
	heap.Insert(3, 300.0)
	heap.Insert(4, 400.0)

	// Diminuisci nodeID=3 a un valore intermedio
	heap.DecreaseKey(3, 150.0)

	// Ordine atteso: 1(100), 2(150 era 3), 2(200), 4(400)
	expectedOrder := []struct {
		nodeID int
		dist   float64
	}{
		{1, 100.0},
		{3, 150.0},
		{2, 200.0},
		{4, 400.0},
	}

	for i, expected := range expectedOrder {
		nodeID, dist := heap.ExtractMin()
		if nodeID != expected.nodeID || dist != expected.dist {
			t.Errorf("Iterazione %d: atteso (%d, %.1f), ottenuto (%d, %.1f)",
				i, expected.nodeID, expected.dist, nodeID, dist)
		}
	}
}

func TestDecreaseKeyMultiple(t *testing.T) {
	heap := CreateBinaryHeap()
	heap.Insert(1, 1000.0)

	// Diminuisci progressivamente
	heap.DecreaseKey(1, 500.0)
	heap.DecreaseKey(1, 250.0)
	heap.DecreaseKey(1, 100.0)
	heap.DecreaseKey(1, 50.0)

	nodeID, dist := heap.ExtractMin()
	if nodeID != 1 || dist != 50.0 {
		t.Errorf("DecreaseKey multipli: atteso (1, 50.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestDecreaseKeyMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinaryHeap()

	// Insert elementi
	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i*10))
	}

	// DecreaseKey su vari nodi
	heap.DecreaseKey(5, 5.0)
	heap.DecreaseKey(8, 12.0)
	heap.DecreaseKey(3, 1.0)

	// Verifica proprietà heap
	for i := 0; i < heap.Size(); i++ {
		leftChild := 2*i + 1
		rightChild := 2*i + 2

		if leftChild < heap.Size() {
			if heap.items[i].distance > heap.items[leftChild].distance {
				t.Error("DecreaseKey ha violato proprietà heap (left child)")
			}
		}

		if rightChild < heap.Size() {
			if heap.items[i].distance > heap.items[rightChild].distance {
				t.Error("DecreaseKey ha violato proprietà heap (right child)")
			}
		}
	}
}

func TestDecreaseKeyUpdatesPosition(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 200.0)
	heap.Insert(30, 300.0)

	// Diminuisci nodeID=30
	heap.DecreaseKey(30, 50.0)

	// Verifica che position map sia aggiornata correttamente
	idx := heap.position[30]
	if heap.items[idx].nodeID != 30 {
		t.Error("Position map non aggiornata correttamente dopo DecreaseKey")
	}

	if heap.items[idx].distance != 50.0 {
		t.Error("Distanza non aggiornata correttamente")
	}
}

// ============================================ //
//      TEST OPERAZIONI MISTE (REALISTICI)     //
// ============================================ //

func TestMixedOperationsDijkstraPattern(t *testing.T) {
	heap := CreateBinaryHeap()

	// Simula pattern Dijkstra
	heap.Insert(0, 0.0) // Source node

	// Scopri vicini
	heap.Insert(1, 10.0)
	heap.Insert(2, 5.0)
	heap.Insert(3, 15.0)

	// Processa nodo con distanza minima (source)
	nodeID, dist := heap.ExtractMin()
	if nodeID != 0 || dist != 0.0 {
		t.Error("Dovrebbe estrarre source node per primo")
	}

	// Processa secondo minimo
	nodeID, _ = heap.ExtractMin()
	if nodeID != 2 { // dist=5.0
		t.Error("Dovrebbe estrarre nodeID=2 (dist=5.0)")
	}

	// Trova percorso migliore verso nodeID=1
	heap.DecreaseKey(1, 7.0)

	// Aggiungi nuovi nodi scoperti
	heap.Insert(4, 12.0)
	heap.Insert(5, 8.0)

	// Continua estrazione
	expectedOrder := []int{1, 5, 4, 3} // Per distanze: 7, 8, 12, 15

	for i, expectedID := range expectedOrder {
		nodeID, _ := heap.ExtractMin()
		if nodeID != expectedID {
			t.Errorf("Iterazione %d: atteso nodeID=%d, ottenuto %d",
				i, expectedID, nodeID)
		}
	}
}

func TestMixedOperationsAlternating(t *testing.T) {
	heap := CreateBinaryHeap()

	// Alterna Insert ed ExtractMin
	heap.Insert(1, 50.0)
	heap.Insert(2, 30.0)
	heap.Insert(3, 70.0)

	heap.ExtractMin() // Estrai 30.0

	heap.Insert(4, 20.0)
	heap.Insert(5, 40.0)

	heap.ExtractMin() // Estrai 20.0

	heap.DecreaseKey(3, 25.0)

	// Ordine finale atteso: 3(25), 5(40), 1(50)
	expectedOrder := []float64{25.0, 40.0, 50.0}

	for i, expectedDist := range expectedOrder {
		_, dist := heap.ExtractMin()
		if dist != expectedDist {
			t.Errorf("Iterazione %d: distanza attesa %.1f, ottenuta %.1f",
				i, expectedDist, dist)
		}
	}
}

func TestMixedOperationsStressTest(t *testing.T) {
	heap := CreateBinaryHeap()

	// Pattern realistico: insert, decrease, extract
	for i := 0; i < 100; i++ {
		heap.Insert(i, float64(i*3))
	}

	// DecreaseKey su alcuni nodi
	for i := 10; i < 20; i++ {
		heap.DecreaseKey(i, float64(i))
	}

	// Estrai alcuni
	for i := 0; i < 30; i++ {
		heap.ExtractMin()
	}

	// Inserisci altri
	for i := 100; i < 120; i++ {
		heap.Insert(i, float64(i))
	}

	// DecreaseKey su nodi recenti
	for i := 100; i < 110; i++ {
		if _, exists := heap.position[i]; exists {
			heap.DecreaseKey(i, float64(i/2))
		}
	}

	// Verifica che l'ordine sia mantenuto
	lastDist := -1.0
	for !heap.IsEmpty() {
		_, dist := heap.ExtractMin()
		if dist < lastDist {
			t.Errorf("Ordine violato: %.1f < %.1f", dist, lastDist)
		}
		lastDist = dist
	}
}

// ============================================ //
//          TEST EDGE CASES                     //
// ============================================ //

func TestFloatingPointPrecision(t *testing.T) {
	heap := CreateBinaryHeap()

	// Distanze con molti decimali (realistiche per A*)
	heap.Insert(1, 10.123456)
	heap.Insert(2, 10.123457)
	heap.Insert(3, 10.123455)

	_, dist1 := heap.ExtractMin()
	_, dist2 := heap.ExtractMin()
	_, dist3 := heap.ExtractMin()

	if dist1 >= dist2 || dist2 >= dist3 {
		t.Error("Ordine non corretto con floating point precision")
	}
}

func TestZeroDistances(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(1, 0.0)
	heap.Insert(2, 10.0)
	heap.Insert(3, 0.0)

	// Primo estratto dovrebbe avere dist=0.0
	_, dist := heap.ExtractMin()
	if dist != 0.0 {
		t.Error("Distanza 0.0 non gestita correttamente")
	}

	// Anche il secondo dovrebbe avere dist=0.0
	_, dist = heap.ExtractMin()
	if dist != 0.0 {
		t.Error("Secondo 0.0 non estratto correttamente")
	}
}

func TestNegativeDistances(t *testing.T) {
	heap := CreateBinaryHeap()

	heap.Insert(1, -10.0)
	heap.Insert(2, 5.0)
	heap.Insert(3, -5.0)
	heap.Insert(4, 0.0)

	expectedOrder := []float64{-10.0, -5.0, 0.0, 5.0}

	for i, expectedDist := range expectedOrder {
		_, dist := heap.ExtractMin()
		if dist != expectedDist {
			t.Errorf("Iterazione %d: distanza attesa %.1f, ottenuta %.1f",
				i, expectedDist, dist)
		}
	}
}

func TestLargeHeap(t *testing.T) {
	heap := CreateBinaryHeap()
	n := 10000

	// Insert n elementi in ordine casuale
	for i := 0; i < n; i++ {
		heap.Insert(i, float64(n-i))
	}

	if heap.Size() != n {
		t.Errorf("Size attesa %d, ottenuta %d", n, heap.Size())
	}

	// ExtractMin deve restituire tutti in ordine
	lastDist := -1.0
	count := 0

	for !heap.IsEmpty() {
		_, dist := heap.ExtractMin()
		if dist < lastDist {
			t.Errorf("Ordine violato: %.1f < %.1f", dist, lastDist)
		}
		lastDist = dist
		count++
	}

	if count != n {
		t.Errorf("Estratti %d elementi, attesi %d", count, n)
	}
}

func TestDuplicateDistances(t *testing.T) {
	heap := CreateBinaryHeap(50.0, 30.0, 50.0, 30.0, 50.0)

	// Conta le occorrenze
	count30 := 0
	count50 := 0

	for !heap.IsEmpty() {
		_, dist := heap.ExtractMin()
		if dist == 30.0 {
			count30++
		} else if dist == 50.0 {
			count50++
		}
	}

	if count30 != 2 || count50 != 3 {
		t.Errorf("Duplicati non gestiti: count30=%d (atteso 2), count50=%d (atteso 3)",
			count30, count50)
	}
}

// ============================================ //
//          TEST SIZE E ISEMPTY                 //
// ============================================ //

func TestSizeTracking(t *testing.T) {
	heap := CreateBinaryHeap()

	if heap.Size() != 0 {
		t.Error("Heap iniziale dovrebbe avere size 0")
	}

	heap.Insert(1, 10.0)
	if heap.Size() != 1 {
		t.Error("Size non incrementata dopo Insert")
	}

	heap.Insert(2, 20.0)
	heap.Insert(3, 30.0)
	if heap.Size() != 3 {
		t.Error("Size non corretta dopo Insert multipli")
	}

	heap.ExtractMin()
	if heap.Size() != 2 {
		t.Error("Size non decrementata dopo ExtractMin")
	}

	heap.DecreaseKey(2, 5.0)
	if heap.Size() != 2 {
		t.Error("DecreaseKey non dovrebbe modificare size")
	}
}

func TestIsEmpty(t *testing.T) {
	heap := CreateBinaryHeap()

	if !heap.IsEmpty() {
		t.Error("Heap nuovo dovrebbe essere vuoto")
	}

	heap.Insert(1, 10.0)
	if heap.IsEmpty() {
		t.Error("Heap con elementi non dovrebbe essere vuoto")
	}

	heap.ExtractMin()
	if !heap.IsEmpty() {
		t.Error("Heap dopo ExtractMin di tutti gli elementi dovrebbe essere vuoto")
	}
}

// ============================================ //
//          TEST POSITION MAP                   //
// ============================================ //

func TestPositionMapConsistency(t *testing.T) {
	heap := CreateBinaryHeap()

	// Insert elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64(i*5))
	}

	// Verifica consistenza: position[nodeID] deve puntare all'indice corretto
	for nodeID, idx := range heap.position {
		if heap.items[idx].nodeID != nodeID {
			t.Errorf("Position map inconsistente: nodeID=%d punta a idx=%d, ma items[%d].nodeID=%d",
				nodeID, idx, idx, heap.items[idx].nodeID)
		}
	}

	// Dopo DecreaseKey
	heap.DecreaseKey(15, 2.0)

	for nodeID, idx := range heap.position {
		if heap.items[idx].nodeID != nodeID {
			t.Error("Position map inconsistente dopo DecreaseKey")
		}
	}

	// Dopo ExtractMin
	heap.ExtractMin()

	for nodeID, idx := range heap.position {
		if heap.items[idx].nodeID != nodeID {
			t.Error("Position map inconsistente dopo ExtractMin")
		}
	}
}

func TestPositionMapSize(t *testing.T) {
	heap := CreateBinaryHeap()

	// La size di position map deve sempre corrispondere a heap.Size()
	if len(heap.position) != heap.Size() {
		t.Error("Position map size iniziale non corrisponde a heap size")
	}

	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i))
		if len(heap.position) != heap.Size() {
			t.Errorf("Dopo Insert %d: position size=%d, heap size=%d",
				i, len(heap.position), heap.Size())
		}
	}

	for i := 0; i < 5; i++ {
		heap.ExtractMin()
		if len(heap.position) != heap.Size() {
			t.Errorf("Dopo ExtractMin %d: position size=%d, heap size=%d",
				i, len(heap.position), heap.Size())
		}
	}
}

// ============================================ //
//              BENCHMARK                       //
// ============================================ //

func BenchmarkBinaryInsert(b *testing.B) {
	heap := CreateBinaryHeap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(i))
	}
}

func BenchmarkBinaryInsertRandom(b *testing.B) {
	heap := CreateBinaryHeap()
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, rand.Float64()*1000)
	}
}

func BenchmarkBinaryExtractMin(b *testing.B) {
	// Prepara heap con b.N elementi
	heap := CreateBinaryHeap()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(b.N-i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.ExtractMin()
	}
}

func BenchmarkBinaryDecreaseKey(b *testing.B) {
	// Prepara heap con 10000 elementi
	heap := CreateBinaryHeap()
	n := 10000
	for i := 0; i < n; i++ {
		heap.Insert(i, float64(i*10))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := i % n
		newDist := float64((i % (n * 5)))
		heap.DecreaseKey(nodeID, newDist)
	}
}

func BenchmarkBinaryMixedDijkstraPattern(b *testing.B) {
	heap := CreateBinaryHeap()
	nodeCounter := 0
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := i % 10

		if op < 6 { // 60% Insert
			heap.Insert(nodeCounter, rand.Float64()*1000)
			nodeCounter++
		} else if op < 8 && !heap.IsEmpty() { // 20% ExtractMin
			heap.ExtractMin()
		} else if !heap.IsEmpty() { // 20% DecreaseKey
			// DecreaseKey su nodo casuale presente
			if heap.Size() > 0 {
				// Scegli un nodeID casuale tra quelli presenti
				targetID := rand.Intn(nodeCounter)
				if _, exists := heap.position[targetID]; exists {
					heap.DecreaseKey(targetID, rand.Float64()*100)
				}
			}
		}
	}
}

func BenchmarkBinaryCreateHeap(b *testing.B) {
	// Benchmark della creazione heap con heapify
	n := 1000
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = float64(n - i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateBinaryHeap(values...)
	}
}

func BenchmarkBinarySize(b *testing.B) {
	heap := CreateBinaryHeap()
	for i := 0; i < 1000; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.Size()
	}
}

func BenchmarkBinaryIsEmpty(b *testing.B) {
	heap := CreateBinaryHeap()
	for i := 0; i < 1000; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.IsEmpty()
	}
}

// ============================================ //
//     BENCHMARK COMPARATIVI (SCALABILITÀ)     //
// ============================================ //

func BenchmarkBinaryInsert_100(b *testing.B) {
	benchmarkInsertN(b, 100)
}

func BenchmarkBinaryInsert_1000(b *testing.B) {
	benchmarkInsertN(b, 1000)
}

func BenchmarkBinaryInsert_10000(b *testing.B) {
	benchmarkInsertN(b, 10000)
}

func benchmarkInsertN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
	}
}

func BenchmarkBinaryExtractMin_100(b *testing.B) {
	benchmarkExtractMinN(b, 100)
}

func BenchmarkBinaryExtractMin_1000(b *testing.B) {
	benchmarkExtractMinN(b, 1000)
}

func BenchmarkBinaryExtractMin_10000(b *testing.B) {
	benchmarkExtractMinN(b, 10000)
}

func benchmarkExtractMinN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.ExtractMin()
		}
	}
}

func BenchmarkBinaryDecreaseKey_100(b *testing.B) {
	benchmarkDecreaseKeyN(b, 100)
}

func BenchmarkBinaryDecreaseKey_1000(b *testing.B) {
	benchmarkDecreaseKeyN(b, 1000)
}

func BenchmarkBinaryDecreaseKey_10000(b *testing.B) {
	benchmarkDecreaseKeyN(b, 10000)
}

func benchmarkDecreaseKeyN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j*10))
		}
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.DecreaseKey(j, float64(j))
		}
	}
}

// ============================================ //
//   BENCHMARK SEQUENZE REALISTICHE DIJKSTRA   //
// ============================================ //

func BenchmarkBinaryDijkstraSmallGraph(b *testing.B) {
	// Simula Dijkstra su grafo piccolo (100 nodi)
	benchmarkDijkstraPattern(b, 100, 300) // 100 nodi, ~300 archi
}

func BenchmarkBinaryDijkstraMediumGraph(b *testing.B) {
	// Simula Dijkstra su grafo medio (1000 nodi)
	benchmarkDijkstraPattern(b, 1000, 5000)
}

func BenchmarkBinaryDijkstraLargeGraph(b *testing.B) {
	// Simula Dijkstra su grafo grande (10000 nodi)
	benchmarkDijkstraPattern(b, 10000, 50000)
}

func benchmarkDijkstraPattern(b *testing.B, numNodes int, numEdges int) {
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap := CreateBinaryHeap()

		// Insert source
		heap.Insert(0, 0.0)
		visited := make(map[int]bool)

		for processed := 0; processed < numNodes && !heap.IsEmpty(); processed++ {
			// ExtractMin
			nodeID, _ := heap.ExtractMin()
			visited[nodeID] = true

			// Simula esplorazione vicini (media 3-5 vicini per nodo)
			numNeighbors := 3 + rand.Intn(3)
			for j := 0; j < numNeighbors; j++ {
				neighborID := rand.Intn(numNodes)
				if visited[neighborID] {
					continue
				}

				newDist := rand.Float64() * 100

				if _, exists := heap.position[neighborID]; exists {
					// Nodo già nell'heap, potenziale DecreaseKey
					heap.DecreaseKey(neighborID, newDist)
				} else {
					// Nuovo nodo
					heap.Insert(neighborID, newDist)
				}
			}
		}
	}
}

// ============================================ //
//        BENCHMARK WORST CASE                  //
// ============================================ //

func BenchmarkBinaryWorstCaseInsert(b *testing.B) {
	// Worst case: inserimenti in ordine crescente
	// Ogni insert richiede bubbleUp fino alla radice
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkBinaryWorstCaseExtractMin(b *testing.B) {
	// Worst case: heap completamente sbilanciato
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		// Insert in ordine decrescente
		for j := 1000; j > 0; j-- {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			heap.ExtractMin()
		}
	}
}

func BenchmarkBinaryWorstCaseDecreaseKey(b *testing.B) {
	// Worst case: DecreaseKey che deve fare bubbleUp fino alla radice
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		// DecreaseKey sull'ultimo elemento a 0 (deve risalire tutto)
		for j := 0; j < 100; j++ {
			targetID := n - 1 - (j % n)
			heap.DecreaseKey(targetID, 0.0)
		}
	}
}

// ============================================ //
//        BENCHMARK BEST CASE                   //
// ============================================ //

func BenchmarkBinaryBestCaseInsert(b *testing.B) {
	// Best case: inserimenti in ordine decrescente
	// Ogni insert va direttamente in fondo (no bubbleUp)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		b.StartTimer()

		for j := 1000; j > 0; j-- {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkBinaryBestCaseDecreaseKey(b *testing.B) {
	// Best case: DecreaseKey che non richiede bubbleUp
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinaryHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		// DecreaseKey di poco (no movimento)
		for j := 0; j < 100; j++ {
			heap.DecreaseKey(j%n, float64(j%n)-0.1)
		}
	}
}

// ============================================ //
//     BENCHMARK MEMORY & POSITION MAP         //
// ============================================ //

func BenchmarkBinaryPositionLookup(b *testing.B) {
	heap := CreateBinaryHeap()
	n := 10000

	for i := 0; i < n; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := i % n
		_ = heap.position[nodeID]
	}
}

func BenchmarkBinarySwap(b *testing.B) {
	heap := CreateBinaryHeap()

	for i := 0; i < 100; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.swap(0, 50)
		heap.swap(50, 0) // Ripristina
	}
}
