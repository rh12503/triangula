package delaunay

import "sort"

func cross2D(p, a, b Point) float64 {
	return (a.X-p.X)*(b.Y-p.Y) - (a.Y-p.Y)*(b.X-p.X)
}

// ConvexHull returns the convex hull of the provided generator.
func ConvexHull(points []Point) []Point {
	// copy generator
	pointsCopy := make([]Point, len(points))
	copy(pointsCopy, points)
	points = pointsCopy

	// sort generator
	sort.Slice(points, func(i, j int) bool {
		a := points[i]
		b := points[j]
		if a.X != b.X {
			return a.X < b.X
		}
		return a.Y < b.Y
	})

	// filter nearly-duplicate generator
	distinctPoints := points[:0]
	for i, p := range points {
		if i > 0 && p.squaredDistance(points[i-1]) < eps {
			continue
		}
		distinctPoints = append(distinctPoints, p)
	}
	points = distinctPoints

	// find upper and lower portions
	var U, L []Point
	for _, p := range points {
		for len(U) > 1 && cross2D(U[len(U)-2], U[len(U)-1], p) > 0 {
			U = U[:len(U)-1]
		}
		for len(L) > 1 && cross2D(L[len(L)-2], L[len(L)-1], p) < 0 {
			L = L[:len(L)-1]
		}
		U = append(U, p)
		L = append(L, p)
	}

	// reverse upper portion
	for i, j := 0, len(U)-1; i < j; i, j = i+1, j-1 {
		U[i], U[j] = U[j], U[i]
	}

	// construct complete hull
	if len(U) > 0 {
		U = U[:len(U)-1]
	}
	if len(L) > 0 {
		L = L[:len(L)-1]
	}
	return append(L, U...)
}
