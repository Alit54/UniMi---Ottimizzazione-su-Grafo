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

func (fn *FlowNetwork) AddEdge(from int, to int, capacity int) {
	fn.validateArcs(from, to)
	if capacity < 0 {
		panic("Capacità deve essere non negativa")
	}
	if from == to {
		panic("Self-loop non ammessi")
	}
	fn.Capacity[from][to] = capacity
}

// PushFlow invia δ unità di flusso lungo l'arco (from → to)
func (fn *FlowNetwork) PushFlow(from int, to int, delta int) {
	fn.validateArcs(from, to)
	if delta < 0 {
		panic("Delta deve essere non negativo")
	}
	residual := fn.Capacity[from][to] - fn.Flow[from][to]
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}
	fn.Flow[from][to] += delta
	fn.Flow[to][from] -= delta
}

func (fn *FlowNetwork) validateArcs(from int, to int) {
	if from < 0 || from >= fn.N {
		panic("Nodo 'from' non valido")
	}
	if to < 0 || to >= fn.N {
		panic("Nodo 'to' non valido")
	}
}
