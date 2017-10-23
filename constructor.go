package constdp

import (
	"simplex/dp"
	"simplex/ctx"
	"simplex/pln"
	"simplex/rng"
	"simplex/opts"
	"simplex/lnr"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/rtree"
)

//Type DP
type ConstDP struct {
	*dp.DouglasPeucker
	ContextDB *rtree.RTree
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point, constraints []geom.Geometry,
	options *opts.Opts, offsetScore lnr.ScoreFn) *ConstDP {
	var instance = (&ConstDP{
		DouglasPeucker: dp.New(coordinates, options, offsetScore),
		ContextDB:      rtree.NewRTree(16),
	}).build_context_db(constraints) //prepare databases

	if len(coordinates) > 1 {
		instance.Pln = pln.New(coordinates)
	}
	return instance
}


func (self *ConstDP) Options() *opts.Opts {
	return self.Opts
}

func (self *ConstDP) Coordinates() []*geom.Point {
	return self.Pln.Coordinates
}

func (self *ConstDP) Polyline() *pln.Polyline {
	return self.Pln
}

func (self *ConstDP) Score(pln lnr.Linear, rg *rng.Range) (int, float64) {
	return self.ScoreFn(pln, rg)
}

//creates constraint db from geometries
func (self *ConstDP) build_context_db(geoms []geom.Geometry) *ConstDP {
	var lst = make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		lst = append(lst, ctx.New(g, 0, -1).AsContextNeighbour())
	}
	self.ContextDB.Clear().Load(lst)
	return self
}
