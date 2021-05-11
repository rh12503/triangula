package polygonation

import (
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/polygonation/voronoi"
	"math"
)

func Polygonate(points normgeom.NormPointGroup, w, h int) []geom.Polygon {
	fW, fH := float64(w), float64(h)

	sites := make([]voronoi.Vertex, len(points))

	for i := range sites {
		p := points[i]
		sites[i] = voronoi.Vertex{
			X: math.Round(p.X * fW),
			Y: math.Round(p.X * fH),
		}
	}

	bounds := voronoi.NewBBox(0, fW, 0, fH)

	diagram := voronoi.ComputeDiagram()
}
