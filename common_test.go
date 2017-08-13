package constdp

import (
	"simplex/geom"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
)

func linear_coords(wkt string) []*geom.Point{
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func create_hulls(indxs [][]int, coords []*geom.Point) []*HullNode {
	n := len(coords) - 1
	pln := ln.NewPolyline(coords)
	hulls := make([]*HullNode, 0)
	pr := rng.NewRange(0, n)

	for _, o := range indxs {
		hulls = append(hulls, NewHullNode(pln, rng.NewRange(o[0], o[1]), pr))
	}
	return hulls
}
