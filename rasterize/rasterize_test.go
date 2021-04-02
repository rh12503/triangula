package rasterize

import (
	"Triangula/geom"
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
		pixels += blockSize*blockSize
		blocks++

	})

	assert.Equal(t, pixels, 1327)
	assert.Equal(t, lines, 77)
	assert.Equal(t, blocks, 129)
}
