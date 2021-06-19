package render

import (
	"github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/rasterize"
)

type PolygonData struct {
	Polygon normgeom.NormPolygon
	Color   color.RGB
}

func PolygonsOnImage(polygons []geom.Polygon, image image.Data) []PolygonData {
	polygonData := make([]PolygonData, len(polygons))

	w, h := image.Size()

	for i, poly := range polygons {
		var color color.AverageRGB

		for i := 2; i < len(poly.Points); i++ {
			tri := geom.Triangle{Points: [3]geom.Point{
				poly.Points[i],
				poly.Points[i-1],
				poly.Points[0],
			}}

			rasterize.DDATriangle(tri, func(x, y int) {
				color.Add(image.RGBAt(x, y))
			})
		}

		if color.Count() == 0 {
			for _, p := range poly.Points {
				x, y := min(p.X, w-1), min(p.Y, h-1)

				color.Add(image.RGBAt(x, y))
			}
		}

		data := PolygonData{
			Polygon: poly.ToNorm(w, h),
			Color:   color.Average(),
		}
		polygonData[i] = data
	}

	return polygonData
}
