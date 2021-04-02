// Package algorithm provides an interface for an optimization algorithm as well as
// an implementation of the algorithm.
package algorithm

import (
	"github.com/RH12503/Triangula/normgeom"
	"time"
)

// An Algorithm is an iterative algorithm for optimizing a group of generator
type Algorithm interface {
	// Step runs one generation of the algorithm
	Step()

	// Best returns the point group with the highest fitness
	Best() normgeom.NormPointGroup

	// Stats returns simple statistics related to the algorithm
	Stats() Stats
}

// Stats contains basic statistics of an Algorithm
type Stats struct {
	BestFitness float64
	Generation  int
	TimeForGen  time.Duration // The time taken for the last generation
}
