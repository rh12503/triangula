// Package delaunay is almost completely identical to fogleman's excellent Delaunay triangulation
// library, with some modifications for utility and speed.
//
// Full credits to: https://github.com/fogleman/delaunay
package delaunay

import (
	"fmt"
	"math"
	"sort"
)

type triangulator struct {
	points           []Point
	squaredDistances []float64
	ids              []int
	center           Point
	triangles        []int
	halfedges        []int
	trianglesLen     int
	hull             *node
	hash             []*node
}

func newTriangulator(points []Point) *triangulator {
	t := &triangulator{points: points}
	n := len(points)
	t.ids = make([]int, n)
	t.squaredDistances = make([]float64, n)
	return t
}

// sorting a triangulator sorts the `ids` such that the referenced generator
// are in order by their distance to `center`
func (a *triangulator) Len() int {
	return len(a.points)
}

func (a *triangulator) Swap(i, j int) {
	a.ids[i], a.ids[j] = a.ids[j], a.ids[i]
}

func (a *triangulator) Less(i, j int) bool {
	d1 := a.squaredDistances[a.ids[i]]
	d2 := a.squaredDistances[a.ids[j]]
	if d1 != d2 {
		return d1 < d2
	}
	p1 := a.points[a.ids[i]]
	p2 := a.points[a.ids[j]]
	if p1.X != p2.X {
		return p1.X < p2.X
	}
	return p1.Y < p2.Y
}

func (tri *triangulator) update() error {
	tri.trianglesLen = 0
	points := tri.points

	n := len(points)
	if n == 0 {
		return nil
	}

	// compute bounds
	x0 := points[0].X
	y0 := points[0].Y
	x1 := points[0].X
	y1 := points[0].Y
	for i, p := range points {
		if p.X < x0 {
			x0 = p.X
		}
		if p.X > x1 {
			x1 = p.X
		}
		if p.Y < y0 {
			y0 = p.Y
		}
		if p.Y > y1 {
			y1 = p.Y
		}
		tri.ids[i] = i
	}

	var i0, i1, i2 int

	// pick a seed point close to midpoint
	m := Point{(x0 + x1) / 2, (y0 + y1) / 2}
	minDist := infinity
	for i, p := range points {
		d := p.squaredDistance(m)
		if d < minDist {
			i0 = i
			minDist = d
		}
	}

	// find point closest to seed point
	minDist = infinity
	for i, p := range points {
		if i == i0 {
			continue
		}
		d := p.squaredDistance(points[i0])
		if d > 0 && d < minDist {
			i1 = i
			minDist = d
		}
	}

	// find the third point which forms the smallest circumcircle
	minRadius := infinity
	for i, p := range points {
		if i == i0 || i == i1 {
			continue
		}
		r := circumradius(points[i0], points[i1], p)
		if r < minRadius {
			i2 = i
			minRadius = r
		}
	}
	if minRadius == infinity {
		return fmt.Errorf("No Delaunay triangulation exists for this input.")
	}

	// swap the order of the seed generator for counter-clockwise orientation
	if area(points[i0], points[i1], points[i2]) < 0 {
		i1, i2 = i2, i1
	}

	tri.center = circumcenter(points[i0], points[i1], points[i2])

	// sort the generator by distance from the seed triangle circumcenter

	for i, p := range points {
		tri.squaredDistances[i] = p.squaredDistance(tri.center)
	}
	sort.Sort(tri)

	// initialize a hash table for storing edges of the advancing convex hull
	hashSize := int(math.Ceil(math.Sqrt(float64(n))))
	tri.hash = make([]*node, hashSize)

	// initialize a circular doubly-linked list that will hold an advancing convex hull
	nodes := make([]node, n)

	e := newNode(nodes, i0, nil)
	e.t = 0
	tri.hashEdge(e)

	e = newNode(nodes, i1, e)
	e.t = 1
	tri.hashEdge(e)

	e = newNode(nodes, i2, e)
	e.t = 2
	tri.hashEdge(e)

	tri.hull = e

	maxTriangles := 2*n - 5
	tri.triangles = make([]int, maxTriangles*3)
	tri.halfedges = make([]int, maxTriangles*3)

	tri.addTriangle(i0, i1, i2, -1, -1, -1)

	pp := Point{infinity, infinity}
	for k := 0; k < n; k++ {
		i := tri.ids[k]
		p := points[i]

		// skip nearly-duplicate generator
		if p.squaredDistance(pp) < eps {
			continue
		}
		pp = p

		// skip seed triangle generator
		if i == i0 || i == i1 || i == i2 {
			continue
		}

		// find a visible edge on the convex hull using edge hash
		var start *node
		key := tri.hashKey(p)
		for j := 0; j < len(tri.hash); j++ {
			start = tri.hash[key]
			if start != nil && start.i >= 0 {
				break
			}
			key++
			if key >= len(tri.hash) {
				key = 0
			}
		}
		start = start.prev

		e := start
		for area(p, points[e.i], points[e.next.i]) >= 0 {
			e = e.next
			if e == start {
				e = nil
				break
			}
		}
		if e == nil {
			// likely a near-duplicate point; skip it
			continue
		}
		walkBack := e == start

		// add the first triangle from the point
		t := tri.addTriangle(e.i, i, e.next.i, -1, -1, e.t)
		e.t = t // keep track of boundary triangles on the hull
		e = newNode(nodes, i, e)

		// recursively flip triangles from the point until they satisfy the Delaunay condition
		e.t = tri.legalize(t + 2)

		// walk forward through the hull, adding more triangles and flipping recursively
		q := e.next
		for area(p, points[q.i], points[q.next.i]) < 0 {
			t = tri.addTriangle(q.i, i, q.next.i, q.prev.t, -1, q.t)
			q.prev.t = tri.legalize(t + 2)
			tri.hull = q.remove()
			q = q.next
		}

		if walkBack {
			// walk backward from the other side, adding more triangles and flipping
			q := e.prev
			for area(p, points[q.prev.i], points[q.i]) < 0 {
				t = tri.addTriangle(q.prev.i, i, q.i, -1, q.t, q.prev.t)
				tri.legalize(t + 2)
				q.prev.t = t
				tri.hull = q.remove()
				q = q.prev
			}
		}

		// save the two new edges in the hash table
		tri.hashEdge(e)
		tri.hashEdge(e.prev)
	}

	tri.triangles = tri.triangles[:tri.trianglesLen]
	tri.halfedges = tri.halfedges[:tri.trianglesLen]

	return nil
}

func (t *triangulator) hashKey(point Point) int {
	d := point.sub(t.center)
	return int(pseudoAngle(d.X, d.Y) * float64(len(t.hash)))
}

func (t *triangulator) hashEdge(e *node) {
	t.hash[t.hashKey(t.points[e.i])] = e
}

func (t *triangulator) addTriangle(i0, i1, i2, a, b, c int) int {
	i := t.trianglesLen
	t.triangles[i] = i0
	t.triangles[i+1] = i1
	t.triangles[i+2] = i2
	t.link(i, a)
	t.link(i+1, b)
	t.link(i+2, c)
	t.trianglesLen += 3
	return i
}

func (t *triangulator) link(a, b int) {
	t.halfedges[a] = b
	if b >= 0 {
		t.halfedges[b] = a
	}
}

func (t *triangulator) legalize(a int) int {
	// if the pair of triangles doesn't satisfy the Delaunay condition
	// (p1 is inside the circumcircle of [p0, pl, pr]), flip them,
	// then do the same check/flip recursively for the new pair of triangles
	//
	//           pl                    pl
	//          /||\                  /  \
	//       al/ || \bl            al/    \a
	//        /  ||  \              /      \
	//       /  a||b  \    flip    /___ar___\
	//     p0\   ||   /p1   =>   p0\---bl---/p1
	//        \  ||  /              \      /
	//       ar\ || /br             b\    /br
	//          \||/                  \  /
	//           pr                    pr

	b := t.halfedges[a]

	a0 := a - a%3
	b0 := b - b%3

	al := a0 + (a+1)%3
	ar := a0 + (a+2)%3
	bl := b0 + (b+2)%3

	if b < 0 {
		return ar
	}

	p0 := t.triangles[ar]
	pr := t.triangles[a]
	pl := t.triangles[al]
	p1 := t.triangles[bl]

	illegal := inCircle(t.points[p0], t.points[pr], t.points[pl], t.points[p1])

	if illegal {
		t.triangles[a] = p1
		t.triangles[b] = p0

		// edge swapped on the other side of the hull (rare)
		// fix the halfedge reference
		if t.halfedges[bl] == -1 {
			e := t.hull
			for {
				if e.t == bl {
					e.t = a
					break
				}
				e = e.next
				if e == t.hull {
					break
				}
			}
		}

		t.link(a, t.halfedges[bl])
		t.link(b, t.halfedges[ar])
		t.link(ar, bl)

		br := b0 + (b+1)%3

		t.legalize(a)
		return t.legalize(br)
	}

	return ar
}

func (t *triangulator) convexHull() []Point {
	var result []Point
	e := t.hull
	for e != nil {
		result = append(result, t.points[e.i])
		e = e.prev
		if e == t.hull {
			break
		}
	}
	return result
}

func (t *triangulator) area() int {
	c := 0
	e := t.hull
	for e != nil {
		c++
		e = e.prev
		if e == t.hull {
			break
		}
	}

	area := 0.

	j := e.next.i
	for i := 0; i < c; i++ {
		pI := t.points[e.i]
		pJ := t.points[j]
		area += (pJ.X + pI.X) * (pJ.Y - pI.Y)
		j = e.i
		e = e.prev
	}

	return int(math.Round(math.Abs(area / 2.)))
}
