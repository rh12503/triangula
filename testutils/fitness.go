package testutils

import (
	"fmt"
	"github.com/RH12503/Triangula/algorithm"
)

// RunAlgorithm runs an algorithm.Algorithm and prints the fitness after a specified number of repetitions.
func RunAlgorithm(algo algorithm.Algorithm, reps int) {
	for {
		//ti := time.Now()
		for i := 0; i < reps; i++ {
			algo.Step()
		}
		stats := algo.Stats()
		fmt.Printf("%v, %v\n", stats.Generation, stats.BestFitness)
		//fmt.Printf("Gen: %v | Fit: %v | Time: %v\n", stats.Generation, stats.BestFitness, float64(time.Since(ti).Microseconds())/(float64(reps)*1000.))
	}
}
