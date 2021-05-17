package geom

import "github.com/RH12503/Triangula/normgeom"

// Triangle represents a triangle with integer coordinates.
type Triangle struct {
	Points [3]Point
}

// ToNorm returns the normalized equivalent of a Triangle given a width and height.
func (t Triangle) ToNorm(w, h int) normgeom.NormTriangle {
	v := t.Points

	return normgeom.NormTriangle{Points: [3]normgeom.NormPoint{
		v[0].ToNorm(w, h),
		v[1].ToNorm(w, h),
		v[2].ToNorm(w, h),
	}}
}

// NewTriangle returns a new Triangle with specified vertex coordinates.
func NewTriangle(x0, y0, x1, y1, x2, y2 int) Triangle {
	return Triangle{[3]Point{{x0, y0}, {x1, y1}, {x2, y2}}}
}
