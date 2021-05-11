// Package render implements utilities for rendering triangles onto an image.
package render

import (
	"github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/rasterize"
)

// TriangleData stores a triangle and its color.
type TriangleData struct {
	Triangle normgeom.NormTriangle
	Color    color.RGB
}

// TrianglesOnImage calculates the optimal color for a group of triangles so the colors of triangles
// are closest to an image.
func TrianglesOnImage(triangles []geom.Triangle, image image.Data) []TriangleData {
	triangleData := make([]TriangleData, len(triangles))

	w, h := image.Size()

	for i, t := range triangles {
		// Calculate the average color of all the pixels in the triangle
		var color color.AverageRGB

		rasterize.DDATriangle(t, func(x, y int) {
			color.Add(image.RGBAt(x, y))
		})

		// If there were no pixels in the triangle, set the color to the nearest pixel (to avoid artifacts)
		if color.Count() == 0 {
			for _, p := range t.Points {
				x, y := min(p.X, w-1), min(p.Y, h-1)

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

