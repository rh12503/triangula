package fitness

import "github.com/RH12503/Triangula/triangulation/incrdelaunay"

// The maximum difference for each pixel
// there can be when compared to the target image. Variance is calculated, so the
// 255 needs to be squared.
const maxPixelDifference = 255 * 255 * 3

// fastRound is an optimized version of math.Round.
func fastRound(n float64) int {
	return int(n+0.5) << 0
}

func createPoint(x, y float64, w, h int) incrdelaunay.Point {
	return incrdelaunay.Point{
		X: int16(fastRound(x * float64(w))),
		Y: int16(fastRound(y * float64(h))),
	}
}

