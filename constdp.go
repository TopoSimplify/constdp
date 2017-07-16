package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/queue"
	"simplex/struct/rtree"
)

//Type DP
type ConstDP struct {
	Simple        *sset.SSet
	Opts          *Opts
	Hulls         *queue.Queue
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
		Simple:        sset.NewSSet(geom.PointCmp),
		Opts:          opts,
		Hulls:         queue.NewQueue(),
		Ints:          sset.NewSSet(geom.PointCmp),
		MaximumOffset: maximum_offset,
		Pln:           NewPolyline(coordinates),
		CtxDB:         rtree.NewRTree(8),
		SegsDB:        rtree.NewRTree(8),
	}
	return cdp.build_segs_db().build_context_db(constraints)
}

func (cdp *ConstDP) Coordinates() []*geom.Point {
	return cdp.Pln.coords
}

//creates constraint db from geometries
func (cdp *ConstDP) build_context_db(geoms []geom.Geometry) *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, g := range geoms {
		cg := NewCtxGeom(g, 0, -1).AsContextNeighbour()
		lst = append(lst, cg)
	}
	cdp.CtxDB.Clear().Load(lst)
	return cdp
}

//creates constraint db from geometries
func (cdp *ConstDP) build_segs_db() *ConstDP {
	lst := make([]rtree.BoxObj, 0)
	for _, s := range cdp.Pln.Segments() {
		lst = append(lst, NewCtxGeom(s, s.I, s.J).AsSelfSegment())
	}
	cdp.SegsDB.Clear().Load(lst)
	return cdp
}
