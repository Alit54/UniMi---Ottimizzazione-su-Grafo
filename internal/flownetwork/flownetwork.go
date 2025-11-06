package flownetwork

// FlowNetwork rappresenta una rete di flusso.
type FlowNetwork struct {
	N        int     // Numero di nodi
	Capacity [][]int // Matrice di Adiacenza delle capacità
	Flow     [][]int // Matrice di Adiacenza del grafo dei residui
	Source   int     // Nodo sorgente
	Sink     int     // Nodo target (pozzo)
}

// NewFlowNetwork crea una rete di flusso con n nodi, senza archi
func NewFlowNetwork(n int, source int, sink int) *FlowNetwork {
	if source < 0 || source >= n {
		panic("Source deve essere tra 0 e n-1")
	}
	if sink < 0 || sink >= n {
		panic("Sink deve essere tra 0 e n-1")
	}
	if source == sink {
		panic("Source e Sink devono essere nodi diversi")
	}

	capacity := make([][]int, n)
	flow := make([][]int, n)
	for i := 0; i < n; i++ {
		capacity[i] = make([]int, n)
		flow[i] = make([]int, n)
	}

	return &FlowNetwork{
		N:        n,
		Capacity: capacity,
		Flow:     flow,
		Source:   source,
		Sink:     sink,
	}
}
