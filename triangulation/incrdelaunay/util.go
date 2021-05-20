package incrdelaunay

import "C"
import (
	"math"
)

// inCircle returns if point d is in the circumcircle of triangle abc.
func inCircle(aX, aY, bX, bY, cX, cY, dX, dY int64) int64 {
	if orientation(aX, aY, bX, bY, cX, cY) < 0 {
		aX, bX = bX, aX
		aY, bY = bY, aY
	}

	a11 := aX - dX
	a21 := bX - dX
	a31 := cX - dX

	a12 := aY - dY
	a22 := bY - dY
	a32 := cY - dY

	return (a11*a11+a12*a12)*(a21*a32-a31*a22) +
		(a21*a21+a22*a22)*(a31*a12-a11*a32) +
		(a31*a31+a32*a32)*(a11*a22-a21*a12)
}

// orientation returns a positive integer if points abc are clockwise.
func orientation(aX, aY, bX, bY, cX, cY int64) int64 {
	return (aX-cX)*(bY-cY) - (bX-cX)*(aY-cY)
}

// calculateScore calculates the score of a Devillers' ear.
// See https://hal.inria.fr/inria-00167201/document.
func calculateScore(a, b, c, d Point) float64 {
	orientation := orientation(int64(a.X), int64(a.Y), int64(b.X), int64(b.Y), int64(c.X), int64(c.Y))

	if orientation <= 0 {
		return math.MaxFloat64
	}

	inCircle := inCircle(int64(a.X), int64(a.Y), int64(b.X), int64(b.Y), int64(c.X), int64(c.Y), int64(d.X), int64(d.Y))

	return float64(inCircle) / float64(orientation)
}
