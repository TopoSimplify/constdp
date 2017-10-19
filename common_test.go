package constdp

import (
	"github.com/intdxdt/geom"
	"simplex/pln"
	"simplex/rng"
)

func linear_coords(wkt string) []*geom.Point{
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func create_hulls(indxs [][]int, coords []*geom.Point) []*HullNode {
	poly := pln.New(coords)
	hulls := make([]*HullNode, 0)
	for _, o := range indxs {
		hulls = append(hulls, NewHullNode(poly, rng.NewRange(o[0], o[1])))
	}
	return hulls
}
