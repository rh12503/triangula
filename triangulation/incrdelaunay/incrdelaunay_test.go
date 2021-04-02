package incrdelaunay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelaunay_NumPoints(t *testing.T) {
	delaunay := NewDelaunay(100,100)
	delaunay.Insert(Point{21, 34})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{32, 21})
	assert.Equal(t, delaunay.NumPoints(), 4)
}

func TestDelaunay_Insert(t *testing.T) {
	delaunay := NewDelaunay(100,100)
	delaunay.Insert(Point{21, 34})
	delaunay.Insert(Point{12, 32})
	delaunay.Insert(Point{32, 21})

	delaunay.IterTriangles(func(tri Triangle) {
		assert.Equal(t, tri.HasVertex(Point{21, 34}), true)
		assert.Equal(t, tri.HasVertex(Point{12, 32}), true)
		assert.Equal(t, tri.HasVertex(Point{32, 21}), true)
	})
}
