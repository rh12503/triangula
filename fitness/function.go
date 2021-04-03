// Package fitness provides an interface for a fitness function as well as
// an implementation of a fitness function that is used in the algorithm.

package fitness

// A Function returns a fitness given a point group
type Function interface {
	// Calculate evaluates a point group and returns a fitness. "Better" generator should have higher fitnesses.zww
	Calculate(data PointsData) float64
}
