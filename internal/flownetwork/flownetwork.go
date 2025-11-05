package flownetwork

// FlowNetwork rappresenta una rete di flusso.
type FlowNetwork struct {
	N        int     // Numero di nodi
	Capacity [][]int // Matrice di Adiacenza delle capacità
	Flow     [][]int // Matrice di Adiacenza del grafo dei residui
	Source   int     // Nodo sorgente
	Sink     int     // Nodo target (pozzo)
}
