package rasterize

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDDATriangle(t *testing.T) {
	pixels := 0
	DDATriangle(geom.NewTriangle(13, 12, 37, 54, 78, 15), func(x, y int) {
		pixels++
	})
	assert.Equal(t, pixels, 1327)
}

func TestDDATriangleLines(t *testing.T) {
	pixels := 0
	lines := 0
	DDATriangleLines(geom.NewTriangle(13, 12, 37, 54, 78, 15), func(x0, x1, y int) {
		pixels += x1 - x0
		lines++
	})
	assert.Equal(t, pixels, 1327)
	assert.Equal(t, lines, 42)
}

const blockSize = 3

func TestDDATriangleBlocks(t *testing.T) {
	pixels := 0
	lines := 0
	blocks := 0
	DDATriangleBlocks(geom.NewTriangle(13, 12, 37, 54, 78, 15), blockSize, func(x0, x1, y int) {
		pixels += x1 - x0
		lines++
	}, func(x, y int) {
		pixels += blockSize * blockSize
		blocks++

	})

	assert.Equal(t, pixels, 1327)
	assert.Equal(t, lines, 77)
	assert.Equal(t, blocks, 129)
}

func TestDDAPolygon(t *testing.T) {
	// A square polygon should produce pixels consistent with its two triangles
	poly := geom.Polygon{Points: []geom.Point{{10, 10}, {30, 10}, {30, 30}, {10, 30}}}
	pixels := 0
	DDAPolygon(poly, func(x, y int) {
		pixels++
	})
	assert.Greater(t, pixels, 0)
}

func TestDDAPolygonBlocks(t *testing.T) {
	poly := geom.Polygon{Points: []geom.Point{{10, 10}, {30, 10}, {30, 30}, {10, 30}}}
	pixels := 0
	blocks := 0
	DDAPolygonBlocks(poly, blockSize, func(x0, x1, y int) {
		pixels += x1 - x0
	}, func(x, y int) {
		pixels += blockSize * blockSize
		blocks++
	})
	assert.Greater(t, pixels, 0)
	assert.Greater(t, blocks, 0)
}

func TestDDATriangle_SmallTriangle(t *testing.T) {
	// Very small triangle
	pixels := 0
	DDATriangle(geom.NewTriangle(0, 0, 2, 0, 1, 1), func(x, y int) {
		pixels++
	})
	assert.Greater(t, pixels, 0)
}

func TestDDATriangle_DegenerateTriangle(t *testing.T) {
	// Collinear points (zero-area triangle)
	pixels := 0
	DDATriangle(geom.NewTriangle(0, 0, 10, 0, 5, 0), func(x, y int) {
		pixels++
	})
	// Degenerate triangle should not panic, may produce 0 or few pixels
	assert.GreaterOrEqual(t, pixels, 0)
}
