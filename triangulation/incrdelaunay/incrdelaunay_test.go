package incrdelaunay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelaunay_NumPoints(t *testing.T) {
	delaunay := NewDelaunay(100, 100)
	delaunay.Insert(Point{21, 34})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{32, 21})
	assert.Equal(t, delaunay.NumPoints(), 4)
}

func TestDelaunay_Insert(t *testing.T) {
	delaunay := NewDelaunay(100, 100)
	delaunay.Insert(Point{21, 34})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{32, 21})

	delaunay.IterTriangles(func(tri Triangle) {
		assert.Equal(t, tri.HasVertex(Point{21, 34}), true)
		assert.Equal(t, tri.HasVertex(Point{12, 32}), true)
		assert.Equal(t, tri.HasVertex(Point{32, 21}), true)
	})
}

func TestSegmentsIntersect(t *testing.T) {
	// Crossing segments
	x, y, ok := segmentsIntersect(0, 0, 10, 10, 0, 10, 10, 0)
	assert.True(t, ok)
	assert.InDelta(t, 5.0, x, 0.01)
	assert.InDelta(t, 5.0, y, 0.01)

	// Parallel segments (no intersection)
	_, _, ok = segmentsIntersect(0, 0, 10, 0, 0, 5, 10, 5)
	assert.False(t, ok)

	// Non-intersecting segments (too short)
	_, _, ok = segmentsIntersect(0, 0, 1, 1, 5, 0, 5, 5)
	assert.False(t, ok)
}

func TestOutside(t *testing.T) {
	assert.False(t, outside(FloatPoint{5, 5}, 10, 10))
	assert.False(t, outside(FloatPoint{0, 0}, 10, 10))
	assert.False(t, outside(FloatPoint{10, 10}, 10, 10))
	assert.True(t, outside(FloatPoint{-1, 5}, 10, 10))
	assert.True(t, outside(FloatPoint{5, -1}, 10, 10))
	assert.True(t, outside(FloatPoint{11, 5}, 10, 10))
	assert.True(t, outside(FloatPoint{5, 11}, 10, 10))
}

func TestEpsilonEquals(t *testing.T) {
	assert.True(t, epsilonEquals(1.0, 1.0))
	assert.True(t, epsilonEquals(1.0, 1.00005))
	assert.False(t, epsilonEquals(1.0, 1.001))
	assert.False(t, epsilonEquals(0.0, 1.0))
}

func TestPolygonPointIntersect(t *testing.T) {
	square := []FloatPoint{{0, 0}, {10, 0}, {10, 10}, {0, 10}}

	// Point inside
	assert.True(t, polygonPointIntersect(square, FloatPoint{5, 5}))
	// Point on vertex
	assert.True(t, polygonPointIntersect(square, FloatPoint{0, 0}))
	// Point outside
	assert.False(t, polygonPointIntersect(square, FloatPoint{15, 15}))
}

func TestTriangle_HasVertex(t *testing.T) {
	tri := NewTriangle(Point{0, 0}, Point{10, 0}, Point{5, 10})
	assert.True(t, tri.HasVertex(Point{0, 0}))
	assert.True(t, tri.HasVertex(Point{10, 0}))
	assert.True(t, tri.HasVertex(Point{5, 10}))
	assert.False(t, tri.HasVertex(Point{1, 1}))
}

func TestPoint_DistSq(t *testing.T) {
	a := Point{0, 0}
	b := Point{3, 4}
	assert.Equal(t, int64(25), a.DistSq(b))
	assert.Equal(t, int64(0), a.DistSq(a))
}

func TestPoint_Hash(t *testing.T) {
	a := Point{3, 4}
	b := Point{3, 4}
	c := Point{5, 6}
	assert.Equal(t, a.Hash(), b.Hash())
	assert.NotEqual(t, a.Hash(), c.Hash())
}

func TestEdge_Equals(t *testing.T) {
	e1 := NewEdge(Point{0, 0}, Point{10, 10})
	e2 := NewEdge(Point{0, 0}, Point{10, 10})
	e3 := NewEdge(Point{5, 5}, Point{10, 10})

	assert.True(t, e1.Equals(e2))
	assert.False(t, e1.Equals(e3))
}

func TestNewEdge_Ordering(t *testing.T) {
	// NewEdge should sort points so A <= B by X, then Y
	e1 := NewEdge(Point{10, 5}, Point{3, 5})
	assert.Equal(t, Point{3, 5}, e1.A)
	assert.Equal(t, Point{10, 5}, e1.B)

	// Same X, different Y
	e2 := NewEdge(Point{5, 10}, Point{5, 3})
	assert.Equal(t, Point{5, 3}, e2.A)
	assert.Equal(t, Point{5, 10}, e2.B)
}

func TestVoronoi_Basic(t *testing.T) {
	delaunay := NewDelaunay(100, 100)
	delaunay.Insert(Point{25, 25})
	delaunay.Insert(Point{75, 25})
	delaunay.Insert(Point{50, 75})

	count := 0
	Voronoi(delaunay, func(points []FloatPoint) {
		count++
		assert.GreaterOrEqual(t, len(points), 3)
	}, 100, 100)

	assert.Equal(t, 3, count)
}
