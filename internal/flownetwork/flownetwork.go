package flownetwork

// FlowEdge rappresenta un arco nel grafo, sia iniziale (Capacity) che dei residui (Flow). Rappresenta un arco (From -> To) e ha l'indice dell'arco inverso (To -> From).
type FlowEdge struct {
	From     int // Nodo sorgente
	To       int // Nodo destinazione
	Capacity int // Capacità dell'arco
	Flow     int // Flusso corrente sull'arco
	Reverse  int // Indice dell'arco inverso
}

// FlowNetwork rappresenta una rete di flusso.
type FlowNetwork struct {
	N        int           // Numero di nodi
	OutStars [][]*FlowEdge // Lista di adiacenza: OutStars[i] = archi uscenti da i
	InStars  [][]*FlowEdge // Lista di adiacenza: InStars[i] = archi entranti in i
	Source   int           // Nodo sorgente
	Sink     int           // Nodo pozzo
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

	outStars := make([][]*FlowEdge, n)
	inStars := make([][]*FlowEdge, n)
	for i := 0; i < n; i++ {
		outStars[i] = []*FlowEdge{}
		inStars[i] = []*FlowEdge{}
	}

	return &FlowNetwork{
		N:        n,
		Source:   source,
		Sink:     sink,
		OutStars: outStars,
		InStars:  inStars,
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
	directArc := FlowEdge{
		From:     from,
		To:       to,
		Capacity: capacity,
		Flow:     0,
		Reverse:  len(fn.OutStars[to]),
	}
	inverseArc := FlowEdge{
		From:     to,
		To:       from,
		Capacity: 0,
		Flow:     0,
		Reverse:  len(fn.OutStars[from]) - 1,
	}
	fn.OutStars[from] = append(fn.OutStars[from], &directArc)
	fn.OutStars[to] = append(fn.OutStars[to], &inverseArc)
	fn.InStars[from] = append(fn.InStars[from], &inverseArc)
	fn.InStars[to] = append(fn.InStars[to], &directArc)
}

// PushFlow invia δ unità di flusso lungo l'arco che va da from a
func (fn *FlowNetwork) PushFlow(from int, to int, delta int) {
	edgeIndex := -1
	for i, edge := range fn.OutStars[from] {
		if edge.To == to {
			edgeIndex = i
		}
	}
	edge := fn.OutStars[from][edgeIndex]
	reverseEdge := fn.OutStars[edge.To][edge.Reverse]

	residual := edge.Capacity - edge.Flow
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}
	edge.Flow += delta
	reverseEdge.Flow -= delta
}

func (fn *FlowNetwork) PushFlowWithIndex(from int, index int, delta int) {
	edge := fn.OutStars[from][index]
	reverseEdge := fn.OutStars[edge.To][edge.Reverse]

	residual := edge.Capacity - edge.Flow
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}
	edge.Flow += delta
	reverseEdge.Flow -= delta
}

// Reset azzera il flusso corrente di ogni arco, riportandolo allo stato di origine.
func (fn *FlowNetwork) Reset() {
	for i := 0; i < fn.N; i++ {
		for j := range fn.OutStars[i] {
			fn.OutStars[i][j].Flow = 0
		}
	}
}

// GetMaxFlowValue calcola il valore del flusso finale da source a sink
func (fn *FlowNetwork) GetMaxFlowValue() int {
	maxFlow := 0
	for _, edge := range fn.OutStars[fn.Source] {
		maxFlow += edge.Flow
	}
	return maxFlow
}

func (fn *FlowNetwork) validateArcs(from int, to int) {
	if from < 0 || from >= fn.N {
		panic("Nodo 'from' non valido")
	}
	if to < 0 || to >= fn.N {
		panic("Nodo 'to' non valido")
	}
}
