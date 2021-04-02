package normgeom

// A simple struct to represent a triangle with normalized generator
type NormTriangle struct {
	Points [3]NormPoint
}

func NewNormTriangle(x0, y0, x1, y1, x2, y2 float64) NormTriangle {
	return NormTriangle{[3]NormPoint{{x0, y0}, {x1, y1}, {x2, y2}}}
}
