package mutation

import (
	"Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGaussianMethod_Mutate(t *testing.T) {
	points := normgeom.NormPointGroup{
		{0.23, 0.12},
		{0.56, 0.34},
		{0.34, 0.12},
	}
	otherPoints := points.Copy()

	a := NewGaussianMethod(0, 0)
	c := 0

	a.Mutate(otherPoints, func(mutation Mutation) {
		c++
	})

	assert.Equal(t, otherPoints, points)
	assert.Equal(t, c, 0)

	otherPoints = points.Copy()

	b := NewGaussianMethod(1, 1)
	c = 0

	b.Mutate(otherPoints, func(mutation Mutation) {
		c++
	})

	assert.NotEqual(t, otherPoints, points)
	assert.Equal(t, c, 3)
}

func TestRandomMethod_Mutate(t *testing.T) {
	points := normgeom.NormPointGroup{
		{0.23, 0.12},
		{0.56, 0.34},
		{0.34, 0.12},
	}
	otherPoints := points.Copy()

	a := NewRandomMethod(0, 0)
	c := 0

	a.Mutate(otherPoints, func(mutation Mutation) {
		c++
	})

	assert.Equal(t, otherPoints, points)
	assert.Equal(t, c, 0)

	otherPoints = points.Copy()

	b := NewRandomMethod(1, 1)
	c = 0

	b.Mutate(otherPoints, func(mutation Mutation) {
		c++
	})

	assert.NotEqual(t, otherPoints, points)
	assert.Equal(t, c, 3)
}