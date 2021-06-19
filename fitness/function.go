// Package fitness provides an interface for a fitness function as well as
// an implementation of a fitness function that is used in the algorithm.
package fitness

import (
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
)

// A Function represents a fitness function to evaluate a point group.
type Function interface {
	// Calculate evaluates a point group and returns a fitness.
	Calculate(data PointsData) float64
}

// PointsData stores data regarding a point group, and is used by the fitness function.
type PointsData struct {
	Points    normgeom.NormPointGroup
	Mutations []mutation.Mutation
}
