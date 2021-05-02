// Package image contains a struct to store the data of an image.
package image

import (
	"github.com/RH12503/Triangula/color"
	"image"
)

// RGBData store pixel data in an image.
type RGBData struct {
	width  int
	height int
	pixels [][]color.RGB
}

// RGBAt returns a color.RGB given a coordinate.
func (data RGBData) RGBAt(x, y int) color.RGB {
	return data.pixels[y][x]
}

// Size returns the dimensions of the image data.
func (data RGBData) Size() (int, int) {
	return data.width, data.height
}

// NewData creates a empty RGBData given a width and height.
func NewData(w, h int) RGBData {
	data := RGBData{width: w, height: h}
	data.pixels = make([][]color.RGB, h)

	for i := range data.pixels {
		data.pixels[i] = make([]color.RGB, w)
	}

	return data
}

// ToData converts an image.Image to a RGBData.
func ToData(image image.Image) RGBData {
	data := NewData(image.Bounds().Max.X, image.Bounds().Max.Y)

	for y := range data.pixels {
		for x := range data.pixels[y] {
			col := &data.pixels[y][x]
			r, g, b, _ := image.At(x, y).RGBA()

			col.R = convertColor(r)
			col.G = convertColor(g)
			col.B = convertColor(b)
		}
	}

	return data
}

// convertColor is a utility function for converting uint32 RGB values to RGB values between 0 and 1.
func convertColor(color uint32) float64 {
	return float64(color >> 8) / 255
}
