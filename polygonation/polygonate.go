package polygonation

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/triangulation/incrdelaunay"
	"math"
)

func Polygonate(points normgeom.NormPointGroup, w, h int) []geom.Polygon {
	fW, fH := float64(w), float64(h)

	triangulation := incrdelaunay.NewDelaunay(w, h)
	for _, p := range points {
		triangulation.Insert(incrdelaunay.Point{
			X: int16(math.Round(p.X * fW)),
			Y: int16(math.Round(p.Y * fH)),
		})
	}

	var polygons []geom.Polygon

	incrdelaunay.Voronoi(triangulation, func(points []incrdelaunay.FloatPoint) {
		var polygon geom.Polygon

		for _, p := range points {
			polygon.Points = append(polygon.Points, geom.Point{
				X: int(math.Round(p.X)),
				Y: int(math.Round(p.Y)),
			})
		}

		polygons = append(polygons, polygon)
	}, w, h)

	return polygons
}
