package heap

import (
	"math/rand"
	"testing"
)

// ============================================ //
//          TEST COSTRUTTORE                    //
// ============================================ //

func TestCreateFibonacciHeap(t *testing.T) {
	values := []float64{50.0, 30.0, 70.0, 10.0, 40.0}
	heap := CreateFibonacciHeap(values...)

	if heap.Size() != 5 {
		t.Errorf("Size attesa 5, ottenuta %d", heap.Size())
	}

	// Verifica che il minimo sia corretto
	if heap.min.distance != 10.0 {
		t.Errorf("Minimo dovrebbe essere 10.0, ottenuto %.1f", heap.min.distance)
	}
}

func TestCreateFibonacciHeapEmpty(t *testing.T) {
	heap := CreateFibonacciHeap()

	if heap.Size() != 0 {
		t.Error("Heap vuoto dovrebbe avere size 0")
	}

	if !heap.IsEmpty() {
		t.Error("Heap vuoto dovrebbe risultare isEmpty")
	}
}

func TestCreateFibonacciHeapSingle(t *testing.T) {
	heap := CreateFibonacciHeap(42.0)

	if heap.Size() != 1 {
		t.Error("Heap con singolo elemento dovrebbe avere size 1")
	}

	nodeID, dist := heap.ExtractMin()
	if nodeID != 0 || dist != 42.0 {
		t.Errorf("Atteso (0, 42.0), ottenuto (%d, %.1f)", nodeID, dist)
	}
}

func TestCreateFibonacciHeapOrdered(t *testing.T) {
	// Test con valori già ordinati
	heap := CreateFibonacciHeap(10.0, 20.0, 30.0, 40.0, 50.0)

	if heap.min.distance != 10.0 {
		t.Error("Heap dovrebbe mantenere 10.0 come minimo")
	}
}

func TestCreateFibonacciHeapReverse(t *testing.T) {
	// Test con valori in ordine inverso
	heap := CreateFibonacciHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	if heap.min.distance != 10.0 {
		t.Error("Heap dovrebbe avere 10.0 come minimo")
	}
}

// ============================================ //
//          TEST INSERT                         //
// ============================================ //

func TestFibonacciInsertSingle(t *testing.T) {
	heap := CreateFibonacciHeap()

	heap.Insert(1, 10.0)

	if heap.Size() != 1 {
		t.Errorf("Size attesa 1, ottenuta %d", heap.Size())
	}

	if heap.min.nodeID != 1 || heap.min.distance != 10.0 {
		t.Error("Elemento inserito non corretto")
	}
}

func TestFibonacciInsertMultiple(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciInsertUpdatePosition(t *testing.T) {
	heap := CreateFibonacciHeap()

	heap.Insert(10, 100.0)
	heap.Insert(20, 200.0)
	heap.Insert(30, 300.0)

	// Verifica che position map sia aggiornata
	for nodeID, node := range heap.position {
		if node.nodeID != nodeID {
			t.Errorf("Position map inconsistente per nodeID=%d", nodeID)
		}
	}
}

func TestFibonacciInsertUpdatesMin(t *testing.T) {
	heap := CreateFibonacciHeap()

	heap.Insert(1, 50.0)
	if heap.min.distance != 50.0 {
		t.Error("Min non aggiornato dopo primo insert")
	}

	heap.Insert(2, 30.0)
	if heap.min.distance != 30.0 {
		t.Error("Min dovrebbe essere 30.0")
	}

	heap.Insert(3, 40.0)
	if heap.min.distance != 30.0 {
		t.Error("Min non dovrebbe cambiare")
	}

	heap.Insert(4, 20.0)
	if heap.min.distance != 20.0 {
		t.Error("Min dovrebbe essere 20.0")
	}
}

// ============================================ //
//          TEST EXTRACTMIN                     //
// ============================================ //

func TestFibonacciExtractMinOrder(t *testing.T) {
	heap := CreateFibonacciHeap(100.0, 50.0, 150.0, 30.0, 70.0, 10.0)

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

func TestFibonacciExtractMinCorrectNodeIDs(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciExtractMinUpdatesPosition(t *testing.T) {
	heap := CreateFibonacciHeap(50.0, 30.0, 70.0, 10.0, 40.0)

	initialSize := heap.Size()
	minNodeID := heap.min.nodeID

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

func TestFibonacciExtractMinSingle(t *testing.T) {
	heap := CreateFibonacciHeap(42.0)

	nodeID, dist := heap.ExtractMin()

	if nodeID != 0 || dist != 42.0 {
		t.Errorf("ExtractMin su singolo elemento: atteso (0, 42.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}

	if !heap.IsEmpty() {
		t.Error("Heap dovrebbe essere vuoto dopo ExtractMin su singolo elemento")
	}
}

func TestFibonacciExtractMinEmpty(t *testing.T) {
	heap := CreateFibonacciHeap()

	// ExtractMin su heap vuoto non dovrebbe crashare
	nodeID, dist := heap.ExtractMin()

	// Dovrebbe restituire valori zero
	if nodeID != 0 || dist != 0.0 {
		t.Errorf("ExtractMin su heap vuoto: atteso (0, 0.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestFibonacciExtractMinAll(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciExtractMinConsolidate(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert molti elementi (tutti nella root list inizialmente)
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64(20-i))
	}

	// Prima di ExtractMin, ci dovrebbero essere 20 alberi nella root list
	// Dopo ExtractMin, consolidate dovrebbe ridurre il numero

	initialMin := heap.min.distance
	heap.ExtractMin()

	// Verifica che il nuovo minimo sia corretto
	if heap.min.distance <= initialMin {
		t.Error("Nuovo minimo dovrebbe essere maggiore del precedente")
	}

	// Verifica che l'heap sia ancora valido
	verifyFibonacciHeapProperty(t, heap)
}

// ============================================ //
//          TEST DECREASEKEY                    //
// ============================================ //

func TestFibonacciDecreaseKeyBasic(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciDecreaseKeyToMinimum(t *testing.T) {
	heap := CreateFibonacciHeap(50.0, 40.0, 30.0, 20.0, 10.0)

	// Diminuisci il nodo 4 (distance=10.0) a 5.0
	heap.DecreaseKey(4, 5.0)

	nodeID, dist := heap.ExtractMin()
	if nodeID != 4 || dist != 5.0 {
		t.Errorf("DecreaseKey a nuovo minimo: atteso (4, 5.0), ottenuto (%d, %.1f)",
			nodeID, dist)
	}
}

func TestFibonacciDecreaseKeyMultiple(t *testing.T) {
	heap := CreateFibonacciHeap()
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

func TestFibonacciDecreaseKeyAfterExtractMin(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert e crea struttura complessa
	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i*10))
	}

	// ExtractMin per consolidare
	heap.ExtractMin()

	// Ora DecreaseKey su un nodo interno
	heap.DecreaseKey(5, 5.0)

	// Il nodo dovrebbe essere stato tagliato e spostato nella root list
	node := heap.position[5]
	if node.parent != nil {
		t.Error("Dopo DecreaseKey che viola heap property, il nodo dovrebbe essere nella root list")
	}

	if node.distance != 5.0 {
		t.Error("Distanza non aggiornata correttamente")
	}
}

func TestFibonacciDecreaseKeyCascadingCut(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Crea una struttura che richiederà cascading cut
	for i := 0; i < 16; i++ {
		heap.Insert(i, float64(i))
	}

	// ExtractMin per consolidare e creare una struttura ad albero
	heap.ExtractMin()

	// DecreaseKey su nodi interni per triggerare cascading cut
	heap.DecreaseKey(10, -10.0)
	heap.DecreaseKey(5, -20.0)

	// Verifica che l'heap sia ancora valido
	verifyFibonacciHeapProperty(t, heap)

	// Il minimo dovrebbe essere -20.0
	if heap.min.distance != -20.0 {
		t.Errorf("Min dovrebbe essere -20.0, ottenuto %.1f", heap.min.distance)
	}
}

func TestFibonacciDecreaseKeyUpdatesMin(t *testing.T) {
	heap := CreateFibonacciHeap()

	heap.Insert(1, 100.0)
	heap.Insert(2, 200.0)
	heap.Insert(3, 300.0)

	initialMin := heap.min.distance

	// DecreaseKey a un valore maggiore del minimo corrente
	heap.DecreaseKey(3, 150.0)
	if heap.min.distance != initialMin {
		t.Error("Min non dovrebbe cambiare")
	}

	// DecreaseKey a un valore minore del minimo corrente
	heap.DecreaseKey(2, 50.0)
	if heap.min.distance != 50.0 {
		t.Errorf("Min dovrebbe essere 50.0, ottenuto %.1f", heap.min.distance)
	}
}

// ============================================ //
//      TEST OPERAZIONI MISTE (REALISTICI)     //
// ============================================ //

func TestFibonacciMixedOperationsDijkstraPattern(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciMixedOperationsAlternating(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciMixedOperationsStressTest(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciHeavyDecreaseKey(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Pattern con molti DecreaseKey (tipico di Dijkstra)
	for i := 0; i < 50; i++ {
		heap.Insert(i, float64(i*100))
	}

	// ExtractMin per consolidare
	heap.ExtractMin()

	// Molti DecreaseKey
	for i := 10; i < 40; i++ {
		heap.DecreaseKey(i, float64(i))
	}

	// Verifica ordine
	lastDist := -1.0
	for !heap.IsEmpty() {
		_, dist := heap.ExtractMin()
		if dist < lastDist {
			t.Errorf("Ordine violato dopo molti DecreaseKey")
		}
		lastDist = dist
	}
}

// ============================================ //
//          TEST EDGE CASES                     //
// ============================================ //

func TestFibonacciFloatingPointPrecision(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciZeroDistances(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciNegativeDistances(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciLargeHeap(t *testing.T) {
	heap := CreateFibonacciHeap()
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

func TestFibonacciDuplicateDistances(t *testing.T) {
	heap := CreateFibonacciHeap(50.0, 30.0, 50.0, 30.0, 50.0)

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

func TestFibonacciAllSameDistances(t *testing.T) {
	heap := CreateFibonacciHeap(7.0, 7.0, 7.0, 7.0)

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
//          TEST SIZE, ISEMPTY, CONTAINS        //
// ============================================ //

func TestFibonacciSizeTracking(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciIsEmpty(t *testing.T) {
	heap := CreateFibonacciHeap()

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

func TestFibonacciPositionMapConsistency(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64(i*5))
	}

	// Verifica consistenza
	for nodeID, node := range heap.position {
		if node.nodeID != nodeID {
			t.Errorf("Position map inconsistente: nodeID=%d, ma node.nodeID=%d",
				nodeID, node.nodeID)
		}
	}

	// Dopo DecreaseKey
	heap.DecreaseKey(15, 2.0)

	for nodeID, node := range heap.position {
		if node.nodeID != nodeID {
			t.Error("Position map inconsistente dopo DecreaseKey")
		}
	}

	// Dopo ExtractMin
	heap.ExtractMin()

	for nodeID, node := range heap.position {
		if node.nodeID != nodeID {
			t.Error("Position map inconsistente dopo ExtractMin")
		}
	}
}

func TestFibonacciPositionMapSize(t *testing.T) {
	heap := CreateFibonacciHeap()

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
//            HELPER FUNCTIONS                  //
// ============================================ //

// verifyFibonacciHeapProperty verifica la proprietà min-heap ricorsivamente
func verifyFibonacciHeapProperty(t *testing.T, heap *FibonacciHeap) {
	if heap.min == nil {
		return
	}

	// Verifica tutti gli alberi nella root list
	current := heap.min
	for {
		verifyFibonacciTree(t, current)
		current = current.right
		if current == heap.min {
			break
		}
	}
}

// verifyFibonacciTree verifica la proprietà heap in un albero
func verifyFibonacciTree(t *testing.T, node *fibonacciNode) {
	if node == nil {
		return
	}

	// Verifica che ogni figlio abbia distanza >= padre
	if node.child != nil {
		child := node.child
		for {
			if child.distance < node.distance {
				t.Errorf("Proprietà heap violata: figlio (%.1f) < padre (%.1f)",
					child.distance, node.distance)
			}
			verifyFibonacciTree(t, child)
			child = child.right
			if child == node.child {
				break
			}
		}
	}
}

// ============================================ //
//              BENCHMARK                       //
// ============================================ //

func BenchmarkFibonacciInsert(b *testing.B) {
	heap := CreateFibonacciHeap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(i))
	}
}

func BenchmarkFibonacciInsertRandom(b *testing.B) {
	heap := CreateFibonacciHeap()
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, rand.Float64()*1000)
	}
}

func BenchmarkFibonacciExtractMin(b *testing.B) {
	// Prepara heap con b.N elementi
	heap := CreateFibonacciHeap()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, float64(b.N-i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.ExtractMin()
	}
}

func BenchmarkFibonacciDecreaseKey(b *testing.B) {
	// Prepara heap con 10000 elementi
	heap := CreateFibonacciHeap()
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

func BenchmarkFibonacciMixedDijkstraPattern(b *testing.B) {
	heap := CreateFibonacciHeap()
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
				if _, exists := heap.position[targetID]; exists {
					heap.DecreaseKey(targetID, rand.Float64()*100)
				}
			}
		}
	}
}

func BenchmarkFibonacciCreateHeap(b *testing.B) {
	n := 1000
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = float64(n - i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateFibonacciHeap(values...)
	}
}

func BenchmarkFibonacciSize(b *testing.B) {
	heap := CreateFibonacciHeap()
	for i := 0; i < 1000; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.Size()
	}
}

func BenchmarkFibonacciIsEmpty(b *testing.B) {
	heap := CreateFibonacciHeap()
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

func BenchmarkFibonacciInsert_100(b *testing.B) {
	benchmarkFibonacciInsertN(b, 100)
}

func BenchmarkFibonacciInsert_1000(b *testing.B) {
	benchmarkFibonacciInsertN(b, 1000)
}

func BenchmarkFibonacciInsert_10000(b *testing.B) {
	benchmarkFibonacciInsertN(b, 10000)
}

func benchmarkFibonacciInsertN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
	}
}

func BenchmarkFibonacciExtractMin_100(b *testing.B) {
	benchmarkFibonacciExtractMinN(b, 100)
}

func BenchmarkFibonacciExtractMin_1000(b *testing.B) {
	benchmarkFibonacciExtractMinN(b, 1000)
}

func BenchmarkFibonacciExtractMin_10000(b *testing.B) {
	benchmarkFibonacciExtractMinN(b, 10000)
}

func benchmarkFibonacciExtractMinN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(n-j))
		}
		b.StartTimer()

		for j := 0; j < n; j++ {
			heap.ExtractMin()
		}
	}
}

func BenchmarkFibonacciDecreaseKey_100(b *testing.B) {
	benchmarkFibonacciDecreaseKeyN(b, 100)
}

func BenchmarkFibonacciDecreaseKey_1000(b *testing.B) {
	benchmarkFibonacciDecreaseKeyN(b, 1000)
}

func BenchmarkFibonacciDecreaseKey_10000(b *testing.B) {
	benchmarkFibonacciDecreaseKeyN(b, 10000)
}

func benchmarkFibonacciDecreaseKeyN(b *testing.B, n int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
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

func BenchmarkFibonacciDijkstraSmallGraph(b *testing.B) {
	// Simula Dijkstra su grafo piccolo (100 nodi)
	benchmarkFibonacciDijkstraPattern(b, 100, 300)
}

func BenchmarkFibonacciDijkstraMediumGraph(b *testing.B) {
	// Simula Dijkstra su grafo medio (1000 nodi)
	benchmarkFibonacciDijkstraPattern(b, 1000, 5000)
}

func BenchmarkFibonacciDijkstraLargeGraph(b *testing.B) {
	// Simula Dijkstra su grafo grande (10000 nodi)
	benchmarkFibonacciDijkstraPattern(b, 10000, 50000)
}

func benchmarkFibonacciDijkstraPattern(b *testing.B, numNodes int, numEdges int) {
	rand.Seed(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap := CreateFibonacciHeap()

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

func BenchmarkFibonacciWorstCaseInsert(b *testing.B) {
	// Worst case: inserimenti in ordine crescente
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkFibonacciWorstCaseExtractMin(b *testing.B) {
	// Worst case: molti alberi nella root list
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		for j := 1000; j > 0; j-- {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			heap.ExtractMin()
		}
	}
}

func BenchmarkFibonacciWorstCaseDecreaseKey(b *testing.B) {
	// Worst case: DecreaseKey che triggerano cascading cut
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		// ExtractMin per creare struttura ad albero
		heap.ExtractMin()
		b.StartTimer()

		// DecreaseKey che causano cascading cut
		for j := 0; j < 100; j++ {
			targetID := (j * 7) % n
			if _, exists := heap.position[targetID]; exists {
				heap.DecreaseKey(targetID, -float64(j))
			}
		}
	}
}

// ============================================ //
//        BENCHMARK BEST CASE                   //
// ============================================ //

func BenchmarkFibonacciBestCaseInsert(b *testing.B) {
	// Best case: inserimenti in ordine decrescente
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		b.StartTimer()

		for j := 1000; j > 0; j-- {
			heap.Insert(j, float64(j))
		}
	}
}

func BenchmarkFibonacciBestCaseDecreaseKey(b *testing.B) {
	// Best case: DecreaseKey su nodi nella root list
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := CreateFibonacciHeap()
		n := 1000
		for j := 0; j < n; j++ {
			heap.Insert(j, float64(j))
		}
		b.StartTimer()

		// DecreaseKey su nodi nella root list (no cut necessario)
		for j := 0; j < 100; j++ {
			heap.DecreaseKey(j%n, float64(j%n)-0.1)
		}
	}
}

// ============================================ //
//     BENCHMARK MEMORY & POSITION MAP         //
// ============================================ //

func BenchmarkFibonacciPositionLookup(b *testing.B) {
	heap := CreateFibonacciHeap()
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

// ============================================ //
//   BENCHMARK VANTAGGIO FIBONACCI: O(1) OPS   //
// ============================================ //

func BenchmarkFibonacciInsertO1(b *testing.B) {
	// Dimostra che Insert è O(1) anche con heap grande
	heap := CreateFibonacciHeap()

	// Pre-popola con molti elementi
	for i := 0; i < 100000; i++ {
		heap.Insert(i, float64(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(100000+i, float64(100000+i))
	}
}

func BenchmarkFibonacciDecreaseKeyO1(b *testing.B) {
	// Dimostra che DecreaseKey è O(1) ammortizzato
	heap := CreateFibonacciHeap()
	n := 10000

	// Pre-popola
	for i := 0; i < n; i++ {
		heap.Insert(i, float64(i*10))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := i % n
		heap.DecreaseKey(nodeID, float64(i%100))
	}
}

// ============================================ //
//     TEST SPECIFICI FIBONACCI HEAP           //
// ============================================ //

func TestFibonacciRootListStructure(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert 5 elementi - tutti dovrebbero essere nella root list inizialmente
	for i := 0; i < 5; i++ {
		heap.Insert(i, float64(5-i))
	}

	// Conta nodi nella root list
	count := 0
	if heap.min != nil {
		current := heap.min
		for {
			count++
			current = current.right
			if current == heap.min {
				break
			}
		}
	}

	if count != 5 {
		t.Errorf("Dopo 5 Insert, attesi 5 nodi in root list, ottenuti %d", count)
	}
}

func TestFibonacciConsolidateReducesTrees(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert molti elementi
	for i := 0; i < 20; i++ {
		heap.Insert(i, float64(20-i))
	}

	// Conta alberi prima di ExtractMin
	countBefore := 0
	current := heap.min
	for {
		countBefore++
		current = current.right
		if current == heap.min {
			break
		}
	}

	// ExtractMin triggerera consolidate
	heap.ExtractMin()

	// Conta alberi dopo ExtractMin
	countAfter := 0
	if heap.min != nil {
		current = heap.min
		for {
			countAfter++
			current = current.right
			if current == heap.min {
				break
			}
		}
	}

	// Dopo consolidate dovrebbero esserci meno alberi
	if countAfter >= countBefore {
		t.Errorf("Consolidate non ha ridotto il numero di alberi: prima=%d, dopo=%d",
			countBefore, countAfter)
	}
}

func TestFibonacciMarkedNodes(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Crea una struttura ad albero
	for i := 0; i < 16; i++ {
		heap.Insert(i, float64(i))
	}

	// ExtractMin per consolidare
	heap.ExtractMin()

	// DecreaseKey per triggerare cut
	heap.DecreaseKey(10, -10.0)

	// Il parent del nodo tagliato dovrebbe essere marked
	// (non possiamo verificarlo direttamente senza esportare marked, ma possiamo verificare che l'heap funzioni)
	verifyFibonacciHeapProperty(t, heap)

	// Secondo DecreaseKey dovrebbe triggerare cascading cut
	heap.DecreaseKey(5, -20.0)

	verifyFibonacciHeapProperty(t, heap)
}

func TestFibonacciCircularListIntegrity(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert elementi
	for i := 0; i < 10; i++ {
		heap.Insert(i, float64(i))
	}

	// Verifica integrità lista circolare
	if heap.min != nil {
		current := heap.min
		visited := make(map[*fibonacciNode]bool)

		for {
			if visited[current] {
				// Abbiamo fatto il giro completo
				break
			}
			visited[current] = true

			// Verifica puntatori left/right
			if current.right.left != current {
				t.Error("Puntatori left/right non consistenti")
			}

			current = current.right
		}

		if len(visited) == 0 {
			t.Error("Lista circolare vuota ma heap.min non nil")
		}
	}
}

func TestFibonacciDegreeConsistency(t *testing.T) {
	heap := CreateFibonacciHeap()

	// Insert e ExtractMin per creare struttura
	for i := 0; i < 16; i++ {
		heap.Insert(i, float64(16-i))
	}

	heap.ExtractMin()

	// Verifica che il degree di ogni nodo corrisponda al numero di figli
	if heap.min != nil {
		current := heap.min
		for {
			// Conta figli
			childCount := 0
			if current.child != nil {
				child := current.child
				for {
					childCount++
					child = child.right
					if child == current.child {
						break
					}
				}
			}

			if current.degree != childCount {
				t.Errorf("Nodo %d: degree=%d, ma ha %d figli",
					current.nodeID, current.degree, childCount)
			}

			current = current.right
			if current == heap.min {
				break
			}
		}
	}
}
