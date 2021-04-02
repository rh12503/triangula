package random

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestIntn(t *testing.T) {
	sum := 0
	for i := 0; i < 10000; i++ {
		sum += Intn(10)
	}
	average := float64(sum) / 10000.
	assert.True(t, math.Abs(average)-4.5 < 0.1)
}
