package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTriangle(t *testing.T) {
	triangle := NewTriangle(14, 15, 92, 65, 35, 89)
	assert.Equal(t, triangle, Triangle{Points: [3]Point{
		{14, 15},
		{92, 65},
		{35, 89},
	}})
}
