// Package render implements utilities for rendering triangles onto an image

package render

import (
	"Triangula/color"
	"Triangula/geom"
	"Triangula/image"
	"Triangula/normgeom"
	"Triangula/rasterize"
)

// TriangleData stores a triangle and it's color
type TriangleData struct {
	Triangle normgeom.NormTriangle
	Color    color.RGB
}

// TrianglesOnImage calculates the optimal color for each triangle so it is the closest to an image
func TrianglesOnImage(triangles []geom.Triangle, image image.Data) []TriangleData {
	triangleData := make([]TriangleData, len(triangles))

	w, h := image.Size()

	for i, t := range triangles {
		// Calculate the average color of all the pixels in the triangle
		var color color.AverageRGB
		c := 0
		rasterize.DDATriangle(t, func(x, y int) {
			color.Add(image.RGBAt(x, y))
			c++
		})

		// If there were no pixels in the triangle, set the color to the nearest pixel (to avoid artifacts)
		if c == 0 {
			for _, p := range t.Points {
				x, y := p.X, p.Y
				if x == w {
					x--
				}

				if y == h {
					y--
				}

				color.Add(image.RGBAt(x, y))
			}
		}

		data := TriangleData{
			Triangle: t.ToNorm(w, h),
			Color:    color.Average(),
		}
		triangleData[i] = data
	}

	return triangleData
}
