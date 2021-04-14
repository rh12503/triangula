package generator

import (
	"github.com/RH12503/Triangula/normgeom"
	"math/rand"
)

// A RandomGenerator generates a point group filled with random points.
type RandomGenerator struct {
}

// Generate returns a set of randomly distributed points.
func (r RandomGenerator) Generate(n int) normgeom.NormPointGroup {
	return randomPoints(n)
}

// randomPoints returns a point group with a specified number of points.
func randomPoints(n int) normgeom.NormPointGroup {
	points := normgeom.NormPointGroup{}

	for i := 0; i < n; i++ {
		points = append(points, normgeom.NormPoint{X: rand.Float64(), Y: rand.Float64()})
	}

	return points
}
