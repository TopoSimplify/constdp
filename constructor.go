package constdp

import (
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/geometry"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/lnr"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/pln"
	"github.com/intdxdt/geom"
)

//ConstDP Type
type ConstDP struct {
	*dp.DouglasPeucker
	ContextDB *hdb.Hdb
}

//NewConstDP creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(
	id int,
	coordinates geom.Coords,
	constraints []geometry.IGeometries,
	options *opts.Opts,
	offsetScore lnr.ScoreFn,
	squareOffsetScore ...lnr.ScoreFn,
) *ConstDP {
	var sqrScore lnr.ScoreFn
	if len(squareOffsetScore) > 0 {
		sqrScore = squareOffsetScore[0]
	}

	var instance = ConstDP{
		DouglasPeucker: dp.New(id, coordinates, options, offsetScore, sqrScore),
		ContextDB:      hdb.NewHdb(),
	}
	instance.BuildContextDB(constraints)

	if coordinates.Len() > 1 {
		instance.Pln = pln.CreatePolyline(coordinates)
	}
	return &instance
}

//BuildContextDB creates constraint db from geometries
func (self *ConstDP) BuildContextDB(geoms []geometry.IGeometries) *ConstDP {
	var lst = make([]node.Node, 0, len(geoms))
	for i := range geoms {
		cg := ctx.New(geoms[i].Geometry(), 0, -1).AsContextNeighbour()
		lst = append(lst, node.Node{
			MBR:      cg.Bounds(),
			Geom:     cg,
			Instance: self,
		})
	}
	self.ContextDB.Clear().Load(lst)
	return self
}
