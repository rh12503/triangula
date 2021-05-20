package incrdelaunay

type polygonMap struct {
	polygons [][]polygonEntry // The table.
}

// newPointMap creates a pointMap with a specified size.
func newPolygonMap(size int) polygonMap {
	pMap := polygonMap{}
	pMap.polygons = make([][]polygonEntry, size)
	return pMap
}

// AddPoint adds a point to the pointMap, accounting for duplicates.
func (pm *polygonMap) AddPolygon(point Point, polygon []FloatPoint) {
	index := point.Hash() % len(pm.polygons)

	pm.polygons[index] = append(pm.polygons[index], polygonEntry{
		point:   point,
		polygon: polygon,
	})
}

func (pm *polygonMap) RemovePolygon(point Point) bool {
	index := point.Hash() % len(pm.polygons)

	for i, p := range pm.polygons[index] {
		if p.point == point {
			pm.polygons[index][i] = pm.polygons[index][len(pm.polygons[index])-1]
			pm.polygons[index] = pm.polygons[index][:len(pm.polygons[index])-1]
			return true
		}
	}

	panic("polygon doesn't exist")
}

// Set sets the pointMap to another pointMap.
func (pm *polygonMap) Set(other *polygonMap) {
	for i := range pm.polygons {
		pm.polygons[i] = pm.polygons[i][:cap(pm.polygons[i])]
		if len(pm.polygons[i]) > len(other.polygons[i]) {
			pm.polygons[i] = pm.polygons[i][:len(other.polygons[i])]
		} else if len(pm.polygons[i]) < len(other.polygons[i]) {
			pm.polygons[i] = make([]polygonEntry, len(other.polygons[i]))
		}

		copy(pm.polygons[i], other.polygons[i])
	}
}

func (pm polygonMap) IterPolygons(polygon func([]FloatPoint)) {
	for i := range pm.polygons {
		for _, p := range pm.polygons[i] {
			polygon(p.polygon)
		}
	}
}

// pointEntry is used to keep track of a point and how many copies of a point there are.
type polygonEntry struct {
	point Point
	polygon []FloatPoint
}
