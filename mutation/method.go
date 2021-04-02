// Package mutation provides an interface for a mutation method and some common implementations.

package mutation

import (
	"github.com/RH12503/Triangula/normgeom"
)

// A Method is used to apply a mutation operator on a point group.
type Method interface {
	// Mutate mutates a normgeom.NormPointGroup.
	// A function mutated is called when a point is mutated
	Mutate(points normgeom.NormPointGroup, mutated func(mutation Mutation))
}
