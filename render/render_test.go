package render

import (
	"Triangula/geom"
	"Triangula/image"
	"Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrianglesOnImage(t *testing.T) {
	data := TrianglesOnImage([]geom.Triangle{
		geom.NewTriangle(12, 32, 65, 43, 23, 87),
	}, image.NewData(100, 100))

	assert.Equal(t, data[0], TriangleData{
		Triangle: normgeom.NewNormTriangle(0.12, 0.32, 0.65, 0.43, 0.23, 0.87),
	})
}
