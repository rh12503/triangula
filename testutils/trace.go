package testutils

import (
	"bufio"
	"github.com/RH12503/Triangula/algorithm"
	"os"
	"runtime/trace"
)

// GenerateTrace generates a trace of an algorithm.
// Used for debugging and evaluating performance.
func GenerateTrace(outputFile string, algo algorithm.Algorithm) {
	dataFile, _ := os.Create(outputFile + ".trace")
	writer := bufio.NewWriter(dataFile)
	for i := 0; i < 2000; i++ {
		algo.Step()
	}

	trace.Start(writer)
	for i := 0; i < 1000; i++ {
		algo.Step()
	}
	trace.Stop()
	dataFile.Close()
}
