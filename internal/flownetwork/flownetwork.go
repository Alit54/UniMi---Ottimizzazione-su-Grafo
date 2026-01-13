package flownetwork

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FlowEdge struct {
	From     int
	To       int
	Capacity int
	Flow     int
	Reverse  int
}

type FlowNetwork struct {
	N        int
	Arcs     int
	OutStars [][]*FlowEdge
	Source   int
	Sink     int
}

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
	for i := 0; i < n; i++ {
		outStars[i] = []*FlowEdge{}
	}

	return &FlowNetwork{
		N:        n,
		Source:   source,
		Sink:     sink,
		OutStars: outStars,
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

	idxFrom := len(fn.OutStars[from])
	idxTo := len(fn.OutStars[to])

	forward := &FlowEdge{
		From:     from,
		To:       to,
		Capacity: capacity,
		Flow:     0,
		Reverse:  idxTo,
	}

	backward := &FlowEdge{
		From:     to,
		To:       from,
		Capacity: 0,
		Flow:     0,
		Reverse:  idxFrom,
	}

	fn.OutStars[from] = append(fn.OutStars[from], forward)
	fn.OutStars[to] = append(fn.OutStars[to], backward)
	fn.Arcs++
}

func (fn *FlowNetwork) PushFlowWithIndex(from int, index int, delta int) {
	edge := fn.OutStars[from][index]
	residual := edge.Capacity - edge.Flow
	if delta > residual {
		panic("Delta supera la capacità residua disponibile")
	}

	edge.Flow += delta
	fn.OutStars[edge.To][edge.Reverse].Flow -= delta
}

func (fn *FlowNetwork) Reset() {
	for i := 0; i < fn.N; i++ {
		for j := range fn.OutStars[i] {
			fn.OutStars[i][j].Flow = 0
		}
	}
}

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
