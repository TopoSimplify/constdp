package constdp

import (
	"simplex/pln"
	"simplex/rng"
	"simplex/opts"
	"simplex/lnr"
	"simplex/ctx"
	"github.com/intdxdt/cmp"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/sset"
	"github.com/intdxdt/deque"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/random"
)

//Type DP
type ConstDP struct {
	Id        string
	Opts      *opts.Opts
	Hulls     *deque.Deque
	Pln       *pln.Polyline
	ContextDB *rtree.RTree
	Meta      map[string]interface{}

	simple *sset.SSet
	score  func(lnr.Linear, *rng.Range) (int, float64)
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point, constraints []geom.Geometry,
	options *opts.Opts, offset_score lnr.ScoreFn) *ConstDP {
	return (&ConstDP{
		Id:    random.String(10),
		Opts:  options,
		Hulls: deque.NewDeque(),
		Pln:   pln.New(coordinates),

		ContextDB: rtree.NewRTree(16),
		Meta:      make(map[string]interface{}, 0),

		simple: sset.NewSSet(cmp.Int),
		score:  offset_score,
	}).build_context_db(constraints) //prepare databases
}

func (self *ConstDP) Simple() *sset.SSet {
	return self.simple
}

func (self *ConstDP) Coordinates() []*geom.Point {
	return self.Pln.Coordinates
}

func (self *ConstDP) Polyline() *pln.Polyline {
	return self.Pln
}

func (self *ConstDP) Score(pln lnr.Linear, rg *rng.Range) (int, float64) {
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
