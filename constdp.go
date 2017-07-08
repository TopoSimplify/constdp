package constdp

import (
	"simplex/geom"
	"simplex/struct/sset"
	"simplex/struct/queue"
	"simplex/struct/rtree"
)

//Type DP
type ConstDP struct {
	Simple *sset.SSet
	Opts   *Opts
	Hulls  *queue.Queue
	Ints   *sset.SSet
	//MaximumOffset  = maximum_offset
	Pln    *Polyline
	CtxDB  *rtree.RTree
	SegsDB *rtree.RTree
}

//Creates a new constrained DP Simplification instance
func NewConstDP(options *Opts, build bool) *ConstDP {
	return &ConstDP{}
}

func (cdp *ConstDP) Coordinates() []*geom.Point {
	return cdp.Pln.coords
}

//class Homotopic(object):
//	"""
//	Homotopic decomposition of linear geometries
//	"""
//
//	def __init__(cdp, coordinates, constraints, opts, maximum_offset):
//		cdp.simple 		= []
//		cdp.opts 			= opts
//		cdp.hulls 			= deque()
//		cdp.ints 			= SSet()
//		cdp.maximum_offset = maximum_offset
//		cdp.pln 			= Polyline(coordinates)
//		cdp.ctxdb 			= RTree(maxEntries=8, attribute=('minx', 'miny', 'maxx', 'maxy'))
//		cdp.segsdb 		= RTree(maxEntries=8, attribute=('minx', 'miny', 'maxx', 'maxy'))
//		cdp.build_segs_db() \
//			.build_context_db(geoms=constraints)
