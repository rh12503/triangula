package polygonation

import (
	"github.com/RH12503/Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolygonate_ReturnsNonEmptyPolygons(t *testing.T) {
	points := normgeom.NormPointGroup{
		{X: 0.2, Y: 0.2},
		{X: 0.8, Y: 0.2},
		{X: 0.5, Y: 0.8},
		{X: 0.5, Y: 0.5},
	}
	polygons := Polygonate(points, 100, 100)

	assert.True(t, len(polygons) > 0)
	for _, p := range polygons {
		assert.True(t, len(p.Points) >= 3, "expected polygon to have at least 3 points")
	}
}

func TestPolygonate_PointsAreWithinBounds(t *testing.T) {
	points := normgeom.NormPointGroup{
		{X: 0.1, Y: 0.1},
		{X: 0.9, Y: 0.1},
		{X: 0.5, Y: 0.9},
		{X: 0.3, Y: 0.5},
		{X: 0.7, Y: 0.5},
	}
	w, h := 100, 100
	polygons := Polygonate(points, w, h)

	for _, poly := range polygons {
		for _, p := range poly.Points {
			assert.True(t, p.X >= 0 && p.X <= w, "x=%d out of [0,%d]", p.X, w)
			assert.True(t, p.Y >= 0 && p.Y <= h, "y=%d out of [0,%d]", p.Y, h)
		}
	}
}

func TestPolygonate_PolygonCountMatchesInputPoints(t *testing.T) {
	points := normgeom.NormPointGroup{
		{X: 0.2, Y: 0.2},
		{X: 0.8, Y: 0.2},
		{X: 0.5, Y: 0.8},
		{X: 0.3, Y: 0.5},
		{X: 0.7, Y: 0.5},
	}
	polygons := Polygonate(points, 100, 100)

	// Each input point should produce exactly one Voronoi cell
	assert.Equal(t, len(polygons), len(points))
}
