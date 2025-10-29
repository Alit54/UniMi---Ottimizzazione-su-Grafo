package heap

type Heap interface {
	/*
		Inserisce un elemento nello heap.
		Precondizione: non deve esistere nodeID già all'interno dello heap
	*/
	Insert(nodeID int, distance float64)
	/*
		Estrae il nodo con distanza minima.
		Precondizione: heap non deve essere vuoto
	*/
	ExtractMin() (nodeID int, distance float64)
	/*
		Diminusce il valore nello heap di un nodo esistente
		Precondizione: la nuova distanza deve essere minore e il nodo deve esistere nello heap
	*/
	DecreaseKey(nodeID int, newDistance float64)
	/*
		Stampa la struttura dello heap
	*/
	PrintHeap()
	/*
		Restituisce il numero di elementi dello heap
	*/
	Size() int
	/*
		Controlla che lo heap sia vuoto
	*/
	IsEmpty() bool
}
