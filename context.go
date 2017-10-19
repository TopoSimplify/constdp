package constdp

import (
	"github.com/intdxdt/rtree"
	"simplex/constdp/igeom"
)

//find context neighbours
func find_context_neighbs(database *rtree.RTree, query igeom.IGeom, dist float64) []rtree.BoxObj {
	return find_knn(database, query.Geometry(), dist, score_fn(query))
}

//find context hulls
func find_context_hulls(hulldb *rtree.RTree, hull *HullNode, dist float64) []rtree.BoxObj {
	predicate := hull_predicate(hull, dist)
	return find_knn(hulldb, hull.Geometry(), dist, score_fn(hull), predicate)
}

//hull predicate within index range i, j.
func hull_predicate(queryhull *HullNode, dist float64) func(*rtree.KObj) (bool, bool) {
	//@formatter:off
	return func(candidate *rtree.KObj) (bool, bool) {
		candhull := candidate.GetItem().(*HullNode)

		qgeom    := queryhull.Geom
		cgeom    := candhull.Geom

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
