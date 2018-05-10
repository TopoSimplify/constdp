package constdp

import (
	"github.com/intdxdt/geom"
	"github.com/intdxdt/rtree"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/lnr"
	"github.com/TopoSimplify/pln"
	"github.com/TopoSimplify/opts"
)

//Type DP
type ConstDP struct {
	*dp.DouglasPeucker
	ContextDB *rtree.RTree
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(
	coordinates []*geom.Point,
	constraints []geom.Geometry,
	options *opts.Opts,
	offsetScore lnr.ScoreFn,
) *ConstDP {
	var instance = (&ConstDP{
		DouglasPeucker: dp.New(coordinates, options, offsetScore),
		ContextDB:      rtree.NewRTree(rtreeBucketSize),
	}).BuildContextDB(constraints) //prepare databases

	if len(coordinates) > 1 {
		instance.Pln = pln.New(coordinates)
	}
	return instance
}

//creates constraint db from geometries
func (self *ConstDP) BuildContextDB(geoms []geom.Geometry) *ConstDP {
	var lst = make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		lst = append(lst, ctx.New(g, 0, -1).AsContextNeighbour())
	}
	self.ContextDB.Clear().Load(lst)
	return self
}
