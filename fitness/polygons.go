package fitness

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/rasterize"
	"github.com/RH12503/Triangula/triangulation/incrdelaunay"
)

type polygonsImageFunction struct {
	target pixelData // pixels data of the target image.

	// Variance data stored in blocks of pixels. The variance of a N*N block can easily be found instead of
	// needing to iterate through N*N pixels.
	targetN   pixelDataN
	blockSize int // The size of each N*N block.

	maxDifference float64 // The maximum difference of all pixels to the target image.

	TriangleCache []CacheData
	nextCache     []CacheData

	// The triangulation used to create the triangles.
	Triangulation *incrdelaunay.IVoronoi
	// The triangulation of the points before being mutated accessed from the
	// fitness function's base.
	Base *incrdelaunay.IVoronoi
}

// Calculate returns the fitness of a group of points.
func (g *polygonsImageFunction) Calculate(data PointsData) float64 {
	points := data.Points

	w, h := g.target.Size()

	if g.Triangulation == nil {
		// If there's no base triangulation, the whole triangulation needs to be recalculated
		g.Triangulation = incrdelaunay.NewVoronoi(w, h)
		for _, p := range points {
			g.Triangulation.Insert(createPoint(p.X, p.Y, w, h))
		}
	} else if g.Base != nil {
		// If there is a base triangulation, set this triangulation to the base
		g.Triangulation.Set(g.Base)

		// And then modify the points that have been mutated
		for _, m := range data.Mutations {
			g.Triangulation.Remove(createPoint(m.Old.X, m.Old.Y, w, h))
		}

		for _, m := range data.Mutations {
			g.Triangulation.Insert(createPoint(m.New.X, m.New.Y, w, h))
		}
	}

	// Prepare for next generation
	g.Base = nil

	g.nextCache = g.nextCache[:0]

	pixels := g.target.pixels
	pixelsN := g.targetN.pixels

	// Calcuate the variance between the target image and current triangles

	var difference float64

	cacheMask := uint64(len(g.TriangleCache)) - 1

	tris := g.TriangleCache

	var polygon geom.Polygon

	var polygonData []int16

	g.Triangulation.IterPolygons(func(points []incrdelaunay.FloatPoint) {

		polygon.Points = polygon.Points[:0]
		polygonData = polygonData[:0]


		for _, p := range points {
			polygon.Points = append(polygon.Points, geom.Point{
				X: fastRound(p.X),
				Y: fastRound(p.Y),
			})
			polygonData = append(polygonData, int16(fastRound(p.X)))
			polygonData = append(polygonData, int16(fastRound(p.Y)))
		}

		polyData := &PolygonCacheData{
			coords:  polygonData,
		}

		hash := polyData.Hash()

		index0 := uint32(hash & cacheMask)

		data := tris[index0]

		// Check if the triangle is in the cache
		if data == nil || !data.Equals(polyData) {
			// The triangle isn't in the hash, so calculate the variance
			// Welford's online algorithm is used:
			// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
			n := 0
			var sR0, sG0, sB0 int
			var sSq int

			rasterize.DDAPolygon(polygon, g.blockSize, func(x0, x1, y int) {
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
				n += g.blockSize * g.blockSize
			})
			var diff float64
			if n != 0 {
				diff = float64(sSq) - float64(sR0*sR0+sG0*sG0+sB0*sB0)/float64(n)
			}
			difference += diff
			polyData.fitness = diff
			polyData.SetCachedHash(index0)
			var newPolyData []int16
			newPolyData = append(newPolyData, polygonData...)

			polyData.coords = newPolyData

			g.nextCache = append(g.nextCache, polyData)
		} else {
			// If the triangle is in the cache, we don't need to recalculate the variance
			difference += data.Data()
			g.nextCache = append(g.nextCache, data)
		}
	})

	g.TriangleCache = g.nextCache

	return 1 - (difference / g.maxDifference)
}

func (g *polygonsImageFunction) SetBase(other CacheFunction) {
	g.Base = other.(*polygonsImageFunction).Triangulation
}

func (g *polygonsImageFunction) Cache() []CacheData {
	return g.TriangleCache
}

func (g *polygonsImageFunction) SetCache(cache []CacheData) {
	g.TriangleCache = cache
}

func PolygonsImageFunctions(target image.Data, blockSize, n int) []CacheFunction {
	w, h := target.Size()

	functions := make([]CacheFunction, n)
	pixels := fromImage(target)
	pixelsN := fromImageN(target, blockSize)

	maxDiff := float64(maxPixelDifference * w * h)

	for i := 0; i < n; i++ {
		function := polygonsImageFunction{
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
