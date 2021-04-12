package incrdelaunay

import (
	"math"
)

// CircumcircleGrid is a data structure that uses spatial partitioning to allowed fast operations
// involving multiple Triangle's and their Circumcircle's
type CircumcircleGrid struct {
	triangles            [][][]uint16 // The grid used to store triangles
	cols, rows           int
	rowPixels, colPixels float64 // The number of pixels per row and column
}

// NewCircumcircleGrid returns a new grid with a set amount of columns and rows
func NewCircumcircleGrid(cols, rows, w, h int) CircumcircleGrid {
	c := CircumcircleGrid{}
	c.rows = rows
	c.cols = cols
	c.colPixels = float64(w) / float64(cols)
	c.rowPixels = float64(h) / float64(rows)

	c.triangles = make([][][]uint16, c.rows)
	for i := range c.triangles {
		c.triangles[i] = make([][]uint16, c.cols)
	}
	return c
}

// AddTriangle adds a Triangle with an index to the grid
func (c *CircumcircleGrid) AddTriangle(t Triangle, index uint16) {
	// Find all the boxes of the grid that the triangle's circumcircle intersects
	radius := t.Circumcircle.Radius + 0.001
	topLeftX := int(float64(t.Circumcircle.cX-radius) / c.colPixels)
	topLeftY := int(float64(t.Circumcircle.cY-radius) / c.rowPixels)
	bottomRightX := int(math.Ceil(float64(t.Circumcircle.cX+radius) / c.colPixels))
	bottomRightY := int(math.Ceil(float64(t.Circumcircle.cY+radius) / c.rowPixels))

	if topLeftX < 0 {
		topLeftX = 0
	}

	if topLeftY < 0 {
		topLeftY = 0
	}

	if bottomRightX > c.cols {
		bottomRightX = c.cols
	}

	if bottomRightY > c.rows {
		bottomRightY = c.rows
	}

	for x := topLeftX; x < bottomRightX; x++ {
		col := c.triangles[x]
		for y := topLeftY; y < bottomRightY; y++ {
			col[y] = append(col[y], index)
		}
	}
}

// RemoveTriangle
func (c *CircumcircleGrid) RemoveTriangle(tri Triangle, index uint16) {
	// Find all the boxes of the grid that the triangle's circumcircle intersects
	radius := tri.Circumcircle.Radius + 0.001
	topLeftX := int(float64(tri.Circumcircle.cX-radius) / c.colPixels)
	topLeftY := int(float64(tri.Circumcircle.cY-radius) / c.rowPixels)
	bottomRightX := int(math.Ceil(float64(tri.Circumcircle.cX+radius) / c.colPixels))
	bottomRightY := int(math.Ceil(float64(tri.Circumcircle.cY+radius) / c.rowPixels))

	if topLeftX < 0 {
		topLeftX = 0
	}

	if topLeftY < 0 {
		topLeftY = 0
	}

	if bottomRightX > c.cols {
		bottomRightX = c.cols
	}

	if bottomRightY > c.rows {
		bottomRightY = c.rows
	}

	for x := topLeftX; x < bottomRightX; x++ {
		col := c.triangles[x]
		for y := topLeftY; y < bottomRightY; y++ {
			for i, t := range col[y] {

				if t == index {
					in := len(col[y]) - 1
					col[y][i] = col[y][in]
					col[y] = col[y][:in]
					break
				}
			}
		}
	}
}

// HasPoint returns if a triangle in the grid has a point
func (c CircumcircleGrid) HasPoint(p Point, triangles []Triangle) bool {
	// Find which box of the grid the point falls into
	x := int(math.Floor(float64(p.X) / c.colPixels))
	y := int(math.Floor(float64(p.Y) / c.rowPixels))

	if x == c.cols {
		x = c.cols - 1
	}

	if y == c.rows {
		y = c.rows - 1
	}

	group := c.triangles[x][y]
	size := len(group)

	for i := 0; i < size; i++ {
		t := group[i]

		tri := triangles[t]
		if tri.A.X == -1 {
			panic("UH OH")
		}
		if tri.HasVertex(p) {
			return true
		}
	}
	return false
}

// RemoveCircumcirclesThatContain removes all triangles whose circumcircle's contain a point
func (c CircumcircleGrid) RemoveCircumcirclesThatContain(p Point, triangles []Triangle, contains func(i uint16)) {
	// Find which box of the grid the point falls into
	x := int(math.Floor(float64(p.X) / c.colPixels))
	y := int(math.Floor(float64(p.Y) / c.rowPixels))

	if x == c.cols {
		x = c.cols - 1
	}

	if y == c.rows {
		y = c.rows - 1
	}

	group := c.triangles[x][y]
	size := len(group)

	for i := 0; i < size; i++ {
		t := group[i]

		tri := triangles[t]
		if tri.A.X == -1 {
			panic("UH OH")
		}
		if inCircle(int64(tri.A.X), int64(tri.A.Y), int64(tri.B.X), int64(tri.B.Y), int64(tri.C.X), int64(tri.C.Y), int64(p.X), int64(p.Y)) >= 0 {
			contains(t)

			c.RemoveTriangle(tri, t)
			i--
			size--
		}
	}
}

// RemoveThatHasVertex removes all triangles that have a point
func (c CircumcircleGrid) RemoveThatHasVertex(p Point, triangles []Triangle, contains func(i uint16)) {
	// Find which box of the grid the point falls into
	x := int(math.Floor(float64(p.X) / c.colPixels))
	y := int(math.Floor(float64(p.Y) / c.rowPixels))

	if x == c.cols {
		x = c.cols - 1
	}

	if y == c.rows {
		y = c.rows - 1
	}

	group := c.triangles[x][y]
	size := len(group)

	for i := 0; i < size; i++ {
		t := group[i]

		tri := triangles[t]
		if tri.HasVertex(p) {
			contains(t)

			c.RemoveTriangle(tri, t)
			i--
			size--
		}
	}
}

// Set sets a CircumcircleGrid to another CircumcircleGrid
func (c *CircumcircleGrid) Set(other *CircumcircleGrid) {
	for x, col := range c.triangles {
		for y := range col {
			c.triangles[x][y] = c.triangles[x][y][:cap(c.triangles[x][y])]
			if len(c.triangles[x][y]) > len(other.triangles[x][y]) {
				c.triangles[x][y] = c.triangles[x][y][:len(other.triangles[x][y])]
			} else if len(c.triangles[x][y]) < len(other.triangles[x][y]) {
				c.triangles[x][y] = make([]uint16, len(other.triangles[x][y]))
			}

			copy(c.triangles[x][y], other.triangles[x][y])
		}
	}
}
