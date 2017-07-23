package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/rtree"
	"simplex/struct/deque"
)
//Opts
type Opts struct {
	Threshold              float64
	MinDist                float64
	RelaxDist              float64
	KeepSelfIntersects     bool
	AvoidNewSelfIntersects bool
	GeomRelation           bool
	DistRelation           bool
	DirRelation            bool
}

//Linear interface
type Linear interface {
	Coordinates() []*geom.Point
}

//Type DP
type ConstDP struct {
	Simple        []*HullNode
	Opts          *Opts
	Hulls         *deque.Deque
	Ints          *sset.SSet
	MaximumOffset func(Linear, *Range) (int, float64)
	Pln           *Polyline
	CtxDB         *rtree.RTree
	SegsDB        *rtree.RTree
}

//Creates a new constrained DP Simplification instance
//	dp decomposition of linear geometries
func NewConstDP(coordinates []*geom.Point, constraints []geom.Geometry, opts *Opts,
	maximum_offset func(Linear, *Range) (int, float64)) *ConstDP {
	cdp := &ConstDP{
		Simple:        []*HullNode{},
		Opts:          opts,
		Hulls:         deque.NewDeque(),
		Ints:          sset.NewSSet(geom.PointCmp),
		MaximumOffset: maximum_offset,
		Pln:           NewPolyline(coordinates),
		CtxDB:         rtree.NewRTree(8),
		SegsDB:        rtree.NewRTree(8),
	}
	return cdp.build_segs_db().build_context_db(constraints)
}

func (self *ConstDP) Coordinates() []*geom.Point {
	return self.Pln.Coords
}

//creates constraint db from geometries
func (self *ConstDP) build_context_db(geoms []geom.Geometry) *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		cg := NewCtxGeom(g, 0, -1).AsContextNeighbour()
		lst = append(lst, cg)
	}
	self.CtxDB.Clear().Load(lst)
	return self
}

//creates constraint db from geometries
func (self *ConstDP) build_segs_db() *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, s := range self.Pln.Segments() {
		lst = append(lst, NewCtxGeom(s, s.I, s.J).AsSelfSegment())
	}
	self.SegsDB.Clear().Load(lst)
	return self
}
