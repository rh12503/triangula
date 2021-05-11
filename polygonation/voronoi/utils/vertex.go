// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Utils for processing voronoi diagrams

package utils

import (
	"github.com/RH12503/Triangula/polygonation/voronoi"
	"math"
)

func Distance(a,b voronoi.Vertex) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx+dy*dy)
}