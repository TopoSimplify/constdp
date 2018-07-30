package constdp

import (
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/lnr"
	"github.com/TopoSimplify/pln"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
)

//Type DP
type ConstDP struct {
	*dp.DouglasPeucker
	ContextDB *hdb.Hdb
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(
	coordinates []geom.Point,
	constraints []geom.Geometry,
	options *opts.Opts,
	offsetScore lnr.ScoreFn,
) *ConstDP {
	var instance = &ConstDP{
		DouglasPeucker: dp.New(coordinates, options, offsetScore),
		ContextDB:      hdb.NewHdb(rtreeBucketSize),
	}
	instance.BuildContextDB(constraints) //prepare databases

	if len(coordinates) > 1 {
		instance.Pln = pln.New(coordinates)
	}
	return instance
}

//creates constraint db from geometries
func (self *ConstDP) BuildContextDB(geoms []geom.Geometry) *ConstDP {
	var lst = make([]node.Node, 0, len(geoms))
	for i := range geoms {
		cg := ctx.New(geoms[i], 0, -1).AsContextNeighbour()
		lst = append(lst, node.Node{
			MBR:      cg.Bounds(),
			Geom:     cg,
			Instance: self,
		})
	}
	self.ContextDB.Clear().Load(lst)
	return self
}
