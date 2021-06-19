// Package incrdelaunay implements a library for incremental Delaunay triangulation, with support for
// dynamically adding and removing points.
package incrdelaunay

import (
	"math"
	"sort"
)

// Delaunay represents a Delaunay triangulation.
type Delaunay struct {
	triangles []Triangle

	grid CircumcircleGrid // For fast detection of circumcircles containing a point.

	pointMap pointMap

	freeTriangles []uint16 // A list of free indexes in the triangles slice.

	superTriangle Triangle // A triangle that contains all points added.

	edges []Edge // For performance purposes.

	hull []Point // For performance purposes.

	ears []ear // For performance purposes.

	numPoints    int // The number of points in the triangulation (including duplicate points).
	uniquePoints int
}

// NewDelaunay returns a new Delaunay triangulation.
func NewDelaunay(w, h int) *Delaunay {
	delaunay := Delaunay{}

	superTriangle := NewSuperTriangle(w, h)

	delaunay.grid = NewCircumcircleGrid(10, 10, w, h)

	delaunay.addTriangle(superTriangle)
	delaunay.superTriangle = superTriangle

	delaunay.pointMap = newPointMap(100)

	return &delaunay
}

// Insert adds a point to the Delaunay triangulation using the Bowyer-Watson algorithm.
// Duplicate points are kept track of.
func (d *Delaunay) Insert(p Point) bool {
	if d.pointMap.AddPoint(p) == 1 {
		d.uniquePoints++
	}

	d.numPoints++

	d.resetEdges()

	// Iterate through all the triangles with circumcircles that contain the point
	d.grid.RemoveCircumcirclesThatContain(p, d.triangles, func(i uint16) {
		t := d.triangles[i]

		d.addEdge(NewEdge(t.A, t.B))
		d.addEdge(NewEdge(t.B, t.C))
		d.addEdge(NewEdge(t.C, t.A))

		d.markFreeTriangle(i) // Remove the triangle
	})

	// Connect the vertices along the hole to the point
	for _, e := range d.edges {
		A := int64(e.B.X - e.A.X)
		B := int64(e.B.Y - e.A.Y)
		G := A*int64(p.Y-e.B.Y) - B*int64(p.X-e.B.X)
		if G != 0 {
			new := NewTriangle(e.A, e.B, p)
			d.addTriangle(new)
		}
	}

	return true
}

// Remove removes a point from the Delaunay Triangulation.
// If there are duplicates of the point only one of the duplicates is removed.
func (d *Delaunay) Remove(p Point) {
	d.numPoints--
	if d.pointMap.RemovePoint(p) != 0 {
		return
	}
	d.uniquePoints--
	d.hull = d.hull[:0]

	// Adds a point to the hull
	addPoint := func(po Point) {
		if po == p {
			return
		}
		found := false
		for _, v := range d.hull {
			if po == v {
				found = true
				break
			}
		}

		if !found {
			d.hull = append(d.hull, po)
		}
	}

	// Iterates through all the triangles that have a connection to point p
	d.grid.RemoveThatHasVertex(p, d.triangles, func(i uint16) {
		t := d.triangles[i]

		addPoint(t.A)
		addPoint(t.B)
		addPoint(t.C)

		d.markFreeTriangle(i)
	})

	// Sort the points in the hull counterclockwise
	sort.Slice(d.hull, func(i, j int) bool {
		a := d.hull[j]
		b := d.hull[i]

		if a.X-p.X >= 0 && b.X-p.X < 0 {
			return true
		}
		if a.X-p.X < 0 && b.X-p.X >= 0 {
			return false
		}
		if a.X-p.X == 0 && b.X-p.X == 0 {
			if a.Y-p.Y >= 0 || b.Y-p.Y >= 0 {
				return a.Y > b.Y
			}
			return b.Y > a.Y
		}

		det := int(a.X-p.X)*int(b.Y-p.Y) - int(b.X-p.X)*int(a.Y-p.Y)
		if det < 0 {
			return true
		}
		if det > 0 {
			return false
		}

		panic("...")
	})

	// Add the ears one by one based on its score
	// using the algorithm described in: https://hal.inria.fr/inria-00167201/document.
	// An ear is a triangle made by three consecutive points along the hull
	d.ears = d.ears[:0]

	// Create the ears of the hull and calculate their scores
	for i := 0; i < len(d.hull); i++ {
		e := ear{a: d.hull[i], b: d.hull[(i+1)%len(d.hull)], c: d.hull[(i+2)%len(d.hull)]}
		e.computeScore(p)
		d.ears = append(d.ears, e)
	}

	for len(d.ears) > 3 {
		// Find the ear with the lowest score
		lowestScore := math.MaxFloat64
		ear := d.ears[0]
		index := 0

		for i, e := range d.ears {
			if e.score < lowestScore {
				lowestScore = e.score
				ear = e
				index = i
			}
		}

		// Add the ear to the triangulation
		d.addTriangle(NewTriangle(ear.a, ear.b, ear.c))

		// Remove the ear from the list of ears, remove it from the hull, and update all other ears
		// accordingly

		// Find the ears before and after the current ear
		before := (index - 1) % len(d.ears)
		after := (index + 1) % len(d.ears)
		if before == -1 {
			before = len(d.ears) - 1
		}

		if after == -1 {
			after = len(d.ears) - 1
		}

		// Connect those ears together as the ear between them has been removed.
		// Then, update the scores of the two modified ears
		d.ears[before].c = d.ears[index].c
		d.ears[before].computeScore(p)
		d.ears[after].a = d.ears[index].a
		d.ears[after].computeScore(p)
		d.ears = append(d.ears[:index], d.ears[index+1:]...)
	}

	// All three remaining ears are identical
	d.addTriangle(NewTriangle(d.ears[0].a, d.ears[0].b, d.ears[0].c))
}

// addEdge adds an edge to the edges if edge e is unique.
func (d *Delaunay) addEdge(e Edge) {
	found := false
	for i, edge := range d.edges {
		if e == edge {
			found = true
			d.edges[i] = d.edges[len(d.edges)-1]
			d.edges = d.edges[:len(d.edges)-1]
			break
		}
	}
	if !found {
		d.edges = append(d.edges, e)
	}
}

// resetEdges empties edges.
func (d *Delaunay) resetEdges() {
	d.edges = d.edges[:0]
}

// addToHull adds a point to hull.
func (d *Delaunay) addToHull(p Point) {
	d.hull = append(d.hull, p)
}

// addTriangle adds a triangle to triangulation.
func (d *Delaunay) addTriangle(t Triangle) {
	// First check if there are any free indexes available
	if len(d.freeTriangles) > 0 {
		// If there is an index available, add the triangle to the index
		index := d.freeTriangles[len(d.freeTriangles)-1]
		d.triangles[index] = t
		d.grid.AddTriangle(t, index) // Update the grid
		d.freeTriangles = d.freeTriangles[:len(d.freeTriangles)-1]
		return
	}

	d.grid.AddTriangle(t, uint16(len(d.triangles))) // Update the grid

	d.triangles = append(d.triangles, t)
}

// markFreeTriangle marks a triangle in the triangles slice as no longer being needed,
// essentially removing the triangle from the triangulation.
func (d *Delaunay) markFreeTriangle(i uint16) {
	d.triangles[i].A.X = -1 // Marks the triangle as invalid
	d.freeTriangles = append(d.freeTriangles, i)
}

// IterTriangles iterates through all the triangles in the triangulation, calling function triangle for each one.
func (d Delaunay) IterTriangles(triangle func(t Triangle)) {
	for _, t := range d.triangles {
		if t.A.X == -1 {
			continue
		}

		// Exclude triangles that are connected to the superTriangle
		if t.A != d.superTriangle.A && t.B != d.superTriangle.A && t.C != d.superTriangle.A &&
			t.A != d.superTriangle.B && t.B != d.superTriangle.B && t.C != d.superTriangle.B &&
			t.A != d.superTriangle.C && t.B != d.superTriangle.C && t.C != d.superTriangle.C {
			triangle(t)
		}
	}
}

// NumPoints returns the number of points in the triangulation, including duplicate points.
func (d Delaunay) NumPoints() int {
	return d.numPoints
}

// Set sets the triangulation to another triangulation.
func (d *Delaunay) Set(other *Delaunay) {
	d.triangles = d.triangles[:cap(d.triangles)]

	if len(d.triangles) > len(other.triangles) {
		d.triangles = d.triangles[:len(other.triangles)]
	} else if len(d.triangles) < len(other.triangles) {
		d.triangles = make([]Triangle, len(other.triangles))
	}

	copy(d.triangles, other.triangles)

	d.freeTriangles = d.freeTriangles[:cap(d.freeTriangles)]
	if len(d.freeTriangles) > len(other.freeTriangles) {
		d.freeTriangles = d.freeTriangles[:len(other.freeTriangles)]
	} else if len(d.freeTriangles) < len(other.freeTriangles) {
		d.freeTriangles = make([]uint16, len(other.freeTriangles))
	}

	copy(d.freeTriangles, other.freeTriangles)

	d.superTriangle = other.superTriangle

	d.grid.Set(&other.grid)
	d.pointMap.Set(&other.pointMap)
}

// HasPoint returns whether the triangulation contains point p.
func (d Delaunay) HasPoint(p Point) bool {
	return d.grid.HasPoint(p, d.triangles)
}

// GetClosestTo returns the closest point in the triangulation to point p.
func (d Delaunay) GetClosestTo(p Point) Point {
	var closest Point
	closestDist := int64(math.MaxInt64)

	for _, t := range d.triangles {
		if t.A.X == -1 {
			continue
		}

		// Check all three of the triangles vertices to see if one is closer

		dist := p.DistSq(t.A)
		if dist < closestDist {
			closestDist = dist
			closest = t.A
		}

		dist = p.DistSq(t.B)
		if dist < closestDist {
			closestDist = dist
			closest = t.B
		}

		dist = p.DistSq(t.C)
		if dist < closestDist {
			closestDist = dist
			closest = t.C
		}
	}
	return closest
}
