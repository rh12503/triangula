package incrdelaunay

import (
	"math"
	"sort"
)

func NewVoronoi(w, h int) *IVoronoi {
	voronoi := IVoronoi{
		delaunay:   NewDelaunay(w, h),
		w:          w,
		h:          h,
		polygonMap: newPolygonMap(6),
	}
	return &voronoi
}

type IVoronoi struct {
	delaunay   *Delaunay
	w, h       int
	points     []FloatPoint
	newPolygon []FloatPoint
	polygonMap polygonMap

	pointsToUpdate []Point
	i              int
}

func (v *IVoronoi) Insert(point Point) {
	if v.delaunay.pointMap.CopiesOf(point) > 0 {
		v.delaunay.Insert(point)
		return
	}

	v.pointsToUpdate = v.pointsToUpdate[:0]
	triangles := v.delaunay.triangles

	v.delaunay.grid.IterCircumcirclesThatContain(point, triangles, func(i uint16) {
		t := triangles[i]
		v.addPointToUpdate(t.A)
		v.addPointToUpdate(t.B)
		v.addPointToUpdate(t.C)
	})

	v.delaunay.Insert(point)

	for _, p := range v.pointsToUpdate {
		if !v.polygonMap.RemovePolygon(p) {
			panic("oh no")
		}
	}

	v.addPointToUpdate(point)
	for _, p := range v.pointsToUpdate {
		v.processPoint(p)
	}
}

func (v *IVoronoi) Remove(point Point) {
	if v.delaunay.pointMap.CopiesOf(point) > 1 {
		v.delaunay.Remove(point)
		return
	}

	v.pointsToUpdate = v.pointsToUpdate[:0]
	triangles := v.delaunay.triangles

	v.delaunay.grid.IterThatHasVertex(point, triangles, func(i uint16) {
		t := triangles[i]

		v.addPointToUpdate(t.A)
		v.addPointToUpdate(t.B)
		v.addPointToUpdate(t.C)
	})

	v.delaunay.Remove(point)

	for _, p := range v.pointsToUpdate {
		if !v.polygonMap.RemovePolygon(p) {
			panic("uh oh")
		}
	}

	for _, p := range v.pointsToUpdate {
		if p != point {
			v.processPoint(p)
		}
	}
}

func (v *IVoronoi) IterPolygons(polygon func([]FloatPoint)) {
	v.polygonMap.IterPolygons(func(poly []FloatPoint) {
		if len(poly) > 0 {
			polygon(poly)
		}
	})
}

func (v *IVoronoi) processPoint(point Point) {
	triangles := v.delaunay.triangles

	v.points = v.points[:0]
	clip := false

	v.delaunay.grid.IterThatHasVertex(point, triangles, func(i uint16) {
		circ := triangles[i].Circumcircle

		new := FloatPoint{
			X: float64(circ.cX),
			Y: float64(circ.cY),
		}

		if outside(new, v.w, v.h) {
			clip = true
		}
		found := false
		for _, p := range v.points {
			if p == new {
				found = true
				break
			}
		}

		if !found {
			v.points = append(v.points, new)
		}
	})

	center := FloatPoint{
		X: float64(point.X),
		Y: float64(point.Y),
	}

	sort.Slice(v.points, func(i, j int) bool {
		a := v.points[j]
		b := v.points[i]

		if a.X-center.X >= 0 && b.X-center.X < 0 {
			return true
		}
		if a.X-center.X < 0 && b.X-center.X >= 0 {
			return false
		}
		if a.X-center.X == 0 && b.X-center.X == 0 {
			if a.Y-center.Y >= 0 || b.Y-center.Y >= 0 {
				return a.Y > b.Y
			}
			return b.Y > a.Y
		}

		det := (a.X-center.X)*(b.Y-center.Y) - (b.X-center.X)*(a.Y-center.Y)
		if det < 0 {
			return true
		}
		if det > 0 {
			return false
		}

		panic("...")
	})

	if clip {
		v.newPolygon = v.newPolygon[:0]

		intersections := [4]int{-1, -1, -1, -1}

		for i := 0; i < len(v.points); i++ {
			a := v.points[i]
			b := v.points[(i+1)%len(v.points)]

			aIn := !outside(a, v.w, v.h)

			if aIn {
				v.newPolygon = append(v.newPolygon, a)
			}

			for j, e := range edges {
				e0 := FloatPoint{
					X: e[0].X * float64(v.w),
					Y: e[0].Y * float64(v.h),
				}
				e1 := FloatPoint{
					X: e[1].X * float64(v.w),
					Y: e[1].Y * float64(v.h),
				}

				x, y, t := segmentsIntersect(
					e[0].X*float64(v.w),
					e[0].Y*float64(v.h),
					e[1].X*float64(v.w),
					e[1].Y*float64(v.h),
					a.X,
					a.Y,
					b.X,
					b.Y,
				)

				if t && !(epsilonEquals(x, a.X) && epsilonEquals(y, a.Y)) && !(epsilonEquals(x, b.X) && epsilonEquals(y, b.Y)) {
					new := FloatPoint{
						X: x,
						Y: y,
					}
					if new != e0 && new != e1 {
						v.newPolygon = append(v.newPolygon, new)
					}
					intersections[j] += 1
				}
			}

		}

		for j, e := range edges {
			i0 := intersections[j]
			i1 := intersections[(j+1)%len(edges)]

			if i0 != -1 && i1 != -1 {
				corner := FloatPoint{
					X: e[1].X * float64(v.w),
					Y: e[1].Y * float64(v.h),
				}
				if polygonPointIntersect(v.points, corner) {
					v.newPolygon = append(v.newPolygon, corner)
				}
			}
		}

		center := FloatPoint{}

		for _, p := range v.newPolygon {
			center.X += p.X
			center.Y += p.Y
		}
		n := float64(len(v.newPolygon))
		center.X /= n
		center.Y /= n

		sort.Slice(v.newPolygon, func(i, j int) bool {
			a := v.newPolygon[j]
			b := v.newPolygon[i]

			if a.X-center.X >= 0 && b.X-center.X < 0 {
				return true
			}
			if a.X-center.X < 0 && b.X-center.X >= 0 {
				return false
			}
			if a.X-center.X == 0 && b.X-center.X == 0 {
				if a.Y-center.Y >= 0 || b.Y-center.Y >= 0 {
					return a.Y > b.Y
				}
				return b.Y > a.Y
			}

			det := (a.X-center.X)*(b.Y-center.Y) - (b.X-center.X)*(a.Y-center.Y)
			if det < 0 {
				return true
			}
			if det > 0 {
				return false
			}

			return false
		})
		//if len(v.newPolygon) != 0 {
		var new []FloatPoint
		new = append(new, v.newPolygon...)

		v.polygonMap.AddPolygon(point, new)
		//}
	} else {
		//if len(v.points) != 0 {
		var new []FloatPoint
		new = append(new, v.points...)

		v.polygonMap.AddPolygon(point, new)
		//}
	}
}

func (v *IVoronoi) addPointToUpdate(point Point) {
	if v.delaunay.superTriangle.HasVertex(point) {
		return
	}

	found := false
	for _, p := range v.pointsToUpdate {
		if p == point {
			found = true
			break
		}
	}

	if !found {
		v.pointsToUpdate = append(v.pointsToUpdate, point)
	}
}

func (v *IVoronoi) Set(other *IVoronoi) {
	v.delaunay.Set(other.delaunay)
	v.polygonMap.Set(&other.polygonMap)
}

// Should not be modified
var edges = [4][2]FloatPoint{
	{FloatPoint{0, 0}, FloatPoint{1, 0}},
	{FloatPoint{1, 0}, FloatPoint{1, 1}},
	{FloatPoint{1, 1}, FloatPoint{0, 1}},
	{FloatPoint{0, 1}, FloatPoint{0, 0}},
}

func Voronoi(delaunay *Delaunay, polygon func([]FloatPoint), w, h int) {
	triangles := delaunay.triangles
	points := make([]FloatPoint, 0, 8)

	newPolygon := make([]FloatPoint, 0, 8)

	delaunay.pointMap.IterPoints(func(point Point) {
		points = points[:0]
		clip := false

		delaunay.grid.IterThatHasVertex(point, triangles, func(i uint16) {
			circ := triangles[i].Circumcircle

			new := FloatPoint{
				X: float64(circ.cX),
				Y: float64(circ.cY),
			}

			if outside(new, w, h) {
				clip = true
			}
			found := false
			for _, p := range points {
				if p == new {
					found = true
					break
				}
			}

			if !found {
				points = append(points, new)
			}
		})

		center := FloatPoint{
			X: float64(point.X),
			Y: float64(point.Y),
		}

		sort.Slice(points, func(i, j int) bool {
			a := points[j]
			b := points[i]

			if a.X-center.X >= 0 && b.X-center.X < 0 {
				return true
			}
			if a.X-center.X < 0 && b.X-center.X >= 0 {
				return false
			}
			if a.X-center.X == 0 && b.X-center.X == 0 {
				if a.Y-center.Y >= 0 || b.Y-center.Y >= 0 {
					return a.Y > b.Y
				}
				return b.Y > a.Y
			}

			det := (a.X-center.X)*(b.Y-center.Y) - (b.X-center.X)*(a.Y-center.Y)
			if det < 0 {
				return true
			}
			if det > 0 {
				return false
			}

			panic("...")
		})

		if clip {
			newPolygon = newPolygon[:0]

			intersections := [4]int{-1, -1, -1, -1}

			for i := 0; i < len(points); i++ {
				a := points[i]
				b := points[(i+1)%len(points)]

				aIn := !outside(a, w, h)

				if aIn {
					newPolygon = append(newPolygon, a)
				}

				for j, e := range edges {
					e0 := FloatPoint{
						X: e[0].X * float64(w),
						Y: e[0].Y * float64(h),
					}
					e1 := FloatPoint{
						X: e[1].X * float64(w),
						Y: e[1].Y * float64(h),
					}

					x, y, t := segmentsIntersect(
						e[0].X*float64(w),
						e[0].Y*float64(h),
						e[1].X*float64(w),
						e[1].Y*float64(h),
						a.X,
						a.Y,
						b.X,
						b.Y,
					)

					if t && !(epsilonEquals(x, a.X) && epsilonEquals(y, a.Y)) && !(epsilonEquals(x, b.X) && epsilonEquals(y, b.Y)) {
						new := FloatPoint{
							X: x,
							Y: y,
						}
						if new != e0 && new != e1 {
							newPolygon = append(newPolygon, new)
						}
						intersections[j] += 1
					}
				}

			}

			for j, e := range edges {
				i0 := intersections[j]
				i1 := intersections[(j+1)%len(edges)]

				if i0 != -1 && i1 != -1 {
					corner := FloatPoint{
						X: e[1].X * float64(w),
						Y: e[1].Y * float64(h),
					}
					if polygonPointIntersect(points, corner) {
						newPolygon = append(newPolygon, corner)
					}
				}
			}

			center := FloatPoint{}

			for _, p := range newPolygon {
				center.X += p.X
				center.Y += p.Y
			}
			n := float64(len(newPolygon))
			center.X /= n
			center.Y /= n

			sort.Slice(newPolygon, func(i, j int) bool {
				a := newPolygon[j]
				b := newPolygon[i]

				if a.X-center.X >= 0 && b.X-center.X < 0 {
					return true
				}
				if a.X-center.X < 0 && b.X-center.X >= 0 {
					return false
				}
				if a.X-center.X == 0 && b.X-center.X == 0 {
					if a.Y-center.Y >= 0 || b.Y-center.Y >= 0 {
						return a.Y > b.Y
					}
					return b.Y > a.Y
				}

				det := (a.X-center.X)*(b.Y-center.Y) - (b.X-center.X)*(a.Y-center.Y)
				if det < 0 {
					return true
				}
				if det > 0 {
					return false
				}

				return false
			})

			if len(newPolygon) != 0 {
				polygon(newPolygon)
			}
		} else {
			if len(points) != 0 {
				polygon(points)
			}
		}
	})
}

const epsilon = 0.0001

func epsilonEquals(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func segmentsIntersect(p0_x, p0_y, p1_x, p1_y,
	p2_x, p2_y, p3_x, p3_y float64) (float64, float64, bool) {
	var s1_x, s1_y, s2_x, s2_y float64
	s1_x = p1_x - p0_x
	s1_y = p1_y - p0_y
	s2_x = p3_x - p2_x
	s2_y = p3_y - p2_y

	var s, t float64
	s = (-s1_y*(p0_x-p2_x) + s1_x*(p0_y-p2_y)) / (-s2_x*s1_y + s1_x*s2_y)
	t = (s2_x*(p0_y-p2_y) - s2_y*(p0_x-p2_x)) / (-s2_x*s1_y + s1_x*s2_y)

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		return p0_x + (t * s1_x), p0_y + (t * s1_y), true
	}

	return 0, 0, false
}

func outside(p FloatPoint, w, h int) bool {
	return p.X < -epsilon || p.Y < -epsilon || p.X > float64(w)+epsilon || p.Y > float64(h)+epsilon
}

func polygonPointIntersect(polygon []FloatPoint, point FloatPoint) bool {
	//n>2 Keep track of cross product sign changes
	pos := 0
	neg := 0

	for i := 0; i < len(polygon); i++ {
		//If point is in the polygon
		if polygon[i] == point {
			return true
		}

		//Form a segment between the i'th point
		var x1 = polygon[i].X
		var y1 = polygon[i].Y

		//And the i+1'th, or if i is the last, with the first point
		var i2 = (i + 1) % len(polygon)

		var x2 = polygon[i2].X
		var y2 = polygon[i2].Y

		var x = point.X
		var y = point.Y

		//Compute the cross product
		var d = (x-x1)*(y2-y1) - (y-y1)*(x2-x1)

		if d > 0 {
			pos++
		}
		if d < 0 {
			neg++
		}

		//If the sign changes, then point is outside
		if pos > 0 && neg > 0 {
			return false
		}
	}

	//If no change in direction, then on same side of all segments, and thus inside
	return true
}
