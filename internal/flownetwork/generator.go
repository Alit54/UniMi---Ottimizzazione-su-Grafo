package flownetwork // o package maxflow, a seconda della tua struttura

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

// GenerateDIMACS genera un file .max ottimizzato per stressare specifici algoritmi.
// mode:
// 0 = Killer (Denso + Cicli + Capacità alte)
// 1 = Anti CS / Scaling (Capacità "bitwise malicious")
// 2 = Anti SAP (Grafo lungo e stretto)
// 3 = Anti Dinic (Denso con molti Back-Edges)
func GenerateDIMACS(outputPath string, nNodes int, density float64, maxCap int, source int, sink int, mode int) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Struttura interna per gli archi
	type Edge struct {
		u, v, cap int
	}
	var edges []Edge
	existingEdges := make(map[string]bool)

	// Helper per aggiungere archi unici
	addEdge := func(u, v, cap int) {
		if u == v {
			return
		} // No self-loops
		key := fmt.Sprintf("%d-%d", u, v)
		if !existingEdges[key] {
			edges = append(edges, Edge{u, v, cap})
			existingEdges[key] = true
		}
	}

	fmt.Printf("Generazione Grafo Mode %d per N=%d e densità %f\n", mode, nNodes, density)

	switch mode {

	// -----------------------------------------------------------------------
	// MODE 0: KILLER (Tutti insieme appassionatamente)
	// Combinazione di profondità, densità e cicli.
	// -----------------------------------------------------------------------
	case 0:
		// 1. Spina dorsale lunga (Anti-SAP)
		backboneNodes := nNodes / 2
		for i := 0; i < backboneNodes-1; i++ {
			// Capacità variabile
			cap := 1 + rng.Intn(maxCap)
			addEdge(i, i+1, cap)
		}
		// Collega fine backbone a Sink
		addEdge(backboneNodes-1, sink, maxCap)

		// 2. Nuvola Densa attorno (Anti-Dinic e Anti-Scaling)
		targetEdges := int(float64(nNodes*nNodes) * density)
		attempts := 0
		maxAttempts := targetEdges * 2

		for len(edges) < targetEdges && attempts < maxAttempts {
			u := rng.Intn(nNodes)
			v := rng.Intn(nNodes)

			if u == v || v == source || u == sink {
				attempts++
				continue
			}

			// Capacità alta e randomica
			cap := 1 + rng.Intn(maxCap)
			addEdge(u, v, cap)
			attempts++
		}

	// -----------------------------------------------------------------------
	// MODE 1: ANTI CS / SCALING
	// Distribuzione capacità bitwise per forzare tutte le fasi di scaling.
	// -----------------------------------------------------------------------
	case 1:
		targetEdges := int(float64(nNodes*nNodes) * density)
		attempts := 0

		for len(edges) < targetEdges && attempts < targetEdges*2 {
			u := rng.Intn(nNodes)
			v := rng.Intn(nNodes)
			if u == v || v == source || u == sink {
				attempts++
				continue
			}

			// Genera capacità che "accendono" bit specifici
			power := rng.Intn(int(math.Log2(float64(maxCap))))
			cap := int(math.Pow(2, float64(power)))

			// Aggiunge rumore su altri bit
			if rng.Float64() > 0.5 && power > 1 {
				cap += int(math.Pow(2, float64(power/2)))
			}
			if cap > maxCap {
				cap = maxCap
			}
			if cap < 1 {
				cap = 1
			}

			addEdge(u, v, cap)
			attempts++
		}
		// Garantire connettività minima S->T
		addEdge(source, 1, maxCap)
		addEdge(nNodes-2, sink, maxCap)

	// -----------------------------------------------------------------------
	// MODE 2: ANTI SAP (Lungo e Stretto)
	// Minimizza la larghezza per massimizzare la lunghezza del cammino.
	// -----------------------------------------------------------------------
	case 2:
		width := 2 // Molto stretto
		if nNodes > 50 {
			width = 3
		}

		// Calcolo layer
		layers := (nNodes - 2) / width
		if layers < 1 {
			layers = 1
		}

		currentNode := 0
		layerNodes := make([][]int, layers)

		// Assegnazione nodi ai layer (esclusi S e T)
		nodesAssigned := 0
		for l := 0; l < layers; l++ {
			for w := 0; w < width; w++ {
				// Saltiamo S e T se capitano nel conteggio lineare
				for currentNode == source || currentNode == sink {
					currentNode++
				}
				if currentNode < nNodes {
					layerNodes[l] = append(layerNodes[l], currentNode)
					currentNode++
					nodesAssigned++
				}
			}
		}

		// Source -> Primo Layer
		if len(layerNodes) > 0 {
			for _, node := range layerNodes[0] {
				addEdge(source, node, maxCap)
			}
		}

		// Layer -> Layer (Pipeline)
		for l := 0; l < len(layerNodes)-1; l++ {
			for _, u := range layerNodes[l] {
				for _, v := range layerNodes[l+1] {
					addEdge(u, v, maxCap) // Capacità massima per evitare saturazione rapida
				}
			}
		}

		// Ultimo Layer -> Sink
		if len(layerNodes) > 0 {
			for _, node := range layerNodes[len(layerNodes)-1] {
				addEdge(node, sink, maxCap)
			}
		}

	// -----------------------------------------------------------------------
	// MODE 3: ANTI DINIC (Denso con Back-Edges)
	// Back-edges invalidano il grafo a livelli di Dinic.
	// -----------------------------------------------------------------------
	case 3:
		// Struttura densa base
		for u := 0; u < nNodes; u++ {
			for v := 0; v < nNodes; v++ {
				if u == v || u == sink || v == source {
					continue
				}

				if rng.Float64() < density {
					cap := 1 + rng.Intn(maxCap)

					// Logica Anti-Dinic: Back Edge (u > v)
					// Capacità piccola per costringere ricalcoli frequenti
					if u > v {
						cap = 1 + rng.Intn(maxCap/100+1)
					}
					addEdge(u, v, cap)
				}
			}
		}
		// Connettività garantita
		addEdge(source, 1, maxCap)
		addEdge(nNodes-2, sink, maxCap)
	}

	// --- Scrittura File ---
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Header DIMACS
	fmt.Fprintf(f, "c Grafo generato da generator.go %d\n", mode)
	fmt.Fprintf(f, "p max %d %d\n", nNodes, len(edges))
	fmt.Fprintf(f, "n %d s\n", source+1)
	fmt.Fprintf(f, "n %d t\n", sink+1)

	for _, e := range edges {
		fmt.Fprintf(f, "a %d %d %d\n", e.u+1, e.v+1, e.cap)
	}

	return nil
}
