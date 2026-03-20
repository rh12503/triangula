package geom

import (
	"github.com/RH12503/Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTriangle(t *testing.T) {
	triangle := NewTriangle(14, 15, 92, 65, 35, 89)
	assert.Equal(t, triangle, Triangle{Points: [3]Point{
		{14, 15},
		{92, 65},
		{35, 89},
	}})
}

func TestPoint_ToNorm(t *testing.T) {
	p := Point{X: 50, Y: 25}
	norm := p.ToNorm(100, 100)
	assert.Equal(t, norm, normgeom.NormPoint{X: 0.5, Y: 0.25})
}

func TestPoint_DistSq(t *testing.T) {
	p := Point{X: 0, Y: 0}
	assert.Equal(t, p.DistSq(Point{X: 3, Y: 4}), 25)
}

func TestPoint_Sub(t *testing.T) {
	p := Point{X: 5, Y: 3}
	assert.Equal(t, p.Sub(Point{X: 2, Y: 1}), Point{X: 3, Y: 2})
}

func TestPolygon_ToNorm(t *testing.T) {
	poly := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 50, Y: 100}}}
	norm := poly.ToNorm(100, 100)
	assert.Equal(t, norm, normgeom.NormPolygon{Points: []normgeom.NormPoint{
		{X: 0.0, Y: 0.0},
		{X: 1.0, Y: 0.0},
		{X: 0.5, Y: 1.0},
	}})
}

func TestPolygon_Triangulate_TriangleCount(t *testing.T) {
	// An N-gon fan-triangulates into N-2 triangles
	quad := Polygon{Points: []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}}
	count := 0
	quad.Triangulate(func(_ Triangle) { count++ })
	assert.Equal(t, count, 2)

	pentagon := Polygon{Points: []Point{{0, 0}, {5, 0}, {8, 4}, {4, 8}, {0, 5}}}
	count = 0
	pentagon.Triangulate(func(_ Triangle) { count++ })
	assert.Equal(t, count, 3)
}

func TestPolygon_Triangulate_SharesFirstPoint(t *testing.T) {
	// Fan triangulation: every triangle shares Points[0]
	poly := Polygon{Points: []Point{{1, 1}, {5, 0}, {9, 1}, {9, 5}}}
	poly.Triangulate(func(tri Triangle) {
		assert.Equal(t, tri.Points[2], poly.Points[0])
	})
}
