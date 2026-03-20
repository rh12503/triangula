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

func TestNormPointGroup_Copy_IsIndependent(t *testing.T) {
	a := NormPointGroup{{0.1, 0.2}, {0.3, 0.4}}
	b := a.Copy()
	b[0] = NormPoint{X: 0.9, Y: 0.9}
	// Modifying the copy must not affect the original
	assert.Equal(t, a[0], NormPoint{X: 0.1, Y: 0.2})
}

func TestNormPoint_Constrain_ClampsAbove1(t *testing.T) {
	p := NormPoint{X: 1.5, Y: 2.0}
	p.Constrain()
	assert.Equal(t, p, NormPoint{X: 1.0, Y: 1.0})
}

func TestNormPoint_Constrain_ClampsBelow0(t *testing.T) {
	p := NormPoint{X: -0.5, Y: -1.0}
	p.Constrain()
	assert.Equal(t, p, NormPoint{X: 0.0, Y: 0.0})
}

func TestNormPoint_Constrain_NoChangeWhenInRange(t *testing.T) {
	p := NormPoint{X: 0.5, Y: 0.7}
	p.Constrain()
	assert.Equal(t, p, NormPoint{X: 0.5, Y: 0.7})
}

func TestDist(t *testing.T) {
	a := NormPoint{X: 0, Y: 0}
	b := NormPoint{X: 3, Y: 4}
	assert.Equal(t, Dist(a, b), 5.0)
}
