package fitness

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/rasterize"
	"github.com/RH12503/Triangula/triangulation/incrdelaunay"
	"math"
)

// The maximum difference for each of the RGB values
// there can be compared to the target image. Variance is calculated, so
// 255 needs to be squared.
const maxPixelDifference = 255 * 255 * 3

// TrianglesImageEvaluator is a fitness function for calculating how optimal a point group is when
// triangulated and "placed" onto a target image.
// It does this by first calculating the Delaunay triangulation of the generator, then calculating the
// optimal color of each triangle. Finally, it iterates through the pixels of the triangles and calculates
// the variance to the target image. (The lower the variance the better)
type TrianglesImageEvaluator struct {
	target pixelData // variance data relating to the pixels of the target image

	// Variance data stored in blocks of pixels. The variance of a N*N block can easily be found instead of
	// needing to iterate through N*N pixels
	targetN   pixelDataN
	blockSize int // The size of each block

	maxDifference float64 // The maximum difference of all pixels compared to the target image

	TriangleCache []TriFit // A cache to store triangles that have already had their variances calculated

	// The variance calculated for each triangle are put here. This means if the triangles don't change
	// in the next generation, they won't need to be reevaluated.
	NextCache []TriFit

	added, removed []normgeom.NormPoint // Lists storing which generator have been modified

	Triangulation *incrdelaunay.Delaunay
	Base          *incrdelaunay.Delaunay // The triangulation which the generator used last generation
}

func (t *TrianglesImageEvaluator) Calculate(points normgeom.NormPointGroup) float64 {
	w, h := t.target.Size()

	// Needs to be cleaned up.
	if t.Triangulation == nil {
		// If there's no base triangulation, the whole triangulation needs to be recalculated
		t.Triangulation = incrdelaunay.NewDelaunay(w, h)
		for _, p := range points {
			t.Triangulation.Insert(incrdelaunay.Point{
				X: int16(fastRound(p.X * float64(w))),
				Y: int16(fastRound(p.Y * float64(h))),
			})
		}
	} else if t.Base != nil {
		// If there is a base triangulation, set this triangulation to the base
		t.Triangulation.Set(t.Base)

		// And then modify the generator that have been mutated
		for _, p := range t.removed {
			t.Triangulation.Remove(incrdelaunay.Point{
				X: int16(fastRound(p.X * float64(w))),
				Y: int16(fastRound(p.Y * float64(h))),
			})
		}

		for _, p := range t.added {
			t.Triangulation.Insert(incrdelaunay.Point{
				X: int16(fastRound(p.X * float64(w))),
				Y: int16(fastRound(p.Y * float64(h))),
			})
		}
	}

	// Prepare for next generation
	t.Base = nil

	t.removed = t.removed[:0]
	t.added = t.added[:0]

	t.NextCache = t.NextCache[:0]

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

		// The total area is taken into account when evaluating the fitness
		area += math.Abs(0.5 * ((float64(b.X-a.X) * float64(c.Y-a.Y)) - (float64(c.X-a.X) * float64(b.Y-a.Y))))

		triData := TriFit{
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
		if !data.Equals(triData) {
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
			triData.OtherHash = index0
			t.NextCache = append(t.NextCache, triData)
		} else {
			// If the triangle is in the cache, we don't need to recalculate the variance
			difference += data.fitness
			t.NextCache = append(t.NextCache, data)
		}
	})

	t.TriangleCache = t.NextCache

	// Lower the fitness based on how many blank pixels there are (the smaller the area)
	// (As the triangles should cover the entire image)
	blank := float64(w*h) - area

	difference += maxPixelDifference * blank

	return 1 - (difference / t.maxDifference)
}

// PointMoved is called to indicate which generator have been modified
func (t *TrianglesImageEvaluator) PointMoved(old, new normgeom.NormPoint) {
	t.removed = append(t.removed, old)
	t.added = append(t.added, new)
}

// TrianglesImageEvaluators returns an array of the fitness evaluator
func TrianglesImageEvaluators(target image.Data, blockSize, n int) []*TrianglesImageEvaluator {
	w, h := target.Size()

	evaluators := make([]*TrianglesImageEvaluator, n)
	pixels := fromImage(target)
	pixelsN := fromImageN(target, blockSize)

	maxDiff := float64(maxPixelDifference * w * h)

	for i := 0; i < n; i++ {
		function := TrianglesImageEvaluator{
			target:        pixels,
			targetN:       pixelsN,
			blockSize:     blockSize,
			maxDifference: maxDiff,
			TriangleCache: make([]TriFit, 2),
		}
		evaluators[i] = &function
	}

	return evaluators
}

// TrianglesImageEvaluators returns a single one of the fitness evaluator
func NewTrianglesImageEvaluator(target image.Data, blockSize int) *TrianglesImageEvaluator {
	w, h := target.Size()

	return &TrianglesImageEvaluator{
		target:        fromImage(target),
		targetN:       fromImageN(target, blockSize),
		blockSize:     blockSize,
		maxDifference: float64(maxPixelDifference * w * h),
		TriangleCache: make([]TriFit, 2),
	}
}

// TriFit stores the triangles vertices and its fitness. The struct is used to cache calculations
type TriFit struct {
	aX, aY    int16
	bX, bY    int16
	cX, cY    int16
	fitness   float64
	OtherHash uint32
}

func (t TriFit) Equals(other TriFit) bool {
	return t.aX == other.aX && t.aY == other.aY &&
		t.bX == other.bX && t.bY == other.bY &&
		t.cX == other.cX && t.cY == other.cY
}

// Hash calculates the hash of a TriFit
func (t TriFit) Hash() uint64 {
	x := int(t.aX) + int(t.bX) + int(t.cX)
	y := int(t.aY) + int(t.bY) + int(t.cY)

	hash := uint64((53+x)*53 + y)

	return hash
}

// fastRound is an optimized version of math.Round
func fastRound(n float64) int {
	return int(n+0.5) << 0
}
