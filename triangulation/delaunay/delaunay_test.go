package delaunay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriangulate(t *testing.T) {
	delaunay, _ := Triangulate([]Point{
		{12, 32},
		{56, 65},
		{98, 37},
		{23, 12},
		{18, 21},
	})

	assert.Equal(t, delaunay.Triangles, []int{1, 2, 3, 3, 4, 1, 4, 0, 1})
}
