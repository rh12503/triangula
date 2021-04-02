package generator

import (
	"github.com/RH12503/Triangula/normgeom"
	"math/rand"
)

// A RandomGenerator generates a normgeom.NormPointGroup of random normalized generator.
type RandomGenerator struct {
}

func (r RandomGenerator) Generate(n int) normgeom.NormPointGroup {
	return randomPoints(n)
}

// randomPoints returns a normgeom.NormPointGroup with n number of random generator.
func randomPoints(n int) normgeom.NormPointGroup {
	points := normgeom.NormPointGroup{}

	for i := 0; i < n; i++ {
		points = append(points, normgeom.NormPoint{X: rand.Float64(), Y: rand.Float64()})
	}

	return points
}
