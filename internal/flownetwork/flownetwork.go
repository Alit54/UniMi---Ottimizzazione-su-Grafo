package flownetwork

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	fmt.Println(n, sink, source)
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

func NewNetworkFromDIMACS(path string) *FlowNetwork {
	file, _ := os.Open(path)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	nNodes := 0
	source := 0
	sink := 0
	fn := &FlowNetwork{}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if line == "" || parts[0] == "c" {
			continue
		}
		if parts[0] == "p" {
			nNodes, _ = strconv.Atoi(parts[2])
			continue
		}
		if parts[0] == "n" {
			if parts[2] == "s" {
				source, _ = strconv.Atoi(parts[1])
			}
			if parts[2] == "t" {
				sink, _ = strconv.Atoi(parts[1])
				fn = NewFlowNetwork(nNodes, source-1, sink-1)
			}
			continue
		}
		if parts[0] == "a" {
			from, _ := strconv.Atoi(parts[1])
			to, _ := strconv.Atoi(parts[2])
			if fn.arcExists(from, to) {
				continue
			}
			capacity, _ := strconv.Atoi(parts[3])
			if capacity > 0 {
				fn.AddEdge(from-1, to-1, capacity)
			}
		}
	}
	return fn
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
	}
	fn.OutStars[from] = append(fn.OutStars[from], &directArc)
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
	residual := edge.Capacity - edge.Flow
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}
	edge.Flow += delta
}

func (fn *FlowNetwork) PushFlowWithIndex(from int, index int, delta int) {
	edge := fn.OutStars[from][index]

	residual := edge.Capacity - edge.Flow
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}
	edge.Flow += delta
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

func (fn *FlowNetwork) arcExists(from int, to int) bool {
	for _, arc := range fn.OutStars[from] {
		if arc.To == to {
			return true
		}
	}
	return false
}
