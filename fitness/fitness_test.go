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

func TestPolygonsImageEvaluator_Calculate_ReturnsValidFitness(t *testing.T) {
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

	evaluators := PolygonsImageFunctions(image2.ToData(img), blockSize, 1)
	result := evaluators[0].Calculate(PointsData{
		Points: normgeom.NormPointGroup{
			{0.12, 0.2},
			{0.73, 0.28},
			{0.57, 0.15},
			{0.23, 0.52},
			{0.13, 0.67},
			{0.34, 0.19},
		},
		Mutations: nil,
	})

	assert.True(t, result >= 0.0 && result <= 1.0, "polygon fitness %v not in [0,1]", result)
}

func TestTrianglesImageEvaluator_FitnessInRange(t *testing.T) {
	random.Seed(1)

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
	result := evaluator.Calculate(PointsData{
		Points: normgeom.NormPointGroup{
			{0.1, 0.1}, {0.9, 0.1}, {0.5, 0.9}, {0.5, 0.5},
		},
		Mutations: nil,
	})

	assert.True(t, result >= 0.0 && result <= 1.0, "triangle fitness %v not in [0,1]", result)
}

func TestTrianglesImageEvaluator_CachePopulatedAfterCalculate(t *testing.T) {
	random.Seed(2)

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

	fn := NewTrianglesImageFunction(image2.ToData(img), blockSize).(*trianglesImageFunction)
	fn.Calculate(PointsData{
		Points: normgeom.NormPointGroup{
			{0.12, 0.2}, {0.73, 0.28}, {0.57, 0.15},
			{0.23, 0.52}, {0.13, 0.67}, {0.34, 0.19},
		},
		Mutations: nil,
	})

	assert.True(t, len(fn.TriangleCache) > 0, "cache should be populated after Calculate")
}
