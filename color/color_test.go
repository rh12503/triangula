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
