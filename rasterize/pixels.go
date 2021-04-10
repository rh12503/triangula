// Package rasterize provides different functions for several use cases to rasterize a triangle.
package rasterize

import (
	"github.com/RH12503/Triangula/geom"
)

// DDATriangle calls function pixel for each pixel a geom.Triangle covers.
func DDATriangle(triangle geom.Triangle, pixel func(x, y int)) {
	DDATriangleLines(triangle, func(x0, x1, y int) {
		for x := x0; x < x1; x++ {
			pixel(x, y)
		}
	})
}
