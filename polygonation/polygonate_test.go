package polygonation

import (
	"github.com/RH12503/Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolygonate_Basic(t *testing.T) {
	points := normgeom.NormPointGroup{
		{0.25, 0.25},
		{0.75, 0.25},
		{0.5, 0.75},
	}

	polygons := Polygonate(points, 100, 100)

	// Should produce polygons (one Voronoi cell per point)
	assert.Equal(t, len(points), len(polygons))

	// Each polygon should have at least 3 points
	for _, poly := range polygons {
		assert.GreaterOrEqual(t, len(poly.Points), 3)
	}
}

func TestPolygonate_SinglePoint(t *testing.T) {
	points := normgeom.NormPointGroup{
		{0.5, 0.5},
	}

	polygons := Polygonate(points, 100, 100)

	// A single point does not produce a valid Voronoi polygon from the
	// Delaunay triangulation (needs at least 3 points for non-degenerate triangles)
	assert.Equal(t, 0, len(polygons))
}

func TestPolygonate_BoundaryPoints(t *testing.T) {
	points := normgeom.NormPointGroup{
		{0.0, 0.0},
		{1.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
		{0.5, 0.5},
	}

	polygons := Polygonate(points, 50, 50)

	assert.Equal(t, len(points), len(polygons))
	for _, poly := range polygons {
		// Boundary Voronoi cells can be clipped, resulting in polygons with 2+ points
		assert.GreaterOrEqual(t, len(poly.Points), 2)
	}
}
