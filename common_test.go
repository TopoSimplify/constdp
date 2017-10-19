package constdp

import (
	"github.com/intdxdt/geom"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
)

func linear_coords(wkt string) []*geom.Point{
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func create_hulls(indxs [][]int, coords []*geom.Point) []*HullNode {
	pln := ln.NewPolyline(coords)
	hulls := make([]*HullNode, 0)
	for _, o := range indxs {
		hulls = append(hulls, NewHullNode(pln, rng.NewRange(o[0], o[1])))
	}
	return hulls
}
