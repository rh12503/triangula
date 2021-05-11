package render

import (
	"github.com/RH12503/Triangula/color"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
)

type PolygonData struct {
	Triangle normgeom.NormPolygon
	Color    color.RGB
}

func PolygonsOnImage(polygons []geom.Polygon, image image.Data) []PolygonData {
	polygonData := make([]PolygonData, len(polygons))

	/*w, h := image.Size()

	for i, p := range polygons {
		var color color.AverageRGB

		rasterize.DDATriangle(t, func(x, y int) {
			color.Add(image.RGBAt(x, y))
		})

		if color.Count() == 0 {
			for _, p := range t.Points {
				x, y := min(p.X, w-1), min(p.Y, h-1)

				color.Add(image.RGBAt(x, y))
			}
		}

		data := TriangleData{
			Triangle: t.ToNorm(w, h),
			Color:    color.Average(),
		}
		triangleData[i] = data
	}*/

	return polygonData
}

