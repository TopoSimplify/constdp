package constdp

import (
	"simplex/struct/rtree"
)

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
