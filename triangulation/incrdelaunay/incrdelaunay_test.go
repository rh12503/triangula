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

func TestDelaunay_Remove_DecreasesPointCount(t *testing.T) {
	d := NewDelaunay(100, 100)
	d.Insert(Point{20, 20})
	d.Insert(Point{80, 20})
	d.Insert(Point{50, 80})
	d.Insert(Point{50, 50})
	assert.Equal(t, d.NumPoints(), 4)

	d.Remove(Point{50, 50})
	assert.Equal(t, d.NumPoints(), 3)
}

func TestDelaunay_Remove_PointAbsentFromTriangles(t *testing.T) {
	d := NewDelaunay(100, 100)
	d.Insert(Point{20, 20})
	d.Insert(Point{80, 20})
	d.Insert(Point{50, 80})
	d.Insert(Point{50, 50})

	removed := Point{50, 50}
	d.Remove(removed)

	d.IterTriangles(func(tri Triangle) {
		assert.False(t, tri.HasVertex(removed), "removed point should not appear in any triangle")
	})
}

func TestDelaunay_Set_CopiesPointCount(t *testing.T) {
	original := NewDelaunay(100, 100)
	original.Insert(Point{20, 20})
	original.Insert(Point{80, 20})
	original.Insert(Point{50, 80})

	copied := NewDelaunay(100, 100)
	copied.Set(original)

	assert.Equal(t, copied.NumPoints(), original.NumPoints())
}

func TestDelaunay_Set_IsIndependent(t *testing.T) {
	original := NewDelaunay(100, 100)
	original.Insert(Point{20, 20})
	original.Insert(Point{80, 20})
	original.Insert(Point{50, 80})
	original.Insert(Point{50, 50})

	copied := NewDelaunay(100, 100)
	copied.Set(original)
	copied.Remove(Point{50, 50})

	// Original should be unchanged
	assert.Equal(t, original.NumPoints(), 4)
}

func TestVoronoi_ProducesPolygons(t *testing.T) {
	d := NewDelaunay(100, 100)
	d.Insert(Point{20, 20})
	d.Insert(Point{80, 20})
	d.Insert(Point{50, 80})
	d.Insert(Point{50, 50})

	count := 0
	Voronoi(d, func(points []FloatPoint) {
		count++
		assert.True(t, len(points) >= 3, "Voronoi cell should have at least 3 vertices")
	}, 100, 100)

	assert.True(t, count > 0, "expected Voronoi to produce at least one polygon")
}

func TestVoronoi_PolygonCountMatchesPointCount(t *testing.T) {
	d := NewDelaunay(100, 100)
	pts := []Point{{20, 20}, {80, 20}, {50, 80}, {20, 70}, {80, 70}}
	for _, p := range pts {
		d.Insert(p)
	}

	count := 0
	Voronoi(d, func(_ []FloatPoint) { count++ }, 100, 100)

	assert.Equal(t, count, len(pts))
}
