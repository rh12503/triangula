// Package image contains a struct to store the data of an image.

package image

import (
	"github.com/RH12503/Triangula/color"
	"image"
)

// Data represents an image
type Data struct {
	width  int
	height int
	pixels [][]color.RGB
}

// RGBAt returns a color.RGB given a coordinate
func (data Data) RGBAt(x, y int) color.RGB {
	return data.pixels[y][x]
}

func (data Data) Size() (int, int) {
	return data.width, data.height
}

// NewData creates a empty Data given a width and height
func NewData(w, h int) Data {
	data := Data{width: w, height: h}
	data.pixels = make([][]color.RGB, h)

	for i := range data.pixels {
		data.pixels[i] = make([]color.RGB, w)
	}

	return data
}

// ToData converts a image.Image to a Data
func ToData(image image.Image) Data {
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

// convertColor is a utility function for converting uint32 RGB values to RGB values between 0 and 1
func convertColor(color uint32) float64 {
	return float64(int(color)/257) / 255
}
