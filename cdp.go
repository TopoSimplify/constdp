package constdp

import (
	"simplex/constdp/ln"
	"simplex/constdp/rng"
	"simplex/constdp/opts"
	"simplex/constdp/ctx"
	"simplex/constdp/cmp"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/deque"
	"github.com/intdxdt/random"
)

//Type DP
type ConstDP struct {
	Id        string
	Simple    *sset.SSet
	Opts      *opts.Opts
	Hulls     *deque.Deque
	Pln       *ln.Polyline
	ContextDB *rtree.RTree
	Meta      map[string]interface{}
	score     func(ln.Linear, *rng.Range) (int, float64)
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point,
	constraints []geom.Geometry, options *opts.Opts,
	offset_score func(ln.Linear, *rng.Range) (int, float64),
) *ConstDP {

	self := &ConstDP{
		Id:     random.String(10),
		Simple: sset.NewSSet(cmp.IntCmp),
		Opts:   options,
		Hulls:  deque.NewDeque(),
		Pln:    ln.NewPolyline(coordinates),

		ContextDB: rtree.NewRTree(8),
		Meta:      make(map[string]interface{}, 0),

		score: offset_score,
	}
	//prepare databases
	return self.build_context_db(constraints)
}

func (self *ConstDP) Coordinates() []*geom.Point {
	return self.Pln.Coords
}

func (self *ConstDP) Polyline() *ln.Polyline {
	return self.Pln
}

func (self *ConstDP) Score(pln ln.Linear, rg *rng.Range) (int, float64) {
	return self.score(pln, rg)
}

//creates constraint db from geometries
func (self *ConstDP) build_context_db(geoms []geom.Geometry) *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		lst = append(lst, ctx.NewCtxGeom(g, 0, -1).AsContextNeighbour())
	}
	self.ContextDB.Clear().Load(lst)
	return self
}
