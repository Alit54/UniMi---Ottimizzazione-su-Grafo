package benchmark

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"OttimizzazioneSuGrafo/internal/maxflow"
	"encoding/csv"
	"fmt"
	"os"
)

func Run(input string, output string) {
	os.MkdirAll("export", os.ModePerm)
	csvFile, _ := os.OpenFile(output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	header := []string{
		"Graph",
		"Nodes",
		"Edges",
		"Algorithm",
		"MaxFlow",
		"Time (nanoseconds)",
		"Augments",
		"Retreats",
		"Advances",
		"Phases",
	}
	csvWriter.Write(header)

	toRun := map[string]maxflow.MaxFlowAlgorithm{}
	toRun["Capacity Scaling"] = &maxflow.CapacityScaling{}
	toRun["Shortest Augmenting Path"] = &maxflow.ShortestAugmentingPath{}
	toRun["Dinic"] = &maxflow.Dinic{}
	toRun["Capacity Scaling con SAP"] = &maxflow.ScalingSAP{}
	toRun["Recursive Dinic"] = &maxflow.RecursiveDinic{}

	networks := LoadSuite(input)
	for _, network := range networks {
		fn := flownetwork.NewNetworkFromDIMACS(network)
		for i, alg := range toRun {
			fn.Reset()
			result := Benchmark(fn, i, alg, network, 1)
			row := []string{
				result.Graph,
				fmt.Sprintf("%d", result.Nodes),
				fmt.Sprintf("%d", result.Edges),
				result.Algorithm,
				fmt.Sprintf("%d", result.MaxFlow),
				fmt.Sprintf("%d", result.Time.Nanoseconds()),
				fmt.Sprintf("%d", result.Augments),
				fmt.Sprintf("%d", result.Retreats),
				fmt.Sprintf("%d", result.Phases),
				fmt.Sprintf("%d", result.Advances),
			}
			csvWriter.Write(row)
			csvWriter.Flush()
			fmt.Printf("-> %s: %.2f s | Flow: %d\n", result.Algorithm, float64(result.Time)/1e9, result.MaxFlow)
		}
	}
}
