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
	MaxOffset func(ln.Linear, *rng.Range) (int, float64)
	Pln       *ln.Polyline
	CtxDB     *rtree.RTree
	SegsDB    *rtree.RTree
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point, constraints []geom.Geometry, options *opts.Opts,
	maximum_offset func(ln.Linear, *rng.Range) (int, float64)) *ConstDP {
	cdp := &ConstDP{
		Simple:    []*HullNode{},
		Opts:      options,
		Hulls:     deque.NewDeque(),
		Ints:      sset.NewSSet(geom.PointCmp),
		MaxOffset: maximum_offset,
		Pln:       ln.NewPolyline(coordinates),
		CtxDB:     rtree.NewRTree(8),
		SegsDB:    rtree.NewRTree(8),
	}
	return cdp.build_segs_db().build_context_db(constraints)
}

func (self *ConstDP) Coordinates() []*geom.Point {
	return self.Pln.Coords
}

func (self *ConstDP) Polyline() *ln.Polyline {
	return self.Pln
}

func (self *ConstDP) MaximumOffset(pln ln.Linear, rg *rng.Range) (int, float64) {
	return self.MaxOffset(pln, rg)
}

//creates constraint db from geometries
func (self *ConstDP) build_context_db(geoms []geom.Geometry) *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		cg := ctx.NewCtxGeom(g, 0, -1).AsContextNeighbour()
		lst = append(lst, cg)
	}
	self.CtxDB.Clear().Load(lst)
	return self
}

//creates constraint db from geometries
func (self *ConstDP) build_segs_db() *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, s := range self.Pln.Segments() {
		lst = append(lst, ctx.NewCtxGeom(s, s.I, s.J).AsSelfSegment())
	}
	self.SegsDB.Clear().Load(lst)
	return self
}
