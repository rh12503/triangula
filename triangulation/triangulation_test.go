package triangulation

import (
	"Triangula/geom"
	"Triangula/normgeom"
	"Triangula/triangulation/delaunay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArea(t *testing.T) {
	area := Area([]delaunay.Point{{0,10}, {2, 4}, {4, 5},{5, 0}, {0,0}})
	assert.Equal(t, area, 26)
}

func TestTriangulatePoints(t *testing.T) {
	tri, count := TriangulatePoints([]delaunay.Point{{0,10}, {2, 4}, {4, 5}})

	assert.Equal(t, tri, []geom.Triangle{geom.NewTriangle(4, 5, 2, 4,0 ,10)})
	assert.Equal(t, count, 7)
}

func TestTriangulate(t *testing.T) {
	tri, count := Triangulate(normgeom.NormPointGroup{{0,1}, {0.2, 0.4}, {0.4, 0.5}}, 10, 10)

	assert.Equal(t, tri, []geom.Triangle{geom.NewTriangle(4, 5, 2, 4,0 ,10)})
	assert.Equal(t, count, 7)
}

