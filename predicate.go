package constdp

import (
	"simplex/struct/rtree"
	"simplex/util/math"
)

//hull_predicate within range i and j
func hull_predicate(queryhull *HullNode, mindist float64) func(*rtree.KObj) (bool, bool) {
	return func(candidate *rtree.KObj) (bool, bool) {
		dist := candidate.Score()
		candhull := candidate.GetItem().(*HullNode)

		qgeom := queryhull.Geom
		cgeom := candhull.Geom

		// same hull
		if candhull.Range.Equals(queryhull.Range) {
			return false, false
		}

		// if intersects or distance from linegen.geom.context neighbours is almost  offset
		if qgeom.Intersects(cgeom) || math.FloatEqual(dist, mindist) {
			return true, false
		}
		return false, true
	}
}
