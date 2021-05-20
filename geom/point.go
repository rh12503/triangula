// Package geom contains structs for representing a Point and a Triangle.
package geom

import (
	"github.com/RH12503/Triangula/normgeom"
)

// Point represents a 2D point with integer coordinates.
type Point struct {
	X int
	Y int
}

// ToNorm returns the normalized equivalent of a Point given a width and height.
func (p Point) ToNorm(w, h int) normgeom.NormPoint {
	return normgeom.NormPoint{X: float64(p.X) / float64(w), Y: float64(p.Y) / float64(h)}
}

func (p Point) DistSq(other Point) int {
	dX := p.X - other.X
	dY := p.Y - other.Y
	return dX*dX + dY*dY
}

func (p Point) Sub(other Point) Point {
	return Point{p.X - other.X, p.Y - other.Y}
}
