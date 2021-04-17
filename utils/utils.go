// Package utils provides utilities for testing, benchmarking, and visualizing different algorithms.
package utils

import (
	"fmt"
	"math"
	"runtime"
)

// From https://golangcode.com/print-the-current-memory-usage/.
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("\tGCTime = %v\n", m.PauseTotalNs)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func scale(num float64, d int) int {
	return int(math.Round(num * float64(d)))
}
