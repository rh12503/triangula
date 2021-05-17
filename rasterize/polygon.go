package rasterize

import "github.com/RH12503/Triangula/geom"

func DDAPolygon(polygon geom.Polygon, blockSize int, line func(x0, x1, y int), block func(x, y int)) {
	polygon.Triangulate(func(triangle geom.Triangle) {
		DDATriangleBlocks(triangle, blockSize, line, block)
	})
}
