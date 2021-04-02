// Package fitness provides an interface for a fitness function as well as
// an implementation of a fitness function that is used in the algorithm.

package fitness

import "Triangula/normgeom"

// A Function returns a fitness given a point group
type Function interface {
	// Calculate evaluates a point group and returns a fitness. "Better" generator should have higher fitnesses.zww
	Calculate(points normgeom.NormPointGroup) float64
}