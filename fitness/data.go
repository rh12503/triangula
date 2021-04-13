package fitness

import (
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
)

// PointsData stores data regarding a point group, and is used by the fitness function.
type PointsData struct {
	Points    normgeom.NormPointGroup
	Mutations []mutation.Mutation
}

// TriFit stores the triangles vertices and its fitness. The struct is used to cache calculations
type TriFit struct {
	aX, aY    int16
	bX, bY    int16
	cX, cY    int16
	fitness   float64
	OtherHash uint32
}

// Equals returns if the TriFit is equal to another
func (t TriFit) Equals(other TriFit) bool {
	return t.aX == other.aX && t.aY == other.aY &&
		t.bX == other.bX && t.bY == other.bY &&
		t.cX == other.cX && t.cY == other.cY
}

// Hash calculates the hash of a TriFit
func (t TriFit) Hash() uint64 {
	x := int(t.aX) + int(t.bX) + int(t.cX)
	y := int(t.aY) + int(t.bY) + int(t.cY)

	hash := uint64((53+x)*53 + y)

	return hash
}

// fastRound is an optimized version of math.Round
func fastRound(n float64) int {
	return int(n+0.5) << 0
}
