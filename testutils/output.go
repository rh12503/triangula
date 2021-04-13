package testutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/RH12503/Triangula/algorithm"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// GenerateAlgorithmOutput runs an algorithm.Algorithm and writes the output points to a file
func GenerateAlgorithmOutput(outputFile string, algo algorithm.Algorithm, reps int) {
	dataFile, _ := os.Create(outputFile + "-stats")
	writer := bufio.NewWriter(dataFile)

	for {
		ti := time.Now()
		for i := 0; i < reps; i++ {
			algo.Step()
		}
		stats := algo.Stats()
		fmt.Printf("Gen: %v | Fit: %v | Time: %v\n", stats.Generation, stats.BestFitness, float64(time.Since(ti).Microseconds())/(float64(reps)*1000.))
		printMemUsage()

		jsonOut, err := json.Marshal(algo.Best())
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(outputFile, jsonOut, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writer.WriteString(fmt.Sprintf("%v, %v\n", stats.Generation, stats.BestFitness))
		writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}

	dataFile.Close()
}
