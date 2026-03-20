package evaluator

import (
	"github.com/RH12503/Triangula/fitness"
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
	fn := mockFunction{id: 1}
	e := NewOne(fn)

	assert.Equal(t, e.Get(0), fn)
	assert.Equal(t, e.Get(1), fn)
	assert.Equal(t, e.Get(99), fn)
}

func TestNewMany_ReturnsDistinctFunctionPerIndex(t *testing.T) {
	counter := 0
	e := NewMany(func() fitness.Function {
		fn := mockFunction{id: counter}
		counter++
		return fn
	}, 3)

	assert.Equal(t, e.Get(0), mockFunction{id: 0})
	assert.Equal(t, e.Get(1), mockFunction{id: 1})
	assert.Equal(t, e.Get(2), mockFunction{id: 2})
}

func TestNewMany_CreatesCorrectCount(t *testing.T) {
	count := 0
	NewMany(func() fitness.Function {
		count++
		return mockFunction{}
	}, 5)
	assert.Equal(t, count, 5)
}

func TestOne_NoOpMethodsDoNotPanic(t *testing.T) {
	e := NewOne(mockFunction{})
	e.SetBase(0, 1)
	e.Swap(0, 1)
	e.Prepare()
	e.Update(0)
}

func TestMany_NoOpMethodsDoNotPanic(t *testing.T) {
	e := NewMany(func() fitness.Function { return mockFunction{} }, 3)
	e.SetBase(0, 1)
	e.Swap(0, 1)
	e.Prepare()
	e.Update(0)
}
