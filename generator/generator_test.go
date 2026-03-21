package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomGenerator_Generate(t *testing.T) {
	gen := RandomGenerator{}
	points := gen.Generate(121)
	assert.Equal(t, len(points), 121)
}

func TestSpacedGenerator_Generate(t *testing.T) {
	gen := NewSpacedGenerator(1)
	points := gen.Generate(121)
	assert.Equal(t, len(points), 121)
}

func TestRandomGenerator_PointsInRange(t *testing.T) {
	gen := RandomGenerator{}
	points := gen.Generate(100)
	for _, p := range points {
		assert.GreaterOrEqual(t, p.X, 0.0)
		assert.LessOrEqual(t, p.X, 1.0)
		assert.GreaterOrEqual(t, p.Y, 0.0)
		assert.LessOrEqual(t, p.Y, 1.0)
	}
}

func TestSpacedGenerator_PointsInRange(t *testing.T) {
	gen := NewSpacedGenerator(1)
	points := gen.Generate(50)
	for _, p := range points {
		assert.GreaterOrEqual(t, p.X, 0.0)
		assert.LessOrEqual(t, p.X, 1.0)
		assert.GreaterOrEqual(t, p.Y, 0.0)
		assert.LessOrEqual(t, p.Y, 1.0)
	}
}

func TestRandomGenerator_GenerateSmall(t *testing.T) {
	gen := RandomGenerator{}
	points := gen.Generate(1)
	assert.Equal(t, 1, len(points))
}

func TestSpacedGenerator_GenerateSmall(t *testing.T) {
	gen := NewSpacedGenerator(1)
	points := gen.Generate(1)
	assert.Equal(t, 1, len(points))
}
