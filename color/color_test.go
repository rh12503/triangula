package color

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAverageRGB_Average(t *testing.T) {
	rgb := AverageRGB{}

	rgb.Add(RGB{1, 0, 1})
	rgb.Add(RGB{0.5, 0.5, 0.5})
	rgb.Add(RGB{0.5, 1, 0.5})
	rgb.Add(RGB{1, 0.5, 0})

	assert.Equal(t, rgb.Average(), RGB{
		R: 0.75,
		G: 0.5,
		B: 0.5,
	})
}

func TestAverageRGB_SingleColor(t *testing.T) {
	rgb := AverageRGB{}
	rgb.Add(RGB{0.3, 0.6, 0.9})
	avg := rgb.Average()
	assert.Equal(t, RGB{0.3, 0.6, 0.9}, avg)
}

func TestAverageRGB_MaxValues(t *testing.T) {
	rgb := AverageRGB{}
	rgb.Add(RGB{1, 1, 1})
	rgb.Add(RGB{1, 1, 1})
	avg := rgb.Average()
	assert.Equal(t, RGB{1, 1, 1}, avg)
}

func TestAverageRGB_Count(t *testing.T) {
	rgb := AverageRGB{}
	assert.Equal(t, uint(0), rgb.Count())
	rgb.Add(RGB{0.5, 0.5, 0.5})
	assert.Equal(t, uint(1), rgb.Count())
	rgb.Add(RGB{0.5, 0.5, 0.5})
	assert.Equal(t, uint(2), rgb.Count())
}

func TestNewRGB(t *testing.T) {
	c := NewRGB(0.1, 0.2, 0.3)
	assert.Equal(t, RGB{0.1, 0.2, 0.3}, c)
}
