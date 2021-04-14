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

	evaluator := NewTrianglesImageEvaluator(image2.ToData(img), blockSize)

	assert.Equal(t, evaluator.Calculate(PointsData{
		Points:    normgeom.NormPointGroup{
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
