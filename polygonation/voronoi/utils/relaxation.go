// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Utils for processing voronoi diagrams

package utils

import (
	"github.com/RH12503/Triangula/polygonation/voronoi"
)


// Apply lloyd relaxation algorithm to the cells.
func LloydRelaxation(cells []*voronoi.Cell) (ret []voronoi.Vertex) {
	ret = make([]voronoi.Vertex, len(cells))
	for id, cell := range cells {
		ret[id] = CellCentroid(cell)
	}
	return
}