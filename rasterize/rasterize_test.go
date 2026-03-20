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

func TestDDAPolygon_ReturnsPositivePixelCount(t *testing.T) {
	// A quad polygon made of two triangles should produce pixels
	quad := geom.Polygon{Points: []geom.Point{{0, 0}, {20, 0}, {20, 20}, {0, 20}}}
	pixels := 0
	DDAPolygon(quad, func(x, y int) { pixels++ })
	assert.True(t, pixels > 0, "expected DDAPolygon to visit pixels for a non-degenerate polygon")
}

func TestDDAPolygonBlocks_ReturnsPositiveCounts(t *testing.T) {
	quad := geom.Polygon{Points: []geom.Point{{0, 0}, {20, 0}, {20, 20}, {0, 20}}}
	total := 0
	DDAPolygonBlocks(quad, blockSize, func(x0, x1, y int) {
		total += x1 - x0
	}, func(x, y int) {
		total += blockSize * blockSize
	})
	assert.True(t, total > 0, "expected DDAPolygonBlocks to visit pixels for a non-degenerate polygon")
}

func TestDDAPolygonBlocks_TotalMatchesDDAPolygon(t *testing.T) {
	// Both methods should cover the same total number of pixels for the same polygon.
	quad := geom.Polygon{Points: []geom.Point{{3, 3}, {30, 3}, {30, 30}, {3, 30}}}

	pixelCount := 0
	DDAPolygon(quad, func(x, y int) { pixelCount++ })

	blockCount := 0
	DDAPolygonBlocks(quad, blockSize, func(x0, x1, y int) {
		blockCount += x1 - x0
	}, func(x, y int) {
		blockCount += blockSize * blockSize
	})

	assert.Equal(t, pixelCount, blockCount)
}
