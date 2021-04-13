package triangulation

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/fogleman/delaunay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArea(t *testing.T) {
	area := Area([]delaunay.Point{{0, 10}, {2, 4}, {4, 5}, {5, 0}, {0, 0}})
	assert.Equal(t, area, 26)
}

func TestTriangulate(t *testing.T) {
	tri := Triangulate(normgeom.NormPointGroup{{0, 1}, {0.2, 0.4}, {0.4, 0.5}}, 10, 10)

	assert.Equal(t, tri, []geom.Triangle{geom.NewTriangle(4, 5, 2, 4, 0, 10)})
}
