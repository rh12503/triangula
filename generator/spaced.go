package generator

import (
	"Triangula/normgeom"
	"math"
	"math/rand"
)

const startTemp = 0.5
const minTemp = 0.02

type spacedGenerator struct {
	iterations int
	decrement  float64
}

func (s spacedGenerator) Generate(n int) normgeom.NormPointGroup {
	points := randomPoints(n)
	temp := startTemp

	for i := 0; i < s.iterations; i++ {
		ran := rand.Intn(n)
		p := points[ran]
		point := &points[ran]

		_, currDist := closestTo(p, points)
		p.X += (rand.Float64() - 0.5) * temp
		p.Y += (rand.Float64() - 0.5) * temp

		p.Constrain()

		_, newDist := closestTo(p, points)

		if newDist > currDist {
			point.X = p.X
			point.Y = p.Y
		}

		temp *= s.decrement
	}
	return points
}

func closestTo(point normgeom.NormPoint, group normgeom.NormPointGroup) (normgeom.NormPoint, float64) {
	closest := group[0]
	dist := -1.

	for _, p := range group {
		if p == point {
			continue
		}

		newDist := normgeom.Dist(p, point)

		if dist == -1 || newDist < dist {
			dist = newDist
			closest = p
		}
	}
	return closest, dist
}

func NewSpacedGenerator(iterations int) spacedGenerator {
	gen := spacedGenerator{iterations: iterations}

	gen.decrement = math.Pow(minTemp/startTemp, 1./float64(iterations))

	return gen
}
