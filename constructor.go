package constdp

//type Homotopic struct {
//	Simple 		    *sset.SSet
//	Opts 			*Opts
//	Hulls 			*queue.Queue
//	Ints 			*sset.SSet
//	//MaximumOffset  = maximum_offset
//	Pln 			*geom.Polyline
//	CtxDB 			*rtree.RTree
//	SegsDB 		    *rtree.RTree
//}

//class Homotopic(object):
//	"""
//	Homotopic decomposition of linear geometries
//	"""
//
//	def __init__(self, coordinates, constraints, opts, maximum_offset):
//		self.simple 		= []
//		self.opts 			= opts
//		self.hulls 			= deque()
//		self.ints 			= SSet()
//		self.maximum_offset = maximum_offset
//		self.pln 			= Polyline(coordinates)
//		self.ctxdb 			= RTree(maxEntries=8, attribute=('minx', 'miny', 'maxx', 'maxy'))
//		self.segsdb 		= RTree(maxEntries=8, attribute=('minx', 'miny', 'maxx', 'maxy'))
//		self.build_segs_db() \
//			.build_context_db(geoms=constraints)
