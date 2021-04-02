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
}
