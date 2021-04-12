package delaunay

type Triangulation struct {
	Points       []Point
	Area         int
	Triangles    []int
	Halfedges    []int
	Triangulator *triangulator
}

// Triangulate returns a Delaunay triangulation of the provided generator.
func Triangulate(points []Point) (*Triangulation, error) {
	t := newTriangulator(points)
	err := t.update()

	return &Triangulation{points, t.area(), t.triangles, t.halfedges, t}, err
}

func (t *Triangulation) Update(points []Point) error {
	t.Triangulator.points = points
	err := t.Triangulator.update()
	t.Points = points
	t.Area = t.Triangulator.area()
	t.Triangles = t.Triangulator.triangles
	t.Halfedges = t.Triangulator.halfedges
	return err
}

func (t *Triangulation) area() float64 {
	var result float64
	points := t.Points
	ts := t.Triangles
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		result += area(p0, p1, p2)
	}
	return result / 2
}

func (t Triangulation) ConvexHull() []Point {
	return t.Triangulator.convexHull()
}
