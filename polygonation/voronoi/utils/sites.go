// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Utils for processing voronoi diagrams

package utils

import (
	"github.com/RH12503/Triangula/polygonation/voronoi"
	"math/rand"
)

// Generate random sites in given bounding box
func RandomSites(bbox voronoi.BBox, count int) []voronoi.Vertex {
	sites := make([]voronoi.Vertex, count)
	w := bbox.Xr - bbox.Xl
	h := bbox.Yb - bbox.Yt
	for j := 0; j < count; j++ {
		sites[j].X = rand.Float64() * w + bbox.Xl
		sites[j].Y = rand.Float64() * h + bbox.Yt
	}
	return sites
}