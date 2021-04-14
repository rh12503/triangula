package incrdelaunay

import (
	"math"
)

// Triangle stores the vertices of a triangle as well as its circumcircle.
type Triangle struct {
	A, B, C      Point
	Circumcircle Circumcircle
}

// NewTriangle return a new Triangle and calculates its circumcircle given three points.
func NewTriangle(a, b, c Point) Triangle {
	return Triangle{
		A:            a,
		B:            b,
		C:            c,
		Circumcircle: calcCircumcircle(a, b, c),
	}
}

// HasVertex returns if the Triangle contains a specified vertex.
func (t Triangle) HasVertex(p Point) bool {
	return t.A == p || t.B == p || t.C == p
}

// NewSuperTriangle returns a Triangle large enough to cover all points within (0, 0) to (w, h).
func NewSuperTriangle(w, h int) Triangle {
	hW := int16(math.Ceil(float64(w) / 2))
	hH := int16(math.Ceil(float64(h) / 2))

	max := int16(w)
	if h > w {
		max = int16(h)
	}
	a := Point{hW - 2*max, hH - max}
	b := Point{hW, hH + 2*max}
	c := Point{hW + 2*max, hH - max}

	return NewTriangle(a, b, c)
}

// Point represents a 2D point, using int16 to optimize space.
type Point struct {
	X, Y int16
}

// DistSq returns the distance squared to another point.
func (p Point) DistSq(b Point) int64 {
	dX := int64(b.X - p.X)
	dY := int64(b.Y - p.Y)

	return dX*dX + dY*dY
}

// Hash returns a hash code for the point.
func (p Point) Hash() int {
	return (53+int(p.X))*53 + int(p.Y)
}

// Circumcircle represents a circumcircle of a Triangle.
type Circumcircle struct {
	cX, cY float32
	Radius float32
}

// calcCircumcircle calculates the circumcircle of three points.
func calcCircumcircle(v0, v1, v2 Point) Circumcircle {
	var circumcircle Circumcircle

	A := int64(v1.X - v0.X)
	B := int64(v1.Y - v0.Y)
	C := int64(v2.X - v0.X)
	D := int64(v2.Y - v0.Y)

	E := A*int64(v0.X+v1.X) + B*int64(v0.Y+v1.Y)
	F := C*int64(v0.X+v2.X) + D*int64(v0.Y+v2.Y)

	G := float64(2 * (A*int64(v2.Y-v1.Y) - B*int64(v2.X-v1.X)))

	cx := float64(D*E-B*F) / G
	cy := float64(A*F-C*E) / G

	circumcircle.cX, circumcircle.cY = float32(cx), float32(cy)

	dx := cx - float64(v0.X)
	dy := cy - float64(v0.Y)

	circumcircle.Radius = float32(math.Sqrt(dx*dx + dy*dy))

	return circumcircle
}

// ear represents a Devillers ear.
type ear struct {
	a, b, c Point
	score   float64
}

// computeScore computes the score of the Devillers ear.
func (e *ear) computeScore(p Point) {
	e.score = calculateScore(e.a, e.b, e.c, p)
}

// Edge represents an edge from point A to B.
type Edge struct {
	A, B Point
}

// Equals returns if the edge is equal to another.
func (e Edge) Equals(b Edge) bool {
	// A and B are ordered, so it isn't necessary to check the other way around
	return e.A == b.A && e.B == b.B
}

// NewEdge returns a new edge with its points sorted.
func NewEdge(a, b Point) Edge {
	// Order the points in the edge
	if a.X > b.X {
		a, b = b, a
	} else if a.X == b.X {
		if a.Y > b.Y {
			a, b = b, a
		}
	}

	return Edge{a, b}
}
