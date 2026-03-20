package fitness

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/rasterize"
	"github.com/RH12503/Triangula/triangulation/incrdelaunay"
	"math"
)

// trianglesImageFunction is a fitness function that calculates how optimal a point group is when
// triangulated and "placed" onto a target image.
// It does this by first calculating the Delaunay triangulation of the points, then iterating through the
// pixels of the triangles and calculating the variance to the target image. (The lower the variance the better)
type trianglesImageFunction struct {
	target pixelData // pixels data of the target image.

	// Variance data stored in blocks of pixels. The variance of a N*N block can easily be found instead of
	// needing to iterate through N*N pixels.
	targetN   pixelDataN
	blockSize int // The size of each N*N block.

	maxDifference float64 // The maximum difference of all pixels to the target image.

	TriangleCache []CacheData // A cache storing triangles that have already had their variances calculated.

	// The variance calculated for each triangle are put here. This means if the triangles don't change
	// in the next generation, they won't need to be reevaluated.
	nextCache []CacheData

	// The triangulation used to create the triangles.
	Triangulation *incrdelaunay.Delaunay
	// The triangulation of the points before being mutated accessed from the
	// fitness function's base.
	Base *incrdelaunay.Delaunay

	// Delta-based fitness fields.
	// baseDifference and baseArea store the total difference and area from the base evaluation,
	// used to compute mutant fitness incrementally.
	baseDifference float64
	baseArea       float64
	// triangleVariances maps each triangle (by normalized vertex key) to its variance and area.
	// Populated during full evaluations (base/initial).
	triangleVariances map[triangleKey]triangleVarData
	// baseTriangleVariances is a read-only reference to the base's triangleVariances map,
	// shared safely across parallel mutant evaluations.
	baseTriangleVariances map[triangleKey]triangleVarData
	// baseSuperTriangle stores the super triangle vertices for filtering.
	baseSuperTriangle incrdelaunay.Triangle
}

// computeTriangleVariance computes the variance of a single triangle against the target image.
func (t *trianglesImageFunction) computeTriangleVariance(a, b, c incrdelaunay.Point) float64 {
	pixels := t.target.pixels
	pixelsN := t.targetN.pixels

	n := 0
	var sR0, sG0, sB0 int
	var sSq int

	tri := geom.NewTriangle(int(a.X), int(a.Y), int(b.X), int(b.Y), int(c.X), int(c.Y))

	rasterize.DDATriangleBlocks(tri, t.blockSize, func(x0, x1, y int) {
		row := pixels[y]
		if x0 >= 0 && x1 <= len(row) {
			for x := x0; x < x1; x++ {
				pixel := row[x]
				sR0 += int(pixel.r)
				sG0 += int(pixel.g)
				sB0 += int(pixel.b)
				sSq += int(pixel.sq)
			}
		}
		n += x1 - x0
	}, func(x, y int) {
		pixel := pixelsN[y][x]
		sR0 += int(pixel.r)
		sG0 += int(pixel.g)
		sB0 += int(pixel.b)
		sSq += int(pixel.sq)
		n += t.blockSize * t.blockSize
	})

	if n != 0 {
		return float64(sSq) - float64(sR0*sR0+sG0*sG0+sB0*sB0)/float64(n)
	}
	return 0
}

// triangleArea computes the area of a triangle from its vertices.
func triangleArea(a, b, c incrdelaunay.Point) float64 {
	return math.Abs(0.5 * ((float64(b.X-a.X) * float64(c.Y-a.Y)) - (float64(c.X-a.X) * float64(b.Y-a.Y))))
}

// isSuperTriangleVertex returns whether point p is a vertex of the super triangle.
func (t *trianglesImageFunction) isSuperTriangleVertex(p incrdelaunay.Point) bool {
	return p == t.baseSuperTriangle.A || p == t.baseSuperTriangle.B || p == t.baseSuperTriangle.C
}

// touchesSuperTriangle returns whether the triangle has a vertex on the super triangle.
func (t *trianglesImageFunction) touchesSuperTriangle(tri incrdelaunay.Triangle) bool {
	return t.isSuperTriangleVertex(tri.A) || t.isSuperTriangleVertex(tri.B) || t.isSuperTriangleVertex(tri.C)
}

// Calculate returns the fitness of a group of points.
func (t *trianglesImageFunction) Calculate(data PointsData) float64 {
	points := data.Points

	w, h := t.target.Size()

	// Delta evaluation path: when we have a base with variance data and mutations to apply,
	// only recompute affected triangles instead of iterating all of them.
	if t.Base != nil && t.baseTriangleVariances != nil && len(data.Mutations) > 0 {
		if t.Triangulation == nil {
			t.Triangulation = incrdelaunay.NewDelaunay(w, h)
		}
		t.Triangulation.Set(t.Base)
		t.baseSuperTriangle = t.Triangulation.SuperTriangle()

		// Enable change tracking, then apply mutations
		t.Triangulation.BeginTrack()

		for _, m := range data.Mutations {
			t.Triangulation.Remove(createPoint(m.Old.X, m.Old.Y, w, h))
		}
		for _, m := range data.Mutations {
			t.Triangulation.Insert(createPoint(m.New.X, m.New.Y, w, h))
		}

		removed, added := t.Triangulation.EndTrack()

		// During Remove/Insert, intermediate triangles may be created then destroyed.
		// Use a counter-based approach to find net changes:
		// +1 for each appearance in 'added', -1 for each appearance in 'removed'.
		// Net > 0 means truly added, net < 0 means truly removed, net == 0 means cancelled.
		type triInfo struct {
			a, b, c incrdelaunay.Point
			net     int
		}
		netChanges := make(map[triangleKey]*triInfo)

		for _, tri := range added {
			if t.touchesSuperTriangle(tri) {
				continue
			}
			key := makeTriangleKey(tri.A.X, tri.A.Y, tri.B.X, tri.B.Y, tri.C.X, tri.C.Y)
			if info, ok := netChanges[key]; ok {
				info.net++
			} else {
				netChanges[key] = &triInfo{a: tri.A, b: tri.B, c: tri.C, net: 1}
			}
		}

		for _, tri := range removed {
			if t.touchesSuperTriangle(tri) {
				continue
			}
			key := makeTriangleKey(tri.A.X, tri.A.Y, tri.B.X, tri.B.Y, tri.C.X, tri.C.Y)
			if info, ok := netChanges[key]; ok {
				info.net--
			} else {
				netChanges[key] = &triInfo{a: tri.A, b: tri.B, c: tri.C, net: -1}
			}
		}

		difference := t.baseDifference
		area := t.baseArea

		// Rebuild triangleVariances from base
		t.triangleVariances = make(map[triangleKey]triangleVarData, len(t.baseTriangleVariances))
		for k, v := range t.baseTriangleVariances {
			t.triangleVariances[k] = v
		}

		for key, info := range netChanges {
			if info.net == 0 {
				continue // cancelled out
			}
			if info.net < 0 {
				// Truly removed: subtract from base totals
				if vd, ok := t.baseTriangleVariances[key]; ok {
					difference -= vd.variance
					area -= vd.area
					delete(t.triangleVariances, key)
				}
			} else {
				// Truly added: compute variance and add
				variance := t.computeTriangleVariance(info.a, info.b, info.c)
				triArea := triangleArea(info.a, info.b, info.c)
				difference += variance
				area += triArea
				t.triangleVariances[key] = triangleVarData{variance: variance, area: triArea}
			}
		}

		t.baseDifference = difference
		t.baseArea = area

		// Prepare for next generation
		t.Base = nil
		t.baseTriangleVariances = nil

		// Also update the old-style cache for compatibility with the parallel evaluator
		t.nextCache = t.nextCache[:0]
		t.Triangulation.IterTriangles(func(triangle incrdelaunay.Triangle) {
			a := triangle.A
			b := triangle.B
			c := triangle.C
			triData := &TriangleCacheData{
				aX: a.X, aY: a.Y,
				bX: b.X, bY: b.Y,
				cX: c.X, cY: c.Y,
			}
			key := makeTriangleKey(a.X, a.Y, b.X, b.Y, c.X, c.Y)
			if vd, ok := t.triangleVariances[key]; ok {
				triData.fitness = vd.variance
			}
			hash := triData.Hash()
			cacheMask := uint64(len(t.TriangleCache)) - 1
			index0 := uint32(hash & cacheMask)
			triData.SetCachedHash(index0)
			t.nextCache = append(t.nextCache, triData)
		})
		t.TriangleCache = t.nextCache

		blank := float64(w*h) - area
		difference += maxPixelDifference * blank

		return 1 - (difference / t.maxDifference)
	}

	// Full evaluation path: when no base (initial eval, or base evaluation itself).
	if t.Triangulation == nil {
		t.Triangulation = incrdelaunay.NewDelaunay(w, h)
		for _, p := range points {
			t.Triangulation.Insert(createPoint(p.X, p.Y, w, h))
		}
		t.baseSuperTriangle = t.Triangulation.SuperTriangle()
	} else if t.Base != nil {
		// Has a base triangulation but no baseTriangleVariances (or no mutations) - fall back to full eval
		t.Triangulation.Set(t.Base)
		t.baseSuperTriangle = t.Triangulation.SuperTriangle()

		for _, m := range data.Mutations {
			t.Triangulation.Remove(createPoint(m.Old.X, m.Old.Y, w, h))
		}
		for _, m := range data.Mutations {
			t.Triangulation.Insert(createPoint(m.New.X, m.New.Y, w, h))
		}
	}

	// Prepare for next generation
	t.Base = nil

	t.nextCache = t.nextCache[:0]

	// Calculate the variance between the target image and current triangles

	var difference float64

	cacheMask := uint64(len(t.TriangleCache)) - 1

	tris := t.TriangleCache

	area := 0.

	// Initialize triangleVariances map for delta evaluation in future generations
	if t.triangleVariances == nil {
		t.triangleVariances = make(map[triangleKey]triangleVarData)
	} else {
		// Clear the map for reuse
		for k := range t.triangleVariances {
			delete(t.triangleVariances, k)
		}
	}

	t.Triangulation.IterTriangles(func(triangle incrdelaunay.Triangle) {
		a := triangle.A
		b := triangle.B
		c := triangle.C

		// The total area is taken into account when calculating the fitness
		triArea := triangleArea(a, b, c)
		area += triArea

		triData := &TriangleCacheData{
			aX: a.X,
			aY: a.Y,
			bX: b.X,
			bY: b.Y,
			cX: c.X,
			cY: c.Y,
		}

		hash := triData.Hash()

		index0 := uint32(hash & cacheMask)

		data := tris[index0]

		var diff float64

		// Check if the triangle is in the cache
		if data == nil || !data.Equals(triData) {
			// The triangle isn't in the hash, so calculate the variance
			diff = t.computeTriangleVariance(a, b, c)
		} else {
			// If the triangle is in the cache, we don't need to recalculate the variance
			diff = data.Data()
		}

		difference += diff
		triData.fitness = diff
		triData.SetCachedHash(index0)
		t.nextCache = append(t.nextCache, triData)

		// Store in triangleVariances map for delta evaluation
		key := makeTriangleKey(a.X, a.Y, b.X, b.Y, c.X, c.Y)
		t.triangleVariances[key] = triangleVarData{variance: diff, area: triArea}
	})

	t.TriangleCache = t.nextCache

	// Store base values for delta evaluation
	t.baseDifference = difference
	t.baseArea = area

	// Lower the fitness based on how many blank pixels there are (the smaller the area)
	// (As the triangles should cover the entire image)
	blank := float64(w*h) - area

	difference += maxPixelDifference * blank

	return 1 - (difference / t.maxDifference)
}

func (t *trianglesImageFunction) SetBase(other CacheFunction) {
	base := other.(*trianglesImageFunction)
	t.Base = base.Triangulation
	t.baseDifference = base.baseDifference
	t.baseArea = base.baseArea
	t.baseTriangleVariances = base.triangleVariances // shared read-only reference
}

func (t *trianglesImageFunction) Cache() []CacheData {
	return t.TriangleCache
}

func (t *trianglesImageFunction) SetCache(cache []CacheData) {
	t.TriangleCache = cache
}

// TrianglesImageFunctions returns an array of fitness functions.
func TrianglesImageFunctions(target image.Data, blockSize, n int) []CacheFunction {
	w, h := target.Size()

	functions := make([]CacheFunction, n)
	pixels := fromImage(target)
	pixelsN := fromImageN(target, blockSize)

	maxDiff := float64(maxPixelDifference * w * h)

	for i := 0; i < n; i++ {
		function := trianglesImageFunction{
			target:        pixels,
			targetN:       pixelsN,
			blockSize:     blockSize,
			maxDifference: maxDiff,
			TriangleCache: make([]CacheData, 2),
		}
		functions[i] = &function
	}

	return functions
}

// NewTrianglesImageFunction returns a new fitness function.
func NewTrianglesImageFunction(target image.Data, blockSize int) CacheFunction {
	w, h := target.Size()

	return &trianglesImageFunction{
		target:        fromImage(target),
		targetN:       fromImageN(target, blockSize),
		blockSize:     blockSize,
		maxDifference: float64(maxPixelDifference * w * h),
		TriangleCache: make([]CacheData, 2),
	}
}
