package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/rtree"
	"simplex/struct/deque"
	"simplex/constdp/ln"
	"simplex/constdp/rng"
	"simplex/constdp/opts"
	"simplex/constdp/ctx"
)

//Type DP
type ConstDP struct {
	Simple    []*HullNode
	Opts      *opts.Opts
	Hulls     *deque.Deque
	Ints      *sset.SSet
	Pln       *ln.Polyline
	ContextDB *rtree.RTree
	SegmentDB *rtree.RTree
	score     func(ln.Linear, *rng.Range) (int, float64)
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point,
	constraints []geom.Geometry, options *opts.Opts,
	offset_score func(ln.Linear, *rng.Range) (int, float64),
) *ConstDP {

	self := &ConstDP{
		Simple: []*HullNode{},
		Opts:   options,
		Hulls:  deque.NewDeque(),
		Ints:   sset.NewSSet(geom.PointCmp),
		Pln:    ln.NewPolyline(coordinates),

		ContextDB: rtree.NewRTree(8),
		SegmentDB: rtree.NewRTree(8),

		score: offset_score,
	}
	//prepare databases
	return self.build_segs_db().build_context_db(constraints)
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

//creates constraint db from geometries
func (self *ConstDP) build_segs_db() *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, s := range self.Pln.Segments() {
		lst = append(lst, ctx.NewCtxGeom(s, s.I, s.J).AsSelfSegment())
	}
	self.SegmentDB.Clear().Load(lst)
	return self
}
