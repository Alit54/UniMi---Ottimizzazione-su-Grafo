package heap

import (
	"math/rand"
	"testing"
)

// ============================================ //
//          TEST COSTRUTTORE                    //
// ============================================ //

func TestCreateBinomialHeap(t *testing.T) {
	values := []float64{50.0, 30.0, 70.0, 10.0, 40.0}
	heap := CreateBinomialHeap(values...)

	if heap.Size() != 5 {
		t.Errorf("Size attesa 5, ottenuta %d", heap.Size())
	}

	// Verifica che il minimo sia corretto
	if heap.min.distance != 10.0 {
		t.Errorf("Minimo dovrebbe essere 10.0, ottenuto %.1f", heap.min.distance)
	}
}

func TestCreateBinomialHeapEmpty(t *testing.T) {
	heap := CreateBinomialHeap()

	if heap.Size() != 0 {
		t.Error("Heap vuoto dovrebbe avere size 0")
	}

	if !heap.IsEmpty() {
		t.Error("Heap vuoto dovrebbe risultare isEmpty")
	}
}

func TestCreateBinomialHeapSingle(t *testing.T) {
	heap := CreateBinomialHeap(42.0)

	if heap.Size() != 1 {
		t.Error("Heap con singolo elemento dovrebbe avere size 1")
	}

	nodeID, dist := heap.ExtractMin()
	if nodeID != 0 || dist != 42.0 {
		t.Errorf("Atteso (0, 42.0), ottenuto (%d, %.1f)", nodeID, dist)
	}
}

func TestCreateBinomialHeapPowerOfTwo(t *testing.T) {
	// Test con potenze di 2 (struttura ottimale per binomial heap)
	for n := 1; n <= 64; n *= 2 {
		heap := CreateBinomialHeap()
		for i := 0; i < n; i++ {
			heap.Insert(i, float64(n-i))
		}

		if heap.Size() != n {
			t.Errorf("Per n=%d, size attesa %d, ottenuta %d", n, n, heap.Size())
		}

		// Verifica che il minimo sia corretto
		if heap.min.distance != 1.0 {
			t.Errorf("Per n=%d, minimo atteso 1.0, ottenuto %.1f", n, heap.min.distance)
		}
	}
}

func TestCreateBinomialHeapOrdered(t *testing.T) {
	// Test con valori già ordinati
	heap := CreateBinomialHeap(10.0, 20.0, 30.0, 40.0, 50.0)

	if heap.min.distance != 10.0 {
		t.Error("Heap dovrebbe mantenere 10.0 come minimo")
	}
}

func TestCreateBinomialHeapReverse(t *testing.T) {
	// Test con valori in ordine inverso
	heap := CreateBinomialHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	if heap.min.distance != 10.0 {
		t.Error("Heap dovrebbe avere 10.0 come minimo")
	}
}

// ============================================ //
//          TEST INSERT                         //
// ============================================ //

func TestBinomialInsertSingle(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(1, 10.0)

	if heap.Size() != 1 {
		t.Errorf("Size attesa 1, ottenuta %d", heap.Size())
	}

	if heap.min.nodeID != 1 || heap.min.distance != 10.0 {
		t.Error("Elemento inserito non corretto")
	}
}

func TestBinomialInsertMultiple(t *testing.T) {
	heap := CreateBinomialHeap()

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
	if heap.min.nodeID != 1 || heap.min.distance != 10.0 {
		t.Errorf("Minimo atteso (1, 10.0), ottenuto (%d, %.1f)",
			heap.min.nodeID, heap.min.distance)
	}
}

func TestBinomialInsertMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinomialHeap()

	// Insert elementi casuali
	values := []float64{50.0, 30.0, 70.0, 20.0, 60.0, 10.0, 80.0}
	for i, val := range values {
		heap.Insert(i, val)
	}

	// Verifica proprietà min-heap: ogni padre <= figli
	verifyBinomialHeapProperty(t, heap)
}

func TestBinomialInsertUpdatePosition(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 200.0)
	heap.Insert(30, 300.0)

	// Verifica che position map sia aggiornata
	for nodeID, node := range heap.pos {
		if node.nodeID != nodeID {
			t.Errorf("Position map inconsistente per nodeID=%d", nodeID)
		}
	}
}

func TestBinomialInsertSequential(t *testing.T) {
	heap := CreateBinomialHeap()

	// Insert 16 elementi (potenza di 2)
	for i := 0; i < 16; i++ {
		heap.Insert(i, float64(16-i))
	}

	if heap.Size() != 16 {
		t.Errorf("Size attesa 16, ottenuta %d", heap.Size())
	}

	// Verifica minimo
	if heap.min.distance != 1.0 {
		t.Errorf("Minimo atteso 1.0, ottenuto %.1f", heap.min.distance)
	}
}

// ============================================ //
//          TEST EXTRACTMIN                     //
// ============================================ //

func TestBinomialExtractMinOrder(t *testing.T) {
	heap := CreateBinomialHeap(100.0, 50.0, 150.0, 30.0, 70.0, 10.0)

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

func TestBinomialExtractMinCorrectNodeIDs(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialExtractMinUpdatesPosition(t *testing.T) {
	heap := CreateBinomialHeap(50.0, 30.0, 70.0, 10.0, 40.0)

	initialSize := heap.Size()
	minNodeID := heap.min.nodeID

	heap.ExtractMin()

	// Verifica che nodeID estratto non sia più in position map
	if _, exists := heap.pos[minNodeID]; exists {
		t.Error("NodeID estratto dovrebbe essere rimosso dalla position map")
	}

	// Verifica che size sia diminuita
	if heap.Size() != initialSize-1 {
		t.Error("Size dovrebbe diminuire dopo ExtractMin")
	}

	// Verifica consistenza position map
	if len(heap.pos) != heap.Size() {
		t.Errorf("Position map size (%d) != heap size (%d)",
			len(heap.pos), heap.Size())
	}
}

func TestBinomialExtractMinSingle(t *testing.T) {
	heap := CreateBinomialHeap(42.0)

	nodeID, dist := heap.ExtractMin()

	if nodeID != 0 || dist != 42.0 {
		t.Errorf("ExtractMin su singolo elemento: atteso (0, 42.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto dopo ExtractMin su singolo elemento")
	}
}

func TestBinomialExtractMinMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinomialHeap()

	// Insert molti elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64((i*7)%20))
	}

	// Estrai alcuni elementi
	for i := 0; i < 10; i++ {
		heap.ExtractMin()
	}

	// Verifica proprietà heap
	verifyBinomialHeapProperty(t, heap)
}

func TestBinomialExtractMinEmpty(t *testing.T) {
	heap := CreateBinomialHeap()

	// ExtractMin su heap vuoto non dovrebbe crashare
	nodeID, dist := heap.ExtractMin()

	// Dovrebbe restituire valori zero
	if nodeID != 0 || dist != 0.0 {
		t.Errorf("ExtractMin su heap vuoto: atteso (0, 0.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestBinomialExtractMinAll(t *testing.T) {
	heap := CreateBinomialHeap()

	n := 10
	for i := 0; i < n; i++ {
		heap.Insert(i, float64(n-i))
	}

	// Estrai tutti e verifica ordine
	lastDist := -1.0
	for i := 0; i < n; i++ {
		_, dist := heap.ExtractMin()
		if dist < lastDist {
			t.Errorf("Ordine violato: %.1f < %.1f", dist, lastDist)
		}
		lastDist = dist
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto")
	}
}

// ============================================ //
//          TEST DECREASEKEY                    //
// ============================================ //

func TestBinomialDecreaseKeyBasic(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialDecreaseKeyToMinimum(t *testing.T) {
	heap := CreateBinomialHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	// Diminuisci il nodo 4 (distance=10.0) a 5.0
	heap.DecreaseKey(4, 5.0)

	nodeID, dist := heap.ExtractMin()
	if nodeID != 4 || dist != 5.0 {
		t.Errorf("DecreaseKey a nuovo minimo: atteso (4, 5.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestBinomialDecreaseKeyMiddle(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(1, 100.0)
	heap.Insert(2, 200.0)
	heap.Insert(3, 300.0)
	heap.Insert(4, 400.0)

	// Diminuisci nodeID=3 a un valore intermedio
	heap.DecreaseKey(3, 150.0)

	// Ordine atteso: 1(100), 3(150), 2(200), 4(400)
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

func TestBinomialDecreaseKeyMultiple(t *testing.T) {
	heap := CreateBinomialHeap()
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

func TestBinomialDecreaseKeyMaintainsHeapProperty(t *testing.T) {
	heap := CreateBinomialHeap()

	// Insert elementi
	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i*10))
	}

	// DecreaseKey su vari nodi
	heap.DecreaseKey(5, 5.0)
	heap.DecreaseKey(8, 12.0)
	heap.DecreaseKey(3, 1.0)

	// Verifica proprietà heap
	verifyBinomialHeapProperty(t, heap)
}

func TestBinomialDecreaseKeyUpdatesPosition(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 200.0)
	heap.Insert(30, 300.0)

	// Diminuisci nodeID=30
	heap.DecreaseKey(30, 50.0)

	// Verifica che position map sia ancora consistente
	node := heap.pos[30]
	if node.nodeID != 30 {
		t.Error("Position map non aggiornata correttamente dopo DecreaseKey")
	}

	if node.distance != 50.0 {
		t.Error("Distanza non aggiornata correttamente")
	}
}

func TestBinomialDecreaseKeyNonExistent(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(1, 100.0)
	heap.Insert(2, 200.0)

	initialSize := heap.Size()

	// DecreaseKey su nodo inesistente non dovrebbe modificare heap
	heap.DecreaseKey(999, 50.0)

	if heap.Size() != initialSize {
		t.Error("DecreaseKey su nodo inesistente ha modificato size")
	}
}

func TestBinomialDecreaseKeyInvalidIncrease(t *testing.T) {
	heap := CreateBinomialHeap()

	heap.Insert(1, 100.0)

	// Tentativo di "aumentare" (non decrease)
	heap.DecreaseKey(1, 200.0)

	// Il valore non dovrebbe essere cambiato
	nodeID, dist := heap.ExtractMin()
	if dist != 100.0 {
		t.Errorf("DecreaseKey con aumento non dovrebbe modificare: atteso 100.0, ottenuto %.1f", dist)
	}

	if nodeID != 1 {
		t.Error("NodeID non corretto")
	}
}

// ============================================ //
//      TEST OPERAZIONI MISTE (REALISTICI)     //
// ============================================ //

func TestBinomialMixedOperationsDijkstraPattern(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialMixedOperationsAlternating(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialMixedOperationsStressTest(t *testing.T) {
	heap := CreateBinomialHeap()

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
		if _, exists := heap.pos[i]; exists {
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

func TestBinomialFloatingPointPrecision(t *testing.T) {
	heap := CreateBinomialHeap()

	// Distanze con molti decimali
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

func TestBinomialZeroDistances(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialNegativeDistances(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialLargeHeap(t *testing.T) {
	heap := CreateBinomialHeap()
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

func TestBinomialDuplicateDistances(t *testing.T) {
	heap := CreateBinomialHeap(50.0, 30.0, 50.0, 30.0, 50.0)

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

func TestBinomialAllSameDistances(t *testing.T) {
	heap := CreateBinomialHeap(7.0, 7.0, 7.0, 7.0)

	for i := 0; i < 4; i++ {
		_, dist := heap.ExtractMin()
		if dist != 7.0 {
			t.Errorf("Atteso 7.0, ottenuto %.1f", dist)
		}
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto")
	}
}

// ============================================ //
//          TEST SIZE, ISEMPTY                  //
// ============================================ //

func TestBinomialSizeTracking(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialIsEmpty(t *testing.T) {
	heap := CreateBinomialHeap()

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

func TestBinomialPositionMapConsistency(t *testing.T) {
	heap := CreateBinomialHeap()

	// Insert elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64(i*5))
	}

	// Verifica consistenza: pos[nodeID] deve puntare al nodo corretto
	for nodeID, node := range heap.pos {
		if node.nodeID != nodeID {
			t.Errorf("Position map inconsistente: nodeID=%d, ma node.nodeID=%d",
				nodeID, node.nodeID)
		}
	}

	// Dopo DecreaseKey
	heap.DecreaseKey(15, 2.0)

	for nodeID, node := range heap.pos {
		if node.nodeID != nodeID {
			t.Error("Position map inconsistente dopo DecreaseKey")
		}
	}

	// Dopo ExtractMin
	heap.ExtractMin()

	for nodeID, node := range heap.pos {
		if node.nodeID != nodeID {
			t.Error("Position map inconsistente dopo ExtractMin")
		}
	}
}

func TestBinomialPositionMapSize(t *testing.T) {
	heap := CreateBinomialHeap()

	// La size di position map deve sempre corrispondere a heap.Size()
	if len(heap.pos) != heap.Size() {
		t.Error("Position map size iniziale non corrisponde a heap size")
	}

	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i))
		if len(heap.pos) != heap.Size() {
			t.Errorf("Dopo Insert %d: position size=%d, heap size=%d",
				i, len(heap.pos), heap.Size())
		}
	}

	for i := 0; i < 5; i++ {
		heap.ExtractMin()
		if len(heap.pos) != heap.Size() {
			t.Errorf("Dopo ExtractMin %d: position size=%d, heap size=%d",
				i, len(heap.pos), heap.Size())
		}
	}
}

// ============================================ //
//            HELPER FUNCTIONS                  //
// ============================================ //

// verifyBinomialHeapProperty verifica ricorsivamente la proprietà min-heap
func verifyBinomialHeapProperty(t *testing.T, heap *BinomialHeap) {
	current := heap.head

	for current != nil {
		verifyBinomialTree(t, current)
		current = current.sibling
	}
}

// verifyBinomialTree verifica la proprietà heap in un albero binomiale
func verifyBinomialTree(t *testing.T, node *binomialNode) {
	if node == nil {
		return
	}

	// Verifica che ogni figlio abbia distanza >= padre
	child := node.child
	for child != nil {
		if child.distance < node.distance {
			t.Errorf("Proprietà heap violata: figlio (%.1f) < padre (%.1f)",
				child.distance, node.distance)
		}
		verifyBinomialTree(t, child)
		child = child.sibling
	}
}

// ============================================ //
//              BENCHMARK                       //
// ============================================ //

func BenchmarkBinomialInsert(b *testing.B) {
	heap := CreateBinomialHeap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(i))
	}
}

func BenchmarkBinomialInsertRandom(b *testing.B) {
	heap := CreateBinomialHeap()
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, rand.Float64()*1000)
	}
}

func BenchmarkBinomialExtractMin(b *testing.B) {
	// Prepara heap con b.N elementi
	heap := CreateBinomialHeap()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(b.N-i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.ExtractMin()
	}
}

func BenchmarkBinomialDecreaseKey(b *testing.B) {
	// Prepara heap con 10000 elementi
	heap := CreateBinomialHeap()
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

func BenchmarkBinomialMixedDijkstraPattern(b *testing.B) {
	heap := CreateBinomialHeap()
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
			if heap.Size() > 0 {
				targetID := rand.Intn(nodeCounter)
				if _, exists := heap.pos[targetID]; exists {
					heap.DecreaseKey(targetID, rand.Float64()*100)
				}
			}
		}
	}
}

func BenchmarkBinomialCreateHeap(b *testing.B) {
	n := 1000
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = float64(n - i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateBinomialHeap(values...)
	}
}

func BenchmarkBinomialSize(b *testing.B) {
	heap := CreateBinomialHeap()
	for i := 0; i < 1000; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.Size()
	}
}

func BenchmarkBinomialIsEmpty(b *testing.B) {
	heap := CreateBinomialHeap()
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

func BenchmarkBinomialInsert_100(b *testing.B) {
	benchmarkBinomialInsertN(b, 100)
}

func BenchmarkBinomialInsert_1000(b *testing.B) {
	benchmarkBinomialInsertN(b, 1000)
}

func BenchmarkBinomialInsert_10000(b *testing.B) {
	benchmarkBinomialInsertN(b, 10000)
}

func benchmarkBinomialInsertN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
	}
}

func BenchmarkBinomialExtractMin_100(b *testing.B) {
	benchmarkBinomialExtractMinN(b, 100)
}

func BenchmarkBinomialExtractMin_1000(b *testing.B) {
	benchmarkBinomialExtractMinN(b, 1000)
}

func BenchmarkBinomialExtractMin_10000(b *testing.B) {
	benchmarkBinomialExtractMinN(b, 10000)
}

func benchmarkBinomialExtractMinN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.ExtractMin()
		}
	}
}

func BenchmarkBinomialDecreaseKey_100(b *testing.B) {
	benchmarkBinomialDecreaseKeyN(b, 100)
}

func BenchmarkBinomialDecreaseKey_1000(b *testing.B) {
	benchmarkBinomialDecreaseKeyN(b, 1000)
}

func BenchmarkBinomialDecreaseKey_10000(b *testing.B) {
	benchmarkBinomialDecreaseKeyN(b, 10000)
}

func benchmarkBinomialDecreaseKeyN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
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

func BenchmarkBinomialDijkstraSmallGraph(b *testing.B) {
	// Simula Dijkstra su grafo piccolo (100 nodi)
	benchmarkBinomialDijkstraPattern(b, 100, 300)
}

func BenchmarkBinomialDijkstraMediumGraph(b *testing.B) {
	// Simula Dijkstra su grafo medio (1000 nodi)
	benchmarkBinomialDijkstraPattern(b, 1000, 5000)
}

func BenchmarkBinomialDijkstraLargeGraph(b *testing.B) {
	// Simula Dijkstra su grafo grande (10000 nodi)
	benchmarkBinomialDijkstraPattern(b, 10000, 50000)
}

func benchmarkBinomialDijkstraPattern(b *testing.B, numNodes int, numEdges int) {
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap := CreateBinomialHeap()

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

				if _, exists := heap.pos[neighborID]; exists {
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

func BenchmarkBinomialWorstCaseInsert(b *testing.B) {
	// Worst case: inserimenti in ordine crescente
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkBinomialWorstCaseExtractMin(b *testing.B) {
	// Worst case: heap con molti alberi
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
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

func BenchmarkBinomialWorstCaseDecreaseKey(b *testing.B) {
	// Worst case: DecreaseKey che deve fare bubbleUp fino alla radice
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		// DecreaseKey su elementi che devono risalire molto
		for j := 0; j < 100; j++ {
			targetID := n - 1 - (j % n)
			if _, exists := heap.pos[targetID]; exists {
				heap.DecreaseKey(targetID, 0.0)
			}
		}
	}
}

// ============================================ //
//        BENCHMARK BEST CASE                   //
// ============================================ //

func BenchmarkBinomialBestCaseInsert(b *testing.B) {
	// Best case: inserimenti che mantengono struttura ottimale
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		b.StartTimer()

		for j := 1000; j > 0; j-- {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkBinomialBestCaseDecreaseKey(b *testing.B) {
	// Best case: DecreaseKey che non richiede bubbleUp
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateBinomialHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		// DecreaseKey di poco (minimo movimento)
		for j := 0; j < 100; j++ {
			if _, exists := heap.pos[j%n]; exists {
				heap.DecreaseKey(j%n, float64(j%n)-0.1)
			}
		}
	}
}

// ============================================ //
//     BENCHMARK MEMORY & POSITION MAP         //
// ============================================ //

func BenchmarkBinomialPositionLookup(b *testing.B) {
	heap := CreateBinomialHeap()
	n := 10000

	for i := 0; i < n; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := i % n
		_ = heap.pos[nodeID]
	}
}

// ============================================ //
//     TEST STRUTTURA BINOMIAL HEAP            //
// ============================================ //

func TestBinomialTreeStructure(b *testing.T) {
	heap := CreateBinomialHeap()

	// Insert 7 elementi (7 = 111 in binario = B0 + B1 + B2)
	for i := 0; i < 7; i++ {
		heap.Insert(i, float64(7-i))
	}

	// Dovrebbero esserci 3 alberi binomiali: B0, B1, B2
	count := 0
	current := heap.head
	degrees := []int{}

	for current != nil {
		count++
		degrees = append(degrees, current.degree)
		current = current.sibling
	}

	if count != 3 {
		b.Errorf("Con 7 elementi, attesi 3 alberi, ottenuti %d", count)
	}

	// Verifica che i gradi siano 0, 1, 2
	expectedDegrees := []int{0, 1, 2}
	for i, degree := range degrees {
		if degree != expectedDegrees[i] {
			b.Errorf("Albero %d: grado atteso %d, ottenuto %d",
				i, expectedDegrees[i], degree)
		}
	}
}

func TestBinomialTreeStructurePowerOfTwo(b *testing.T) {
	heap := CreateBinomialHeap()

	// Insert 8 elementi (8 = 1000 in binario = solo B3)
	for i := 0; i < 8; i++ {
		heap.Insert(i, float64(8-i))
	}

	// Dovrebbe esserci un solo albero binomiale di grado 3
	count := 0
	current := heap.head

	for current != nil {
		count++
		if current.degree != 3 {
			b.Errorf("Con 8 elementi, atteso grado 3, ottenuto %d", current.degree)
		}
		current = current.sibling
	}

	if count != 1 {
		b.Errorf("Con 8 elementi (potenza di 2), atteso 1 albero, ottenuti %d", count)
	}
}
