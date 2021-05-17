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
	Base          *incrdelaunay.Delaunay
}

// Calculate returns the fitness of a group of points.
func (t *trianglesImageFunction) Calculate(data PointsData) float64 {
	points := data.Points

	w, h := t.target.Size()

	if t.Triangulation == nil {
		// If there's no base triangulation, the whole triangulation needs to be recalculated
		t.Triangulation = incrdelaunay.NewDelaunay(w, h)
		for _, p := range points {
			t.Triangulation.Insert(createPoint(p.X, p.Y, w, h))
		}
	} else if t.Base != nil {
		// If there is a base triangulation, set this triangulation to the base
		t.Triangulation.Set(t.Base)

		// And then modify the points that have been mutated
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

	pixels := t.target.pixels
	pixelsN := t.targetN.pixels

	// Calcuate the variance between the target image and current triangles

	var difference float64

	cacheMask := uint64(len(t.TriangleCache)) - 1

	tris := t.TriangleCache

	area := 0.

	t.Triangulation.IterTriangles(func(triangle incrdelaunay.Triangle) {
		a := triangle.A
		b := triangle.B
		c := triangle.C

		// The total area is taken into account when calculating the fitness
		area += math.Abs(0.5 * ((float64(b.X-a.X) * float64(c.Y-a.Y)) - (float64(c.X-a.X) * float64(b.Y-a.Y))))

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

		// Check if the triangle is in the cache
		if data == nil || !data.Equals(triData) {
			// The triangle isn't in the hash, so calculate the variance
			// Welford's online algorithm is used:
			// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
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
			var diff float64
			if n != 0 {
				diff = float64(sSq) - float64(sR0*sR0+sG0*sG0+sB0*sB0)/float64(n)
			}
			difference += diff
			triData.fitness = diff
			triData.SetCachedHash(index0)
			t.nextCache = append(t.nextCache, triData)
		} else {
			// If the triangle is in the cache, we don't need to recalculate the variance
			difference += data.Data()
			t.nextCache = append(t.nextCache, data)
		}
	})

	t.TriangleCache = t.nextCache

	// Lower the fitness based on how many blank pixels there are (the smaller the area)
	// (As the triangles should cover the entire image)
	blank := float64(w*h) - area

	difference += maxPixelDifference * blank

	return 1 - (difference / t.maxDifference)
}

func (t *trianglesImageFunction) SetBase(other CacheFunction) {
	t.Base = other.(*trianglesImageFunction).Triangulation
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
