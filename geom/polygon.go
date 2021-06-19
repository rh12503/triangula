package geom

import "github.com/RH12503/Triangula/normgeom"

// Polygon represents a n-gon with integer coordinates.
type Polygon struct {
	Points []Point
}

func (p Polygon) ToNorm(w, h int) normgeom.NormPolygon {
	v := p.Points

	new := normgeom.NormPolygon{Points: make([]normgeom.NormPoint, len(v))}

	for i := range p.Points {
		new.Points[i] = p.Points[i].ToNorm(w, h)
	}

	return new
}

func (p Polygon) Triangulate(triangle func(Triangle)) {
	for i := 2; i < len(p.Points); i++ {
		tri := Triangle{[3]Point{
			p.Points[i],
			p.Points[i-1],
			p.Points[0],
		}}

		triangle(tri)
	}
}
