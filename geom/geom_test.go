package geom

import (
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

func TestPolygon_Triangulate_Triangle(t *testing.T) {
	// A 3-point polygon should yield exactly 1 triangle
	poly := Polygon{Points: []Point{{0, 0}, {10, 0}, {5, 10}}}
	var triangles []Triangle
	poly.Triangulate(func(tri Triangle) {
		triangles = append(triangles, tri)
	})
	assert.Equal(t, 1, len(triangles))
	assert.Equal(t, Triangle{Points: [3]Point{{5, 10}, {10, 0}, {0, 0}}}, triangles[0])
}

func TestPolygon_Triangulate_Quad(t *testing.T) {
	// A 4-point polygon should yield 2 triangles (fan decomposition)
	poly := Polygon{Points: []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}}
	var triangles []Triangle
	poly.Triangulate(func(tri Triangle) {
		triangles = append(triangles, tri)
	})
	assert.Equal(t, 2, len(triangles))
	// All triangles should share the first vertex
	for _, tri := range triangles {
		assert.Equal(t, Point{0, 0}, tri.Points[2])
	}
}

func TestPolygon_Triangulate_Pentagon(t *testing.T) {
	poly := Polygon{Points: []Point{{0, 0}, {5, 0}, {8, 4}, {4, 8}, {-1, 4}}}
	count := 0
	poly.Triangulate(func(tri Triangle) {
		count++
	})
	assert.Equal(t, 3, count) // n-2 triangles for n-gon
}

func TestPolygon_ToNorm(t *testing.T) {
	poly := Polygon{Points: []Point{{50, 25}, {100, 50}}}
	norm := poly.ToNorm(100, 50)
	assert.Equal(t, 2, len(norm.Points))
	assert.Equal(t, 0.5, norm.Points[0].X)
	assert.Equal(t, 0.5, norm.Points[0].Y)
	assert.Equal(t, 1.0, norm.Points[1].X)
	assert.Equal(t, 1.0, norm.Points[1].Y)
}

func TestPoint_ToNorm(t *testing.T) {
	p := Point{25, 50}
	norm := p.ToNorm(100, 200)
	assert.Equal(t, 0.25, norm.X)
	assert.Equal(t, 0.25, norm.Y)
}

func TestPoint_DistSq(t *testing.T) {
	a := Point{0, 0}
	b := Point{3, 4}
	assert.Equal(t, 25, a.DistSq(b))
	assert.Equal(t, 0, a.DistSq(a))
}

func TestPoint_Sub(t *testing.T) {
	a := Point{10, 20}
	b := Point{3, 5}
	assert.Equal(t, Point{7, 15}, a.Sub(b))
}
