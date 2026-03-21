package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/stretchr/testify/assert"
	"testing"
)

// mockFunction is a minimal fitness.Function for testing.
type mockFunction struct {
	id int
}

func (m mockFunction) Calculate(_ fitness.PointsData) float64 {
	return float64(m.id)
}

func TestNewOne_ReturnsSameFunctionForAllIndices(t *testing.T) {
	f := mockFunction{id: 42}
	eval := NewOne(f)

	// Get should always return the same function regardless of index
	for i := 0; i < 5; i++ {
		result := eval.Get(i)
		assert.Equal(t, 42.0, result.Calculate(fitness.PointsData{
			Points: normgeom.NormPointGroup{{0.5, 0.5}},
		}))
	}
}

func TestNewOne_NoOpMethods(t *testing.T) {
	f := mockFunction{id: 1}
	eval := NewOne(f)

	// These should not panic
	eval.Prepare()
	eval.Update(0)
	eval.SetBase(0, 1)
	eval.Swap(0, 1)
}

func TestNewMany_ReturnsDistinctFunctions(t *testing.T) {
	callCount := 0
	eval := NewMany(func() fitness.Function {
		callCount++
		return mockFunction{id: callCount}
	}, 3)

	assert.Equal(t, 3, callCount)

	// Each index should return a different function
	r0 := eval.Get(0).Calculate(fitness.PointsData{})
	r1 := eval.Get(1).Calculate(fitness.PointsData{})
	r2 := eval.Get(2).Calculate(fitness.PointsData{})

	assert.Equal(t, 1.0, r0)
	assert.Equal(t, 2.0, r1)
	assert.Equal(t, 3.0, r2)
}

func TestNewMany_NoOpMethods(t *testing.T) {
	eval := NewMany(func() fitness.Function {
		return mockFunction{id: 0}
	}, 2)

	// These should not panic
	eval.Prepare()
	eval.Update(0)
	eval.SetBase(0, 1)
	eval.Swap(0, 1)
}
