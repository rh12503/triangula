package fitness

import (
	"github.com/RH12503/Triangula/image"
)

// pixelData stores the data relating to the pixels of an image, and is used in TrianglesImageEvaluator.
type pixelData struct {
	pixels        [][]pixel
	width, height int
}

// Size returns the width and height of an pixelData
func (p pixelData) Size() (int, int) {
	return p.width, p.height
}

// pixel stores RGB values of a pixel as well as the sum of the squares of them. These values are used in TrianglesImageEvaluator.
type pixel struct {
	r, g, b uint8
	sq      uint32
}

// newPixelData creates a new pixelData given a width and height
func newPixelData(w, h int) pixelData {
	data := pixelData{width: w, height: h}
	data.pixels = make([][]pixel, h)

	for i := range data.pixels {
		data.pixels[i] = make([]pixel, w)
	}

	return data
}

// fromImage creates a pixelData from an image.Data
func fromImage(image image.Data) pixelData {
	w, h := image.Size()
	data := newPixelData(w, h)

	for y := range data.pixels {
		for x := range data.pixels[y] {
			col := &data.pixels[y][x]
			rgb := image.RGBAt(x, y)

			col.r = uint8(rgb.R * 255)
			col.g = uint8(rgb.G * 255)
			col.b = uint8(rgb.B * 255)
			col.sq = uint32(rgb.R*255*rgb.R*255 + rgb.G*255*rgb.G*255 + rgb.B*255*rgb.B*255)
		}
	}

	return data
}

// pixelDataN stores the sum of RGB values of pixels in a N*N block.
// This speeds up performance as the variation can be calculated in blocks instead of each individual pixel
type pixelDataN struct {
	pixels [][]pixelN
}

// pixel stores RGB values of a pixel as well as the sum of the squares of them in a N*N block.
type pixelN struct {
	r, g, b uint16
	sq      uint32
}

// newPixelDataN creates a new pixelDataN given a width and height
func newPixelDataN(w, h, n int) pixelDataN {
	data := pixelDataN{}
	data.pixels = make([][]pixelN, h-n+1)

	for i := range data.pixels {
		data.pixels[i] = make([]pixelN, w-n+1)
	}

	return data
}

// fromImageN creates a pixelDataN from an image.Data. The block size n also needs to be specified.
func fromImageN(image image.Data, n int) pixelDataN {
	w, h := image.Size()
	data := newPixelDataN(w, h, n)

	for y := range data.pixels {
		for x := range data.pixels[y] {
			col := &data.pixels[y][x]
			// Loop through an n*n block and add the values
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					rgb := image.RGBAt(x+i, y+j)
					col.r += uint16(rgb.R * 255)
					col.g += uint16(rgb.G * 255)
					col.b += uint16(rgb.B * 255)
					col.sq += uint32(rgb.R*255*rgb.R*255 + rgb.G*255*rgb.G*255 + rgb.B*255*rgb.B*255)
				}
			}
		}
	}

	return data
}
