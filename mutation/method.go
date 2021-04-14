// Package mutation provides an interface for a mutation method and some implementations.
package mutation

import (
	"github.com/RH12503/Triangula/normgeom"
)

// A Method is used to apply mutations on a point group.
type Method interface {
	// Mutate mutates a normgeom.NormPointGroup.
	// A function, mutated, is called when a point is mutated.
	Mutate(points normgeom.NormPointGroup, mutated func(mutation Mutation))
}
