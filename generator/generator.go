// Package generator provides an interface and implementations of generators to create a group of points.
package generator

import "github.com/RH12503/Triangula/normgeom"

// A Generator is used to generate a group of points.
type Generator interface {
	// Generate generates and returns a point group with a specified number of points.
	Generate(numPoints int) normgeom.NormPointGroup
}
