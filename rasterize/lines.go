package rasterize

import "github.com/RH12503/Triangula/geom"

// DDATriangleLines calls function line for each horizontal line a geom.Triangle covers
// using a digital differential analyzing algorithm.
func DDATriangleLines(triangle geom.Triangle, line func(x0, x1, y int)) {
	p0 := triangle.Points[0]
	p1 := triangle.Points[1]
	p2 := triangle.Points[2]

	// Sort vertices by height, where y0 has the lowest y value

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

	if y1 == y2 {
		// If the bottom 2 vertices are equal, rasterize a triangle with a flat bottom

		if x2 < x1 {
			x1, x2 = x2, x1
			y1 = y2
		}
		bottomTriangleLines(x1, x2, y1, x0, y0, line)
	} else if y0 == y1 {
		// If the top 2 vertices are equal, rasterize a triangle with a flat top

		if x0 > x1 {
			x0, x1 = x1, x0
			y0 = y1
		}
		topTriangleLines(x0, x1, y0, x2, y2, line)
	} else {
		// If all the y values are different, rasterize the triangle normally

		normalTriangleLines(x0, y0, x1, y1, x2, y2, line)
	}
}

// normalTriangleLines rasterizes a triangle with different y values.
// The y values must be sorted where y0 has the lowest value.
func normalTriangleLines(x0, y0, x1, y1, x2, y2 int, line func(x0, x1, y int)) {

	// Calculate the slopes of the first two lines
	m0 := float64(x2-x0) / float64(y2-y0)
	m1 := float64(x2-x1) / float64(y2-y1)

	// Swap the slopes so m0 is the slope of the left line and m1 is the slope of the right line
	swap := m0 > m1
	if swap {
		m0, m1 = m1, m0
	}

	// Start from the top vertex
	b0 := float64(x2)
	b1 := float64(x2)

	var nX0, nX1 float64

	for i := y2; i < y1; i++ {
		nX0 = m0*float64(i-y2) + b0
		nX1 = m1*float64(i-y2) + b1
		line(int(nX0), int(nX1), i)
	}

	var d0, d1 int

	// One slope will always remain the same, and the second one needs to be calculated
	if swap {
		m0 = float64(x1-x0) / float64(y1-y0)
		b0 = float64(x1)
		d1 = y1 - y2
	} else {
		m1 = float64(x1-x0) / float64(y1-y0)
		b1 = float64(x1)
		d0 = y1 - y2
	}

	for i := y1; i < y0; i++ {
		nX0 = m0*float64(i-y1+d0) + b0
		nX1 = m1*float64(i-y1+d1) + b1
		line(int(nX0), int(nX1), i)
	}
}

// bottomTriangleLines rasterizes a triangle with a flat bottom of coordinate y.
func bottomTriangleLines(x0, x1, y, x2, y2 int, line func(x0, x1, y int)) {
	flatTriangleLines(x2, y2, x0, y, x1-x2, x1, line)
}

// topTriangleLines rasterizes a triangle with a flat top of coordinate y.
func topTriangleLines(x0, x1, y, x2, y2 int, line func(x0, x1, y int)) {
	flatTriangleLines(x0, y, x2, y2, x2-x1, x2, line)
}

// flatTriangleLines rasterizes a triangle with a flat top or bottom.
func flatTriangleLines(x0, y, x2, y2, i2, p int, line func(x0, x1, y int)) {
	m0 := float64(x2-x0) / float64(y2-y)
	m1 := float64(i2) / float64(y2-y)

	fillTriangleLines(y2, y, float64(x2), float64(p), m0, m1, line)
}

// fillTriangleLines rasterizes a flat triangle given the linear equations of its two lines.
func fillTriangleLines(minY, maxY int, lX0, lX1, m0, m1 float64, line func(x0, x1, y int)) {
	for i := minY; i < maxY; i++ {
		nX0 := m0*float64(i-minY) + lX0
		nX1 := m1*float64(i-minY) + lX1
		line(int(nX0), int(nX1), i)
	}
}
