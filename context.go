package constdp

import (
	"simplex/igeom"
	"simplex/node"
	"github.com/intdxdt/rtree"
)

//find context neighbours
func find_context_neighbs(database *rtree.RTree, query igeom.IGeom, dist float64) []rtree.BoxObj {
	return find_knn(database, query.Geometry(), dist, score_fn(query))
}

//find context hulls
func find_context_hulls(hulldb *rtree.RTree, hull *node.Node, dist float64) []rtree.BoxObj {
	return find_knn(hulldb, hull.Geometry(), dist, score_fn(hull), hull_predicate(hull, dist))
}

//hull predicate within index range i, j.
func hull_predicate(queryhull *node.Node, dist float64) func(*rtree.KObj) (bool, bool) {
	//@formatter:off
	return func(candidate *rtree.KObj) (bool, bool) {
		var candhull = candidate.GetItem().(*node.Node)
		var qgeom    = queryhull.Geom
		var cgeom    = candhull.Geom

		// same hull
		if candhull.Range.Equals(queryhull.Range) {
			return false, false
		}

		// if intersects or distance from context neighbours is within dist
		if qgeom.Intersects(cgeom) || (candidate.Score() <= dist) {
			return true, false
		}
		return false, true
	}
}
