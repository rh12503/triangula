package utils

import (
	"github.com/RH12503/Triangula/algorithm"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

// GenerateProfile creates a CPU profile by running the algorithm.
func GenerateProfile(outputFile string, algo algorithm.Algorithm, seconds int) {
	f, err := os.Create(outputFile + ".prof")
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	pprof.StartCPUProfile(f)
	for time.Since(t).Seconds() < float64(seconds) {
		for i := 0; i < 5; i++ {
			algo.Step()
		}
	}

	defer pprof.StopCPUProfile()
}
