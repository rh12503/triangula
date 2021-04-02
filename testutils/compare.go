package testutils

import (
	"Triangula/algorithm"
	"Triangula/random"
	"fmt"
	"math/rand"
	"time"
)

// CompareAlgorithms compares the effectiveness of two algorithms
func CompareAlgorithms(a, b func() algorithm.Algorithm) {
	var averagesA []int
	var averagesB []int
	count := 0

	for false {
		rand.Seed(time.Now().UnixNano())
		random.Seed(time.Now().UnixNano())

		algoA := a()
		algoB := b()

		for algoA.Stats().Generation < 10000 {
			algoA.Step()
		}

		for algoB.Stats().BestFitness < algoA.Stats().BestFitness {
			algoB.Step()
		}

		averagesA = append(averagesA, algoA.Stats().Generation)
		averagesB = append(averagesB, algoB.Stats().Generation)
		count++

		totalA := 0
		totalB := 0
		for i := 0; i < count; i++ {
			totalA += averagesA[i]
			totalB += averagesB[i]
		}

		fmt.Println(float64(totalA)/float64(count), float64(totalB)/float64(count), count)
	}
}