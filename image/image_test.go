package image

import (
	"github.com/RH12503/Triangula/color"
	"github.com/stretchr/testify/assert"
	stdimage "image"
	stdcolor "image/color"
	"testing"
)

func TestNewData(t *testing.T) {
	data := NewData(100, 50)
	w, h := data.Size()
	assert.Equal(t, w, 100)
	assert.Equal(t, h, 50)
}

func TestToDataBounds(t *testing.T) {
	for _, bounds := range []stdimage.Rectangle{
		stdimage.Rect(0, 0, 2, 3),
		stdimage.Rect(10, 20, 12, 23),
		stdimage.Rect(-3, -4, -1, -1),
		stdimage.Rect(-1, -2, 1, 1),
	} {
		t.Run(bounds.String(), func(t *testing.T) {
			img := stdimage.NewRGBA(bounds)
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					img.SetRGBA(x, y, stdcolor.RGBA{R: uint8((x-bounds.Min.X)*80 + 20), G: uint8((y-bounds.Min.Y)*60 + 10), B: 150, A: 255})
				}
			}
			data := ToData(img)
			w, h := data.Size()
			if w != bounds.Dx() || h != bounds.Dy() {
				t.Fatalf("Size() = (%d, %d), want (%d, %d)", w, h, bounds.Dx(), bounds.Dy())
			}
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					pixel := img.RGBAAt(x+bounds.Min.X, y+bounds.Min.Y)
					want := color.RGB{R: float64(pixel.R) / 255, G: float64(pixel.G) / 255, B: float64(pixel.B) / 255}
					assert.Equal(t, want, data.RGBAt(x, y))
				}
			}
		})
	}
}
