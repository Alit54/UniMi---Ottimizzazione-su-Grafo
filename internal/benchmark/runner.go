package benchmark

import (
	"OttimizzazioneSuGrafo/internal/flownetwork"
	"OttimizzazioneSuGrafo/internal/maxflow"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func Run() {
	os.MkdirAll("export", os.ModePerm)
	csvFile, _ := os.Create("export/benchmark_results.csv")
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	header := []string{
		"Graph",
		"Nodes",
		"Edges",
		"Algorithm",
		"MaxFlow",
		"Time",
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

	networks := LoadSuite("data/flownetwork/BVZ-tsukuba")
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
				strconv.FormatInt(int64(result.Time), 10),
				fmt.Sprintf("%d", result.Augments),
				fmt.Sprintf("%d", result.Retreats),
				fmt.Sprintf("%d", result.Phases),
				fmt.Sprintf("%d", result.Advances),
			}
			csvWriter.Write(row)
			csvWriter.Flush()
			fmt.Printf("-> %s: %d ms | Flow: %d\n", result.Algorithm, result.Time, result.MaxFlow)
		}
	}
}
