package fitness

import (
	image2 "github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/random"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"math"
	"testing"
)

const width, height = 100, 100

const blockSize = 3

func TestTrianglesImageEvaluator_Calculate(t *testing.T) {
	random.Seed(0)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(random.Intn(math.MaxUint8)),
				G: uint8(random.Intn(math.MaxUint8)),
				B: uint8(random.Intn(math.MaxUint8)),
				A: math.MaxUint8,
			})
		}
	}

	evaluator := NewTrianglesImageFunction(image2.ToData(img), blockSize)

	assert.Equal(t, evaluator.Calculate(PointsData{
		Points: normgeom.NormPointGroup{
			{0.12, 0.2},
			{0.73, 0.28},
			{0.57, 0.15},
			{0.23, 0.52},
			{0.13, 0.67},
			{0.34, 0.19},
		},
		Mutations: nil,
	}), 0.16173023698665479)
}

func TestPixels(t *testing.T) {
	pixels := newPixelData(100, 50)
	assert.Equal(t, len(pixels.pixels), 50)
	assert.Equal(t, len(pixels.pixels[0]), 100)

	pixelsN := newPixelDataN(100, 50, blockSize)
	assert.Equal(t, len(pixelsN.pixels), 50-blockSize+1)
	assert.Equal(t, len(pixelsN.pixels[0]), 100-blockSize+1)
}

func TestTriangleCacheData_Equals(t *testing.T) {
	a := &TriangleCacheData{aX: 1, aY: 2, bX: 3, bY: 4, cX: 5, cY: 6, fitness: 0.5}
	b := &TriangleCacheData{aX: 1, aY: 2, bX: 3, bY: 4, cX: 5, cY: 6, fitness: 0.9}
	c := &TriangleCacheData{aX: 10, aY: 20, bX: 3, bY: 4, cX: 5, cY: 6}

	// Same vertices, different fitness — should be equal (equality is vertex-based)
	assert.True(t, a.Equals(b))
	// Different vertices — should not be equal
	assert.False(t, a.Equals(c))
}

func TestTriangleCacheData_Hash(t *testing.T) {
	a := &TriangleCacheData{aX: 1, aY: 2, bX: 3, bY: 4, cX: 5, cY: 6}
	b := &TriangleCacheData{aX: 1, aY: 2, bX: 3, bY: 4, cX: 5, cY: 6}
	c := &TriangleCacheData{aX: 10, aY: 20, bX: 30, bY: 40, cX: 50, cY: 60}

	// Same data should produce same hash
	assert.Equal(t, a.Hash(), b.Hash())
	// Different data should (likely) produce different hash
	assert.NotEqual(t, a.Hash(), c.Hash())
}

func TestTriangleCacheData_DataAndCachedHash(t *testing.T) {
	a := &TriangleCacheData{fitness: 0.75}
	assert.Equal(t, 0.75, a.Data())

	a.SetCachedHash(42)
	assert.Equal(t, uint32(42), a.CachedHash())
}

func TestPolygonCacheData_Equals(t *testing.T) {
	a := &PolygonCacheData{coords: []int16{1, 2, 3, 4}, fitness: 0.5}
	b := &PolygonCacheData{coords: []int16{1, 2, 3, 4}, fitness: 0.9}
	c := &PolygonCacheData{coords: []int16{1, 2, 3, 4, 5}, fitness: 0.5}
	d := &PolygonCacheData{coords: []int16{10, 20, 30, 40}, fitness: 0.5}

	assert.True(t, a.Equals(b))   // Same coords, different fitness
	assert.False(t, a.Equals(c))  // Different lengths
	assert.False(t, a.Equals(d))  // Same length, different values
}

func TestPolygonCacheData_Hash(t *testing.T) {
	a := &PolygonCacheData{coords: []int16{1, 2, 3, 4}}
	b := &PolygonCacheData{coords: []int16{1, 2, 3, 4}}
	c := &PolygonCacheData{coords: []int16{10, 20, 30, 40}}

	assert.Equal(t, a.Hash(), b.Hash())
	assert.NotEqual(t, a.Hash(), c.Hash())
}

func TestPolygonCacheData_DataAndCachedHash(t *testing.T) {
	a := &PolygonCacheData{fitness: 0.33}
	assert.Equal(t, 0.33, a.Data())

	a.SetCachedHash(99)
	assert.Equal(t, uint32(99), a.CachedHash())
}
