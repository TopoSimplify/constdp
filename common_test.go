package constdp

import (
	"simplex/dp"
	"simplex/pln"
	"simplex/rng"
	"simplex/node"
	"github.com/intdxdt/geom"
)

func linear_coords(wkt string) []*geom.Point{
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func create_hulls(indxs [][]int, coords []*geom.Point) []*node.Node {
	poly := pln.New(coords)
	hulls := make([]*node.Node, 0)
	for _, o := range indxs {
		hulls = append(hulls, node.New(poly, rng.NewRange(o[0], o[1]), dp.NodeGeometry))
	}
	return hulls
}
