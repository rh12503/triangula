// Package normgeom contains similar structs to the geom package, except all coordinates
// are normalized between 0 and 1.

package normgeom

import (
	"math"
)

// NormPoint represents a 2D point with coordinates between 0 and 1
type NormPoint struct {
	X float64
	Y float64
}

// Constrain constrains a NormPoint's X and Y values to between 0 and 1
func (p *NormPoint) Constrain() {
	if p.X > 1 {
		p.X = 1
	} else if p.X < 0 {
		p.X = 0
	}

	if p.Y > 1 {
		p.Y = 1
	} else if p.Y < 0 {
		p.Y = 0
	}
}

// Dist calculates the distance between 2 NormPoint's
func Dist(a, b NormPoint) float64 {
	dX := a.X - b.X
	dY := a.Y - b.Y

	return math.Sqrt(dX*dX + dY*dY)
}

// NormPointGroup represents a group of NormPoint's
type NormPointGroup []NormPoint

// Set sets a NormPointGroup to another NormPointGroup.
// The lengths of the NormPointGroup's cannot be different
func (p NormPointGroup) Set(other NormPointGroup) {
	if len(p) == len(other) {
		copy(p, other)
	} else {
		panic("lengths of point groups need to be the same")
	}
}

// Copy creates a deep copy of a NormPointGroup
func (p NormPointGroup) Copy() NormPointGroup {
	newGroup := make(NormPointGroup, len(p))
	copy(newGroup, p)

	return newGroup
}
