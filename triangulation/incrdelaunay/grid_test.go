package incrdelaunay

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircumcircleGrid_Allocation_NonSquare(t *testing.T) {
	testCases := []struct {
		name string
		cols int
		rows int
		w    int
		h    int
	}{
		{
			name: "more columns than rows",
			cols: 25,
			rows: 6,
			w:    500,
			h:    120,
		},
		{
			name: "more rows than columns",
			cols: 7,
			rows: 24,
			w:    140,
			h:    480,
		},
		{
			name: "single column multi row",
			cols: 1,
			rows: 15,
			w:    50,
			h:    300,
		},
		{
			name: "single row multi column",
			cols: 15,
			rows: 1,
			w:    300,
			h:    50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grid := NewCircumcircleGrid(tc.cols, tc.rows, tc.w, tc.h)
			require.Equal(t, tc.cols, len(grid.triangles), "outer dimension must equal cols")
			for x := 0; x < tc.cols; x++ {
				assert.Equal(t, tc.rows, len(grid.triangles[x]), "inner dimension must equal rows")
			}
		})
	}
}

func TestCircumcircleGrid_ColsGreaterThanRows_Operations(t *testing.T) {
	cols := 20
	rows := 5
	w := 200
	h := 50
	grid := NewCircumcircleGrid(cols, rows, w, h)

	// Triangle in column index ~18 (18 >= rows, would panic if outer dimension was rows)
	tri1 := NewTriangle(Point{180, 20}, Point{190, 20}, Point{185, 30})
	// Triangle near origin in column index ~1
	tri2 := NewTriangle(Point{10, 10}, Point{20, 10}, Point{15, 20})
	triangles := []Triangle{tri1, tri2}

	grid.AddTriangle(tri1, 0)
	grid.AddTriangle(tri2, 1)

	// Verify HasPoint for points inside and outside
	assert.True(t, grid.HasPoint(Point{180, 20}, triangles))
	assert.True(t, grid.HasPoint(Point{10, 10}, triangles))
	assert.False(t, grid.HasPoint(Point{100, 25}, triangles))

	// Verify IterCircumcirclesThatContain
	var foundTriangles []uint16
	grid.IterCircumcirclesThatContain(Point{185, 23}, triangles, func(i uint16) {
		foundTriangles = append(foundTriangles, i)
	})
	assert.Contains(t, foundTriangles, uint16(0))

	// Verify IterThatHasVertex
	var vertexTriangles []uint16
	grid.IterThatHasVertex(Point{180, 20}, triangles, func(i uint16) {
		vertexTriangles = append(vertexTriangles, i)
	})
	assert.Contains(t, vertexTriangles, uint16(0))

	// Verify RemoveCircumcirclesThatContain
	var removedCircumcircles []uint16
	grid.RemoveCircumcirclesThatContain(Point{185, 23}, triangles, func(i uint16) {
		removedCircumcircles = append(removedCircumcircles, i)
	})
	assert.Contains(t, removedCircumcircles, uint16(0))
	assert.False(t, grid.HasPoint(Point{180, 20}, triangles))

	// Verify RemoveThatHasVertex on the other triangle
	var removedVertices []uint16
	grid.RemoveThatHasVertex(Point{10, 10}, triangles, func(i uint16) {
		removedVertices = append(removedVertices, i)
	})
	assert.Contains(t, removedVertices, uint16(1))
	assert.False(t, grid.HasPoint(Point{10, 10}, triangles))
}

func TestCircumcircleGrid_RowsGreaterThanCols_Operations(t *testing.T) {
	cols := 5
	rows := 20
	w := 50
	h := 200
	grid := NewCircumcircleGrid(cols, rows, w, h)

	// Triangle in row index ~18 (18 >= cols, would panic if inner dimension was cols)
	tri1 := NewTriangle(Point{20, 180}, Point{30, 180}, Point{25, 190})
	tri2 := NewTriangle(Point{10, 10}, Point{20, 10}, Point{15, 20})
	triangles := []Triangle{tri1, tri2}

	grid.AddTriangle(tri1, 0)
	grid.AddTriangle(tri2, 1)

	assert.True(t, grid.HasPoint(Point{20, 180}, triangles))
	assert.True(t, grid.HasPoint(Point{10, 10}, triangles))
	assert.False(t, grid.HasPoint(Point{25, 100}, triangles))

	// Test RemoveTriangle directly
	grid.RemoveTriangle(tri1, 0)
	assert.False(t, grid.HasPoint(Point{20, 180}, triangles))
	assert.True(t, grid.HasPoint(Point{10, 10}, triangles))
}

func TestCircumcircleGrid_Set_NonSquare(t *testing.T) {
	cols := 16
	rows := 8
	w := 160
	h := 80

	grid1 := NewCircumcircleGrid(cols, rows, w, h)
	grid2 := NewCircumcircleGrid(cols, rows, w, h)

	tri := NewTriangle(Point{140, 60}, Point{150, 60}, Point{145, 70})
	triangles := []Triangle{tri}

	grid1.AddTriangle(tri, 0)
	assert.True(t, grid1.HasPoint(Point{140, 60}, triangles))
	assert.False(t, grid2.HasPoint(Point{140, 60}, triangles))

	grid2.Set(&grid1)
	assert.True(t, grid2.HasPoint(Point{140, 60}, triangles))
}

func TestCircumcircleGrid_BoundaryPoints_NonSquare(t *testing.T) {
	cols := 12
	rows := 4
	w := 120
	h := 40
	grid := NewCircumcircleGrid(cols, rows, w, h)

	// Triangles with vertices on boundaries: (0, 0), (w, h), and (w-1, h-1)
	triOrigin := NewTriangle(Point{0, 0}, Point{5, 0}, Point{0, 5})
	triCorner := NewTriangle(Point{int16(w - 5), int16(h)}, Point{int16(w), int16(h - 5)}, Point{int16(w), int16(h)})
	triInner := NewTriangle(Point{int16(w - 2), int16(h - 1)}, Point{int16(w - 1), int16(h - 2)}, Point{int16(w - 1), int16(h - 1)})
	triangles := []Triangle{triOrigin, triCorner, triInner}

	grid.AddTriangle(triOrigin, 0)
	grid.AddTriangle(triCorner, 1)
	grid.AddTriangle(triInner, 2)

	assert.True(t, grid.HasPoint(Point{0, 0}, triangles))
	assert.True(t, grid.HasPoint(Point{int16(w), int16(h)}, triangles))
	assert.True(t, grid.HasPoint(Point{int16(w - 1), int16(h - 1)}, triangles))
}
