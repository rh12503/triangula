package draw

type point struct {
	x, y int
}

// ellipseArea returns a list of consecutive point pairs. Each pair lies on a
// horizontal line (i.e. both generator have the same y position) and if you draw
// horizontal pixels lines for all pairs, you will have the requested ellipse.
//
//         c···d
//      g·········h
//     i···········j
//      e·········f
//         a···b
//
func ellipseArea(x, y, w, h int) (p []point) {
	quarter := quaterEllipsePoints(w, h)
	xPivot, yPivot := 0, 0
	if w%2 == 0 {
		xPivot = 1
	}
	if h%2 == 0 {
		yPivot = 1
	}
	dx, dy := x+w/2, y+h/2
	for i := 0; i < len(quarter); i++ {
		if i == len(quarter)-1 || quarter[i].y != quarter[i+1].y {
			p = append(p,
				// this line
				point{
					x: -quarter[i].x - xPivot + dx,
					y: quarter[i].y + dy,
				},
				point{
					x: quarter[i].x + dx,
					y: quarter[i].y + dy,
				},
				// the line mirrored in y
				point{
					x: -quarter[i].x - xPivot + dx,
					y: -quarter[i].y - yPivot + dy,
				},
				point{
					x: quarter[i].x + dx,
					y: -quarter[i].y - yPivot + dy,
				},
			)
		}
	}
	// remove the last line if it is contained twice at the end
	n := len(p)
	if n >= 4 && p[n-1] == p[n-3] {
		p = p[:n-2]
	}
	return
}

// ellipseOutline returns a list of pixel positions that mark the outline of the
// requested ellipse.
//
//       jih
//     lk   gf
//    m       e
//     no   cd
//       pab
//
func ellipseOutline(x, y, w, h int) []point {
	quarter := quaterEllipsePoints(w, h)
	xPivot, yPivot := 0, 0
	if w%2 == 0 {
		xPivot = 1
	}
	if h%2 == 0 {
		yPivot = 1
	}
	dx, dy := x+w/2, y+h/2
	p := make([]point, 0, len(quarter)*4)
	for i := range quarter {
		p = append(p, point{
			x: quarter[i].x + dx,
			y: quarter[i].y + dy,
		})
	}
	for i := len(quarter) - 1 - (1 - yPivot); i >= 0; i-- {
		p = append(p, point{
			x: quarter[i].x + dx,
			y: -quarter[i].y - yPivot + dy,
		})
	}
	for i := 1 - xPivot; i < len(quarter); i++ {
		p = append(p, point{
			x: -quarter[i].x - xPivot + dx,
			y: -quarter[i].y - yPivot + dy,
		})
	}
	for i := len(quarter) - 1 - (1 - yPivot); i >= 1-xPivot; i-- {
		p = append(p, point{
			x: -quarter[i].x - xPivot + dx,
			y: quarter[i].y + dy,
		})
	}
	return p
}

func quaterEllipsePoints(w, h int) (p []point) {
	if w <= 0 || h <= 0 {
		return nil
	}

	a, b := (w-1)/2, (h-1)/2
	x, y := 0, b
	a2, b2 := a*a, b*b

	crit1 := -(a2/4 + a%2 + b2)
	crit2 := -(b2/4 + b%2 + a2)
	crit3 := -(b2/4 + b%2)
	t := -a2 * y
	dxt := 2 * b2 * x
	dyt := -2 * a2 * y
	d2xt := 2 * b2
	d2yt := 2 * a2

	for y >= 0 && x <= a {
		p = append(p, point{x: x, y: y})
		if t+b2*x <= crit1 || t+a2*y <= crit3 {
			x++
			dxt += d2xt
			t += dxt
		} else if t-a2*y > crit2 {
			y--
			dyt += d2yt
			t += dyt
		} else {
			x++
			dxt += d2xt
			t += dxt
			y--
			dyt += d2yt
			t += dyt
		}
	}
	return
}
