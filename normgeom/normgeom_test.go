package normgeom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNormTriangle(t *testing.T) {
	triangle := NewNormTriangle(0.14, 0.15, 0.92, 0.65, 0.35, 0.89)
	assert.Equal(t, triangle, NormTriangle{Points: [3]NormPoint{
		{0.14, 0.15},
		{0.92, 0.65},
		{0.35, 0.89},
	}})
}

func TestNormPointGroup_Set(t *testing.T) {
	a := NormPointGroup{{0, 0.12}, {0.34, 0.54}, {0.21, 0.45}}
	b := NormPointGroup{{}, {}, {}}
	a.Set(b)

	assert.Equal(t, a, b)
}

func TestNormPointGroup_Copy(t *testing.T) {
	a := NormPointGroup{{0, 0.12}, {0.34, 0.54}, {0.21, 0.45}}
	b := a.Copy()

	assert.Equal(t, a, b)

	// Modify copy and ensure original unchanged
	b[0].X = 0.99
	assert.NotEqual(t, a[0].X, b[0].X)
}

func TestNormPolygon(t *testing.T) {
	poly := NormPolygon{Points: []NormPoint{{0.1, 0.2}, {0.3, 0.4}, {0.5, 0.6}}}
	assert.Equal(t, 3, len(poly.Points))
	assert.Equal(t, 0.1, poly.Points[0].X)
	assert.Equal(t, 0.2, poly.Points[0].Y)
}

func TestNormPointGroup_Set_Independence(t *testing.T) {
	a := NormPointGroup{{0.1, 0.2}, {0.3, 0.4}}
	b := NormPointGroup{{0.5, 0.6}, {0.7, 0.8}}
	a.Set(b)
	assert.Equal(t, a, b)

	// Modify b and ensure a is unchanged
	b[0].X = 0.99
	assert.NotEqual(t, a[0].X, b[0].X)
}
