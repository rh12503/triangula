package render

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrianglesOnImage(t *testing.T) {
	data := TrianglesOnImage([]geom.Triangle{
		geom.NewTriangle(12, 32, 65, 43, 23, 87),
	}, image.NewData(100, 100))

	assert.Equal(t, data[0], TriangleData{
		Triangle: normgeom.NewNormTriangle(0.12, 0.32, 0.65, 0.43, 0.23, 0.87),
	})
}

func TestTrianglesOnImage_Multiple(t *testing.T) {
	data := TrianglesOnImage([]geom.Triangle{
		geom.NewTriangle(10, 10, 50, 10, 30, 50),
		geom.NewTriangle(60, 60, 90, 60, 75, 90),
	}, image.NewData(100, 100))

	assert.Equal(t, 2, len(data))
}

func TestPolygonsOnImage(t *testing.T) {
	polygons := []geom.Polygon{
		{Points: []geom.Point{{10, 10}, {50, 10}, {50, 50}, {10, 50}}},
	}
	data := PolygonsOnImage(polygons, image.NewData(100, 100))

	assert.Equal(t, 1, len(data))
	assert.Equal(t, 4, len(data[0].Polygon.Points))
}

func TestPolygonsOnImage_Multiple(t *testing.T) {
	polygons := []geom.Polygon{
		{Points: []geom.Point{{10, 10}, {30, 10}, {30, 30}, {10, 30}}},
		{Points: []geom.Point{{50, 50}, {70, 50}, {70, 70}, {50, 70}}},
	}
	data := PolygonsOnImage(polygons, image.NewData(100, 100))

	assert.Equal(t, 2, len(data))
}
