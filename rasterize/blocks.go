package rasterize

import (
	"Triangula/geom"
	"math"
)

// DDATriangle calls function line for each line a geom.Triangle covers.
// It calls function block for each blockSize by blockSize block the triangle covers
func DDATriangleBlocks(triangle geom.Triangle, blockSize int, line func(x0, x1, y int), block func(x, y int)) {
	p0 := triangle.Points[0]
	p1 := triangle.Points[1]
	p2 := triangle.Points[2]

	// Sort vertices by y value, where y0 has the lowest value

	x0, y0 := p0.X, p0.Y
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y

	if y1 > y0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	if y2 > y1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1

		if y1 > y0 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}
	}

	normalTriangleBlocks(x0, y0, x1, y1, x2, y2, blockSize, line, block)
}

// normalTriangleBlocks rasterizes the lines of a triangle, while trying to rasterize in blocks when possible
func normalTriangleBlocks(x0, y0, x1, y1, x2, y2 int, blockSize int, line func(x0, x1, y int), block func(x, y int)) {

	// Calculate the slopes of the first two lines
	m0 := float64(x2-x0) / float64(y2-y0)
	m1 := float64(x2-x1) / float64(y2-y1)

	// Swap the slopes so m0 is the slope of the left line and m1 is the slope of the right line
	swap := m0 > m1
	if swap {
		m0, m1 = m1, m0
	}

	// Start from the top vertex
	lX0 := float64(x2)
	lX1 := float64(x2)

	// Starting from the bottom y coordinate, iterate upwards through the pixels using incrementing by blockSize
	i := y1 - 1
	for ; i > y2; i -= blockSize {
		top := i - blockSize + 1

		bottomX := m0*float64(i-y2) + lX0
		topX := m0*float64((top)-y2) + lX0
		maxX := math.Max(bottomX, topX)

		bottomX = m1*float64(i-y2) + lX1
		topX = m1*float64((top)-y2) + lX1
		minX := math.Min(bottomX, topX)

		// Leave the loop if the remaining triangle isn't wide enough to rasterize blocks
		if float64(int(maxX)+blockSize) >= minX {
			break
		}

		// Fill in the left section of the triangle where blocks can't be rasterized
		for y := 0; y < blockSize; y++ {
			pX0 := m0*float64((i-y)-y2) + lX0
			line(int(pX0), int(maxX), i-y)
		}

		// Rasterize the middle section of the triangle in blocks
		x := int(maxX)
		for ; float64(x+blockSize) < minX; x += blockSize {
			block(x, i-blockSize+1)
		}

		// Fill in the right section of the triangle where blocks can't be rasterized
		for y := 0; y < blockSize; y++ {
			pX1 := m1*float64((i-y)-y2) + lX1
			line(x, int(pX1), i-y)
		}
	}

	// Rasterize the remaining part of the top triangle with pixels
	for ; i > y2; i-- {
		pX0 := m0*float64(i-y2) + lX0
		pX1 := m1*float64(i-y2) + lX1

		line(int(pX0), int(pX1), i)
	}

	// Calculate the new slope for the line that changed, and repeat the process above

	var d0, d1 int

	if swap {
		m0 = float64(x1-x0) / float64(y1-y0)
		lX0 = float64(x1)
		d1 = y1 - y2
	} else {
		m1 = float64(x1-x0) / float64(y1-y0)
		lX1 = float64(x1)
		d0 = y1 - y2
	}

	if y1 == y2 {
		lX0 = float64(x2)
		lX1 = float64(x1)

		if lX0 > lX1 {
			lX0, lX1 = lX1, lX0
		}
		if m0 < m1 {
			m0, m1 = m1, m0
		}
	}

	i = y1

	// Starting from the top y coordinate, iterate downwards through the pixels using incrementing by blockSize
	for ; i+blockSize < y0; i += blockSize {
		top := i + blockSize - 1

		bottomX := m0*float64(i-y1+d0) + lX0
		topX := m0*float64(top-y1+d0) + lX0
		maxX := math.Max(bottomX, topX)

		bottomX = m1*float64(i-y1+d1) + lX1
		topX = m1*float64(top-y1+d1) + lX1
		minX := math.Min(bottomX, topX)

		// Leave the loop if the remaining triangle isn't wide enough to rasterize blocks
		if float64(int(maxX)+blockSize) >= minX {
			break
		}

		// Fill in the right section of the triangle where blocks can't be rasterized
		for y := 0; y < blockSize; y++ {
			pX0 := m0*float64((i+y)-y1+d0) + lX0
			line(int(pX0), int(maxX), i+y)
		}

		// Rasterize the middle section of the triangle in blocks
		x := int(maxX)
		for ; float64(x+blockSize) < minX; x += blockSize {
			block(x, i)
		}

		// Fill in the right section of the triangle where blocks can't be rasterized
		for y := 0; y < blockSize; y++ {
			pX1 := m1*float64((i+y)-y1+d1) + lX1
			line(x, int(pX1), i+y)
		}
	}

	// Rasterize the remaining part of the bottom triangle with pixels
	for ; i < y0; i++ {
		pX0 := m0*float64(i-y1+d0) + lX0
		pX1 := m1*float64(i-y1+d1) + lX1
		line(int(pX0), int(pX1), i)
	}
}
