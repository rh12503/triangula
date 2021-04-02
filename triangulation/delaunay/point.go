package delaunay

import "math"

type Point struct {
	X, Y float64
}

func (a Point) squaredDistance(b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (a Point) distance(b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

func (a Point) sub(b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}
