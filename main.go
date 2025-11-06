package main

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"fmt"
	"os"
)

func main() {
	// Genera un FlowNetwork casuale
	fn := flownetwork.NewFlowNetwork(
		10, // 10 nodi
		0,  // 30% di densità
		6,  // source = 0
	)

	fn.GenerateRandomArcs(0.3, 0, 10)

	// Esporta in JSON
	jsonStr := fn.ToJSON()

	// Salva su file
	err := os.WriteFile("export/flownetwork.json", []byte(jsonStr), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("FlowNetwork generato e salvato in flownetwork.json")
	fmt.Printf("Nodi: %d, Archi generati: circa %.0f\n", fn.N, 0.3*float64(fn.N*(fn.N-1)))
}
