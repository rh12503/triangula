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

