package image

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewData(t *testing.T) {
	data := NewData(100, 50)
	w, h := data.Size()
	assert.Equal(t, w, 100)
	assert.Equal(t, h, 50)
}
