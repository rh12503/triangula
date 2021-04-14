package incrdelaunay

// pointMap is an optimized hashset like structure to indicate if and how many copies of a point exists.
type pointMap struct {
	points [][]pointEntry // The table.
}

// newPointMap creates a pointMap with a specified size.
func newPointMap(size int) pointMap {
	pMap := pointMap{}
	pMap.points = make([][]pointEntry, size)
	return pMap
}

// AddPoint adds a point to the pointMap, accounting for duplicates.
func (pm *pointMap) AddPoint(point Point) uint16 {
	x, y := uint16(point.X), uint16(point.Y)
	index := point.Hash() % len(pm.points)

	// Check if the point already exists
	for i, p := range pm.points[index] {
		if x == p.x && y == p.y {
			pm.points[index][i].count++
			return pm.points[index][i].count
		}
	}

	// Point doesn't exist; Create a new entry
	pm.points[index] = append(pm.points[index], pointEntry{
		x:     x,
		y:     y,
		count: 1,
	})
	return 1
}

// RemovePoint removes a point from the pointMap, accounting for duplicates,
// and returning the number of the same point left.
func (pm *pointMap) RemovePoint(point Point) uint16 {
	x, y := uint16(point.X), uint16(point.Y)
	index := point.Hash() % len(pm.points)

	for i, p := range pm.points[index] {
		// Find the point and decrement its count
		if x == p.x && y == p.y {
			pm.points[index][i].count--

			// If there are no copies of the point left, remove its entry
			if pm.points[index][i].count != 0 {
				return pm.points[index][i].count
			}

			pm.points[index][i] = pm.points[index][len(pm.points[index])-1]
			pm.points[index] = pm.points[index][:len(pm.points[index])-1]
			return 0
		}
	}
	panic("point doesn't exist")
}

// NumPoints returns the number of points in the pointMap.
func (pm *pointMap) NumPoints() int {
	total := 0
	for _, p := range pm.points {
		for _, e := range p {
			total += int(e.count)
		}
	}
	return total
}

// Set sets the pointMap to another pointMap.
func (pm *pointMap) Set(other *pointMap) {
	for i := range pm.points {
		pm.points[i] = pm.points[i][:cap(pm.points[i])]
		if len(pm.points[i]) > len(other.points[i]) {
			pm.points[i] = pm.points[i][:len(other.points[i])]
		} else if len(pm.points[i]) < len(other.points[i]) {
			pm.points[i] = make([]pointEntry, len(other.points[i]))
		}

		copy(pm.points[i], other.points[i])
	}
}

// pointEntry is used to keep track of a point and how many copies of a point there are.
type pointEntry struct {
	x, y  uint16
	count uint16 // the number of a point the hash table contains.
}
